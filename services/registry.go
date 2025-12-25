package services

import (
	"manabu-service/repositories"
	jlptLevelService "manabu-service/services/jlpt_level"
	services "manabu-service/services/user"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
}

type IServiceRegistry interface {
	GetUser() services.IUserService
	GetJlptLevel() jlptLevelService.IJlptLevelService
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
