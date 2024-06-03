// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	models "microservices/user/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// UserStorager is an autogenerated mock type for the UserStorager type
type UserStorager struct {
	mock.Mock
}

// Create provides a mock function with given fields: user
func (_m *UserStorager) Create(user models.UserDTO) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.UserDTO) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByEmail provides a mock function with given fields: email
func (_m *UserStorager) GetByEmail(email string) (models.UserDTO, error) {
	ret := _m.Called(email)

	var r0 models.UserDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (models.UserDTO, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) models.UserDTO); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(models.UserDTO)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields:
func (_m *UserStorager) List() ([]models.User, error) {
	ret := _m.Called()

	var r0 []models.User
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.User, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserStorager creates a new instance of UserStorager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserStorager(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserStorager {
	mock := &UserStorager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
