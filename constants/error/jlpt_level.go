package error

import "errors"

var (
	ErrJlptLevelNotFound   = errors.New("jlpt level not found")
	ErrJlptLevelCodeExist  = errors.New("jlpt level code already exist")
	ErrJlptLevelOrderExist = errors.New("jlpt level order already exist")
)

var JlptLevelErrors = []error{
	ErrJlptLevelNotFound,
	ErrJlptLevelCodeExist,
	ErrJlptLevelOrderExist,
}
