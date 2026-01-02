package repositories

import (
	"context"
	"errors"
	errWrap "manabu-service/common/error"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

type ICategoryRepository interface {
	Create(context.Context, *dto.CreateCategoryRequest) (*models.Category, error)
	GetAll(context.Context, *dto.PaginationRequest) ([]models.Category, int64, error)
	GetByID(context.Context, uint) (*models.Category, error)
	GetByJlptLevelID(context.Context, uint, *dto.PaginationRequest) ([]models.Category, int64, error)
	GetByNameAndJlptLevel(context.Context, string, uint) (*models.Category, error)
	Update(context.Context, *dto.UpdateCategoryRequest, uint) (*models.Category, error)
	Delete(context.Context, uint) error
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*models.Category, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		JlptLevelID: req.JlptLevelID,
	}

	err := r.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Load JlptLevel relationship using Association (more efficient than separate query)
	err = r.db.WithContext(ctx).Model(&category).Association("JlptLevel").Find(&category.JlptLevel)
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &category, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context, pagination *dto.PaginationRequest) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Build query with pagination
	query := r.db.WithContext(ctx).Preload("JlptLevel").Order("created_at DESC")

	if pagination != nil && pagination.Limit > 0 {
		// Defensive validation: ensure page is at least 1
		page := pagination.Page
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * pagination.Limit
		// Ensure offset is never negative
		if offset < 0 {
			offset = 0
		}

		query = query.Limit(pagination.Limit).Offset(offset)
	}

	err := query.Find(&categories).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return categories, total, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Where("id = ?", id).
		First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrCategoryNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &category, nil
}

func (r *CategoryRepository) GetByJlptLevelID(ctx context.Context, jlptLevelID uint, pagination *dto.PaginationRequest) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total records for this JLPT level
	if err := r.db.WithContext(ctx).Model(&models.Category{}).Where("jlpt_level_id = ?", jlptLevelID).Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Build query with pagination
	query := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Where("jlpt_level_id = ?", jlptLevelID).
		Order("created_at DESC")

	if pagination != nil && pagination.Limit > 0 {
		// Defensive validation: ensure page is at least 1
		page := pagination.Page
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * pagination.Limit
		// Ensure offset is never negative
		if offset < 0 {
			offset = 0
		}

		query = query.Limit(pagination.Limit).Offset(offset)
	}

	err := query.Find(&categories).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return categories, total, nil
}

func (r *CategoryRepository) GetByNameAndJlptLevel(ctx context.Context, name string, jlptLevelID uint) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Where("name = ? AND jlpt_level_id = ?", name, jlptLevelID).
		First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrCategoryNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, req *dto.UpdateCategoryRequest, id uint) (*models.Category, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		JlptLevelID: req.JlptLevelID,
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&category)

	if result.Error != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrCategoryNotFound
	}

	// Fetch the updated record with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Where("id = ?", id).
		First(&category).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &category, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Category{}).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
