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

type ExerciseQuestionController struct {
	service services.IServiceRegistry
}

type IExerciseQuestionController interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	UpdatePublishStatus(*gin.Context)
	GetByExerciseID(*gin.Context)
}

func NewExerciseQuestionController(service services.IServiceRegistry) IExerciseQuestionController {
	return &ExerciseQuestionController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *ExerciseQuestionController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrExerciseQuestionNotFound:
		return http.StatusNotFound
	case errConstant.ErrDuplicateQuestionOrderIndex:
		return http.StatusConflict
	case errConstant.ErrInvalidExerciseIDQuestion, errConstant.ErrInvalidQuestionText,
		errConstant.ErrInvalidCorrectAnswer, errConstant.ErrInvalidQuestionPoints,
		errConstant.ErrInvalidQuestionOrderIndex, errConstant.ErrInvalidQuestionType,
		errConstant.ErrInvalidQuestionOptions, errConstant.ErrInvalidQuestionExplanation,
		errConstant.ErrInvalidQuestionAudioURL, errConstant.ErrInvalidQuestionImageURL:
		return http.StatusUnprocessableEntity
	case errConstant.ErrExerciseQuestionAlreadyPublished, errConstant.ErrExerciseQuestionNotPublished:
		return http.StatusBadRequest
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary      Create Exercise Question
// @Description  Create a new question for an exercise (admin only)
// @Tags         Exercise Questions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateExerciseQuestionRequest true "Question details"
// @Success      201 {object} dto.ExerciseQuestionSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response "Question with this order_index already exists for this exercise"
// @Failure      422 {object} response.Response "Invalid exercise ID, question type, points, or order_index"
// @Failure      500 {object} response.Response
// @Router       /exercise-questions [post]
func (c *ExerciseQuestionController) Create(ctx *gin.Context) {
	request := &dto.CreateExerciseQuestionRequest{}

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

	question, err := c.service.GetExerciseQuestion().Create(ctx, request)
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
		Data: question,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all Exercise Questions (Public)
// @Description  Retrieve exercise questions with advanced filtering, search, sorting, and pagination. Note: CorrectAnswer and Explanation are hidden for security.
// @Tags         Exercise Questions
// @Produce      json
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        exerciseId query int false "Filter by Exercise ID" example(1)
// @Param        questionType query string false "Filter by Question Type (multiple_choice, fill_blank, matching, listening, speaking)" example("multiple_choice")
// @Param        isPublished query bool false "Filter by publication status" example(true)
// @Param        search query string false "Search in question text and explanation" example("hiragana")
// @Param        sortBy query string false "Sort by field (order_index, question_text, created_at, points)" default(order_index) example("order_index")
// @Param        sortOrder query string false "Sort order (asc, desc)" default(asc) example("asc")
// @Success      200 {object} dto.ExerciseQuestionListSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /exercise-questions [get]
func (c *ExerciseQuestionController) GetAll(ctx *gin.Context) {
	filter := &dto.ExerciseQuestionFilterRequest{}

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

	// Use public method to hide CorrectAnswer and Explanation
	questions, err := c.service.GetExerciseQuestion().GetAllPublic(ctx, filter)
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
		"pagination": questions.Pagination,
		"status":     "success",
		"data":       questions.Data,
	})
}

// GetByID godoc
// @Summary      Get Exercise Question by ID (Public)
// @Description  Retrieve a specific exercise question by ID. Note: CorrectAnswer and Explanation are hidden for security.
// @Tags         Exercise Questions
// @Produce      json
// @Param        id path int true "Exercise Question ID"
// @Success      200 {object} dto.ExerciseQuestionSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response "Exercise question not found"
// @Failure      500 {object} response.Response
// @Router       /exercise-questions/{id} [get]
func (c *ExerciseQuestionController) GetByID(ctx *gin.Context) {
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

	// Use public method to hide CorrectAnswer and Explanation
	question, err := c.service.GetExerciseQuestion().GetByIDPublic(ctx, uint(id))
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
		Data: question,
		Gin:  ctx,
	})
}

// Update godoc
// @Summary      Update Exercise Question
// @Description  Update an existing exercise question by ID (admin only)
// @Tags         Exercise Questions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Exercise Question ID"
// @Param        request body dto.UpdateExerciseQuestionRequest true "Updated question details"
// @Success      200 {object} dto.ExerciseQuestionSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Exercise question not found"
// @Failure      409 {object} response.Response "Question with this order_index already exists for this exercise"
// @Failure      422 {object} response.Response "Invalid exercise ID, question type, points, or order_index"
// @Failure      500 {object} response.Response
// @Router       /exercise-questions/{id} [put]
func (c *ExerciseQuestionController) Update(ctx *gin.Context) {
	request := &dto.UpdateExerciseQuestionRequest{}
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

	question, err := c.service.GetExerciseQuestion().Update(ctx, request, uint(id))
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
		Data: question,
		Gin:  ctx,
	})
}

// Delete godoc
// @Summary      Delete Exercise Question
// @Description  Delete an exercise question by ID (admin only)
// @Tags         Exercise Questions
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Exercise Question ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Exercise question not found"
// @Failure      500 {object} response.Response
// @Router       /exercise-questions/{id} [delete]
func (c *ExerciseQuestionController) Delete(ctx *gin.Context) {
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

	err = c.service.GetExerciseQuestion().Delete(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	successMessage := "Exercise question deleted successfully"
	response.HttpResponse(response.ParamHTTPResp{
		Code:    http.StatusOK,
		Message: &successMessage,
		Gin:     ctx,
	})
}

// UpdatePublishStatus godoc
// @Summary      Publish/Unpublish Exercise Question
// @Description  Update the publication status of an exercise question by ID (admin only) - PATCH request with isPublished boolean
// @Tags         Exercise Questions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Exercise Question ID"
// @Param        request body dto.PublishExerciseQuestionRequest true "Publish status (true to publish, false to unpublish)"
// @Success      200 {object} dto.ExerciseQuestionSwaggerResponse
// @Failure      400 {object} response.Response "Question is already in the requested state"
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Exercise question not found"
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /exercise-questions/{id}/publish [patch]
func (c *ExerciseQuestionController) UpdatePublishStatus(ctx *gin.Context) {
	request := &dto.PublishExerciseQuestionRequest{}
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

	question, err := c.service.GetExerciseQuestion().UpdatePublishStatus(ctx, uint(id), *request.IsPublished)
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
		Data: question,
		Gin:  ctx,
	})
}

// GetByExerciseID godoc
// @Summary      Get Questions by Exercise ID (Public)
// @Description  Retrieve all questions for a specific exercise, ordered by order_index. Note: CorrectAnswer and Explanation are hidden for security.
// @Tags         Exercises
// @Produce      json
// @Param        id path int true "Exercise ID"
// @Success      200 {object} response.Response{data=[]dto.ExerciseQuestionPublicResponse}
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response "Invalid exercise ID"
// @Failure      500 {object} response.Response
// @Router       /exercises/{id}/questions [get]
func (c *ExerciseQuestionController) GetByExerciseID(ctx *gin.Context) {
	exerciseIDParam := ctx.Param("id")
	if exerciseIDParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	exerciseID, err := strconv.ParseUint(exerciseIDParam, 10, 32)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	// Use public method to hide CorrectAnswer and Explanation
	questions, err := c.service.GetExerciseQuestion().GetByExerciseIDPublic(ctx, uint(exerciseID))
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
		Data: questions,
		Gin:  ctx,
	})
}
