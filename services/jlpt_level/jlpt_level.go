package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/repositories"
)

type JlptLevelService struct {
	repository repositories.IRepositoryRegistry
}

type IJlptLevelService interface {
	Create(context.Context, *dto.CreateJlptLevelRequest) (*dto.JlptLevelResponse, error)
	GetAll(context.Context) ([]dto.JlptLevelResponse, error)
	GetByID(context.Context, uint) (*dto.JlptLevelResponse, error)
	Update(context.Context, *dto.UpdateJlptLevelRequest, uint) (*dto.JlptLevelResponse, error)
	Delete(context.Context, uint) error
}

func NewJlptLevelService(repository repositories.IRepositoryRegistry) IJlptLevelService {
	return &JlptLevelService{repository: repository}
}

func (s *JlptLevelService) isCodeExist(ctx context.Context, code string) bool {
	jlptLevel, err := s.repository.GetJlptLevel().GetByCode(ctx, code)
	if err != nil {
		return false
	}
	return jlptLevel != nil
}

func (s *JlptLevelService) isLevelOrderExist(ctx context.Context, levelOrder int) bool {
	jlptLevel, err := s.repository.GetJlptLevel().GetByLevelOrder(ctx, levelOrder)
	if err != nil {
		return false
	}
	return jlptLevel != nil
}

func (s *JlptLevelService) Create(ctx context.Context, req *dto.CreateJlptLevelRequest) (*dto.JlptLevelResponse, error) {
	// Check if code already exists
	if s.isCodeExist(ctx, req.Code) {
		return nil, errConstant.ErrJlptLevelCodeExist
	}

	// Check if level order already exists
	if s.isLevelOrderExist(ctx, req.LevelOrder) {
		return nil, errConstant.ErrJlptLevelOrderExist
	}

	jlptLevel, err := s.repository.GetJlptLevel().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	response := &dto.JlptLevelResponse{
		ID:          jlptLevel.ID,
		Code:        jlptLevel.Code,
		Name:        jlptLevel.Name,
		Description: jlptLevel.Description,
		LevelOrder:  jlptLevel.LevelOrder,
	}

	return response, nil
}

func (s *JlptLevelService) GetAll(ctx context.Context) ([]dto.JlptLevelResponse, error) {
	jlptLevels, err := s.repository.GetJlptLevel().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.JlptLevelResponse
	for _, jlptLevel := range jlptLevels {
		responses = append(responses, dto.JlptLevelResponse{
			ID:          jlptLevel.ID,
			Code:        jlptLevel.Code,
			Name:        jlptLevel.Name,
			Description: jlptLevel.Description,
			LevelOrder:  jlptLevel.LevelOrder,
		})
	}

	return responses, nil
}

func (s *JlptLevelService) GetByID(ctx context.Context, id uint) (*dto.JlptLevelResponse, error) {
	jlptLevel, err := s.repository.GetJlptLevel().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := &dto.JlptLevelResponse{
		ID:          jlptLevel.ID,
		Code:        jlptLevel.Code,
		Name:        jlptLevel.Name,
		Description: jlptLevel.Description,
		LevelOrder:  jlptLevel.LevelOrder,
	}

	return response, nil
}

func (s *JlptLevelService) Update(ctx context.Context, req *dto.UpdateJlptLevelRequest, id uint) (*dto.JlptLevelResponse, error) {
	// Check if JLPT level exists
	existingLevel, err := s.repository.GetJlptLevel().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if code already exists (excluding current record)
	if s.isCodeExist(ctx, req.Code) && existingLevel.Code != req.Code {
		checkCode, err := s.repository.GetJlptLevel().GetByCode(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		if checkCode != nil {
			return nil, errConstant.ErrJlptLevelCodeExist
		}
	}

	// Check if level order already exists (excluding current record)
	if s.isLevelOrderExist(ctx, req.LevelOrder) && existingLevel.LevelOrder != req.LevelOrder {
		checkOrder, err := s.repository.GetJlptLevel().GetByLevelOrder(ctx, req.LevelOrder)
		if err != nil {
			return nil, err
		}
		if checkOrder != nil {
			return nil, errConstant.ErrJlptLevelOrderExist
		}
	}

	jlptLevel, err := s.repository.GetJlptLevel().Update(ctx, req, id)
	if err != nil {
		return nil, err
	}

	response := &dto.JlptLevelResponse{
		ID:          jlptLevel.ID,
		Code:        jlptLevel.Code,
		Name:        jlptLevel.Name,
		Description: jlptLevel.Description,
		LevelOrder:  jlptLevel.LevelOrder,
	}

	return response, nil
}

func (s *JlptLevelService) Delete(ctx context.Context, id uint) error {
	// Check if JLPT level exists
	_, err := s.repository.GetJlptLevel().GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repository.GetJlptLevel().Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
