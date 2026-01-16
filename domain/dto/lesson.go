package dto

type CreateLessonRequest struct {
	CourseID         uint   `json:"courseId" validate:"required,min=1" example:"1"`
	Title            string `json:"title" validate:"required,min=3,max=255" example:"Introduction to Hiragana"`
	Content          string `json:"content" validate:"omitempty" example:"Learn the basics of Hiragana characters..."`
	OrderIndex       int    `json:"orderIndex" validate:"required,min=0" example:"1"`
	EstimatedMinutes int    `json:"estimatedMinutes" validate:"omitempty,min=0" example:"30"`
}

type UpdateLessonRequest struct {
	CourseID         uint   `json:"courseId" validate:"required,min=1" example:"1"`
	Title            string `json:"title" validate:"required,min=3,max=255" example:"Introduction to Hiragana"`
	Content          string `json:"content" validate:"omitempty" example:"Learn the basics of Hiragana characters..."`
	OrderIndex       int    `json:"orderIndex" validate:"required,min=0" example:"1"`
	EstimatedMinutes int    `json:"estimatedMinutes" validate:"omitempty,min=0" example:"30"`
}

type LessonResponse struct {
	ID               uint            `json:"id" example:"1"`
	CourseID         uint            `json:"courseId" example:"1"`
	Title            string          `json:"title" example:"Introduction to Hiragana"`
	Content          string          `json:"content" example:"Learn the basics of Hiragana characters..."`
	OrderIndex       int             `json:"orderIndex" example:"1"`
	EstimatedMinutes int             `json:"estimatedMinutes" example:"30"`
	IsPublished      bool            `json:"isPublished" example:"true"`
	PublishedAt      *string         `json:"publishedAt,omitempty" example:"2024-01-15T10:30:00Z"`
	Course           *CourseResponse `json:"course,omitempty"`
}

type LessonListResponse struct {
	Data       []LessonResponse   `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type LessonFilterRequest struct {
	CourseID    uint   `form:"courseId" validate:"omitempty,min=1" example:"1"`
	IsPublished *bool  `form:"isPublished" validate:"omitempty" example:"true"`
	Search      string `form:"search" validate:"omitempty,max=100" example:"hiragana"`
	SortBy      string `form:"sortBy" validate:"omitempty,oneof=order_index title created_at" example:"order_index"`
	SortOrder   string `form:"sortOrder" validate:"omitempty,oneof=asc desc" example:"asc"`
	PaginationRequest
}

// Swagger response wrappers (without token field)
type LessonSwaggerResponse struct {
	Message string         `json:"message" example:"Lesson created successfully"`
	Status  string         `json:"status" example:"success"`
	Data    LessonResponse `json:"data"`
}

type LessonListSwaggerResponse struct {
	Message    string             `json:"message" example:"Lessons retrieved successfully"`
	Pagination PaginationResponse `json:"pagination"`
	Status     string             `json:"status" example:"success"`
	Data       []LessonResponse   `json:"data"`
}
