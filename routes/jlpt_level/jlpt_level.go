package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type JlptLevelRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IJlptLevelRoute interface {
	Run()
}

func NewJlptLevelRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IJlptLevelRoute {
	return &JlptLevelRoute{controller: controller, group: group}
}

func (r *JlptLevelRoute) Run() {
	group := r.group.Group("/jlpt-levels")
	group.GET("", r.controller.GetJlptLevelController().GetAll)
	group.GET("/:id", r.controller.GetJlptLevelController().GetByID)
	group.POST("", middlewares.Authenticate(), r.controller.GetJlptLevelController().Create)
	group.PUT("/:id", middlewares.Authenticate(), r.controller.GetJlptLevelController().Update)
	group.DELETE("/:id", middlewares.Authenticate(), r.controller.GetJlptLevelController().Delete)
}
