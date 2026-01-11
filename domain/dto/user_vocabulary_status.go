package dto

import "time"

// CreateUserVocabStatusRequest represents the request to start learning a vocabulary
type CreateUserVocabStatusRequest struct {
	VocabularyID uint `json:"vocabularyId" validate:"required,min=1" example:"1"`
}

// UserVocabStatusResponse represents the response for user vocabulary status
type UserVocabStatusResponse struct {
	ID             uint                `json:"id" example:"1"`
	UserID         string              `json:"userId" example:"1"`
	VocabularyID   uint                `json:"vocabularyId" example:"1"`
	Vocabulary     *VocabularyResponse `json:"vocabulary,omitempty"`
	Status         string              `json:"status" example:"learning"`
	Repetitions    int                 `json:"repetitions" example:"0"`
	LastReviewedAt *time.Time          `json:"lastReviewedAt,omitempty" example:"2024-01-08T10:00:00Z"`
	CreatedAt      time.Time           `json:"createdAt" example:"2024-01-08T10:00:00Z"`
	UpdatedAt      time.Time           `json:"updatedAt" example:"2024-01-08T10:00:00Z"`
}

// UserVocabStatusListRequest represents the request for listing user vocabulary statuses
type UserVocabStatusListRequest struct {
	Page   int    `form:"page" validate:"omitempty,min=1" example:"1"`
	Limit  int    `form:"limit" validate:"omitempty,min=1,max=100" example:"10"`
	Sort   string `form:"sort" validate:"omitempty,oneof=id created_at next_review_date status" example:"next_review_date"`
	Order  string `form:"order" validate:"omitempty,oneof=asc desc" example:"asc"`
	Status string `form:"status" validate:"omitempty,oneof=learning completed" example:"learning"`
}

// UserVocabStatusListResponse represents the response for listing user vocabulary statuses
type UserVocabStatusListResponse struct {
	Data       []UserVocabStatusResponse `json:"data"`
	Pagination PaginationResponse        `json:"pagination"`
}

// UserVocabStatusSwaggerResponse is used for Swagger documentation
type UserVocabStatusSwaggerResponse struct {
	Status  string                  `json:"status" example:"success"`
	Message string                  `json:"message" example:"OK"`
	Data    UserVocabStatusResponse `json:"data"`
}

// UserVocabStatusListSwaggerResponse is used for Swagger documentation
type UserVocabStatusListSwaggerResponse struct {
	Status     string                    `json:"status" example:"success"`
	Message    string                    `json:"message" example:"OK"`
	Data       []UserVocabStatusResponse `json:"data"`
	Pagination PaginationResponse        `json:"pagination"`
}

// UserVocabStatusDueSwaggerResponse is used for Swagger documentation
type UserVocabStatusDueSwaggerResponse struct {
	Status  string                    `json:"status" example:"success"`
	Message string                    `json:"message" example:"OK"`
	Data    []UserVocabStatusResponse `json:"data"`
}

// ReviewUserVocabStatusRequest represents the request to review a vocabulary
type ReviewUserVocabStatusRequest struct {
	IsCorrect bool `json:"isCorrect" validate:"required" example:"true"`
}

// ReviewUserVocabStatusSwaggerResponse is used for Swagger documentation
type ReviewUserVocabStatusSwaggerResponse struct {
	Status  string                  `json:"status" example:"success"`
	Message string                  `json:"message" example:"OK"`
	Data    UserVocabStatusResponse `json:"data"`
}
