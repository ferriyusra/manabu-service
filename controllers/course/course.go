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

type CourseController struct {
	service services.IServiceRegistry
}

type ICourseController interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	Publish(*gin.Context)
	Unpublish(*gin.Context)
	GetPublished(*gin.Context)
}

func NewCourseController(service services.IServiceRegistry) ICourseController {
	return &CourseController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *CourseController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrCourseNotFound:
		return http.StatusNotFound
	case errConstant.ErrCourseDuplicate:
		return http.StatusConflict
	case errConstant.ErrInvalidJlptLevelIDCourse, errConstant.ErrInvalidCourseDifficulty, errConstant.ErrInvalidCourseEstimatedHours:
		return http.StatusUnprocessableEntity
	case errConstant.ErrCourseAlreadyPublished, errConstant.ErrCourseNotPublished:
		return http.StatusBadRequest
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary      Create Course
// @Description  Create a new course entry (admin only)
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateCourseRequest true "Course details"
// @Success      201 {object} dto.CourseSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response "Course with this title already exists for this JLPT level"
// @Failure      422 {object} response.Response "Invalid JLPT level ID or difficulty"
// @Failure      500 {object} response.Response
// @Router       /courses [post]
func (c *CourseController) Create(ctx *gin.Context) {
	request := &dto.CreateCourseRequest{}

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

	course, err := c.service.GetCourse().Create(ctx, request)
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
		Data: course,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all Courses
// @Description  Retrieve courses with advanced filtering, search, sorting, and pagination
// @Tags         Courses
// @Produce      json
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        jlptLevelId query int false "Filter by JLPT Level ID" example(5)
// @Param        difficulty query int false "Filter by difficulty (1-5)" minimum(1) maximum(5)
// @Param        isPublished query bool false "Filter by publication status" example(true)
// @Param        search query string false "Search in title or description" example("japanese")
// @Param        sortBy query string false "Sort by field (title, difficulty, created_at)" default(created_at) example("title")
// @Param        sortOrder query string false "Sort order (asc, desc)" default(desc) example("asc")
// @Success      200 {object} dto.CourseListSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /courses [get]
func (c *CourseController) GetAll(ctx *gin.Context) {
	filter := &dto.CourseFilterRequest{}

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

	courses, err := c.service.GetCourse().GetAll(ctx, filter)
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
		"pagination": courses.Pagination,
		"status":     "success",
		"data":       courses.Data,
	})
}

// GetByID godoc
// @Summary      Get Course by ID
// @Description  Retrieve a specific course entry by ID
// @Tags         Courses
// @Produce      json
// @Param        id path int true "Course ID"
// @Success      200 {object} dto.CourseSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response "Course not found"
// @Failure      500 {object} response.Response
// @Router       /courses/{id} [get]
func (c *CourseController) GetByID(ctx *gin.Context) {
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

	course, err := c.service.GetCourse().GetByID(ctx, uint(id))
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
		Data: course,
		Gin:  ctx,
	})
}

// Update godoc
// @Summary      Update Course
// @Description  Update an existing course entry by ID (admin only)
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Course ID"
// @Param        request body dto.UpdateCourseRequest true "Updated course details"
// @Success      200 {object} dto.CourseSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Course not found"
// @Failure      409 {object} response.Response "Course with this title already exists for this JLPT level"
// @Failure      422 {object} response.Response "Invalid JLPT level ID or difficulty"
// @Failure      500 {object} response.Response
// @Router       /courses/{id} [put]
func (c *CourseController) Update(ctx *gin.Context) {
	request := &dto.UpdateCourseRequest{}
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

	course, err := c.service.GetCourse().Update(ctx, request, uint(id))
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
		Data: course,
		Gin:  ctx,
	})
}

// Delete godoc
// @Summary      Delete Course
// @Description  Delete a course entry by ID (admin only)
// @Tags         Courses
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Course ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Course not found"
// @Failure      500 {object} response.Response
// @Router       /courses/{id} [delete]
func (c *CourseController) Delete(ctx *gin.Context) {
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

	err = c.service.GetCourse().Delete(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	successMessage := "Course deleted successfully"
	response.HttpResponse(response.ParamHTTPResp{
		Code:    http.StatusOK,
		Message: &successMessage,
		Gin:     ctx,
	})
}

// Publish godoc
// @Summary      Publish Course
// @Description  Publish a course by ID (admin only) - sets is_published=true and published_at=now
// @Tags         Courses
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Course ID"
// @Success      200 {object} dto.CourseSwaggerResponse
// @Failure      400 {object} response.Response "Course is already published"
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Course not found"
// @Failure      500 {object} response.Response
// @Router       /courses/{id}/publish [post]
func (c *CourseController) Publish(ctx *gin.Context) {
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

	course, err := c.service.GetCourse().Publish(ctx, uint(id))
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
		Data: course,
		Gin:  ctx,
	})
}

// Unpublish godoc
// @Summary      Unpublish Course
// @Description  Unpublish a course by ID (admin only) - sets is_published=false
// @Tags         Courses
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Course ID"
// @Success      200 {object} dto.CourseSwaggerResponse
// @Failure      400 {object} response.Response "Course is not published"
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Course not found"
// @Failure      500 {object} response.Response
// @Router       /courses/{id}/unpublish [post]
func (c *CourseController) Unpublish(ctx *gin.Context) {
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

	course, err := c.service.GetCourse().Unpublish(ctx, uint(id))
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
		Data: course,
		Gin:  ctx,
	})
}

// GetPublished godoc
// @Summary      Get Published Courses
// @Description  Retrieve only published courses with filtering, search, sorting, and pagination
// @Tags         Courses
// @Produce      json
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        jlptLevelId query int false "Filter by JLPT Level ID" example(5)
// @Param        difficulty query int false "Filter by difficulty (1-5)" minimum(1) maximum(5)
// @Param        search query string false "Search in title or description" example("japanese")
// @Param        sortBy query string false "Sort by field (title, difficulty, created_at)" default(created_at) example("title")
// @Param        sortOrder query string false "Sort order (asc, desc)" default(desc) example("asc")
// @Success      200 {object} dto.CourseListSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /courses/published [get]
func (c *CourseController) GetPublished(ctx *gin.Context) {
	filter := &dto.CourseFilterRequest{}

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

	courses, err := c.service.GetCourse().GetPublished(ctx, filter)
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
		"pagination": courses.Pagination,
		"status":     "success",
		"data":       courses.Data,
	})
}
