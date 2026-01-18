package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type UserCourseProgressRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IUserCourseProgressRoute interface {
	Run()
}

func NewUserCourseProgressRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IUserCourseProgressRoute {
	return &UserCourseProgressRoute{controller: controller, group: group}
}

func (r *UserCourseProgressRoute) Run() {
	// User course progress routes (all require authentication)
	progressGroup := r.group.Group("/user-course-progress")
	progressGroup.Use(middlewares.Authenticate())

	// All endpoints require authentication
	progressGroup.POST("", r.controller.GetUserCourseProgressController().Create)
	progressGroup.GET("", r.controller.GetUserCourseProgressController().GetAll)
	progressGroup.GET("/:id", r.controller.GetUserCourseProgressController().GetByID)
	progressGroup.PUT("/:id", r.controller.GetUserCourseProgressController().Update)
}
