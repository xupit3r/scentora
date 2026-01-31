# Testing & Quality Assurance Plan - Phase 9

**Date:** January 31, 2026  
**Status:** In Planning  
**Priority:** HIGH - Critical for Production Readiness

---

## Overview

Implement comprehensive testing coverage for the Scentora application to ensure reliability, prevent regressions, and enable confident deployment to production.

## Current Testing Status

### Backend (Go)
- **Existing:** 5 test files
  - `database_test.go` - Database connection tests
  - `user_repo_test.go` - User repository tests (8 tests, passing)
  - `invitation_repo_test.go` - Invitation tests
  - `accord_repo_test.go` - Accord repository tests
  - `predefined_tag_repo_test.go` - Tag tests
- **Coverage:** Repository layer partially covered
- **Missing:** Service layer, handler layer, middleware tests
- **Framework:** Built-in Go testing + testify

### Frontend (Vue.js)
- **Existing:** No tests
- **Script:** `"test": "echo \"Error: no test specified\" && exit 1"`
- **Missing:** Unit tests, component tests, E2E tests
- **Framework:** None installed (need Vitest)

## Testing Strategy

### Priorities (In Order)

1. **Backend Unit Tests** (Critical - 1 week)
   - Repository layer (complete existing)
   - Service layer (NEW)
   - Utility functions (NEW)

2. **Backend Integration Tests** (High - 3 days)
   - Handler/Controller tests
   - Middleware tests
   - End-to-end API tests

3. **Frontend Unit Tests** (High - 1 week)
   - Component tests
   - Store tests (Pinia)
   - Composable tests
   - Utility function tests

4. **E2E Tests** (Medium - 3 days)
   - Critical user flows
   - Authentication flow
   - CRUD operations

5. **Security & Load Testing** (Medium - 2 days)
   - Security audit
   - Load testing
   - Performance benchmarks

## Testing Goals

### Coverage Targets
- **Backend:** 80%+ code coverage
- **Frontend:** 70%+ code coverage
- **E2E:** 100% of critical paths covered

### Success Criteria
- [ ] All existing tests pass consistently
- [ ] Repository layer: 100% coverage
- [ ] Service layer: 90%+ coverage
- [ ] Handler layer: 85%+ coverage
- [ ] Middleware: 100% coverage
- [ ] Frontend components: 70%+ coverage
- [ ] Frontend stores: 90%+ coverage
- [ ] E2E: Auth + CRUD flows covered
- [ ] CI/CD integration
- [ ] Test documentation complete

---

## Phase 9 Implementation Plan

### Phase 9.1: Backend Repository Tests (2-3 days)

**Goal:** Complete and fix all repository tests

**Tasks:**
- [ ] Fix failing repository tests
- [ ] Complete accord repository tests
- [ ] Complete invitation repository tests
- [ ] Complete predefined tag repository tests
- [ ] Add refresh token repository tests
- [ ] Test all edge cases and error conditions
- [ ] Achieve 100% repository coverage

**Files:**
- `internal/repository/accord_repo_test.go`
- `internal/repository/invitation_repo_test.go`
- `internal/repository/predefined_tag_repo_test.go`
- `internal/repository/refresh_token_repo_test.go` (NEW)

### Phase 9.2: Backend Service Tests (3-4 days)

**Goal:** Test business logic layer

**Tasks:**
- [ ] Create `accord_service_test.go`
  - CreateAccord validation
  - UpdateAccord logic
  - DeleteAccord cascading
  - ListAccords filtering
  - SearchAccords functionality
  - GetStatistics calculations
- [ ] Create `auth_service_test.go`
  - Registration flow
  - Login flow
  - Token generation
  - Token refresh
  - Password hashing
- [ ] Create `invitation_service_test.go`
- [ ] Create `tag_service_test.go`
- [ ] Mock repository dependencies

**Files:**
- `internal/services/accord_service_test.go` (NEW)
- `internal/services/auth_service_test.go` (NEW)
- `internal/services/invitation_service_test.go` (NEW)
- `internal/services/tag_service_test.go` (NEW)

### Phase 9.3: Backend Handler/Integration Tests (2-3 days)

**Goal:** Test HTTP handlers and middleware

**Tasks:**
- [ ] Create `accord_handler_test.go`
  - POST /api/accords
  - GET /api/accords
  - GET /api/accords/:id
  - PUT /api/accords/:id
  - DELETE /api/accords/:id
  - Tag operations
- [ ] Create `auth_handler_test.go`
  - Registration endpoint
  - Login endpoint
  - Refresh token endpoint
  - Logout endpoints
- [ ] Create `stats_handler_test.go`
- [ ] Create `export_handler_test.go`
- [ ] Test middleware
  - JWT authentication
  - Rate limiting
  - CORS
  - Error handling

**Files:**
- `internal/handlers/accord_test.go` (NEW)
- `internal/handlers/auth_test.go` (NEW)
- `internal/handlers/stats_test.go` (NEW)
- `internal/handlers/export_test.go` (NEW)
- `internal/middleware/auth_test.go` (NEW)
- `internal/middleware/rate_limit_test.go` (NEW)

### Phase 9.4: Frontend Testing Setup (1 day)

**Goal:** Set up Vitest testing framework

**Tasks:**
- [ ] Install Vitest and dependencies
- [ ] Install @vue/test-utils
- [ ] Install happy-dom or jsdom
- [ ] Configure vitest.config.ts
- [ ] Set up test utilities
- [ ] Create test helpers
- [ ] Update package.json scripts

**Dependencies to Install:**
```bash
npm install -D vitest @vue/test-utils happy-dom
npm install -D @vitest/ui @vitest/coverage-v8
```

**Files:**
- `vitest.config.ts` (NEW)
- `src/test-utils/` (NEW directory)

### Phase 9.5: Frontend Component Tests (3-4 days)

**Goal:** Test Vue components

**Tasks:**
- [ ] Test `AccordCard.vue`
  - Renders correct data
  - Position badge displays
  - Volume formatting
  - Low stock warnings
  - Action button clicks
- [ ] Test `AccordForm.vue`
  - Form validation
  - Submit handling
  - Edit mode
  - Tag selection
- [ ] Test `AccordFilters.vue`
  - Filter changes emit correctly
  - Clear filters works
  - Mobile mode
- [ ] Test `TagSelector.vue`
  - Autocomplete functionality
  - Tag selection
  - Custom tag creation
- [ ] Test `AppSidebar.vue`
- [ ] Test `SkeletonCard.vue`
- [ ] Test `EmptyState.vue`

**Files:**
- `src/components/AccordCard.test.ts` (NEW)
- `src/components/AccordForm.test.ts` (NEW)
- `src/components/AccordFilters.test.ts` (NEW)
- `src/components/TagSelector.test.ts` (NEW)
- `src/components/layout/AppSidebar.test.ts` (NEW)
- `src/components/ui/SkeletonCard.test.ts` (NEW)
- `src/components/ui/EmptyState.test.ts` (NEW)

### Phase 9.6: Frontend Store Tests (1-2 days)

**Goal:** Test Pinia stores

**Tasks:**
- [ ] Test `auth.store.ts`
  - Login action
  - Logout action
  - Register action
  - Token refresh
  - State persistence
- [ ] Mock API calls
- [ ] Test error handling

**Files:**
- `src/stores/auth.test.ts` (NEW)

### Phase 9.7: Frontend Utility Tests (1 day)

**Goal:** Test utility functions

**Tasks:**
- [ ] Test `useKeyboard.ts` composable
- [ ] Test `useConfirmDialog.ts` composable
- [ ] Test `volume.ts` utilities
  - mlToDrops conversion
  - dropsToMl conversion
  - formatVolume
  - getStockLevel
  - getStockWarning

**Files:**
- `src/composables/useKeyboard.test.ts` (NEW)
- `src/composables/useConfirmDialog.test.ts` (NEW)
- `src/utils/volume.test.ts` (NEW)

### Phase 9.8: E2E Tests (2-3 days)

**Goal:** Test critical user journeys

**Tasks:**
- [ ] Install Playwright or Cypress
- [ ] Set up E2E test environment
- [ ] Write authentication tests
  - Login flow
  - Logout flow
  - Session persistence
- [ ] Write accord CRUD tests
  - Create accord
  - Edit accord
  - Delete accord (with confirmation)
  - Search and filter
- [ ] Write statistics tests
  - View statistics
  - Export data

**Framework Choice:** Playwright (recommended)

**Files:**
- `e2e/auth.spec.ts` (NEW)
- `e2e/accords.spec.ts` (NEW)
- `e2e/statistics.spec.ts` (NEW)
- `playwright.config.ts` (NEW)

### Phase 9.9: CI/CD Integration (1 day)

**Goal:** Automate testing in CI pipeline

**Tasks:**
- [ ] Create GitHub Actions workflow
- [ ] Run backend tests on push
- [ ] Run frontend tests on push
- [ ] Run E2E tests on PR
- [ ] Code coverage reporting
- [ ] Status badges in README

**Files:**
- `.github/workflows/test.yml` (NEW)
- `.github/workflows/e2e.yml` (NEW)

### Phase 9.10: Documentation & Polish (1 day)

**Goal:** Document testing approach

**Tasks:**
- [ ] Write testing guide
- [ ] Document how to run tests
- [ ] Document how to write tests
- [ ] Add coverage reports
- [ ] Update README with test status

**Files:**
- `docs/TESTING_GUIDE.md` (NEW)
- `docs/phases/PHASE9_COMPLETE.md` (NEW)

---

## Testing Tools & Frameworks

### Backend (Go)
- **Framework:** Built-in `testing` package
- **Assertions:** `github.com/stretchr/testify`
- **Mocking:** `github.com/stretchr/testify/mock`
- **Coverage:** `go test -cover`
- **Database:** Test database with migrations

### Frontend (Vue.js)
- **Framework:** Vitest
- **Component Testing:** @vue/test-utils
- **DOM:** happy-dom (faster than jsdom)
- **Coverage:** @vitest/coverage-v8
- **UI:** @vitest/ui (optional test viewer)

### E2E Testing
- **Framework:** Playwright
- **Browser:** Chromium (primary), Firefox, WebKit
- **Features:** Auto-wait, screenshots, video recording

### CI/CD
- **Platform:** GitHub Actions
- **Triggers:** Push, Pull Request
- **Coverage:** Codecov or Coveralls

---

## Test Execution Commands

### Backend
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/services

# Run specific test
go test -run TestAccordService_Create ./internal/services

# Verbose output
go test -v ./...
```

### Frontend
```bash
# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch

# Run UI mode
npm run test:ui

# Run specific test
npm test AccordCard
```

### E2E
```bash
# Run E2E tests
npm run test:e2e

# Run in headed mode
npm run test:e2e:headed

# Run specific spec
npm run test:e2e -- auth.spec.ts
```

---

## Success Metrics

### Code Coverage
- Backend repositories: 100%
- Backend services: 90%+
- Backend handlers: 85%+
- Frontend components: 70%+
- Frontend stores: 90%+
- Overall backend: 85%+
- Overall frontend: 70%+

### Test Quality
- All tests pass consistently
- Tests run in < 2 minutes (backend)
- Tests run in < 1 minute (frontend)
- E2E tests run in < 5 minutes
- No flaky tests
- Clear test descriptions

### Documentation
- Testing guide complete
- All test commands documented
- CI/CD setup documented
- Coverage reports generated

---

## Timeline

**Total Estimated Time:** 3-4 weeks

- Week 1: Backend repository & service tests
- Week 2: Backend handlers & frontend setup
- Week 3: Frontend component & store tests
- Week 4: E2E tests, CI/CD, documentation

**Can be accelerated to 2 weeks with focused effort.**

---

## Risks & Mitigation

### Risks
1. **Test Database Setup:** Complex setup for integration tests
   - Mitigation: Use Docker for test database, clear seed data

2. **Frontend Async Testing:** Vue component async behavior
   - Mitigation: Use proper async utilities from @vue/test-utils

3. **E2E Test Flakiness:** Network timing issues
   - Mitigation: Use Playwright's auto-wait, proper selectors

4. **Time Overrun:** Testing takes longer than estimated
   - Mitigation: Prioritize critical paths first, incremental coverage

---

## Next Steps

1. Review and approve this testing plan
2. Start with Phase 9.1 (Backend Repository Tests)
3. Fix existing failing tests
4. Implement new tests incrementally
5. Maintain >80% coverage throughout

**Ready to begin implementation!** 🚀
