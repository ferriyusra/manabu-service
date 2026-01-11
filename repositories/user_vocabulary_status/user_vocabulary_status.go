package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "manabu-service/common/error"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"time"

	"gorm.io/gorm"
)

type UserVocabularyStatusRepository struct {
	db *gorm.DB
}

// IUserVocabularyStatusRepository defines the interface for user vocabulary status repository operations
type IUserVocabularyStatusRepository interface {
	Create(context.Context, *models.UserVocabularyStatus) (*models.UserVocabularyStatus, error)
	GetByID(context.Context, uint) (*models.UserVocabularyStatus, error)
	GetByUserAndVocabulary(context.Context, string, uint) (*models.UserVocabularyStatus, error)
	GetByUserID(context.Context, string, *dto.UserVocabStatusListRequest) ([]*models.UserVocabularyStatus, int64, error)
	GetDueForReview(context.Context, string) ([]*models.UserVocabularyStatus, error)
	Update(context.Context, *models.UserVocabularyStatus) (*models.UserVocabularyStatus, error)
}

func NewUserVocabularyStatusRepository(db *gorm.DB) IUserVocabularyStatusRepository {
	return &UserVocabularyStatusRepository{db: db}
}

// Create creates a new user vocabulary status record
func (r *UserVocabularyStatusRepository) Create(ctx context.Context, status *models.UserVocabularyStatus) (*models.UserVocabularyStatus, error) {
	err := r.db.WithContext(ctx).Create(status).Error
	if err != nil {
		// Check for unique constraint violation
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			(err != nil && (errors.Is(err, gorm.ErrCheckConstraintViolated) ||
				err.Error() == "duplicate key value violates unique constraint \"idx_user_vocabulary\"")) {
			return nil, errConstant.ErrVocabAlreadyLearning
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return status, nil
}

// GetByID retrieves a user vocabulary status by ID
func (r *UserVocabularyStatusRepository) GetByID(ctx context.Context, id uint) (*models.UserVocabularyStatus, error) {
	var status models.UserVocabularyStatus
	err := r.db.WithContext(ctx).
		Preload("Vocabulary").
		Where("id = ?", id).
		First(&status).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserVocabStatusNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &status, nil
}

// GetByUserAndVocabulary retrieves a user vocabulary status by user ID and vocabulary ID
func (r *UserVocabularyStatusRepository) GetByUserAndVocabulary(ctx context.Context, userID string, vocabularyID uint) (*models.UserVocabularyStatus, error) {
	var status models.UserVocabularyStatus
	err := r.db.WithContext(ctx).
		Where("user_id = ?::uuid AND vocabulary_id = ?", userID, vocabularyID).
		First(&status).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserVocabStatusNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &status, nil
}

// GetByUserID retrieves all vocabulary statuses for a user with pagination and filtering
func (r *UserVocabularyStatusRepository) GetByUserID(ctx context.Context, userID string, params *dto.UserVocabStatusListRequest) ([]*models.UserVocabularyStatus, int64, error) {
	var statuses []*models.UserVocabularyStatus
	var total int64

	query := r.db.WithContext(ctx).
		Preload("Vocabulary").
		Where("user_id = ?::uuid", userID)

	// Apply status filter if provided
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Count total
	countQuery := r.db.WithContext(ctx).Model(&models.UserVocabularyStatus{}).
		Where("user_id = ?::uuid", userID)
	if params.Status != "" {
		countQuery = countQuery.Where("status = ?", params.Status)
	}
	countQuery.Count(&total)

	// Apply sorting
	sort := "next_review_date"
	if params.Sort != "" {
		sort = params.Sort
	}
	order := "asc"
	if params.Order != "" {
		order = params.Order
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	err := query.Order(fmt.Sprintf("%s %s", sort, order)).
		Limit(params.Limit).
		Offset(offset).
		Find(&statuses).Error

	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return statuses, total, nil
}

// GetDueForReview retrieves all vocabulary statuses that are due for review
func (r *UserVocabularyStatusRepository) GetDueForReview(ctx context.Context, userID string) ([]*models.UserVocabularyStatus, error) {
	var statuses []*models.UserVocabularyStatus

	err := r.db.WithContext(ctx).
		Preload("Vocabulary").
		Where("user_id = ?::uuid AND next_review_date <= ?", userID, time.Now()).
		Order("next_review_date ASC").
		Find(&statuses).Error

	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return statuses, nil
}

// Update updates an existing user vocabulary status
func (r *UserVocabularyStatusRepository) Update(ctx context.Context, status *models.UserVocabularyStatus) (*models.UserVocabularyStatus, error) {
	result := r.db.WithContext(ctx).Save(status)
	if result.Error != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	if result.RowsAffected == 0 {
		return nil, errConstant.ErrUserVocabStatusNotFound
	}

	return status, nil
}
