package service

import (
	"context"
	"microservices/user/internal/models"
	"microservices/user/internal/modules/user/storage"

	"go.uber.org/zap"
)

type User struct {
	storage storage.UserStorager
	logger  *zap.Logger
}

func NewUser(storage storage.UserStorager, logger *zap.Logger) Userer {
	return &User{
		storage: storage,
		logger:  logger,
	}
}

func (u *User) Profile(email string) (models.UserDTO, error) {
	user, err := u.storage.GetByEmail(email)
	if err != nil {
		return models.UserDTO{}, err
	}

	return user, nil
}

func (u *User) Create(user models.UserDTO) error {
	err := u.storage.Create(user)
	if err != nil {
		u.logger.Error("user: error create user", zap.Error(err))
	}

	return err
}

func (u *User) List() ([]models.User, error) {
	return u.storage.List()
}

func (u *User) GeyByID(ctx context.Context, id int) (models.User, error) {
	return u.storage.GetByID(ctx, id)
}
