package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type LessonRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type ILessonRoute interface {
	Run()
}

func NewLessonRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) ILessonRoute {
	return &LessonRoute{controller: controller, group: group}
}

func (r *LessonRoute) Run() {
	// Main lessons routes
	lessonGroup := r.group.Group("/lessons")

	// Public endpoints
	lessonGroup.GET("", r.controller.GetLessonController().GetAll)
	lessonGroup.GET("/:id", r.controller.GetLessonController().GetByID)

	// Admin endpoints (require authentication)
	lessonGroup.POST("", middlewares.Authenticate(), r.controller.GetLessonController().Create)
	lessonGroup.PUT("/:id", middlewares.Authenticate(), r.controller.GetLessonController().Update)
	lessonGroup.DELETE("/:id", middlewares.Authenticate(), r.controller.GetLessonController().Delete)
	lessonGroup.POST("/:id/publish", middlewares.Authenticate(), r.controller.GetLessonController().Publish)
	lessonGroup.POST("/:id/unpublish", middlewares.Authenticate(), r.controller.GetLessonController().Unpublish)
}
