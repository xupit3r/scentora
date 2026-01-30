package models

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID           string    `json:"_id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type Perfume struct {
	ID            string         `json:"_id" db:"id"`
	UserID        string         `json:"userId" db:"user_id"`
	Name          string         `json:"name" db:"name"`
	Designer      string         `json:"designer" db:"designer"`
	Year          *int           `json:"year,omitempty" db:"year"`
	Concentration *string        `json:"concentration,omitempty" db:"concentration"`
	TopNotes      pq.StringArray `json:"-" db:"top_notes"`
	MiddleNotes   pq.StringArray `json:"-" db:"middle_notes"`
	BaseNotes     pq.StringArray `json:"-" db:"base_notes"`
	Description   *string        `json:"description,omitempty" db:"description"`
	ImageURL      *string        `json:"imageUrl,omitempty" db:"image_url"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
}

type Pyramid struct {
	Top    []string `json:"top"`
	Middle []string `json:"middle"`
	Base   []string `json:"base"`
}

type PerfumeResponse struct {
	ID            string    `json:"_id"`
	UserID        string    `json:"userId"`
	Name          string    `json:"name"`
	Designer      string    `json:"designer"`
	Year          *int      `json:"year,omitempty"`
	Concentration *string   `json:"concentration,omitempty"`
	Pyramid       Pyramid   `json:"pyramid"`
	Description   *string   `json:"description,omitempty"`
	ImageURL      *string   `json:"imageUrl,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (p *Perfume) ToResponse() *PerfumeResponse {
	return &PerfumeResponse{
		ID:            p.ID,
		UserID:        p.UserID,
		Name:          p.Name,
		Designer:      p.Designer,
		Year:          p.Year,
		Concentration: p.Concentration,
		Pyramid: Pyramid{
			Top:    p.TopNotes,
			Middle: p.MiddleNotes,
			Base:   p.BaseNotes,
		},
		Description: p.Description,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type JournalEntry struct {
	ID         string     `json:"_id" db:"id"`
	UserID     string     `json:"userId" db:"user_id"`
	PerfumeID  string     `json:"perfumeId" db:"perfume_id"`
	Date       string     `json:"date" db:"date"`
	Content    string     `json:"content" db:"content"`
	Rating     *int       `json:"rating,omitempty" db:"rating"`
	Occasion   *string    `json:"occasion,omitempty" db:"occasion"`
	Weather    *string    `json:"weather,omitempty" db:"weather"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
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

type CreatePerfumeRequest struct {
	Name          string   `json:"name" validate:"required"`
	Designer      string   `json:"designer" validate:"required"`
	Year          *int     `json:"year"`
	Concentration *string  `json:"concentration"`
	Pyramid       Pyramid  `json:"pyramid" validate:"required"`
	Description   *string  `json:"description"`
	ImageURL      *string  `json:"imageUrl"`
}

type UpdatePerfumeRequest struct {
	Name          *string  `json:"name"`
	Designer      *string  `json:"designer"`
	Year          *int     `json:"year"`
	Concentration *string  `json:"concentration"`
	Pyramid       *Pyramid `json:"pyramid"`
	Description   *string  `json:"description"`
	ImageURL      *string  `json:"imageUrl"`
}

type CreateJournalRequest struct {
	PerfumeID string  `json:"perfumeId" validate:"required"`
	Date      string  `json:"date" validate:"required"`
	Content   string  `json:"content" validate:"required"`
	Rating    *int    `json:"rating" validate:"omitempty,min=1,max=10"`
	Occasion  *string `json:"occasion"`
	Weather   *string `json:"weather"`
}

type UpdateJournalRequest struct {
	Date     *string `json:"date"`
	Content  *string `json:"content"`
	Rating   *int    `json:"rating" validate:"omitempty,min=1,max=10"`
	Occasion *string `json:"occasion"`
	Weather  *string `json:"weather"`
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

type ImportCollectionRequest struct {
	Version        string           `json:"version"`
	ExportDate     string           `json:"exportDate"`
	Perfumes       []ImportPerfume  `json:"perfumes" validate:"required"`
	JournalEntries []ImportJournal  `json:"journalEntries"`
}

type ImportPerfume struct {
	Name          string   `json:"name"`
	Designer      string   `json:"designer"`
	Year          *int     `json:"year"`
	Concentration *string  `json:"concentration"`
	Pyramid       Pyramid  `json:"pyramid"`
	Description   *string  `json:"description"`
	ImageURL      *string  `json:"imageUrl"`
}

type ImportJournal struct {
	PerfumeID string  `json:"perfumeId"`
	Date      string  `json:"date"`
	Content   string  `json:"content"`
	Rating    *int    `json:"rating"`
	Occasion  *string `json:"occasion"`
	Weather   *string `json:"weather"`
}

type ImportResult struct {
	PerfumesImported       int      `json:"perfumesImported"`
	JournalEntriesImported int      `json:"journalEntriesImported"`
	Errors                 []string `json:"errors"`
}

type ExportResponse struct {
	Version        string              `json:"version"`
	ExportDate     string              `json:"exportDate"`
	Perfumes       []*PerfumeResponse  `json:"perfumes"`
	JournalEntries []*JournalEntry     `json:"journalEntries"`
}
