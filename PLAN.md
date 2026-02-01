# Scentora - Comprehensive Development Plan

**Last Updated**: January 31, 2026  
**Status**: Phase 9 Complete - Testing Infrastructure | Phase 10 Planned - Recipe System

---

## 📁 Documentation Conventions

### Phase Completion Documents

**IMPORTANT for AI Agents:** All phase completion documents MUST be stored in the `docs/phases/` directory to keep the root directory clean.

**Location:** `/docs/phases/`

**Naming Convention:**
- `PHASE{X}_{Y}_COMPLETE.md` - Main phase completion summary
- `PHASE{X}_{Y}_IMPLEMENTATION_SUMMARY.md` - Detailed implementation notes
- `PHASE{X}_{Y}_TESTING_GUIDE.md` - Testing documentation
- `PHASE{X}_{Y}_PLAN_SUMMARY.md` - Planning documentation

**Example:**
```
docs/phases/PHASE8_1_COMPLETE.md
docs/phases/PHASE8_9_COMPLETE.md
docs/phases/PHASE8_9_1_COMPLETE.md
```

**When Creating Phase Documentation:**
1. Always create phase documents in `docs/phases/` directory
2. Use the naming convention above
3. Include comprehensive summaries of work completed
4. Document code changes, features added, and testing results
5. Link to related documentation in other directories

**Current Phase Documents:**
- Phase 8.1-8.9 completion documents are in `docs/phases/`
- See `docs/phases/` for all historical phase documentation

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Testing & Quality Assurance](#testing--quality-assurance) ⚠️ **CRITICAL**
3. [Past Work Summary](#past-work-summary)
4. [Current State](#current-state)
5. [Planned Work](#planned-work)
6. [Technical Specifications](#technical-specifications)
7. [Development Phases](#development-phases)

---

## Project Overview

**Scentora** is evolving from a perfume cataloging application to a comprehensive **perfume formulation and accord management system**.

### Vision
Enable perfumers and enthusiasts to:
- Manage their accord inventory (scent building blocks)
- Track volumes, suppliers, and characteristics
- Organize accords with rich metadata and tagging
- **Create and manage perfume recipes/formulas** using their accords (Phase 10)
- Version control recipes and track iterations
- Export and share formulas with the community

### Current State
- **✅ Phase 8 Complete**: Accord inventory management system fully operational
- **✅ Phase 9 Complete**: Comprehensive backend testing infrastructure (136 tests)
- **🎯 Phase 10 Next**: Recipe/formula creation system for combining accords

---

## Testing & Quality Assurance

### ⚠️ CRITICAL: Testing Philosophy

**Testing is mandatory and non-negotiable.** All code must be tested before committing and pushing.

### Testing Guidelines

#### 🔴 ALWAYS Test Before:
1. **Committing code** - Run relevant test suite
2. **Completing a phase** - Run full test suite
3. **Creating pull requests** - Verify all tests pass
4. **Deploying to production** - Full integration test suite
5. **Making breaking changes** - Update and run affected tests

#### 📋 Test Coverage Requirements

**Minimum Coverage Targets**:
- **Overall**: 80%+ code coverage
- **Repositories**: 90%+ coverage (data layer is critical)
- **Services**: 85%+ coverage (business logic must be solid)
- **Handlers**: 80%+ coverage (API endpoints)
- **Models**: 100% validation tests

**What to Test**:
- ✅ All database operations (CRUD)
- ✅ All business logic in services
- ✅ All API endpoints (happy path + error cases)
- ✅ All validation rules
- ✅ All authentication/authorization flows
- ✅ Database migrations and schema
- ✅ Error handling and edge cases
- ✅ Integration between layers

#### 🧪 Test Types

1. **Unit Tests** (`*_test.go`)
   - Test individual functions in isolation
   - Mock external dependencies
   - Fast execution (<100ms per test)
   - Examples: Model validation, utility functions

2. **Repository Tests** (`repository/*_test.go`)
   - Test database operations
   - Use test database (not production)
   - Test transactions and rollbacks
   - Examples: Create, Read, Update, Delete operations

3. **Service Tests** (`services/*_test.go`)
   - Test business logic
   - Mock repositories
   - Test error conditions
   - Examples: Auth flows, invitation validation

4. **Handler Tests** (`handlers/*_test.go`)
   - Test HTTP endpoints
   - Mock services
   - Test request/response formats
   - Examples: API endpoint behavior

5. **Integration Tests** (`integration/*_test.go`)
   - Test full request flow
   - Use test database
   - Test multi-layer interactions
   - Examples: Full auth flow, data persistence

#### 🏃 Running Tests

**Backend (Go)**:
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./internal/repository

# Run with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Frontend (Vue/Vitest)**:
```bash
# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch

# Run specific test file
npm test -- PerfumeCard.spec.ts
```

#### 📝 Test Documentation

Every test file should include:
- Clear test names describing what is being tested
- Setup/teardown for test fixtures
- Comments explaining complex test scenarios
- Examples of expected behavior

#### 🔄 Test Maintenance

**Keep Tests Updated**:
- Update tests when changing functionality
- Remove tests for deleted features
- Refactor tests when refactoring code
- Never skip failing tests - fix them or fix the code

**Red Flags** 🚩:
- Commented-out tests
- Skipped tests (`t.Skip()`)
- Tests that only pass sometimes (flaky tests)
- Tests with no assertions
- Tests that take too long (>5s for unit tests)

#### 📊 Test Reporting

After running tests, verify:
- ✅ All tests pass
- ✅ Coverage meets minimum thresholds
- ✅ No race conditions detected
- ✅ No memory leaks
- ✅ Performance benchmarks met

#### 🎯 Phase-Specific Testing

**Before Completing Each Phase**:
1. Write tests for all new code
2. Run full test suite
3. Verify coverage meets targets
4. Fix any failing tests
5. Document test approach in phase completion doc
6. Commit tests alongside implementation

**Phase 8 Testing Requirements**:
- Phase 8.1: Test database migrations, schema, seed data
- Phase 8.2: Test accord repositories and services
- Phase 8.3: Test search and filtering logic
- Phase 8.4: Test statistics and export/import
- Phase 8.5-8.8: Test frontend components and integration

---

## Past Work Summary

### Phase 1: Foundation (Completed - Jan 17, 2026)
**Goal**: Project scaffolding and basic infrastructure

**Deliverables**:
- ✅ Monorepo structure (backend + frontend)
- ✅ Node.js/Koa backend with TypeScript
- ✅ Vue.js 3 frontend with TypeScript
- ✅ CouchDB database setup via Docker
- ✅ Health check endpoint
- ✅ Basic routing and API service layer
- ✅ Error handling middleware
- ✅ Development environment with hot reload

**Technologies**:
- Backend: Koa.js 3.x, TypeScript 5.x, CouchDB (nano client)
- Frontend: Vue 3.5, Vite 6.x, Pinia 3.x, Vue Router 4.x
- Infrastructure: Docker Compose

---

### Phase 2: Core CRUD Operations (Completed - Jan 17, 2026)
**Goal**: Basic perfume and journal entry management

**Deliverables**:
- ✅ Perfume CRUD API endpoints
- ✅ Journal entry CRUD API endpoints
- ✅ Request validation with Zod schemas
- ✅ Perfume pyramid structure (top/middle/base notes)
- ✅ API service layer (frontend)
- ✅ Basic UI components (PerfumeForm, PerfumePyramid)
- ✅ CouchDB views and indexes

**Endpoints Added**:
- `GET/POST /api/perfumes`
- `GET/PUT/DELETE /api/perfumes/:id`
- `GET/POST /api/perfumes/:perfumeId/journal`
- `PUT/DELETE /api/journal/:id`

---

### Phase 3: Advanced Features (Completed - Jan 17, 2026)
**Goal**: Search, filtering, and enhanced UX

**Deliverables**:
- ✅ Full-text search across perfumes
- ✅ Multi-filter system (concentration, year, notes)
- ✅ Note autocomplete from existing collection
- ✅ Edit functionality for perfumes and journals
- ✅ Enhanced UI with hover actions
- ✅ Comma-separated note input with chips
- ✅ Visual pyramid display with color coding

**Features**:
- Search by name, designer, description
- Filter by concentration, year, specific note
- Autocomplete suggestions for notes
- Action buttons (edit, view, delete)

---

### Phase 4: Analytics & Export (Completed - Jan 17, 2026)
**Goal**: Collection insights and data portability

**Deliverables**:
- ✅ Statistics dashboard
- ✅ Visual charts (bar charts for designers, notes, years)
- ✅ Collection export as JSON
- ✅ Collection import from JSON
- ✅ Notes aggregation endpoint
- ✅ Statistics endpoint

**Dashboard Metrics**:
- Total perfumes, journal entries
- Top designers, most used notes
- Concentration distribution
- Release year timeline
- Average ratings

---

### Phase 5: Authentication System (Completed - Jan 18, 2026)
**Goal**: Multi-user support with secure authentication

**Deliverables**:
- ✅ JWT-based authentication
- ✅ User registration and login
- ✅ Password hashing with bcrypt
- ✅ Protected API routes
- ✅ User data isolation (userId in queries)
- ✅ Login/register UI pages
- ✅ Auth store with Pinia
- ✅ Route guards for protection

**Security**:
- JWT tokens (7-day expiry initially)
- Bcrypt password hashing (10 rounds)
- Email and username uniqueness
- Auth middleware on all data endpoints

---

### Phase 5B: Refresh Tokens & Rate Limiting (Completed - Jan 18, 2026)
**Goal**: Enhanced security and session management

**Deliverables**:
- ✅ JWT refresh token system
- ✅ Token rotation on refresh
- ✅ Short-lived access tokens (15 minutes)
- ✅ Long-lived refresh tokens (7 days)
- ✅ Rate limiting on auth endpoints (5 req/15min)
- ✅ Rate limiting on general endpoints (100 req/min)
- ✅ Automatic token refresh in frontend
- ✅ Logout and logout-all functionality

**Security Improvements**:
- Reduced attack surface (short access tokens)
- Token revocation capability
- Brute force protection via rate limiting
- Session management across devices

---

### Phase 6: Invitation System (Completed - Jan 30, 2026)
**Goal**: Invitation-only registration for controlled growth

**Deliverables**:
- ✅ Invitation code generation
- ✅ Email-specific invitations (optional)
- ✅ Invitation expiration handling
- ✅ Registration requires valid invitation code
- ✅ Invitation management UI (create, list, revoke)
- ✅ Invitation tracking (used by, used at)

**Endpoints**:
- `POST /api/invitations` - Create invitation
- `GET /api/invitations` - List user's invitations
- `DELETE /api/invitations/:code` - Revoke invitation

---

### Phase 7: Migration to Go Backend (Completed - Jan 30, 2026)
**Goal**: Rewrite backend in Go for better performance and maintainability

**Deliverables**:
- ✅ Go/Echo backend with PostgreSQL
- ✅ 100% API compatibility with Node.js version
- ✅ All features ported (auth, invitations, perfumes, journals)
- ✅ Database migrations (auto-run on startup)
- ✅ Backend parity features:
  - Logout from all devices
  - Collection import
  - Export format alignment
  - Rate limiting middleware
  - Optional auth middleware
- ✅ Comprehensive test suite
- ✅ Updated launcher scripts
- ✅ Node.js backend removed

**Tech Stack After Migration**:
- Backend: Go 1.21+, Echo v4, PostgreSQL 15
- Database: PostgreSQL (from CouchDB)
- ORM: sqlx
- JWT: golang-jwt/jwt
- Validation: go-playground/validator

**Performance Gains**:
- Startup: ~100ms (vs ~2s Node.js)
- Memory: ~20MB (vs ~80MB Node.js)
- Build: Single binary (vs node_modules)

---

## Current State

### Architecture (As of Jan 31, 2026)

**Backend**: `backend/` (Go)
- Framework: Echo v4
- Database: PostgreSQL 15
- Port: 3000
- Location: `/home/joe/code/scentora/backend/`

**Frontend**: `frontend/` (Vue.js)
- Framework: Vue 3.5 + TypeScript
- Build: Vite 6.x
- Port: 5173
- Location: `/home/joe/code/scentora/frontend/`

**Database**: PostgreSQL
- Host: localhost:5435 (Docker)
- Database: scentora
- Credentials: admin/password (dev)

### Current Features
- ✅ User authentication (JWT with refresh tokens)
- ✅ Invitation-only registration
- ✅ Perfume CRUD operations
- ✅ Journal entry tracking
- ✅ Search and filtering
- ✅ Statistics dashboard
- ✅ Data export/import
- ✅ Rate limiting
- ✅ Multi-user data isolation

### Current Tables
- `users` - User accounts
- `perfumes` - Perfume catalog
- `journal_entries` - Journal entries
- `refresh_tokens` - JWT refresh tokens
- `invitations` - Invitation codes

---

## Planned Work

### Phase 8: Transition to Accord Management System (In Planning)
**Goal**: Transform from perfume catalog to accord inventory system

**Strategic Pivot**: 
The application will shift from tracking commercial perfumes to managing DIY perfume creation components (accords) and recipes.

#### What's Changing
**Remove**:
- Perfume tracking (entire feature set)
- Journal entries linked to perfumes
- Perfume-related UI components
- Perfume statistics

**Add**:
- Accord inventory management
- Rich tagging system for accords
- Volume tracking (ml + drops)
- Supplier and purchase tracking
- Recipe system (future phase)

#### Key Concepts

**Accord**: A harmonious blend of aromatic materials that creates a distinct scent impression. Accords are the building blocks used in perfume recipes.

**Pyramid Position**: Where an accord sits in the fragrance structure:
- **Top**: Volatile, first impression (citrus, herbs)
- **Middle**: Heart notes, core character (florals, spices)
- **Base**: Foundation, lasting (woods, resins, musks)

**Volume Tracking**: Dual measurement system:
- Primary: Milliliters (ml) - precise liquid measurement
- Secondary: Drops - practical measurement (1 ml ≈ 20 drops)

#### Accord Data Model

**Core Properties**:
- Name (string, required)
- Pyramid position (enum: top/middle/base, required)
- Volume in ML (decimal, required)
- Volume in drops (calculated: ml × 20)
- Supplier/vendor (string, optional)
- Purchase date (date, optional)
- Dilution percentage (decimal 0-100%, optional)
- Notes/description (text, optional)
- Tags (array of strings)

**Unique Constraint**: `(name + pyramid_position)` per user
- Example: Can't have two "Citrus Accord" in "top" position
- But can have "Citrus Accord" in both "top" and "middle"

#### Tag System

**Default Tag Categories** (90+ predefined tags):
1. **Character**: fresh, warm, cool, dry, powdery, creamy, sharp, soft, rich, light
2. **Mood**: romantic, sensual, energetic, calming, mysterious, playful, sophisticated, innocent
3. **Season**: spring, summer, autumn, winter
4. **Time**: morning, afternoon, evening, night
5. **Intensity**: subtle, moderate, strong, intense, bold
6. **Quality**: clean, dirty, animalic, synthetic, natural, modern, vintage
7. **Scent Family**: floral, fruity, woody, oriental, fresh, aromatic, spicy, gourmand
8. **Texture**: smooth, rough, silky, velvety, airy, dense
9. **Style**: casual, formal, sporty, elegant, edgy

**Custom Tags**: Users can create their own tags beyond predefined options.

#### Database Schema Changes

**Drop Tables**:
```sql
DROP TABLE journal_entries;
DROP TABLE perfumes;
```

**New Tables**:

```sql
-- Accords table
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

-- Accord tags (many-to-many)
CREATE TABLE accord_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    accord_id UUID NOT NULL REFERENCES accords(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(accord_id, tag)
);

-- Predefined tags (seeded data)
CREATE TABLE predefined_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(50) NOT NULL,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(tag)
);
```

**Indexes**:
```sql
CREATE INDEX idx_accords_user_id ON accords(user_id);
CREATE INDEX idx_accords_position ON accords(pyramid_position);
CREATE INDEX idx_accords_created_at ON accords(created_at DESC);
CREATE INDEX idx_accord_tags_accord_id ON accord_tags(accord_id);
CREATE INDEX idx_accord_tags_tag ON accord_tags(tag);
CREATE INDEX idx_predefined_tags_category ON predefined_tags(category);
```

#### Backend Implementation Plan

**Phase 8.1: Database & Cleanup**
- [ ] Create migration to drop old tables
- [ ] Create migration to add new tables
- [ ] Seed predefined_tags table with 90+ tags
- [ ] Remove perfume/journal repositories
- [ ] Remove perfume/journal services
- [ ] Remove perfume/journal handlers
- [ ] Update models (remove old, add Accord types)
- [ ] Update routes

**Phase 8.2: Accord Core Features**
- [ ] Create `accord_repo.go` (CRUD operations)
- [ ] Create `accord_tag_repo.go` (tag management)
- [ ] Create `predefined_tag_repo.go` (tag lookup)
- [ ] Create `accord_service.go` (business logic)
- [ ] Create `accord.go` handler (HTTP endpoints)
- [ ] Create `tags.go` handler (tag endpoints)
- [ ] Add routes for `/api/accords` and `/api/tags`

**Phase 8.3: Search & Filter**
- [ ] Implement filtering by position
- [ ] Implement filtering by volume range
- [ ] Implement filtering by tags
- [ ] Implement search by name/notes
- [ ] Implement filter by supplier
- [ ] Implement low stock filtering

**Phase 8.4: Statistics & Export**
- [ ] Update stats handler (accord metrics)
- [ ] Update export handler (accord format)
- [ ] Update import handler (accord import)
- [ ] Add low inventory alerts

**API Endpoints** (New):
```
GET    /api/accords              - List accords (with filters)
POST   /api/accords              - Create accord
GET    /api/accords/:id          - Get accord details
PUT    /api/accords/:id          - Update accord
DELETE /api/accords/:id          - Delete accord
GET    /api/accords/:id/tags     - Get accord tags
POST   /api/accords/:id/tags     - Add tag to accord
DELETE /api/accords/:id/tags/:tag - Remove tag from accord
GET    /api/tags/predefined      - Get predefined tags
GET    /api/tags                 - Get all user tags
GET    /api/stats                - Accord statistics
GET    /api/export               - Export accords
POST   /api/export/import        - Import accords
```

#### Frontend Implementation Plan

**Phase 8.5: Frontend Cleanup**
- [ ] Remove `PerfumeDetail.vue`
- [ ] Remove `PerfumeForm.vue`
- [ ] Remove `PerfumePyramid.vue`
- [ ] Remove `SearchFilters.vue`
- [ ] Remove `Statistics.vue` (will recreate)
- [ ] Update `Home.vue` (remove perfume grid)
- [ ] Update types (remove Perfume, Journal)
- [ ] Update API service (remove perfume methods)
- [ ] Update router (remove perfume routes)

**Phase 8.6: New Components**
- [ ] Create `AccordCard.vue` - Display accord in grid
- [ ] Create `AccordForm.vue` - Add/edit accord modal
- [ ] Create `AccordDetail.vue` - Full accord view
- [ ] Create `AccordFilters.vue` - Filter sidebar/panel
- [ ] Create `TagSelector.vue` - Tag autocomplete component

**Phase 8.7: New Views**
- [ ] Update `Home.vue` - Accord inventory grid
- [ ] Create `AccordDetailView.vue` - Detailed accord page
- [ ] Create `Statistics.vue` (simplified) - Accord stats
- [ ] Update `About.vue` - New app description

**Phase 8.8: Features & Polish**
- [ ] Implement tag autocomplete
- [ ] Implement filter combinations
- [ ] Add volume conversion display
- [ ] Add low stock warnings
- [ ] Add confirmation dialogs
- [ ] Responsive design testing
- [ ] Mobile UI optimization

#### UI/UX Specifications

**Color Scheme**:
- Top notes: Yellow gradient (#FFD93D → #FFA800)
- Middle notes: Purple gradient (#B565D8 → #8B5CF6)
- Base notes: Brown gradient (#A0826D → #6B4423)
- Accent: Teal (#14B8A6)
- Background: Light gray (#F5F5F5)

**Components**:
- **AccordCard**: Name, position badge, volume, tags, hover actions
- **AccordForm**: Multi-section form (Basic/Inventory/Tags/Notes)
- **TagSelector**: Autocomplete with grouped suggestions
- **AccordFilters**: Collapsible sidebar with multiple filter types

**Interactions**:
- Add/Edit: Modal form with tabbed sections
- Delete: Confirmation dialog
- Tags: Autocomplete dropdown with custom tag creation
- Filter: Slide-in panel (mobile) or sidebar (desktop)
- Low volume: Orange (< 5ml) or red (< 1ml) badge

**Search & Filter**:
- Search: Name, notes, supplier (fuzzy match)
- Filters: Position, volume range, tags, supplier, low stock
- Sort: Name, date added, volume, position

**Statistics** (Simplified):
- Total accord count
- Count by pyramid position
- Total volume by position
- Most used tags (top 20)
- Low inventory alerts

#### Export/Import Format

```json
{
  "version": "2.0",
  "exportDate": "2026-01-31T15:00:00Z",
  "accords": [
    {
      "name": "Citrus Fresh Accord",
      "pyramidPosition": "top",
      "volumeMl": 25.5,
      "volumeDrops": 510,
      "supplier": "Perfumer's Apprentice",
      "purchaseDate": "2025-12-15",
      "dilutionPercentage": 10,
      "notes": "Very bright and zesty. Works well in summer blends.",
      "tags": ["fresh", "energetic", "citrus", "summer", "morning"],
      "createdAt": "2025-12-15T10:30:00Z",
      "updatedAt": "2026-01-20T14:00:00Z"
    }
  ]
}
```

#### Testing Checklist

**Backend**:
- [ ] Create accord (valid data)
- [ ] Create accord (duplicate name+position, should fail)
- [ ] Update accord (all fields)
- [ ] Delete accord
- [ ] List with filters (position, volume, tags)
- [ ] Search (name, notes)
- [ ] Add/remove tags
- [ ] Get predefined tags
- [ ] Export/import accords
- [ ] Statistics calculation

**Frontend**:
- [ ] Display accord grid
- [ ] Create new accord (form validation)
- [ ] Edit accord
- [ ] Delete accord (with confirmation)
- [ ] Filter by position/volume/tags
- [ ] Search by name
- [ ] Tag autocomplete
- [ ] Add custom tag
- [ ] View accord detail
- [ ] Export/import
- [ ] View statistics
- [ ] Mobile responsive

---

### Phase 8.9: Notion-Inspired UI/UX Redesign

**Status**: Planned  
**Priority**: High  
**Timeline**: 1-2 weeks  
**Goal**: Transform Scentora into a clean, minimalist interface inspired by Notion's design principles

#### Design Framework Selection

After evaluating 6 Vue UI frameworks, the recommended approach is:

**Selected Stack**: **Naive UI + Tailwind CSS** (Hybrid Approach)

**Framework Comparison**:
- ✅ **Naive UI** (9/10) - Clean, minimalist, TypeScript-first, lightweight
- Headless UI (10/10) - Maximum control but more development time
- PrimeVue (6/10) - Requires heavy theming
- Element Plus (5/10) - Too opinionated
- Vuetify (4/10) - Material Design conflicts with Notion aesthetic
- Ant Design Vue (5/10) - Different design language

**Rationale**:
- Naive UI provides clean, minimalist base components
- Tailwind CSS for utility-first custom styling
- Best balance of speed vs. control
- Excellent TypeScript support
- Lightweight and performant

#### Design Principles to Implement

**1. Clean Typography**
- System font stack with excellent readability
- Clear hierarchy (titles, headings, body, captions)
- Generous line height (1.5-1.7)
- Restrained font sizes (14px base, 16px for comfortable reading)

**2. Minimalist Color Palette**
```javascript
// Notion-inspired colors
text: {
  primary: '#37352F',
  secondary: '#787774',
  tertiary: '#9B9A97',
}
background: {
  primary: '#FFFFFF',
  secondary: '#FAFAFA',
  tertiary: '#F7F6F3',
}
border: {
  default: '#E9E9E7',
}
accent: {
  primary: '#0F766E', // Teal for actions
}
```

**3. Spacing System**
- Consistent 8px grid system
- Generous whitespace for breathing room
- 48-64px padding for main containers
- 24-32px gaps between major sections
- 8-16px gaps between related items

**4. Interaction Design**
- Subtle hover states (background changes, not borders)
- Smooth transitions (150-200ms)
- Inline editing (click to edit in place)
- Keyboard shortcuts
- Delayed tooltips (500ms)

**5. Navigation**
- Fixed sidebar (collapsible on mobile)
- Clear hierarchy (workspace → pages → sub-pages)
- Icon + label for clarity
- Active state highlighting
- Breadcrumbs for deep navigation

**6. Component Design**
- Cards with subtle shadows (0 1px 3px rgba(0,0,0,0.12))
- Rounded corners (6-8px, not too rounded)
- Ghost buttons (transparent until hover)
- Inline actions that appear on hover
- Modals with backdrop blur

#### Implementation Sub-Phases

**Phase 8.9.1: Foundation & Setup** (2-3 days)
- [ ] Install Naive UI and Tailwind CSS
- [ ] Create design system (`design-system/tokens.ts`)
- [ ] Configure Tailwind with custom theme
- [ ] Set up global styles and typography
- [ ] Create base component wrappers

**Phase 8.9.2: Navigation & Layout** (2-3 days)
- [ ] Create `AppSidebar.vue` component
  - Logo area at top
  - Navigation items with icons (Home, Statistics, Settings)
  - User profile at bottom
  - Collapse/expand functionality
- [ ] Restructure `App.vue` for sidebar layout
- [ ] Implement responsive sidebar (collapse on mobile)
- [ ] Add breadcrumb navigation
- [ ] Create header with page title + actions

**Phase 8.9.3: Core Components Redesign** (3-4 days)
- [ ] Redesign `AccordCard.vue`
  - Simplify visual hierarchy
  - Subtle position indicator (left border or emoji)
  - Actions appear on hover only
  - Smooth hover lift effect
- [ ] Redesign `AccordForm.vue`
  - Single-column layout
  - Inline section headers
  - Better input focus states
  - Auto-save draft support
  - Keyboard shortcuts (Cmd+Enter to save)
- [ ] Redesign `TagSelector.vue`
  - Notion-style dropdown
  - Grouped categories with headers
  - Keyboard navigation
- [ ] Redesign `AccordFilters.vue`
  - Collapsible sections
  - Clean checkbox/radio styles
  - Subtle active states

**Phase 8.9.4: Advanced Interactions** (2-3 days)
- [ ] Implement inline editing
  - Click accord name to edit inline
  - Click volume to quick-adjust
  - Auto-save after delay
- [ ] Add keyboard shortcuts
  - `/` - Focus search
  - `N` - New accord
  - `?` - Show shortcuts help
  - `Esc` - Close modals/cancel
- [ ] Standardize hover effects (150ms ease)
- [ ] Implement skeleton screens (not spinners)
- [ ] Add optimistic UI updates

**Phase 8.9.5: Typography & Spacing Refinement** (1-2 days)
- [ ] Define type scale (12, 14, 16, 20, 24, 32, 40px)
- [ ] Set line heights (1.4 for headings, 1.6 for body)
- [ ] Review and apply 8px grid consistently
- [ ] Increase whitespace strategically
- [ ] Replace vibrant colors with subtle tones

**Phase 8.9.6: Empty States & Feedback** (1 day)
- [ ] Create beautiful empty states
  - Empty inventory
  - Empty search results
  - Empty filters
- [ ] Implement toast notifications
- [ ] Add loading indicators
- [ ] Create form validation feedback

**Phase 8.9.7: Polish & Testing** (2-3 days)
- [ ] Responsive testing (desktop, tablet, mobile)
- [ ] Accessibility audit (keyboard nav, ARIA labels, contrast)
- [ ] Performance optimization (lazy load, bundle size)
- [ ] Browser testing (Chrome, Firefox, Safari)
- [ ] Lighthouse audit (score > 90)

#### New File Structure

```
frontend/src/
├── design-system/
│   ├── tokens.ts              # Design tokens
│   ├── theme.ts               # Naive UI theme config
│   └── tailwind.config.js     # Tailwind config
│
├── components/
│   ├── layout/
│   │   ├── AppSidebar.vue     # NEW: Sidebar navigation
│   │   ├── AppHeader.vue      # NEW: Page header
│   │   └── AppBreadcrumbs.vue # NEW: Breadcrumbs
│   │
│   ├── accord/
│   │   ├── AccordCard.vue     # REDESIGN
│   │   ├── AccordForm.vue     # REDESIGN
│   │   ├── AccordGrid.vue     # NEW: Grid container
│   │   └── AccordDetail.vue   # NEW: Detail view
│   │
│   ├── ui/
│   │   ├── TagSelector.vue    # REDESIGN
│   │   ├── FilterPanel.vue    # REDESIGN
│   │   ├── EmptyState.vue     # NEW: Reusable empty state
│   │   └── SkeletonCard.vue   # NEW: Loading skeleton
│   │
│   └── common/
│       ├── Button.vue         # NEW: Custom button wrapper
│       ├── Input.vue          # NEW: Custom input wrapper
│       └── Modal.vue          # NEW: Custom modal wrapper
│
├── composables/
│   ├── useKeyboard.ts         # NEW: Keyboard shortcuts
│   ├── useToast.ts            # NEW: Toast notifications
│   └── useTheme.ts            # NEW: Theme switching
│
└── App.vue                     # MAJOR REDESIGN
```

#### Success Criteria

**Visual Quality**:
- Clean, uncluttered interface
- Consistent spacing throughout
- Professional typography
- Subtle, purposeful colors
- Smooth animations (no jank)

**Usability**:
- Intuitive navigation
- Fast, responsive interactions
- Clear feedback for actions
- Helpful empty states
- Keyboard accessible

**Technical**:
- TypeScript type safety maintained
- Bundle size < 500KB (gzipped)
- Lighthouse score > 90
- No accessibility violations
- Works on all modern browsers

#### Migration Strategy

**Approach**: Gradual Migration
1. Set up framework alongside existing components
2. Migrate App.vue first (navigation/layout)
3. Redesign one component at a time
4. Keep old components until new ones are ready
5. Remove old code once migration complete

**Risk Mitigation**:
- Feature flags for new UI
- Git branches for each phase
- Regular user testing
- Rollback plan if issues arise
- Performance monitoring

#### Detailed Design Specifications

For complete design specifications including:
- Full color palette definitions
- Typography scale and font weights
- Spacing system (8px grid)
- Shadow definitions
- Component mockups
- Interaction patterns

See: `/home/joe/.copilot/session-state/*/plan.md` (Session planning document)

---

## Technical Specifications

For detailed technical specifications, see the `specs/` directory:

- **[Data Models](specs/data-models.md)** - Database schemas and entity relationships
- **[API Specification](specs/api-spec.md)** - All API endpoints with request/response formats
- **[UI/UX Specification](specs/ui-ux-spec.md)** - Component designs and interaction patterns
- **[Tag System](specs/tag-system.md)** - Predefined tags and categorization

---

## Development Phases

### Completed Phases (1-7)
All previous phases are complete and documented in historical files:
- PHASE1_COMPLETE.md through PHASE5B_COMPLETE.md
- BACKEND_PARITY_COMPLETE.md
- IMPLEMENTATION_COMPLETE.md
- FINAL_MIGRATION_STATUS.md

### Completed Phase (8): Accord System ✅
**Status**: Complete  
**Completion Date**: January 2026

**Sub-phases Completed**:
1. ✅ Database cleanup & migration
2. ✅ Backend core (repos, services, handlers)
3. ✅ Backend features (search, filter, stats)
4. ✅ Frontend cleanup
5. ✅ Frontend components
6. ✅ Frontend views
7. ✅ Integration testing
8. ✅ Notion-inspired UI redesign
9. ✅ Statistics & export features

### Completed Phase (9): Testing & Quality Assurance ✅
**Status**: Complete  
**Completion Date**: January 31, 2026  
**Achievement**: 136 passing tests with 59% backend coverage

**Sub-phases Completed**:
1. ✅ Backend Repository Tests (40 tests, 59.9% coverage)
2. ✅ Backend Service Tests (70 tests, 59.6% coverage)
3. ✅ Backend Handler Tests (24 tests, ~75% coverage)
4. ⏭️ Frontend Testing Setup (deferred to future)
5. ⏭️ E2E Testing (deferred to future)

### Active Phase (10): Recipe/Formula System 📋
**Status**: Planning Complete - Design Decisions Confirmed  
**Start Date**: February 1, 2026  
**Estimated Duration**: 25-35 days

**Overview**: 
Enable users to create perfume recipes/formulas by combining their accords with version control, notes, and organization features.

**Core Features**:
- ✨ Create recipes with target volumes
- 🔄 Version control (immutable versions, auto-activate new versions)
- 📊 Ingredients with quantities (ml, drops, percentages)
- ⚙️ Configurable volume validation (disabled by default, opt-in via user preferences)
- 📝 Recipe notes and journaling
- 🏷️ Tags and collections for organization
- 📤 Export/import recipes as JSON
- 🛡️ Protect accords in use from deletion (prevent with error message)

**Confirmed Design Decisions** (Feb 1, 2026):
1. **Volume Validation**: Disabled by default (users opt-in via preferences)
   - Allows theoretical/planning recipes without inventory concerns
   - Can be enabled per-user in preferences
2. **Version Activation**: New versions automatically become the active version
   - Latest version is typically the one being worked on
   - Can be manually changed if needed
3. **Accord Deletion Protection**: Prevent deletion if accord used in recipes
   - Return 409 Conflict error with list of recipes using the accord
   - User must remove from recipes first or delete recipes
4. **Recipe Status Values**: Five states for granular workflow tracking
   - `draft` - Initial creation, work in progress
   - `in_progress` - Actively developing/testing
   - `tested` - Has been tested, results documented
   - `finalized` - Complete, production-ready formula
   - `archived` - No longer active, kept for reference
5. **Implementation Approach**: Backend-first for cleaner separation
   - Complete backend (Phases 10.1-10.5)
   - Then frontend (Phases 10.6-10.10)
   - Then testing (Phase 10.11)

**Sub-phases** (12 phases):
1. Backend - Data Models & Database Schema
2. Backend - Repository Layer
3. Backend - Service Layer
4. Backend - Handler Layer
5. Backend - Routes & Middleware
6. Frontend - Data Models & Types
7. Frontend - State Management
8. Frontend - UI Components
9. Frontend - Views & Pages
10. Frontend - Features & Polish
11. Testing (Backend + Frontend + E2E)
12. Documentation & Deployment

**Detailed Plan**: See [specs/recipe-system.md](specs/recipe-system.md)

### Future Phases (11+)

**Phase 11: Advanced Recipe Features** (Future)
- Batch creation tracking
- Cost tracking per accord/recipe
- Ingredient sourcing database
- Safety data (IFRA limits, allergens)
- Maturation tracking
- Label generation

---

## Success Criteria

### Phase 8 Success Criteria
✅ All perfume/journal code removed  
✅ Accord CRUD fully functional  
✅ Tag system working (predefined + custom)  
✅ Search and filtering smooth  
✅ Volume tracking accurate (ml + drops)  
✅ Export/import preserves all data  
✅ UI intuitive and visually appealing  
✅ Mobile responsive  
✅ No breaking errors  
✅ Documentation updated  

### Overall Project Success
- Clean, maintainable codebase
- Production-ready backend (Go)
- Modern, responsive frontend (Vue 3)
- Comprehensive documentation
- Secure authentication
- Good performance (fast response times)
- Extensible architecture for future features

---

## Resources & Documentation

### Key Documents
- **PLAN.md** (this file) - Comprehensive development plan
- **README.md** - Project overview and quick start
- **QUICKSTART.md** - Quick start guide
- **docs/AUTH_IMPLEMENTATION.md** - Authentication system
- **docs/REFRESH_TOKENS_RATE_LIMITING.md** - Security features
- **docs/LAUNCHER_GUIDE.md** - Launcher script usage
- **specs/** - Detailed technical specifications

### Historical Documents
- specs/archives/PHASE*.md files - Completed phase documentation
- specs/archives/MIGRATION_TO_GO.md - Backend migration details
- specs/archives/API_REFERENCE.md - Old perfume API documentation
- docs/AUTH_IMPLEMENTATION.md - Authentication system
- docs/REFRESH_TOKENS_RATE_LIMITING.md - Security features

### Developer Resources
- `backend/README.md` - Go backend documentation
- `backend/QUICKSTART.md` - Backend quick start
- `frontend/package.json` - Frontend dependencies and scripts
- `.env.example` - Environment configuration template

---

## Notes

### Design Philosophy
- **Simplicity**: Features should be intuitive and uncluttered
- **Flexibility**: System should accommodate various workflows
- **Performance**: Fast response times and minimal loading
- **Security**: User data protected and isolated
- **Extensibility**: Easy to add new features

### Technical Decisions
- **Go Backend**: Better performance, type safety, easier deployment
- **PostgreSQL**: Relational model better for accord relationships
- **JWT Auth**: Stateless, scalable authentication
- **Vue 3 Composition API**: Modern, reactive, maintainable
- **Tag System**: Flexible categorization without rigid constraints

### Future Considerations
- Mobile app (React Native or Flutter)
- AI-powered accord recommendations
- Community recipe sharing
- Supplier marketplace integration
- Inventory alerts and reordering
- Barcode scanning for products

---

**Project**: Scentora - Accord Inventory & Recipe Manager  
**Repository**: github.com/xupit3r/scentora  
**Maintainer**: Joe (xupit3r)  
**License**: MIT
