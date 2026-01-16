package repositories

import (
	"context"
	"errors"
	errWrap "manabu-service/common/error"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ExerciseRepository struct {
	db *gorm.DB
}

// IExerciseRepository defines the contract for exercise data access operations.
type IExerciseRepository interface {
	// Create inserts a new exercise entry into the database.
	Create(context.Context, *dto.CreateExerciseRequest) (*models.Exercise, error)

	// GetAll retrieves all exercises with optional filtering and pagination.
	// Returns the list of exercises and total count.
	GetAll(context.Context, *dto.ExerciseFilterRequest) ([]models.Exercise, int64, error)

	// GetByID retrieves a single exercise entry by its ID.
	GetByID(context.Context, uint) (*models.Exercise, error)

	// GetByLessonIDAndOrderIndex retrieves an exercise entry by lesson ID and order index.
	// Used for duplicate detection.
	GetByLessonIDAndOrderIndex(context.Context, uint, int) (*models.Exercise, error)

	// GetByLessonID retrieves all exercises for a specific lesson, ordered by order_index.
	GetByLessonID(context.Context, uint) ([]models.Exercise, error)

	// Update modifies an existing exercise entry by ID.
	Update(context.Context, *dto.UpdateExerciseRequest, uint) (*models.Exercise, error)

	// Delete removes an exercise entry by ID.
	Delete(context.Context, uint) error

	// Publish sets an exercise as published with the current timestamp.
	Publish(context.Context, uint) (*models.Exercise, error)

	// Unpublish sets an exercise as unpublished.
	Unpublish(context.Context, uint) (*models.Exercise, error)
}

func NewExerciseRepository(db *gorm.DB) IExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) Create(ctx context.Context, req *dto.CreateExerciseRequest) (*models.Exercise, error) {
	exercise := models.Exercise{
		LessonID:         req.LessonID,
		Title:            req.Title,
		Description:      req.Description,
		ExerciseType:     req.ExerciseType,
		OrderIndex:       req.OrderIndex,
		DifficultyLevel:  req.DifficultyLevel,
		EstimatedMinutes: req.EstimatedMinutes,
		IsPublished:      false,
	}

	err := r.db.WithContext(ctx).Create(&exercise).Error
	if err != nil {
		// Check for unique constraint violation on order_index within lesson
		if strings.Contains(err.Error(), "idx_exercise_lesson_order") ||
			strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, errConstant.ErrDuplicateExerciseOrderIndex
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Load relationships using Preload for efficiency
	err = r.db.WithContext(ctx).
		Preload("Lesson").
		Preload("Lesson.Course").
		Preload("Lesson.Course.JlptLevel").
		First(&exercise, exercise.ID).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &exercise, nil
}

func (r *ExerciseRepository) GetAll(ctx context.Context, filter *dto.ExerciseFilterRequest) ([]models.Exercise, int64, error) {
	var exercises []models.Exercise
	var total int64

	// Build base query with filters
	query := r.db.WithContext(ctx).Model(&models.Exercise{})

	// Apply filters
	if filter != nil {
		if filter.LessonID > 0 {
			query = query.Where("lesson_id = ?", filter.LessonID)
		}
		if filter.ExerciseType != "" {
			query = query.Where("exercise_type = ?", filter.ExerciseType)
		}
		if filter.IsPublished != nil {
			query = query.Where("is_published = ?", *filter.IsPublished)
		}
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchPattern, searchPattern)
		}
	}

	// Count total records with filters applied
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Apply sorting with whitelist validation (defense in depth)
	allowedSortFields := map[string]string{
		"order_index":      "order_index",
		"title":            "title",
		"created_at":       "created_at",
		"difficulty_level": "difficulty_level",
	}
	allowedSortOrders := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	sortBy := "order_index"
	sortOrder := "ASC"
	if filter != nil {
		if filter.SortBy != "" {
			if validField, ok := allowedSortFields[filter.SortBy]; ok {
				sortBy = validField
			}
		}
		if filter.SortOrder != "" {
			if validOrder, ok := allowedSortOrders[filter.SortOrder]; ok {
				sortOrder = validOrder
			}
		}
	}

	query = query.Preload("Lesson").
		Preload("Lesson.Course").
		Preload("Lesson.Course.JlptLevel").
		Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filter != nil && filter.Limit > 0 {
		// Defensive validation: ensure page is at least 1
		page := filter.Page
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * filter.Limit
		// Ensure offset is never negative
		if offset < 0 {
			offset = 0
		}

		query = query.Limit(filter.Limit).Offset(offset)
	}

	err := query.Find(&exercises).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return exercises, total, nil
}

func (r *ExerciseRepository) GetByID(ctx context.Context, id uint) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.db.WithContext(ctx).
		Preload("Lesson").
		Preload("Lesson.Course").
		Preload("Lesson.Course.JlptLevel").
		Where("id = ?", id).
		First(&exercise).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrExerciseNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &exercise, nil
}

func (r *ExerciseRepository) GetByLessonIDAndOrderIndex(ctx context.Context, lessonID uint, orderIndex int) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.db.WithContext(ctx).
		Where("lesson_id = ? AND order_index = ?", lessonID, orderIndex).
		First(&exercise).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrExerciseNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &exercise, nil
}

func (r *ExerciseRepository) GetByLessonID(ctx context.Context, lessonID uint) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.WithContext(ctx).
		Preload("Lesson").
		Preload("Lesson.Course").
		Preload("Lesson.Course.JlptLevel").
		Where("lesson_id = ?", lessonID).
		Order("order_index ASC").
		Find(&exercises).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return exercises, nil
}

func (r *ExerciseRepository) Update(ctx context.Context, req *dto.UpdateExerciseRequest, id uint) (*models.Exercise, error) {
	exercise := models.Exercise{
		LessonID:         req.LessonID,
		Title:            req.Title,
		Description:      req.Description,
		ExerciseType:     req.ExerciseType,
		OrderIndex:       req.OrderIndex,
		DifficultyLevel:  req.DifficultyLevel,
		EstimatedMinutes: req.EstimatedMinutes,
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&exercise)

	if result.Error != nil {
		// Check for unique constraint violation on order_index within lesson
		if strings.Contains(result.Error.Error(), "idx_exercise_lesson_order") ||
			strings.Contains(result.Error.Error(), "duplicate key") ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint") {
			return nil, errConstant.ErrDuplicateExerciseOrderIndex
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrExerciseNotFound
	}

	// Fetch the updated record with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("Lesson").
		Preload("Lesson.Course").
		Preload("Lesson.Course.JlptLevel").
		Where("id = ?", id).
		First(&exercise).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &exercise, nil
}

func (r *ExerciseRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Exercise{})

	if result.Error != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errConstant.ErrExerciseNotFound
	}

	return nil
}

func (r *ExerciseRepository) Publish(ctx context.Context, id uint) (*models.Exercise, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.Exercise{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_published": true,
			"published_at": now,
		})

	if result.Error != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrExerciseNotFound
	}

	// Fetch the updated record with preloaded relationships
	var exercise models.Exercise
	err := r.db.WithContext(ctx).
		Preload("Lesson").
		Preload("Lesson.Course").
		Preload("Lesson.Course.JlptLevel").
		Where("id = ?", id).
		First(&exercise).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &exercise, nil
}

func (r *ExerciseRepository) Unpublish(ctx context.Context, id uint) (*models.Exercise, error) {
	result := r.db.WithContext(ctx).
		Model(&models.Exercise{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_published": false,
			"published_at": nil,
		})

	if result.Error != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrExerciseNotFound
	}

	// Fetch the updated record with preloaded relationships
	var exercise models.Exercise
	err := r.db.WithContext(ctx).
		Preload("Lesson").
		Preload("Lesson.Course").
		Preload("Lesson.Course.JlptLevel").
		Where("id = ?", id).
		First(&exercise).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &exercise, nil
}
