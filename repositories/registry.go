package repositories

import (
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
