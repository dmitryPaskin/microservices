package modules

import (
	"microservices/geo/internal/infrastructure/component"
	geoservice "microservices/geo/internal/modules/geo/service"
	"microservices/geo/internal/storages"
)

type Services struct {
	Geo geoservice.Georer
}

func NewServices(storages *storages.Storages, components *component.Components) *Services {
	geoService := geoservice.NewGeo(storages.Geo, components.Logger, components.RateLimit, components.MQ)

	return &Services{
		Geo: geoService,
	}
}
