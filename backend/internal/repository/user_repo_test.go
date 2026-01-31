package repository

import (
"testing"

"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestUserRepository_Create(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

user, err := repo.Create("test@example.com", "testuser", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

if user.ID == "" {
t.Error("Expected user ID to be set, got empty string")
}

if user.Email != "test@example.com" {
t.Errorf("Expected email test@example.com, got %s", user.Email)
}

if user.Username != "testuser" {
t.Errorf("Expected username testuser, got %s", user.Username)
}
}

func TestUserRepository_CreateDuplicateEmail(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

// Create first user
_, err := repo.Create("test@example.com", "user1", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create first user: %v", err)
}

// Try to create second user with same email
_, err = repo.Create("test@example.com", "user2", "hashedpassword")
if err == nil {
t.Error("Expected error for duplicate email, got nil")
}
}

func TestUserRepository_GetByEmail(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

// Create a user
created, err := repo.Create("test@example.com", "testuser", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

// Get by email
found, err := repo.GetByEmail("test@example.com")
if err != nil {
t.Fatalf("Failed to get user by email: %v", err)
}

if found.ID != created.ID {
t.Errorf("Expected user ID %s, got %s", created.ID, found.ID)
}

if found.Email != "test@example.com" {
t.Errorf("Expected email test@example.com, got %s", found.Email)
}
}

func TestUserRepository_GetByEmailNotFound(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

_, err := repo.GetByEmail("nonexistent@example.com")
if err == nil {
t.Error("Expected error for nonexistent email, got nil")
}
}

func TestUserRepository_GetByID(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

// Create a user
created, err := repo.Create("test@example.com", "testuser", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

// Get by ID
found, err := repo.GetByID(created.ID)
if err != nil {
t.Fatalf("Failed to get user by ID: %v", err)
}

if found.ID != created.ID {
t.Errorf("Expected user ID %s, got %s", created.ID, found.ID)
}

if found.Email != "test@example.com" {
t.Errorf("Expected email test@example.com, got %s", found.Email)
}
}

func TestUserRepository_GetByIDNotFound(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

_, err := repo.GetByID("00000000-0000-0000-0000-000000000000")
if err == nil {
t.Error("Expected error for nonexistent ID, got nil")
}
}
