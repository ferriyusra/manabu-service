package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"manabu-service/constants"
	errConstant "manabu-service/constants/error"
	"manabu-service/domain/dto"
	"manabu-service/domain/models"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	repository IUserRepository
	db         *gorm.DB
	sqlMock    sqlmock.Sqlmock
	ctx        context.Context
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) SetupTest() {
	var err error
	sqlDB, mock, err := sqlmock.New()
	s.Require().NoError(err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	s.db, err = gorm.Open(dialector, &gorm.Config{})
	s.Require().NoError(err)

	s.sqlMock = mock
	s.repository = NewUserRepository(s.db)
	s.ctx = context.Background()
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	s.Require().NoError(s.sqlMock.ExpectationsWereMet())
}

// Test Register - Happy Path
func (s *UserRepositoryTestSuite) TestRegister_Success() {
	// Arrange
	req := &dto.RegisterRequest{
		Name:     "John Doe",
		Username: "johndoe",
		Password: "hashedpassword",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(
			sqlmock.AnyArg(), // UUID
			req.Name,
			req.Username,
			req.Password,
			req.Email,
			req.RoleID,
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.sqlMock.ExpectCommit()

	// Act
	result, err := s.repository.Register(s.ctx, req)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Name, result.Name)
	assert.Equal(s.T(), req.Username, result.Username)
	assert.Equal(s.T(), req.Email, result.Email)
	assert.Equal(s.T(), req.RoleID, result.RoleID)
	assert.NotEqual(s.T(), uuid.Nil, result.UUID)
}

// Test Register - SQL Error
func (s *UserRepositoryTestSuite) TestRegister_SQLError() {
	// Arrange
	req := &dto.RegisterRequest{
		Name:     "John Doe",
		Username: "johndoe",
		Password: "hashedpassword",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectQuery(`INSERT INTO "users"`).
		WillReturnError(errors.New("database error"))
	s.sqlMock.ExpectRollback()

	// Act
	result, err := s.repository.Register(s.ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "database server failed to execute query")
}

// Test FindByUsername - Success
func (s *UserRepositoryTestSuite) TestFindByUsername_Success() {
	// Arrange
	username := "johndoe"
	userUUID := uuid.New()
	expectedUser := &models.User{
		ID:       1,
		UUID:     userUUID,
		Name:     "John Doe",
		Username: username,
		Password: "hashedpassword",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	userRows := sqlmock.NewRows([]string{
		"id", "uuid", "name", "username", "password", "email", "role_id",
	}).AddRow(
		expectedUser.ID,
		expectedUser.UUID,
		expectedUser.Name,
		expectedUser.Username,
		expectedUser.Password,
		expectedUser.Email,
		expectedUser.RoleID,
	)

	roleRows := sqlmock.NewRows([]string{
		"id", "code", "name",
	}).AddRow(
		constants.User,
		"USER",
		"User",
	)

	s.sqlMock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(username, 1).
		WillReturnRows(userRows)

	s.sqlMock.ExpectQuery(`SELECT \* FROM "roles"`).
		WithArgs(constants.User).
		WillReturnRows(roleRows)

	// Act
	result, err := s.repository.FindByUsername(s.ctx, username)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.UUID, result.UUID)
	assert.Equal(s.T(), expectedUser.Username, result.Username)
	assert.Equal(s.T(), expectedUser.Email, result.Email)
}

// Test FindByUsername - Not Found
func (s *UserRepositoryTestSuite) TestFindByUsername_NotFound() {
	// Arrange
	username := "nonexistent"

	s.sqlMock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(username, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := s.repository.FindByUsername(s.ctx, username)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrUserNotFound, err)
}

// Test FindByUsername - SQL Error
func (s *UserRepositoryTestSuite) TestFindByUsername_SQLError() {
	// Arrange
	username := "johndoe"

	s.sqlMock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(username, 1).
		WillReturnError(errors.New("database error"))

	// Act
	result, err := s.repository.FindByUsername(s.ctx, username)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "database server failed to execute query")
}

// Test FindByEmail - Success
func (s *UserRepositoryTestSuite) TestFindByEmail_Success() {
	// Arrange
	email := "john@example.com"
	userUUID := uuid.New()
	expectedUser := &models.User{
		ID:       1,
		UUID:     userUUID,
		Name:     "John Doe",
		Username: "johndoe",
		Password: "hashedpassword",
		Email:    email,
		RoleID:   constants.User,
	}

	userRows := sqlmock.NewRows([]string{
		"id", "uuid", "name", "username", "password", "email", "role_id",
	}).AddRow(
		expectedUser.ID,
		expectedUser.UUID,
		expectedUser.Name,
		expectedUser.Username,
		expectedUser.Password,
		expectedUser.Email,
		expectedUser.RoleID,
	)

	roleRows := sqlmock.NewRows([]string{
		"id", "code", "name",
	}).AddRow(
		constants.User,
		"USER",
		"User",
	)

	s.sqlMock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(email, 1).
		WillReturnRows(userRows)

	s.sqlMock.ExpectQuery(`SELECT \* FROM "roles"`).
		WithArgs(constants.User).
		WillReturnRows(roleRows)

	// Act
	result, err := s.repository.FindByEmail(s.ctx, email)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.Email, result.Email)
}

// Test FindByEmail - Not Found
func (s *UserRepositoryTestSuite) TestFindByEmail_NotFound() {
	// Arrange
	email := "nonexistent@example.com"

	s.sqlMock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(email, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := s.repository.FindByEmail(s.ctx, email)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrUserNotFound, err)
}

// Test FindByUUID - Success
func (s *UserRepositoryTestSuite) TestFindByUUID_Success() {
	// Arrange
	userUUID := uuid.New()
	expectedUser := &models.User{
		ID:       1,
		UUID:     userUUID,
		Name:     "John Doe",
		Username: "johndoe",
		Password: "hashedpassword",
		Email:    "john@example.com",
		RoleID:   constants.User,
	}

	userRows := sqlmock.NewRows([]string{
		"id", "uuid", "name", "username", "password", "email", "role_id",
	}).AddRow(
		expectedUser.ID,
		expectedUser.UUID,
		expectedUser.Name,
		expectedUser.Username,
		expectedUser.Password,
		expectedUser.Email,
		expectedUser.RoleID,
	)

	roleRows := sqlmock.NewRows([]string{
		"id", "code", "name",
	}).AddRow(
		constants.User,
		"USER",
		"User",
	)

	s.sqlMock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(userUUID.String(), 1).
		WillReturnRows(userRows)

	s.sqlMock.ExpectQuery(`SELECT \* FROM "roles"`).
		WithArgs(constants.User).
		WillReturnRows(roleRows)

	// Act
	result, err := s.repository.FindByUUID(s.ctx, userUUID.String())

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.UUID, result.UUID)
}

// Test FindByUUID - Not Found
func (s *UserRepositoryTestSuite) TestFindByUUID_NotFound() {
	// Arrange
	userUUID := uuid.New()

	s.sqlMock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(userUUID.String(), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := s.repository.FindByUUID(s.ctx, userUUID.String())

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), errConstant.ErrUserNotFound, err)
}

// Test Update - Success
func (s *UserRepositoryTestSuite) TestUpdate_Success() {
	// Arrange
	userUUID := uuid.New().String()
	password := "newhashedpassword"
	req := &dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johnupdated",
		Password: &password,
		Email:    "johnupdated@example.com",
	}

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectExec(`UPDATE "users"`).
		WithArgs(
			req.Name,
			req.Username,
			*req.Password,
			req.Email,
			sqlmock.AnyArg(), // UpdatedAt
			userUUID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.sqlMock.ExpectCommit()

	// Act
	result, err := s.repository.Update(s.ctx, req, userUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Name, result.Name)
	assert.Equal(s.T(), req.Username, result.Username)
	assert.Equal(s.T(), req.Email, result.Email)
}

// Test Update - SQL Error
func (s *UserRepositoryTestSuite) TestUpdate_SQLError() {
	// Arrange
	userUUID := uuid.New().String()
	password := "newhashedpassword"
	req := &dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johnupdated",
		Password: &password,
		Email:    "johnupdated@example.com",
	}

	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectExec(`UPDATE "users"`).
		WillReturnError(errors.New("database error"))
	s.sqlMock.ExpectRollback()

	// Act
	result, err := s.repository.Update(s.ctx, req, userUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "database server failed to execute query")
}

// Test Update - Empty Password (edge case)
// Note: When password is empty string, GORM will not include it in the UPDATE statement
// because Updates() ignores zero values. This test verifies that behavior.
func (s *UserRepositoryTestSuite) TestUpdate_EmptyPassword() {
	// Arrange
	userUUID := uuid.New().String()
	emptyPassword := ""
	req := &dto.UpdateRequest{
		Name:     "John Updated",
		Username: "johnupdated",
		Password: &emptyPassword,
		Email:    "johnupdated@example.com",
	}

	// GORM Updates() ignores zero values, so empty password won't be in the query
	s.sqlMock.ExpectBegin()
	s.sqlMock.ExpectExec(`UPDATE "users"`).
		WithArgs(
			req.Name,
			req.Username,
			req.Email,
			sqlmock.AnyArg(), // UpdatedAt
			userUUID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.sqlMock.ExpectCommit()

	// Act
	result, err := s.repository.Update(s.ctx, req, userUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), emptyPassword, result.Password)
}
