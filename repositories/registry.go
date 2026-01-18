package repositories

import (
	categoryRepo "manabu-service/repositories/category"
	courseRepo "manabu-service/repositories/course"
	exerciseRepo "manabu-service/repositories/exercise"
	exerciseQuestionRepo "manabu-service/repositories/exercise_question"
	jlptLevelRepo "manabu-service/repositories/jlpt_level"
	lessonRepo "manabu-service/repositories/lesson"
	tagRepo "manabu-service/repositories/tag"
	repositories "manabu-service/repositories/user"
	userCourseProgressRepo "manabu-service/repositories/user_course_progress"
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
	GetCourse() courseRepo.ICourseRepository
	GetLesson() lessonRepo.ILessonRepository
	GetExercise() exerciseRepo.IExerciseRepository
	GetExerciseQuestion() exerciseQuestionRepo.IExerciseQuestionRepository
	GetUserCourseProgress() userCourseProgressRepo.IUserCourseProgressRepository
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

func (r *Registry) GetCourse() courseRepo.ICourseRepository {
	return courseRepo.NewCourseRepository(r.db)
}

func (r *Registry) GetLesson() lessonRepo.ILessonRepository {
	return lessonRepo.NewLessonRepository(r.db)
}

func (r *Registry) GetExercise() exerciseRepo.IExerciseRepository {
	return exerciseRepo.NewExerciseRepository(r.db)
}

func (r *Registry) GetExerciseQuestion() exerciseQuestionRepo.IExerciseQuestionRepository {
	return exerciseQuestionRepo.NewExerciseQuestionRepository(r.db)
}

func (r *Registry) GetUserCourseProgress() userCourseProgressRepo.IUserCourseProgressRepository {
	return userCourseProgressRepo.NewUserCourseProgressRepository(r.db)
}
