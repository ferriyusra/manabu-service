package services

import (
	"manabu-service/repositories"
	categoryService "manabu-service/services/category"
	jlptLevelService "manabu-service/services/jlpt_level"
	tagService "manabu-service/services/tag"
	services "manabu-service/services/user"
	vocabularyService "manabu-service/services/vocabulary"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
}

type IServiceRegistry interface {
	GetUser() services.IUserService
	GetJlptLevel() jlptLevelService.IJlptLevelService
	GetCategory() categoryService.ICategoryService
	GetVocabulary() vocabularyService.IVocabularyService
	GetTag() tagService.ITagService
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

func (r *Registry) GetVocabulary() vocabularyService.IVocabularyService {
	return vocabularyService.NewVocabularyService(r.repository)
}

func (r *Registry) GetTag() tagService.ITagService {
	return tagService.NewTagService(r.repository)
}
