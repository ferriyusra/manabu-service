package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"math"
)

type CategoryService struct {
	repository repositories.IRepositoryRegistry
}

type ICategoryService interface {
	Create(context.Context, *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetAll(context.Context, *dto.PaginationRequest) (*dto.CategoryListResponse, error)
	GetByID(context.Context, uint) (*dto.CategoryResponse, error)
	GetByJlptLevelID(context.Context, uint, *dto.PaginationRequest) (*dto.CategoryListResponse, error)
	Update(context.Context, *dto.UpdateCategoryRequest, uint) (*dto.CategoryResponse, error)
	Delete(context.Context, uint) error
}

func NewCategoryService(repository repositories.IRepositoryRegistry) ICategoryService {
	return &CategoryService{repository: repository}
}

// toCategoryResponse converts a Category model to CategoryResponse DTO
func (s *CategoryService) toCategoryResponse(category *models.Category) *dto.CategoryResponse {
	response := &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		JlptLevelID: category.JlptLevelID,
	}

	if category.JlptLevel.ID > 0 {
		response.JlptLevel = &dto.JlptLevelResponse{
			ID:          category.JlptLevel.ID,
			Code:        category.JlptLevel.Code,
			Name:        category.JlptLevel.Name,
			Description: category.JlptLevel.Description,
			LevelOrder:  category.JlptLevel.LevelOrder,
		}
	}

	return response
}

func (s *CategoryService) isCategoryNameExist(ctx context.Context, name string, jlptLevelID uint) bool {
	category, err := s.repository.GetCategory().GetByNameAndJlptLevel(ctx, name, jlptLevelID)
	if err != nil {
		return false
	}
	return category != nil
}

func (s *CategoryService) isJlptLevelExist(ctx context.Context, jlptLevelID uint) bool {
	jlptLevel, err := s.repository.GetJlptLevel().GetByID(ctx, jlptLevelID)
	if err != nil {
		return false
	}
	return jlptLevel != nil
}

func (s *CategoryService) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Validate JLPT level exists
	if !s.isJlptLevelExist(ctx, req.JlptLevelID) {
		return nil, errConstant.ErrInvalidJlptLevelID
	}

	// Check if category name already exists for this JLPT level
	if s.isCategoryNameExist(ctx, req.Name, req.JlptLevelID) {
		return nil, errConstant.ErrCategoryNameExist
	}

	category, err := s.repository.GetCategory().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *CategoryService) GetAll(ctx context.Context, pagination *dto.PaginationRequest) (*dto.CategoryListResponse, error) {
	// Set default pagination values
	if pagination == nil {
		pagination = &dto.PaginationRequest{
			Page:  1,
			Limit: 10,
		}
	}
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 {
		pagination.Limit = 10
	}
	if pagination.Limit > 100 {
		pagination.Limit = 100
	}

	categories, total, err := s.repository.GetCategory().GetAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, *s.toCategoryResponse(&category))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))

	return &dto.CategoryListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id uint) (*dto.CategoryResponse, error) {
	category, err := s.repository.GetCategory().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *CategoryService) GetByJlptLevelID(ctx context.Context, jlptLevelID uint, pagination *dto.PaginationRequest) (*dto.CategoryListResponse, error) {
	// Validate JLPT level exists
	if !s.isJlptLevelExist(ctx, jlptLevelID) {
		return nil, errConstant.ErrInvalidJlptLevelID
	}

	// Set default pagination values
	if pagination == nil {
		pagination = &dto.PaginationRequest{
			Page:  1,
			Limit: 10,
		}
	}
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 {
		pagination.Limit = 10
	}
	if pagination.Limit > 100 {
		pagination.Limit = 100
	}

	categories, total, err := s.repository.GetCategory().GetByJlptLevelID(ctx, jlptLevelID, pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, *s.toCategoryResponse(&category))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))

	return &dto.CategoryListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *CategoryService) Update(ctx context.Context, req *dto.UpdateCategoryRequest, id uint) (*dto.CategoryResponse, error) {
	// Check if category exists
	existingCategory, err := s.repository.GetCategory().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate JLPT level exists
	if !s.isJlptLevelExist(ctx, req.JlptLevelID) {
		return nil, errConstant.ErrInvalidJlptLevelID
	}

	// Check if category name already exists for this JLPT level (excluding current record)
	// Only check if name or JLPT level is being changed
	if existingCategory.Name != req.Name || existingCategory.JlptLevelID != req.JlptLevelID {
		checkCategory, err := s.repository.GetCategory().GetByNameAndJlptLevel(ctx, req.Name, req.JlptLevelID)
		if err != nil && err != errConstant.ErrCategoryNotFound {
			return nil, err // Only propagate real errors, not "not found"
		}
		if checkCategory != nil && checkCategory.ID != id {
			return nil, errConstant.ErrCategoryNameExist
		}
	}

	category, err := s.repository.GetCategory().Update(ctx, req, id)
	if err != nil {
		return nil, err
	}

	return s.toCategoryResponse(category), nil
}

func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	// Check if category exists
	_, err := s.repository.GetCategory().GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repository.GetCategory().Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
