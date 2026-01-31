package models

import (
	"time"
)

type User struct {
	ID           string    `json:"_id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// Accord represents a perfume ingredient/accord
type Accord struct {
	ID                 string     `json:"_id" db:"id"`
	UserID             string     `json:"userId" db:"user_id"`
	Name               string     `json:"name" db:"name"`
	PyramidPosition    string     `json:"pyramidPosition" db:"pyramid_position"`
	VolumeMl           float64    `json:"volumeMl" db:"volume_ml"`
	VolumeDrops        int        `json:"volumeDrops" db:"volume_drops"`
	Supplier           *string    `json:"supplier,omitempty" db:"supplier"`
	PurchaseDate       *time.Time `json:"purchaseDate,omitempty" db:"purchase_date"`
	DilutionPercentage *float64   `json:"dilutionPercentage,omitempty" db:"dilution_percentage"`
	Notes              *string    `json:"notes,omitempty" db:"notes"`
	Tags               []string   `json:"tags" db:"-"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time  `json:"updatedAt" db:"updated_at"`
}

// AccordTag represents a tag associated with an accord
type AccordTag struct {
	ID        string    `json:"_id" db:"id"`
	AccordID  string    `json:"accordId" db:"accord_id"`
	Tag       string    `json:"tag" db:"tag"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// PredefinedTag represents a system-defined tag
type PredefinedTag struct {
	ID        string    `json:"_id" db:"id"`
	Category  string    `json:"category" db:"category"`
	Tag       string    `json:"tag" db:"tag"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// AccordResponse includes tags in the response
type AccordResponse struct {
	ID                 string     `json:"_id"`
	UserID             string     `json:"userId"`
	Name               string     `json:"name"`
	PyramidPosition    string     `json:"pyramidPosition"`
	VolumeMl           float64    `json:"volumeMl"`
	VolumeDrops        int        `json:"volumeDrops"`
	Supplier           *string    `json:"supplier,omitempty"`
	PurchaseDate       *time.Time `json:"purchaseDate,omitempty"`
	DilutionPercentage *float64   `json:"dilutionPercentage,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
	Tags               []string   `json:"tags"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

type RefreshToken struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	CreatedAt time.Time `db:"created_at"`
}

type Invitation struct {
	ID        string     `json:"_id" db:"id"`
	Code      string     `json:"code" db:"code"`
	Email     *string    `json:"email,omitempty" db:"email"`
	CreatedBy string     `json:"createdBy" db:"created_by"`
	ExpiresAt time.Time  `json:"expiresAt" db:"expires_at"`
	Used      bool       `json:"used" db:"used"`
	UsedAt    *time.Time `json:"usedAt,omitempty" db:"used_at"`
	UsedBy    *string    `json:"usedBy,omitempty" db:"used_by"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
}

// Request types for Accords
type CreateAccordRequest struct {
	Name               string     `json:"name" validate:"required"`
	PyramidPosition    string     `json:"pyramidPosition" validate:"required,oneof=top middle base"`
	VolumeMl           float64    `json:"volumeMl" validate:"required,gte=0"`
	Supplier           *string    `json:"supplier"`
	PurchaseDate       *time.Time `json:"purchaseDate"`
	DilutionPercentage *float64   `json:"dilutionPercentage" validate:"omitempty,gte=0,lte=100"`
	Notes              *string    `json:"notes"`
	Tags               []string   `json:"tags"`
}

type UpdateAccordRequest struct {
	Name               *string    `json:"name"`
	PyramidPosition    *string    `json:"pyramidPosition" validate:"omitempty,oneof=top middle base"`
	VolumeMl           *float64   `json:"volumeMl" validate:"omitempty,gte=0"`
	Supplier           *string    `json:"supplier"`
	PurchaseDate       *time.Time `json:"purchaseDate"`
	DilutionPercentage *float64   `json:"dilutionPercentage" validate:"omitempty,gte=0,lte=100"`
	Notes              *string    `json:"notes"`
	Tags               *[]string  `json:"tags"`
}

type AddTagRequest struct {
	Tag string `json:"tag" validate:"required,min=1,max=50"`
}

type RegisterRequest struct {
	Email          string `json:"email" validate:"required,email"`
	Username       string `json:"username" validate:"required,min=3"`
	Password       string `json:"password" validate:"required,min=4"`
	InvitationCode string `json:"invitationCode" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         *User  `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type CreateInvitationRequest struct {
	Email         *string `json:"email" validate:"omitempty,email"`
	ExpiresInDays int     `json:"expiresInDays" validate:"omitempty,min=1,max=365"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Export/Import types
type ExportResponse struct {
	Version    string            `json:"version"`
	ExportDate string            `json:"exportDate"`
	Accords    []*AccordResponse `json:"accords"`
}

type ImportCollectionRequest struct {
	Version    string          `json:"version"`
	ExportDate string          `json:"exportDate"`
	Accords    []ImportAccord  `json:"accords" validate:"required"`
}

type ImportAccord struct {
	Name               string     `json:"name"`
	PyramidPosition    string     `json:"pyramidPosition"`
	VolumeMl           float64    `json:"volumeMl"`
	Supplier           *string    `json:"supplier"`
	PurchaseDate       *time.Time `json:"purchaseDate"`
	DilutionPercentage *float64   `json:"dilutionPercentage"`
	Notes              *string    `json:"notes"`
	Tags               []string   `json:"tags"`
}

type ImportResult struct {
	AccordsImported int      `json:"accordsImported"`
	Errors          []string `json:"errors"`
}
