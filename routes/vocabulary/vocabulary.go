package routes

import (
	"manabu-service/controllers"
	"manabu-service/middlewares"

	"github.com/gin-gonic/gin"
)

type VocabularyRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IVocabularyRoute interface {
	Run()
}

func NewVocabularyRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IVocabularyRoute {
	return &VocabularyRoute{controller: controller, group: group}
}

func (r *VocabularyRoute) Run() {
	group := r.group.Group("/vocabularies")
	group.GET("", r.controller.GetVocabularyController().GetAll)
	group.GET("/:id", r.controller.GetVocabularyController().GetByID)
	group.POST("", middlewares.Authenticate(), r.controller.GetVocabularyController().Create)
	group.PUT("/:id", middlewares.Authenticate(), r.controller.GetVocabularyController().Update)
	group.DELETE("/:id", middlewares.Authenticate(), r.controller.GetVocabularyController().Delete)
}
