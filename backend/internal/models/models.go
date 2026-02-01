package models

import (
	"time"
)

type User struct {
	ID                    string    `json:"_id" db:"id"`
	Email                 string    `json:"email" db:"email"`
	Username              string    `json:"username" db:"username"`
	PasswordHash          string    `json:"-" db:"password_hash"`
	ValidateRecipeVolumes bool      `json:"validateRecipeVolumes" db:"validate_recipe_volumes"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
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

// Statistics models
type AccordStatistics struct {
	Overview       OverviewStats       `json:"overview"`
	PyramidStats   PyramidStats        `json:"pyramidStats"`
	TagStats       []TagStat           `json:"tagStats"`
	SupplierStats  []SupplierStat      `json:"supplierStats"`
	VolumeStats    VolumeStats         `json:"volumeStats"`
	LowInventory   []LowInventoryItem  `json:"lowInventory"`
}

type OverviewStats struct {
	TotalAccords   int     `json:"totalAccords"`
	TotalVolume    float64 `json:"totalVolume"`
	TotalSuppliers int     `json:"totalSuppliers"`
	TotalTags      int     `json:"totalTags"`
}

type PyramidStats struct {
	TopCount    int     `json:"topCount"`
	MiddleCount int     `json:"middleCount"`
	BaseCount   int     `json:"baseCount"`
	TopVolume   float64 `json:"topVolume"`
	MiddleVolume float64 `json:"middleVolume"`
	BaseVolume  float64 `json:"baseVolume"`
}

type TagStat struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

type SupplierStat struct {
	Supplier string `json:"supplier"`
	Count    int    `json:"count"`
}

type VolumeStats struct {
	AverageVolume float64 `json:"averageVolume"`
	MinVolume     float64 `json:"minVolume"`
	MaxVolume     float64 `json:"maxVolume"`
}

type LowInventoryItem struct {
	AccordID   string  `json:"accordId"`
	Name       string  `json:"name"`
	VolumeML   float64 `json:"volumeMl"`
	Supplier   *string `json:"supplier,omitempty"`
}

// Export/Import additional types
type ExportData struct {
	Version    string            `json:"version"`
	ExportedAt time.Time         `json:"exportedAt"`
	Accords    []*AccordResponse `json:"accords"`
}

type ImportData struct {
	Version string         `json:"version"`
	Accords []ImportAccord `json:"accords"`
}

type ImportResult2 struct {
	TotalRecords    int      `json:"totalRecords"`
	ImportedRecords int      `json:"importedRecords"`
	FailedRecords   int      `json:"failedRecords"`
	Errors          []string `json:"errors,omitempty"`
}

// ========== PHASE 10: RECIPE SYSTEM MODELS ==========

// Recipe represents a perfume formula/recipe
type Recipe struct {
	ID              string     `json:"_id" db:"id"`
	UserID          string     `json:"userId" db:"user_id"`
	Name            string     `json:"name" db:"name"`
	Description     *string    `json:"description,omitempty" db:"description"`
	TargetVolumeMl  float64    `json:"targetVolumeMl" db:"target_volume_ml"`
	Status          string     `json:"status" db:"status"` // draft, in_progress, tested, finalized, archived
	ActiveVersionID *string    `json:"activeVersionId,omitempty" db:"active_version_id"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
}

// RecipeVersion represents a version of a recipe (immutable)
type RecipeVersion struct {
	ID            string    `json:"_id" db:"id"`
	RecipeID      string    `json:"recipeId" db:"recipe_id"`
	VersionNumber int       `json:"versionNumber" db:"version_number"`
	Name          string    `json:"name" db:"name"`
	Notes         *string   `json:"notes,omitempty" db:"notes"`
	IsActive      bool      `json:"isActive" db:"is_active"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

// RecipeIngredient represents an accord used in a recipe version
type RecipeIngredient struct {
	ID            string   `json:"_id" db:"id"`
	VersionID     string   `json:"versionId" db:"version_id"`
	AccordID      string   `json:"accordId" db:"accord_id"`
	QuantityMl    float64  `json:"quantityMl" db:"quantity_ml"`
	QuantityDrops int      `json:"quantityDrops" db:"quantity_drops"`
	Percentage    *float64 `json:"percentage,omitempty" db:"percentage"`
	Notes         *string  `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

// RecipeNote represents a note attached to a recipe or version
type RecipeNote struct {
	ID        string     `json:"_id" db:"id"`
	RecipeID  string     `json:"recipeId" db:"recipe_id"`
	VersionID *string    `json:"versionId,omitempty" db:"version_id"`
	Content   string     `json:"content" db:"content"`
	NoteType  string     `json:"noteType" db:"note_type"` // general, testing, observation
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}

// RecipeTag represents a tag associated with a recipe
type RecipeTag struct {
	ID        string    `json:"_id" db:"id"`
	RecipeID  string    `json:"recipeId" db:"recipe_id"`
	Tag       string    `json:"tag" db:"tag"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// RecipeCollection represents a folder/collection of recipes
type RecipeCollection struct {
	ID          string    `json:"_id" db:"id"`
	UserID      string    `json:"userId" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// RecipeCollectionMember represents the join table between recipes and collections
type RecipeCollectionMember struct {
	ID           string    `json:"_id" db:"id"`
	CollectionID string    `json:"collectionId" db:"collection_id"`
	RecipeID     string    `json:"recipeId" db:"recipe_id"`
	AddedAt      time.Time `json:"addedAt" db:"added_at"`
}

// Request types for Recipes
type CreateRecipeRequest struct {
	Name           string  `json:"name" validate:"required"`
	Description    *string `json:"description"`
	TargetVolumeMl float64 `json:"targetVolumeMl" validate:"required,gt=0"`
	Status         *string `json:"status" validate:"omitempty,oneof=draft in_progress tested finalized archived"`
}

type UpdateRecipeRequest struct {
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	TargetVolumeMl *float64 `json:"targetVolumeMl" validate:"omitempty,gt=0"`
	Status         *string  `json:"status" validate:"omitempty,oneof=draft in_progress tested finalized archived"`
}

type CreateRecipeVersionRequest struct {
	Name  string  `json:"name" validate:"required"`
	Notes *string `json:"notes"`
}

type AddIngredientRequest struct {
	AccordID   string   `json:"accordId" validate:"required"`
	QuantityMl float64  `json:"quantityMl" validate:"required,gt=0"`
	Percentage *float64 `json:"percentage" validate:"omitempty,gte=0,lte=100"`
	Notes      *string  `json:"notes"`
}

type UpdateIngredientRequest struct {
	QuantityMl *float64 `json:"quantityMl" validate:"omitempty,gt=0"`
	Percentage *float64 `json:"percentage" validate:"omitempty,gte=0,lte=100"`
	Notes      *string  `json:"notes"`
}

type CreateRecipeNoteRequest struct {
	VersionID *string `json:"versionId"`
	Content   string  `json:"content" validate:"required"`
	NoteType  *string `json:"noteType" validate:"omitempty,oneof=general testing observation"`
}

type CreateRecipeCollectionRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type UpdateRecipeCollectionRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CreateRecipeIngredientRequest struct {
	AccordID   string  `json:"accordId" validate:"required"`
	QuantityMl float64 `json:"quantityMl" validate:"required,gt=0"`
}

type UpdateRecipeIngredientRequest struct {
	QuantityMl float64 `json:"quantityMl" validate:"required,gt=0"`
}

type UpdateRecipeNoteRequest struct {
	Content string `json:"content" validate:"required"`
}

type RecipeStats struct {
	TotalCount      int `json:"totalCount"`
	DraftCount      int `json:"draftCount"`
	InProgressCount int `json:"inProgressCount"`
	TestedCount     int `json:"testedCount"`
	FinalizedCount  int `json:"finalizedCount"`
	ArchivedCount   int `json:"archivedCount"`
}

type TagCount struct {
	Tag   string `json:"tag" db:"tag"`
	Count int    `json:"count" db:"count"`
}

// Response types for Recipes
type RecipeResponse struct {
	ID              string                   `json:"_id"`
	UserID          string                   `json:"userId"`
	Name            string                   `json:"name"`
	Description     *string                  `json:"description,omitempty"`
	TargetVolumeMl  float64                  `json:"targetVolumeMl"`
	Status          string                   `json:"status"`
	ActiveVersionID *string                  `json:"activeVersionId,omitempty"`
	ActiveVersion   *RecipeVersionResponse   `json:"activeVersion,omitempty"`
	Versions        []RecipeVersionResponse  `json:"versions,omitempty"`
	Tags            []string                 `json:"tags,omitempty"`
	CreatedAt       time.Time                `json:"createdAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
}

type RecipeVersionResponse struct {
	ID            string                      `json:"_id"`
	RecipeID      string                      `json:"recipeId"`
	VersionNumber int                         `json:"versionNumber"`
	Name          string                      `json:"name"`
	Notes         *string                     `json:"notes,omitempty"`
	IsActive      bool                        `json:"isActive"`
	Ingredients   []RecipeIngredientResponse  `json:"ingredients,omitempty"`
	CreatedAt     time.Time                   `json:"createdAt"`
}

type RecipeIngredientResponse struct {
	ID            string   `json:"_id"`
	VersionID     string   `json:"versionId"`
	AccordID      string   `json:"accordId"`
	AccordName    string   `json:"accordName,omitempty"` // Populated via join
	QuantityMl    float64  `json:"quantityMl"`
	QuantityDrops int      `json:"quantityDrops"`
	Percentage    *float64 `json:"percentage,omitempty"`
	Notes         *string  `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}
