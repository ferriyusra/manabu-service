package services

import (
	"context"
	"manabu-service/constants"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"time"
)

type UserVocabularyStatusService struct {
	repository repositories.IRepositoryRegistry
}

// IUserVocabularyStatusService defines the interface for user vocabulary status service operations
type IUserVocabularyStatusService interface {
	Create(context.Context, *dto.CreateUserVocabStatusRequest) (*dto.UserVocabStatusResponse, error)
	GetByID(context.Context, uint) (*dto.UserVocabStatusResponse, error)
	GetAll(context.Context, *dto.UserVocabStatusListRequest) (*dto.UserVocabStatusListResponse, error)
	GetDueForReview(context.Context) ([]dto.UserVocabStatusResponse, error)
	Review(context.Context, uint, *dto.ReviewUserVocabStatusRequest) (*dto.UserVocabStatusResponse, error)
}

func NewUserVocabularyStatusService(repository repositories.IRepositoryRegistry) IUserVocabularyStatusService {
	return &UserVocabularyStatusService{repository: repository}
}


// Create starts learning a new vocabulary for the user
func (s *UserVocabularyStatusService) Create(ctx context.Context, req *dto.CreateUserVocabStatusRequest) (*dto.UserVocabStatusResponse, error) {
	// Get user from context
	userLogin := ctx.Value(constants.UserLogin).(*dto.UserResponse)
	if userLogin == nil {
		return nil, errConstant.ErrUnauthorized
	}

	// Validate vocabulary exists
	vocabulary, err := s.repository.GetVocabulary().GetByID(ctx, req.VocabularyID)
	if err != nil {
		if err == errConstant.ErrVocabularyNotFound {
			return nil, errConstant.ErrVocabularyNotFoundForLearning
		}
		return nil, err
	}

	// Check if user is already learning this vocabulary
	existingStatus, _ := s.repository.GetUserVocabularyStatus().GetByUserAndVocabulary(ctx, userLogin.UUID.String(), req.VocabularyID)
	if existingStatus != nil {
		return nil, errConstant.ErrVocabAlreadyLearning
	}

	// Initialize simple tracking values
	const (
		initialRepetitions = 0
		initialStatus      = "learning"
	)

	// Create user vocabulary status
	status := &models.UserVocabularyStatus{
		UserID:         userLogin.UUID,
		VocabularyID:   vocabulary.ID,
		Status:         initialStatus,
		Repetitions:    initialRepetitions,
		LastReviewedAt: nil,
	}

	createdStatus, err := s.repository.GetUserVocabularyStatus().Create(ctx, status)
	if err != nil {
		return nil, err
	}

	// Map to response DTO (use helper for consistent mapping)
	return s.mapStatusToResponse(createdStatus, userLogin.UUID.String()), nil
}

// GetByID retrieves user vocabulary status by ID
func (s *UserVocabularyStatusService) GetByID(ctx context.Context, id uint) (*dto.UserVocabStatusResponse, error) {
	// Get user from context
	userLogin := ctx.Value(constants.UserLogin).(*dto.UserResponse)
	if userLogin == nil {
		return nil, errConstant.ErrUnauthorized
	}

	// Retrieve status
	status, err := s.repository.GetUserVocabularyStatus().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify the status belongs to the authenticated user
	if status.UserID != userLogin.UUID {
		return nil, errConstant.ErrForbidden
	}

	// Map to response DTO using helper function (includes vocabulary data)
	return s.mapStatusToResponse(status, userLogin.UUID.String()), nil
}

// mapStatusToResponse maps a model to response DTO with vocabulary data
func (s *UserVocabularyStatusService) mapStatusToResponse(status *models.UserVocabularyStatus, userUUID string) *dto.UserVocabStatusResponse {
	response := &dto.UserVocabStatusResponse{
		ID:             status.ID,
		UserID:         userUUID,
		VocabularyID:   status.VocabularyID,
		Status:         status.Status,
		Repetitions:    status.Repetitions,
		LastReviewedAt: status.LastReviewedAt,
		CreatedAt:      *status.CreatedAt,
		UpdatedAt:      *status.UpdatedAt,
	}

	// Map vocabulary if preloaded (check if ID is not zero)
	if status.Vocabulary.ID != 0 {
		response.Vocabulary = &dto.VocabularyResponse{
			ID:                     status.Vocabulary.ID,
			Word:                   status.Vocabulary.Word,
			Reading:                status.Vocabulary.Reading,
			Meaning:                status.Vocabulary.Meaning,
			PartOfSpeech:           status.Vocabulary.PartOfSpeech,
			JlptLevelID:            status.Vocabulary.JlptLevelID,
			CategoryID:             status.Vocabulary.CategoryID,
			ExampleSentence:        status.Vocabulary.ExampleSentence,
			ExampleSentenceReading: status.Vocabulary.ExampleSentenceReading,
			ExampleSentenceMeaning: status.Vocabulary.ExampleSentenceMeaning,
			AudioURL:               status.Vocabulary.AudioURL,
			ImageURL:               status.Vocabulary.ImageURL,
			Difficulty:             status.Vocabulary.Difficulty,
		}
	}

	return response
}

// GetAll retrieves all vocabulary statuses for the authenticated user
func (s *UserVocabularyStatusService) GetAll(ctx context.Context, req *dto.UserVocabStatusListRequest) (*dto.UserVocabStatusListResponse, error) {
	// Get user from context
	userLogin := ctx.Value(constants.UserLogin).(*dto.UserResponse)
	if userLogin == nil {
		return nil, errConstant.ErrUnauthorized
	}

	// Apply default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = "next_review_date"
	}
	if req.Order == "" {
		req.Order = "asc"
	}

	// Retrieve statuses from repository
	statuses, total, err := s.repository.GetUserVocabularyStatus().GetByUserID(ctx, userLogin.UUID.String(), req)
	if err != nil {
		return nil, err
	}

	// Map to response DTOs
	userUUID := userLogin.UUID.String()
	responses := make([]dto.UserVocabStatusResponse, 0, len(statuses))
	for _, status := range statuses {
		responses = append(responses, *s.mapStatusToResponse(status, userUUID))
	}

	// Calculate pagination
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.UserVocabStatusListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

// GetDueForReview retrieves all vocabulary statuses that are due for review
func (s *UserVocabularyStatusService) GetDueForReview(ctx context.Context) ([]dto.UserVocabStatusResponse, error) {
	// Get user from context
	userLogin := ctx.Value(constants.UserLogin).(*dto.UserResponse)
	if userLogin == nil {
		return nil, errConstant.ErrUnauthorized
	}

	// Retrieve due statuses from repository
	statuses, err := s.repository.GetUserVocabularyStatus().GetDueForReview(ctx, userLogin.UUID.String())
	if err != nil {
		return nil, err
	}

	// Map to response DTOs
	userUUID := userLogin.UUID.String()
	responses := make([]dto.UserVocabStatusResponse, 0, len(statuses))
	for _, status := range statuses {
		responses = append(responses, *s.mapStatusToResponse(status, userUUID))
	}

	return responses, nil
}


// Review processes a vocabulary review with simple progress tracking
func (s *UserVocabularyStatusService) Review(ctx context.Context, vocabularyID uint, req *dto.ReviewUserVocabStatusRequest) (*dto.UserVocabStatusResponse, error) {
	// Get user from context
	userLogin := ctx.Value(constants.UserLogin).(*dto.UserResponse)
	if userLogin == nil {
		return nil, errConstant.ErrUnauthorized
	}

	// Retrieve the status by user and vocabulary
	status, err := s.repository.GetUserVocabularyStatus().GetByUserAndVocabulary(ctx, userLogin.UUID.String(), vocabularyID)
	if err != nil {
		return nil, err
	}

	// Apply simple review logic
	now := time.Now()

	if req.IsCorrect {
		// Correct answer: increment progress
		status.Repetitions++

		// After 5 correct reviews, mark as completed
		if status.Repetitions >= 5 {
			status.Status = "completed"
		}
	} else {
		// Incorrect answer: reset progress
		status.Repetitions = 0
		status.Status = "learning"
	}

	// Update last reviewed timestamp
	status.LastReviewedAt = &now

	// Save vocabulary relation before update (will be lost after Save operation)
	vocabulary := status.Vocabulary

	// Save updates to database
	updatedStatus, err := s.repository.GetUserVocabularyStatus().Update(ctx, status)
	if err != nil {
		return nil, err
	}

	// Restore vocabulary relation for response
	updatedStatus.Vocabulary = vocabulary

	// Map to response DTO
	return s.mapStatusToResponse(updatedStatus, userLogin.UUID.String()), nil
}
