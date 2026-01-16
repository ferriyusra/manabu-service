package error

import "errors"

var (
	ErrLessonNotFound             = errors.New("lesson not found")
	ErrInvalidCourseIDLesson      = errors.New("invalid course ID for lesson")
	ErrInvalidLessonTitle         = errors.New("lesson title must be between 3 and 255 characters")
	ErrInvalidLessonOrderIndex    = errors.New("order_index must be a non-negative integer")
	ErrInvalidLessonEstimatedTime = errors.New("estimated_minutes must be a non-negative integer")
	ErrLessonAlreadyPublished     = errors.New("lesson is already published")
	ErrLessonNotPublished         = errors.New("lesson is not published")
	ErrDuplicateOrderIndex        = errors.New("a lesson with this order_index already exists for this course")
)

var LessonErrors = []error{
	ErrLessonNotFound,
	ErrInvalidCourseIDLesson,
	ErrInvalidLessonTitle,
	ErrInvalidLessonOrderIndex,
	ErrInvalidLessonEstimatedTime,
	ErrLessonAlreadyPublished,
	ErrLessonNotPublished,
	ErrDuplicateOrderIndex,
}
