package error

import "errors"

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategoryNameExist  = errors.New("category name already exist for this JLPT level")
	ErrInvalidJlptLevelID = errors.New("invalid JLPT level ID")
)

var CategoryErrors = []error{
	ErrCategoryNotFound,
	ErrCategoryNameExist,
	ErrInvalidJlptLevelID,
}
