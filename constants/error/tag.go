package error

import "errors"

var (
	ErrTagNotFound    = errors.New("tag not found")
	ErrTagDuplicate   = errors.New("tag with this name already exists")
	ErrInvalidColor   = errors.New("color must be a valid hex color code (e.g., #FF5733)")
	ErrInvalidTagName = errors.New("tag name is required and cannot be empty")
)

var TagErrors = []error{
	ErrTagNotFound,
	ErrTagDuplicate,
	ErrInvalidColor,
	ErrInvalidTagName,
}
