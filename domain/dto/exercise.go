package dto

type CreateExerciseRequest struct {
	LessonID         uint   `json:"lessonId" validate:"required,min=1" example:"1"`
	Title            string `json:"title" validate:"required,min=3,max=200" example:"Fill in the Hiragana"`
	Description      string `json:"description" validate:"omitempty,max=1000" example:"Complete the sentences by filling in the correct Hiragana character"`
	ExerciseType     string `json:"exerciseType" validate:"required,oneof=multiple_choice fill_blank matching listening speaking" example:"fill_blank"`
	OrderIndex       int    `json:"orderIndex" validate:"required,min=0" example:"1"`
	DifficultyLevel  int    `json:"difficultyLevel" validate:"omitempty,min=1,max=5" example:"2"`
	EstimatedMinutes int    `json:"estimatedMinutes" validate:"omitempty,min=1,max=60" example:"10"`
}

type UpdateExerciseRequest struct {
	LessonID         uint   `json:"lessonId" validate:"required,min=1" example:"1"`
	Title            string `json:"title" validate:"required,min=3,max=200" example:"Fill in the Hiragana"`
	Description      string `json:"description" validate:"omitempty,max=1000" example:"Complete the sentences by filling in the correct Hiragana character"`
	ExerciseType     string `json:"exerciseType" validate:"required,oneof=multiple_choice fill_blank matching listening speaking" example:"fill_blank"`
	OrderIndex       int    `json:"orderIndex" validate:"required,min=0" example:"1"`
	DifficultyLevel  int    `json:"difficultyLevel" validate:"omitempty,min=1,max=5" example:"2"`
	EstimatedMinutes int    `json:"estimatedMinutes" validate:"omitempty,min=1,max=60" example:"10"`
}

type ExerciseResponse struct {
	ID               uint            `json:"id" example:"1"`
	LessonID         uint            `json:"lessonId" example:"1"`
	Title            string          `json:"title" example:"Fill in the Hiragana"`
	Description      string          `json:"description" example:"Complete the sentences by filling in the correct Hiragana character"`
	ExerciseType     string          `json:"exerciseType" example:"fill_blank"`
	OrderIndex       int             `json:"orderIndex" example:"1"`
	DifficultyLevel  int             `json:"difficultyLevel" example:"2"`
	EstimatedMinutes int             `json:"estimatedMinutes" example:"10"`
	IsPublished      bool            `json:"isPublished" example:"true"`
	PublishedAt      *string         `json:"publishedAt,omitempty" example:"2024-01-15T10:30:00Z"`
	Lesson           *LessonResponse `json:"lesson,omitempty"`
}

type ExerciseListResponse struct {
	Data       []ExerciseResponse `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type ExerciseFilterRequest struct {
	LessonID     uint   `form:"lessonId" validate:"omitempty,min=1" example:"1"`
	ExerciseType string `form:"exerciseType" validate:"omitempty,oneof=multiple_choice fill_blank matching listening speaking" example:"fill_blank"`
	IsPublished  *bool  `form:"isPublished" validate:"omitempty" example:"true"`
	Search       string `form:"search" validate:"omitempty,max=100" example:"hiragana"`
	SortBy       string `form:"sortBy" validate:"omitempty,oneof=order_index title created_at difficulty_level" example:"order_index"`
	SortOrder    string `form:"sortOrder" validate:"omitempty,oneof=asc desc" example:"asc"`
	PaginationRequest
}

type PublishExerciseRequest struct {
	IsPublished *bool `json:"isPublished" validate:"required" example:"true"`
}

// Swagger response wrappers (without token field)
type ExerciseSwaggerResponse struct {
	Message string           `json:"message" example:"Exercise created successfully"`
	Status  string           `json:"status" example:"success"`
	Data    ExerciseResponse `json:"data"`
}

type ExerciseListSwaggerResponse struct {
	Message    string             `json:"message" example:"Exercises retrieved successfully"`
	Pagination PaginationResponse `json:"pagination"`
	Status     string             `json:"status" example:"success"`
	Data       []ExerciseResponse `json:"data"`
}
