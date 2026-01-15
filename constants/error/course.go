package error

import "errors"

var (
	ErrCourseNotFound              = errors.New("course not found")
	ErrCourseDuplicate             = errors.New("course with this title already exists for this JLPT level")
	ErrInvalidJlptLevelIDCourse    = errors.New("invalid JLPT level ID for course")
	ErrInvalidCourseDifficulty     = errors.New("course difficulty must be between 1 and 5")
	ErrInvalidCourseEstimatedHours = errors.New("estimated hours must be a positive number")
	ErrCourseAlreadyPublished      = errors.New("course is already published")
	ErrCourseNotPublished          = errors.New("course is not published")
)

var CourseErrors = []error{
	ErrCourseNotFound,
	ErrCourseDuplicate,
	ErrInvalidJlptLevelIDCourse,
	ErrInvalidCourseDifficulty,
	ErrInvalidCourseEstimatedHours,
	ErrCourseAlreadyPublished,
	ErrCourseNotPublished,
}
