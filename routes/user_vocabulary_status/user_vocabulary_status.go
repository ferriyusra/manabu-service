package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type UserVocabularyStatusRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IUserVocabularyStatusRoute interface {
	Run()
}

func NewUserVocabularyStatusRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IUserVocabularyStatusRoute {
	return &UserVocabularyStatusRoute{controller: controller, group: group}
}

func (r *UserVocabularyStatusRoute) Run() {
	group := r.group.Group("/user-vocabulary-status")
	group.POST("", middlewares.Authenticate(), r.controller.GetUserVocabularyStatusController().Create)
	group.GET("", middlewares.Authenticate(), r.controller.GetUserVocabularyStatusController().GetAll)
	group.GET("/due", middlewares.Authenticate(), r.controller.GetUserVocabularyStatusController().GetDueForReview)
	group.GET("/:id", middlewares.Authenticate(), r.controller.GetUserVocabularyStatusController().GetByID)
	group.POST("/:vocabulary_id/review", middlewares.Authenticate(), r.controller.GetUserVocabularyStatusController().Review)
}
