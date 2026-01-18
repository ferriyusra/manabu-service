package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"math"
)

type ExerciseQuestionService struct {
	repository repositories.IRepositoryRegistry
}

// IExerciseQuestionService defines the contract for exercise question business logic operations.
type IExerciseQuestionService interface {
	// Create validates and creates a new exercise question entry.
	// Validates exercise existence and checks for duplicate order_index.
	Create(context.Context, *dto.CreateExerciseQuestionRequest) (*dto.ExerciseQuestionResponse, error)

	// GetAll retrieves all exercise questions with filtering, sorting, and pagination.
	// Returns full response with CorrectAnswer (for admin use).
	GetAll(context.Context, *dto.ExerciseQuestionFilterRequest) (*dto.ExerciseQuestionListResponse, error)

	// GetAllPublic retrieves all exercise questions for public endpoints.
	// Hides CorrectAnswer and Explanation fields.
	GetAllPublic(context.Context, *dto.ExerciseQuestionFilterRequest) (*dto.ExerciseQuestionPublicListResponse, error)

	// GetByID retrieves a single exercise question entry by its ID.
	// Returns full response with CorrectAnswer (for admin use).
	GetByID(context.Context, uint) (*dto.ExerciseQuestionResponse, error)

	// GetByIDPublic retrieves a single exercise question entry by its ID.
	// Hides CorrectAnswer and Explanation fields (for public use).
	GetByIDPublic(context.Context, uint) (*dto.ExerciseQuestionPublicResponse, error)

	// GetByExerciseID retrieves all exercise questions for a specific exercise, ordered by order_index.
	// Returns full response with CorrectAnswer (for admin use).
	GetByExerciseID(context.Context, uint) ([]dto.ExerciseQuestionResponse, error)

	// GetByExerciseIDPublic retrieves all exercise questions for a specific exercise.
	// Hides CorrectAnswer and Explanation fields (for public use).
	GetByExerciseIDPublic(context.Context, uint) ([]dto.ExerciseQuestionPublicResponse, error)

	// Update validates and updates an existing exercise question entry.
	// Validates exercise existence and checks for duplicate order_index.
	Update(context.Context, *dto.UpdateExerciseQuestionRequest, uint) (*dto.ExerciseQuestionResponse, error)

	// Delete removes an exercise question entry by ID if it exists.
	Delete(context.Context, uint) error

	// UpdatePublishStatus marks an exercise question as published or unpublished.
	UpdatePublishStatus(context.Context, uint, bool) (*dto.ExerciseQuestionResponse, error)
}

func NewExerciseQuestionService(repository repositories.IRepositoryRegistry) IExerciseQuestionService {
	return &ExerciseQuestionService{repository: repository}
}

// toExerciseQuestionResponse converts an ExerciseQuestion model to ExerciseQuestionResponse DTO
func (s *ExerciseQuestionService) toExerciseQuestionResponse(question *models.ExerciseQuestion) *dto.ExerciseQuestionResponse {
	response := &dto.ExerciseQuestionResponse{
		ID:            question.ID,
		ExerciseID:    question.ExerciseID,
		QuestionText:  question.QuestionText,
		QuestionType:  question.QuestionType,
		Options:       question.Options,
		CorrectAnswer: question.CorrectAnswer,
		Explanation:   question.Explanation,
		AudioURL:      question.AudioURL,
		ImageURL:      question.ImageURL,
		OrderIndex:    question.OrderIndex,
		Points:        question.Points,
		IsPublished:   question.IsPublished,
	}

	// Format PublishedAt as string if present
	if question.PublishedAt != nil {
		publishedAtStr := question.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		response.PublishedAt = &publishedAtStr
	}

	// Include exercise data if available
	if question.Exercise.ID > 0 {
		exerciseResponse := &dto.ExerciseResponse{
			ID:               question.Exercise.ID,
			LessonID:         question.Exercise.LessonID,
			Title:            question.Exercise.Title,
			Description:      question.Exercise.Description,
			ExerciseType:     question.Exercise.ExerciseType,
			OrderIndex:       question.Exercise.OrderIndex,
			DifficultyLevel:  question.Exercise.DifficultyLevel,
			EstimatedMinutes: question.Exercise.EstimatedMinutes,
			IsPublished:      question.Exercise.IsPublished,
		}

		// Format Exercise PublishedAt as string if present
		if question.Exercise.PublishedAt != nil {
			exercisePublishedAtStr := question.Exercise.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
			exerciseResponse.PublishedAt = &exercisePublishedAtStr
		}

		// Include lesson data if available
		if question.Exercise.Lesson.ID > 0 {
			lessonResponse := &dto.LessonResponse{
				ID:               question.Exercise.Lesson.ID,
				CourseID:         question.Exercise.Lesson.CourseID,
				Title:            question.Exercise.Lesson.Title,
				Content:          question.Exercise.Lesson.Content,
				OrderIndex:       question.Exercise.Lesson.OrderIndex,
				EstimatedMinutes: question.Exercise.Lesson.EstimatedMinutes,
				IsPublished:      question.Exercise.Lesson.IsPublished,
			}

			// Format Lesson PublishedAt as string if present
			if question.Exercise.Lesson.PublishedAt != nil {
				lessonPublishedAtStr := question.Exercise.Lesson.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
				lessonResponse.PublishedAt = &lessonPublishedAtStr
			}

			// Include course data if available
			if question.Exercise.Lesson.Course.ID > 0 {
				courseResponse := &dto.CourseResponse{
					ID:             question.Exercise.Lesson.Course.ID,
					Title:          question.Exercise.Lesson.Course.Title,
					Description:    question.Exercise.Lesson.Course.Description,
					JlptLevelID:    question.Exercise.Lesson.Course.JlptLevelID,
					ThumbnailURL:   question.Exercise.Lesson.Course.ThumbnailURL,
					Difficulty:     question.Exercise.Lesson.Course.Difficulty,
					EstimatedHours: question.Exercise.Lesson.Course.EstimatedHours,
					IsPublished:    question.Exercise.Lesson.Course.IsPublished,
				}

				// Format Course PublishedAt as string if present
				if question.Exercise.Lesson.Course.PublishedAt != nil {
					coursePublishedAtStr := question.Exercise.Lesson.Course.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
					courseResponse.PublishedAt = &coursePublishedAtStr
				}

				// Include JLPT Level data if available
				if question.Exercise.Lesson.Course.JlptLevel.ID > 0 {
					courseResponse.JlptLevel = &dto.JlptLevelResponse{
						ID:          question.Exercise.Lesson.Course.JlptLevel.ID,
						Code:        question.Exercise.Lesson.Course.JlptLevel.Code,
						Name:        question.Exercise.Lesson.Course.JlptLevel.Name,
						Description: question.Exercise.Lesson.Course.JlptLevel.Description,
						LevelOrder:  question.Exercise.Lesson.Course.JlptLevel.LevelOrder,
					}
				}

				lessonResponse.Course = courseResponse
			}

			exerciseResponse.Lesson = lessonResponse
		}

		response.Exercise = exerciseResponse
	}

	return response
}

// toExerciseQuestionPublicResponse converts an ExerciseQuestion model to ExerciseQuestionPublicResponse DTO
// This hides sensitive fields like CorrectAnswer and Explanation for public endpoints
func (s *ExerciseQuestionService) toExerciseQuestionPublicResponse(question *models.ExerciseQuestion) *dto.ExerciseQuestionPublicResponse {
	response := &dto.ExerciseQuestionPublicResponse{
		ID:           question.ID,
		ExerciseID:   question.ExerciseID,
		QuestionText: question.QuestionText,
		QuestionType: question.QuestionType,
		Options:      question.Options,
		AudioURL:     question.AudioURL,
		ImageURL:     question.ImageURL,
		OrderIndex:   question.OrderIndex,
		Points:       question.Points,
		IsPublished:  question.IsPublished,
	}

	// Format PublishedAt as string if present
	if question.PublishedAt != nil {
		publishedAtStr := question.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		response.PublishedAt = &publishedAtStr
	}

	// Include exercise data if available
	if question.Exercise.ID > 0 {
		exerciseResponse := &dto.ExerciseResponse{
			ID:               question.Exercise.ID,
			LessonID:         question.Exercise.LessonID,
			Title:            question.Exercise.Title,
			Description:      question.Exercise.Description,
			ExerciseType:     question.Exercise.ExerciseType,
			OrderIndex:       question.Exercise.OrderIndex,
			DifficultyLevel:  question.Exercise.DifficultyLevel,
			EstimatedMinutes: question.Exercise.EstimatedMinutes,
			IsPublished:      question.Exercise.IsPublished,
		}

		if question.Exercise.PublishedAt != nil {
			exercisePublishedAtStr := question.Exercise.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
			exerciseResponse.PublishedAt = &exercisePublishedAtStr
		}

		response.Exercise = exerciseResponse
	}

	return response
}

func (s *ExerciseQuestionService) isExerciseExist(ctx context.Context, exerciseID uint) bool {
	exercise, err := s.repository.GetExercise().GetByID(ctx, exerciseID)
	if err != nil {
		return false
	}
	return exercise != nil
}

func (s *ExerciseQuestionService) validateOrderIndex(orderIndex int) error {
	if orderIndex < 0 {
		return errConstant.ErrInvalidQuestionOrderIndex
	}
	return nil
}

func (s *ExerciseQuestionService) validatePoints(points int) error {
	if points < 1 || points > 100 {
		return errConstant.ErrInvalidQuestionPoints
	}
	return nil
}

func (s *ExerciseQuestionService) validateQuestionType(questionType string) error {
	validTypes := map[string]bool{
		"multiple_choice": true,
		"fill_blank":      true,
		"matching":        true,
		"listening":       true,
		"speaking":        true,
	}
	if !validTypes[questionType] {
		return errConstant.ErrInvalidQuestionType
	}
	return nil
}

func (s *ExerciseQuestionService) Create(ctx context.Context, req *dto.CreateExerciseQuestionRequest) (*dto.ExerciseQuestionResponse, error) {
	// Validate exercise exists
	if !s.isExerciseExist(ctx, req.ExerciseID) {
		return nil, errConstant.ErrInvalidExerciseIDQuestion
	}

	// Validate order_index
	if err := s.validateOrderIndex(req.OrderIndex); err != nil {
		return nil, err
	}

	// Validate question_type
	if err := s.validateQuestionType(req.QuestionType); err != nil {
		return nil, err
	}

	// Validate points
	if err := s.validatePoints(req.Points); err != nil {
		return nil, err
	}

	// Check if question with same order_index exists for this exercise
	existingQuestion, err := s.repository.GetExerciseQuestion().GetByExerciseIDAndOrderIndex(ctx, req.ExerciseID, req.OrderIndex)
	if err != nil && err != errConstant.ErrExerciseQuestionNotFound {
		return nil, err
	}
	if existingQuestion != nil {
		return nil, errConstant.ErrDuplicateQuestionOrderIndex
	}

	question, err := s.repository.GetExerciseQuestion().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.toExerciseQuestionResponse(question), nil
}

func (s *ExerciseQuestionService) GetAll(ctx context.Context, filter *dto.ExerciseQuestionFilterRequest) (*dto.ExerciseQuestionListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.ExerciseQuestionFilterRequest{
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

	questions, total, err := s.repository.GetExerciseQuestion().GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ExerciseQuestionResponse, 0, len(questions))
	for _, question := range questions {
		responses = append(responses, *s.toExerciseQuestionResponse(&question))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.ExerciseQuestionListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *ExerciseQuestionService) GetByID(ctx context.Context, id uint) (*dto.ExerciseQuestionResponse, error) {
	question, err := s.repository.GetExerciseQuestion().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toExerciseQuestionResponse(question), nil
}

func (s *ExerciseQuestionService) GetByExerciseID(ctx context.Context, exerciseID uint) ([]dto.ExerciseQuestionResponse, error) {
	// Validate exercise exists
	if !s.isExerciseExist(ctx, exerciseID) {
		return nil, errConstant.ErrInvalidExerciseIDQuestion
	}

	questions, err := s.repository.GetExerciseQuestion().GetByExerciseID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ExerciseQuestionResponse, 0, len(questions))
	for _, question := range questions {
		responses = append(responses, *s.toExerciseQuestionResponse(&question))
	}

	return responses, nil
}

func (s *ExerciseQuestionService) GetAllPublic(ctx context.Context, filter *dto.ExerciseQuestionFilterRequest) (*dto.ExerciseQuestionPublicListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.ExerciseQuestionFilterRequest{
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

	questions, total, err := s.repository.GetExerciseQuestion().GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ExerciseQuestionPublicResponse, 0, len(questions))
	for _, question := range questions {
		responses = append(responses, *s.toExerciseQuestionPublicResponse(&question))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.ExerciseQuestionPublicListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *ExerciseQuestionService) GetByIDPublic(ctx context.Context, id uint) (*dto.ExerciseQuestionPublicResponse, error) {
	question, err := s.repository.GetExerciseQuestion().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toExerciseQuestionPublicResponse(question), nil
}

func (s *ExerciseQuestionService) GetByExerciseIDPublic(ctx context.Context, exerciseID uint) ([]dto.ExerciseQuestionPublicResponse, error) {
	// Validate exercise exists
	if !s.isExerciseExist(ctx, exerciseID) {
		return nil, errConstant.ErrInvalidExerciseIDQuestion
	}

	questions, err := s.repository.GetExerciseQuestion().GetByExerciseID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ExerciseQuestionPublicResponse, 0, len(questions))
	for _, question := range questions {
		responses = append(responses, *s.toExerciseQuestionPublicResponse(&question))
	}

	return responses, nil
}

func (s *ExerciseQuestionService) Update(ctx context.Context, req *dto.UpdateExerciseQuestionRequest, id uint) (*dto.ExerciseQuestionResponse, error) {
	// Check if question exists
	existingQuestion, err := s.repository.GetExerciseQuestion().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate exercise exists
	if !s.isExerciseExist(ctx, req.ExerciseID) {
		return nil, errConstant.ErrInvalidExerciseIDQuestion
	}

	// Validate order_index
	if err := s.validateOrderIndex(req.OrderIndex); err != nil {
		return nil, err
	}

	// Validate question_type
	if err := s.validateQuestionType(req.QuestionType); err != nil {
		return nil, err
	}

	// Validate points
	if err := s.validatePoints(req.Points); err != nil {
		return nil, err
	}

	// Check if question with same order_index exists for this exercise (excluding current record)
	// Only check if exercise_id or order_index is being changed
	if existingQuestion.ExerciseID != req.ExerciseID || existingQuestion.OrderIndex != req.OrderIndex {
		checkQuestion, err := s.repository.GetExerciseQuestion().GetByExerciseIDAndOrderIndex(ctx, req.ExerciseID, req.OrderIndex)
		if err != nil && err != errConstant.ErrExerciseQuestionNotFound {
			return nil, err // Only propagate real errors, not "not found"
		}
		if checkQuestion != nil && checkQuestion.ID != id {
			return nil, errConstant.ErrDuplicateQuestionOrderIndex
		}
	}

	question, err := s.repository.GetExerciseQuestion().Update(ctx, req, id)
	if err != nil {
		return nil, err
	}

	return s.toExerciseQuestionResponse(question), nil
}

func (s *ExerciseQuestionService) Delete(ctx context.Context, id uint) error {
	return s.repository.GetExerciseQuestion().Delete(ctx, id)
}

func (s *ExerciseQuestionService) UpdatePublishStatus(ctx context.Context, id uint, isPublished bool) (*dto.ExerciseQuestionResponse, error) {
	// Check if question exists
	existingQuestion, err := s.repository.GetExerciseQuestion().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if already in the requested state
	if isPublished {
		if existingQuestion.IsPublished {
			return nil, errConstant.ErrExerciseQuestionAlreadyPublished
		}
		question, err := s.repository.GetExerciseQuestion().Publish(ctx, id)
		if err != nil {
			return nil, err
		}
		return s.toExerciseQuestionResponse(question), nil
	} else {
		if !existingQuestion.IsPublished {
			return nil, errConstant.ErrExerciseQuestionNotPublished
		}
		question, err := s.repository.GetExerciseQuestion().Unpublish(ctx, id)
		if err != nil {
			return nil, err
		}
		return s.toExerciseQuestionResponse(question), nil
	}
}
