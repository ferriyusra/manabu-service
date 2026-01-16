package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type CourseRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type ICourseRoute interface {
	Run()
}

func NewCourseRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) ICourseRoute {
	return &CourseRoute{controller: controller, group: group}
}

func (r *CourseRoute) Run() {
	group := r.group.Group("/courses")

	// Public endpoints
	group.GET("", r.controller.GetCourseController().GetAll)
	group.GET("/published", r.controller.GetCourseController().GetPublished)
	group.GET("/:id", r.controller.GetCourseController().GetByID)
	group.GET("/:id/lessons", r.controller.GetLessonController().GetByCourseID)

	// Admin endpoints (require authentication)
	group.POST("", middlewares.Authenticate(), r.controller.GetCourseController().Create)
	group.PUT("/:id", middlewares.Authenticate(), r.controller.GetCourseController().Update)
	group.DELETE("/:id", middlewares.Authenticate(), r.controller.GetCourseController().Delete)
	group.POST("/:id/publish", middlewares.Authenticate(), r.controller.GetCourseController().Publish)
	group.POST("/:id/unpublish", middlewares.Authenticate(), r.controller.GetCourseController().Unpublish)
}
