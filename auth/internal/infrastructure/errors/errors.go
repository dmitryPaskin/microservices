package errors

import "fmt"

var (
	ErrNotFound  = fmt.Errorf("not found")
	ErrWrongPass = fmt.Errorf("wrong password")
)
