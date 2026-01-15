package controllers

import (
	categoryController "manabu-service/controllers/category"
	courseController "manabu-service/controllers/course"
	jlptLevelController "manabu-service/controllers/jlpt_level"
	tagController "manabu-service/controllers/tag"
	controllers "manabu-service/controllers/user"
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
