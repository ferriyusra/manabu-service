package controllers

import (
	errWrap "manabu-service/common/error"
	"manabu-service/common/response"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VocabularyController struct {
	service services.IServiceRegistry
}

type IVocabularyController interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

func NewVocabularyController(service services.IServiceRegistry) IVocabularyController {
	return &VocabularyController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *VocabularyController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrVocabularyNotFound:
		return http.StatusNotFound
	case errConstant.ErrVocabularyDuplicate:
		return http.StatusConflict
	case errConstant.ErrInvalidJlptLevelID, errConstant.ErrInvalidCategoryID, errConstant.ErrInvalidDifficulty, errConstant.ErrInvalidPartOfSpeech:
		return http.StatusUnprocessableEntity
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary      Create Vocabulary
// @Description  Create a new vocabulary entry (admin only)
// @Tags         Vocabularies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateVocabularyRequest true "Vocabulary details"
// @Success      201 {object} response.Response{data=dto.VocabularyResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response "Vocabulary already exists for this JLPT level"
// @Failure      422 {object} response.Response "Invalid JLPT level ID or Category ID"
// @Failure      500 {object} response.Response
// @Router       /vocabularies [post]
func (c *VocabularyController) Create(ctx *gin.Context) {
	request := &dto.CreateVocabularyRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	vocabulary, err := c.service.GetVocabulary().Create(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: vocabulary,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all Vocabularies
// @Description  Retrieve vocabularies with advanced filtering, search, sorting, and pagination
// @Tags         Vocabularies
// @Produce      json
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        jlpt_level_id query int false "Filter by JLPT Level ID" example(5)
// @Param        category_id query int false "Filter by Category ID" example(1)
// @Param        part_of_speech query string false "Filter by part of speech" example("noun")
// @Param        difficulty query int false "Filter by difficulty (1-5)" minimum(1) maximum(5)
// @Param        search query string false "Search in word, reading, or meaning" example("dog")
// @Param        sort_by query string false "Sort by field (word, difficulty, created_at)" default(created_at) example("word")
// @Param        sort_order query string false "Sort order (asc, desc)" default(desc) example("asc")
// @Success      200 {object} response.Response{data=[]dto.VocabularyResponse,pagination=dto.PaginationResponse}
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /vocabularies [get]
func (c *VocabularyController) GetAll(ctx *gin.Context) {
	filter := &dto.VocabularyFilterRequest{}

	if err := ctx.ShouldBindQuery(filter); err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(filter)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	vocabularies, err := c.service.GetVocabulary().GetAll(ctx, filter)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Return response with data array and pagination at same level
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    http.StatusText(http.StatusOK),
		"data":       vocabularies.Data,
		"pagination": vocabularies.Pagination,
	})
}

// GetByID godoc
// @Summary      Get Vocabulary by ID
// @Description  Retrieve a specific vocabulary entry by ID
// @Tags         Vocabularies
// @Produce      json
// @Param        id path int true "Vocabulary ID"
// @Success      200 {object} response.Response{data=dto.VocabularyResponse}
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response "Vocabulary not found"
// @Failure      500 {object} response.Response
// @Router       /vocabularies/{id} [get]
func (c *VocabularyController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	vocabulary, err := c.service.GetVocabulary().GetByID(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: vocabulary,
		Gin:  ctx,
	})
}

// Update godoc
// @Summary      Update Vocabulary
// @Description  Update an existing vocabulary entry by ID (admin only)
// @Tags         Vocabularies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Vocabulary ID"
// @Param        request body dto.UpdateVocabularyRequest true "Updated vocabulary details"
// @Success      200 {object} response.Response{data=dto.VocabularyResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Vocabulary not found"
// @Failure      409 {object} response.Response "Vocabulary already exists for this JLPT level"
// @Failure      422 {object} response.Response "Invalid JLPT level ID or Category ID"
// @Failure      500 {object} response.Response
// @Router       /vocabularies/{id} [put]
func (c *VocabularyController) Update(ctx *gin.Context) {
	request := &dto.UpdateVocabularyRequest{}
	idParam := ctx.Param("id")
	if idParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	err = ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	vocabulary, err := c.service.GetVocabulary().Update(ctx, request, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: vocabulary,
		Gin:  ctx,
	})
}

// Delete godoc
// @Summary      Delete Vocabulary
// @Description  Delete a vocabulary entry by ID (admin only)
// @Tags         Vocabularies
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Vocabulary ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Vocabulary not found"
// @Failure      500 {object} response.Response
// @Router       /vocabularies/{id} [delete]
func (c *VocabularyController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	err = c.service.GetVocabulary().Delete(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	successMessage := "Vocabulary deleted successfully"
	response.HttpResponse(response.ParamHTTPResp{
		Code:    http.StatusOK,
		Message: &successMessage,
		Gin:     ctx,
	})
}
