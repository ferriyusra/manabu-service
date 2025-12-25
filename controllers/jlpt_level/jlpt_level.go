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

type JlptLevelController struct {
	service services.IServiceRegistry
}

type IJlptLevelController interface {
	Create(*gin.Context)
	GetAll(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

func NewJlptLevelController(service services.IServiceRegistry) IJlptLevelController {
	return &JlptLevelController{service: service}
}

// Create godoc
// @Summary      Create JLPT Level
// @Description  Create a new JLPT level (admin only)
// @Tags         JLPT Levels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateJlptLevelRequest true "JLPT Level details"
// @Success      201 {object} response.JlptLevelResponse{data=dto.JlptLevelResponse}
// @Failure      400 {object} response.JlptLevelResponse
// @Failure      401 {object} response.JlptLevelResponse
// @Failure      422 {object} response.JlptLevelResponse
// @Router       /jlpt-levels [post]
func (c *JlptLevelController) Create(ctx *gin.Context) {
	request := &dto.CreateJlptLevelRequest{}

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

	jlptLevel, err := c.service.GetJlptLevel().Create(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: jlptLevel,
		Gin:  ctx,
	})
}

// GetAll godoc
// @Summary      Get all JLPT Levels
// @Description  Retrieve all JLPT levels ordered by level
// @Tags         JLPT Levels
// @Produce      json
// @Success      200 {object} response.JlptLevelResponse{data=[]dto.JlptLevelResponse}
// @Failure      400 {object} response.JlptLevelResponse
// @Router       /jlpt-levels [get]
func (c *JlptLevelController) GetAll(ctx *gin.Context) {
	jlptLevels, err := c.service.GetJlptLevel().GetAll(ctx)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: jlptLevels,
		Gin:  ctx,
	})
}

// GetByID godoc
// @Summary      Get JLPT Level by ID
// @Description  Retrieve a specific JLPT level by ID
// @Tags         JLPT Levels
// @Produce      json
// @Param        id path int true "JLPT Level ID"
// @Success      200 {object} response.JlptLevelResponse{data=dto.JlptLevelResponse}
// @Failure      400 {object} response.JlptLevelResponse
// @Router       /jlpt-levels/{id} [get]
func (c *JlptLevelController) GetByID(ctx *gin.Context) {
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

	jlptLevel, err := c.service.GetJlptLevel().GetByID(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: jlptLevel,
		Gin:  ctx,
	})
}

// Update godoc
// @Summary      Update JLPT Level
// @Description  Update an existing JLPT level by ID (admin only)
// @Tags         JLPT Levels
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "JLPT Level ID"
// @Param        request body dto.UpdateJlptLevelRequest true "Updated JLPT Level details"
// @Success      200 {object} response.JlptLevelResponse{data=dto.JlptLevelResponse}
// @Failure      400 {object} response.JlptLevelResponse
// @Failure      401 {object} response.JlptLevelResponse
// @Failure      422 {object} response.JlptLevelResponse
// @Router       /jlpt-levels/{id} [put]
func (c *JlptLevelController) Update(ctx *gin.Context) {
	request := &dto.UpdateJlptLevelRequest{}
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

	jlptLevel, err := c.service.GetJlptLevel().Update(ctx, request, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: jlptLevel,
		Gin:  ctx,
	})
}

// Delete godoc
// @Summary      Delete JLPT Level
// @Description  Delete a JLPT level by ID (admin only)
// @Tags         JLPT Levels
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "JLPT Level ID"
// @Success      200 {object} response.JlptLevelResponse
// @Failure      400 {object} response.JlptLevelResponse
// @Failure      401 {object} response.JlptLevelResponse
// @Router       /jlpt-levels/{id} [delete]
func (c *JlptLevelController) Delete(ctx *gin.Context) {
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

	err = c.service.GetJlptLevel().Delete(ctx, uint(id))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	successMessage := "JLPT level deleted successfully"
	response.HttpResponse(response.ParamHTTPResp{
		Code:    http.StatusOK,
		Message: &successMessage,
		Gin:     ctx,
	})
}
