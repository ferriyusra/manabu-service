package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type ExerciseRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IExerciseRoute interface {
	Run()
}

func NewExerciseRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IExerciseRoute {
	return &ExerciseRoute{controller: controller, group: group}
}

func (r *ExerciseRoute) Run() {
	// Main exercises routes
	exerciseGroup := r.group.Group("/exercises")

	// Public endpoints
	exerciseGroup.GET("", r.controller.GetExerciseController().GetAll)
	exerciseGroup.GET("/:id", r.controller.GetExerciseController().GetByID)
	// Nested route: Get questions by exercise ID
	exerciseGroup.GET("/:id/questions", r.controller.GetExerciseQuestionController().GetByExerciseID)

	// Admin endpoints (require authentication)
	exerciseGroup.POST("", middlewares.Authenticate(), r.controller.GetExerciseController().Create)
	exerciseGroup.PUT("/:id", middlewares.Authenticate(), r.controller.GetExerciseController().Update)
	exerciseGroup.DELETE("/:id", middlewares.Authenticate(), r.controller.GetExerciseController().Delete)
	exerciseGroup.PATCH("/:id/publish", middlewares.Authenticate(), r.controller.GetExerciseController().UpdatePublishStatus)
}
