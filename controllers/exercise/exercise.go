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

type ExerciseController struct {
	service services.IServiceRegistry
}

type IExerciseController interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	UpdatePublishStatus(*gin.Context)
	GetByLessonID(*gin.Context)
}

func NewExerciseController(service services.IServiceRegistry) IExerciseController {
	return &ExerciseController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *ExerciseController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrExerciseNotFound:
		return http.StatusNotFound
	case errConstant.ErrDuplicateExerciseOrderIndex:
		return http.StatusConflict
	case errConstant.ErrInvalidLessonIDExercise, errConstant.ErrInvalidExerciseTitle,
		errConstant.ErrInvalidExerciseType, errConstant.ErrInvalidExerciseOrderIndex,
		errConstant.ErrInvalidExerciseDifficulty, errConstant.ErrInvalidExerciseEstimatedMinutes:
		return http.StatusUnprocessableEntity
	case errConstant.ErrExerciseAlreadyPublished, errConstant.ErrExerciseNotPublished:
		return http.StatusBadRequest
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary      Create Exercise
// @Description  Create a new exercise entry for a lesson (admin only)
// @Tags         Exercises
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateExerciseRequest true "Exercise details"
// @Success      201 {object} dto.ExerciseSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response "Exercise with this order_index already exists for this lesson"
// @Failure      422 {object} response.Response "Invalid lesson ID, exercise type, difficulty level, or order_index"
// @Failure      500 {object} response.Response
// @Router       /exercises [post]
func (c *ExerciseController) Create(ctx *gin.Context) {
	request := &dto.CreateExerciseRequest{}

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

	exercise, err := c.service.GetExercise().Create(ctx, request)
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
		Data: exercise,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all Exercises
// @Description  Retrieve exercises with advanced filtering, search, sorting, and pagination
// @Tags         Exercises
// @Produce      json
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        lessonId query int false "Filter by Lesson ID" example(1)
// @Param        exerciseType query string false "Filter by exercise type (multiple_choice, fill_blank, matching, listening, speaking)" example("fill_blank")
// @Param        isPublished query bool false "Filter by publication status" example(true)
// @Param        search query string false "Search in title and description" example("hiragana")
// @Param        sortBy query string false "Sort by field (order_index, title, created_at, difficulty_level)" default(order_index) example("order_index")
// @Param        sortOrder query string false "Sort order (asc, desc)" default(asc) example("asc")
// @Success      200 {object} dto.ExerciseListSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /exercises [get]
func (c *ExerciseController) GetAll(ctx *gin.Context) {
	filter := &dto.ExerciseFilterRequest{}

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

	exercises, err := c.service.GetExercise().GetAll(ctx, filter)
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
		"pagination": exercises.Pagination,
		"status":     "success",
		"data":       exercises.Data,
	})
}

// GetByID godoc
// @Summary      Get Exercise by ID
// @Description  Retrieve a specific exercise entry by ID
// @Tags         Exercises
// @Produce      json
// @Param        id path int true "Exercise ID"
// @Success      200 {object} dto.ExerciseSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response "Exercise not found"
// @Failure      500 {object} response.Response
// @Router       /exercises/{id} [get]
func (c *ExerciseController) GetByID(ctx *gin.Context) {
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

	exercise, err := c.service.GetExercise().GetByID(ctx, uint(id))
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
		Data: exercise,
		Gin:  ctx,
	})
}

// Update godoc
// @Summary      Update Exercise
// @Description  Update an existing exercise entry by ID (admin only)
// @Tags         Exercises
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Exercise ID"
// @Param        request body dto.UpdateExerciseRequest true "Updated exercise details"
// @Success      200 {object} dto.ExerciseSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Exercise not found"
// @Failure      409 {object} response.Response "Exercise with this order_index already exists for this lesson"
// @Failure      422 {object} response.Response "Invalid lesson ID, exercise type, difficulty level, or order_index"
// @Failure      500 {object} response.Response
// @Router       /exercises/{id} [put]
func (c *ExerciseController) Update(ctx *gin.Context) {
	request := &dto.UpdateExerciseRequest{}
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

	exercise, err := c.service.GetExercise().Update(ctx, request, uint(id))
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
		Data: exercise,
		Gin:  ctx,
	})
}

// Delete godoc
// @Summary      Delete Exercise
// @Description  Delete an exercise entry by ID (admin only)
// @Tags         Exercises
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Exercise ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Exercise not found"
// @Failure      500 {object} response.Response
// @Router       /exercises/{id} [delete]
func (c *ExerciseController) Delete(ctx *gin.Context) {
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

	err = c.service.GetExercise().Delete(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	successMessage := "Exercise deleted successfully"
	response.HttpResponse(response.ParamHTTPResp{
		Code:    http.StatusOK,
		Message: &successMessage,
		Gin:     ctx,
	})
}

// UpdatePublishStatus godoc
// @Summary      Publish/Unpublish Exercise
// @Description  Update the publication status of an exercise by ID (admin only) - PATCH request with isPublished boolean
// @Tags         Exercises
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Exercise ID"
// @Param        request body dto.PublishExerciseRequest true "Publish status (true to publish, false to unpublish)"
// @Success      200 {object} dto.ExerciseSwaggerResponse
// @Failure      400 {object} response.Response "Exercise is already in the requested state"
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Exercise not found"
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /exercises/{id}/publish [patch]
func (c *ExerciseController) UpdatePublishStatus(ctx *gin.Context) {
	request := &dto.PublishExerciseRequest{}
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

	exercise, err := c.service.GetExercise().UpdatePublishStatus(ctx, uint(id), *request.IsPublished)
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
		Data: exercise,
		Gin:  ctx,
	})
}

// GetByLessonID godoc
// @Summary      Get Exercises by Lesson ID
// @Description  Retrieve all exercises for a specific lesson, ordered by order_index
// @Tags         Lessons
// @Produce      json
// @Param        id path int true "Lesson ID"
// @Success      200 {object} response.Response{data=[]dto.ExerciseResponse}
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response "Invalid lesson ID"
// @Failure      500 {object} response.Response
// @Router       /lessons/{id}/exercises [get]
func (c *ExerciseController) GetByLessonID(ctx *gin.Context) {
	lessonIDParam := ctx.Param("id")
	if lessonIDParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	lessonID, err := strconv.ParseUint(lessonIDParam, 10, 32)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	exercises, err := c.service.GetExercise().GetByLessonID(ctx, uint(lessonID))
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
		Data: exercises,
		Gin:  ctx,
	})
}
