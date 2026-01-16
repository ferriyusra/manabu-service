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

type LessonController struct {
	service services.IServiceRegistry
}

type ILessonController interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	Publish(*gin.Context)
	Unpublish(*gin.Context)
	GetByCourseID(*gin.Context)
}

func NewLessonController(service services.IServiceRegistry) ILessonController {
	return &LessonController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *LessonController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrLessonNotFound:
		return http.StatusNotFound
	case errConstant.ErrDuplicateOrderIndex:
		return http.StatusConflict
	case errConstant.ErrInvalidCourseIDLesson, errConstant.ErrInvalidLessonTitle,
		errConstant.ErrInvalidLessonOrderIndex, errConstant.ErrInvalidLessonEstimatedTime:
		return http.StatusUnprocessableEntity
	case errConstant.ErrLessonAlreadyPublished, errConstant.ErrLessonNotPublished:
		return http.StatusBadRequest
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary      Create Lesson
// @Description  Create a new lesson entry for a course (admin only)
// @Tags         Lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateLessonRequest true "Lesson details"
// @Success      201 {object} dto.LessonSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response "Lesson with this order_index already exists for this course"
// @Failure      422 {object} response.Response "Invalid course ID or order_index"
// @Failure      500 {object} response.Response
// @Router       /lessons [post]
func (c *LessonController) Create(ctx *gin.Context) {
	request := &dto.CreateLessonRequest{}

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

	lesson, err := c.service.GetLesson().Create(ctx, request)
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
		Data: lesson,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all Lessons
// @Description  Retrieve lessons with advanced filtering, search, sorting, and pagination
// @Tags         Lessons
// @Produce      json
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        courseId query int false "Filter by Course ID" example(1)
// @Param        isPublished query bool false "Filter by publication status" example(true)
// @Param        search query string false "Search in title" example("hiragana")
// @Param        sortBy query string false "Sort by field (order_index, title, created_at)" default(order_index) example("order_index")
// @Param        sortOrder query string false "Sort order (asc, desc)" default(asc) example("asc")
// @Success      200 {object} dto.LessonListSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /lessons [get]
func (c *LessonController) GetAll(ctx *gin.Context) {
	filter := &dto.LessonFilterRequest{}

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

	lessons, err := c.service.GetLesson().GetAll(ctx, filter)
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
		"pagination": lessons.Pagination,
		"status":     "success",
		"data":       lessons.Data,
	})
}

// GetByID godoc
// @Summary      Get Lesson by ID
// @Description  Retrieve a specific lesson entry by ID
// @Tags         Lessons
// @Produce      json
// @Param        id path int true "Lesson ID"
// @Success      200 {object} dto.LessonSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response
// @Router       /lessons/{id} [get]
func (c *LessonController) GetByID(ctx *gin.Context) {
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

	lesson, err := c.service.GetLesson().GetByID(ctx, uint(id))
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
		Data: lesson,
		Gin:  ctx,
	})
}

// Update godoc
// @Summary      Update Lesson
// @Description  Update an existing lesson entry by ID (admin only)
// @Tags         Lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Lesson ID"
// @Param        request body dto.UpdateLessonRequest true "Updated lesson details"
// @Success      200 {object} dto.LessonSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      409 {object} response.Response "Lesson with this order_index already exists for this course"
// @Failure      422 {object} response.Response "Invalid course ID or order_index"
// @Failure      500 {object} response.Response
// @Router       /lessons/{id} [put]
func (c *LessonController) Update(ctx *gin.Context) {
	request := &dto.UpdateLessonRequest{}
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

	lesson, err := c.service.GetLesson().Update(ctx, request, uint(id))
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
		Data: lesson,
		Gin:  ctx,
	})
}

// Delete godoc
// @Summary      Delete Lesson
// @Description  Delete a lesson entry by ID (admin only)
// @Tags         Lessons
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Lesson ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response
// @Router       /lessons/{id} [delete]
func (c *LessonController) Delete(ctx *gin.Context) {
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

	err = c.service.GetLesson().Delete(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	successMessage := "Lesson deleted successfully"
	response.HttpResponse(response.ParamHTTPResp{
		Code:    http.StatusOK,
		Message: &successMessage,
		Gin:     ctx,
	})
}

// Publish godoc
// @Summary      Publish Lesson
// @Description  Publish a lesson by ID (admin only) - sets is_published=true and published_at=now
// @Tags         Lessons
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Lesson ID"
// @Success      200 {object} dto.LessonSwaggerResponse
// @Failure      400 {object} response.Response "Lesson is already published"
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response
// @Router       /lessons/{id}/publish [post]
func (c *LessonController) Publish(ctx *gin.Context) {
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

	lesson, err := c.service.GetLesson().Publish(ctx, uint(id))
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
		Data: lesson,
		Gin:  ctx,
	})
}

// Unpublish godoc
// @Summary      Unpublish Lesson
// @Description  Unpublish a lesson by ID (admin only) - sets is_published=false
// @Tags         Lessons
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Lesson ID"
// @Success      200 {object} dto.LessonSwaggerResponse
// @Failure      400 {object} response.Response "Lesson is not published"
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response
// @Router       /lessons/{id}/unpublish [post]
func (c *LessonController) Unpublish(ctx *gin.Context) {
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

	lesson, err := c.service.GetLesson().Unpublish(ctx, uint(id))
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
		Data: lesson,
		Gin:  ctx,
	})
}

// GetByCourseID godoc
// @Summary      Get Lessons by Course ID
// @Description  Retrieve all lessons for a specific course, ordered by order_index
// @Tags         Lessons
// @Produce      json
// @Param        id path int true "Course ID"
// @Success      200 {object} response.Response{data=[]dto.LessonResponse}
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response "Invalid course ID"
// @Failure      500 {object} response.Response
// @Router       /courses/{id}/lessons [get]
func (c *LessonController) GetByCourseID(ctx *gin.Context) {
	courseIDParam := ctx.Param("id")
	if courseIDParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	courseID, err := strconv.ParseUint(courseIDParam, 10, 32)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	lessons, err := c.service.GetLesson().GetByCourseID(ctx, uint(courseID))
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
		Data: lessons,
		Gin:  ctx,
	})
}
