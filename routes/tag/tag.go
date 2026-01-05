package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type TagRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type ITagRoute interface {
	Run()
}

func NewTagRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) ITagRoute {
	return &TagRoute{controller: controller, group: group}
}

func (r *TagRoute) Run() {
	group := r.group.Group("/tags")
	group.GET("", r.controller.GetTagController().GetAll)
	group.GET("/:id", r.controller.GetTagController().GetByID)
	group.GET("/search", r.controller.GetTagController().GetByName)
	group.POST("", middlewares.Authenticate(), r.controller.GetTagController().Create)
	group.PUT("/:id", middlewares.Authenticate(), r.controller.GetTagController().Update)
	group.DELETE("/:id", middlewares.Authenticate(), r.controller.GetTagController().Delete)
}
