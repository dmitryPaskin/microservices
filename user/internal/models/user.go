package models

type UserDTO struct {
	ID       int
	Name     string
	Email    string
	Phone    string
	Password string
}

type User struct {
	Name  string
	Email string
	Phone string
}
