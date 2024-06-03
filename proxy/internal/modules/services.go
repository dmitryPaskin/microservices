package modules

import (
	aservice "proxy/internal/modules/auth/service"
	geoservice "proxy/internal/modules/geo/service"
	userservice "proxy/internal/modules/user/service"
)

type Services struct {
	Auth aservice.Auther
	Geo  geoservice.Georer
	User userservice.Userer
}
