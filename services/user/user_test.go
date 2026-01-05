package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"manabu-service/constants"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
	mockRepo "manabu-service/mocks/repositories"
	mockUserRepo "manabu-service/mocks/repositories/user"
)

type UserServiceTestSuite struct {
	suite.Suite
	service           IUserService
	mockRepoRegistry  *mockRepo.MockIRepositoryRegistry
	mockUserRepo      *mockUserRepo.MockIUserRepository
	ctx               context.Context
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) SetupTest() {
	s.mockRepoRegistry = mockRepo.NewMockIRepositoryRegistry(s.T())
	s.mockUserRepo = mockUserRepo.NewMockIUserRepository(s.T())
	s.service = NewUserService(s.mockRepoRegistry)
	s.ctx = context.Background()
}

// ==================== Login Tests ====================

// Test Login - Success
func (s *UserServiceTestSuite) TestLogin_Success() {
	// Arrange
	req := &dto.LoginRequest{
		Username: "johndoe",
		Password: "password123",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	userUUID := uuid.New()

	mockUser := &models.User{
		ID:       1,
		UUID:     userUUID,
		Name:     "John Doe",
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    "john@example.com",
		RoleID:   constants.User,
		Role: models.Role{
			ID:   constants.User,
			Code: "USER",
			Name: "User",
		},
	}

	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(mockUser, nil)

	// Act
	result, err := s.service.Login(s.ctx, req)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), mockUser.UUID, result.User.UUID)
	assert.Equal(s.T(), mockUser.Name, result.User.Name)
	assert.Equal(s.T(), mockUser.Username, result.User.Username)
	assert.Equal(s.T(), mockUser.Email, result.User.Email)
	assert.Equal(s.T(), "user", result.User.Role)
	assert.NotEmpty(s.T(), result.Token)
}

// Test Login - User Not Found
func (s *UserServiceTestSuite) TestLogin_UserNotFound() {
	// Arrange
	req := &dto.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Act
	result, err := s.service.Login(s.ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrUserNotFound, err)
}

// Test Login - Wrong Password
func (s *UserServiceTestSuite) TestLogin_WrongPassword() {
	// Arrange
	req := &dto.LoginRequest{
		Username: "johndoe",
		Password: "wrongpassword",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	userUUID := uuid.New()

	mockUser := &models.User{
		ID:       1,
		UUID:     userUUID,
		Name:     "John Doe",
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    "john@example.com",
		RoleID:   constants.User,
		Role: models.Role{
			ID:   constants.User,
			Code: "USER",
			Name: "User",
		},
	}

	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(mockUser, nil)

	// Act
	result, err := s.service.Login(s.ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), bcrypt.ErrMismatchedHashAndPassword, err)
}

// ==================== Register Tests ====================

// Test Register - Success
func (s *UserServiceTestSuite) TestRegister_Success() {
	// Arrange
	req := &dto.RegisterRequest{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "password123",
		ConfirmPassword: "password123",
		Email:           "john@example.com",
	}

	userUUID := uuid.New()
	mockUser := &models.User{
		ID:       1,
		UUID:     userUUID,
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		RoleID:   constants.User,
	}

	// Mock username check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Mock email check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByEmail(s.ctx, req.Email).Return(nil, errConstant.ErrUserNotFound)

	// Mock register
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().Register(s.ctx, mock.AnythingOfType("*dto.RegisterRequest")).Return(mockUser, nil)

	// Act
	result, err := s.service.Register(s.ctx, req)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), mockUser.UUID, result.User.UUID)
	assert.Equal(s.T(), mockUser.Name, result.User.Name)
	assert.Equal(s.T(), mockUser.Username, result.User.Username)
	assert.Equal(s.T(), mockUser.Email, result.User.Email)
}

// Test Register - Username Already Exists
func (s *UserServiceTestSuite) TestRegister_UsernameExists() {
	// Arrange
	req := &dto.RegisterRequest{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "password123",
		ConfirmPassword: "password123",
		Email:           "john@example.com",
	}

	existingUser := &models.User{
		ID:       1,
		UUID:     uuid.New(),
		Username: req.Username,
	}

	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(existingUser, nil)

	// Act
	result, err := s.service.Register(s.ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrUsernameExist, err)
}

// Test Register - Email Already Exists
func (s *UserServiceTestSuite) TestRegister_EmailExists() {
	// Arrange
	req := &dto.RegisterRequest{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "password123",
		ConfirmPassword: "password123",
		Email:           "john@example.com",
	}

	existingUser := &models.User{
		ID:    1,
		UUID:  uuid.New(),
		Email: req.Email,
	}

	// Mock username check (pass)
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Mock email check (fail - exists)
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByEmail(s.ctx, req.Email).Return(existingUser, nil)

	// Act
	result, err := s.service.Register(s.ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrEmailExist, err)
}

// Test Register - Password Does Not Match
func (s *UserServiceTestSuite) TestRegister_PasswordMismatch() {
	// Arrange
	req := &dto.RegisterRequest{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "password123",
		ConfirmPassword: "differentpassword",
		Email:           "john@example.com",
	}

	// Mock username check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Mock email check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByEmail(s.ctx, req.Email).Return(nil, errConstant.ErrUserNotFound)

	// Act
	result, err := s.service.Register(s.ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrPasswordDoesNotMatch, err)
}

// Test Register - Repository Error
func (s *UserServiceTestSuite) TestRegister_RepositoryError() {
	// Arrange
	req := &dto.RegisterRequest{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "password123",
		ConfirmPassword: "password123",
		Email:           "john@example.com",
	}

	// Mock username check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Mock email check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByEmail(s.ctx, req.Email).Return(nil, errConstant.ErrUserNotFound)

	// Mock register error
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().Register(s.ctx, mock.AnythingOfType("*dto.RegisterRequest")).Return(nil, errors.New("database error"))

	// Act
	result, err := s.service.Register(s.ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

// ==================== Update Tests ====================

// Test Update - Success
func (s *UserServiceTestSuite) TestUpdate_Success() {
	// Arrange
	userUUID := uuid.New().String()
	password := "newpassword123"
	confirmPassword := "newpassword123"

	req := &dto.UpdateRequest{
		Name:            "John Updated",
		Username:        "johnupdated",
		Password:        &password,
		ConfirmPassword: &confirmPassword,
		Email:           "johnupdated@example.com",
	}

	existingUser := &models.User{
		ID:       1,
		UUID:     uuid.MustParse(userUUID),
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	updatedUser := &models.User{
		ID:       1,
		UUID:     uuid.MustParse(userUUID),
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		RoleID:   constants.User,
	}

	// Mock find by UUID
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUUID(s.ctx, userUUID).Return(existingUser, nil)

	// Mock username check (not exist)
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Mock email check (not exist)
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByEmail(s.ctx, req.Email).Return(nil, errConstant.ErrUserNotFound)

	// Mock update
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().Update(s.ctx, mock.AnythingOfType("*dto.UpdateRequest"), userUUID).Return(updatedUser, nil)

	// Act
	result, err := s.service.Update(s.ctx, req, userUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), updatedUser.UUID, result.UUID)
	assert.Equal(s.T(), updatedUser.Name, result.Name)
	assert.Equal(s.T(), updatedUser.Username, result.Username)
	assert.Equal(s.T(), updatedUser.Email, result.Email)
}

// Test Update - User Not Found
func (s *UserServiceTestSuite) TestUpdate_UserNotFound() {
	// Arrange
	userUUID := uuid.New().String()
	req := &dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johnupdated",
		Email:    "johnupdated@example.com",
	}

	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUUID(s.ctx, userUUID).Return(nil, errConstant.ErrUserNotFound)

	// Act
	result, err := s.service.Update(s.ctx, req, userUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrUserNotFound, err)
}

// Test Update - Username Already Exists (different user)
func (s *UserServiceTestSuite) TestUpdate_UsernameExistsForDifferentUser() {
	// Arrange
	userUUID := uuid.New().String()
	req := &dto.UpdateRequest{
		Name:     "John Updated",
		Username: "existinguser",
		Email:    "john@example.com",
	}

	existingUser := &models.User{
		ID:       1,
		UUID:     uuid.MustParse(userUUID),
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	anotherUser := &models.User{
		ID:       2,
		UUID:     uuid.New(),
		Username: req.Username,
	}

	// Mock find by UUID
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUUID(s.ctx, userUUID).Return(existingUser, nil)

	// Mock username check (exists for different user)
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(anotherUser, nil).Times(2)

	// Act
	result, err := s.service.Update(s.ctx, req, userUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrUsernameExist, err)
}

// Test Update - Email Already Exists (different user)
func (s *UserServiceTestSuite) TestUpdate_EmailExistsForDifferentUser() {
	// Arrange
	userUUID := uuid.New().String()
	req := &dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johndoe",
		Email:    "existing@example.com",
	}

	existingUser := &models.User{
		ID:       1,
		UUID:     uuid.MustParse(userUUID),
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	anotherUser := &models.User{
		ID:    2,
		UUID:  uuid.New(),
		Email: req.Email,
	}

	// Mock find by UUID
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUUID(s.ctx, userUUID).Return(existingUser, nil)

	// Mock username check (same username, no conflict)
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Mock email check (exists for different user)
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByEmail(s.ctx, req.Email).Return(anotherUser, nil).Times(2)

	// Act
	result, err := s.service.Update(s.ctx, req, userUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrEmailExist, err)
}

// Test Update - Password Mismatch
func (s *UserServiceTestSuite) TestUpdate_PasswordMismatch() {
	// Arrange
	userUUID := uuid.New().String()
	password := "newpassword123"
	confirmPassword := "differentpassword"

	req := &dto.UpdateRequest{
		Name:            "John Updated",
		Username:        "johndoe",
		Password:        &password,
		ConfirmPassword: &confirmPassword,
		Email:           "john@example.com",
	}

	existingUser := &models.User{
		ID:       1,
		UUID:     uuid.MustParse(userUUID),
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	// Mock find by UUID
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUUID(s.ctx, userUUID).Return(existingUser, nil)

	// Mock username check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Mock email check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByEmail(s.ctx, req.Email).Return(nil, errConstant.ErrUserNotFound)

	// Act
	result, err := s.service.Update(s.ctx, req, userUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrPasswordDoesNotMatch, err)
}

// Test Update - Without Password
func (s *UserServiceTestSuite) TestUpdate_WithoutPassword() {
	// Arrange
	userUUID := uuid.New().String()

	req := &dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johnupdated",
		Password: nil,
		Email:    "johnupdated@example.com",
	}

	existingUser := &models.User{
		ID:       1,
		UUID:     uuid.MustParse(userUUID),
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	updatedUser := &models.User{
		ID:       1,
		UUID:     uuid.MustParse(userUUID),
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		RoleID:   constants.User,
	}

	// Mock find by UUID
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUUID(s.ctx, userUUID).Return(existingUser, nil)

	// Mock username check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(nil, errConstant.ErrUserNotFound)

	// Mock email check
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByEmail(s.ctx, req.Email).Return(nil, errConstant.ErrUserNotFound)

	// Mock update
	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().Update(s.ctx, mock.AnythingOfType("*dto.UpdateRequest"), userUUID).Return(updatedUser, nil)

	// Act
	result, err := s.service.Update(s.ctx, req, userUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
}

// ==================== GetUserLogin Tests ====================

// Test GetUserLogin - Success
func (s *UserServiceTestSuite) TestGetUserLogin_Success() {
	// Arrange
	userUUID := uuid.New()
	userLogin := &dto.UserResponse{
		UUID:     userUUID,
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Role:     "user",
	}

	ctx := context.WithValue(s.ctx, constants.UserLogin, userLogin)

	// Act
	result, err := s.service.GetUserLogin(ctx)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), userLogin.UUID, result.UUID)
	assert.Equal(s.T(), userLogin.Name, result.Name)
	assert.Equal(s.T(), userLogin.Username, result.Username)
	assert.Equal(s.T(), userLogin.Email, result.Email)
	assert.Equal(s.T(), userLogin.Role, result.Role)
}

// ==================== GetUserByUUID Tests ====================

// Test GetUserByUUID - Success
func (s *UserServiceTestSuite) TestGetUserByUUID_Success() {
	// Arrange
	userUUID := uuid.New()
	mockUser := &models.User{
		ID:       1,
		UUID:     userUUID,
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUUID(s.ctx, userUUID.String()).Return(mockUser, nil)

	// Act
	result, err := s.service.GetUserByUUID(s.ctx, userUUID.String())

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), mockUser.UUID, result.UUID)
	assert.Equal(s.T(), mockUser.Name, result.Name)
	assert.Equal(s.T(), mockUser.Username, result.Username)
	assert.Equal(s.T(), mockUser.Email, result.Email)
}

// Test GetUserByUUID - Not Found
func (s *UserServiceTestSuite) TestGetUserByUUID_NotFound() {
	// Arrange
	userUUID := uuid.New()

	s.mockRepoRegistry.EXPECT().GetUser().Return(s.mockUserRepo)
	s.mockUserRepo.EXPECT().FindByUUID(s.ctx, userUUID.String()).Return(nil, errConstant.ErrUserNotFound)

	// Act
	result, err := s.service.GetUserByUUID(s.ctx, userUUID.String())

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrUserNotFound, err)
}
