package controllers

import (
	categoryController "manabu-service/controllers/category"
	jlptLevelController "manabu-service/controllers/jlpt_level"
	controllers "manabu-service/controllers/user"
	"manabu-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetUserController() controllers.IUserController
	GetJlptLevelController() jlptLevelController.IJlptLevelController
	GetCategoryController() categoryController.ICategoryController
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
