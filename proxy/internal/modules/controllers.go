package modules

import (
	"proxy/internal/infrastructure/component"
	acontroller "proxy/internal/modules/auth/controller"
	aservice "proxy/internal/modules/auth/service"
	geocontroller "proxy/internal/modules/geo/controller"
	geoservice "proxy/internal/modules/geo/service"
	usercontroller "proxy/internal/modules/user/controller"
	userservice "proxy/internal/modules/user/service"
)

type Controllers struct {
	Auth acontroller.Auther
	Geo  geocontroller.Georer
	User usercontroller.Userer
}

func NewControllers(authRPC aservice.Auther, geoRPC geoservice.Georer, userRPC userservice.Userer, components *component.Components) *Controllers {
	authcontroller := acontroller.NewAuth(authRPC, components)
	geoController := geocontroller.NewGeoController(geoRPC, components)
	userController := usercontroller.NewUserController(userRPC, components)

	return &Controllers{
		Auth: authcontroller,
		Geo:  geoController,
		User: userController,
	}
}
