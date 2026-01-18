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

type ExerciseQuestionRepository struct {
	db *gorm.DB
}

// IExerciseQuestionRepository defines the contract for exercise question data access operations.
type IExerciseQuestionRepository interface {
	// Create inserts a new exercise question entry into the database.
	Create(context.Context, *dto.CreateExerciseQuestionRequest) (*models.ExerciseQuestion, error)

	// GetAll retrieves all exercise questions with optional filtering and pagination.
	// Returns the list of exercise questions and total count.
	GetAll(context.Context, *dto.ExerciseQuestionFilterRequest) ([]models.ExerciseQuestion, int64, error)

	// GetByID retrieves a single exercise question entry by its ID.
	GetByID(context.Context, uint) (*models.ExerciseQuestion, error)

	// GetByExerciseIDAndOrderIndex retrieves an exercise question by exercise ID and order index.
	// Used for duplicate detection.
	GetByExerciseIDAndOrderIndex(context.Context, uint, int) (*models.ExerciseQuestion, error)

	// GetByExerciseID retrieves all exercise questions for a specific exercise, ordered by order_index.
	GetByExerciseID(context.Context, uint) ([]models.ExerciseQuestion, error)

	// Update modifies an existing exercise question entry by ID.
	Update(context.Context, *dto.UpdateExerciseQuestionRequest, uint) (*models.ExerciseQuestion, error)

	// Delete removes an exercise question entry by ID.
	Delete(context.Context, uint) error

	// Publish sets an exercise question as published with the current timestamp.
	Publish(context.Context, uint) (*models.ExerciseQuestion, error)

	// Unpublish sets an exercise question as unpublished.
	Unpublish(context.Context, uint) (*models.ExerciseQuestion, error)
}

func NewExerciseQuestionRepository(db *gorm.DB) IExerciseQuestionRepository {
	return &ExerciseQuestionRepository{db: db}
}

func (r *ExerciseQuestionRepository) Create(ctx context.Context, req *dto.CreateExerciseQuestionRequest) (*models.ExerciseQuestion, error) {
	question := models.ExerciseQuestion{
		ExerciseID:    req.ExerciseID,
		QuestionText:  req.QuestionText,
		QuestionType:  req.QuestionType,
		Options:       req.Options,
		CorrectAnswer: req.CorrectAnswer,
		Explanation:   req.Explanation,
		AudioURL:      req.AudioURL,
		ImageURL:      req.ImageURL,
		OrderIndex:    req.OrderIndex,
		Points:        req.Points,
		IsPublished:   false,
	}

	err := r.db.WithContext(ctx).Create(&question).Error
	if err != nil {
		// Check for unique constraint violation on order_index within exercise
		if strings.Contains(err.Error(), "idx_question_exercise_order") ||
			strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, errConstant.ErrDuplicateQuestionOrderIndex
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Load relationships using Preload for efficiency
	err = r.db.WithContext(ctx).
		Preload("Exercise").
		Preload("Exercise.Lesson").
		Preload("Exercise.Lesson.Course").
		Preload("Exercise.Lesson.Course.JlptLevel").
		First(&question, question.ID).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &question, nil
}

func (r *ExerciseQuestionRepository) GetAll(ctx context.Context, filter *dto.ExerciseQuestionFilterRequest) ([]models.ExerciseQuestion, int64, error) {
	var questions []models.ExerciseQuestion
	var total int64

	// Build base query with filters
	query := r.db.WithContext(ctx).Model(&models.ExerciseQuestion{})

	// Apply filters
	if filter != nil {
		if filter.ExerciseID > 0 {
			query = query.Where("exercise_id = ?", filter.ExerciseID)
		}
		if filter.QuestionType != "" {
			query = query.Where("question_type = ?", filter.QuestionType)
		}
		if filter.IsPublished != nil {
			query = query.Where("is_published = ?", *filter.IsPublished)
		}
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(question_text) LIKE ? OR LOWER(explanation) LIKE ?", searchPattern, searchPattern)
		}
	}

	// Count total records with filters applied
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Apply sorting with whitelist validation (defense in depth)
	allowedSortFields := map[string]string{
		"order_index":   "order_index",
		"question_text": "question_text",
		"created_at":    "created_at",
		"points":        "points",
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

	query = query.Preload("Exercise").
		Preload("Exercise.Lesson").
		Preload("Exercise.Lesson.Course").
		Preload("Exercise.Lesson.Course.JlptLevel").
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

	err := query.Find(&questions).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return questions, total, nil
}

func (r *ExerciseQuestionRepository) GetByID(ctx context.Context, id uint) (*models.ExerciseQuestion, error) {
	var question models.ExerciseQuestion
	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Preload("Exercise.Lesson").
		Preload("Exercise.Lesson.Course").
		Preload("Exercise.Lesson.Course.JlptLevel").
		Where("id = ?", id).
		First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrExerciseQuestionNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &question, nil
}

func (r *ExerciseQuestionRepository) GetByExerciseIDAndOrderIndex(ctx context.Context, exerciseID uint, orderIndex int) (*models.ExerciseQuestion, error) {
	var question models.ExerciseQuestion
	err := r.db.WithContext(ctx).
		Where("exercise_id = ? AND order_index = ?", exerciseID, orderIndex).
		First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrExerciseQuestionNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &question, nil
}

func (r *ExerciseQuestionRepository) GetByExerciseID(ctx context.Context, exerciseID uint) ([]models.ExerciseQuestion, error) {
	var questions []models.ExerciseQuestion
	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Preload("Exercise.Lesson").
		Preload("Exercise.Lesson.Course").
		Preload("Exercise.Lesson.Course.JlptLevel").
		Where("exercise_id = ?", exerciseID).
		Order("order_index ASC").
		Find(&questions).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return questions, nil
}

func (r *ExerciseQuestionRepository) Update(ctx context.Context, req *dto.UpdateExerciseQuestionRequest, id uint) (*models.ExerciseQuestion, error) {
	question := models.ExerciseQuestion{
		ExerciseID:    req.ExerciseID,
		QuestionText:  req.QuestionText,
		QuestionType:  req.QuestionType,
		Options:       req.Options,
		CorrectAnswer: req.CorrectAnswer,
		Explanation:   req.Explanation,
		AudioURL:      req.AudioURL,
		ImageURL:      req.ImageURL,
		OrderIndex:    req.OrderIndex,
		Points:        req.Points,
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&question)

	if result.Error != nil {
		// Check for unique constraint violation on order_index within exercise
		if strings.Contains(result.Error.Error(), "idx_question_exercise_order") ||
			strings.Contains(result.Error.Error(), "duplicate key") ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint") {
			return nil, errConstant.ErrDuplicateQuestionOrderIndex
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrExerciseQuestionNotFound
	}

	// Fetch the updated record with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Preload("Exercise.Lesson").
		Preload("Exercise.Lesson.Course").
		Preload("Exercise.Lesson.Course.JlptLevel").
		Where("id = ?", id).
		First(&question).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &question, nil
}

func (r *ExerciseQuestionRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.ExerciseQuestion{})

	if result.Error != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errConstant.ErrExerciseQuestionNotFound
	}

	return nil
}

func (r *ExerciseQuestionRepository) Publish(ctx context.Context, id uint) (*models.ExerciseQuestion, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.ExerciseQuestion{}).
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
		return nil, errConstant.ErrExerciseQuestionNotFound
	}

	// Fetch the updated record with preloaded relationships
	var question models.ExerciseQuestion
	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Preload("Exercise.Lesson").
		Preload("Exercise.Lesson.Course").
		Preload("Exercise.Lesson.Course.JlptLevel").
		Where("id = ?", id).
		First(&question).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &question, nil
}

func (r *ExerciseQuestionRepository) Unpublish(ctx context.Context, id uint) (*models.ExerciseQuestion, error) {
	result := r.db.WithContext(ctx).
		Model(&models.ExerciseQuestion{}).
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
		return nil, errConstant.ErrExerciseQuestionNotFound
	}

	// Fetch the updated record with preloaded relationships
	var question models.ExerciseQuestion
	err := r.db.WithContext(ctx).
		Preload("Exercise").
		Preload("Exercise.Lesson").
		Preload("Exercise.Lesson.Course").
		Preload("Exercise.Lesson.Course.JlptLevel").
		Where("id = ?", id).
		First(&question).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &question, nil
}
