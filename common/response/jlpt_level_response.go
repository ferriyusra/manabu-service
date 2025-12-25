package response

// JlptLevelResponse represents response without token field for non-auth endpoints
type JlptLevelResponse struct {
	Status  string      `json:"status"`
	Message any         `json:"message"`
	Data    interface{} `json:"data"`
}
