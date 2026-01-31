package repository

import (
"testing"
"time"

"github.com/yourusername/scentora-backend/internal/testutil"
)

func TestInvitationRepository_Create(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

userRepo := NewUserRepository(tdb.DB)
inviteRepo := NewInvitationRepository(tdb.DB)

// Create a user to create invitation
user, err := userRepo.Create("test@example.com", "testuser", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

email := "invite@example.com"
expiresAt := time.Now().Add(7 * 24 * time.Hour)

invitation, err := inviteRepo.Create("TEST123", &email, user.ID, expiresAt)
if err != nil {
t.Fatalf("Failed to create invitation: %v", err)
}

if invitation.Code != "TEST123" {
t.Errorf("Expected code TEST123, got %s", invitation.Code)
}

if invitation.Email == nil || *invitation.Email != email {
t.Errorf("Expected email %s, got %v", email, invitation.Email)
}

if invitation.CreatedBy != user.ID {
t.Errorf("Expected createdBy %s, got %s", user.ID, invitation.CreatedBy)
}
}

func TestInvitationRepository_GetByCode(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

userRepo := NewUserRepository(tdb.DB)
inviteRepo := NewInvitationRepository(tdb.DB)

user, err := userRepo.Create("test@example.com", "testuser", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

expiresAt := time.Now().Add(7 * 24 * time.Hour)
created, err := inviteRepo.Create("TEST123", nil, user.ID, expiresAt)
if err != nil {
t.Fatalf("Failed to create invitation: %v", err)
}

found, err := inviteRepo.GetByCode("TEST123")
if err != nil {
t.Fatalf("Failed to get invitation by code: %v", err)
}

if found.ID != created.ID {
t.Errorf("Expected invitation ID %s, got %s", created.ID, found.ID)
}

if found.Code != "TEST123" {
t.Errorf("Expected code TEST123, got %s", found.Code)
}
}

func TestInvitationRepository_GetByCodeNotFound(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

inviteRepo := NewInvitationRepository(tdb.DB)

_, err := inviteRepo.GetByCode("NONEXISTENT")
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

creator, err := userRepo.Create("creator@example.com", "creator", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create creator: %v", err)
}

recipient, err := userRepo.Create("recipient@example.com", "recipient", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create recipient: %v", err)
}

expiresAt := time.Now().Add(7 * 24 * time.Hour)
invitation, err := inviteRepo.Create("TEST123", nil, creator.ID, expiresAt)
if err != nil {
t.Fatalf("Failed to create invitation: %v", err)
}

// Mark as used
err = inviteRepo.MarkAsUsed(invitation.Code, recipient.ID)
if err != nil {
t.Fatalf("Failed to mark invitation as used: %v", err)
}

// Verify it was marked as used
updated, err := inviteRepo.GetByCode("TEST123")
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

user1, err := userRepo.Create("user1@example.com", "user1", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create user1: %v", err)
}

user2, err := userRepo.Create("user2@example.com", "user2", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create user2: %v", err)
}

expiresAt := time.Now().Add(7 * 24 * time.Hour)

// Create invitations for user1
_, err = inviteRepo.Create("USER1_CODE1", nil, user1.ID, expiresAt)
if err != nil {
t.Fatalf("Failed to create invitation 1: %v", err)
}

_, err = inviteRepo.Create("USER1_CODE2", nil, user1.ID, expiresAt)
if err != nil {
t.Fatalf("Failed to create invitation 2: %v", err)
}

// Create invitation for user2
_, err = inviteRepo.Create("USER2_CODE1", nil, user2.ID, expiresAt)
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

func TestInvitationRepository_Delete(t *testing.T) {
tdb := testutil.SetupTestDB(t)
defer tdb.Teardown(t)
defer tdb.CleanupTables(t)

userRepo := NewUserRepository(tdb.DB)
inviteRepo := NewInvitationRepository(tdb.DB)

user, err := userRepo.Create("test@example.com", "testuser", "hashedpassword")
if err != nil {
t.Fatalf("Failed to create user: %v", err)
}

expiresAt := time.Now().Add(7 * 24 * time.Hour)
_, err = inviteRepo.Create("TEST123", nil, user.ID, expiresAt)
if err != nil {
t.Fatalf("Failed to create invitation: %v", err)
}

// Delete the invitation
err = inviteRepo.Delete("TEST123")
if err != nil {
t.Fatalf("Failed to delete invitation: %v", err)
}

// Verify it was deleted
_, err = inviteRepo.GetByCode("TEST123")
if err == nil {
t.Error("Expected error for deleted invitation, got nil")
}
}
