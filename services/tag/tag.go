package services

import (
	"context"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	"manabu-service/repositories"
	"math"
	"regexp"
	"strings"
)

type TagService struct {
	repository repositories.IRepositoryRegistry
}

// ITagService defines the contract for tag business logic operations.
type ITagService interface {
	// Create validates and creates a new tag.
	// Validates name uniqueness and color format.
	Create(context.Context, *dto.CreateTagRequest) (*dto.TagResponse, error)

	// GetAll retrieves all tags with search and pagination.
	GetAll(context.Context, *dto.TagFilterRequest) (*dto.TagListResponse, error)

	// GetByID retrieves a single tag by its ID.
	GetByID(context.Context, uint) (*dto.TagResponse, error)

	// GetByName retrieves a single tag by its name (case-insensitive exact match).
	GetByName(context.Context, string) (*dto.TagResponse, error)

	// Update validates and updates an existing tag.
	// Validates name uniqueness and color format.
	Update(context.Context, *dto.UpdateTagRequest, uint) (*dto.TagResponse, error)

	// Delete removes a tag by ID if it exists.
	Delete(context.Context, uint) error
}

func NewTagService(repository repositories.IRepositoryRegistry) ITagService {
	return &TagService{repository: repository}
}

// toTagResponse converts a Tag model to TagResponse DTO
func (s *TagService) toTagResponse(tag *models.Tag) *dto.TagResponse {
	return &dto.TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Description: tag.Description,
		Color:       tag.Color,
	}
}

// isTagExist checks if a tag with the given name already exists (case-insensitive)
func (s *TagService) isTagExist(ctx context.Context, name string) bool {
	tag, err := s.repository.GetTag().GetByName(ctx, name)
	if err != nil {
		return false
	}
	return tag != nil
}

// validateColor validates hex color code format (#RRGGBB)
func (s *TagService) validateColor(color string) error {
	if color == "" {
		return nil // Color is optional
	}

	// Regex pattern for hex color code: #RRGGBB
	hexColorPattern := `^#[0-9A-Fa-f]{6}$`
	matched, err := regexp.MatchString(hexColorPattern, color)
	if err != nil || !matched {
		return errConstant.ErrInvalidColor
	}

	return nil
}

func (s *TagService) Create(ctx context.Context, req *dto.CreateTagRequest) (*dto.TagResponse, error) {
	// Trim and validate name
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, errConstant.ErrInvalidTagName
	}

	// Check if tag with the same name already exists (case-insensitive)
	if s.isTagExist(ctx, req.Name) {
		return nil, errConstant.ErrTagDuplicate
	}

	// Validate color format if provided
	if err := s.validateColor(req.Color); err != nil {
		return nil, err
	}

	tag, err := s.repository.GetTag().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.toTagResponse(tag), nil
}

func (s *TagService) GetAll(ctx context.Context, filter *dto.TagFilterRequest) (*dto.TagListResponse, error) {
	// Set default pagination values
	if filter == nil {
		filter = &dto.TagFilterRequest{
			PaginationRequest: dto.PaginationRequest{
				Page:  1,
				Limit: 10,
			},
		}
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	tags, total, err := s.repository.GetTag().GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		responses = append(responses, *s.toTagResponse(&tag))
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.TagListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}, nil
}

func (s *TagService) GetByID(ctx context.Context, id uint) (*dto.TagResponse, error) {
	tag, err := s.repository.GetTag().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toTagResponse(tag), nil
}

func (s *TagService) GetByName(ctx context.Context, name string) (*dto.TagResponse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errConstant.ErrInvalidTagName
	}

	tag, err := s.repository.GetTag().GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return s.toTagResponse(tag), nil
}

func (s *TagService) Update(ctx context.Context, req *dto.UpdateTagRequest, id uint) (*dto.TagResponse, error) {
	// Check if tag exists
	existingTag, err := s.repository.GetTag().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate and check name uniqueness if name is being updated
	if req.Name != "" {
		req.Name = strings.TrimSpace(req.Name)

		// Check if a different tag with this name already exists (case-insensitive)
		if existingTag.Name != req.Name {
			checkTag, err := s.repository.GetTag().GetByName(ctx, req.Name)
			if err != nil && err != errConstant.ErrTagNotFound {
				return nil, err // Only propagate real errors, not "not found"
			}
			if checkTag != nil && checkTag.ID != id {
				return nil, errConstant.ErrTagDuplicate
			}
		}
	}

	// Validate color format if provided
	if err := s.validateColor(req.Color); err != nil {
		return nil, err
	}

	tag, err := s.repository.GetTag().Update(ctx, req, id)
	if err != nil {
		return nil, err
	}

	return s.toTagResponse(tag), nil
}

func (s *TagService) Delete(ctx context.Context, id uint) error {
	// Check if tag exists
	_, err := s.repository.GetTag().GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repository.GetTag().Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
