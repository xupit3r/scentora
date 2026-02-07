# Data Models Specification

**Last Updated**: February 7, 2026
**Version**: 3.0 (Accord + Recipe System - Koa.js/TypeScript/Prisma)

---

## Overview

This document defines all data models and database schemas for the Scentora accord management system.

---

## Current Data Models (Phase 8+)

### Accord

The core entity representing a perfume accord (scent building block).

#### Database Schema

```sql
CREATE TABLE accords (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    pyramid_position VARCHAR(10) NOT NULL CHECK (pyramid_position IN ('top', 'middle', 'base')),
    volume_ml DECIMAL(10,2) NOT NULL CHECK (volume_ml >= 0),
    volume_drops INTEGER GENERATED ALWAYS AS (ROUND(volume_ml * 20)) STORED,
    supplier VARCHAR(255),
    purchase_date DATE,
    dilution_percentage DECIMAL(5,2) CHECK (dilution_percentage >= 0 AND dilution_percentage <= 100),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name, pyramid_position)
);

CREATE INDEX idx_accords_user_id ON accords(user_id);
CREATE INDEX idx_accords_position ON accords(pyramid_position);
CREATE INDEX idx_accords_created_at ON accords(created_at DESC);
```

#### Prisma Model

See `backend/prisma/schema.prisma` for the authoritative schema. Fields use camelCase with `@map("snake_case")` annotations.

#### TypeScript Interface

```typescript
interface Accord {
  _id: string;
  userId: string;
  name: string;
  pyramidPosition: 'top' | 'middle' | 'base';
  volumeMl: number;
  volumeDrops: number;
  supplier?: string;
  purchaseDate?: string; // ISO 8601 date
  dilutionPercentage?: number; // 0-100
  notes?: string;
  createdAt: string; // ISO 8601 timestamp
  updatedAt: string; // ISO 8601 timestamp
}
```

#### Field Descriptions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | UUID | Yes | Unique identifier |
| `user_id` | UUID | Yes | Owner of the accord |
| `name` | string(255) | Yes | Accord name (e.g., "Citrus Fresh") |
| `pyramid_position` | enum | Yes | Position: `top`, `middle`, or `base` |
| `volume_ml` | decimal(10,2) | Yes | Current volume in milliliters |
| `volume_drops` | integer | Auto | Calculated: ml × 20 |
| `supplier` | string(255) | No | Vendor/supplier name |
| `purchase_date` | date | No | When acquired |
| `dilution_percentage` | decimal(5,2) | No | Dilution (0-100%) |
| `notes` | text | No | Personal observations |
| `created_at` | timestamp | Yes | Creation timestamp |
| `updated_at` | timestamp | Yes | Last update timestamp |

#### Constraints

- **Unique**: `(user_id, name, pyramid_position)` - No duplicate accord names in same position per user
- **Check**: `volume_ml >= 0` - Volume cannot be negative
- **Check**: `pyramid_position IN ('top', 'middle', 'base')` - Valid position only
- **Check**: `dilution_percentage >= 0 AND dilution_percentage <= 100` - Valid percentage
- **Foreign Key**: `user_id` references `users(id)` with CASCADE delete

---

### AccordTag

Links accords to their descriptive tags (many-to-many relationship).

#### Database Schema

```sql
CREATE TABLE accord_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    accord_id UUID NOT NULL REFERENCES accords(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(accord_id, tag)
);

CREATE INDEX idx_accord_tags_accord_id ON accord_tags(accord_id);
CREATE INDEX idx_accord_tags_tag ON accord_tags(tag);
```

#### Constraints

- **Unique**: `(accord_id, tag)` - Same tag cannot be added twice to an accord
- **Foreign Key**: `accord_id` references `accords(id)` with CASCADE delete

---

### PredefinedTag

Predefined tag suggestions organized by category.

#### Database Schema

```sql
CREATE TABLE predefined_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(50) NOT NULL,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tag)
);

CREATE INDEX idx_predefined_tags_category ON predefined_tags(category);
```

#### TypeScript Interface

```typescript
interface PredefinedTag {
  _id: string;
  category: string;
  tag: string;
  createdAt: string;
}

interface PredefinedTagsByCategory {
  [category: string]: string[];
}
```

#### Seed Data

See [Tag System Specification](tag-system.md) for complete list of predefined tags.

---

## Response Models

### AccordResponse

Accord with tags array included (for API responses).

```typescript
interface AccordResponse extends Accord {
  tags: string[];
}
```

---

## Request DTOs

### CreateAccordRequest

Validated via Zod schemas in `backend/src/schemas/`.

```typescript
interface CreateAccordRequest {
  name: string;
  pyramidPosition: 'top' | 'middle' | 'base';
  volumeMl: number;
  supplier?: string;
  purchaseDate?: string; // YYYY-MM-DD
  dilutionPercentage?: number; // 0-100
  notes?: string;
  tags?: string[];
}
```

### UpdateAccordRequest

All fields optional (partial update).

```typescript
interface UpdateAccordRequest {
  name?: string;
  pyramidPosition?: 'top' | 'middle' | 'base';
  volumeMl?: number;
  supplier?: string;
  purchaseDate?: string;
  dilutionPercentage?: number;
  notes?: string;
  tags?: string[]; // Replaces all tags
}
```

---

## Existing Models (Retained)

### User

User accounts with authentication and preferences.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    validate_recipe_volumes BOOLEAN DEFAULT FALSE,  -- Phase 10: Volume validation preference
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**New in Phase 10:**
- `validate_recipe_volumes` - Global user preference for recipe volume validation
  - When `true`: Validates accord availability before saving recipe versions
  - When `false`: Allows theoretical/planning recipes without inventory checks
  - Default: `false`

### RefreshToken

JWT refresh tokens (unchanged).

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Invitation

Invitation codes (unchanged).

```sql
CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,
    used_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## Removed Models (Phase 8)

The following models were removed during the transition to accord system:

- **Perfume** - Replaced by Accord
- **JournalEntry** - Removed (was linked to perfumes)
- **Pyramid** - Removed (accords have position, not structure)

---

## Entity Relationships

```
┌─────────┐
│  users  │
└────┬────┘
     │
     │ 1:N
     │
     ├─────────┬─────────────┬──────────────┐
     │         │             │              │
     ▼         ▼             ▼              ▼
┌─────────┐ ┌─────────┐ ┌────────────┐ ┌──────────────┐
│ accords │ │ refresh │ │invitations │ │ invitations  │
│         │ │ tokens  │ │ (created)  │ │ (used_by)    │
└────┬────┘ └─────────┘ └────────────┘ └──────────────┘
     │
     │ 1:N
     │
     ▼
┌─────────────┐      N:1    ┌─────────────────┐
│ accord_tags │─────────────▶│ predefined_tags │
└─────────────┘              └─────────────────┘
```

---

## Validation Rules

### Accord Validation

- **Name**: Required, 1-255 characters, no special validation
- **Pyramid Position**: Required, must be 'top', 'middle', or 'base'
- **Volume ML**: Required, must be ≥ 0, up to 2 decimal places
- **Supplier**: Optional, max 255 characters
- **Purchase Date**: Optional, must be valid date (YYYY-MM-DD)
- **Dilution Percentage**: Optional, 0-100, up to 2 decimal places
- **Notes**: Optional, no length limit (TEXT field)
- **Tags**: Optional array, each tag 1-50 characters

### Business Rules

1. **Uniqueness**: `(name + position)` must be unique per user
2. **Volume**: Cannot be negative
3. **Tags**: Maximum 50 tags per accord (practical limit)
4. **Position Change**: Changing position requires uniqueness check
5. **Deletion**: Deleting accord cascades to all tags

---

## Implemented Models (Phase 10)

### Recipe

Complete perfume formulas created by combining multiple accords.

#### Database Schema

```sql
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
```

#### TypeScript Interface

```typescript
interface Recipe {
  _id: string;
  userId: string;
  name: string;
  description?: string;
  targetVolumeMl: number;
  status: 'draft' | 'tested' | 'finalized' | 'archived';
  activeVersionId?: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
}
```

#### Field Descriptions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | UUID | Yes | Unique identifier |
| `user_id` | UUID | Yes | Owner of the recipe |
| `name` | string(255) | Yes | Recipe name (unique per user) |
| `description` | text | No | Recipe description/inspiration |
| `target_volume_ml` | decimal(10,2) | Yes | Target volume for recipe |
| `status` | enum | Yes | Recipe status (draft, tested, finalized, archived) |
| `active_version_id` | UUID | No | Currently active version |
| `created_at` | timestamp | Yes | Creation timestamp |
| `updated_at` | timestamp | Yes | Last update timestamp |

---

### RecipeVersion

Immutable versions of recipes for iteration and version control.

#### Database Schema

```sql
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

-- Add foreign key after both tables exist
ALTER TABLE recipes 
    ADD CONSTRAINT fk_recipes_active_version 
    FOREIGN KEY (active_version_id) 
    REFERENCES recipe_versions(id) ON DELETE SET NULL;
```

#### TypeScript Interface

```typescript
interface RecipeVersion {
  _id: string;
  recipeId: string;
  versionNumber: number;
  name: string;
  notes?: string;
  isActive: boolean;
  createdAt: string;
}
```

---

### RecipeIngredient

Links accords to recipe versions with quantities.

#### Database Schema

```sql
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
```

#### TypeScript Interface

```typescript
interface RecipeIngredient {
  _id: string;
  versionId: string;
  accordId: string;
  quantityMl: number;
  quantityDrops: number;
  percentage: number;
  notes?: string;
  createdAt: string;
}
```

---

### RecipeNote

Notes and journal entries for recipes.

#### Database Schema

```sql
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
```

---

### RecipeCollection

Group recipes into collections/folders.

#### Database Schema

```sql
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

## Future Models (Phase 12+)

Potential additions: batch tracking, cost tracking, safety data (IFRA limits), maturation tracking.

---

## Migration Strategy

### From Perfume System to Accord System

**Drop Tables** (destructive):
```sql
DROP TABLE journal_entries CASCADE;
DROP TABLE perfumes CASCADE;
```

**Create Tables** (new):
```sql
CREATE TABLE accords (...);
CREATE TABLE accord_tags (...);
CREATE TABLE predefined_tags (...);
```

**Seed Data**:
```sql
INSERT INTO predefined_tags (category, tag) VALUES
    ('character', 'fresh'),
    ('character', 'warm'),
    -- ... (90+ tags)
```

**No Data Migration**: Fresh start with accord system.

---

## Notes

- All UUIDs generated via `gen_random_uuid()` (PostgreSQL) / Prisma `@default(uuid())`
- All timestamps in UTC
- Prisma schema uses camelCase fields with `@map("snake_case")` annotations
- JSON responses transform `id` → `_id`, Prisma Decimal → JS number, null omission (matching Go-era `omitempty`)
- API response shapes: accords wrapped `{ accord }`, recipes returned directly (no wrapper)
- Generated columns (`volume_drops`, `quantity_drops`) are DB-computed; API uses `Math.round(ml * 20)` fallback
- Soft deletes not implemented (hard deletes with CASCADE)
- Authoritative schema: `backend/prisma/schema.prisma` (14 models)
