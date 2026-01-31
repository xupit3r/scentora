# Data Models Specification

**Last Updated**: January 31, 2026  
**Version**: 2.0 (Accord System)

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

#### Go Struct

```go
type Accord struct {
    ID                 string    `json:"_id" db:"id"`
    UserID             string    `json:"userId" db:"user_id"`
    Name               string    `json:"name" db:"name"`
    PyramidPosition    string    `json:"pyramidPosition" db:"pyramid_position"`
    VolumeMl           float64   `json:"volumeMl" db:"volume_ml"`
    VolumeDrops        int       `json:"volumeDrops" db:"volume_drops"`
    Supplier           *string   `json:"supplier,omitempty" db:"supplier"`
    PurchaseDate       *string   `json:"purchaseDate,omitempty" db:"purchase_date"`
    DilutionPercentage *float64  `json:"dilutionPercentage,omitempty" db:"dilution_percentage"`
    Notes              *string   `json:"notes,omitempty" db:"notes"`
    CreatedAt          time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt          time.Time `json:"updatedAt" db:"updated_at"`
}
```

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

#### Go Struct

```go
type AccordTag struct {
    ID        string    `json:"_id" db:"id"`
    AccordID  string    `json:"accordId" db:"accord_id"`
    Tag       string    `json:"tag" db:"tag"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
```

#### TypeScript Interface

```typescript
interface AccordTag {
  _id: string;
  accordId: string;
  tag: string;
  createdAt: string;
}
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

#### Go Struct

```go
type PredefinedTag struct {
    ID        string    `json:"_id" db:"id"`
    Category  string    `json:"category" db:"category"`
    Tag       string    `json:"tag" db:"tag"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
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

```go
type AccordResponse struct {
    Accord
    Tags []string `json:"tags"`
}
```

```typescript
interface AccordResponse extends Accord {
  tags: string[];
}
```

---

## Request DTOs

### CreateAccordRequest

```go
type CreateAccordRequest struct {
    Name               string   `json:"name" validate:"required,min=1,max=255"`
    PyramidPosition    string   `json:"pyramidPosition" validate:"required,oneof=top middle base"`
    VolumeMl           float64  `json:"volumeMl" validate:"required,gte=0"`
    Supplier           *string  `json:"supplier" validate:"omitempty,max=255"`
    PurchaseDate       *string  `json:"purchaseDate" validate:"omitempty,datetime=2006-01-02"`
    DilutionPercentage *float64 `json:"dilutionPercentage" validate:"omitempty,gte=0,lte=100"`
    Notes              *string  `json:"notes" validate:"omitempty"`
    Tags               []string `json:"tags" validate:"omitempty,dive,min=1,max=50"`
}
```

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

```go
type UpdateAccordRequest struct {
    Name               *string  `json:"name" validate:"omitempty,min=1,max=255"`
    PyramidPosition    *string  `json:"pyramidPosition" validate:"omitempty,oneof=top middle base"`
    VolumeMl           *float64 `json:"volumeMl" validate:"omitempty,gte=0"`
    Supplier           *string  `json:"supplier" validate:"omitempty,max=255"`
    PurchaseDate       *string  `json:"purchaseDate" validate:"omitempty,datetime=2006-01-02"`
    DilutionPercentage *float64 `json:"dilutionPercentage" validate:"omitempty,gte=0,lte=100"`
    Notes              *string  `json:"notes" validate:"omitempty"`
    Tags               []string `json:"tags" validate:"omitempty,dive,min=1,max=50"`
}
```

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

User accounts (unchanged from previous phases).

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

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

## Future Models (Phase 9+)

### Recipe (Planned)

```sql
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_volume_ml DECIMAL(10,2),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### RecipeIngredient (Planned)

```sql
CREATE TABLE recipe_ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    accord_id UUID NOT NULL REFERENCES accords(id) ON DELETE RESTRICT,
    quantity_ml DECIMAL(10,2) NOT NULL,
    quantity_drops INTEGER,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

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

- All UUIDs generated via `gen_random_uuid()` (PostgreSQL)
- All timestamps in UTC
- JSON responses use camelCase for compatibility with frontend
- Database columns use snake_case per PostgreSQL convention
- Soft deletes not implemented (hard deletes with CASCADE)
- Audit logging not implemented (future consideration)
