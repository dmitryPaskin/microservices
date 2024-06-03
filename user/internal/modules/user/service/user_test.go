package service

import (
	"microservices/user/internal/models"
)

var (
	testName  = "name"
	testEmail = "email"
	testUser  = models.User{
		Name:  testName,
		Email: testEmail,
	}
)

// func TestUser_Profile(t *testing.T) {
// 	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
// 	mockStorage := mocks.NewUserStorager(t)
// 	mockStorage.On("GetByEmail", testEmail).Return(testUser)
// 	user := NewUser(mockStorage, logger)

// 	got, err := user.Profile(testEmail)

// 	assert.Nil(t, err)
// 	assert.Equal(t, testUser, got)
// }
