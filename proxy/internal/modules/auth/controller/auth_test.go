package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"proxy/config"
	"proxy/internal/infrastructure/component"
	"proxy/internal/infrastructure/logs"
	"proxy/internal/infrastructure/responder"
	"proxy/internal/modules/auth/service"
	"proxy/internal/modules/auth/service/mocks"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth_Login_BadRequest(t *testing.T) {
	conf := config.NewAppConf()
	logger := logs.NewLogger(conf, os.Stdout)
	respond := responder.NewResponder(logger)
	components := component.NewComponents(conf, respond, logger)
	auth := NewAuth(mocks.NewAuther(t), components)

	req := map[string]interface{}{"username": 123, "password": 123}
	reqJSON, _ := json.Marshal(req)

	s := httptest.NewServer(http.HandlerFunc(auth.Login))
	defer s.Close()

	resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		t.Fatal("ошибка при выполнении тестового запроса:", err.Error())
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuth_Login(t *testing.T) {
	conf := config.NewAppConf()
	logger := logs.NewLogger(conf, os.Stdout)
	respond := responder.NewResponder(logger)
	components := component.NewComponents(conf, respond, logger)
	authMock := mocks.NewAuther(t)

	authMock.On("Login", mock.MatchedBy(func(in service.LoginIn) bool {
		return in.Email == "test" && in.Password == "123"
	})).Return(service.LoginOut{Success: true, Message: "token"})

	auth := NewAuth(authMock, components)

	logindReq := LoginRequest{
		Email:    "test",
		Password: "123",
	}

	reqBody, _ := json.Marshal(logindReq)

	s := httptest.NewServer(http.HandlerFunc(auth.Login))
	defer s.Close()

	resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var logResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&logResp)
	if err != nil {
		t.Fatal("ошибка при декодировании тестового ответа:", err.Error())
	}

	assert.True(t, logResp.Success)
	assert.NotEmpty(t, logResp.Message)
}

func TestAuth_Register_BadRequest(t *testing.T) {
	conf := config.NewAppConf()
	logger := logs.NewLogger(conf, os.Stdout)
	respond := responder.NewResponder(logger)
	components := component.NewComponents(conf, respond, logger)
	mockAuth := mocks.NewAuther(t)
	auth := NewAuth(mockAuth, components)

	req := map[string]interface{}{"username": 123, "password": 123}
	reqJSON, _ := json.Marshal(req)

	s := httptest.NewServer(http.HandlerFunc(auth.Register))
	defer s.Close()

	resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		t.Fatal("ошибка при выполнении тестового запроса:", err.Error())
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuth_Register(t *testing.T) {
	conf := config.NewAppConf()
	logger := logs.NewLogger(conf, os.Stdout)
	respond := responder.NewResponder(logger)
	components := component.NewComponents(conf, respond, logger)

	authMock := mocks.NewAuther(t)
	authMock.On("Register", service.RegisterIn{Name: "test", Email: "email", Password: "123"}).Return(service.RegisterOut{Status: http.StatusOK, Message: ""})

	auth := NewAuth(authMock, components)

	req := RegisterRequest{
		Name:     "test",
		Email:    "email",
		Password: "123",
	}

	reqBody, _ := json.Marshal(req)

	s := httptest.NewServer(http.HandlerFunc(auth.Register))
	defer s.Close()

	resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var regResp RegisterReponse
	err = json.NewDecoder(resp.Body).Decode(&regResp)
	if err != nil {
		t.Fatal("ошибка при декодировании тестового ответа:", err.Error())
	}

	assert.True(t, regResp.Success)
	assert.Equal(t, "Пользователь успешно зарегистрирован", regResp.Message)
}

func TestAtuh_Register_Error(t *testing.T) {
	conf := config.NewAppConf()
	logger := logs.NewLogger(conf, os.Stdout)
	respond := responder.NewResponder(logger)
	components := component.NewComponents(conf, respond, logger)
	authMock := mocks.NewAuther(t)

	authMock.On("Register", service.RegisterIn{Name: "test", Email: "email", Password: "123"}).Return(service.RegisterOut{Status: http.StatusConflict, Message: "пользователь с таким именем уже зарегестрирован"})

	auth := NewAuth(authMock, components)

	req := RegisterRequest{
		Name:     "test",
		Email:    "email",
		Password: "123",
	}

	reqBody, _ := json.Marshal(req)

	s := httptest.NewServer(http.HandlerFunc(auth.Register))
	defer s.Close()

	resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}
