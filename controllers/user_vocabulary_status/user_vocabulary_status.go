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

type UserVocabularyStatusController struct {
	service services.IServiceRegistry
}

type IUserVocabularyStatusController interface {
	Create(*gin.Context)
	GetByID(*gin.Context)
	GetAll(*gin.Context)
	GetDueForReview(*gin.Context)
	Review(*gin.Context)
}

func NewUserVocabularyStatusController(service services.IServiceRegistry) IUserVocabularyStatusController {
	return &UserVocabularyStatusController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *UserVocabularyStatusController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrUserVocabStatusNotFound:
		return http.StatusNotFound
	case errConstant.ErrVocabAlreadyLearning:
		return http.StatusConflict
	case errConstant.ErrInvalidVocabularyID, errConstant.ErrInvalidUserVocabStatusID:
		return http.StatusUnprocessableEntity
	case errConstant.ErrVocabularyNotFoundForLearning:
		return http.StatusUnprocessableEntity
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	case errConstant.ErrUnauthorized:
		return http.StatusUnauthorized
	case errConstant.ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary      Start learning a vocabulary
// @Description  Create a new user vocabulary status to start learning a vocabulary word with simple progress tracking
// @Tags         User Vocabulary Status
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateUserVocabStatusRequest true "Vocabulary ID to start learning"
// @Success      201 {object} dto.UserVocabStatusSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response "Vocabulary already being learned by user"
// @Failure      422 {object} response.Response "Invalid vocabulary ID or vocabulary not found"
// @Failure      500 {object} response.Response
// @Router       /user-vocabulary-status [post]
func (c *UserVocabularyStatusController) Create(ctx *gin.Context) {
	request := &dto.CreateUserVocabStatusRequest{}

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

	status, err := c.service.GetUserVocabularyStatus().Create(ctx.Request.Context(), request)
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
		Data: status,
		Gin:  ctx,
	})
}

// GetByID godoc
// @Summary      Get vocabulary learning status
// @Description  Retrieve a specific user vocabulary status by ID including progress tracking data
// @Tags         User Vocabulary Status
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User Vocabulary Status ID"
// @Success      200 {object} dto.UserVocabStatusSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response "Forbidden - status belongs to another user"
// @Failure      404 {object} response.Response "User vocabulary status not found"
// @Failure      500 {object} response.Response
// @Router       /user-vocabulary-status/{id} [get]
func (c *UserVocabularyStatusController) GetByID(ctx *gin.Context) {
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

	status, err := c.service.GetUserVocabularyStatus().GetByID(ctx.Request.Context(), uint(id))
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
		Data: status,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all user's vocabulary learning statuses
// @Description  Retrieve all vocabulary that the authenticated user is learning with pagination and filtering
// @Tags         User Vocabulary Status
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        sort query string false "Sort field" Enums(id, created_at, next_review_date, status) default(next_review_date)
// @Param        order query string false "Sort order" Enums(asc, desc) default(asc)
// @Param        status query string false "Filter by status" Enums(learning, completed)
// @Success      200 {object} dto.UserVocabStatusListSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /user-vocabulary-status [get]
func (c *UserVocabularyStatusController) GetAll(ctx *gin.Context) {
	request := &dto.UserVocabStatusListRequest{}

	err := ctx.ShouldBindQuery(request)
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

	result, err := c.service.GetUserVocabularyStatus().GetAll(ctx.Request.Context(), request)
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
		"message":    http.StatusText(http.StatusOK),
		"pagination": result.Pagination,
		"status":     "success",
		"data":       result.Data,
	})
}

// GetDueForReview godoc
// @Summary      Get vocabularies due for review
// @Description  Retrieve all vocabulary that need to be reviewed (based on next_review_date). Note: In simple system, user decides when to review.
// @Tags         User Vocabulary Status
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.UserVocabStatusDueSwaggerResponse
// @Failure      401 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /user-vocabulary-status/due [get]
func (c *UserVocabularyStatusController) GetDueForReview(ctx *gin.Context) {
	statuses, err := c.service.GetUserVocabularyStatus().GetDueForReview(ctx.Request.Context())
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
		Data: statuses,
		Gin:  ctx,
	})
}

// Review godoc
// @Summary      Review a vocabulary
// @Description  Submit review result with simple progress tracking. If correct, increment repetitions (5 correct = completed). If incorrect, reset to 0. User controls when to review.
// @Tags         User Vocabulary Status
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        vocabulary_id path int true "Vocabulary ID"
// @Param        request body dto.ReviewUserVocabStatusRequest true "Review result (isCorrect: true/false)"
// @Success      200 {object} dto.ReviewUserVocabStatusSwaggerResponse
// @Failure      400 {object} response.Response "Invalid vocabulary ID"
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response "Forbidden - status belongs to another user"
// @Failure      404 {object} response.Response "User vocabulary status not found"
// @Failure      422 {object} response.Response "Validation errors"
// @Failure      500 {object} response.Response
// @Router       /user-vocabulary-status/{vocabulary_id}/review [post]
func (c *UserVocabularyStatusController) Review(ctx *gin.Context) {
	// Parse vocabulary_id from URL parameter
	vocabularyIDParam := ctx.Param("vocabulary_id")
	if vocabularyIDParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidVocabularyID,
			Gin:  ctx,
		})
		return
	}

	vocabularyID, err := strconv.ParseUint(vocabularyIDParam, 10, 32)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidVocabularyID,
			Gin:  ctx,
		})
		return
	}

	// Parse request body
	request := &dto.ReviewUserVocabStatusRequest{}
	err = ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Validate request
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

	// Process review
	status, err := c.service.GetUserVocabularyStatus().Review(ctx.Request.Context(), uint(vocabularyID), request)
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
		Data: status,
		Gin:  ctx,
	})
}
