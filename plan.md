# Scentora - Comprehensive Development Plan

**Last Updated**: February 10, 2026
**Status**: Phase 13 In Progress - CLI Console (Planning/Documentation)

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
- **✅ Phase 9 Complete**: Backend rewritten to Koa.js/TypeScript/Prisma, 70 integration tests
- **✅ Phase 10 Complete**: Recipe/formula system (all 12 sub-phases)
- **✅ Phase 11 Complete**: Production deployment to DigitalOcean (https://scentora.thejoeshow.net)
- **✅ Phase 12 Complete**: Upgraded to Koa.js 3.1.1 and latest ecosystem packages
- **🔄 Phase 13 In Progress**: CLI Console (interactive REPL for administration)

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

1. **Integration Tests** (`backend/tests/*.test.ts`)
   - Test full request flow via Supertest
   - Use test PostgreSQL database (port 5435 dev, 5432 CI)
   - DB cleaned between tests in FK-safe order
   - 70 tests across 7 test files
   - Examples: Auth flows, CRUD operations, recipe system

#### 🏃 Running Tests

**Backend (Koa.js/TypeScript/Vitest)**:
```bash
# Start test database
docker compose up -d

# Run all tests
cd backend && npm test

# Run specific test file
npm test -- tests/auth.test.ts
```

**Frontend (Vue/Vitest)**:
```bash
# Run all tests
cd frontend && npm test
```

**CI**: Tests run automatically via GitHub Actions on PRs and pushes to main.
Rate limiting is disabled in test environment (`NODE_ENV=test`) to prevent 429 errors.

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

### Phase 7: Migration to Go Backend (Completed then Superseded - Jan 30, 2026)
**Goal**: Rewrite backend in Go for better performance and maintainability

**Note**: The Go backend was later replaced with Koa.js/TypeScript/Prisma during Phase 9 for better developer experience and ecosystem alignment with the TypeScript frontend.

**Deliverables** (historical):
- ✅ Go/Echo backend with PostgreSQL
- ✅ 100% API compatibility with Node.js version
- ✅ All features ported (auth, invitations, perfumes, journals)
- ✅ Database migrations (auto-run on startup)

---

## Current State

### Architecture (As of Feb 7, 2026)

**Backend**: `backend/` (Koa.js + TypeScript)
- Framework: Koa.js 3.1.1 + TypeScript + Prisma + Zod
- Database: PostgreSQL 15
- Port: 3000
- Architecture: Routes → Services → Prisma (no repository layer)
- ESM project with `.js` import extensions

**Frontend**: `frontend/` (Vue.js)
- Framework: Vue 3.5 + TypeScript + Naive UI + Tailwind CSS v4
- Build: Vite 6.x
- Port: 5173 (dev), 8080 (prod nginx container)

**Database**: PostgreSQL
- Dev: localhost:5435 (Docker)
- Prod: Docker container on scentora-network
- Database: scentora
- Schema: 14 Prisma models with camelCase fields + `@map("snake_case")`

**Production**: https://scentora.thejoeshow.net
- DigitalOcean droplet (Ubuntu 24.04)
- 3 Docker containers: frontend, backend, postgres
- Host nginx reverse proxy + Let's Encrypt SSL
- GitHub Actions CI/CD (test on PR, deploy on push to main)

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

### Phase 12: Koa.js 3.1.1 Upgrade (✅ Complete)
**Goal**: Upgrade Koa.js and ecosystem packages to latest stable versions
**Status**: ✅ Complete
**Completed**: February 10, 2026

**Objective**: Update backend framework dependencies to modern versions for improved performance, security, and compatibility with latest Node.js features.

**Package Upgrades Completed**:
- `koa`: 2.16.3 → 3.1.1 ✅
- `@koa/router`: 13.1.1 → 15.3.0 ✅
- `@types/koa`: 2.15.0 → 3.0.1 ✅
- `@types/koa__router`: 12.0.4 → 12.0.5 ✅

**Results**:
- ✅ Zero code changes required (already Koa 3.x compatible)
- ✅ TypeScript compilation successful
- ✅ All endpoints tested and functional
- ✅ Production build successful
- ✅ 30-minute completion time

**Documentation**: `docs/phases/PHASE12_COMPLETE.md`

---

### Future Work

### Phase 13: CLI Console (In Progress)
**Goal**: Interactive command-line interface for administration and data management
**Status**: 📝 Planning / Documentation Phase
**Started**: February 10, 2026

**Objective**: Create an interactive REPL-style CLI console for managing users, accords, recipes, and performing database operations without the web UI.

#### Technology Stack

- **Framework**: Vorpal.js - Interactive CLI/REPL framework
- **Runtime**: tsx - TypeScript execution matching backend stack
- **UI Libraries**:
  - chalk - Terminal colors and styling
  - cli-table3 - Formatted ASCII tables
  - ora - Loading spinners and progress indicators
  - inquirer - Interactive prompts (integrated with Vorpal)

#### Architecture

```
backend/cli/
├── index.ts              # Entry point, bootstrap REPL
├── context.ts            # Application context (Prisma, services, config)
├── commands/
│   ├── user.ts          # User management commands
│   ├── accord.ts        # Accord commands (future)
│   ├── recipe.ts        # Recipe commands (future)
│   └── db.ts            # Database operations (future)
├── utils/
│   ├── output.ts        # Formatted output helpers
│   ├── validation.ts    # Input validation
│   └── prompts.ts       # Reusable prompt definitions
└── types/
    └── context.ts       # Context type definitions
```

#### Phase 13.1: User Management Commands (Initial Release)

**Commands to Implement**:
1. `create-user` - Create new user with interactive prompts
2. `list-users` - Display all users in formatted table
3. `delete-user <email|id>` - Delete user with confirmation
4. `reset-password <email|id>` - Reset user password
5. `show-user <email|id>` - Display detailed user information

**Features**:
- ✅ Auto-complete for commands
- ✅ Command history (persisted across sessions)
- ✅ Colored output (success/error/warning/info)
- ✅ Interactive prompts with validation
- ✅ Password input hiding
- ✅ Confirmation dialogs for destructive actions
- ✅ Formatted tables for list views
- ✅ Loading spinners for operations

**Usage Example**:
```bash
cd backend
npm run console

🌸 Scentora CLI Console (v1.0.0)
Environment: development
Database: localhost:5432/scentora
Connected: ✅

scentora> create-user
? Email: admin@example.com
? Password: ********
? Username: admin
✅ User created successfully!
```

#### Phase 13.2: Accord Management (Future)
- `import-accords <file>` - Import accords from JSON
- `export-accords <email|id> [file]` - Export user's accords
- `list-accords <email|id>` - List all accords for user

#### Phase 13.3: Recipe Management (Future)
- `import-recipes <file>` - Import recipes from JSON
- `export-recipes <email|id> [file]` - Export user's recipes
- `clone-recipe <id> <target-user>` - Clone recipe to another user

#### Phase 13.4: Database Operations (Future)
- `db:migrate` - Run pending migrations
- `db:seed` - Seed test data
- `db:backup [file]` - Create database backup
- `db:restore <file>` - Restore from backup

#### Implementation Tasks (Phase 13.1)

- [ ] Setup Vorpal.js with TypeScript
- [ ] Create CLI entry point and context loader
- [ ] Install dependencies (chalk, cli-table3, ora, inquirer)
- [ ] Implement `create-user` command with validation
- [ ] Implement `list-users` command with table formatting
- [ ] Implement `delete-user` command with confirmation
- [ ] Implement `reset-password` command
- [ ] Implement `show-user` command
- [ ] Add output utilities (colors, tables, spinners)
- [ ] Write tests for all commands (80%+ coverage)
- [ ] Update package.json with `console` script
- [ ] Document CLI in `docs/CLI_CONSOLE.md`

#### Success Criteria (Phase 13.1)

- ✅ CLI starts without errors
- ✅ All 5 user commands functional
- ✅ Auto-complete works for commands
- ✅ Command history persists across sessions
- ✅ Colored output renders correctly
- ✅ Tables display properly formatted
- ✅ Password prompts hide input
- ✅ Confirmation prompts work correctly
- ✅ Error messages are clear and helpful
- ✅ Tests pass for all commands (80%+ coverage)
- ✅ Documentation complete and accurate

#### Documentation

- **Spec**: `specs/cli-console.md` - Detailed technical specification
- **User Guide**: `docs/CLI_CONSOLE.md` - User-facing documentation
- **API Examples**: Included in spec for all commands

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

### Completed Phase (9): Backend Rewrite & Testing ✅
**Status**: Complete
**Completion Date**: February 2026
**Achievement**: Backend rewritten from Go/Echo to Koa.js/TypeScript/Prisma; 70 integration tests passing

**Key Changes**:
- Rewrote Go/Echo/sqlx backend to Koa.js/TypeScript/Prisma
- ESM project with `.js` import extensions
- Routes → Services → Prisma architecture (no repository layer)
- Prisma schema with camelCase fields + `@map("snake_case")`
- 70 integration tests via Vitest + Supertest across 7 test files
- Test config supports CI override via `process.env` fallbacks

### Completed Phase (10): Recipe/Formula System ✅
**Status**: Complete
**Completion Date**: February 2026

**Overview**:
Recipe/formula system enabling users to create perfume recipes by combining accords with version control, notes, collections, and tagging.

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

### Completed Phase (11): Production Deployment to DigitalOcean ✅
**Status**: Complete
**Completion Date**: February 7, 2026
**Live URL**: https://scentora.thejoeshow.net

**Overview**:
Automated deployment pipeline to DigitalOcean droplet with Ubuntu 24.04, GitHub Actions CI/CD, and production configuration.

**Deliverables**:
- ✅ `deploy/setup-droplet.sh` - Automated droplet configuration script
- ✅ `.github/workflows/deploy.yml` - GitHub Actions deployment pipeline
- ✅ `.github/workflows/test.yml` - Automated testing workflow (reusable via `workflow_call`)
- ✅ `docker-compose.prod.yml` - Production Docker Compose (3 containers)
- ✅ `docs/DEPLOYMENT.md` - Comprehensive deployment guide

**Infrastructure**:
- DigitalOcean droplet (Ubuntu 24.04)
- 3 Docker containers: frontend (nginx:alpine), backend (Node.js), postgres (15-alpine)
- Host nginx reverse proxy with Let's Encrypt SSL (auto-renewal via certbot timer)
- UFW firewall (ports 22, 80, 443)
- Daily database backups at 2 AM (30-day retention)
- Automated security updates via unattended-upgrades

**CI/CD Pipeline**:
- `test.yml`: Runs on PRs + callable via `workflow_call`
- `deploy.yml`: Runs on push to main, calls test.yml first, then SSHs to droplet to rebuild and deploy
- Health check with retry loop (60s initial wait + 5 retries)

**Key Debugging Lessons**:
- Container health checks must use `127.0.0.1` not `localhost` (IPv6 resolution fails in Alpine)
- Nginx configs in shell heredocs need quoted heredoc (`<<'EOF'`) + sed for domain substitution
- Docker Compose healthcheck overrides Dockerfile HEALTHCHECK
- Rate limiting must be disabled in test env to prevent CI 429 errors

### Future Phases (12+)

**Phase 12: Advanced Recipe Features** (Planned)
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
- Clean, maintainable full-stack TypeScript codebase
- Production backend deployed (Koa.js + Prisma)
- Modern, responsive frontend (Vue 3 + Naive UI)
- Comprehensive documentation
- Secure authentication with refresh token rotation
- Automated CI/CD pipeline with GitHub Actions
- Production deployment at https://scentora.thejoeshow.net
- 70 integration tests passing in CI
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
- **Koa.js + TypeScript Backend**: Full-stack TypeScript for consistency, Prisma ORM for type-safe DB access
- **PostgreSQL**: Relational model for accord/recipe relationships, generated columns for computed fields
- **JWT Auth**: Stateless, scalable authentication with refresh token rotation
- **Vue 3 Composition API**: Modern, reactive, maintainable
- **Naive UI + Tailwind CSS**: Notion-inspired minimalist design
- **Tag System**: Flexible categorization without rigid constraints
- **Docker + GitHub Actions**: Containerized deployment with automated CI/CD

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
