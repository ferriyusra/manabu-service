package dto

type CreateCourseRequest struct {
	Title          string `json:"title" validate:"required,min=3,max=200" example:"Introduction to Japanese"`
	Description    string `json:"description" validate:"required,min=10" example:"A comprehensive course for beginners learning Japanese language"`
	JlptLevelID    uint   `json:"jlptLevelId" validate:"required,min=1" example:"5"`
	ThumbnailURL   string `json:"thumbnailUrl" validate:"omitempty,url,max=255" example:"https://example.com/images/course-thumbnail.jpg"`
	Difficulty     int    `json:"difficulty" validate:"omitempty,min=1,max=5" example:"1"`
	EstimatedHours int    `json:"estimatedHours" validate:"omitempty,min=1" example:"40"`
}

type UpdateCourseRequest struct {
	Title          string `json:"title" validate:"required,min=3,max=200" example:"Introduction to Japanese"`
	Description    string `json:"description" validate:"required,min=10" example:"A comprehensive course for beginners learning Japanese language"`
	JlptLevelID    uint   `json:"jlptLevelId" validate:"required,min=1" example:"5"`
	ThumbnailURL   string `json:"thumbnailUrl" validate:"omitempty,url,max=255" example:"https://example.com/images/course-thumbnail.jpg"`
	Difficulty     int    `json:"difficulty" validate:"omitempty,min=1,max=5" example:"1"`
	EstimatedHours int    `json:"estimatedHours" validate:"omitempty,min=1" example:"40"`
}

type CourseResponse struct {
	ID             uint               `json:"id" example:"1"`
	Title          string             `json:"title" example:"Introduction to Japanese"`
	Description    string             `json:"description" example:"A comprehensive course for beginners learning Japanese language"`
	JlptLevelID    uint               `json:"jlptLevelId" example:"5"`
	ThumbnailURL   string             `json:"thumbnailUrl" example:"https://example.com/images/course-thumbnail.jpg"`
	Difficulty     int                `json:"difficulty" example:"1"`
	EstimatedHours int                `json:"estimatedHours" example:"40"`
	IsPublished    bool               `json:"isPublished" example:"true"`
	PublishedAt    *string            `json:"publishedAt,omitempty" example:"2024-01-15T10:30:00Z"`
	JlptLevel      *JlptLevelResponse `json:"jlptLevel,omitempty"`
}

type CourseListResponse struct {
	Data       []CourseResponse   `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type CourseFilterRequest struct {
	JlptLevelID uint   `form:"jlptLevelId" validate:"omitempty,min=1" example:"5"`
	Difficulty  int    `form:"difficulty" validate:"omitempty,min=1,max=5" example:"1"`
	IsPublished *bool  `form:"isPublished" validate:"omitempty" example:"true"`
	Search      string `form:"search" validate:"omitempty,max=100" example:"japanese"`
	SortBy      string `form:"sortBy" validate:"omitempty,oneof=title difficulty created_at" example:"title"`
	SortOrder   string `form:"sortOrder" validate:"omitempty,oneof=asc desc" example:"asc"`
	PaginationRequest
}

// Swagger response wrappers (without token field)
type CourseSwaggerResponse struct {
	Message string         `json:"message" example:"Course created successfully"`
	Status  string         `json:"status" example:"success"`
	Data    CourseResponse `json:"data"`
}

type CourseListSwaggerResponse struct {
	Message    string             `json:"message" example:"Courses retrieved successfully"`
	Pagination PaginationResponse `json:"pagination"`
	Status     string             `json:"status" example:"success"`
	Data       []CourseResponse   `json:"data"`
}
