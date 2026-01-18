package dto

type CreateExerciseQuestionRequest struct {
	ExerciseID    uint   `json:"exerciseId" validate:"required,min=1" example:"1"`
	QuestionText  string `json:"questionText" validate:"required,min=3,max=1000" example:"What is the correct Hiragana for 'a'?"`
	QuestionType  string `json:"questionType" validate:"required,oneof=multiple_choice fill_blank matching listening speaking" example:"multiple_choice"`
	Options       string `json:"options" validate:"omitempty,max=2000" example:"{\"a\":\"あ\",\"b\":\"い\",\"c\":\"う\",\"d\":\"え\"}"`
	CorrectAnswer string `json:"correctAnswer" validate:"required,min=1,max=500" example:"あ"`
	Explanation   string `json:"explanation" validate:"omitempty,max=1000" example:"The Hiragana character for 'a' is あ"`
	AudioURL      string `json:"audioUrl" validate:"omitempty,url,max=500" example:"https://example.com/audio/question1.mp3"`
	ImageURL      string `json:"imageUrl" validate:"omitempty,url,max=500" example:"https://example.com/images/question1.jpg"`
	OrderIndex    int    `json:"orderIndex" validate:"required,min=0" example:"1"`
	Points        int    `json:"points" validate:"required,min=1,max=100" example:"10"`
}

type UpdateExerciseQuestionRequest struct {
	ExerciseID    uint   `json:"exerciseId" validate:"required,min=1" example:"1"`
	QuestionText  string `json:"questionText" validate:"required,min=3,max=1000" example:"What is the correct Hiragana for 'a'?"`
	QuestionType  string `json:"questionType" validate:"required,oneof=multiple_choice fill_blank matching listening speaking" example:"multiple_choice"`
	Options       string `json:"options" validate:"omitempty,max=2000" example:"{\"a\":\"あ\",\"b\":\"い\",\"c\":\"う\",\"d\":\"え\"}"`
	CorrectAnswer string `json:"correctAnswer" validate:"required,min=1,max=500" example:"あ"`
	Explanation   string `json:"explanation" validate:"omitempty,max=1000" example:"The Hiragana character for 'a' is あ"`
	AudioURL      string `json:"audioUrl" validate:"omitempty,url,max=500" example:"https://example.com/audio/question1.mp3"`
	ImageURL      string `json:"imageUrl" validate:"omitempty,url,max=500" example:"https://example.com/images/question1.jpg"`
	OrderIndex    int    `json:"orderIndex" validate:"required,min=0" example:"1"`
	Points        int    `json:"points" validate:"required,min=1,max=100" example:"10"`
}

type ExerciseQuestionResponse struct {
	ID            uint              `json:"id" example:"1"`
	ExerciseID    uint              `json:"exerciseId" example:"1"`
	QuestionText  string            `json:"questionText" example:"What is the correct Hiragana for 'a'?"`
	QuestionType  string            `json:"questionType" example:"multiple_choice"`
	Options       string            `json:"options,omitempty" example:"{\"a\":\"あ\",\"b\":\"い\",\"c\":\"う\",\"d\":\"え\"}"`
	CorrectAnswer string            `json:"correctAnswer" example:"あ"`
	Explanation   string            `json:"explanation,omitempty" example:"The Hiragana character for 'a' is あ"`
	AudioURL      string            `json:"audioUrl,omitempty" example:"https://example.com/audio/question1.mp3"`
	ImageURL      string            `json:"imageUrl,omitempty" example:"https://example.com/images/question1.jpg"`
	OrderIndex    int               `json:"orderIndex" example:"1"`
	Points        int               `json:"points" example:"10"`
	IsPublished   bool              `json:"isPublished" example:"true"`
	PublishedAt   *string           `json:"publishedAt,omitempty" example:"2024-01-15T10:30:00Z"`
	Exercise      *ExerciseResponse `json:"exercise,omitempty"`
}

// ExerciseQuestionPublicResponse is used for public endpoints to hide sensitive fields like CorrectAnswer
type ExerciseQuestionPublicResponse struct {
	ID           uint              `json:"id" example:"1"`
	ExerciseID   uint              `json:"exerciseId" example:"1"`
	QuestionText string            `json:"questionText" example:"What is the correct Hiragana for 'a'?"`
	QuestionType string            `json:"questionType" example:"multiple_choice"`
	Options      string            `json:"options,omitempty" example:"{\"a\":\"あ\",\"b\":\"い\",\"c\":\"う\",\"d\":\"え\"}"`
	AudioURL     string            `json:"audioUrl,omitempty" example:"https://example.com/audio/question1.mp3"`
	ImageURL     string            `json:"imageUrl,omitempty" example:"https://example.com/images/question1.jpg"`
	OrderIndex   int               `json:"orderIndex" example:"1"`
	Points       int               `json:"points" example:"10"`
	IsPublished  bool              `json:"isPublished" example:"true"`
	PublishedAt  *string           `json:"publishedAt,omitempty" example:"2024-01-15T10:30:00Z"`
	Exercise     *ExerciseResponse `json:"exercise,omitempty"`
}

type ExerciseQuestionListResponse struct {
	Data       []ExerciseQuestionResponse `json:"data"`
	Pagination PaginationResponse         `json:"pagination"`
}

type ExerciseQuestionPublicListResponse struct {
	Data       []ExerciseQuestionPublicResponse `json:"data"`
	Pagination PaginationResponse               `json:"pagination"`
}

type ExerciseQuestionFilterRequest struct {
	ExerciseID   uint   `form:"exerciseId" validate:"omitempty,min=1" example:"1"`
	QuestionType string `form:"questionType" validate:"omitempty,oneof=multiple_choice fill_blank matching listening speaking" example:"multiple_choice"`
	IsPublished  *bool  `form:"isPublished" validate:"omitempty" example:"true"`
	Search       string `form:"search" validate:"omitempty,max=100" example:"hiragana"`
	SortBy       string `form:"sortBy" validate:"omitempty,oneof=order_index question_text created_at points" example:"order_index"`
	SortOrder    string `form:"sortOrder" validate:"omitempty,oneof=asc desc" example:"asc"`
	PaginationRequest
}

type PublishExerciseQuestionRequest struct {
	IsPublished *bool `json:"isPublished" validate:"required" example:"true"`
}

// Swagger response wrappers
type ExerciseQuestionSwaggerResponse struct {
	Message string                   `json:"message" example:"Exercise question created successfully"`
	Status  string                   `json:"status" example:"success"`
	Data    ExerciseQuestionResponse `json:"data"`
}

type ExerciseQuestionListSwaggerResponse struct {
	Message    string                     `json:"message" example:"Exercise questions retrieved successfully"`
	Pagination PaginationResponse         `json:"pagination"`
	Status     string                     `json:"status" example:"success"`
	Data       []ExerciseQuestionResponse `json:"data"`
}
