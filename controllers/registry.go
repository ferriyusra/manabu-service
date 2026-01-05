package controllers

import (
	categoryController "manabu-service/controllers/category"
	jlptLevelController "manabu-service/controllers/jlpt_level"
	tagController "manabu-service/controllers/tag"
	controllers "manabu-service/controllers/user"
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
