package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	mockService "manabu-service/mocks/services"
	mockUserService "manabu-service/mocks/services/user"
)

type UserControllerTestSuite struct {
	suite.Suite
	controller        IUserController
	mockServiceReg    *mockService.MockIServiceRegistry
	mockUserService   *mockUserService.MockIUserService
	router            *gin.Engine
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (s *UserControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockServiceReg = mockService.NewMockIServiceRegistry(s.T())
	s.mockUserService = mockUserService.NewMockIUserService(s.T())
	s.controller = NewUserController(s.mockServiceReg)
	s.router = gin.New()
}

// ==================== Login Tests ====================

// Test Login - Success
func (s *UserControllerTestSuite) TestLogin_Success() {
	// Arrange
	loginReq := dto.LoginRequest{
		Username: "johndoe",
		Password: "password123",
	}

	userUUID := uuid.New()
	loginResp := &dto.LoginResponse{
		User: dto.UserResponse{
			UUID:     userUUID,
			Name:     "John Doe",
			Username: "johndoe",
			Email:    "john@example.com",
			Role:     "user",
		},
		Token: "mock-jwt-token",
	}

	reqBody, _ := json.Marshal(loginReq)

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().Login(mock.Anything, &loginReq).Return(loginResp, nil)

	// Setup route
	s.router.POST("/auth/login", s.controller.Login)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["status"])
	assert.NotNil(s.T(), response["data"])
	assert.NotNil(s.T(), response["token"])
}

// Test Login - Invalid JSON
func (s *UserControllerTestSuite) TestLogin_InvalidJSON() {
	// Arrange
	invalidJSON := []byte(`{"username": "johndoe", "password":}`)

	s.router.POST("/auth/login", s.controller.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// Test Login - Validation Error (Missing Username)
func (s *UserControllerTestSuite) TestLogin_ValidationError_MissingUsername() {
	// Arrange
	loginReq := dto.LoginRequest{
		Password: "password123",
	}

	reqBody, _ := json.Marshal(loginReq)

	s.router.POST("/auth/login", s.controller.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "error", response["status"])
}

// Test Login - Service Error (User Not Found)
func (s *UserControllerTestSuite) TestLogin_ServiceError_UserNotFound() {
	// Arrange
	loginReq := dto.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	reqBody, _ := json.Marshal(loginReq)

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().Login(mock.Anything, &loginReq).Return(nil, errConstant.ErrUserNotFound)

	s.router.POST("/auth/login", s.controller.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// ==================== Register Tests ====================

// Test Register - Success
func (s *UserControllerTestSuite) TestRegister_Success() {
	// Arrange
	registerReq := dto.RegisterRequest{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "password123",
		ConfirmPassword: "password123",
		Email:           "john@example.com",
	}

	userUUID := uuid.New()
	registerResp := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:     userUUID,
			Name:     "John Doe",
			Username: "johndoe",
			Email:    "john@example.com",
		},
	}

	reqBody, _ := json.Marshal(registerReq)

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().Register(mock.Anything, &registerReq).Return(registerResp, nil)

	s.router.POST("/auth/register", s.controller.Register)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["status"])
	assert.NotNil(s.T(), response["data"])
}

// Test Register - Invalid Email
func (s *UserControllerTestSuite) TestRegister_ValidationError_InvalidEmail() {
	// Arrange
	registerReq := dto.RegisterRequest{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "password123",
		ConfirmPassword: "password123",
		Email:           "invalid-email",
	}

	reqBody, _ := json.Marshal(registerReq)

	s.router.POST("/auth/register", s.controller.Register)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusUnprocessableEntity, w.Code)
}

// Test Register - Username Already Exists
func (s *UserControllerTestSuite) TestRegister_UsernameExists() {
	// Arrange
	registerReq := dto.RegisterRequest{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "password123",
		ConfirmPassword: "password123",
		Email:           "john@example.com",
	}

	reqBody, _ := json.Marshal(registerReq)

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().Register(mock.Anything, &registerReq).Return(nil, errConstant.ErrUsernameExist)

	s.router.POST("/auth/register", s.controller.Register)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// Test Register - Missing Required Fields
func (s *UserControllerTestSuite) TestRegister_ValidationError_MissingFields() {
	// Arrange
	registerReq := dto.RegisterRequest{
		Username: "johndoe",
		// Missing Name, Password, ConfirmPassword, Email
	}

	reqBody, _ := json.Marshal(registerReq)

	s.router.POST("/auth/register", s.controller.Register)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusUnprocessableEntity, w.Code)
}

// ==================== Update Tests ====================

// Test Update - Success
func (s *UserControllerTestSuite) TestUpdate_Success() {
	// Arrange
	userUUID := uuid.New().String()
	updateReq := dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johnupdated",
		Email:    "johnupdated@example.com",
	}

	updatedUser := &dto.UserResponse{
		UUID:     uuid.MustParse(userUUID),
		Name:     "John Updated",
		Username: "johnupdated",
		Email:    "johnupdated@example.com",
	}

	reqBody, _ := json.Marshal(updateReq)

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().Update(mock.Anything, &updateReq, userUUID).Return(updatedUser, nil)

	s.router.PUT("/auth/:uuid", s.controller.Update)

	req := httptest.NewRequest(http.MethodPut, "/auth/"+userUUID, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["status"])
	assert.NotNil(s.T(), response["data"])
}

// Test Update - Invalid Email
func (s *UserControllerTestSuite) TestUpdate_ValidationError_InvalidEmail() {
	// Arrange
	userUUID := uuid.New().String()
	updateReq := dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johnupdated",
		Email:    "invalid-email",
	}

	reqBody, _ := json.Marshal(updateReq)

	s.router.PUT("/auth/:uuid", s.controller.Update)

	req := httptest.NewRequest(http.MethodPut, "/auth/"+userUUID, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusUnprocessableEntity, w.Code)
}

// Test Update - User Not Found
func (s *UserControllerTestSuite) TestUpdate_UserNotFound() {
	// Arrange
	userUUID := uuid.New().String()
	updateReq := dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johnupdated",
		Email:    "johnupdated@example.com",
	}

	reqBody, _ := json.Marshal(updateReq)

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().Update(mock.Anything, &updateReq, userUUID).Return(nil, errConstant.ErrUserNotFound)

	s.router.PUT("/auth/:uuid", s.controller.Update)

	req := httptest.NewRequest(http.MethodPut, "/auth/"+userUUID, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// Test Update - Invalid JSON
func (s *UserControllerTestSuite) TestUpdate_InvalidJSON() {
	// Arrange
	userUUID := uuid.New().String()
	invalidJSON := []byte(`{"name": "John Updated", "email":}`)

	s.router.PUT("/auth/:uuid", s.controller.Update)

	req := httptest.NewRequest(http.MethodPut, "/auth/"+userUUID, bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// ==================== GetUserLogin Tests ====================

// Test GetUserLogin - Success
func (s *UserControllerTestSuite) TestGetUserLogin_Success() {
	// Arrange
	userUUID := uuid.New()
	userResp := &dto.UserResponse{
		UUID:     userUUID,
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Role:     "user",
	}

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().GetUserLogin(mock.Anything).Return(userResp, nil)

	s.router.GET("/auth/user", s.controller.GetUserLogin)

	req := httptest.NewRequest(http.MethodGet, "/auth/user", nil)
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["status"])
	assert.NotNil(s.T(), response["data"])
}

// Test GetUserLogin - Service Error
func (s *UserControllerTestSuite) TestGetUserLogin_ServiceError() {
	// Arrange
	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().GetUserLogin(mock.Anything).Return(nil, errors.New("service error"))

	s.router.GET("/auth/user", s.controller.GetUserLogin)

	req := httptest.NewRequest(http.MethodGet, "/auth/user", nil)
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// ==================== GetUserByUUID Tests ====================

// Test GetUserByUUID - Success
func (s *UserControllerTestSuite) TestGetUserByUUID_Success() {
	// Arrange
	userUUID := uuid.New().String()
	userResp := &dto.UserResponse{
		UUID:     uuid.MustParse(userUUID),
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
	}

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().GetUserByUUID(mock.Anything, userUUID).Return(userResp, nil)

	s.router.GET("/auth/:uuid", s.controller.GetUserByUUID)

	req := httptest.NewRequest(http.MethodGet, "/auth/"+userUUID, nil)
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "success", response["status"])
	assert.NotNil(s.T(), response["data"])
}

// Test GetUserByUUID - User Not Found
func (s *UserControllerTestSuite) TestGetUserByUUID_UserNotFound() {
	// Arrange
	userUUID := uuid.New().String()

	s.mockServiceReg.EXPECT().GetUser().Return(s.mockUserService)
	s.mockUserService.EXPECT().GetUserByUUID(mock.Anything, userUUID).Return(nil, errConstant.ErrUserNotFound)

	s.router.GET("/auth/:uuid", s.controller.GetUserByUUID)

	req := httptest.NewRequest(http.MethodGet, "/auth/"+userUUID, nil)
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// ==================== Edge Cases ====================

// Test Login - Empty Request Body
func (s *UserControllerTestSuite) TestLogin_EmptyRequestBody() {
	// Arrange
	s.router.POST("/auth/login", s.controller.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusUnprocessableEntity, w.Code)
}

// Test Register - Empty Request Body
func (s *UserControllerTestSuite) TestRegister_EmptyRequestBody() {
	// Arrange
	s.router.POST("/auth/register", s.controller.Register)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusUnprocessableEntity, w.Code)
}

// Test Update - Empty Request Body
func (s *UserControllerTestSuite) TestUpdate_EmptyRequestBody() {
	// Arrange
	userUUID := uuid.New().String()
	s.router.PUT("/auth/:uuid", s.controller.Update)

	req := httptest.NewRequest(http.MethodPut, "/auth/"+userUUID, bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(s.T(), http.StatusUnprocessableEntity, w.Code)
}
