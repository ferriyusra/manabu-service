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

type JlptLevelRepository struct {
	db *gorm.DB
}

type IJlptLevelRepository interface {
	Create(context.Context, *dto.CreateJlptLevelRequest) (*models.JlptLevel, error)
	GetAll(context.Context) ([]models.JlptLevel, error)
	GetByID(context.Context, uint) (*models.JlptLevel, error)
	GetByCode(context.Context, string) (*models.JlptLevel, error)
	GetByLevelOrder(context.Context, int) (*models.JlptLevel, error)
	Update(context.Context, *dto.UpdateJlptLevelRequest, uint) (*models.JlptLevel, error)
	Delete(context.Context, uint) error
}

func NewJlptLevelRepository(db *gorm.DB) IJlptLevelRepository {
	return &JlptLevelRepository{db: db}
}

func (r *JlptLevelRepository) Create(ctx context.Context, req *dto.CreateJlptLevelRequest) (*models.JlptLevel, error) {
	jlptLevel := models.JlptLevel{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		LevelOrder:  req.LevelOrder,
	}

	err := r.db.WithContext(ctx).Create(&jlptLevel).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &jlptLevel, nil
}

func (r *JlptLevelRepository) GetAll(ctx context.Context) ([]models.JlptLevel, error) {
	var jlptLevels []models.JlptLevel
	err := r.db.WithContext(ctx).
		Order("level_order ASC").
		Find(&jlptLevels).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return jlptLevels, nil
}

func (r *JlptLevelRepository) GetByID(ctx context.Context, id uint) (*models.JlptLevel, error) {
	var jlptLevel models.JlptLevel
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&jlptLevel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrJlptLevelNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &jlptLevel, nil
}

func (r *JlptLevelRepository) GetByCode(ctx context.Context, code string) (*models.JlptLevel, error) {
	var jlptLevel models.JlptLevel
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&jlptLevel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrJlptLevelNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &jlptLevel, nil
}

func (r *JlptLevelRepository) GetByLevelOrder(ctx context.Context, levelOrder int) (*models.JlptLevel, error) {
	var jlptLevel models.JlptLevel
	err := r.db.WithContext(ctx).
		Where("level_order = ?", levelOrder).
		First(&jlptLevel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrJlptLevelNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &jlptLevel, nil
}

func (r *JlptLevelRepository) Update(ctx context.Context, req *dto.UpdateJlptLevelRequest, id uint) (*models.JlptLevel, error) {
	jlptLevel := models.JlptLevel{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		LevelOrder:  req.LevelOrder,
	}

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&jlptLevel).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Fetch the updated record to return complete data
	err = r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&jlptLevel).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &jlptLevel, nil
}

func (r *JlptLevelRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.JlptLevel{}).Error
	if err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
