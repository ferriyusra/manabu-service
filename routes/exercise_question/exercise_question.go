package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type ExerciseQuestionRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IExerciseQuestionRoute interface {
	Run()
}

func NewExerciseQuestionRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IExerciseQuestionRoute {
	return &ExerciseQuestionRoute{controller: controller, group: group}
}

func (r *ExerciseQuestionRoute) Run() {
	// Main exercise questions routes
	questionGroup := r.group.Group("/exercise-questions")

	// Public endpoints
	questionGroup.GET("", r.controller.GetExerciseQuestionController().GetAll)
	questionGroup.GET("/:id", r.controller.GetExerciseQuestionController().GetByID)

	// Admin endpoints (require authentication)
	questionGroup.POST("", middlewares.Authenticate(), r.controller.GetExerciseQuestionController().Create)
	questionGroup.PUT("/:id", middlewares.Authenticate(), r.controller.GetExerciseQuestionController().Update)
	questionGroup.DELETE("/:id", middlewares.Authenticate(), r.controller.GetExerciseQuestionController().Delete)
	questionGroup.PATCH("/:id/publish", middlewares.Authenticate(), r.controller.GetExerciseQuestionController().UpdatePublishStatus)
}
