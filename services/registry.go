package services

import (
	"manabu-service/repositories"
	categoryService "manabu-service/services/category"
	courseService "manabu-service/services/course"
	jlptLevelService "manabu-service/services/jlpt_level"
	lessonService "manabu-service/services/lesson"
	tagService "manabu-service/services/tag"
	services "manabu-service/services/user"
	userVocabStatusService "manabu-service/services/user_vocabulary_status"
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
	GetUserVocabularyStatus() userVocabStatusService.IUserVocabularyStatusService
	GetCourse() courseService.ICourseService
	GetLesson() lessonService.ILessonService
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

func (r *Registry) GetUserVocabularyStatus() userVocabStatusService.IUserVocabularyStatusService {
	return userVocabStatusService.NewUserVocabularyStatusService(r.repository)
}

func (r *Registry) GetCourse() courseService.ICourseService {
	return courseService.NewCourseService(r.repository)
}

func (r *Registry) GetLesson() lessonService.ILessonService {
	return lessonService.NewLessonService(r.repository)
}
