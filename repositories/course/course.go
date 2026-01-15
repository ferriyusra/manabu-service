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

type CourseRepository struct {
	db *gorm.DB
}

// ICourseRepository defines the contract for course data access operations.
type ICourseRepository interface {
	// Create inserts a new course entry into the database.
	Create(context.Context, *dto.CreateCourseRequest) (*models.Course, error)

	// GetAll retrieves all courses with optional filtering and pagination.
	// Returns the list of courses and total count.
	GetAll(context.Context, *dto.CourseFilterRequest) ([]models.Course, int64, error)

	// GetByID retrieves a single course entry by its ID.
	GetByID(context.Context, uint) (*models.Course, error)

	// GetByTitleAndJlptLevel retrieves a course entry by title and JLPT level.
	// Used for duplicate detection.
	GetByTitleAndJlptLevel(context.Context, string, uint) (*models.Course, error)

	// Update modifies an existing course entry by ID.
	Update(context.Context, *dto.UpdateCourseRequest, uint) (*models.Course, error)

	// Delete removes a course entry by ID.
	Delete(context.Context, uint) error

	// Publish sets a course as published with the current timestamp.
	Publish(context.Context, uint) (*models.Course, error)

	// Unpublish sets a course as unpublished.
	Unpublish(context.Context, uint) (*models.Course, error)

	// GetPublished retrieves only published courses with optional filtering and pagination.
	GetPublished(context.Context, *dto.CourseFilterRequest) ([]models.Course, int64, error)
}

func NewCourseRepository(db *gorm.DB) ICourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(ctx context.Context, req *dto.CreateCourseRequest) (*models.Course, error) {
	// Set default difficulty if not provided
	difficulty := req.Difficulty
	if difficulty == 0 {
		difficulty = 1
	}

	course := models.Course{
		Title:          req.Title,
		Description:    req.Description,
		JlptLevelID:    req.JlptLevelID,
		ThumbnailURL:   req.ThumbnailURL,
		Difficulty:     difficulty,
		EstimatedHours: req.EstimatedHours,
		IsPublished:    false,
	}

	err := r.db.WithContext(ctx).Create(&course).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Load relationships using Preload for efficiency
	err = r.db.WithContext(ctx).
		Preload("JlptLevel").
		First(&course, course.ID).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &course, nil
}

func (r *CourseRepository) GetAll(ctx context.Context, filter *dto.CourseFilterRequest) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	// Build base query with filters
	query := r.db.WithContext(ctx).Model(&models.Course{})

	// Apply filters
	if filter != nil {
		if filter.JlptLevelID > 0 {
			query = query.Where("jlpt_level_id = ?", filter.JlptLevelID)
		}
		if filter.Difficulty > 0 {
			query = query.Where("difficulty = ?", filter.Difficulty)
		}
		if filter.IsPublished != nil {
			query = query.Where("is_published = ?", *filter.IsPublished)
		}
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
				searchPattern, searchPattern)
		}
	}

	// Count total records with filters applied
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Apply sorting with whitelist validation (defense in depth)
	allowedSortFields := map[string]string{
		"title":      "title",
		"difficulty": "difficulty",
		"created_at": "created_at",
	}
	allowedSortOrders := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	sortBy := "created_at"
	sortOrder := "DESC"
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

	query = query.Preload("JlptLevel").
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

	err := query.Find(&courses).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return courses, total, nil
}

func (r *CourseRepository) GetByID(ctx context.Context, id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Where("id = ?", id).
		First(&course).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrCourseNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &course, nil
}

func (r *CourseRepository) GetByTitleAndJlptLevel(ctx context.Context, title string, jlptLevelID uint) (*models.Course, error) {
	var course models.Course
	err := r.db.WithContext(ctx).
		Where("LOWER(title) = LOWER(?) AND jlpt_level_id = ?", title, jlptLevelID).
		First(&course).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrCourseNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &course, nil
}

func (r *CourseRepository) Update(ctx context.Context, req *dto.UpdateCourseRequest, id uint) (*models.Course, error) {
	// Set default difficulty if not provided
	difficulty := req.Difficulty
	if difficulty == 0 {
		difficulty = 1
	}

	course := models.Course{
		Title:          req.Title,
		Description:    req.Description,
		JlptLevelID:    req.JlptLevelID,
		ThumbnailURL:   req.ThumbnailURL,
		Difficulty:     difficulty,
		EstimatedHours: req.EstimatedHours,
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&course)

	if result.Error != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrCourseNotFound
	}

	// Fetch the updated record with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Where("id = ?", id).
		First(&course).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &course, nil
}

func (r *CourseRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Course{})

	if result.Error != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errConstant.ErrCourseNotFound
	}

	return nil
}

func (r *CourseRepository) Publish(ctx context.Context, id uint) (*models.Course, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.Course{}).
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
		return nil, errConstant.ErrCourseNotFound
	}

	// Fetch the updated record with preloaded relationships
	var course models.Course
	err := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Where("id = ?", id).
		First(&course).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &course, nil
}

func (r *CourseRepository) Unpublish(ctx context.Context, id uint) (*models.Course, error) {
	result := r.db.WithContext(ctx).
		Model(&models.Course{}).
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
		return nil, errConstant.ErrCourseNotFound
	}

	// Fetch the updated record with preloaded relationships
	var course models.Course
	err := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Where("id = ?", id).
		First(&course).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &course, nil
}

func (r *CourseRepository) GetPublished(ctx context.Context, filter *dto.CourseFilterRequest) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	// Build base query with is_published = true filter
	query := r.db.WithContext(ctx).Model(&models.Course{}).Where("is_published = ?", true)

	// Apply additional filters
	if filter != nil {
		if filter.JlptLevelID > 0 {
			query = query.Where("jlpt_level_id = ?", filter.JlptLevelID)
		}
		if filter.Difficulty > 0 {
			query = query.Where("difficulty = ?", filter.Difficulty)
		}
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
				searchPattern, searchPattern)
		}
	}

	// Count total records with filters applied
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Apply sorting with whitelist validation (defense in depth)
	allowedSortFields := map[string]string{
		"title":      "title",
		"difficulty": "difficulty",
		"created_at": "created_at",
	}
	allowedSortOrders := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	sortBy := "created_at"
	sortOrder := "DESC"
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

	query = query.Preload("JlptLevel").
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

	err := query.Find(&courses).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return courses, total, nil
}
