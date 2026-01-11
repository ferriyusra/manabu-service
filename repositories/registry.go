package repositories

import (
	categoryRepo "manabu-service/repositories/category"
	jlptLevelRepo "manabu-service/repositories/jlpt_level"
	tagRepo "manabu-service/repositories/tag"
	repositories "manabu-service/repositories/user"
	userVocabStatusRepo "manabu-service/repositories/user_vocabulary_status"
	vocabularyRepo "manabu-service/repositories/vocabulary"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() repositories.IUserRepository
	GetJlptLevel() jlptLevelRepo.IJlptLevelRepository
	GetCategory() categoryRepo.ICategoryRepository
	GetVocabulary() vocabularyRepo.IVocabularyRepository
	GetTag() tagRepo.ITagRepository
	GetUserVocabularyStatus() userVocabStatusRepo.IUserVocabularyStatusRepository
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

func (r *Registry) GetVocabulary() vocabularyRepo.IVocabularyRepository {
	return vocabularyRepo.NewVocabularyRepository(r.db)
}

func (r *Registry) GetTag() tagRepo.ITagRepository {
	return tagRepo.NewTagRepository(r.db)
}

func (r *Registry) GetUserVocabularyStatus() userVocabStatusRepo.IUserVocabularyStatusRepository {
	return userVocabStatusRepo.NewUserVocabularyStatusRepository(r.db)
}
