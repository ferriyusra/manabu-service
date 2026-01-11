package error

import "errors"

var (
	ErrUserVocabStatusNotFound       = errors.New("user vocabulary status not found")
	ErrVocabAlreadyLearning          = errors.New("vocabulary already being learned by user")
	ErrInvalidVocabularyID           = errors.New("invalid vocabulary ID")
	ErrInvalidUserVocabStatusID      = errors.New("invalid user vocabulary status ID")
	ErrVocabularyNotFoundForLearning = errors.New("vocabulary not found, cannot start learning")
)

var UserVocabularyStatusErrors = []error{
	ErrUserVocabStatusNotFound,
	ErrVocabAlreadyLearning,
	ErrInvalidVocabularyID,
	ErrInvalidUserVocabStatusID,
	ErrVocabularyNotFoundForLearning,
}
