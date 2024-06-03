package service

import (
	"context"
	"microservices/auth/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=Userer
type Userer interface {
	Profile(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, user models.User) error
}

type ProfileIn struct {
	Email string
}
type PrfileOut struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

type CreateIn struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

type CreateOut struct {
	Success bool
}
