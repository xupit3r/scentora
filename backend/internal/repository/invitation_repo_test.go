package repository

import (
"testing"
"time"

"github.com/yourusername/scentora-backend/internal/models"
"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestInvitationRepository_Create(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

userRepo := NewUserRepository(tdb.DB)
inviteRepo := NewInvitationRepository(tdb.DB)

// Create a user to create invitation
user := &models.User{
Email:        "test@example.com",
Username:     "testuser",
PasswordHash: "hashedpassword",
}
err := userRepo.Create(user)
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

email := "invite@example.com"
expiresAt := time.Now().Add(7 * 24 * time.Hour)

invitation := &models.Invitation{
Code:      "TEST123",
Email:     &email,
CreatedBy: user.ID,
ExpiresAt: expiresAt,
Used:      false,
}

err = inviteRepo.Create(invitation)
if err != nil {
t.Fatalf("Failed to create invitation: %v", err)
}

if invitation.ID == "" {
t.Error("Expected invitation ID to be set")
}

if invitation.Code != "TEST123" {
t.Errorf("Expected code TEST123, got %s", invitation.Code)
}

if invitation.Email == nil || *invitation.Email != email {
t.Errorf("Expected email %s, got %v", email, invitation.Email)
}
}

func TestInvitationRepository_FindByCode(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

userRepo := NewUserRepository(tdb.DB)
inviteRepo := NewInvitationRepository(tdb.DB)

user := &models.User{
Email:        "test@example.com",
Username:     "testuser",
PasswordHash: "hashedpassword",
}
err := userRepo.Create(user)
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

expiresAt := time.Now().Add(7 * 24 * time.Hour)
created := &models.Invitation{
Code:      "TEST123",
CreatedBy: user.ID,
ExpiresAt: expiresAt,
Used:      false,
}
err = inviteRepo.Create(created)
if err != nil {
t.Fatalf("Failed to create invitation: %v", err)
}

found, err := inviteRepo.FindByCode("TEST123")
if err != nil {
t.Fatalf("Failed to find invitation by code: %v", err)
}

if found.ID != created.ID {
t.Errorf("Expected invitation ID %s, got %s", created.ID, found.ID)
}

if found.Code != "TEST123" {
t.Errorf("Expected code TEST123, got %s", found.Code)
}
}

func TestInvitationRepository_FindByCodeNotFound(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

inviteRepo := NewInvitationRepository(tdb.DB)

_, err := inviteRepo.FindByCode("NONEXISTENT")
if err == nil {
t.Error("Expected error for nonexistent code, got nil")
}
}

func TestInvitationRepository_MarkAsUsed(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

userRepo := NewUserRepository(tdb.DB)
inviteRepo := NewInvitationRepository(tdb.DB)

creator := &models.User{
Email:        "creator@example.com",
Username:     "creator",
PasswordHash: "hashedpassword",
}
err := userRepo.Create(creator)
if err != nil {
t.Fatalf("Failed to create creator: %v", err)
}

recipient := &models.User{
Email:        "recipient@example.com",
Username:     "recipient",
PasswordHash: "hashedpassword",
}
err = userRepo.Create(recipient)
if err != nil {
t.Fatalf("Failed to create recipient: %v", err)
}

expiresAt := time.Now().Add(7 * 24 * time.Hour)
invitation := &models.Invitation{
Code:      "TEST123",
CreatedBy: creator.ID,
ExpiresAt: expiresAt,
Used:      false,
}
err = inviteRepo.Create(invitation)
if err != nil {
t.Fatalf("Failed to create invitation: %v", err)
}

// Mark as used
err = inviteRepo.MarkAsUsed(invitation.ID, recipient.ID)
if err != nil {
t.Fatalf("Failed to mark invitation as used: %v", err)
}

// Verify it was marked as used
updated, err := inviteRepo.FindByCode("TEST123")
if err != nil {
t.Fatalf("Failed to get updated invitation: %v", err)
}

if !updated.Used {
t.Error("Expected invitation to be marked as used")
}

if updated.UsedBy == nil || *updated.UsedBy != recipient.ID {
t.Errorf("Expected usedBy to be %s, got %v", recipient.ID, updated.UsedBy)
}

if updated.UsedAt == nil {
t.Error("Expected usedAt to be set")
}
}

func TestInvitationRepository_ListByCreator(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

userRepo := NewUserRepository(tdb.DB)
inviteRepo := NewInvitationRepository(tdb.DB)

user1 := &models.User{
Email:        "user1@example.com",
Username:     "user1",
PasswordHash: "hashedpassword",
}
err := userRepo.Create(user1)
if err != nil {
t.Fatalf("Failed to create user1: %v", err)
}

user2 := &models.User{
Email:        "user2@example.com",
Username:     "user2",
PasswordHash: "hashedpassword",
}
err = userRepo.Create(user2)
if err != nil {
t.Fatalf("Failed to create user2: %v", err)
}

expiresAt := time.Now().Add(7 * 24 * time.Hour)

// Create invitations for user1
invite1 := &models.Invitation{
Code:      "USER1_CODE1",
CreatedBy: user1.ID,
ExpiresAt: expiresAt,
Used:      false,
}
err = inviteRepo.Create(invite1)
if err != nil {
t.Fatalf("Failed to create invitation 1: %v", err)
}

invite2 := &models.Invitation{
Code:      "USER1_CODE2",
CreatedBy: user1.ID,
ExpiresAt: expiresAt,
Used:      false,
}
err = inviteRepo.Create(invite2)
if err != nil {
t.Fatalf("Failed to create invitation 2: %v", err)
}

// Create invitation for user2
invite3 := &models.Invitation{
Code:      "USER2_CODE1",
CreatedBy: user2.ID,
ExpiresAt: expiresAt,
Used:      false,
}
err = inviteRepo.Create(invite3)
if err != nil {
t.Fatalf("Failed to create invitation 3: %v", err)
}

// List invitations for user1
invitations, err := inviteRepo.ListByCreator(user1.ID)
if err != nil {
t.Fatalf("Failed to list invitations: %v", err)
}

if len(invitations) != 2 {
t.Errorf("Expected 2 invitations for user1, got %d", len(invitations))
}
}

func TestInvitationRepository_Revoke(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

userRepo := NewUserRepository(tdb.DB)
inviteRepo := NewInvitationRepository(tdb.DB)

user := &models.User{
Email:        "test@example.com",
Username:     "testuser",
PasswordHash: "hashedpassword",
}
err := userRepo.Create(user)
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

expiresAt := time.Now().Add(7 * 24 * time.Hour)
invitation := &models.Invitation{
Code:      "TEST123",
CreatedBy: user.ID,
ExpiresAt: expiresAt,
Used:      false,
}
err = inviteRepo.Create(invitation)
if err != nil {
t.Fatalf("Failed to create invitation: %v", err)
}

// Revoke the invitation
err = inviteRepo.Revoke("TEST123", user.ID)
if err != nil {
t.Fatalf("Failed to revoke invitation: %v", err)
}

// Verify it was revoked (marked as used)
updated, err := inviteRepo.FindByCode("TEST123")
if err != nil {
t.Fatalf("Failed to get updated invitation: %v", err)
}

if !updated.Used {
t.Error("Expected invitation to be marked as used (revoked)")
}
}
