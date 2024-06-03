package controller

import (
	"encoding/json"
	"net/http"
	"proxy/internal/infrastructure/component"
	"proxy/internal/infrastructure/responder"
	"proxy/internal/modules/auth/service"

	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	authRegisterRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_register_requests_total",
		Help: "Total number of requests",
	})
	authLoginRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_login_request_total",
		Help: "Total number of requests",
	})
	authRegisterDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "auth_register_duration_seconds",
		Help: "Reqeust duration in seconds",
	})
	authLoginDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "auth_login_duration_seconds",
		Help: "Request duration in seconds",
	})
)

func init() {
	prometheus.MustRegister(authRegisterRequestsTotal)
	prometheus.MustRegister(authLoginRequestsTotal)
	prometheus.MustRegister(authRegisterDuration)
	prometheus.MustRegister(authLoginDuration)
}

type Auther interface {
	Register(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

type Auth struct {
	auth service.Auther
	responder.Responder
}

func NewAuth(service service.Auther, components *component.Components) Auther {
	return &Auth{
		auth:      service,
		Responder: components.Responder,
	}
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	authLoginRequestsTotal.Inc()

	var logReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
		a.Responder.ErrorBadRequest(w, err)
		return
	}

	out := a.auth.Login(service.LoginIn{Ctx: r.Context(), Email: logReq.Email, Password: logReq.Password})

	logResp := LoginResponse{
		Success: out.Success,
		Message: out.Message,
	}

	a.OutputJSON(w, logResp)

	duration := time.Since(startTime).Seconds()
	authLoginDuration.Observe(duration)
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	authRegisterRequestsTotal.Inc()

	var regReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&regReq); err != nil {
		a.ErrorBadRequest(w, err)
		return
	}

	out := a.auth.Register(service.RegisterIn{Name: regReq.Name, Email: regReq.Email, Password: regReq.Password, Phone: regReq.Phone})
	switch out.Status {
	case http.StatusConflict:
		http.Error(w, out.Message, http.StatusConflict)
		return
	case http.StatusInternalServerError:
		http.Error(w, out.Message, http.StatusInternalServerError)
		return
	default:
	}

	// log.Println(out.Status)
	// switch out.Status {
	// case int(codes.AlreadyExists):
	// 	http.Error(w, out.Message, http.StatusConflict)
	// 	return
	// case int(codes.Internal):
	// 	http.Error(w, out.Message, http.StatusInternalServerError)
	// 	return
	// }

	regResp := RegisterReponse{
		Success: true,
		Message: "Пользователь успешно зарегистрирован",
	}

	a.OutputJSON(w, regResp)

	duration := time.Since(startTime).Seconds()
	authRegisterDuration.Observe(duration)
}
