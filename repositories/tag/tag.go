package repositories

import (
	"context"
	"errors"
	errWrap "manabu-service/common/error"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"strings"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

// ITagRepository defines the contract for tag data access operations.
type ITagRepository interface {
	// Create inserts a new tag into the database.
	Create(context.Context, *dto.CreateTagRequest) (*models.Tag, error)

	// GetAll retrieves all tags with optional search and pagination.
	// Returns the list of tags and total count.
	GetAll(context.Context, *dto.TagFilterRequest) ([]models.Tag, int64, error)

	// GetByID retrieves a single tag by its ID.
	GetByID(context.Context, uint) (*models.Tag, error)

	// GetByName retrieves a tag by its name (case-insensitive).
	// Used for duplicate detection.
	GetByName(context.Context, string) (*models.Tag, error)

	// Update modifies an existing tag by ID.
	Update(context.Context, *dto.UpdateTagRequest, uint) (*models.Tag, error)

	// Delete removes a tag by ID.
	Delete(context.Context, uint) error
}

func NewTagRepository(db *gorm.DB) ITagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, req *dto.CreateTagRequest) (*models.Tag, error) {
	tag := models.Tag{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}

	err := r.db.WithContext(ctx).Create(&tag).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &tag, nil
}

func (r *TagRepository) GetAll(ctx context.Context, filter *dto.TagFilterRequest) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64

	// Build base query
	query := r.db.WithContext(ctx).Model(&models.Tag{})

	// Apply search filter (case-insensitive search on name)
	if filter != nil && filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?)", searchPattern)
	}

	// Count total records with filters applied
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Apply default ordering by creation date (newest first)
	query = query.Order("created_at DESC")

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

	err := query.Find(&tags).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return tags, total, nil
}

func (r *TagRepository) GetByID(ctx context.Context, id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrTagNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &tag, nil
}

func (r *TagRepository) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	var tag models.Tag
	// Case-insensitive search using LOWER
	err := r.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?)", strings.TrimSpace(name)).
		First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrTagNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &tag, nil
}

func (r *TagRepository) Update(ctx context.Context, req *dto.UpdateTagRequest, id uint) (*models.Tag, error) {
	// Prepare update map to only update non-empty fields
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}

	result := r.db.WithContext(ctx).
		Model(&models.Tag{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrTagNotFound
	}

	// Fetch the updated record
	var tag models.Tag
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&tag).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &tag, nil
}

func (r *TagRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Tag{})

	if result.Error != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errConstant.ErrTagNotFound
	}

	return nil
}
