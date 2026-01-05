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

type CategoryController struct {
	service services.IServiceRegistry
}

type ICategoryController interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
	GetByID(*gin.Context)
	GetByJlptLevelID(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

func NewCategoryController(service services.IServiceRegistry) ICategoryController {
	return &CategoryController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *CategoryController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrCategoryNotFound:
		return http.StatusNotFound
	case errConstant.ErrCategoryNameExist:
		return http.StatusConflict
	case errConstant.ErrInvalidJlptLevelID:
		return http.StatusUnprocessableEntity
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary      Create Category
// @Description  Create a new category (admin only)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateCategoryRequest true "Category details"
// @Success      201 {object} response.Response{data=dto.CategoryResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      422 {object} response.Response
// @Router       /categories [post]
func (c *CategoryController) Create(ctx *gin.Context) {
	request := &dto.CreateCategoryRequest{}

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

	category, err := c.service.GetCategory().Create(ctx, request)
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
		Data: category,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all Categories
// @Description  Retrieve all categories with pagination
// @Tags         Categories
// @Produce      json
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Success      200 {object} response.Response{data=[]dto.CategoryResponse,pagination=dto.PaginationResponse}
// @Failure      400 {object} response.Response
// @Router       /categories [get]
func (c *CategoryController) GetAll(ctx *gin.Context) {
	pagination := &dto.PaginationRequest{}

	if err := ctx.ShouldBindQuery(pagination); err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(pagination)
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

	categories, err := c.service.GetCategory().GetAll(ctx, pagination)
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
		"pagination": categories.Pagination,
		"status":     "success",
		"data":       categories.Data,
	})
}

// GetByID godoc
// @Summary      Get Category by ID
// @Description  Retrieve a specific category by ID
// @Tags         Categories
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} response.Response{data=dto.CategoryResponse}
// @Failure      400 {object} response.Response
// @Router       /categories/{id} [get]
func (c *CategoryController) GetByID(ctx *gin.Context) {
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

	category, err := c.service.GetCategory().GetByID(ctx, uint(id))
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
		Data: category,
		Gin:  ctx,
	})
}

// GetByJlptLevelID godoc
// @Summary      Get Categories by JLPT Level
// @Description  Retrieve all categories for a specific JLPT level with pagination
// @Tags         Categories
// @Produce      json
// @Param        jlpt_level_id path int true "JLPT Level ID"
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Success      200 {object} response.Response{data=[]dto.CategoryResponse,pagination=dto.PaginationResponse}
// @Failure      400 {object} response.Response
// @Router       /categories/jlpt/{jlpt_level_id} [get]
func (c *CategoryController) GetByJlptLevelID(ctx *gin.Context) {
	idParam := ctx.Param("jlpt_level_id")
	if idParam == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	jlptLevelID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID,
			Gin:  ctx,
		})
		return
	}

	pagination := &dto.PaginationRequest{}

	if err := ctx.ShouldBindQuery(pagination); err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(pagination)
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

	categories, err := c.service.GetCategory().GetByJlptLevelID(ctx, uint(jlptLevelID), pagination)
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
		"pagination": categories.Pagination,
		"status":     "success",
		"data":       categories.Data,
	})
}

// Update godoc
// @Summary      Update Category
// @Description  Update an existing category by ID (admin only)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Param        request body dto.UpdateCategoryRequest true "Updated category details"
// @Success      200 {object} response.Response{data=dto.CategoryResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      422 {object} response.Response
// @Router       /categories/{id} [put]
func (c *CategoryController) Update(ctx *gin.Context) {
	request := &dto.UpdateCategoryRequest{}
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

	category, err := c.service.GetCategory().Update(ctx, request, uint(id))
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
		Data: category,
		Gin:  ctx,
	})
}

// Delete godoc
// @Summary      Delete Category
// @Description  Delete a category by ID (admin only)
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /categories/{id} [delete]
func (c *CategoryController) Delete(ctx *gin.Context) {
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

	err = c.service.GetCategory().Delete(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	successMessage := "Category deleted successfully"
	response.HttpResponse(response.ParamHTTPResp{
		Code:    http.StatusOK,
		Message: &successMessage,
		Gin:     ctx,
	})
}
