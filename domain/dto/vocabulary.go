package dto

type CreateVocabularyRequest struct {
	Word                   string `json:"word" validate:"required,min=1,max=255" example:"犬"`
	Reading                string `json:"reading" validate:"omitempty,max=255" example:"いぬ"`
	Meaning                string `json:"meaning" validate:"required,min=1,max=500" example:"dog"`
	PartOfSpeech           string `json:"partOfSpeech" validate:"omitempty,max=50" example:"noun"`
	JlptLevelID            uint   `json:"jlptLevelId" validate:"required,min=1" example:"5"`
	CategoryID             uint   `json:"categoryId" validate:"required,min=1" example:"1"`
	ExampleSentence        string `json:"exampleSentence" validate:"omitempty" example:"犬が好きです"`
	ExampleSentenceReading string `json:"exampleSentenceReading" validate:"omitempty" example:"いぬがすきです"`
	ExampleSentenceMeaning string `json:"exampleSentenceMeaning" validate:"omitempty" example:"I like dogs"`
	AudioURL               string `json:"audioUrl" validate:"omitempty,url,max=255" example:"https://example.com/audio/inu.mp3"`
	ImageURL               string `json:"imageUrl" validate:"omitempty,url,max=255" example:"https://example.com/images/dog.jpg"`
	Difficulty             int    `json:"difficulty" validate:"omitempty,min=1,max=5" example:"1"`
}

type UpdateVocabularyRequest struct {
	Word                   string `json:"word" validate:"required,min=1,max=255" example:"犬"`
	Reading                string `json:"reading" validate:"omitempty,max=255" example:"いぬ"`
	Meaning                string `json:"meaning" validate:"required,min=1,max=500" example:"dog"`
	PartOfSpeech           string `json:"partOfSpeech" validate:"omitempty,max=50" example:"noun"`
	JlptLevelID            uint   `json:"jlptLevelId" validate:"required,min=1" example:"5"`
	CategoryID             uint   `json:"categoryId" validate:"required,min=1" example:"1"`
	ExampleSentence        string `json:"exampleSentence" validate:"omitempty" example:"犬が好きです"`
	ExampleSentenceReading string `json:"exampleSentenceReading" validate:"omitempty" example:"いぬがすきです"`
	ExampleSentenceMeaning string `json:"exampleSentenceMeaning" validate:"omitempty" example:"I like dogs"`
	AudioURL               string `json:"audioUrl" validate:"omitempty,url,max=255" example:"https://example.com/audio/inu.mp3"`
	ImageURL               string `json:"imageUrl" validate:"omitempty,url,max=255" example:"https://example.com/images/dog.jpg"`
	Difficulty             int    `json:"difficulty" validate:"omitempty,min=1,max=5" example:"1"`
}

type VocabularyResponse struct {
	ID                     uint               `json:"id" example:"1"`
	Word                   string             `json:"word" example:"犬"`
	Reading                string             `json:"reading" example:"いぬ"`
	Meaning                string             `json:"meaning" example:"dog"`
	PartOfSpeech           string             `json:"partOfSpeech" example:"noun"`
	JlptLevelID            uint               `json:"jlptLevelId" example:"5"`
	CategoryID             uint               `json:"categoryId" example:"1"`
	ExampleSentence        string             `json:"exampleSentence" example:"犬が好きです"`
	ExampleSentenceReading string             `json:"exampleSentenceReading" example:"いぬがすきです"`
	ExampleSentenceMeaning string             `json:"exampleSentenceMeaning" example:"I like dogs"`
	AudioURL               string             `json:"audioUrl" example:"https://example.com/audio/inu.mp3"`
	ImageURL               string             `json:"imageUrl" example:"https://example.com/images/dog.jpg"`
	Difficulty             int                `json:"difficulty" example:"1"`
	JlptLevel              *JlptLevelResponse `json:"jlptLevel,omitempty"`
	Category               *CategoryResponse  `json:"category,omitempty"`
}

type VocabularyListResponse struct {
	Data       []VocabularyResponse `json:"data"`
	Pagination PaginationResponse   `json:"pagination"`
}

type VocabularyFilterRequest struct {
	JlptLevelID  uint   `form:"jlpt_level_id" validate:"omitempty,min=1" example:"5"`
	CategoryID   uint   `form:"category_id" validate:"omitempty,min=1" example:"1"`
	PartOfSpeech string `form:"part_of_speech" validate:"omitempty,max=50" example:"noun"`
	Difficulty   int    `form:"difficulty" validate:"omitempty,min=1,max=5" example:"1"`
	Search       string `form:"search" validate:"omitempty,max=100" example:"dog"`
	SortBy       string `form:"sort_by" validate:"omitempty,oneof=word difficulty created_at" example:"word"`
	SortOrder    string `form:"sort_order" validate:"omitempty,oneof=asc desc" example:"asc"`
	PaginationRequest
}

// Swagger response wrappers (without token field)
type VocabularySwaggerResponse struct {
	Status  string             `json:"status" example:"success"`
	Message string             `json:"message" example:"Vocabulary created successfully"`
	Data    VocabularyResponse `json:"data"`
}

type VocabularyListSwaggerResponse struct {
	Status  string                  `json:"status" example:"success"`
	Message string                  `json:"message" example:"Vocabularies retrieved successfully"`
	Data    VocabularyListResponse  `json:"data"`
}
