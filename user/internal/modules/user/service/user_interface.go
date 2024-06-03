package service

import (
	"context"
	"microservices/user/internal/models"
)

type Userer interface {
	Profile(email string) (models.UserDTO, error)
	GeyByID(ctx context.Context, id int) (models.User, error)
	Create(user models.UserDTO) error
	List() ([]models.User, error)
}
