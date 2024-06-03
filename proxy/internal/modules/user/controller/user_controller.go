package controller

import (
	"encoding/json"
	"net/http"
	"proxy/internal/infrastructure/component"
	"proxy/internal/infrastructure/errors"
	"proxy/internal/infrastructure/responder"
	"proxy/internal/modules/user/service"
)

type Userer interface {
	Profile(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	user service.Userer
	responder.Responder
}

func NewUserController(service service.Userer, components *component.Components) Userer {
	return &UserController{
		user:      service,
		Responder: components.Responder,
	}
}

func (u *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	var req ProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.ErrorBadRequest(w, err)
		return
	}

	user, err := u.user.Profile(r.Context(), req.Email)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			u.ErrorNotFound(w, err)
			return
		default:
			u.ErrorInternal(w, err)
			return
		}
	}

	u.OutputJSON(w, user)
}

func (u *UserController) List(w http.ResponseWriter, r *http.Request) {
	users, err := u.user.List(r.Context())
	if err != nil {
		u.ErrorInternal(w, err)
	}

	u.OutputJSON(w, users)
}
