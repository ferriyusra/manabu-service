package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"math"
)

type VocabularyService struct {
	repository repositories.IRepositoryRegistry
}

// IVocabularyService defines the contract for vocabulary business logic operations.
type IVocabularyService interface {
	// Create validates and creates a new vocabulary entry.
	// Validates JLPT level and category existence, checks for duplicates.
	Create(context.Context, *dto.CreateVocabularyRequest) (*dto.VocabularyResponse, error)

	// GetAll retrieves all vocabularies with filtering, sorting, and pagination.
	GetAll(context.Context, *dto.VocabularyFilterRequest) (*dto.VocabularyListResponse, error)

	// GetByID retrieves a single vocabulary entry by its ID.
	GetByID(context.Context, uint) (*dto.VocabularyResponse, error)

	// Update validates and updates an existing vocabulary entry.
	// Validates JLPT level and category existence, checks for duplicates.
	Update(context.Context, *dto.UpdateVocabularyRequest, uint) (*dto.VocabularyResponse, error)

	// Delete removes a vocabulary entry by ID if it exists.
	Delete(context.Context, uint) error
}

func NewVocabularyService(repository repositories.IRepositoryRegistry) IVocabularyService {
	return &VocabularyService{repository: repository}
}

// toVocabularyResponse converts a Vocabulary model to VocabularyResponse DTO
func (s *VocabularyService) toVocabularyResponse(vocabulary *models.Vocabulary) *dto.VocabularyResponse {
	response := &dto.VocabularyResponse{
		ID:                     vocabulary.ID,
		Word:                   vocabulary.Word,
		Reading:                vocabulary.Reading,
		Meaning:                vocabulary.Meaning,
		PartOfSpeech:           vocabulary.PartOfSpeech,
		JlptLevelID:            vocabulary.JlptLevelID,
		CategoryID:             vocabulary.CategoryID,
		ExampleSentence:        vocabulary.ExampleSentence,
		ExampleSentenceReading: vocabulary.ExampleSentenceReading,
		ExampleSentenceMeaning: vocabulary.ExampleSentenceMeaning,
		AudioURL:               vocabulary.AudioURL,
		ImageURL:               vocabulary.ImageURL,
		Difficulty:             vocabulary.Difficulty,
	}

	if vocabulary.JlptLevel.ID > 0 {
		response.JlptLevel = &dto.JlptLevelResponse{
			ID:          vocabulary.JlptLevel.ID,
			Code:        vocabulary.JlptLevel.Code,
			Name:        vocabulary.JlptLevel.Name,
			Description: vocabulary.JlptLevel.Description,
			LevelOrder:  vocabulary.JlptLevel.LevelOrder,
		}
	}

	if vocabulary.Category.ID > 0 {
		response.Category = &dto.CategoryResponse{
			ID:          vocabulary.Category.ID,
			Name:        vocabulary.Category.Name,
			Description: vocabulary.Category.Description,
			JlptLevelID: vocabulary.Category.JlptLevelID,
		}

		// Include nested JlptLevel in Category if available
		if vocabulary.Category.JlptLevel.ID > 0 {
			response.Category.JlptLevel = &dto.JlptLevelResponse{
				ID:          vocabulary.Category.JlptLevel.ID,
				Code:        vocabulary.Category.JlptLevel.Code,
				Name:        vocabulary.Category.JlptLevel.Name,
				Description: vocabulary.Category.JlptLevel.Description,
				LevelOrder:  vocabulary.Category.JlptLevel.LevelOrder,
			}
		}
	}

	return response
}

func (s *VocabularyService) isVocabularyExist(ctx context.Context, word string, jlptLevelID uint) bool {
	vocabulary, err := s.repository.GetVocabulary().GetByWordAndJlptLevel(ctx, word, jlptLevelID)
	if err != nil {
		return false
	}
	return vocabulary != nil
}

func (s *VocabularyService) isJlptLevelExist(ctx context.Context, jlptLevelID uint) bool {
	jlptLevel, err := s.repository.GetJlptLevel().GetByID(ctx, jlptLevelID)
	if err != nil {
		return false
	}
	return jlptLevel != nil
}

func (s *VocabularyService) isCategoryExist(ctx context.Context, categoryID uint) bool {
	category, err := s.repository.GetCategory().GetByID(ctx, categoryID)
	if err != nil {
		return false
	}
	return category != nil
}

func (s *VocabularyService) validateDifficulty(difficulty int) error {
	if difficulty < 1 || difficulty > 5 {
		return errConstant.ErrInvalidDifficulty
	}
	return nil
}

func (s *VocabularyService) Create(ctx context.Context, req *dto.CreateVocabularyRequest) (*dto.VocabularyResponse, error) {
	// Validate JLPT level exists
	if !s.isJlptLevelExist(ctx, req.JlptLevelID) {
		return nil, errConstant.ErrInvalidJlptLevelID
	}

	// Validate Category exists
	if !s.isCategoryExist(ctx, req.CategoryID) {
		return nil, errConstant.ErrInvalidCategoryID
	}

	// Validate difficulty range
	if req.Difficulty > 0 {
		if err := s.validateDifficulty(req.Difficulty); err != nil {
			return nil, err
		}
	}

	// Check if vocabulary word already exists for this JLPT level
	if s.isVocabularyExist(ctx, req.Word, req.JlptLevelID) {
		return nil, errConstant.ErrVocabularyDuplicate
	}

	vocabulary, err := s.repository.GetVocabulary().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.toVocabularyResponse(vocabulary), nil
}

func (s *VocabularyService) GetAll(ctx context.Context, filter *dto.VocabularyFilterRequest) (*dto.VocabularyListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.VocabularyFilterRequest{
			PaginationRequest: dto.PaginationRequest{
				Page:  1,
				Limit: 10,
			},
		}
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	// Validate difficulty if provided
	if filter.Difficulty > 0 {
		if err := s.validateDifficulty(filter.Difficulty); err != nil {
			return nil, err
		}
	}

	vocabularies, total, err := s.repository.GetVocabulary().GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.VocabularyResponse, 0, len(vocabularies))
	for _, vocabulary := range vocabularies {
		responses = append(responses, *s.toVocabularyResponse(&vocabulary))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.VocabularyListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *VocabularyService) GetByID(ctx context.Context, id uint) (*dto.VocabularyResponse, error) {
	vocabulary, err := s.repository.GetVocabulary().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toVocabularyResponse(vocabulary), nil
}

func (s *VocabularyService) Update(ctx context.Context, req *dto.UpdateVocabularyRequest, id uint) (*dto.VocabularyResponse, error) {
	// Check if vocabulary exists
	existingVocabulary, err := s.repository.GetVocabulary().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate JLPT level exists
	if !s.isJlptLevelExist(ctx, req.JlptLevelID) {
		return nil, errConstant.ErrInvalidJlptLevelID
	}

	// Validate Category exists
	if !s.isCategoryExist(ctx, req.CategoryID) {
		return nil, errConstant.ErrInvalidCategoryID
	}

	// Validate difficulty range
	if req.Difficulty > 0 {
		if err := s.validateDifficulty(req.Difficulty); err != nil {
			return nil, err
		}
	}

	// Check if vocabulary word already exists for this JLPT level (excluding current record)
	// Only check if word or JLPT level is being changed
	if existingVocabulary.Word != req.Word || existingVocabulary.JlptLevelID != req.JlptLevelID {
		checkVocabulary, err := s.repository.GetVocabulary().GetByWordAndJlptLevel(ctx, req.Word, req.JlptLevelID)
		if err != nil && err != errConstant.ErrVocabularyNotFound {
			return nil, err // Only propagate real errors, not "not found"
		}
		if checkVocabulary != nil && checkVocabulary.ID != id {
			return nil, errConstant.ErrVocabularyDuplicate
		}
	}

	vocabulary, err := s.repository.GetVocabulary().Update(ctx, req, id)
	if err != nil {
		return nil, err
	}

	return s.toVocabularyResponse(vocabulary), nil
}

func (s *VocabularyService) Delete(ctx context.Context, id uint) error {
	// Check if vocabulary exists
	_, err := s.repository.GetVocabulary().GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repository.GetVocabulary().Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
