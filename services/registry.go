package services

import (
	"manabu-service/repositories"
	categoryService "manabu-service/services/category"
	jlptLevelService "manabu-service/services/jlpt_level"
	services "manabu-service/services/user"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
}

type IServiceRegistry interface {
	GetUser() services.IUserService
	GetJlptLevel() jlptLevelService.IJlptLevelService
	GetCategory() categoryService.ICategoryService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repository}
}

func (r *Registry) GetUser() services.IUserService {
	return services.NewUserService(r.repository)
}

func (r *Registry) GetJlptLevel() jlptLevelService.IJlptLevelService {
	return jlptLevelService.NewJlptLevelService(r.repository)
}

func (r *Registry) GetCategory() categoryService.ICategoryService {
	return categoryService.NewCategoryService(r.repository)
}
