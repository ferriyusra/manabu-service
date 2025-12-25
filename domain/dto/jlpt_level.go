package dto

type CreateJlptLevelRequest struct {
	Code        string `json:"code" validate:"required" example:"N5"`
	Name        string `json:"name" validate:"required" example:"JLPT N5"`
	Description string `json:"description" validate:"required" example:"Basic level"`
	LevelOrder  int    `json:"levelOrder" validate:"required,min=1,max=5" example:"5"`
}

type UpdateJlptLevelRequest struct {
	Code        string `json:"code" validate:"required" example:"N5"`
	Name        string `json:"name" validate:"required" example:"JLPT N5"`
	Description string `json:"description" validate:"required" example:"Basic level"`
	LevelOrder  int    `json:"levelOrder" validate:"required,min=1,max=5" example:"5"`
}

type JlptLevelResponse struct {
	ID          uint   `json:"id" example:"1"`
	Code        string `json:"code" example:"N5"`
	Name        string `json:"name" example:"JLPT N5"`
	Description string `json:"description" example:"Basic level"`
	LevelOrder  int    `json:"levelOrder" example:"5"`
}
