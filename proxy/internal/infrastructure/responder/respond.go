package responder

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
	ErrorNotFound(w http.ResponseWriter, err error)
	ErrorToManyRequests(w http.ResponseWriter, err error)
}

type Respond struct {
	log *zap.Logger
}

func NewResponder(logger *zap.Logger) Responder {
	return &Respond{log: logger}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		r.log.Error("responder json encode error", zap.Error(err))
		r.ErrorInternal(w, err)
	}
}

func (r *Respond) ErrorNotFound(w http.ResponseWriter, err error) {
	r.log.Info("http response Not Found")
	message := extractGrpcErrorMessage(err)
	http.Error(w, message, http.StatusNotFound)
}

func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.log.Warn("http response Unauthorized", zap.Error(err))
	message := extractGrpcErrorMessage(err)
	http.Error(w, message, http.StatusUnauthorized)
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	message := extractGrpcErrorMessage(err)
	http.Error(w, message, http.StatusBadRequest)
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	r.log.Error("http response internal error", zap.Error(err))
	message := extractGrpcErrorMessage(err)
	http.Error(w, message, http.StatusInternalServerError)
}

func (r *Respond) ErrorToManyRequests(w http.ResponseWriter, err error) {
	r.log.Error("http response to many requests status cide", zap.Error(err))
	message := extractGrpcErrorMessage(err)
	http.Error(w, message, http.StatusTooManyRequests)
}
