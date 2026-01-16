package error

import "errors"

var (
	ErrExerciseNotFound                = errors.New("exercise not found")
	ErrInvalidLessonIDExercise         = errors.New("invalid lesson ID for exercise")
	ErrInvalidExerciseTitle            = errors.New("exercise title must be between 3 and 200 characters")
	ErrInvalidExerciseType             = errors.New("exercise type must be one of: multiple_choice, fill_blank, matching, listening, speaking")
	ErrInvalidExerciseOrderIndex       = errors.New("order_index must be a non-negative integer")
	ErrInvalidExerciseDifficulty       = errors.New("difficulty_level must be between 1 and 5")
	ErrInvalidExerciseEstimatedMinutes = errors.New("estimated_minutes must be between 1 and 60")
	ErrExerciseAlreadyPublished        = errors.New("exercise is already published")
	ErrExerciseNotPublished            = errors.New("exercise is not published")
	ErrDuplicateExerciseOrderIndex     = errors.New("an exercise with this order_index already exists for this lesson")
)

var ExerciseErrors = []error{
	ErrExerciseNotFound,
	ErrInvalidLessonIDExercise,
	ErrInvalidExerciseTitle,
	ErrInvalidExerciseType,
	ErrInvalidExerciseOrderIndex,
	ErrInvalidExerciseDifficulty,
	ErrInvalidExerciseEstimatedMinutes,
	ErrExerciseAlreadyPublished,
	ErrExerciseNotPublished,
	ErrDuplicateExerciseOrderIndex,
}
