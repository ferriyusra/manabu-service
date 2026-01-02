package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100" example:"Animals"`
	Description string `json:"description" validate:"omitempty,max=255" example:"Category for animal-related vocabulary"`
	JlptLevelID uint   `json:"jlptLevelId" validate:"required,min=1" example:"1"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100" example:"Animals"`
	Description string `json:"description" validate:"omitempty,max=255" example:"Category for animal-related vocabulary"`
	JlptLevelID uint   `json:"jlptLevelId" validate:"required,min=1" example:"1"`
}

type CategoryResponse struct {
	ID          uint               `json:"id" example:"1"`
	Name        string             `json:"name" example:"Animals"`
	Description string             `json:"description" example:"Category for animal-related vocabulary"`
	JlptLevelID uint               `json:"jlptLevelId" example:"1"`
	JlptLevel   *JlptLevelResponse `json:"jlptLevel,omitempty"`
}

type PaginationRequest struct {
	Page  int `form:"page" validate:"omitempty,min=1" example:"1"`
	Limit int `form:"limit" validate:"omitempty,min=1,max=100" example:"10"`
}

type PaginationResponse struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	TotalPages int   `json:"totalPages" example:"5"`
	TotalItems int64 `json:"totalItems" example:"50"`
}

type CategoryListResponse struct {
	Data       []CategoryResponse `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
