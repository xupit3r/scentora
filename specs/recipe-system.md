# Recipe/Formula System Specification

**Created**: January 31, 2026  
**Updated**: February 1, 2026  
**Phase**: 10  
**Status**: Design Decisions Confirmed - Ready for Implementation  
**Version**: 1.1

---

## Overview

This document specifies the Recipe/Formula system for Scentora - a comprehensive feature that enables users to create perfume recipes by combining their accords with version control, notes, and organization capabilities.

### Purpose

Enable users to:
- Create formulas that specify which accords to combine and in what proportions
- Version their recipes as they iterate and refine them
- Organize recipes with tags and collections
- Document their blending process with notes
- Export/share recipes with others

### Key Design Principles

1. **Non-inventory tracking** - Recipes don't deduct from accord inventory (accords remain independent)
2. **Version control** - Multiple versions of recipes, track iterations (new versions auto-activate)
3. **Configurable validation** - Optional volume validation (disabled by default, opt-in per user)
4. **Rich organization** - Tags, collections, search, filtering
5. **Comprehensive notes** - Journal entries per recipe and per version
6. **Export/share** - JSON export for sharing formulas
7. **Accord protection** - Prevent deletion of accords used in recipes (409 error with recipe list)

---

## ✅ Confirmed Design Decisions (February 1, 2026)

These decisions have been confirmed with the project stakeholders and should be implemented as specified:

### 1. Volume Validation Default
**Decision**: ❌ **Disabled by default** (users must opt-in)

**Rationale**: 
- Allows users to create theoretical/planning recipes without worrying about current inventory
- Users who want strict inventory tracking can enable it in their preferences
- More flexible for experimentation and formula development

**Implementation**:
- Add `validate_recipe_volumes` boolean field to `users` table (default: `false`)
- Provide API endpoint to toggle this preference: `PUT /api/preferences`
- When enabled, validate accord volumes during recipe ingredient creation/update
- When disabled, allow any quantities without validation

### 2. Recipe Version Activation
**Decision**: ✅ **New versions automatically become the active version**

**Rationale**:
- Latest version is typically the one being actively worked on
- Reduces friction in the development workflow
- Users can manually set a different version as active if needed

**Implementation**:
- When creating a new version, set `is_active = true` for new version
- Set `is_active = false` for all other versions of the same recipe
- Provide API endpoint to manually change active version: `PUT /api/recipes/:id/versions/:versionNumber/activate`

### 3. Accord Deletion Protection
**Decision**: 🛡️ **Prevent deletion with error message**

**Rationale**:
- Safer approach - prevents accidental data loss
- Forces user to be explicit about cleanup (remove from recipes first)
- Clear error message guides user on next steps

**Implementation**:
- Before deleting an accord, check `recipe_ingredients` table for references
- If references exist, return `409 Conflict` with:
  ```json
  {
    "error": {
      "message": "Cannot delete accord - it is used in recipes",
      "details": "This accord is used in 3 recipes. Remove it from recipes or delete the recipes first.",
      "recipes": [
        {"id": "uuid", "name": "Summer Citrus Blend"},
        {"id": "uuid", "name": "Fresh Morning Cologne"},
        {"id": "uuid", "name": "Experimental Mix v2"}
      ]
    }
  }
  ```
- Add foreign key constraint on `recipe_ingredients.accord_id` with `ON DELETE RESTRICT`

### 4. Recipe Status Values
**Decision**: 📊 **Five-state workflow** (draft, in_progress, tested, finalized, archived)

**Rationale**:
- More granular than 4 states, providing better workflow tracking
- `in_progress` distinguishes active development from initial drafts
- Clear progression: draft → in_progress → tested → finalized → archived

**Implementation**:
- Database: `status VARCHAR(20) CHECK (status IN ('draft', 'in_progress', 'tested', 'finalized', 'archived'))`
- Default status on creation: `draft`
- API filtering: Support `?status=tested` query parameter
- Frontend: Display status badges with appropriate colors
  - draft: gray
  - in_progress: blue
  - tested: yellow
  - finalized: green
  - archived: gray (dimmed)

### 5. Implementation Approach
**Decision**: 🔄 **Backend-first** (Phases 10.1-10.5, then 10.6-10.10, then 10.11)

**Rationale**:
- Cleaner separation of concerns
- Easier to track progress (complete backend before frontend)
- API can be tested independently before UI development
- Follows established pattern from Phase 8 (Accord System)

**Phases**:
1. **Backend** (Phases 10.1-10.5): Models, repos, services, handlers, routes
2. **Frontend** (Phases 10.6-10.10): Types, stores, components, views, features
3. **Testing** (Phase 10.11): Comprehensive test suite for both layers

---

## Feature Requirements

### Must-Have Features (Phase 10.1-10.5)

✅ **Core Recipe Model**
- Recipe name, description, total target volume
- Creator/owner (user_id for isolation)
- Status (draft, in_progress, tested, finalized, archived) ⭐ **UPDATED**
- Created/updated timestamps

✅ **Recipe Versions**
- Each recipe can have multiple versions (v1, v2, v3...)
- Each version is immutable once created
- Version-specific ingredients, notes, and metadata
- Active version designation (new versions auto-activate) ⭐ **UPDATED**

✅ **Recipe Ingredients**
- Link to accords (many-to-many)
- Quantity (ml and drops) per accord
- Percentage of total formula
- Optional notes per ingredient
- Foreign key with RESTRICT to prevent accord deletion ⭐ **UPDATED**

✅ **Volume Validation (Configurable)**
- Optional setting to validate accord availability (disabled by default) ⭐ **UPDATED**
- Warning/error if accord volume insufficient (when enabled)
- Can be disabled for planning/theoretical recipes (default behavior) ⭐ **UPDATED**
- Per-user preference stored in user_preferences table

✅ **Recipe Notes/Journaling**
- General recipe notes (description, inspiration)
- Version-specific notes (changes made, results)
- Testing notes (batch notes, observations)

✅ **Tags & Organization**
- Tag recipes (similar to accord tags)
- Search by name, description, notes
- Filter by tags, status, date

✅ **Accord Protection**
- Prevent deletion of accords used in recipes ⭐ **UPDATED**
- Return 409 Conflict with list of recipes using the accord ⭐ **UPDATED**
- User must remove from recipes first ⭐ **UPDATED**

### Should-Have Features (Phase 2)

📋 **Recipe Collections**
- Group recipes into collections/folders
- Examples: "Summer Scents", "Experiments", "Client Projects"
- Collection-level tags and metadata

📋 **Export/Share**
- Export recipe (single) as JSON
- Export collection as JSON
- Import recipes from JSON
- Share link generation (optional future: public recipes)

📋 **Advanced Features**
- Calculate total volume from ingredients
- Validate percentages sum to 100%
- Recipe duplication (copy and modify)
- Batch scaling (scale recipe to different volumes)

### Nice-to-Have Features (Future)

💭 **Statistics & Analytics**
- Most-used accords
- Recipe complexity metrics
- Tag distribution across recipes

💭 **Collaboration** (Future)
- Share recipes between users
- Public recipe database
- Recipe ratings/favorites

💭 **Cost Tracking** (Future - if accords get pricing)
- Calculate recipe cost based on accord prices
- Cost per ml for finished perfume

---

## Workplan

### Phase 10.1: Backend - Data Models & Database Schema

- [ ] **Update User Model**
  - [ ] Add `validate_recipe_volumes` boolean field to User struct
  - [ ] Add migration to alter users table (default: false)
  - [ ] Update UserResponse to include setting

- [ ] **Create Recipe Models**
  - [ ] Define `Recipe` struct in Go
  - [ ] Define `RecipeVersion` struct
  - [ ] Define `RecipeIngredient` struct
  - [ ] Define `RecipeNote` struct
  - [ ] Define `RecipeCollection` struct
  - [ ] Define `RecipeCollectionMember` struct (join table)
  - [ ] Add request/response DTOs

- [ ] **Database Migrations**
  - [ ] Alter `users` table - add `validate_recipe_volumes` column
  - [ ] Create `recipes` table
  - [ ] Create `recipe_versions` table
  - [ ] Create `recipe_ingredients` table (with RESTRICT on accord_id)
  - [ ] Create `recipe_notes` table
  - [ ] Create `recipe_tags` table (similar to accord_tags)
  - [ ] Create `recipe_collections` table
  - [ ] Create `recipe_collection_members` table
  - [ ] Add indexes for performance
  - [ ] Add foreign key constraints with CASCADE/RESTRICT

- [ ] **Update Documentation**
  - [ ] Add models to `specs/data-models.md`
  - [ ] Create `specs/recipe-system.md` with full specification
  - [ ] Document validation rules
  - [ ] Document version control strategy
  - [ ] Document accord deletion protection behavior

### Phase 10.2: Backend - Repository Layer

- [ ] **Recipe Repository**
  - [ ] `CreateRecipe(recipe *Recipe) error`
  - [ ] `GetRecipe(recipeID, userID string) (*Recipe, error)`
  - [ ] `ListRecipes(userID string, filters RecipeFilters) ([]*Recipe, error)`
  - [ ] `UpdateRecipe(recipe *Recipe) error`
  - [ ] `DeleteRecipe(recipeID, userID string) error`
  - [ ] `SearchRecipes(userID, query string) ([]*Recipe, error)`

- [ ] **Recipe Version Repository**
  - [ ] `CreateVersion(version *RecipeVersion) error`
  - [ ] `GetVersion(versionID, userID string) (*RecipeVersion, error)`
  - [ ] `ListVersions(recipeID, userID string) ([]*RecipeVersion, error)`
  - [ ] `GetActiveVersion(recipeID, userID string) (*RecipeVersion, error)`
  - [ ] `SetActiveVersion(recipeID, versionID, userID string) error`
  - [ ] Version immutability checks

- [ ] **Accord Repository Updates**
  - [ ] Add `GetRecipesUsingAccord(accordID string) ([]*Recipe, error)`
  - [ ] Update `DeleteAccord` to check for recipe usage (RESTRICT behavior)
  - [ ] Return helpful error with recipe list if deletion blocked

- [ ] **Recipe Ingredient Repository**
  - [ ] `AddIngredient(ingredient *RecipeIngredient) error`
  - [ ] `GetIngredients(versionID string) ([]*RecipeIngredient, error)`
  - [ ] `UpdateIngredient(ingredient *RecipeIngredient) error`
  - [ ] `RemoveIngredient(ingredientID string) error`
  - [ ] Validate accord existence

- [ ] **Recipe Note Repository**
  - [ ] `CreateNote(note *RecipeNote) error`
  - [ ] `GetNotes(recipeID, userID string) ([]*RecipeNote, error)`
  - [ ] `UpdateNote(note *RecipeNote) error`
  - [ ] `DeleteNote(noteID, userID string) error`

- [ ] **Recipe Tag Repository**
  - [ ] `AddTag(recipeID, userID, tag string) error`
  - [ ] `RemoveTag(recipeID, userID, tag string) error`
  - [ ] `GetTags(recipeID string) ([]string, error)`
  - [ ] Use predefined tags system

- [ ] **Recipe Collection Repository**
  - [ ] `CreateCollection(collection *RecipeCollection) error`
  - [ ] `GetCollection(collectionID, userID string) (*RecipeCollection, error)`
  - [ ] `ListCollections(userID string) ([]*RecipeCollection, error)`
  - [ ] `UpdateCollection(collection *RecipeCollection) error`
  - [ ] `DeleteCollection(collectionID, userID string) error`
  - [ ] `AddRecipeToCollection(collectionID, recipeID, userID string) error`
  - [ ] `RemoveRecipeFromCollection(collectionID, recipeID, userID string) error`
  - [ ] `GetCollectionRecipes(collectionID, userID string) ([]*Recipe, error)`

- [ ] **Repository Tests** (Target: 90%+ coverage)
  - [ ] Recipe CRUD tests
  - [ ] Version management tests
  - [ ] Ingredient management tests
  - [ ] Note CRUD tests
  - [ ] Tag operations tests
  - [ ] Collection CRUD tests
  - [ ] User isolation tests
  - [ ] Cascade deletion tests

### Phase 10.3: Backend - Service Layer

- [ ] **Recipe Service**
  - [ ] `CreateRecipe(userID string, req *CreateRecipeRequest) (*Recipe, error)`
  - [ ] `GetRecipe(recipeID, userID string) (*RecipeResponse, error)`
  - [ ] `ListRecipes(userID string, filters) ([]*RecipeResponse, error)`
  - [ ] `UpdateRecipe(recipeID, userID string, req *UpdateRecipeRequest) (*Recipe, error)`
  - [ ] `DeleteRecipe(recipeID, userID string) error`
  - [ ] `SearchRecipes(userID, query string) ([]*RecipeResponse, error)`
  - [ ] Business logic for recipe validation

- [ ] **Recipe Version Service**
  - [ ] `CreateVersion(recipeID, userID string, req *CreateVersionRequest) (*RecipeVersion, error)`
  - [ ] `GetVersion(versionID, userID string) (*RecipeVersionResponse, error)`
  - [ ] `ListVersions(recipeID, userID string) ([]*RecipeVersionResponse, error)`
  - [ ] `SetActiveVersion(recipeID, versionID, userID string) error`
  - [ ] `DuplicateVersion(versionID, userID string) (*RecipeVersion, error)`
  - [ ] Validate ingredients sum to target volume
  - [ ] Calculate percentages
  - [ ] Check user's `validate_recipe_volumes` preference
  - [ ] If enabled, validate accord availability for all ingredients
  - [ ] Return helpful error if accord volume insufficient

- [ ] **Recipe Note Service**
  - [ ] `CreateNote(recipeID, userID string, req *CreateNoteRequest) (*RecipeNote, error)`
  - [ ] `GetNotes(recipeID, userID string) ([]*RecipeNote, error)`
  - [ ] `UpdateNote(noteID, userID string, req *UpdateNoteRequest) (*RecipeNote, error)`
  - [ ] `DeleteNote(noteID, userID string) error`

- [ ] **Recipe Collection Service**
  - [ ] `CreateCollection(userID string, req *CreateCollectionRequest) (*RecipeCollection, error)`
  - [ ] `GetCollection(collectionID, userID string) (*RecipeCollectionResponse, error)`
  - [ ] `ListCollections(userID string) ([]*RecipeCollection, error)`
  - [ ] `UpdateCollection(collectionID, userID string, req *UpdateCollectionRequest) (*RecipeCollection, error)`
  - [ ] `DeleteCollection(collectionID, userID string) error`
  - [ ] `AddRecipeToCollection(collectionID, recipeID, userID string) error`
  - [ ] `RemoveRecipeFromCollection(collectionID, recipeID, userID string) error`
  - [ ] `GetCollectionRecipes(collectionID, userID string) ([]*RecipeResponse, error)`

- [ ] **Export/Import Service**
  - [ ] `ExportRecipe(recipeID, userID string) (*RecipeExport, error)`
  - [ ] `ExportCollection(collectionID, userID string) (*CollectionExport, error)`
  - [ ] `ImportRecipe(userID string, data *RecipeExport) (*Recipe, error)`
  - [ ] `ImportCollection(userID string, data *CollectionExport) (*RecipeCollection, error)`
  - [ ] JSON serialization/deserialization
  - [ ] Version handling on import
  - [ ] **Name conflict detection**: Check if recipe name exists
  - [ ] Return error with conflict details if name exists
  - [ ] Suggest alternative name (append " (2)", " (3)", etc.)
  - [ ] Require user to provide new name to proceed

- [ ] **User Settings Service**
  - [ ] `UpdateVolumeValidationSetting(userID string, enabled bool) error`
  - [ ] `GetUserSettings(userID string) (*UserSettings, error)`

- [ ] **Service Tests** (Target: 85%+ coverage)
  - [ ] Recipe business logic tests
  - [ ] Version control tests
  - [ ] Validation tests (volumes, percentages)
  - [ ] Note management tests
  - [ ] Collection management tests
  - [ ] Export/import tests
  - [ ] Error handling tests

### Phase 10.4: Backend - Handler Layer

- [ ] **Recipe Handlers**
  - [ ] `POST /api/recipes` - Create recipe
  - [ ] `GET /api/recipes` - List recipes (with filtering)
  - [ ] `GET /api/recipes/:id` - Get single recipe
  - [ ] `PUT /api/recipes/:id` - Update recipe
  - [ ] `DELETE /api/recipes/:id` - Delete recipe
  - [ ] `GET /api/recipes/search?q=...` - Search recipes

- [ ] **Recipe Version Handlers**
  - [ ] `POST /api/recipes/:id/versions` - Create new version
  - [ ] `GET /api/recipes/:id/versions` - List versions
  - [ ] `GET /api/recipes/:id/versions/:versionId` - Get specific version
  - [ ] `POST /api/recipes/:id/versions/:versionId/activate` - Set as active
  - [ ] `POST /api/recipes/:id/versions/:versionId/duplicate` - Duplicate version

- [ ] **Recipe Ingredient Handlers**
  - [ ] `POST /api/recipes/:id/versions/:versionId/ingredients` - Add ingredient
  - [ ] `PUT /api/recipes/:id/versions/:versionId/ingredients/:ingredientId` - Update ingredient
  - [ ] `DELETE /api/recipes/:id/versions/:versionId/ingredients/:ingredientId` - Remove ingredient

- [ ] **Recipe Note Handlers**
  - [ ] `POST /api/recipes/:id/notes` - Create note
  - [ ] `GET /api/recipes/:id/notes` - List notes
  - [ ] `PUT /api/recipes/:id/notes/:noteId` - Update note
  - [ ] `DELETE /api/recipes/:id/notes/:noteId` - Delete note

- [ ] **Recipe Tag Handlers**
  - [ ] `POST /api/recipes/:id/tags` - Add tag
  - [ ] `DELETE /api/recipes/:id/tags/:tag` - Remove tag

- [ ] **Recipe Collection Handlers**
  - [ ] `POST /api/collections` - Create collection
  - [ ] `GET /api/collections` - List collections
  - [ ] `GET /api/collections/:id` - Get collection
  - [ ] `PUT /api/collections/:id` - Update collection
  - [ ] `DELETE /api/collections/:id` - Delete collection
  - [ ] `POST /api/collections/:id/recipes` - Add recipe to collection
  - [ ] `DELETE /api/collections/:id/recipes/:recipeId` - Remove recipe from collection

- [ ] **Export/Import Handlers**
  - [ ] `GET /api/recipes/:id/export` - Export recipe as JSON
  - [ ] `POST /api/recipes/import` - Import recipe from JSON (with name conflict handling)
  - [ ] `GET /api/collections/:id/export` - Export collection
  - [ ] `POST /api/collections/import` - Import collection

- [ ] **User Settings Handlers**
  - [ ] `GET /api/settings` - Get user settings
  - [ ] `PUT /api/settings/volume-validation` - Toggle volume validation preference

- [ ] **Handler Tests** (Target: 80%+ coverage)
  - [ ] HTTP integration tests for all endpoints
  - [ ] Authentication tests
  - [ ] Validation tests
  - [ ] Error response tests
  - [ ] User isolation tests

### Phase 10.5: Backend - Routes & Middleware

- [ ] **Route Registration**
  - [ ] Register recipe routes in `routes/routes.go`
  - [ ] Apply authentication middleware to all recipe endpoints
  - [ ] Apply rate limiting if needed
  - [ ] Add CORS configuration for recipe endpoints

- [ ] **API Documentation**
  - [ ] Create `specs/recipe-api.md` with all endpoints
  - [ ] Document request/response formats
  - [ ] Document error codes
  - [ ] Add examples for each endpoint

### Phase 10.6: Frontend - Data Models & Types

- [ ] **TypeScript Interfaces**
  - [ ] Define `Recipe` interface
  - [ ] Define `RecipeVersion` interface
  - [ ] Define `RecipeIngredient` interface
  - [ ] Define `RecipeNote` interface
  - [ ] Define `RecipeCollection` interface
  - [ ] Define request/response types
  - [ ] Define filter types

- [ ] **API Service**
  - [ ] Create `frontend/src/services/recipe.ts`
  - [ ] Implement all recipe API calls
  - [ ] Error handling
  - [ ] Type safety

### Phase 10.7: Frontend - State Management

- [ ] **Recipe Store (Pinia)**
  - [ ] Create `frontend/src/stores/recipe.ts`
  - [ ] State: recipes list, active recipe, loading states
  - [ ] Actions: CRUD operations, search, filtering
  - [ ] Getters: filtered recipes, statistics
  - [ ] Cache management

- [ ] **Recipe Version Store**
  - [ ] State: versions list, active version
  - [ ] Actions: create version, switch version, duplicate
  - [ ] Version comparison logic

- [ ] **Recipe Collection Store**
  - [ ] State: collections list, active collection
  - [ ] Actions: CRUD operations, add/remove recipes
  - [ ] Collection filtering

### Phase 10.8: Frontend - UI Components

- [ ] **Recipe Components**
  - [ ] `RecipeCard.vue` - List view card
  - [ ] `RecipeForm.vue` - Create/edit recipe
  - [ ] `RecipeDetail.vue` - Full recipe view
  - [ ] `RecipeList.vue` - Filtered list of recipes
  - [ ] `RecipeFilters.vue` - Search and filter UI

- [ ] **Version Components**
  - [ ] `VersionSelector.vue` - Dropdown to switch versions
  - [ ] `VersionTimeline.vue` - Visual version history
  - [ ] `VersionForm.vue` - Create new version
  - [ ] `VersionComparison.vue` - Compare two versions

- [ ] **Ingredient Components**
  - [ ] `IngredientPicker.vue` - Select accords
  - [ ] `IngredientList.vue` - List of ingredients with quantities
  - [ ] `IngredientForm.vue` - Add/edit ingredient
  - [ ] `IngredientCalculator.vue` - Calculate ml/drops/percentages

- [ ] **Note Components**
  - [ ] `RecipeNoteEditor.vue` - Rich text editor for notes
  - [ ] `RecipeNoteList.vue` - List of notes
  - [ ] `RecipeNoteCard.vue` - Single note display

- [ ] **Collection Components**
  - [ ] `CollectionCard.vue` - Collection display
  - [ ] `CollectionForm.vue` - Create/edit collection
  - [ ] `CollectionGrid.vue` - Grid of collections
  - [ ] `CollectionRecipeList.vue` - Recipes in a collection

### Phase 10.9: Frontend - Views & Pages

- [ ] **Recipe Views**
  - [ ] `RecipesView.vue` - Main recipes page with list
  - [ ] `RecipeDetailView.vue` - Single recipe detail page
  - [ ] `RecipeCreateView.vue` - Create new recipe
  - [ ] `RecipeEditView.vue` - Edit existing recipe

- [ ] **Collection Views**
  - [ ] `CollectionsView.vue` - List of collections
  - [ ] `CollectionDetailView.vue` - Single collection with recipes

- [ ] **Router Configuration**
  - [ ] Add recipe routes to `router/index.ts`
  - [ ] Route guards for authentication
  - [ ] Breadcrumb configuration

### Phase 10.10: Frontend - Features & Polish

- [ ] **Search & Filtering**
  - [ ] Full-text search across recipes
  - [ ] Filter by tags
  - [ ] Filter by status (draft, tested, etc.)
  - [ ] Filter by date range
  - [ ] Sort options (name, date, volume)

- [ ] **Export/Import UI**
  - [ ] Export button with download
  - [ ] Import file upload
  - [ ] Preview before import
  - [ ] **Conflict resolution modal**: If name exists, prompt for new name
  - [ ] Show suggested name with auto-increment
  - [ ] Validate new name before submitting

- [ ] **Calculations & Validations**
  - [ ] Real-time percentage calculations
  - [ ] Check user's volume validation preference from settings
  - [ ] If enabled, show accord availability warnings
  - [ ] Warn if percentages don't sum to 100%
  - [ ] Batch scaling calculator

- [ ] **User Settings UI**
  - [ ] Add Settings view/page
  - [ ] Volume validation toggle switch
  - [ ] Explain what validation does
  - [ ] Save preference to backend

- [ ] **Empty States & Loading**
  - [ ] Skeleton loaders for recipes
  - [ ] Empty state for no recipes
  - [ ] Empty state for no collections
  - [ ] Loading indicators

- [ ] **Responsive Design**
  - [ ] Mobile-friendly recipe cards
  - [ ] Tablet layout optimizations
  - [ ] Touch-friendly ingredient picker

### Phase 10.11: Testing

- [ ] **Backend Testing**
  - [ ] Run full backend test suite
  - [ ] Verify 80%+ coverage
  - [ ] Fix any failing tests
  - [ ] Add missing test cases

- [ ] **Frontend Testing** (if Vitest is set up)
  - [ ] Component tests for recipe components
  - [ ] Store tests for recipe/collection stores
  - [ ] Integration tests for API service

- [ ] **E2E Testing** (if Playwright/Cypress set up)
  - [ ] User journey: Create recipe
  - [ ] User journey: Add ingredients
  - [ ] User journey: Create version
  - [ ] User journey: Create collection
  - [ ] User journey: Export/import

### Phase 10.12: Documentation & Deployment

- [ ] **Documentation**
  - [ ] Update README with recipe features
  - [ ] Create `docs/RECIPE_GUIDE.md` user guide
  - [ ] Update API documentation
  - [ ] Create phase completion document
  - [ ] Update data model documentation

- [ ] **Final Testing**
  - [ ] Full integration test
  - [ ] Performance testing with large recipe sets
  - [ ] User acceptance testing

- [ ] **Deployment**
  - [ ] Database migrations in production
  - [ ] Deploy backend changes
  - [ ] Deploy frontend changes
  - [ ] Smoke tests in production

---

## Data Model Design

### Recipe
```go
type Recipe struct {
    ID          string    `json:"_id" db:"id"`
    UserID      string    `json:"userId" db:"user_id"`
    Name        string    `json:"name" db:"name"`
    Description *string   `json:"description,omitempty" db:"description"`
    TargetVolumeMl float64 `json:"targetVolumeMl" db:"target_volume_ml"`
    Status      string    `json:"status" db:"status"` // draft, tested, finalized, archived
    ActiveVersionID *string `json:"activeVersionId,omitempty" db:"active_version_id"`
    Tags        []string  `json:"tags" db:"-"`
    CreatedAt   time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
```

### RecipeVersion
```go
type RecipeVersion struct {
    ID          string    `json:"_id" db:"id"`
    RecipeID    string    `json:"recipeId" db:"recipe_id"`
    VersionNumber int     `json:"versionNumber" db:"version_number"`
    Name        string    `json:"name" db:"name"` // e.g., "v1", "Summer 2026", "Reduced Jasmine"
    Notes       *string   `json:"notes,omitempty" db:"notes"`
    IsActive    bool      `json:"isActive" db:"is_active"`
    CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}
```

### RecipeIngredient
```go
type RecipeIngredient struct {
    ID          string   `json:"_id" db:"id"`
    VersionID   string   `json:"versionId" db:"version_id"`
    AccordID    string   `json:"accordId" db:"accord_id"`
    QuantityMl  float64  `json:"quantityMl" db:"quantity_ml"`
    QuantityDrops int    `json:"quantityDrops" db:"quantity_drops"`
    Percentage  float64  `json:"percentage" db:"percentage"` // % of total
    Notes       *string  `json:"notes,omitempty" db:"notes"`
    CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}
```

### RecipeNote
```go
type RecipeNote struct {
    ID        string    `json:"_id" db:"id"`
    RecipeID  string    `json:"recipeId" db:"recipe_id"`
    VersionID *string   `json:"versionId,omitempty" db:"version_id"` // null = recipe-level note
    Content   string    `json:"content" db:"content"`
    NoteType  string    `json:"noteType" db:"note_type"` // general, testing, observation
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
```

### RecipeCollection
```go
type RecipeCollection struct {
    ID          string    `json:"_id" db:"id"`
    UserID      string    `json:"userId" db:"user_id"`
    Name        string    `json:"name" db:"name"`
    Description *string   `json:"description,omitempty" db:"description"`
    Tags        []string  `json:"tags" db:"-"`
    CreatedAt   time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
```

---

## Database Schema

```sql
-- Recipes table
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    target_volume_ml DECIMAL(10,2) NOT NULL CHECK (target_volume_ml > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' 
        CHECK (status IN ('draft', 'tested', 'finalized', 'archived')),
    active_version_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name)
);

CREATE INDEX idx_recipes_user_id ON recipes(user_id);
CREATE INDEX idx_recipes_status ON recipes(status);
CREATE INDEX idx_recipes_created_at ON recipes(created_at DESC);

-- Recipe versions table
CREATE TABLE recipe_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    notes TEXT,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(recipe_id, version_number)
);

CREATE INDEX idx_recipe_versions_recipe_id ON recipe_versions(recipe_id);
CREATE INDEX idx_recipe_versions_is_active ON recipe_versions(is_active);

-- Add foreign key constraint after both tables exist
ALTER TABLE recipes 
    ADD CONSTRAINT fk_recipes_active_version 
    FOREIGN KEY (active_version_id) 
    REFERENCES recipe_versions(id) ON DELETE SET NULL;

-- Recipe ingredients table
-- CRITICAL: ON DELETE RESTRICT prevents accord deletion if used in recipe
CREATE TABLE recipe_ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version_id UUID NOT NULL REFERENCES recipe_versions(id) ON DELETE CASCADE,
    accord_id UUID NOT NULL REFERENCES accords(id) ON DELETE RESTRICT,
    quantity_ml DECIMAL(10,2) NOT NULL CHECK (quantity_ml > 0),
    quantity_drops INTEGER GENERATED ALWAYS AS (ROUND(quantity_ml * 20)) STORED,
    percentage DECIMAL(5,2) CHECK (percentage >= 0 AND percentage <= 100),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(version_id, accord_id)
);

CREATE INDEX idx_recipe_ingredients_version_id ON recipe_ingredients(version_id);
CREATE INDEX idx_recipe_ingredients_accord_id ON recipe_ingredients(accord_id);

-- Helper query for accord deletion protection:
-- SELECT r.id, r.name FROM recipes r 
-- JOIN recipe_versions rv ON rv.recipe_id = r.id 
-- JOIN recipe_ingredients ri ON ri.version_id = rv.id 
-- WHERE ri.accord_id = :accord_id;

-- Recipe notes table
CREATE TABLE recipe_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    version_id UUID REFERENCES recipe_versions(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    note_type VARCHAR(20) DEFAULT 'general' 
        CHECK (note_type IN ('general', 'testing', 'observation')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recipe_notes_recipe_id ON recipe_notes(recipe_id);
CREATE INDEX idx_recipe_notes_version_id ON recipe_notes(version_id);

-- Recipe tags table (similar to accord_tags)
CREATE TABLE recipe_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(recipe_id, tag)
);

CREATE INDEX idx_recipe_tags_recipe_id ON recipe_tags(recipe_id);
CREATE INDEX idx_recipe_tags_tag ON recipe_tags(tag);

-- Recipe collections table
CREATE TABLE recipe_collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name)
);

CREATE INDEX idx_recipe_collections_user_id ON recipe_collections(user_id);

-- Recipe collection members (many-to-many)
CREATE TABLE recipe_collection_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collection_id UUID NOT NULL REFERENCES recipe_collections(id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(collection_id, recipe_id)
);

CREATE INDEX idx_recipe_collection_members_collection_id 
    ON recipe_collection_members(collection_id);
CREATE INDEX idx_recipe_collection_members_recipe_id 
    ON recipe_collection_members(recipe_id);
```

---

## API Endpoints Summary

### Recipes
- `POST /api/recipes` - Create recipe
- `GET /api/recipes` - List recipes (filterable)
- `GET /api/recipes/:id` - Get recipe
- `PUT /api/recipes/:id` - Update recipe
- `DELETE /api/recipes/:id` - Delete recipe
- `GET /api/recipes/search?q=...` - Search

### Versions
- `POST /api/recipes/:id/versions` - Create version
- `GET /api/recipes/:id/versions` - List versions
- `GET /api/recipes/:id/versions/:versionId` - Get version
- `POST /api/recipes/:id/versions/:versionId/activate` - Activate
- `POST /api/recipes/:id/versions/:versionId/duplicate` - Duplicate

### Ingredients
- `POST /api/recipes/:id/versions/:versionId/ingredients` - Add
- `PUT /api/recipes/:id/versions/:versionId/ingredients/:ingredientId` - Update
- `DELETE /api/recipes/:id/versions/:versionId/ingredients/:ingredientId` - Remove

### Notes
- `POST /api/recipes/:id/notes` - Create note
- `GET /api/recipes/:id/notes` - List notes
- `PUT /api/recipes/:id/notes/:noteId` - Update note
- `DELETE /api/recipes/:id/notes/:noteId` - Delete note

### Tags
- `POST /api/recipes/:id/tags` - Add tag
- `DELETE /api/recipes/:id/tags/:tag` - Remove tag

### Collections
- `POST /api/collections` - Create
- `GET /api/collections` - List
- `GET /api/collections/:id` - Get
- `PUT /api/collections/:id` - Update
- `DELETE /api/collections/:id` - Delete
- `POST /api/collections/:id/recipes` - Add recipe
- `DELETE /api/collections/:id/recipes/:recipeId` - Remove recipe

### Export/Import
- `GET /api/recipes/:id/export` - Export recipe
- `POST /api/recipes/import` - Import recipe (returns error if name conflict)
- `GET /api/collections/:id/export` - Export collection
- `POST /api/collections/import` - Import collection

### User Settings
- `GET /api/settings` - Get user settings (includes validate_recipe_volumes)
- `PUT /api/settings/volume-validation` - Update volume validation preference

---

## Technical Considerations

### Version Control Strategy
- **Immutable versions**: Once created, versions cannot be edited
- **Create new version**: To make changes, duplicate and modify
- **Active version**: One version marked as active per recipe
- **Version numbering**: Auto-increment (v1, v2, v3...)
- **Custom names**: Users can name versions ("Summer 2026")

### Volume Validation (Global User Preference)
- **Setting location**: `users.validate_recipe_volumes` boolean column
- **Default**: `false` (validation disabled)
- **When enabled**: 
  - Check accord availability before saving recipe version
  - Query accord's current `volume_ml`
  - Compare against ingredient's `quantity_ml`
  - If insufficient, return error with details
  - Error message: "Insufficient volume for accord 'X': need Y ml, have Z ml"
- **When disabled**: 
  - Skip all validation
  - Allow theoretical/planning recipes
  - Don't check accord availability at all
- **UI behavior**:
  - Settings page has toggle for "Validate accord volumes"
  - Recipe form shows warning icon if validation enabled and volume low
  - Can still save if user acknowledges warning

### Percentage Calculation
- **Automatic**: Calculate percentage from ml quantities
- **Validation**: Warn if percentages don't sum to 100%
- **Tolerance**: Allow ±1% tolerance for rounding
- **Display**: Show both ml and % in UI

### Cascade Deletes & Restrictions
- **Delete recipe** → CASCADE to all versions, ingredients, notes, tags (clean removal)
- **Delete version** → CASCADE to all ingredients for that version
- **Delete accord** → **RESTRICT** (prevent if used in any recipe ingredient)
  - Check `recipe_ingredients` table for accord_id usage
  - If found, return error: "Cannot delete accord: used in X recipe(s)"
  - Include list of recipe names using the accord
  - User must remove accord from all recipes before deletion
- **Delete collection** → Remove recipes from collection (don't delete the recipes themselves)

### Performance Optimization
- **Indexes**: On user_id, recipe_id, version_id, created_at
- **Pagination**: List endpoints support pagination
- **Caching**: Cache frequently accessed recipes in frontend
- **Lazy loading**: Load versions/ingredients on demand

---

## Notes and Considerations

### User Experience
- **Intuitive versioning**: Clear visual indicators for versions
- **Easy ingredient picker**: Autocomplete accord search
- **Real-time calculations**: Update percentages as quantities change
- **Helpful validations**: Guide users to valid recipes
- **Quick duplication**: One-click to create new version

### Data Integrity
- **Uniqueness**: Recipe names unique per user
- **Referential integrity**: Foreign keys with appropriate CASCADE/RESTRICT
- **Validation**: Server-side validation for all inputs
- **User isolation**: All queries filtered by user_id

### Future Enhancements
- **Recipe scaling**: Scale recipe to different volumes
- **Cost calculation**: If accords have prices
- **Sharing**: Public recipe database
- **Ratings**: Community ratings for recipes
- **Photos**: Add photos to recipes/versions
- **Aging tracking**: Track maturation time (maybe later)

### Migration Path
- No migration needed from existing data
- Fresh implementation alongside accord system
- Can reference existing accords
- No breaking changes to accord APIs

---

## Success Criteria

✅ **Backend Complete When:**
- All database tables created with proper constraints
- All repositories tested with 90%+ coverage
- All services tested with 85%+ coverage
- All handlers tested with 80%+ coverage
- API documentation complete
- All endpoints working and tested

✅ **Frontend Complete When:**
- All views and components implemented
- State management working correctly
- API integration complete
- Responsive design works on all screen sizes
- Empty states and loading states implemented
- Export/import functionality working

✅ **Feature Complete When:**
- Users can create and manage recipes
- Version control working correctly
- Ingredient management working
- Notes/journaling functional
- Collections working
- Export/import working
- All tests passing
- Documentation updated

---

## Estimated Timeline

This is a substantial feature addition. Rough estimates:

- **Phase 10.1**: 1-2 days (Models & schema)
- **Phase 10.2**: 3-4 days (Repositories + tests)
- **Phase 10.3**: 3-4 days (Services + tests)
- **Phase 10.4**: 2-3 days (Handlers + tests)
- **Phase 10.5**: 1 day (Routes & docs)
- **Phase 10.6**: 1 day (Frontend types)
- **Phase 10.7**: 2-3 days (Stores)
- **Phase 10.8**: 4-5 days (UI components)
- **Phase 10.9**: 2-3 days (Views)
- **Phase 10.10**: 2-3 days (Features & polish)
- **Phase 10.11**: 2-3 days (Testing)
- **Phase 10.12**: 1-2 days (Docs & deployment)

**Total**: ~25-35 days of focused development

---

## Design Decisions (Confirmed)

✅ **Volume Validation**: Global user preference (enable/disable for all recipes)
- Add `validate_recipe_volumes` boolean to users table
- Default: `false` (validation disabled)
- Users can toggle in settings
- When enabled, check accord availability before saving recipe version

✅ **Import Conflict Resolution**: Prompt user to rename
- On import, check for name conflicts
- If conflict exists, return error with suggested name (e.g., "Recipe Name (2)")
- User must provide new name to complete import
- No automatic renaming

✅ **No Limits**: 
- Unlimited recipes per user
- Unlimited versions per recipe
- Unlimited recipes per collection
- Unlimited collections per user
- Trust users to manage their own data

✅ **Accord Deletion Protection**: RESTRICT (prevent deletion)
- If accord is used in any recipe ingredient, deletion is blocked
- Return error message: "Cannot delete accord: used in X recipe(s)"
- List recipes using the accord
- User must remove accord from all recipes before deletion

✅ **Recipe Status**: User-defined only
- No approval workflows
- Simple status enum: draft, tested, finalized, archived
- Users control their own workflow

---

**Status**: Design decisions confirmed. Ready to begin implementation.
