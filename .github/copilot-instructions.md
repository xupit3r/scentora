# Scentora - GitHub Copilot Instructions

## Project Overview

Scentora is a perfume cataloging application that allows users to track scent profiles, notes, and personal journal entries. It's built with a modern full-stack architecture emphasizing security and user data isolation.

## Architecture

### Tech Stack
- **Backend**: Koa.js (Node.js framework) + TypeScript + CouchDB
- **Frontend**: Vue.js 3 + TypeScript + Pinia (state management)
- **Database**: CouchDB (NoSQL document database)
- **Authentication**: JWT with refresh tokens

### Project Structure
```
scentora/
├── backend/                  # Koa.js REST API
│   ├── src/
│   │   ├── config/          # Database and JWT configuration
│   │   ├── controllers/     # Request handlers
│   │   ├── middleware/      # Auth, error handling, rate limiting
│   │   ├── models/          # Data models (User, Perfume, etc.)
│   │   └── routes/          # API route definitions
│   ├── package.json
│   └── tsconfig.json
├── frontend/                 # Vue.js SPA
│   ├── src/
│   │   ├── components/      # Reusable Vue components
│   │   ├── views/           # Page components
│   │   ├── stores/          # Pinia stores
│   │   ├── services/        # API service layer
│   │   ├── router/          # Vue Router configuration
│   │   └── types/           # TypeScript types
│   ├── package.json
│   └── vite.config.ts
└── docker-compose.yml        # CouchDB setup
```

## Development Setup

### Prerequisites
- Node.js 20+
- Docker & Docker Compose (for CouchDB)
- npm

### Quick Start
```bash
# Automated (recommended)
npm start

# Or manually
docker-compose up -d        # Start CouchDB
cd backend && npm install && npm run dev
cd frontend && npm install && npm run dev
```

### Backend Development
- **Port**: 3000
- **Dev Command**: `npm run dev` (uses ts-node-dev with hot reload)
- **Build Command**: `npm run build` (compiles TypeScript to dist/)
- **Start Command**: `npm start` (runs compiled dist/index.js)

### Frontend Development
- **Port**: 5173
- **Dev Command**: `npm run dev` (Vite dev server with HMR)
- **Build Command**: `npm run build` (Vue TSC + Vite build)
- **Preview Command**: `npm run preview` (preview production build)

### Database
- **CouchDB URL**: http://localhost:5984
- **Admin UI**: http://localhost:5984/_utils
- **Credentials**: admin/password (dev environment)

## Code Style & Conventions

### TypeScript
- **Strict Mode**: Enabled in both frontend and backend
- **Target**: ES2020
- **Module System**: CommonJS (backend), ES Modules (frontend)
- Always define types explicitly; avoid `any` type
- Use interfaces for object shapes and types for unions/primitives

### Backend (Koa.js)
- Use async/await for all asynchronous operations
- Controllers should handle request/response logic
- Models contain business logic and data validation
- Middleware for cross-cutting concerns (auth, error handling, rate limiting)
- Use Zod for runtime validation of request bodies
- Follow RESTful API conventions

### Frontend (Vue.js)
- **Composition API**: Use `<script setup>` syntax
- **State Management**: Pinia stores for global state
- **Routing**: Vue Router with route guards for protected routes
- **API Calls**: Use the axios-based service layer in `src/services/`
- **Type Safety**: Define TypeScript interfaces for all API responses

### File Naming
- **Components**: PascalCase (e.g., `PerfumeCard.vue`)
- **Routes/Controllers**: kebab-case or descriptive names (e.g., `perfume-routes.ts`)
- **Stores**: kebab-case with `.store.ts` suffix (e.g., `auth.store.ts`)
- **Types**: PascalCase for interfaces/types (e.g., `Perfume`, `User`)

### Comments
- Use JSDoc comments for functions and complex logic
- Avoid obvious comments; code should be self-documenting
- Document WHY, not WHAT, when adding comments

## Security Practices

### Authentication & Authorization
- All protected endpoints require valid JWT access token
- Access tokens expire in 15 minutes
- Refresh tokens expire in 7 days
- Rate limiting on auth endpoints (5 requests per 15 minutes)
- Never log or expose tokens or secrets
- All user data is isolated per user (userId in database queries)

### Password Handling
- Use bcrypt for password hashing (10 rounds)
- Never store or log plain-text passwords
- Minimum password requirements in validation

### Environment Variables
- **NEVER** commit `.env` files
- Use `.env.example` as a template
- Critical: `JWT_SECRET` must be changed in production
- Database credentials should be secured in production

### Data Validation
- Validate all user input using Zod schemas
- Sanitize data before storing in database
- Use TypeScript for compile-time type safety

## Testing

### Backend Testing
- Currently: `npm test` returns "no test specified"
- When adding tests: Use Jest or Mocha
- Test structure: Unit tests for models/controllers, integration tests for routes
- Mock CouchDB for unit tests

### Frontend Testing
- Currently: `npm test` returns "no test specified"
- When adding tests: Use Vitest (built-in with Vite)
- Test structure: Unit tests for components/stores, E2E tests with Playwright/Cypress

## API Conventions

### Endpoints
- Base URL: `/api`
- Auth: `/api/auth/*` (public)
- Perfumes: `/api/perfumes/*` (protected)
- Journal: `/api/journal/*` (protected)
- Utility: `/api/notes`, `/api/stats`, `/api/export/*` (protected)

### Response Format
```typescript
// Success
{
  success: true,
  data: { /* response data */ }
}

// Error
{
  success: false,
  error: "Error message"
}
```

### Status Codes
- 200: Success
- 201: Created
- 400: Bad Request (validation errors)
- 401: Unauthorized (missing/invalid token)
- 404: Not Found
- 429: Too Many Requests (rate limit)
- 500: Internal Server Error

## Database Patterns

### CouchDB Documents
- Each document has `_id` and `_rev` (managed by CouchDB)
- User-specific data includes `userId` field for isolation
- Use compound keys for complex queries (e.g., `userId:perfumeId`)

### Models
- User: Authentication and profile data
- Perfume: Perfume details with notes pyramid
- Journal: Journal entries linked to perfumes
- RefreshToken: Token rotation and session management

## Dependencies

### Backend Core
- `koa`: Web framework
- `@koa/router`: Routing
- `koa-bodyparser`: Parse request bodies
- `@koa/cors`: CORS middleware
- `koa-ratelimit`: Rate limiting
- `nano`: CouchDB client
- `jsonwebtoken`: JWT auth
- `bcryptjs`: Password hashing
- `zod`: Schema validation
- `dotenv`: Environment variables

### Frontend Core
- `vue`: UI framework
- `vue-router`: Routing
- `pinia`: State management
- `axios`: HTTP client
- `vite`: Build tool

## Documentation

- `README.md`: Getting started and overview
- `API_REFERENCE.md`: Complete API documentation
- `AUTH_IMPLEMENTATION.md`: Authentication details
- `QUICKSTART.md`: Quick start guide
- `LAUNCHER_GUIDE.md`: Launcher script usage

## Common Tasks

### Adding a New API Endpoint
1. Create controller function in `backend/src/controllers/`
2. Add route in `backend/src/routes/`
3. Add authentication middleware if protected
4. Update specs/api-spec.md

### Adding a New Feature
1. Backend: Create controller, model, and route
2. Frontend: Create service method in `src/services/`
3. Add Pinia store if managing state
4. Create/update Vue components and views
5. Update router if adding new pages
6. Update specs/api-spec.md

### Database Migration
1. CouchDB is schema-less; no migrations needed
2. Add version field to documents for schema evolution
3. Handle multiple document versions in code

## Git Workflow

- Main branch: `main`
- Feature branches: `feature/description`
- Bug fixes: `fix/description`
- Commit messages: Use conventional commits (feat:, fix:, docs:, etc.)

## Deployment

- See `deploy/` directory for deployment guides
- Ensure environment variables are properly set
- Change `JWT_SECRET` in production
- Enable HTTPS
- Configure production-grade rate limiting (consider Redis)
- Review security checklist in docs/AUTH_IMPLEMENTATION.md

## Tips for Copilot

- When creating controllers, follow the pattern in existing controllers
- Authentication middleware (`authMiddleware`) should be applied to protected routes
- Use the existing validation patterns with Zod
- Frontend API calls should use the service layer, not direct axios
- State management should use Pinia stores, not local component state for global data
- Follow the monorepo structure: backend and frontend are separate but coordinated
