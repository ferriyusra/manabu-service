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

type LessonRepository struct {
	db *gorm.DB
}

// ILessonRepository defines the contract for lesson data access operations.
type ILessonRepository interface {
	// Create inserts a new lesson entry into the database.
	Create(context.Context, *dto.CreateLessonRequest) (*models.Lesson, error)

	// GetAll retrieves all lessons with optional filtering and pagination.
	// Returns the list of lessons and total count.
	GetAll(context.Context, *dto.LessonFilterRequest) ([]models.Lesson, int64, error)

	// GetByID retrieves a single lesson entry by its ID.
	GetByID(context.Context, uint) (*models.Lesson, error)

	// GetByCourseIDAndOrderIndex retrieves a lesson entry by course ID and order index.
	// Used for duplicate detection.
	GetByCourseIDAndOrderIndex(context.Context, uint, int) (*models.Lesson, error)

	// GetByCourseID retrieves all lessons for a specific course, ordered by order_index.
	GetByCourseID(context.Context, uint) ([]models.Lesson, error)

	// Update modifies an existing lesson entry by ID.
	Update(context.Context, *dto.UpdateLessonRequest, uint) (*models.Lesson, error)

	// Delete removes a lesson entry by ID.
	Delete(context.Context, uint) error

	// Publish sets a lesson as published with the current timestamp.
	Publish(context.Context, uint) (*models.Lesson, error)

	// Unpublish sets a lesson as unpublished.
	Unpublish(context.Context, uint) (*models.Lesson, error)
}

func NewLessonRepository(db *gorm.DB) ILessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(ctx context.Context, req *dto.CreateLessonRequest) (*models.Lesson, error) {
	lesson := models.Lesson{
		CourseID:         req.CourseID,
		Title:            req.Title,
		Content:          req.Content,
		OrderIndex:       req.OrderIndex,
		EstimatedMinutes: req.EstimatedMinutes,
		IsPublished:      false,
	}

	err := r.db.WithContext(ctx).Create(&lesson).Error
	if err != nil {
		// Check for unique constraint violation on order_index within course
		if strings.Contains(err.Error(), "idx_lesson_course_order") ||
			strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, errConstant.ErrDuplicateOrderIndex
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Load relationships using Preload for efficiency
	err = r.db.WithContext(ctx).
		Preload("Course").
		Preload("Course.JlptLevel").
		First(&lesson, lesson.ID).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &lesson, nil
}

func (r *LessonRepository) GetAll(ctx context.Context, filter *dto.LessonFilterRequest) ([]models.Lesson, int64, error) {
	var lessons []models.Lesson
	var total int64

	// Build base query with filters
	query := r.db.WithContext(ctx).Model(&models.Lesson{})

	// Apply filters
	if filter != nil {
		if filter.CourseID > 0 {
			query = query.Where("course_id = ?", filter.CourseID)
		}
		if filter.IsPublished != nil {
			query = query.Where("is_published = ?", *filter.IsPublished)
		}
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(title) LIKE ?", searchPattern)
		}
	}

	// Count total records with filters applied
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Apply sorting with whitelist validation (defense in depth)
	allowedSortFields := map[string]string{
		"order_index": "order_index",
		"title":       "title",
		"created_at":  "created_at",
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

	query = query.Preload("Course").
		Preload("Course.JlptLevel").
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

	err := query.Find(&lessons).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return lessons, total, nil
}

func (r *LessonRepository) GetByID(ctx context.Context, id uint) (*models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Course.JlptLevel").
		Where("id = ?", id).
		First(&lesson).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrLessonNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &lesson, nil
}

func (r *LessonRepository) GetByCourseIDAndOrderIndex(ctx context.Context, courseID uint, orderIndex int) (*models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.WithContext(ctx).
		Where("course_id = ? AND order_index = ?", courseID, orderIndex).
		First(&lesson).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrLessonNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &lesson, nil
}

func (r *LessonRepository) GetByCourseID(ctx context.Context, courseID uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Course.JlptLevel").
		Where("course_id = ?", courseID).
		Order("order_index ASC").
		Find(&lessons).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return lessons, nil
}

func (r *LessonRepository) Update(ctx context.Context, req *dto.UpdateLessonRequest, id uint) (*models.Lesson, error) {
	lesson := models.Lesson{
		CourseID:         req.CourseID,
		Title:            req.Title,
		Content:          req.Content,
		OrderIndex:       req.OrderIndex,
		EstimatedMinutes: req.EstimatedMinutes,
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&lesson)

	if result.Error != nil {
		// Check for unique constraint violation on order_index within course
		if strings.Contains(result.Error.Error(), "idx_lesson_course_order") ||
			strings.Contains(result.Error.Error(), "duplicate key") ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint") {
			return nil, errConstant.ErrDuplicateOrderIndex
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrLessonNotFound
	}

	// Fetch the updated record with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Course.JlptLevel").
		Where("id = ?", id).
		First(&lesson).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &lesson, nil
}

func (r *LessonRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Lesson{})

	if result.Error != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errConstant.ErrLessonNotFound
	}

	return nil
}

func (r *LessonRepository) Publish(ctx context.Context, id uint) (*models.Lesson, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.Lesson{}).
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
		return nil, errConstant.ErrLessonNotFound
	}

	// Fetch the updated record with preloaded relationships
	var lesson models.Lesson
	err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Course.JlptLevel").
		Where("id = ?", id).
		First(&lesson).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &lesson, nil
}

func (r *LessonRepository) Unpublish(ctx context.Context, id uint) (*models.Lesson, error) {
	result := r.db.WithContext(ctx).
		Model(&models.Lesson{}).
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
		return nil, errConstant.ErrLessonNotFound
	}

	// Fetch the updated record with preloaded relationships
	var lesson models.Lesson
	err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Course.JlptLevel").
		Where("id = ?", id).
		First(&lesson).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &lesson, nil
}
