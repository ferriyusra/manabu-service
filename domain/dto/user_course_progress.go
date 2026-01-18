package dto

type CreateUserCourseProgressRequest struct {
	CourseID uint `json:"courseId" validate:"required,min=1" example:"1"`
}

type UpdateUserCourseProgressRequest struct {
	CompletedLessons int `json:"completedLessons" validate:"required,min=0" example:"5"`
}

type UserCourseProgressResponse struct {
	ID                 string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID             uint            `json:"userId" example:"1"`
	CourseID           uint            `json:"courseId" example:"1"`
	Status             string          `json:"status" example:"in_progress"`
	ProgressPercentage float64         `json:"progressPercentage" example:"50.00"`
	CompletedLessons   int             `json:"completedLessons" example:"5"`
	TotalLessons       int             `json:"totalLessons" example:"10"`
	StartedAt          *string         `json:"startedAt,omitempty" example:"2024-01-15T10:30:00Z"`
	CompletedAt        *string         `json:"completedAt,omitempty" example:"2024-02-15T10:30:00Z"`
	LastAccessedAt     *string         `json:"lastAccessedAt,omitempty" example:"2024-01-20T14:30:00Z"`
	Course             *CourseResponse `json:"course,omitempty"`
}

type UserCourseProgressListResponse struct {
	Data       []UserCourseProgressResponse `json:"data"`
	Pagination PaginationResponse           `json:"pagination"`
}

type UserCourseProgressFilterRequest struct {
	Status   string `form:"status" validate:"omitempty,oneof=not_started in_progress completed" example:"in_progress"`
	CourseID uint   `form:"courseId" validate:"omitempty,min=1" example:"1"`
	SortBy   string `form:"sortBy" validate:"omitempty,oneof=last_accessed_at progress_percentage started_at" example:"last_accessed_at"`
	SortOrder string `form:"sortOrder" validate:"omitempty,oneof=asc desc" example:"desc"`
	PaginationRequest
}

// Swagger response wrappers
type UserCourseProgressSwaggerResponse struct {
	Message string                     `json:"message" example:"User course progress created successfully"`
	Status  string                     `json:"status" example:"success"`
	Data    UserCourseProgressResponse `json:"data"`
}

type UserCourseProgressListSwaggerResponse struct {
	Message    string                       `json:"message" example:"User course progress retrieved successfully"`
	Pagination PaginationResponse           `json:"pagination"`
	Status     string                       `json:"status" example:"success"`
	Data       []UserCourseProgressResponse `json:"data"`
}
