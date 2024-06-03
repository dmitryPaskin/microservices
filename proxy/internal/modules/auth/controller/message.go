package controller

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type RegisterReponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
