# Backend Parity Plan: Add Node.js Features to Go Backend

## Problem Statement
The Go backend (PostgreSQL) is missing several features that exist in the Node.js backend (CouchDB). This plan identifies the missing features and provides a checklist to achieve complete parity.

## Approach
Add missing features to the Go backend to match all functionality present in the Node.js/Koa backend, ensuring the Go backend can serve as a complete drop-in replacement.

---

## Feature Comparison & Parity Checklist

### ✅ Already Implemented (Verified)
- [x] Authentication system (register, login, refresh, logout, me)
- [x] Invitation system (create, list, revoke)
- [x] Perfume CRUD operations
- [x] Journal entry CRUD operations
- [x] Notes aggregation endpoint
- [x] Statistics endpoint
- [x] Export collection endpoint (GET /api/export)
- [x] JWT authentication middleware
- [x] Request validation
- [x] Error handling
- [x] CORS configuration

### ❌ Missing Features in Go Backend

#### 1. Auth Endpoints
- [ ] **POST /api/auth/logout-all** - Revoke all refresh tokens for a user (logout from all devices)
  - Node.js location: `backend/src/controllers/authController.ts:238`
  - Node.js route: `backend/src/routes/auth.ts:19`
  - Requires: Protected endpoint with authentication

#### 2. Data Management
- [ ] **POST /api/export/import** - Import collection data
  - Node.js location: `backend/src/controllers/exportController.ts:33`
  - Node.js route: `backend/src/routes/export.ts:12`
  - Functionality: Import perfumes and journal entries from JSON

#### 3. Middleware Features
- [ ] **Optional Authentication Middleware** - Allow endpoints to work with or without auth
  - Node.js location: `backend/src/middleware/auth.ts:41`
  - Use case: Future features that might need optional user context

#### 4. Rate Limiting
- [ ] **Rate Limiting Middleware Implementation**
  - Node.js: Uses koa-ratelimit with in-memory store
  - Auth endpoints: 5 requests per 15 minutes
  - General endpoints: 100 requests per minute
  - Config exists in Go (`config.go:34-35`) but not implemented
  - Need to implement actual middleware

#### 5. Export Format Consistency
- [ ] **Export data format alignment**
  - Node.js includes: `version`, `exportDate`, `perfumes`, `journalEntries`
  - Go only includes: `perfumes`, `journal`
  - Make Go format match Node.js exactly

---

## Implementation Work Plan

### Phase 1: Auth Enhancements
- [ ] Add `LogoutAll` handler to revoke all user tokens
  - [ ] Add method in `services.AuthService`
  - [ ] Add handler in `handlers.AuthHandler`
  - [ ] Add route `POST /api/auth/logout-all` (protected)
  - [ ] Test logout-all functionality

### Phase 2: Import Collection Feature
- [ ] Implement collection import functionality
  - [ ] Add `ImportCollection` method to `handlers.ExportHandler`
  - [ ] Create validation for import data structure
  - [ ] Handle perfumes import with UUID remapping
  - [ ] Handle journal entries import with perfume ID mapping
  - [ ] Return import results (counts and errors)
  - [ ] Add route `POST /api/export/import`
  - [ ] Test import with various data scenarios

### Phase 3: Export Format Alignment
- [ ] Update export response to match Node.js format
  - [ ] Add `version` field (e.g., "1.0")
  - [ ] Add `exportDate` field (ISO 8601 timestamp)
  - [ ] Rename `journal` to `journalEntries`
  - [ ] Test export format matches Node.js exactly

### Phase 4: Rate Limiting Implementation
- [ ] Add rate limiting middleware
  - [ ] Research Go rate limiting libraries (e.g., `golang.org/x/time/rate`, `tollbooth-echo`)
  - [ ] Implement auth rate limiter (5 req/15min)
  - [ ] Implement general rate limiter (100 req/1min)
  - [ ] Apply to auth routes
  - [ ] Test rate limit enforcement

### Phase 5: Optional Auth Middleware
- [ ] Create optional authentication middleware
  - [ ] Add `OptionalJWTAuth` middleware function
  - [ ] Set user context if valid token present
  - [ ] Continue without error if no token or invalid token
  - [ ] Document usage patterns

### Phase 6: Testing & Validation
- [ ] Create comprehensive test suite
  - [ ] Test logout-all with multiple tokens
  - [ ] Test import with various data formats
  - [ ] Test export format matches Node.js
  - [ ] Test rate limiting triggers correctly
  - [ ] Test optional auth middleware behavior
- [ ] Compare API responses between Node.js and Go backends
- [ ] Document any remaining differences

### Phase 7: Documentation
- [ ] Update IMPLEMENTATION_COMPLETE.md with new features
- [ ] Update API documentation
- [ ] Document database differences (CouchDB vs PostgreSQL)
- [ ] Create migration guide for switching backends

---

## Key Considerations

### Database Differences
- **Node.js**: CouchDB (document store, no schema)
- **Go**: PostgreSQL (relational, with schema)
- Import/export must handle ID mapping due to different ID generation strategies

### Error Response Format
Both backends should return errors in this format:
```json
{
  "error": {
    "message": "Error message",
    "details": "Optional details"
  }
}
```

### Authentication Compatibility
- JWT structure must be identical
- Token expiration times must match
- Refresh token rotation must work the same way

### API Route Compatibility
All routes must match exactly:
- `/api/auth/*`
- `/api/perfumes/*`
- `/api/journal/*`
- `/api/invitations/*`
- `/api/notes`
- `/api/stats`
- `/api/export/*`

---

## Success Criteria

✅ Go backend passes all tests that Node.js backend passes
✅ Frontend can switch between backends without code changes
✅ All API endpoints exist in both backends
✅ Response formats match exactly
✅ Error handling behaves identically
✅ Rate limiting protects both backends equally

---

## Notes

- The Go backend will continue using PostgreSQL (not switching to CouchDB)
- Feature parity means functional equivalence, not implementation equivalence
- Some internal differences are acceptable as long as the API contract is identical
- Performance characteristics may differ between backends (that's okay)
