package routes

import (
	"manabu-service/controllers"
	jlptLevelRoute "manabu-service/routes/jlpt_level"
	routes "manabu-service/routes/user"

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
}

func (r *Registry) userRoute() routes.IUserRoute {
	return routes.NewUserRoute(r.controller, r.group)
}

func (r *Registry) jlptLevelRoute() jlptLevelRoute.IJlptLevelRoute {
	return jlptLevelRoute.NewJlptLevelRoute(r.controller, r.group)
}
