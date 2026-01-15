package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type CategoryRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type ICategoryRoute interface {
	Run()
}

func NewCategoryRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) ICategoryRoute {
	return &CategoryRoute{controller: controller, group: group}
}

func (r *CategoryRoute) Run() {
	group := r.group.Group("/categories")
	group.GET("", r.controller.GetCategoryController().GetAll)
	group.GET("/:id", r.controller.GetCategoryController().GetByID)
	group.GET("/jlpt/:jlptLevelId", r.controller.GetCategoryController().GetByJlptLevelID)
	group.POST("", middlewares.Authenticate(), r.controller.GetCategoryController().Create)
	group.PUT("/:id", middlewares.Authenticate(), r.controller.GetCategoryController().Update)
	group.DELETE("/:id", middlewares.Authenticate(), r.controller.GetCategoryController().Delete)
}
