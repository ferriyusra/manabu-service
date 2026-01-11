package routes

import (
	"manabu-service/controllers"
	categoryRoute "manabu-service/routes/category"
	jlptLevelRoute "manabu-service/routes/jlpt_level"
	tagRoute "manabu-service/routes/tag"
	routes "manabu-service/routes/user"
	userVocabStatusRoute "manabu-service/routes/user_vocabulary_status"
	vocabularyRoute "manabu-service/routes/vocabulary"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IRouteRegister interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRouteRegister {
	return &Registry{controller: controller, group: group}
}

func (r *Registry) Serve() {
	r.userRoute().Run()
	r.jlptLevelRoute().Run()
	r.categoryRoute().Run()
	r.vocabularyRoute().Run()
	r.tagRoute().Run()
	r.userVocabularyStatusRoute().Run()
}

func (r *Registry) userRoute() routes.IUserRoute {
	return routes.NewUserRoute(r.controller, r.group)
}

func (r *Registry) jlptLevelRoute() jlptLevelRoute.IJlptLevelRoute {
	return jlptLevelRoute.NewJlptLevelRoute(r.controller, r.group)
}

func (r *Registry) categoryRoute() categoryRoute.ICategoryRoute {
	return categoryRoute.NewCategoryRoute(r.controller, r.group)
}

func (r *Registry) vocabularyRoute() vocabularyRoute.IVocabularyRoute {
	return vocabularyRoute.NewVocabularyRoute(r.controller, r.group)
}

func (r *Registry) tagRoute() tagRoute.ITagRoute {
	return tagRoute.NewTagRoute(r.controller, r.group)
}

func (r *Registry) userVocabularyStatusRoute() userVocabStatusRoute.IUserVocabularyStatusRoute {
	return userVocabStatusRoute.NewUserVocabularyStatusRoute(r.controller, r.group)
}
