package repositories

import (
	categoryRepo "manabu-service/repositories/category"
	jlptLevelRepo "manabu-service/repositories/jlpt_level"
	repositories "manabu-service/repositories/user"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() repositories.IUserRepository
	GetJlptLevel() jlptLevelRepo.IJlptLevelRepository
	GetCategory() categoryRepo.ICategoryRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetUser() repositories.IUserRepository {
	return repositories.NewUserRepository(r.db)
}

func (r *Registry) GetJlptLevel() jlptLevelRepo.IJlptLevelRepository {
	return jlptLevelRepo.NewJlptLevelRepository(r.db)
}

func (r *Registry) GetCategory() categoryRepo.ICategoryRepository {
	return categoryRepo.NewCategoryRepository(r.db)
}
