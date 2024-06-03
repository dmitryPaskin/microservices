package service

import (
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=Auther
type Auther interface {
	Register(in RegisterIn) RegisterOut
	Login(in LoginIn) LoginOut
}

type LoginIn struct {
	Ctx      context.Context
	Email    string
	Password string
}

type LoginOut struct {
	Success bool
	Message string
}

type RegisterIn struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

type RegisterOut struct {
	Status  int
	Message string
}
