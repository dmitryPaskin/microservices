package service

import (
	"context"
	"fmt"
	"os"

	"microservices/auth/internal/infrastructure/tools/cryptography"
	"microservices/auth/internal/models"
	"microservices/auth/internal/modules/user/service"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte(os.Getenv("MY_SECRET")), nil)
}

type Auth struct {
	user   service.Userer
	logger *zap.Logger
}

func NewAuth(userService service.Userer, logger *zap.Logger) Auther {
	return &Auth{
		user:   userService,
		logger: logger,
	}
}

func (a *Auth) Register(in RegisterIn) RegisterOut {
	hashPassword, err := cryptography.HashPassword(in.Password)
	if err != nil {
		a.logger.Error("error hashing password", zap.Error(err))
		return RegisterOut{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	// _, err = a.user.Profile(context.Background(), in.Email)
	// if err == nil {
	// 	a.logger.Info("user already exists")
	// 	return RegisterOut{
	// 		Status: http.StatusConflict,
	// 		Error:  fmt.Errorf("user already exists"),
	// 	}
	// } else if err != errors.ErrNotFound {
	// 	a.logger.Error("register", zap.Error(err))
	// 	return RegisterOut{
	// 		Status: http.StatusInternalServerError,
	// 		Error:  err,
	// 	}
	// }

	_, err = a.user.Profile(context.Background(), in.Email)
	if err == nil {
		a.logger.Info("user already exists")
		return RegisterOut{
			Status: http.StatusConflict,
			Error:  status.Error(codes.AlreadyExists, "user already exists"),
		}
	} else {
		st, ok := status.FromError(err)
		if ok && st.Code() != codes.NotFound {
			a.logger.Error("register", zap.Error(err))
			return RegisterOut{
				Status: http.StatusInternalServerError,
				Error:  status.Error(codes.Internal, err.Error()),
			}
		}
	}

	user := models.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: hashPassword,
		Phone:    in.Phone,
	}

	if err := a.user.Create(context.Background(), user); err != nil {
		a.logger.Error("create user error", zap.Error(err))
		return RegisterOut{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	a.logger.Info("user crerated")
	return RegisterOut{
		Status: http.StatusOK,
		Error:  nil,
	}
}

func (a *Auth) Login(in LoginIn) LoginOut {
	user, err := a.user.Profile(in.Ctx, in.Email)
	if err != nil {
		return LoginOut{
			Success: false,
			Message: err.Error(),
		}
	}

	if !cryptography.CheckPassword(user.Password, in.Password) {
		return LoginOut{
			Success: false,
			Message: "Неверный пароль",
		}
	}

	// _, claims, _ := jwtauth.FromContext(in.Ctx)

	claims := jwt.MapClaims{
		"id": fmt.Sprintf("%d", user.ID),
	}
	_, tokenString, _ := TokenAuth.Encode(claims)

	return LoginOut{
		Success: true,
		Message: tokenString,
	}
}
