package repositories

import (
	"context"
	"errors"
	errWrap "manabu-service/common/error"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"

	"gorm.io/gorm"
)

type VocabularyRepository struct {
	db *gorm.DB
}

// IVocabularyRepository defines the contract for vocabulary data access operations.
type IVocabularyRepository interface {
	// Create inserts a new vocabulary entry into the database.
	Create(context.Context, *dto.CreateVocabularyRequest) (*models.Vocabulary, error)

	// GetAll retrieves all vocabularies with optional filtering and pagination.
	// Returns the list of vocabularies and total count.
	GetAll(context.Context, *dto.VocabularyFilterRequest) ([]models.Vocabulary, int64, error)

	// GetByID retrieves a single vocabulary entry by its ID.
	GetByID(context.Context, uint) (*models.Vocabulary, error)

	// GetByWordAndJlptLevel retrieves a vocabulary entry by word and JLPT level.
	// Used for duplicate detection.
	GetByWordAndJlptLevel(context.Context, string, uint) (*models.Vocabulary, error)

	// Update modifies an existing vocabulary entry by ID.
	Update(context.Context, *dto.UpdateVocabularyRequest, uint) (*models.Vocabulary, error)

	// Delete removes a vocabulary entry by ID.
	Delete(context.Context, uint) error
}

func NewVocabularyRepository(db *gorm.DB) IVocabularyRepository {
	return &VocabularyRepository{db: db}
}

func (r *VocabularyRepository) Create(ctx context.Context, req *dto.CreateVocabularyRequest) (*models.Vocabulary, error) {
	// Set default difficulty if not provided
	difficulty := req.Difficulty
	if difficulty == 0 {
		difficulty = 1
	}

	vocabulary := models.Vocabulary{
		Word:                   req.Word,
		Reading:                req.Reading,
		Meaning:                req.Meaning,
		PartOfSpeech:           req.PartOfSpeech,
		JlptLevelID:            req.JlptLevelID,
		CategoryID:             req.CategoryID,
		ExampleSentence:        req.ExampleSentence,
		ExampleSentenceReading: req.ExampleSentenceReading,
		ExampleSentenceMeaning: req.ExampleSentenceMeaning,
		AudioURL:               req.AudioURL,
		ImageURL:               req.ImageURL,
		Difficulty:             difficulty,
	}

	err := r.db.WithContext(ctx).Create(&vocabulary).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Load relationships using Preload for efficiency
	err = r.db.WithContext(ctx).
		Preload("JlptLevel").
		Preload("Category").
		Preload("Category.JlptLevel").
		First(&vocabulary, vocabulary.ID).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &vocabulary, nil
}

func (r *VocabularyRepository) GetAll(ctx context.Context, filter *dto.VocabularyFilterRequest) ([]models.Vocabulary, int64, error) {
	var vocabularies []models.Vocabulary
	var total int64

	// Build base query with filters
	query := r.db.WithContext(ctx).Model(&models.Vocabulary{})

	// Apply filters
	if filter != nil {
		if filter.JlptLevelID > 0 {
			query = query.Where("jlpt_level_id = ?", filter.JlptLevelID)
		}
		if filter.CategoryID > 0 {
			query = query.Where("category_id = ?", filter.CategoryID)
		}
		if filter.PartOfSpeech != "" {
			query = query.Where("part_of_speech = ?", filter.PartOfSpeech)
		}
		if filter.Difficulty > 0 {
			query = query.Where("difficulty = ?", filter.Difficulty)
		}
		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			query = query.Where("word LIKE ? OR reading LIKE ? OR meaning LIKE ?",
				searchPattern, searchPattern, searchPattern)
		}
	}

	// Count total records with filters applied
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Apply sorting
	sortBy := "created_at"
	sortOrder := "DESC"
	if filter != nil {
		if filter.SortBy != "" {
			sortBy = filter.SortBy
		}
		if filter.SortOrder != "" {
			sortOrder = filter.SortOrder
		}
	}

	query = query.Preload("JlptLevel").
		Preload("Category").
		Preload("Category.JlptLevel").
		Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filter != nil && filter.Limit > 0 {
		// Defensive validation: ensure page is at least 1
		page := filter.Page
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * filter.Limit
		// Ensure offset is never negative
		if offset < 0 {
			offset = 0
		}

		query = query.Limit(filter.Limit).Offset(offset)
	}

	err := query.Find(&vocabularies).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return vocabularies, total, nil
}

func (r *VocabularyRepository) GetByID(ctx context.Context, id uint) (*models.Vocabulary, error) {
	var vocabulary models.Vocabulary
	err := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Preload("Category").
		Preload("Category.JlptLevel").
		Where("id = ?", id).
		First(&vocabulary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrVocabularyNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &vocabulary, nil
}

func (r *VocabularyRepository) GetByWordAndJlptLevel(ctx context.Context, word string, jlptLevelID uint) (*models.Vocabulary, error) {
	var vocabulary models.Vocabulary
	err := r.db.WithContext(ctx).
		Where("word = ? AND jlpt_level_id = ?", word, jlptLevelID).
		First(&vocabulary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrVocabularyNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &vocabulary, nil
}

func (r *VocabularyRepository) Update(ctx context.Context, req *dto.UpdateVocabularyRequest, id uint) (*models.Vocabulary, error) {
	// Set default difficulty if not provided
	difficulty := req.Difficulty
	if difficulty == 0 {
		difficulty = 1
	}

	vocabulary := models.Vocabulary{
		Word:                   req.Word,
		Reading:                req.Reading,
		Meaning:                req.Meaning,
		PartOfSpeech:           req.PartOfSpeech,
		JlptLevelID:            req.JlptLevelID,
		CategoryID:             req.CategoryID,
		ExampleSentence:        req.ExampleSentence,
		ExampleSentenceReading: req.ExampleSentenceReading,
		ExampleSentenceMeaning: req.ExampleSentenceMeaning,
		AudioURL:               req.AudioURL,
		ImageURL:               req.ImageURL,
		Difficulty:             difficulty,
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&vocabulary)

	if result.Error != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return nil, errConstant.ErrVocabularyNotFound
	}

	// Fetch the updated record with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("JlptLevel").
		Preload("Category").
		Preload("Category.JlptLevel").
		Where("id = ?", id).
		First(&vocabulary).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &vocabulary, nil
}

func (r *VocabularyRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Vocabulary{})

	if result.Error != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errConstant.ErrVocabularyNotFound
	}

	return nil
}
