package docs

import (
	"proxy/internal/models"
	"proxy/internal/modules/user/controller"
)

// swagger:route POST /api/user/profile user ProfileRequest
// Получить профиль пользователя по email.
// security:
//   - Bearer: []
// responses:
//  200: ProfileResponse

// swagger:parameters ProfileRequest
type ProfileRequest struct {
	// in:body
	// required: true
	Body controller.ProfileRequest
}

// swagger:response ProfileResponse
type ProfileResponse struct {
	// in:body
	Body models.User
}

// swagger:route GET /api/user/list user ListRequest
// Получить список пользователей.
// security:
//   - Bearer: []
// responses:
//  200: ListResponse

// swagger:response ListResponse
type ListResponse struct {
	// in:body
	List []models.User
}
