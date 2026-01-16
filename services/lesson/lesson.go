package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"math"
)

type LessonService struct {
	repository repositories.IRepositoryRegistry
}

// ILessonService defines the contract for lesson business logic operations.
type ILessonService interface {
	// Create validates and creates a new lesson entry.
	// Validates course existence and checks for duplicate order_index.
	Create(context.Context, *dto.CreateLessonRequest) (*dto.LessonResponse, error)

	// GetAll retrieves all lessons with filtering, sorting, and pagination.
	GetAll(context.Context, *dto.LessonFilterRequest) (*dto.LessonListResponse, error)

	// GetByID retrieves a single lesson entry by its ID.
	GetByID(context.Context, uint) (*dto.LessonResponse, error)

	// GetByCourseID retrieves all lessons for a specific course, ordered by order_index.
	GetByCourseID(context.Context, uint) ([]dto.LessonResponse, error)

	// Update validates and updates an existing lesson entry.
	// Validates course existence and checks for duplicate order_index.
	Update(context.Context, *dto.UpdateLessonRequest, uint) (*dto.LessonResponse, error)

	// Delete removes a lesson entry by ID if it exists.
	Delete(context.Context, uint) error

	// Publish marks a lesson as published with the current timestamp.
	Publish(context.Context, uint) (*dto.LessonResponse, error)

	// Unpublish marks a lesson as unpublished.
	Unpublish(context.Context, uint) (*dto.LessonResponse, error)
}

func NewLessonService(repository repositories.IRepositoryRegistry) ILessonService {
	return &LessonService{repository: repository}
}

// toLessonResponse converts a Lesson model to LessonResponse DTO
func (s *LessonService) toLessonResponse(lesson *models.Lesson) *dto.LessonResponse {
	response := &dto.LessonResponse{
		ID:               lesson.ID,
		CourseID:         lesson.CourseID,
		Title:            lesson.Title,
		Content:          lesson.Content,
		OrderIndex:       lesson.OrderIndex,
		EstimatedMinutes: lesson.EstimatedMinutes,
		IsPublished:      lesson.IsPublished,
	}

	// Format PublishedAt as string if present
	if lesson.PublishedAt != nil {
		publishedAtStr := lesson.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		response.PublishedAt = &publishedAtStr
	}

	// Include course data if available
	if lesson.Course.ID > 0 {
		courseResponse := &dto.CourseResponse{
			ID:             lesson.Course.ID,
			Title:          lesson.Course.Title,
			Description:    lesson.Course.Description,
			JlptLevelID:    lesson.Course.JlptLevelID,
			ThumbnailURL:   lesson.Course.ThumbnailURL,
			Difficulty:     lesson.Course.Difficulty,
			EstimatedHours: lesson.Course.EstimatedHours,
			IsPublished:    lesson.Course.IsPublished,
		}

		// Format Course PublishedAt as string if present
		if lesson.Course.PublishedAt != nil {
			coursePublishedAtStr := lesson.Course.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
			courseResponse.PublishedAt = &coursePublishedAtStr
		}

		// Include JLPT Level data if available
		if lesson.Course.JlptLevel.ID > 0 {
			courseResponse.JlptLevel = &dto.JlptLevelResponse{
				ID:          lesson.Course.JlptLevel.ID,
				Code:        lesson.Course.JlptLevel.Code,
				Name:        lesson.Course.JlptLevel.Name,
				Description: lesson.Course.JlptLevel.Description,
				LevelOrder:  lesson.Course.JlptLevel.LevelOrder,
			}
		}

		response.Course = courseResponse
	}

	return response
}

func (s *LessonService) isCourseExist(ctx context.Context, courseID uint) bool {
	course, err := s.repository.GetCourse().GetByID(ctx, courseID)
	if err != nil {
		return false
	}
	return course != nil
}

func (s *LessonService) validateOrderIndex(orderIndex int) error {
	if orderIndex < 0 {
		return errConstant.ErrInvalidLessonOrderIndex
	}
	return nil
}

func (s *LessonService) validateEstimatedMinutes(minutes int) error {
	if minutes < 0 {
		return errConstant.ErrInvalidLessonEstimatedTime
	}
	return nil
}

func (s *LessonService) Create(ctx context.Context, req *dto.CreateLessonRequest) (*dto.LessonResponse, error) {
	// Validate course exists
	if !s.isCourseExist(ctx, req.CourseID) {
		return nil, errConstant.ErrInvalidCourseIDLesson
	}

	// Validate order_index
	if err := s.validateOrderIndex(req.OrderIndex); err != nil {
		return nil, err
	}

	// Validate estimated_minutes
	if err := s.validateEstimatedMinutes(req.EstimatedMinutes); err != nil {
		return nil, err
	}

	// Check if lesson with same order_index exists for this course
	existingLesson, err := s.repository.GetLesson().GetByCourseIDAndOrderIndex(ctx, req.CourseID, req.OrderIndex)
	if err != nil && err != errConstant.ErrLessonNotFound {
		return nil, err
	}
	if existingLesson != nil {
		return nil, errConstant.ErrDuplicateOrderIndex
	}

	lesson, err := s.repository.GetLesson().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.toLessonResponse(lesson), nil
}

func (s *LessonService) GetAll(ctx context.Context, filter *dto.LessonFilterRequest) (*dto.LessonListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.LessonFilterRequest{
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

	lessons, total, err := s.repository.GetLesson().GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.LessonResponse, 0, len(lessons))
	for _, lesson := range lessons {
		responses = append(responses, *s.toLessonResponse(&lesson))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.LessonListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *LessonService) GetByID(ctx context.Context, id uint) (*dto.LessonResponse, error) {
	lesson, err := s.repository.GetLesson().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toLessonResponse(lesson), nil
}

func (s *LessonService) GetByCourseID(ctx context.Context, courseID uint) ([]dto.LessonResponse, error) {
	// Validate course exists
	if !s.isCourseExist(ctx, courseID) {
		return nil, errConstant.ErrInvalidCourseIDLesson
	}

	lessons, err := s.repository.GetLesson().GetByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.LessonResponse, 0, len(lessons))
	for _, lesson := range lessons {
		responses = append(responses, *s.toLessonResponse(&lesson))
	}

	return responses, nil
}

func (s *LessonService) Update(ctx context.Context, req *dto.UpdateLessonRequest, id uint) (*dto.LessonResponse, error) {
	// Check if lesson exists
	existingLesson, err := s.repository.GetLesson().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate course exists
	if !s.isCourseExist(ctx, req.CourseID) {
		return nil, errConstant.ErrInvalidCourseIDLesson
	}

	// Validate order_index
	if err := s.validateOrderIndex(req.OrderIndex); err != nil {
		return nil, err
	}

	// Validate estimated_minutes
	if err := s.validateEstimatedMinutes(req.EstimatedMinutes); err != nil {
		return nil, err
	}

	// Check if lesson with same order_index exists for this course (excluding current record)
	// Only check if course_id or order_index is being changed
	if existingLesson.CourseID != req.CourseID || existingLesson.OrderIndex != req.OrderIndex {
		checkLesson, err := s.repository.GetLesson().GetByCourseIDAndOrderIndex(ctx, req.CourseID, req.OrderIndex)
		if err != nil && err != errConstant.ErrLessonNotFound {
			return nil, err // Only propagate real errors, not "not found"
		}
		if checkLesson != nil && checkLesson.ID != id {
			return nil, errConstant.ErrDuplicateOrderIndex
		}
	}

	lesson, err := s.repository.GetLesson().Update(ctx, req, id)
	if err != nil {
		return nil, err
	}

	return s.toLessonResponse(lesson), nil
}

func (s *LessonService) Delete(ctx context.Context, id uint) error {
	// Check if lesson exists
	_, err := s.repository.GetLesson().GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repository.GetLesson().Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *LessonService) Publish(ctx context.Context, id uint) (*dto.LessonResponse, error) {
	// Check if lesson exists
	existingLesson, err := s.repository.GetLesson().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if already published
	if existingLesson.IsPublished {
		return nil, errConstant.ErrLessonAlreadyPublished
	}

	lesson, err := s.repository.GetLesson().Publish(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toLessonResponse(lesson), nil
}

func (s *LessonService) Unpublish(ctx context.Context, id uint) (*dto.LessonResponse, error) {
	// Check if lesson exists
	existingLesson, err := s.repository.GetLesson().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if not published
	if !existingLesson.IsPublished {
		return nil, errConstant.ErrLessonNotPublished
	}

	lesson, err := s.repository.GetLesson().Unpublish(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toLessonResponse(lesson), nil
}
