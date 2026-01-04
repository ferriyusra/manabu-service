package error

import "errors"

var (
	ErrVocabularyNotFound  = errors.New("vocabulary not found")
	ErrVocabularyDuplicate = errors.New("vocabulary word already exists for this JLPT level")
	ErrInvalidCategoryID   = errors.New("invalid category ID")
	ErrInvalidDifficulty   = errors.New("difficulty must be between 1 and 5")
	ErrInvalidPartOfSpeech = errors.New("invalid part of speech")
)

var VocabularyErrors = []error{
	ErrVocabularyNotFound,
	ErrVocabularyDuplicate,
	ErrInvalidCategoryID,
	ErrInvalidDifficulty,
	ErrInvalidPartOfSpeech,
}
