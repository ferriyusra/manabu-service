package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"math"
)

type ExerciseService struct {
	repository repositories.IRepositoryRegistry
}

// IExerciseService defines the contract for exercise business logic operations.
type IExerciseService interface {
	// Create validates and creates a new exercise entry.
	// Validates lesson existence and checks for duplicate order_index.
	Create(context.Context, *dto.CreateExerciseRequest) (*dto.ExerciseResponse, error)

	// GetAll retrieves all exercises with filtering, sorting, and pagination.
	GetAll(context.Context, *dto.ExerciseFilterRequest) (*dto.ExerciseListResponse, error)

	// GetByID retrieves a single exercise entry by its ID.
	GetByID(context.Context, uint) (*dto.ExerciseResponse, error)

	// GetByLessonID retrieves all exercises for a specific lesson, ordered by order_index.
	GetByLessonID(context.Context, uint) ([]dto.ExerciseResponse, error)

	// Update validates and updates an existing exercise entry.
	// Validates lesson existence and checks for duplicate order_index.
	Update(context.Context, *dto.UpdateExerciseRequest, uint) (*dto.ExerciseResponse, error)

	// Delete removes an exercise entry by ID if it exists.
	Delete(context.Context, uint) error

	// UpdatePublishStatus marks an exercise as published or unpublished.
	UpdatePublishStatus(context.Context, uint, bool) (*dto.ExerciseResponse, error)
}

func NewExerciseService(repository repositories.IRepositoryRegistry) IExerciseService {
	return &ExerciseService{repository: repository}
}

// toExerciseResponse converts an Exercise model to ExerciseResponse DTO
func (s *ExerciseService) toExerciseResponse(exercise *models.Exercise) *dto.ExerciseResponse {
	response := &dto.ExerciseResponse{
		ID:               exercise.ID,
		LessonID:         exercise.LessonID,
		Title:            exercise.Title,
		Description:      exercise.Description,
		ExerciseType:     exercise.ExerciseType,
		OrderIndex:       exercise.OrderIndex,
		DifficultyLevel:  exercise.DifficultyLevel,
		EstimatedMinutes: exercise.EstimatedMinutes,
		IsPublished:      exercise.IsPublished,
	}

	// Format PublishedAt as string if present
	if exercise.PublishedAt != nil {
		publishedAtStr := exercise.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		response.PublishedAt = &publishedAtStr
	}

	// Include lesson data if available
	if exercise.Lesson.ID > 0 {
		lessonResponse := &dto.LessonResponse{
			ID:               exercise.Lesson.ID,
			CourseID:         exercise.Lesson.CourseID,
			Title:            exercise.Lesson.Title,
			Content:          exercise.Lesson.Content,
			OrderIndex:       exercise.Lesson.OrderIndex,
			EstimatedMinutes: exercise.Lesson.EstimatedMinutes,
			IsPublished:      exercise.Lesson.IsPublished,
		}

		// Format Lesson PublishedAt as string if present
		if exercise.Lesson.PublishedAt != nil {
			lessonPublishedAtStr := exercise.Lesson.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
			lessonResponse.PublishedAt = &lessonPublishedAtStr
		}

		// Include course data if available
		if exercise.Lesson.Course.ID > 0 {
			courseResponse := &dto.CourseResponse{
				ID:             exercise.Lesson.Course.ID,
				Title:          exercise.Lesson.Course.Title,
				Description:    exercise.Lesson.Course.Description,
				JlptLevelID:    exercise.Lesson.Course.JlptLevelID,
				ThumbnailURL:   exercise.Lesson.Course.ThumbnailURL,
				Difficulty:     exercise.Lesson.Course.Difficulty,
				EstimatedHours: exercise.Lesson.Course.EstimatedHours,
				IsPublished:    exercise.Lesson.Course.IsPublished,
			}

			// Format Course PublishedAt as string if present
			if exercise.Lesson.Course.PublishedAt != nil {
				coursePublishedAtStr := exercise.Lesson.Course.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
				courseResponse.PublishedAt = &coursePublishedAtStr
			}

			// Include JLPT Level data if available
			if exercise.Lesson.Course.JlptLevel.ID > 0 {
				courseResponse.JlptLevel = &dto.JlptLevelResponse{
					ID:          exercise.Lesson.Course.JlptLevel.ID,
					Code:        exercise.Lesson.Course.JlptLevel.Code,
					Name:        exercise.Lesson.Course.JlptLevel.Name,
					Description: exercise.Lesson.Course.JlptLevel.Description,
					LevelOrder:  exercise.Lesson.Course.JlptLevel.LevelOrder,
				}
			}

			lessonResponse.Course = courseResponse
		}

		response.Lesson = lessonResponse
	}

	return response
}

func (s *ExerciseService) isLessonExist(ctx context.Context, lessonID uint) bool {
	lesson, err := s.repository.GetLesson().GetByID(ctx, lessonID)
	if err != nil {
		return false
	}
	return lesson != nil
}

func (s *ExerciseService) validateOrderIndex(orderIndex int) error {
	if orderIndex < 0 {
		return errConstant.ErrInvalidExerciseOrderIndex
	}
	return nil
}

func (s *ExerciseService) validateDifficultyLevel(level int) error {
	if level < 1 || level > 5 {
		return errConstant.ErrInvalidExerciseDifficulty
	}
	return nil
}

func (s *ExerciseService) validateEstimatedMinutes(minutes int) error {
	if minutes < 1 || minutes > 60 {
		return errConstant.ErrInvalidExerciseEstimatedMinutes
	}
	return nil
}

func (s *ExerciseService) validateExerciseType(exerciseType string) error {
	validTypes := map[string]bool{
		"multiple_choice": true,
		"fill_blank":      true,
		"matching":        true,
		"listening":       true,
		"speaking":        true,
	}
	if !validTypes[exerciseType] {
		return errConstant.ErrInvalidExerciseType
	}
	return nil
}

func (s *ExerciseService) Create(ctx context.Context, req *dto.CreateExerciseRequest) (*dto.ExerciseResponse, error) {
	// Validate lesson exists
	if !s.isLessonExist(ctx, req.LessonID) {
		return nil, errConstant.ErrInvalidLessonIDExercise
	}

	// Validate order_index
	if err := s.validateOrderIndex(req.OrderIndex); err != nil {
		return nil, err
	}

	// Validate exercise_type
	if err := s.validateExerciseType(req.ExerciseType); err != nil {
		return nil, err
	}

	// Validate difficulty_level if provided
	if req.DifficultyLevel > 0 {
		if err := s.validateDifficultyLevel(req.DifficultyLevel); err != nil {
			return nil, err
		}
	}

	// Validate estimated_minutes if provided
	if req.EstimatedMinutes > 0 {
		if err := s.validateEstimatedMinutes(req.EstimatedMinutes); err != nil {
			return nil, err
		}
	}

	// Check if exercise with same order_index exists for this lesson
	existingExercise, err := s.repository.GetExercise().GetByLessonIDAndOrderIndex(ctx, req.LessonID, req.OrderIndex)
	if err != nil && err != errConstant.ErrExerciseNotFound {
		return nil, err
	}
	if existingExercise != nil {
		return nil, errConstant.ErrDuplicateExerciseOrderIndex
	}

	exercise, err := s.repository.GetExercise().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.toExerciseResponse(exercise), nil
}

func (s *ExerciseService) GetAll(ctx context.Context, filter *dto.ExerciseFilterRequest) (*dto.ExerciseListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.ExerciseFilterRequest{
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

	exercises, total, err := s.repository.GetExercise().GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ExerciseResponse, 0, len(exercises))
	for _, exercise := range exercises {
		responses = append(responses, *s.toExerciseResponse(&exercise))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.ExerciseListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *ExerciseService) GetByID(ctx context.Context, id uint) (*dto.ExerciseResponse, error) {
	exercise, err := s.repository.GetExercise().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toExerciseResponse(exercise), nil
}

func (s *ExerciseService) GetByLessonID(ctx context.Context, lessonID uint) ([]dto.ExerciseResponse, error) {
	// Validate lesson exists
	if !s.isLessonExist(ctx, lessonID) {
		return nil, errConstant.ErrInvalidLessonIDExercise
	}

	exercises, err := s.repository.GetExercise().GetByLessonID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ExerciseResponse, 0, len(exercises))
	for _, exercise := range exercises {
		responses = append(responses, *s.toExerciseResponse(&exercise))
	}

	return responses, nil
}

func (s *ExerciseService) Update(ctx context.Context, req *dto.UpdateExerciseRequest, id uint) (*dto.ExerciseResponse, error) {
	// Check if exercise exists
	existingExercise, err := s.repository.GetExercise().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate lesson exists
	if !s.isLessonExist(ctx, req.LessonID) {
		return nil, errConstant.ErrInvalidLessonIDExercise
	}

	// Validate order_index
	if err := s.validateOrderIndex(req.OrderIndex); err != nil {
		return nil, err
	}

	// Validate exercise_type
	if err := s.validateExerciseType(req.ExerciseType); err != nil {
		return nil, err
	}

	// Validate difficulty_level if provided
	if req.DifficultyLevel > 0 {
		if err := s.validateDifficultyLevel(req.DifficultyLevel); err != nil {
			return nil, err
		}
	}

	// Validate estimated_minutes if provided
	if req.EstimatedMinutes > 0 {
		if err := s.validateEstimatedMinutes(req.EstimatedMinutes); err != nil {
			return nil, err
		}
	}

	// Check if exercise with same order_index exists for this lesson (excluding current record)
	// Only check if lesson_id or order_index is being changed
	if existingExercise.LessonID != req.LessonID || existingExercise.OrderIndex != req.OrderIndex {
		checkExercise, err := s.repository.GetExercise().GetByLessonIDAndOrderIndex(ctx, req.LessonID, req.OrderIndex)
		if err != nil && err != errConstant.ErrExerciseNotFound {
			return nil, err // Only propagate real errors, not "not found"
		}
		if checkExercise != nil && checkExercise.ID != id {
			return nil, errConstant.ErrDuplicateExerciseOrderIndex
		}
	}

	exercise, err := s.repository.GetExercise().Update(ctx, req, id)
	if err != nil {
		return nil, err
	}

	return s.toExerciseResponse(exercise), nil
}

func (s *ExerciseService) Delete(ctx context.Context, id uint) error {
	// Check if exercise exists
	_, err := s.repository.GetExercise().GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repository.GetExercise().Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExerciseService) UpdatePublishStatus(ctx context.Context, id uint, isPublished bool) (*dto.ExerciseResponse, error) {
	// Check if exercise exists
	existingExercise, err := s.repository.GetExercise().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if already in the requested state
	if isPublished {
		if existingExercise.IsPublished {
			return nil, errConstant.ErrExerciseAlreadyPublished
		}
		exercise, err := s.repository.GetExercise().Publish(ctx, id)
		if err != nil {
			return nil, err
		}
		return s.toExerciseResponse(exercise), nil
	} else {
		if !existingExercise.IsPublished {
			return nil, errConstant.ErrExerciseNotPublished
		}
		exercise, err := s.repository.GetExercise().Unpublish(ctx, id)
		if err != nil {
			return nil, err
		}
		return s.toExerciseResponse(exercise), nil
	}
}
