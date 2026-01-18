package error

import "errors"

var (
	ErrUserCourseProgressNotFound        = errors.New("user course progress not found")
	ErrUserCourseProgressAlreadyExists   = errors.New("user already enrolled in this course")
	ErrInvalidCourseIDProgress           = errors.New("invalid course ID for progress")
	ErrInvalidUserIDProgress             = errors.New("invalid user ID for progress")
	ErrInvalidProgressStatus             = errors.New("status must be one of: not_started, in_progress, completed")
	ErrInvalidCompletedLessons           = errors.New("completed lessons must be between 0 and total lessons")
	ErrInvalidProgressPercentage         = errors.New("progress percentage must be between 0 and 100")
	ErrCannotUpdateCompletedCourse       = errors.New("cannot update progress for a completed course")
	ErrCompletedLessonsExceedTotal       = errors.New("completed lessons cannot exceed total lessons in the course")
)

var UserCourseProgressErrors = []error{
	ErrUserCourseProgressNotFound,
	ErrUserCourseProgressAlreadyExists,
	ErrInvalidCourseIDProgress,
	ErrInvalidUserIDProgress,
	ErrInvalidProgressStatus,
	ErrInvalidCompletedLessons,
	ErrInvalidProgressPercentage,
	ErrCannotUpdateCompletedCourse,
	ErrCompletedLessonsExceedTotal,
}
