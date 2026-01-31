# Scentora Backend - Testing Guide

**Last Updated**: January 31, 2026  
**Test Coverage Target**: 80%+

---

## Overview

This document describes the testing strategy, tools, and best practices for the Scentora backend. **All code must be tested before committing.**

---

## Quick Start

### Running Tests

```bash
# Run all tests
./run-tests.sh

# Run specific package tests
go test ./internal/repository

# Run with coverage report
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Database

Tests use a separate `scentora_test` database to avoid affecting development data.

**Connection**: `postgres://admin:password@localhost:5435/scentora_test`

The test database is automatically created when running `./run-tests.sh`.

---

## Test Structure

### Test Organization

```
backend/
├── internal/
│   ├── config/
│   │   ├── database.go
│   │   └── database_test.go         # Database migration tests
│   ├── repository/
│   │   ├── user_repo.go
│   │   ├── user_repo_test.go        # Repository tests
│   │   ├── invitation_repo.go
│   │   └── invitation_repo_test.go
│   ├── services/
│   │   ├── auth_service.go
│   │   └── auth_service_test.go     # Service tests
│   ├── handlers/
│   │   ├── auth.go
│   │   └── auth_test.go             # Handler tests
│   └── testutil/
│       └── db.go                    # Test utilities
└── run-tests.sh                      # Test runner script
```

---

## Test Utilities

### testutil.SetupTestDB

Creates a test database connection and runs migrations:

```go
func TestMyFunction(t *testing.T) {
    tdb := testutil.SetupTestDB(t)
    defer tdb.Teardown(t)
    defer tdb.CleanupTables(t)
    
    // Your test code here
}
```

**Methods**:
- `SetupTestDB(t)` - Connect and migrate test database
- `Teardown(t)` - Close database connection
- `CleanupTables(t)` - Delete all data (for test isolation)

---

## Writing Tests

### Test Naming Convention

```go
// Function: CreateUser
// Tests:
func TestUserRepository_CreateUser(t *testing.T) {}          // Happy path
func TestUserRepository_CreateUserDuplicateEmail(t *testing.T) {}  // Error case
func TestUserRepository_CreateUserInvalidData(t *testing.T) {}    // Edge case
```

### Test Structure

```go
func TestFeatureName(t *testing.T) {
    // Setup
    tdb := testutil.SetupTestDB(t)
    defer tdb.Teardown(t)
    defer tdb.CleanupTables(t)
    
    // Arrange
    repo := repository.NewUserRepository(tdb.DB)
    user := &models.User{
        Email: "test@example.com",
        Username: "testuser",
        PasswordHash: "hashed",
    }
    
    // Act
    err := repo.Create(user)
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    if user.ID == "" {
        t.Error("Expected ID to be set")
    }
}
```

---

## Test Categories

### 1. Database Migration Tests

**File**: `internal/config/database_test.go`

Tests that database schema is created correctly:

- ✅ All tables exist
- ✅ Correct column types
- ✅ Constraints are enforced
- ✅ Indexes are created
- ✅ Generated columns work
- ✅ Cascade deletes work

**Example**:
```go
func TestDatabaseMigrations(t *testing.T) {
    tdb := testutil.SetupTestDB(t)
    defer tdb.Teardown(t)
    
    // Test that tables exist
    expectedTables := []string{"users", "accords", "accord_tags"}
    for _, table := range expectedTables {
        var exists bool
        err := tdb.DB.QueryRow(`
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_name = $1
            )
        `, table).Scan(&exists)
        
        if !exists {
            t.Errorf("Table %s does not exist", table)
        }
    }
}
```

### 2. Repository Tests

**Files**: `internal/repository/*_test.go`

Tests data access layer:

- ✅ CRUD operations
- ✅ Queries return correct data
- ✅ Unique constraints enforced
- ✅ Foreign key constraints enforced
- ✅ Error handling

**Example**:
```go
func TestUserRepository_Create(t *testing.T) {
    tdb := testutil.SetupTestDB(t)
    defer tdb.Teardown(t)
    defer tdb.CleanupTables(t)
    
    repo := repository.NewUserRepository(tdb.DB)
    user := &models.User{
        Email: "test@example.com",
        Username: "testuser",
        PasswordHash: "hashed",
    }
    
    err := repo.Create(user)
    if err != nil {
        t.Fatalf("Create failed: %v", err)
    }
    
    // Verify ID was set
    if user.ID == "" {
        t.Error("ID not set after create")
    }
}
```

### 3. Service Tests

**Files**: `internal/services/*_test.go`

Tests business logic layer:

- ✅ Business rules enforced
- ✅ Data validation
- ✅ Error handling
- ✅ Integration between repositories

**Example**:
```go
func TestAuthService_Register(t *testing.T) {
    // Mock repositories
    userRepo := &MockUserRepository{}
    tokenRepo := &MockTokenRepository{}
    inviteRepo := &MockInvitationRepository{}
    
    service := services.NewAuthService(userRepo, tokenRepo, inviteRepo, cfg)
    
    // Test registration logic
    result, err := service.Register("test@example.com", "testuser", "password", "INVITE123")
    
    if err != nil {
        t.Fatalf("Registration failed: %v", err)
    }
    if result.AccessToken == "" {
        t.Error("Expected access token")
    }
}
```

### 4. Handler Tests

**Files**: `internal/handlers/*_test.go`

Tests HTTP endpoints:

- ✅ Request parsing
- ✅ Response formatting
- ✅ Status codes
- ✅ Error responses
- ✅ Authentication

**Example**:
```go
func TestAuthHandler_Login(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{
        "email": "test@example.com",
        "password": "password"
    }`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    
    handler := handlers.NewAuthHandler(mockService)
    err := handler.Login(c)
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

---

## Test Data

### Creating Test Fixtures

```go
func createTestUser(t *testing.T, db *sqlx.DB) *models.User {
    user := &models.User{
        Email: fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
        Username: "testuser",
        PasswordHash: "hashed",
    }
    
    repo := repository.NewUserRepository(db)
    err := repo.Create(user)
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    
    return user
}
```

### Cleanup Between Tests

```go
func TestMultipleTests(t *testing.T) {
    tdb := testutil.SetupTestDB(t)
    defer tdb.Teardown(t)
    
    t.Run("Test1", func(t *testing.T) {
        defer tdb.CleanupTables(t)
        // Test code
    })
    
    t.Run("Test2", func(t *testing.T) {
        defer tdb.CleanupTables(t)
        // Test code
    })
}
```

---

## Coverage Requirements

### Minimum Coverage Targets

- **Overall**: 80%+
- **Repositories**: 90%+ (data layer is critical)
- **Services**: 85%+ (business logic must be solid)
- **Handlers**: 80%+ (API endpoints)
- **Models**: 100% (validation tests)

### Generating Coverage Reports

```bash
# Terminal report
go test -cover ./...

# HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Per-package report
go test -coverprofile=coverage.out ./internal/repository
go tool cover -func=coverage.out
```

---

## Best Practices

### ✅ DO

- Write tests BEFORE committing code
- Test both happy path and error cases
- Use descriptive test names
- Clean up test data (use `defer tdb.CleanupTables(t)`)
- Test edge cases and boundary conditions
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Keep tests fast (<5s for unit tests)

### ❌ DON'T

- Skip failing tests
- Comment out tests
- Write flaky tests (that pass sometimes)
- Test against production database
- Hard-code IDs or timestamps
- Leave test data in database
- Create tests without assertions
- Test framework code (e.g., Echo routing)

---

## Common Test Patterns

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "notanemail", true},
        {"empty email", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateEmail(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Subtests

```go
func TestUserOperations(t *testing.T) {
    tdb := testutil.SetupTestDB(t)
    defer tdb.Teardown(t)
    
    t.Run("Create", func(t *testing.T) {
        defer tdb.CleanupTables(t)
        // Create test
    })
    
    t.Run("Update", func(t *testing.T) {
        defer tdb.CleanupTables(t)
        // Update test
    })
}
```

---

## Continuous Integration

### Pre-Commit Checklist

Before committing:

1. ✅ Run `./run-tests.sh` - all tests pass
2. ✅ Check coverage - meets minimum thresholds
3. ✅ No commented-out tests
4. ✅ No skipped tests
5. ✅ New code has tests

### CI Pipeline (Future)

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: go test -v -cover ./...
```

---

## Debugging Tests

### Running Single Test

```bash
go test -v -run TestUserRepository_Create ./internal/repository
```

### Verbose Output

```bash
go test -v ./...
```

### Race Detection

```bash
go test -race ./...
```

### Print Debugging

```go
t.Logf("User ID: %s", user.ID)  // Only prints if test fails
fmt.Printf("Debug: %v\n", data) // Always prints
```

---

## Current Test Status

### Phase 8.1 Tests

**Database Migration Tests**: ✅ 5/6 passing (83%)
- Table creation
- Schema validation
- Constraints
- Generated columns
- Cascade deletes

**Repository Tests**: 🚧 In Progress
- User repository
- Invitation repository
- Refresh token repository

**Service Tests**: ⏳ Pending
**Handler Tests**: ⏳ Pending
**Integration Tests**: ⏳ Pending

---

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Library](https://github.com/stretchr/testify) - Assertions
- [sqlx Testing](https://github.com/jmoiron/sqlx)
- [Echo Testing Guide](https://echo.labstack.com/guide/testing/)

---

## Support

If you encounter issues with tests:

1. Check test database is running: `docker ps | grep postgres`
2. Verify connection: `psql postgres://admin:password@localhost:5435/scentora_test`
3. Check logs: Tests output detailed error messages
4. Run with verbose: `go test -v`

---

**Remember**: Tests are documentation. Write tests that clearly show how the code should be used.
