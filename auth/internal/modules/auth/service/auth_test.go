package service

import (
	"context"
	"fmt"
	"microservices/auth/config"
	"microservices/auth/internal/infrastructure/logs"
	"microservices/auth/internal/infrastructure/tools/cryptography"
	"microservices/auth/internal/models"
	"microservices/auth/internal/modules/user/service/mocks"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	testName     = "name"
	testEmail    = "email"
	testPassword = "password"
)

func TestAuter_Register_AlreadyExists(t *testing.T) {
	in := RegisterIn{
		Name:     testName,
		Email:    testEmail,
		Password: testPassword,
	}

	mockUserer := mocks.NewUserer(t)
	mockUserer.On("Profile", mock.Anything, testEmail).Return(models.User{}, nil)

	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
	auth := NewAuth(mockUserer, logger)

	got := auth.Register(in)

	assert.Equal(t, http.StatusConflict, got.Status)
}

func TestAuter_Register(t *testing.T) {
	in := RegisterIn{
		Name:     testName,
		Email:    testEmail,
		Password: testPassword,
	}

	mockUserer := mocks.NewUserer(t)
	mockUserer.On("Profile", mock.Anything, testEmail).Return(models.User{}, status.Error(codes.NotFound, "not found"))
	mockUserer.On("Create", mock.Anything, mock.AnythingOfType("models.User")).Return(nil)

	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
	auth := NewAuth(mockUserer, logger)

	got := auth.Register(in)

	assert.Equal(t, http.StatusOK, got.Status)
}

func TestAuter_Register_Internal(t *testing.T) {
	in := RegisterIn{
		Name:     testName,
		Email:    testEmail,
		Password: testPassword,
	}

	mockUserer := mocks.NewUserer(t)
	mockUserer.On("Profile", mock.Anything, testEmail).Return(models.User{}, status.Error(codes.Internal, "some error"))

	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
	auth := NewAuth(mockUserer, logger)

	got := auth.Register(in)

	assert.Equal(t, http.StatusInternalServerError, got.Status)
}

func TestAuter_Register_Create_Error(t *testing.T) {
	in := RegisterIn{
		Name:     testName,
		Email:    testEmail,
		Password: testPassword,
	}

	mockUserer := mocks.NewUserer(t)
	mockUserer.On("Profile", mock.Anything, testEmail).Return(models.User{}, status.Error(codes.NotFound, "not found"))
	mockUserer.On("Create", mock.Anything, mock.AnythingOfType("models.User")).Return(fmt.Errorf("some error"))

	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
	auth := NewAuth(mockUserer, logger)

	got := auth.Register(in)

	assert.Equal(t, http.StatusInternalServerError, got.Status)
}

func TestAuterLogin_Register_ProfileError(t *testing.T) {
	in := LoginIn{
		Ctx:      context.Background(),
		Email:    testEmail,
		Password: testPassword,
	}

	mockUserer := mocks.NewUserer(t)
	mockUserer.On("Profile", mock.Anything, testEmail).Return(models.User{}, fmt.Errorf("some error"))

	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
	auth := NewAuth(mockUserer, logger)

	got := auth.Login(in)

	assert.False(t, got.Success)
}

func TestAuterLogin_Register(t *testing.T) {
	in := LoginIn{
		Ctx:      context.Background(),
		Email:    testEmail,
		Password: testPassword,
	}

	hashPassword, _ := cryptography.HashPassword(testPassword)
	mockUserer := mocks.NewUserer(t)
	mockUserer.On("Profile", mock.Anything, testEmail).Return(models.User{Password: hashPassword}, nil)

	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
	auth := NewAuth(mockUserer, logger)

	got := auth.Login(in)

	assert.True(t, got.Success)
}

func TestAuterLogin_Register_WrongPass(t *testing.T) {
	in := LoginIn{
		Ctx:      context.Background(),
		Email:    testEmail,
		Password: testPassword,
	}

	hashPassword, _ := cryptography.HashPassword("122")
	mockUserer := mocks.NewUserer(t)
	mockUserer.On("Profile", mock.Anything, testEmail).Return(models.User{Password: hashPassword}, nil)

	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
	auth := NewAuth(mockUserer, logger)

	got := auth.Login(in)

	assert.False(t, got.Success)
}
