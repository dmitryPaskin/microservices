package service

import (
	"context"
	"proxy/internal/models"
)

type Userer interface {
	Profile(ctx context.Context, email string) (models.User, error)
	List(ctx context.Context) ([]models.User, error)
}

type ProfileIn struct {
	Email string
}

type ProfileOut struct {
	Name     string
	Email    string
	Password string
}

type ListIn struct{}

type ListOut struct {
	Users []models.User
}
