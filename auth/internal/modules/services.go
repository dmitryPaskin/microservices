package modules

import (
	"microservices/auth/internal/infrastructure/component"
	aservice "microservices/auth/internal/modules/auth/service"
	userservice "microservices/auth/internal/modules/user/service"
)

type Services struct {
	Auth aservice.Auther
}

func NewServices(userService userservice.Userer, components *component.Components) *Services {
	authService := aservice.NewAuth(userService, components.Logger)

	return &Services{
		Auth: authService,
	}
}
