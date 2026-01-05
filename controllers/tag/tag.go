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

type TagController struct {
	service services.IServiceRegistry
}

type ITagController interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
	GetByID(*gin.Context)
	GetByName(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

func NewTagController(service services.IServiceRegistry) ITagController {
	return &TagController{service: service}
}

// getStatusCode maps errors to appropriate HTTP status codes
func (c *TagController) getStatusCode(err error) int {
	switch err {
	case errConstant.ErrTagNotFound:
		return http.StatusNotFound
	case errConstant.ErrTagDuplicate:
		return http.StatusConflict
	case errConstant.ErrInvalidColor:
		return http.StatusUnprocessableEntity
	case errConstant.ErrInvalidTagName:
		return http.StatusBadRequest
	case errConstant.ErrInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Create godoc
// @Summary      Create Tag
// @Description  Create a new tag (admin only)
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateTagRequest true "Tag details"
// @Success      201 {object} dto.TagSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      409 {object} response.Response "Tag with this name already exists"
// @Failure      422 {object} response.Response "Invalid color format"
// @Failure      500 {object} response.Response
// @Router       /tags [post]
func (c *TagController) Create(ctx *gin.Context) {
	request := &dto.CreateTagRequest{}

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

	tag, err := c.service.GetTag().Create(ctx, request)
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
		Data: tag,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all Tags
// @Description  Retrieve all tags with optional search and pagination
// @Tags         Tags
// @Produce      json
// @Param        page query int false "Page number" default(1) minimum(1)
// @Param        limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param        search query string false "Search in tag name" example("grammar")
// @Success      200 {object} dto.TagListSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      422 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /tags [get]
func (c *TagController) GetAll(ctx *gin.Context) {
	filter := &dto.TagFilterRequest{}

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

	tags, err := c.service.GetTag().GetAll(ctx, filter)
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
		"pagination": tags.Pagination,
		"status":     "success",
		"data":       tags.Data,
	})
}

// GetByID godoc
// @Summary      Get Tag by ID
// @Description  Retrieve a specific tag by ID
// @Tags         Tags
// @Produce      json
// @Param        id path int true "Tag ID"
// @Success      200 {object} dto.TagSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response "Tag not found"
// @Failure      500 {object} response.Response
// @Router       /tags/{id} [get]
func (c *TagController) GetByID(ctx *gin.Context) {
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

	tag, err := c.service.GetTag().GetByID(ctx, uint(id))
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
		Data: tag,
		Gin:  ctx,
	})
}

// GetByName godoc
// @Summary      Search Tag by name
// @Description  Retrieve a tag by exact name match (case-insensitive)
// @Tags         Tags
// @Produce      json
// @Param        name query string true "Tag name" example("Grammar")
// @Success      200 {object} dto.TagSwaggerResponse
// @Failure      400 {object} response.Response "Name parameter is required"
// @Failure      404 {object} response.Response "Tag not found"
// @Failure      500 {object} response.Response
// @Router       /tags/search [get]
func (c *TagController) GetByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  errConstant.ErrInvalidID, // Reuse for missing parameter
			Gin:  ctx,
		})
		return
	}

	tag, err := c.service.GetTag().GetByName(ctx, name)
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
		Data: tag,
		Gin:  ctx,
	})
}

// Update godoc
// @Summary      Update Tag
// @Description  Update an existing tag by ID (admin only)
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Tag ID"
// @Param        request body dto.UpdateTagRequest true "Updated tag details"
// @Success      200 {object} dto.TagSwaggerResponse
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Tag not found"
// @Failure      409 {object} response.Response "Tag with this name already exists"
// @Failure      422 {object} response.Response "Invalid color format"
// @Failure      500 {object} response.Response
// @Router       /tags/{id} [put]
func (c *TagController) Update(ctx *gin.Context) {
	request := &dto.UpdateTagRequest{}
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

	tag, err := c.service.GetTag().Update(ctx, request, uint(id))
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
		Data: tag,
		Gin:  ctx,
	})
}

// Delete godoc
// @Summary      Delete Tag
// @Description  Delete a tag by ID (admin only)
// @Tags         Tags
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Tag ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response "Tag not found"
// @Failure      500 {object} response.Response
// @Router       /tags/{id} [delete]
func (c *TagController) Delete(ctx *gin.Context) {
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

	err = c.service.GetTag().Delete(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: c.getStatusCode(err),
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	successMessage := "Tag deleted successfully"
	response.HttpResponse(response.ParamHTTPResp{
		Code:    http.StatusOK,
		Message: &successMessage,
		Gin:     ctx,
	})
}
