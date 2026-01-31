package repository

import (
"testing"

"github.com/yourusername/scentora-backend/internal/models"
"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestUserRepository_Create(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

user := &models.User{
Email:        "test@example.com",
Username:     "testuser",
PasswordHash: "hashedpassword",
}

err := repo.Create(user)
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
user1 := &models.User{
Email:        "test@example.com",
Username:     "user1",
PasswordHash: "hashedpassword",
}
err := repo.Create(user1)
if err != nil {
t.Fatalf("Failed to create first user: %v", err)
}

// Try to create second user with same email
user2 := &models.User{
Email:        "test@example.com",
Username:     "user2",
PasswordHash: "hashedpassword",
}
err = repo.Create(user2)
if err == nil {
t.Error("Expected error for duplicate email, got nil")
}
}

func TestUserRepository_FindByEmail(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

// Create a user
created := &models.User{
Email:        "test@example.com",
Username:     "testuser",
PasswordHash: "hashedpassword",
}
err := repo.Create(created)
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

// Find by email
found, err := repo.FindByEmail("test@example.com")
if err != nil {
t.Fatalf("Failed to find user by email: %v", err)
}

if found.ID != created.ID {
t.Errorf("Expected user ID %s, got %s", created.ID, found.ID)
}

if found.Email != "test@example.com" {
t.Errorf("Expected email test@example.com, got %s", found.Email)
}
}

func TestUserRepository_FindByEmailNotFound(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

_, err := repo.FindByEmail("nonexistent@example.com")
if err == nil {
t.Error("Expected error for nonexistent email, got nil")
}
}

func TestUserRepository_FindByID(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

// Create a user
created := &models.User{
Email:        "test@example.com",
Username:     "testuser",
PasswordHash: "hashedpassword",
}
err := repo.Create(created)
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

// Find by ID
found, err := repo.FindByID(created.ID)
if err != nil {
t.Fatalf("Failed to find user by ID: %v", err)
}

if found.ID != created.ID {
t.Errorf("Expected user ID %s, got %s", created.ID, found.ID)
}

if found.Email != "test@example.com" {
t.Errorf("Expected email test@example.com, got %s", found.Email)
}
}

func TestUserRepository_FindByIDNotFound(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

_, err := repo.FindByID("00000000-0000-0000-0000-000000000000")
if err == nil {
t.Error("Expected error for nonexistent ID, got nil")
}
}

func TestUserRepository_EmailExists(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

repo := NewUserRepository(tdb.DB)

// Create a user
user := &models.User{
Email:        "test@example.com",
Username:     "testuser",
PasswordHash: "hashedpassword",
}
err := repo.Create(user)
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

// Check existing email
exists, err := repo.EmailExists("test@example.com")
if err != nil {
t.Fatalf("Failed to check email exists: %v", err)
}
if !exists {
t.Error("Expected email to exist, got false")
}

// Check non-existing email
exists, err = repo.EmailExists("nonexistent@example.com")
if err != nil {
t.Fatalf("Failed to check email exists: %v", err)
}
if exists {
t.Error("Expected email to not exist, got true")
}
}
