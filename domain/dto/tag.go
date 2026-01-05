package dto

// CreateTagRequest represents the request body for creating a new tag
type CreateTagRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=50" example:"Grammar"`
	Description string `json:"description" validate:"omitempty,max=255" example:"Tags for grammar-related vocabulary"`
	Color       string `json:"color" validate:"omitempty,hexcolor,len=7" example:"#FF5733"`
}

// UpdateTagRequest represents the request body for updating an existing tag
type UpdateTagRequest struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=50" example:"Grammar"`
	Description string `json:"description" validate:"omitempty,max=255" example:"Tags for grammar-related vocabulary"`
	Color       string `json:"color" validate:"omitempty,hexcolor,len=7" example:"#FF5733"`
}

// TagResponse represents the response structure for a single tag
type TagResponse struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Grammar"`
	Description string `json:"description" example:"Tags for grammar-related vocabulary"`
	Color       string `json:"color" example:"#FF5733"`
}

// TagListResponse represents the response structure for a list of tags with pagination
type TagListResponse struct {
	Data       []TagResponse      `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// TagFilterRequest represents query parameters for filtering and searching tags
type TagFilterRequest struct {
	Search string `form:"search" validate:"omitempty,max=100" example:"grammar"`
	PaginationRequest
}

// Swagger response wrappers (without token field)
type TagSwaggerResponse struct {
	Message string      `json:"message" example:"Tag created successfully"`
	Status  string      `json:"status" example:"success"`
	Data    TagResponse `json:"data"`
}

type TagListSwaggerResponse struct {
	Message    string             `json:"message" example:"Tags retrieved successfully"`
	Pagination PaginationResponse `json:"pagination"`
	Status     string             `json:"status" example:"success"`
	Data       []TagResponse      `json:"data"`
}
