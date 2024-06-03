package service

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
)

func NewToken() *jwtauth.JWTAuth {
	secret := os.Getenv("MY_SECRET")
	return jwtauth.New("HS256", []byte(secret), nil)
}
