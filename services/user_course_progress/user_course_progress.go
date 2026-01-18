package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"math"

	"github.com/google/uuid"
)

type UserCourseProgressService struct {
	repository repositories.IRepositoryRegistry
}

// IUserCourseProgressService defines the contract for user course progress business logic operations.
type IUserCourseProgressService interface {
	// Create validates and creates a new user course progress entry (enrolls user in course).
	// Validates course existence and checks if user is already enrolled.
	Create(context.Context, uint, *dto.CreateUserCourseProgressRequest) (*dto.UserCourseProgressResponse, error)

	// GetAll retrieves all user course progress entries with filtering, sorting, and pagination.
	GetAll(context.Context, uint, *dto.UserCourseProgressFilterRequest) (*dto.UserCourseProgressListResponse, error)

	// GetByID retrieves a single user course progress entry by its UUID.
	// Verifies that the progress record belongs to the specified user.
	GetByID(context.Context, uuid.UUID, uint) (*dto.UserCourseProgressResponse, error)

	// Update validates and updates an existing user course progress entry.
	// Validates completed lessons count and auto-updates status and progress percentage.
	Update(context.Context, *dto.UpdateUserCourseProgressRequest, uuid.UUID, uint) (*dto.UserCourseProgressResponse, error)
}

func NewUserCourseProgressService(repository repositories.IRepositoryRegistry) IUserCourseProgressService {
	return &UserCourseProgressService{repository: repository}
}

// toUserCourseProgressResponse converts a UserCourseProgress model to UserCourseProgressResponse DTO
func (s *UserCourseProgressService) toUserCourseProgressResponse(progress *models.UserCourseProgress) *dto.UserCourseProgressResponse {
	response := &dto.UserCourseProgressResponse{
		ID:                 progress.ID.String(),
		UserID:             progress.UserID,
		CourseID:           progress.CourseID,
		Status:             progress.Status,
		ProgressPercentage: progress.ProgressPercentage,
		CompletedLessons:   progress.CompletedLessons,
		TotalLessons:       progress.TotalLessons,
	}

	// Format timestamps as strings if present
	if progress.StartedAt != nil {
		startedAtStr := progress.StartedAt.Format("2006-01-02T15:04:05Z07:00")
		response.StartedAt = &startedAtStr
	}
	if progress.CompletedAt != nil {
		completedAtStr := progress.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.CompletedAt = &completedAtStr
	}
	if progress.LastAccessedAt != nil {
		lastAccessedAtStr := progress.LastAccessedAt.Format("2006-01-02T15:04:05Z07:00")
		response.LastAccessedAt = &lastAccessedAtStr
	}

	// Include course data if available
	if progress.Course.ID > 0 {
		courseResponse := &dto.CourseResponse{
			ID:             progress.Course.ID,
			Title:          progress.Course.Title,
			Description:    progress.Course.Description,
			JlptLevelID:    progress.Course.JlptLevelID,
			ThumbnailURL:   progress.Course.ThumbnailURL,
			Difficulty:     progress.Course.Difficulty,
			EstimatedHours: progress.Course.EstimatedHours,
			IsPublished:    progress.Course.IsPublished,
		}

		// Format Course PublishedAt as string if present
		if progress.Course.PublishedAt != nil {
			coursePublishedAtStr := progress.Course.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
			courseResponse.PublishedAt = &coursePublishedAtStr
		}

		// Include JLPT Level data if available
		if progress.Course.JlptLevel.ID > 0 {
			courseResponse.JlptLevel = &dto.JlptLevelResponse{
				ID:          progress.Course.JlptLevel.ID,
				Code:        progress.Course.JlptLevel.Code,
				Name:        progress.Course.JlptLevel.Name,
				Description: progress.Course.JlptLevel.Description,
				LevelOrder:  progress.Course.JlptLevel.LevelOrder,
			}
		}

		response.Course = courseResponse
	}

	return response
}

func (s *UserCourseProgressService) isCourseExist(ctx context.Context, courseID uint) bool {
	course, err := s.repository.GetCourse().GetByID(ctx, courseID)
	if err != nil {
		return false
	}
	return course != nil
}

func (s *UserCourseProgressService) validateCompletedLessons(completedLessons, totalLessons int) error {
	if completedLessons < 0 {
		return errConstant.ErrInvalidCompletedLessons
	}
	if completedLessons > totalLessons {
		return errConstant.ErrCompletedLessonsExceedTotal
	}
	return nil
}

func (s *UserCourseProgressService) Create(ctx context.Context, userID uint, req *dto.CreateUserCourseProgressRequest) (*dto.UserCourseProgressResponse, error) {
	// Validate course exists
	if !s.isCourseExist(ctx, req.CourseID) {
		return nil, errConstant.ErrInvalidCourseIDProgress
	}

	// Check if user is already enrolled in this course
	existingProgress, err := s.repository.GetUserCourseProgress().GetByUserIDAndCourseID(ctx, userID, req.CourseID)
	if err != nil && err != errConstant.ErrUserCourseProgressNotFound {
		return nil, err
	}
	if existingProgress != nil {
		return nil, errConstant.ErrUserCourseProgressAlreadyExists
	}

	progress, err := s.repository.GetUserCourseProgress().Create(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return s.toUserCourseProgressResponse(progress), nil
}

func (s *UserCourseProgressService) GetAll(ctx context.Context, userID uint, filter *dto.UserCourseProgressFilterRequest) (*dto.UserCourseProgressListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.UserCourseProgressFilterRequest{
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

	// Validate course existence when filtering by course ID
	if filter.CourseID > 0 {
		course, err := s.repository.GetCourse().GetByID(ctx, filter.CourseID)
		if err != nil || course == nil {
			return nil, errConstant.ErrInvalidCourseIDProgress
		}
	}

	progressList, total, err := s.repository.GetUserCourseProgress().GetAll(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserCourseProgressResponse, 0, len(progressList))
	for _, progress := range progressList {
		responses = append(responses, *s.toUserCourseProgressResponse(&progress))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.UserCourseProgressListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *UserCourseProgressService) GetByID(ctx context.Context, id uuid.UUID, userID uint) (*dto.UserCourseProgressResponse, error) {
	progress, err := s.repository.GetUserCourseProgress().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership - return not found error if the progress doesn't belong to the user
	if progress.UserID != userID {
		return nil, errConstant.ErrUserCourseProgressNotFound
	}

	return s.toUserCourseProgressResponse(progress), nil
}

func (s *UserCourseProgressService) Update(ctx context.Context, req *dto.UpdateUserCourseProgressRequest, id uuid.UUID, userID uint) (*dto.UserCourseProgressResponse, error) {
	// Check if progress exists and belongs to user
	existingProgress, err := s.repository.GetUserCourseProgress().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify that the progress belongs to the requesting user
	if existingProgress.UserID != userID {
		return nil, errConstant.ErrUserCourseProgressNotFound
	}

	// Don't allow updates if course is already completed
	if existingProgress.Status == models.ProgressStatusCompleted {
		return nil, errConstant.ErrCannotUpdateCompletedCourse
	}

	// Validate completed lessons
	if err := s.validateCompletedLessons(req.CompletedLessons, existingProgress.TotalLessons); err != nil {
		return nil, err
	}

	progress, err := s.repository.GetUserCourseProgress().Update(ctx, req, id, existingProgress.TotalLessons)
	if err != nil {
		return nil, err
	}

	return s.toUserCourseProgressResponse(progress), nil
}
