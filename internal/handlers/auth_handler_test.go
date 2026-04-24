package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coolmate/ecommerce-backend/internal/services"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *services.RegisterRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(req *services.LoginRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Refresh(refreshToken string) (*services.AuthResponse, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Logout(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	authResp := &services.AuthResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		User: &services.UserDTO{
			ID:    1,
			Email: "test@example.com",
			Role:  "customer",
		},
	}

	registerReq := &services.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	mockSvc.On("Register", mock.MatchedBy(func(req *services.RegisterRequest) bool {
		return req.Email == "test@example.com"
	})).Return(authResp, nil)

	handler := NewAuthHandler(mockSvc)

	body, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	handler := NewAuthHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	mockSvc.On("Register", mock.Anything).Return(nil, assert.AnError)

	handler := NewAuthHandler(mockSvc)

	registerReq := &services.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	body, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	authResp := &services.AuthResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		User: &services.UserDTO{
			ID:    1,
			Email: "test@example.com",
			Role:  "customer",
		},
	}

	loginReq := &services.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockSvc.On("Login", mock.MatchedBy(func(req *services.LoginRequest) bool {
		return req.Email == "test@example.com"
	})).Return(authResp, nil)

	handler := NewAuthHandler(mockSvc)

	body, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	handler := NewAuthHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	mockSvc.On("Login", mock.Anything).Return(nil, assert.AnError)

	handler := NewAuthHandler(mockSvc)

	loginReq := &services.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Refresh_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	authResp := &services.AuthResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}

	mockSvc.On("Refresh", "refresh_token").Return(authResp, nil)

	handler := NewAuthHandler(mockSvc)

	refreshReq := map[string]string{"refreshToken": "refresh_token"}
	body, _ := json.Marshal(refreshReq)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Refresh(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Refresh_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	handler := NewAuthHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Refresh(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Refresh_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	mockSvc.On("Refresh", "invalid_token").Return(nil, assert.AnError)

	handler := NewAuthHandler(mockSvc)

	refreshReq := map[string]string{"refreshToken": "invalid_token"}
	body, _ := json.Marshal(refreshReq)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Refresh(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	mockSvc.On("Logout", "refresh_token").Return(nil)

	handler := NewAuthHandler(mockSvc)

	logoutReq := map[string]string{"refreshToken": "refresh_token"}
	body, _ := json.Marshal(logoutReq)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/logout", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Logout_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	mockSvc.On("Logout", "refresh_token").Return(assert.AnError)

	handler := NewAuthHandler(mockSvc)

	logoutReq := map[string]string{"refreshToken": "refresh_token"}
	body, _ := json.Marshal(logoutReq)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/logout", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Logout(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_RegisterVendor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	handler := NewAuthHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/vendor/register", nil)

	handler.RegisterVendor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
