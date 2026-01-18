package error

import "errors"

var (
	ErrExerciseQuestionNotFound           = errors.New("exercise question not found")
	ErrInvalidExerciseIDQuestion          = errors.New("invalid exercise ID for question")
	ErrInvalidQuestionText                = errors.New("question text must be between 3 and 1000 characters")
	ErrInvalidCorrectAnswer               = errors.New("correct answer must be between 1 and 500 characters")
	ErrInvalidQuestionPoints              = errors.New("points must be between 1 and 100")
	ErrInvalidQuestionOrderIndex          = errors.New("order_index must be a non-negative integer")
	ErrDuplicateQuestionOrderIndex        = errors.New("a question with this order_index already exists for this exercise")
	ErrExerciseQuestionAlreadyPublished   = errors.New("exercise question is already published")
	ErrExerciseQuestionNotPublished       = errors.New("exercise question is not published")
	ErrInvalidQuestionType                = errors.New("question type must be one of: multiple_choice, fill_blank, matching, listening, speaking")
	ErrInvalidQuestionOptions             = errors.New("options must not exceed 2000 characters")
	ErrInvalidQuestionExplanation         = errors.New("explanation must not exceed 1000 characters")
	ErrInvalidQuestionAudioURL            = errors.New("audio URL must be a valid URL and not exceed 500 characters")
	ErrInvalidQuestionImageURL            = errors.New("image URL must be a valid URL and not exceed 500 characters")
)

var ExerciseQuestionErrors = []error{
	ErrExerciseQuestionNotFound,
	ErrInvalidExerciseIDQuestion,
	ErrInvalidQuestionText,
	ErrInvalidCorrectAnswer,
	ErrInvalidQuestionPoints,
	ErrInvalidQuestionOrderIndex,
	ErrDuplicateQuestionOrderIndex,
	ErrExerciseQuestionAlreadyPublished,
	ErrExerciseQuestionNotPublished,
	ErrInvalidQuestionType,
	ErrInvalidQuestionOptions,
	ErrInvalidQuestionExplanation,
	ErrInvalidQuestionAudioURL,
	ErrInvalidQuestionImageURL,
}
