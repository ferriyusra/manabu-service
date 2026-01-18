package controllers

import (
	errWrap "manabu-service/common/error"
	"manabu-service/common/response"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserCourseProgressController struct {
	service services.IServiceRegistry
}

// IUserCourseProgressController defines the contract for user course progress HTTP handlers.
type IUserCourseProgressController interface {
	// Create handles POST requests to enroll a user in a course.
	Create(*gin.Context)
	// GetAll handles GET requests to retrieve all progress entries for the authenticated user.
	GetAll(*gin.Context)
	// GetByID handles GET requests to retrieve a specific progress entry by UUID.
	GetByID(*gin.Context)
	// Update handles PUT requests to update completed lessons count.
	Update(*gin.Context)
}

func NewUserCourseProgressController(service services.IServiceRegistry) IUserCourseProgressController {
	return &UserCourseProgressController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *UserCourseProgressController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrUserCourseProgressNotFound:
		return http.StatusNotFound
	case errConstant.ErrUserCourseProgressAlreadyExists:
		return http.StatusConflict
	case errConstant.ErrInvalidCourseIDProgress, errConstant.ErrInvalidUserIDProgress,
		errConstant.ErrInvalidProgressStatus, errConstant.ErrInvalidCompletedLessons,
		errConstant.ErrInvalidProgressPercentage, errConstant.ErrCompletedLessonsExceedTotal:
		return http.StatusUnprocessableEntity
	case errConstant.ErrCannotUpdateCompletedCourse:
		return http.StatusBadRequest
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// getUserIDFromContext extracts the user ID from the authenticated context
func (c *UserCourseProgressController) getUserIDFromContext(ctx *gin.Context) (uint, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, errConstant.ErrInvalidUserIDProgress
	}

	// Handle different possible types for user_id
	switch v := userID.(type) {
	case uint:
		return v, nil
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	default:
		return 0, errConstant.ErrInvalidUserIDProgress
	}
}

// Create godoc
// @Summary      Enroll in Course (Start Progress)
// @Description  Create a new progress record by enrolling the authenticated user in a course
// @Tags         User Course Progress
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateUserCourseProgressRequest true "Course enrollment details"
// @Success      201 {object} dto.UserCourseProgressSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response "User already enrolled in this course"
// @Failure      422 {object} response.Response "Invalid course ID"
// @Failure      500 {object} response.Response
// @Router       /user-course-progress [post]
func (c *UserCourseProgressController) Create(ctx *gin.Context) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	request := &dto.CreateUserCourseProgressRequest{}

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

	progress, err := c.service.GetUserCourseProgress().Create(ctx, userID, request)
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
		Data: progress,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get All User Course Progress
// @Description  Retrieve all course progress for the authenticated user with filtering, sorting, and pagination
// @Tags         User Course Progress
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        status query string false "Filter by status (not_started, in_progress, completed)" example("in_progress")
// @Param        courseId query int false "Filter by Course ID" example(1)
// @Param        sortBy query string false "Sort by field (last_accessed_at, progress_percentage, started_at)" default(last_accessed_at) example("last_accessed_at")
// @Param        sortOrder query string false "Sort order (asc, desc)" default(desc) example("desc")
// @Success      200 {object} dto.UserCourseProgressListSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /user-course-progress [get]
func (c *UserCourseProgressController) GetAll(ctx *gin.Context) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	filter := &dto.UserCourseProgressFilterRequest{}

	if err := ctx.ShouldBindQuery(filter); err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(filter)
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

	progressList, err := c.service.GetUserCourseProgress().GetAll(ctx, userID, filter)
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
		"pagination": progressList.Pagination,
		"status":     "success",
		"data":       progressList.Data,
	})
}

// GetByID godoc
// @Summary      Get User Course Progress by ID
// @Description  Retrieve a specific course progress record by UUID
// @Tags         User Course Progress
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User Course Progress UUID" format(uuid)
// @Success      200 {object} dto.UserCourseProgressSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "User course progress not found"
// @Failure      500 {object} response.Response
// @Router       /user-course-progress/{id} [get]
func (c *UserCourseProgressController) GetByID(ctx *gin.Context) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	idParam := ctx.Param("id")
	if idParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	progress, err := c.service.GetUserCourseProgress().GetByID(ctx, id, userID)
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
		Data: progress,
		Gin:  ctx,
	})
}

// Update godoc
// @Summary      Update User Course Progress
// @Description  Update progress by recording completed lessons count. Status and percentage are auto-calculated.
// @Tags         User Course Progress
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User Course Progress UUID" format(uuid)
// @Param        request body dto.UpdateUserCourseProgressRequest true "Updated progress details"
// @Success      200 {object} dto.UserCourseProgressSwaggerResponse
// @Failure      400 {object} response.Response "Cannot update completed course"
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "User course progress not found"
// @Failure      422 {object} response.Response "Invalid completed lessons count"
// @Failure      500 {object} response.Response
// @Router       /user-course-progress/{id} [put]
func (c *UserCourseProgressController) Update(ctx *gin.Context) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	request := &dto.UpdateUserCourseProgressRequest{}
	idParam := ctx.Param("id")
	if idParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	id, err := uuid.Parse(idParam)
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

	progress, err := c.service.GetUserCourseProgress().Update(ctx, request, id, userID)
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
		Data: progress,
		Gin:  ctx,
	})
}
