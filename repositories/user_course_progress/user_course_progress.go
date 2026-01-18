package repositories

import (
	"context"
	"errors"
	errWrap "manabu-service/common/error"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserCourseProgressRepository struct {
	db *gorm.DB
}

// IUserCourseProgressRepository defines the contract for user course progress data access operations.
type IUserCourseProgressRepository interface {
	// Create inserts a new user course progress entry into the database.
	Create(context.Context, uint, *dto.CreateUserCourseProgressRequest) (*models.UserCourseProgress, error)

	// GetAll retrieves all user course progress entries with optional filtering and pagination.
	// Returns the list of progress entries and total count.
	GetAll(context.Context, uint, *dto.UserCourseProgressFilterRequest) ([]models.UserCourseProgress, int64, error)

	// GetByID retrieves a single user course progress entry by its UUID.
	GetByID(context.Context, uuid.UUID) (*models.UserCourseProgress, error)

	// GetByUserIDAndCourseID retrieves user course progress by user ID and course ID.
	// Used to check if user is already enrolled and for duplicate detection.
	GetByUserIDAndCourseID(context.Context, uint, uint) (*models.UserCourseProgress, error)

	// Update modifies an existing user course progress entry by UUID.
	Update(context.Context, *dto.UpdateUserCourseProgressRequest, uuid.UUID, int) (*models.UserCourseProgress, error)
}

func NewUserCourseProgressRepository(db *gorm.DB) IUserCourseProgressRepository {
	return &UserCourseProgressRepository{db: db}
}

func (r *UserCourseProgressRepository) Create(ctx context.Context, userID uint, req *dto.CreateUserCourseProgressRequest) (*models.UserCourseProgress, error) {
	// First, get the total number of lessons for the course
	var lessonCount int64
	err := r.db.WithContext(ctx).
		Model(&models.Lesson{}).
		Where("course_id = ?", req.CourseID).
		Count(&lessonCount).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	progress := models.UserCourseProgress{
		UserID:             userID,
		CourseID:           req.CourseID,
		Status:             models.ProgressStatusNotStarted,
		ProgressPercentage: 0.00,
		CompletedLessons:   0,
		TotalLessons:       int(lessonCount),
	}

	err = r.db.WithContext(ctx).Create(&progress).Error
	if err != nil {
		// Check for unique constraint violation on user_id and course_id
		if strings.Contains(err.Error(), "idx_user_course_progress_user_course") ||
			strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, errConstant.ErrUserCourseProgressAlreadyExists
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Load relationships using Preload for efficiency
	err = r.db.WithContext(ctx).
		Preload("Course").
		Preload("Course.JlptLevel").
		First(&progress, "id = ?", progress.ID).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &progress, nil
}

func (r *UserCourseProgressRepository) GetAll(ctx context.Context, userID uint, filter *dto.UserCourseProgressFilterRequest) ([]models.UserCourseProgress, int64, error) {
	var progressList []models.UserCourseProgress
	var total int64

	// Build base query with user filter
	query := r.db.WithContext(ctx).Model(&models.UserCourseProgress{}).Where("user_id = ?", userID)

	// Apply filters
	if filter != nil {
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.CourseID > 0 {
			query = query.Where("course_id = ?", filter.CourseID)
		}
	}

	// Count total records with filters applied
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Apply sorting with whitelist validation (defense in depth)
	allowedSortFields := map[string]string{
		"last_accessed_at":    "last_accessed_at",
		"progress_percentage": "progress_percentage",
		"started_at":          "started_at",
	}
	allowedSortOrders := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	sortBy := "last_accessed_at"
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

	// Handle NULL values in sorting (put nulls last)
	orderClause := sortBy + " " + sortOrder + " NULLS LAST"

	query = query.Preload("Course").
		Preload("Course.JlptLevel").
		Order(orderClause)

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

	err := query.Find(&progressList).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return progressList, total, nil
}

func (r *UserCourseProgressRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserCourseProgress, error) {
	var progress models.UserCourseProgress
	err := r.db.WithContext(ctx).
		Preload("Course").
		Preload("Course.JlptLevel").
		Where("id = ?", id).
		First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserCourseProgressNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &progress, nil
}

func (r *UserCourseProgressRepository) GetByUserIDAndCourseID(ctx context.Context, userID uint, courseID uint) (*models.UserCourseProgress, error) {
	var progress models.UserCourseProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserCourseProgressNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &progress, nil
}

func (r *UserCourseProgressRepository) Update(ctx context.Context, req *dto.UpdateUserCourseProgressRequest, id uuid.UUID, totalLessons int) (*models.UserCourseProgress, error) {
	// Use transaction to ensure atomicity and prevent race conditions
	var progress models.UserCourseProgress
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the row for update to prevent concurrent modifications (SELECT ... FOR UPDATE)
		var currentProgress models.UserCourseProgress
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", id).
			First(&currentProgress).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errConstant.ErrUserCourseProgressNotFound
			}
			return errWrap.WrapError(errConstant.ErrSQLError)
		}

		// Calculate progress percentage
		progressPercentage := 0.0
		if totalLessons > 0 {
			progressPercentage = (float64(req.CompletedLessons) / float64(totalLessons)) * 100
		}

		// Determine status based on completed lessons
		status := models.ProgressStatusNotStarted
		if req.CompletedLessons > 0 && req.CompletedLessons < totalLessons {
			status = models.ProgressStatusInProgress
		} else if req.CompletedLessons >= totalLessons && totalLessons > 0 {
			status = models.ProgressStatusCompleted
		}

		// Build update map
		updateData := map[string]interface{}{
			"completed_lessons":   req.CompletedLessons,
			"progress_percentage": progressPercentage,
			"status":              status,
		}

		// Set started_at if transitioning from not_started
		if currentProgress.Status == models.ProgressStatusNotStarted && req.CompletedLessons > 0 {
			now := gorm.Expr("CURRENT_TIMESTAMP")
			updateData["started_at"] = now
		}

		// Set completed_at if status is completed
		if status == models.ProgressStatusCompleted && currentProgress.Status != models.ProgressStatusCompleted {
			now := gorm.Expr("CURRENT_TIMESTAMP")
			updateData["completed_at"] = now
		}

		// Always update last_accessed_at
		updateData["last_accessed_at"] = gorm.Expr("CURRENT_TIMESTAMP")

		result := tx.Model(&models.UserCourseProgress{}).
			Where("id = ?", id).
			Updates(updateData)

		if result.Error != nil {
			return errWrap.WrapError(errConstant.ErrSQLError)
		}

		// Check if any rows were affected
		if result.RowsAffected == 0 {
			return errConstant.ErrUserCourseProgressNotFound
		}

		// Fetch the updated record with preloaded relationships
		err = tx.Preload("Course").
			Preload("Course.JlptLevel").
			Where("id = ?", id).
			First(&progress).Error
		if err != nil {
			return errWrap.WrapError(errConstant.ErrSQLError)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &progress, nil
}
