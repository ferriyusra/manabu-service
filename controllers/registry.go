package controllers

import (
	categoryController "manabu-service/controllers/category"
	courseController "manabu-service/controllers/course"
	exerciseController "manabu-service/controllers/exercise"
	exerciseQuestionController "manabu-service/controllers/exercise_question"
	jlptLevelController "manabu-service/controllers/jlpt_level"
	lessonController "manabu-service/controllers/lesson"
	tagController "manabu-service/controllers/tag"
	controllers "manabu-service/controllers/user"
	userCourseProgressController "manabu-service/controllers/user_course_progress"
	userVocabStatusController "manabu-service/controllers/user_vocabulary_status"
	vocabularyController "manabu-service/controllers/vocabulary"
	"manabu-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetUserController() controllers.IUserController
	GetJlptLevelController() jlptLevelController.IJlptLevelController
	GetCategoryController() categoryController.ICategoryController
	GetVocabularyController() vocabularyController.IVocabularyController
	GetTagController() tagController.ITagController
	GetUserVocabularyStatusController() userVocabStatusController.IUserVocabularyStatusController
	GetCourseController() courseController.ICourseController
	GetLessonController() lessonController.ILessonController
	GetExerciseController() exerciseController.IExerciseController
	GetExerciseQuestionController() exerciseQuestionController.IExerciseQuestionController
	GetUserCourseProgressController() userCourseProgressController.IUserCourseProgressController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{service: service}
}

func (u *Registry) GetUserController() controllers.IUserController {
	return controllers.NewUserController(u.service)
}

func (u *Registry) GetJlptLevelController() jlptLevelController.IJlptLevelController {
	return jlptLevelController.NewJlptLevelController(u.service)
}

func (u *Registry) GetCategoryController() categoryController.ICategoryController {
	return categoryController.NewCategoryController(u.service)
}

func (u *Registry) GetVocabularyController() vocabularyController.IVocabularyController {
	return vocabularyController.NewVocabularyController(u.service)
}

func (u *Registry) GetTagController() tagController.ITagController {
	return tagController.NewTagController(u.service)
}

func (u *Registry) GetUserVocabularyStatusController() userVocabStatusController.IUserVocabularyStatusController {
	return userVocabStatusController.NewUserVocabularyStatusController(u.service)
}

func (u *Registry) GetCourseController() courseController.ICourseController {
	return courseController.NewCourseController(u.service)
}

func (u *Registry) GetLessonController() lessonController.ILessonController {
	return lessonController.NewLessonController(u.service)
}

func (u *Registry) GetExerciseController() exerciseController.IExerciseController {
	return exerciseController.NewExerciseController(u.service)
}

func (u *Registry) GetExerciseQuestionController() exerciseQuestionController.IExerciseQuestionController {
	return exerciseQuestionController.NewExerciseQuestionController(u.service)
}

func (u *Registry) GetUserCourseProgressController() userCourseProgressController.IUserCourseProgressController {
	return userCourseProgressController.NewUserCourseProgressController(u.service)
}
