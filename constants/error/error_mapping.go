package error

func ErrMapping(err error) bool {
	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors[:]...)
	allErrors = append(allErrors, UserErrors[:]...)
	allErrors = append(allErrors, JlptLevelErrors[:]...)
	allErrors = append(allErrors, CategoryErrors[:]...)
	allErrors = append(allErrors, VocabularyErrors[:]...)
	allErrors = append(allErrors, TagErrors[:]...)
	allErrors = append(allErrors, UserVocabularyStatusErrors[:]...)
	allErrors = append(allErrors, CourseErrors[:]...)
	allErrors = append(allErrors, LessonErrors[:]...)
	allErrors = append(allErrors, ExerciseErrors[:]...)
	allErrors = append(allErrors, ExerciseQuestionErrors[:]...)
	allErrors = append(allErrors, UserCourseProgressErrors[:]...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
