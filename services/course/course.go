package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"math"
)

type CourseService struct {
	repository repositories.IRepositoryRegistry
}

// ICourseService defines the contract for course business logic operations.
type ICourseService interface {
	// Create validates and creates a new course entry.
	// Validates JLPT level existence and checks for duplicates.
	Create(context.Context, *dto.CreateCourseRequest) (*dto.CourseResponse, error)

	// GetAll retrieves all courses with filtering, sorting, and pagination.
	GetAll(context.Context, *dto.CourseFilterRequest) (*dto.CourseListResponse, error)

	// GetByID retrieves a single course entry by its ID.
	GetByID(context.Context, uint) (*dto.CourseResponse, error)

	// Update validates and updates an existing course entry.
	// Validates JLPT level existence and checks for duplicates.
	Update(context.Context, *dto.UpdateCourseRequest, uint) (*dto.CourseResponse, error)

	// Delete removes a course entry by ID if it exists.
	Delete(context.Context, uint) error

	// Publish marks a course as published with the current timestamp.
	Publish(context.Context, uint) (*dto.CourseResponse, error)

	// Unpublish marks a course as unpublished.
	Unpublish(context.Context, uint) (*dto.CourseResponse, error)

	// GetPublished retrieves only published courses with filtering, sorting, and pagination.
	GetPublished(context.Context, *dto.CourseFilterRequest) (*dto.CourseListResponse, error)
}

func NewCourseService(repository repositories.IRepositoryRegistry) ICourseService {
	return &CourseService{repository: repository}
}

// toCourseResponse converts a Course model to CourseResponse DTO
func (s *CourseService) toCourseResponse(course *models.Course) *dto.CourseResponse {
	response := &dto.CourseResponse{
		ID:             course.ID,
		Title:          course.Title,
		Description:    course.Description,
		JlptLevelID:    course.JlptLevelID,
		ThumbnailURL:   course.ThumbnailURL,
		Difficulty:     course.Difficulty,
		EstimatedHours: course.EstimatedHours,
		IsPublished:    course.IsPublished,
	}

	// Format PublishedAt as string if present
	if course.PublishedAt != nil {
		publishedAtStr := course.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		response.PublishedAt = &publishedAtStr
	}

	if course.JlptLevel.ID > 0 {
		response.JlptLevel = &dto.JlptLevelResponse{
			ID:          course.JlptLevel.ID,
			Code:        course.JlptLevel.Code,
			Name:        course.JlptLevel.Name,
			Description: course.JlptLevel.Description,
			LevelOrder:  course.JlptLevel.LevelOrder,
		}
	}

	return response
}

func (s *CourseService) isCourseExist(ctx context.Context, title string, jlptLevelID uint) bool {
	course, err := s.repository.GetCourse().GetByTitleAndJlptLevel(ctx, title, jlptLevelID)
	if err != nil {
		return false
	}
	return course != nil
}

func (s *CourseService) isJlptLevelExist(ctx context.Context, jlptLevelID uint) bool {
	jlptLevel, err := s.repository.GetJlptLevel().GetByID(ctx, jlptLevelID)
	if err != nil {
		return false
	}
	return jlptLevel != nil
}

func (s *CourseService) validateDifficulty(difficulty int) error {
	if difficulty < 1 || difficulty > 5 {
		return errConstant.ErrInvalidCourseDifficulty
	}
	return nil
}

func (s *CourseService) validateEstimatedHours(hours int) error {
	if hours < 0 {
		return errConstant.ErrInvalidCourseEstimatedHours
	}
	return nil
}

func (s *CourseService) Create(ctx context.Context, req *dto.CreateCourseRequest) (*dto.CourseResponse, error) {
	// Validate JLPT level exists
	if !s.isJlptLevelExist(ctx, req.JlptLevelID) {
		return nil, errConstant.ErrInvalidJlptLevelIDCourse
	}

	// Validate difficulty range
	if req.Difficulty > 0 {
		if err := s.validateDifficulty(req.Difficulty); err != nil {
			return nil, err
		}
	}

	// Validate estimated hours if provided
	if req.EstimatedHours > 0 {
		if err := s.validateEstimatedHours(req.EstimatedHours); err != nil {
			return nil, err
		}
	}

	// Check if course title already exists for this JLPT level
	if s.isCourseExist(ctx, req.Title, req.JlptLevelID) {
		return nil, errConstant.ErrCourseDuplicate
	}

	course, err := s.repository.GetCourse().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.toCourseResponse(course), nil
}

func (s *CourseService) GetAll(ctx context.Context, filter *dto.CourseFilterRequest) (*dto.CourseListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.CourseFilterRequest{
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

	courses, total, err := s.repository.GetCourse().GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CourseResponse, 0, len(courses))
	for _, course := range courses {
		responses = append(responses, *s.toCourseResponse(&course))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.CourseListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *CourseService) GetByID(ctx context.Context, id uint) (*dto.CourseResponse, error) {
	course, err := s.repository.GetCourse().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toCourseResponse(course), nil
}

func (s *CourseService) Update(ctx context.Context, req *dto.UpdateCourseRequest, id uint) (*dto.CourseResponse, error) {
	// Check if course exists
	existingCourse, err := s.repository.GetCourse().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate JLPT level exists
	if !s.isJlptLevelExist(ctx, req.JlptLevelID) {
		return nil, errConstant.ErrInvalidJlptLevelIDCourse
	}

	// Validate difficulty range
	if req.Difficulty > 0 {
		if err := s.validateDifficulty(req.Difficulty); err != nil {
			return nil, err
		}
	}

	// Validate estimated hours if provided
	if req.EstimatedHours > 0 {
		if err := s.validateEstimatedHours(req.EstimatedHours); err != nil {
			return nil, err
		}
	}

	// Check if course title already exists for this JLPT level (excluding current record)
	// Only check if title or JLPT level is being changed
	if existingCourse.Title != req.Title || existingCourse.JlptLevelID != req.JlptLevelID {
		checkCourse, err := s.repository.GetCourse().GetByTitleAndJlptLevel(ctx, req.Title, req.JlptLevelID)
		if err != nil && err != errConstant.ErrCourseNotFound {
			return nil, err // Only propagate real errors, not "not found"
		}
		if checkCourse != nil && checkCourse.ID != id {
			return nil, errConstant.ErrCourseDuplicate
		}
	}

	course, err := s.repository.GetCourse().Update(ctx, req, id)
	if err != nil {
		return nil, err
	}

	return s.toCourseResponse(course), nil
}

func (s *CourseService) Delete(ctx context.Context, id uint) error {
	// Check if course exists
	_, err := s.repository.GetCourse().GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repository.GetCourse().Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CourseService) Publish(ctx context.Context, id uint) (*dto.CourseResponse, error) {
	// Check if course exists
	existingCourse, err := s.repository.GetCourse().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if already published
	if existingCourse.IsPublished {
		return nil, errConstant.ErrCourseAlreadyPublished
	}

	course, err := s.repository.GetCourse().Publish(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toCourseResponse(course), nil
}

func (s *CourseService) Unpublish(ctx context.Context, id uint) (*dto.CourseResponse, error) {
	// Check if course exists
	existingCourse, err := s.repository.GetCourse().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if not published
	if !existingCourse.IsPublished {
		return nil, errConstant.ErrCourseNotPublished
	}

	course, err := s.repository.GetCourse().Unpublish(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toCourseResponse(course), nil
}

func (s *CourseService) GetPublished(ctx context.Context, filter *dto.CourseFilterRequest) (*dto.CourseListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.CourseFilterRequest{
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

	courses, total, err := s.repository.GetCourse().GetPublished(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CourseResponse, 0, len(courses))
	for _, course := range courses {
		responses = append(responses, *s.toCourseResponse(&course))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.CourseListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}
