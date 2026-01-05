# Testing Documentation - Manabu Service

Dokumentasi lengkap mengenai unit testing untuk User API di manabu-service.

## Table of Contents
- [Overview](#overview)
- [Testing Stack](#testing-stack)
- [Project Structure](#project-structure)
- [Mocking Strategy](#mocking-strategy)
- [Running Tests](#running-tests)
- [Test Coverage](#test-coverage)
- [Best Practices](#best-practices)

## Overview

Project ini menggunakan comprehensive unit testing dengan coverage tinggi untuk semua layer aplikasi:
- **Controller Layer**: HTTP request/response testing
- **Service Layer**: Business logic testing
- **Repository Layer**: Database interaction testing

## Testing Stack

### Dependencies
```bash
# Testing Framework
github.com/stretchr/testify         # Assertions dan test suite
github.com/vektra/mockery/v2        # Mock generator

# Database Mocking
github.com/DATA-DOG/go-sqlmock      # SQL mock untuk testing GORM
```

### Installation
```bash
# Install dependencies
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock

# Install mockery tool
go install github.com/vektra/mockery/v2@latest

# Verify installation
mockery --version
```

## Project Structure

```
manabu-service/
├── .mockery.yaml                    # Mockery configuration
├── controllers/
│   └── user/
│       ├── user.go                  # User controller implementation
│       └── user_test.go             # Controller tests (96.8% coverage)
├── services/
│   └── user/
│       ├── user.go                  # User service implementation
│       └── user_test.go             # Service tests (90.0% coverage)
├── repositories/
│   └── user/
│       ├── user.go                  # User repository implementation
│       └── user_test.go             # Repository tests (93.8% coverage)
└── mocks/                           # Auto-generated mocks
    ├── repositories/
    │   ├── mock_i_repository_registry.go
    │   └── user/
    │       └── mock_i_user_repository.go
    └── services/
        ├── mock_i_service_registry.go
        └── user/
            └── mock_i_user_service.go
```

## Mocking Strategy

### What to Mock and Why

#### 1. Repository Layer Testing
**Mocking**: Database (GORM)
**Tool**: `go-sqlmock`
**Why**: Kami tidak ingin bergantung pada database real saat testing. Dengan mock, kita bisa:
- Test lebih cepat (no I/O overhead)
- Kontrol penuh atas database responses
- Test error scenarios yang sulit di-reproduce dengan database real

**Example**:
```go
// Setup mock database
sqlDB, mock, err := sqlmock.New()
dialector := postgres.New(postgres.Config{
    Conn:       sqlDB,
    DriverName: "postgres",
})
db, err := gorm.Open(dialector, &gorm.Config{})

// Mock query expectation
mock.ExpectQuery(`SELECT \* FROM "users"`).
    WithArgs(username, 1).
    WillReturnRows(userRows)
```

#### 2. Service Layer Testing
**Mocking**: Repository Layer (IUserRepository, IRepositoryRegistry)
**Tool**: `mockery`
**Why**: Service layer hanya perlu test business logic, bukan database operations:
- Isolate business logic dari data access
- Test business rules tanpa database dependency
- Mock repository untuk test berbagai scenarios (success, error, edge cases)

**Example**:
```go
// Setup mocks
mockRepoRegistry := mockRepo.NewMockIRepositoryRegistry(t)
mockUserRepo := mockUserRepo.NewMockIUserRepository(t)

// Set expectations
mockRepoRegistry.EXPECT().GetUser().Return(mockUserRepo)
mockUserRepo.EXPECT().FindByUsername(ctx, username).Return(user, nil)
```

#### 3. Controller Layer Testing
**Mocking**: Service Layer (IUserService, IServiceRegistry)
**Tool**: `mockery` + `httptest`
**Why**: Controller hanya perlu test HTTP handling, bukan business logic:
- Test HTTP request/response flow
- Validate JSON binding dan validation
- Test error responses dan status codes
- Mock service untuk isolate HTTP logic

**Example**:
```go
// Setup mocks
mockServiceReg := mockService.NewMockIServiceRegistry(t)
mockUserService := mockUserService.NewMockIUserService(t)

// Set expectations
mockServiceReg.EXPECT().GetUser().Return(mockUserService)
mockUserService.EXPECT().Login(mock.Anything, &loginReq).Return(loginResp, nil)

// Test HTTP request
req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
```

### Mockery Configuration

File `.mockery.yaml` mengatur auto-generation untuk semua interfaces:

```yaml
with-expecter: true
dir: "mocks/{{.InterfaceDirRelative}}"
mockname: "Mock{{.InterfaceName}}"
outpkg: "{{.PackageName}}"
filename: "mock_{{.InterfaceNameSnake}}.go"

packages:
  # Repository interfaces
  manabu-service/repositories/user:
    interfaces:
      IUserRepository:
        config:
          with-expecter: true

  # Service interfaces
  manabu-service/services/user:
    interfaces:
      IUserService:
        config:
          with-expecter: true

  # Registry interfaces
  manabu-service/repositories:
    interfaces:
      IRepositoryRegistry:

  manabu-service/services:
    interfaces:
      IServiceRegistry:
```

### Regenerate Mocks

Setelah mengubah interface, regenerate mocks dengan:
```bash
mockery
```

## Running Tests

### Run All Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific package
go test ./controllers/user/... -v
go test ./services/user/... -v
go test ./repositories/user/... -v
```

### Run Specific Test
```bash
# Run specific test suite
go test ./controllers/user/... -v -run TestUserControllerTestSuite

# Run specific test case
go test ./controllers/user/... -v -run "TestUserControllerTestSuite/TestLogin_Success"
```

### Generate Coverage Report
```bash
# Generate coverage profile
go test ./controllers/user/... ./services/user/... ./repositories/user/... -coverprofile=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Test Coverage

### Current Coverage Stats

| Layer | Package | Coverage | Test Count |
|-------|---------|----------|------------|
| Controller | `controllers/user` | **96.8%** | 19 tests |
| Service | `services/user` | **90.0%** | 17 tests |
| Repository | `repositories/user` | **93.8%** | 12 tests |

### Coverage Breakdown

#### UserController Tests (19 tests)
- ✅ Login - Success
- ✅ Login - Invalid JSON
- ✅ Login - Validation Error (Missing Username)
- ✅ Login - Service Error (User Not Found)
- ✅ Login - Empty Request Body
- ✅ Register - Success
- ✅ Register - Validation Error (Invalid Email)
- ✅ Register - Validation Error (Missing Fields)
- ✅ Register - Username Already Exists
- ✅ Register - Empty Request Body
- ✅ Update - Success
- ✅ Update - Validation Error (Invalid Email)
- ✅ Update - User Not Found
- ✅ Update - Invalid JSON
- ✅ Update - Empty Request Body
- ✅ GetUserLogin - Success
- ✅ GetUserLogin - Service Error
- ✅ GetUserByUUID - Success
- ✅ GetUserByUUID - User Not Found

#### UserService Tests (17 tests)
- ✅ Login - Success (with JWT generation)
- ✅ Login - User Not Found
- ✅ Login - Wrong Password
- ✅ Register - Success (with password hashing)
- ✅ Register - Username Already Exists
- ✅ Register - Email Already Exists
- ✅ Register - Password Mismatch
- ✅ Register - Repository Error
- ✅ Update - Success (with password hashing)
- ✅ Update - User Not Found
- ✅ Update - Username Exists (different user)
- ✅ Update - Email Exists (different user)
- ✅ Update - Password Mismatch
- ✅ Update - Without Password
- ✅ GetUserLogin - Success
- ✅ GetUserByUUID - Success
- ✅ GetUserByUUID - Not Found

#### UserRepository Tests (12 tests)
- ✅ Register - Success
- ✅ Register - SQL Error
- ✅ FindByUsername - Success
- ✅ FindByUsername - Not Found
- ✅ FindByUsername - SQL Error
- ✅ FindByEmail - Success
- ✅ FindByEmail - Not Found
- ✅ FindByUUID - Success
- ✅ FindByUUID - Not Found
- ✅ Update - Success
- ✅ Update - SQL Error
- ✅ Update - Empty Password (edge case)

## Best Practices

### 1. Test Organization
- Gunakan **testify suite** untuk organize tests
- Setup dan teardown yang konsisten
- Group related tests dalam satu suite

```go
type UserServiceTestSuite struct {
    suite.Suite
    service          IUserService
    mockRepoRegistry *mockRepo.MockIRepositoryRegistry
    mockUserRepo     *mockUserRepo.MockIUserRepository
    ctx              context.Context
}

func (s *UserServiceTestSuite) SetupTest() {
    // Setup untuk setiap test
    s.mockRepoRegistry = mockRepo.NewMockIRepositoryRegistry(s.T())
    s.mockUserRepo = mockUserRepo.NewMockIUserRepository(s.T())
    s.service = NewUserService(s.mockRepoRegistry)
    s.ctx = context.Background()
}
```

### 2. Test Naming Convention
- Format: `Test<Function>_<Scenario>`
- Jelas dan deskriptif
- Examples:
  - `TestLogin_Success`
  - `TestRegister_UsernameExists`
  - `TestUpdate_PasswordMismatch`

### 3. Test Structure (AAA Pattern)
```go
func (s *UserServiceTestSuite) TestLogin_Success() {
    // Arrange - Setup data dan expectations
    req := &dto.LoginRequest{
        Username: "johndoe",
        Password: "password123",
    }
    mockUser := &models.User{...}
    s.mockUserRepo.EXPECT().FindByUsername(s.ctx, req.Username).Return(mockUser, nil)

    // Act - Execute function being tested
    result, err := s.service.Login(s.ctx, req)

    // Assert - Verify results
    assert.NoError(s.T(), err)
    assert.NotNil(s.T(), result)
    assert.Equal(s.T(), mockUser.UUID, result.User.UUID)
}
```

### 4. Test Coverage Goals
- **Minimum**: 80% coverage untuk semua packages
- **Target**: 90%+ coverage
- **Critical paths**: 100% coverage (authentication, authorization)

### 5. Mock Best Practices
- Mock external dependencies (database, HTTP clients, third-party APIs)
- Don't mock internal logic
- Use `mock.Anything` untuk parameters yang tidak penting
- Be specific dengan expectations ketika values matter

### 6. Testing Edge Cases
Selalu test:
- ✅ Happy path scenarios
- ✅ Error scenarios
- ✅ Validation errors
- ✅ Empty inputs
- ✅ Nil values
- ✅ Boundary conditions

### 7. Database Testing with GORM
- Always use `sqlmock.AnyArg()` untuk dynamic values (timestamps, UUIDs)
- Remember GORM adds `LIMIT 1` untuk `.First()` queries
- Account untuk `Preload()` queries

Example:
```go
// GORM will add LIMIT 1 automatically
s.sqlMock.ExpectQuery(`SELECT \* FROM "users"`).
    WithArgs(username, 1).  // Include the LIMIT value
    WillReturnRows(userRows)
```

### 8. HTTP Testing Best Practices
- Test all HTTP methods (GET, POST, PUT, DELETE)
- Test different status codes
- Validate JSON responses structure
- Test request validation

### 9. Continuous Testing
```bash
# Watch mode (using external tool like entr)
ls **/*.go | entr -c go test ./...

# Pre-commit hook
git add .git/hooks/pre-commit
#!/bin/bash
go test ./... -cover
```

## Common Issues & Solutions

### Issue 1: Mock expectations not met
**Problem**: `Received unexpected error: there is a remaining expectation`
**Solution**: Ensure all EXPECT() calls are matched with actual function calls

### Issue 2: GORM query arguments mismatch
**Problem**: `arguments do not match: expected 1, but got 2 arguments`
**Solution**: GORM adds LIMIT automatically. Include it in WithArgs:
```go
// Wrong
WithArgs(username)
// Correct
WithArgs(username, 1)
```

### Issue 3: Context matching issues
**Problem**: Context values don't match in assertions
**Solution**: Use `mock.Anything` untuk context parameters:
```go
mockUserService.EXPECT().Login(mock.Anything, &loginReq).Return(loginResp, nil)
```

## Additional Resources

- [Testify Documentation](https://github.com/stretchr/testify)
- [Mockery Documentation](https://vektra.github.io/mockery/)
- [Go-SQLMock Documentation](https://github.com/DATA-DOG/go-sqlmock)
- [Go Testing Best Practices](https://golang.org/doc/tutorial/add-a-test)

## Contributing

Saat menambahkan features baru:

1. **Write tests first** (TDD approach)
2. Ensure **90%+ coverage** untuk code baru
3. Update mocks jika ada perubahan interface
4. Run all tests sebelum commit
5. Update dokumentasi jika perlu

```bash
# Workflow
mockery                                # Regenerate mocks
go test ./... -cover                   # Run all tests
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out       # Check coverage
```

---

**Last Updated**: 2026-01-05
**Maintainer**: Manabu Development Team
