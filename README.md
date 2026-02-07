# Scentora

A perfume formulation and accord management system for DIY perfumers and enthusiasts.

## Tech Stack

- **Backend**: Koa.js + TypeScript + Prisma + Zod + JWT Authentication
- **Frontend**: Vue.js 3 + TypeScript + Pinia + Naive UI + Tailwind CSS v4
- **Database**: PostgreSQL
- **Testing**: Vitest + Supertest (backend), Vitest (frontend)
- **Design System**: Notion-inspired minimalist UI

## Features

### Authentication & Security
- JWT-based authentication with refresh tokens
- Automatic token refresh (15-min access, 7-day refresh)
- Rate limiting on auth endpoints (brute force protection)
- Multi-user support with complete data isolation
- Secure password hashing (bcrypt)
- Persistent sessions with logout-all capability
- Token rotation for enhanced security
- Invitation-only registration system

### Accord Management
- Manage perfume accords and essential oils
- Track pyramid position (top, middle, base notes)
- Inventory management (volume tracking in ml and drops)
- Rich tagging system with 57+ predefined tags
- Advanced filtering by position, volume, supplier, tags
- Full-text search across names, notes, and metadata
- Low stock warnings and inventory alerts

### Recipe System
- Create and manage perfume recipes with versioning
- Recipe versions with ingredient tracking
- Ingredient volume validation against accord inventory
- Recipe notes (general, testing, observation, adjustment, reminder)
- Recipe tagging and search
- Recipe collections for organization
- Version duplication for iterating on formulas

### Analytics & Export
- Collection statistics dashboard
- Pyramid position distribution analysis
- Tag and supplier usage statistics
- Volume analytics (min/max/average)
- Low inventory alerts (< 10ml threshold)
- Export collection as JSON
- Import accords from JSON backup

### UI/UX
- Notion-inspired clean, minimalist interface
- Fully responsive design (desktop, tablet, mobile)
- Keyboard shortcuts
- Collapsible sidebar navigation
- Skeleton loading states with shimmer animation
- Toast notifications for user feedback

## Getting Started

### Prerequisites

- Node.js 18+
- Docker & Docker Compose (for PostgreSQL)
- npm

### Quick Start

1. **Start PostgreSQL**:
   ```bash
   docker compose up -d
   ```

2. **Backend Setup**:
   ```bash
   cd backend
   npm install
   cp .env.example .env
   # Edit .env to set JWT_SECRET for production
   npx prisma generate
   npx prisma db push
   npm run dev
   ```
   Backend runs at http://localhost:3000

3. **Frontend Setup** (in a new terminal):
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
   Frontend runs at http://localhost:5173

4. **Access the app**:
   - Open http://localhost:5173
   - You'll need an invitation code to register
   - See [Creating Invitations](#creating-invitations) below

## Environment Variables

Create `backend/.env` from the example:
```env
PORT=3000
NODE_ENV=development
DATABASE_URL=postgres://admin:password@localhost:5435/scentora?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_EXPIRES_IN=15m
JWT_REFRESH_EXPIRES_IN=7d
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

## Creating Invitations

Scentora uses an invitation-only registration system.

### Method 1: Direct Database (First User)

```bash
# Connect to PostgreSQL
docker exec -it scentora-postgres psql -U admin -d scentora

# Create a first user manually (password: 'password')
INSERT INTO users (email, username, password_hash)
VALUES ('admin@example.com', 'admin', '$2a$10$your-bcrypt-hash');

# Create an invitation
INSERT INTO invitations (code, created_by, expires_at)
VALUES ('my-invite-code', (SELECT id FROM users LIMIT 1), NOW() + INTERVAL '7 days');
```

### Method 2: API (Authenticated Users)

```bash
curl -X POST "http://localhost:3000/api/invitations" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"email": "friend@example.com", "expiresInDays": 7}'
```

## Project Structure

```
scentora/
├── backend/                # Koa.js + TypeScript API
│   ├── prisma/
│   │   └── schema.prisma   # Database schema (14 models)
│   ├── src/
│   │   ├── middleware/      # Auth, error handling, rate limiting
│   │   ├── routes/          # API route handlers
│   │   ├── services/        # Business logic
│   │   ├── schemas/         # Zod validation schemas
│   │   ├── utils/           # Errors, response helpers
│   │   ├── types/           # TypeScript types
│   │   ├── app.ts           # Koa app setup
│   │   └── index.ts         # Entry point
│   └── tests/               # Vitest integration tests
├── frontend/                # Vue.js 3 SPA
│   ├── src/
│   │   ├── components/      # Vue components
│   │   ├── views/           # Page views
│   │   ├── stores/          # Pinia stores
│   │   ├── services/        # API client
│   │   └── types/           # TypeScript interfaces
│   └── ...
├── specs/                   # API and design specifications
├── deploy/                  # Deployment scripts (setup-droplet.sh)
├── .github/workflows/       # CI/CD (test.yml, deploy.yml)
├── docker-compose.yml       # Dev PostgreSQL container
├── docker-compose.prod.yml  # Production Docker Compose
└── README.md
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register (requires invitation code)
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh tokens
- `POST /api/auth/logout` - Logout (revoke refresh token)
- `GET /api/auth/me` - Get current user
- `POST /api/auth/logout-all` - Revoke all sessions

### Invitations (Protected)
- `POST /api/invitations` - Create invitation
- `GET /api/invitations` - List invitations
- `DELETE /api/invitations/:code` - Revoke invitation

### Accords (Protected)
- `POST /api/accords` - Create accord
- `GET /api/accords` - List accords (with filters)
- `GET /api/accords/:id` - Get accord
- `PUT /api/accords/:id` - Update accord
- `DELETE /api/accords/:id` - Delete accord
- `POST /api/accords/:id/tags` - Add tag
- `DELETE /api/accords/:id/tags/:tag` - Remove tag

### Tags (Public)
- `GET /api/tags` - List all predefined tags
- `GET /api/tags/search?q=...` - Search tags
- `GET /api/tags/categories` - List categories
- `GET /api/tags/grouped` - Tags grouped by category
- `GET /api/tags/category/:category` - Tags by category

### Statistics (Protected)
- `GET /api/stats` - Collection statistics

### Export/Import (Protected)
- `GET /api/export` - Export accords as JSON
- `POST /api/export/import` - Import accords from JSON

### Recipes (Protected) - 21 endpoints
- `POST /api/recipes` - Create recipe
- `GET /api/recipes` - List recipes
- `GET /api/recipes/search?q=...` - Search recipes
- `GET /api/recipes/:id` - Get recipe
- `PUT /api/recipes/:id` - Update recipe
- `DELETE /api/recipes/:id` - Delete recipe
- `POST /api/recipes/:id/versions` - Create version
- `GET /api/recipes/:id/versions` - List versions
- `GET /api/recipes/:id/versions/:versionId` - Get version
- `POST /api/recipes/:id/versions/:versionNumber/activate` - Activate version
- `POST /api/recipes/:id/versions/:versionId/duplicate` - Duplicate version
- `POST /api/recipes/:id/versions/:versionId/ingredients` - Add ingredient
- `PUT /api/recipes/:id/versions/:versionId/ingredients/:ingredientId` - Update ingredient
- `DELETE /api/recipes/:id/versions/:versionId/ingredients/:ingredientId` - Remove ingredient
- `POST /api/recipes/:id/notes` - Create note
- `GET /api/recipes/:id/notes` - List notes
- `PUT /api/recipes/:id/notes/:noteId` - Update note
- `DELETE /api/recipes/:id/notes/:noteId` - Delete note
- `POST /api/recipes/:id/tags` - Add tag
- `DELETE /api/recipes/:id/tags/:tag` - Remove tag
- `GET /api/recipes/tags/popular` - Popular tags

### Collections (Protected)
- `POST /api/collections` - Create collection
- `GET /api/collections` - List collections
- `GET /api/collections/:id` - Get collection
- `PUT /api/collections/:id` - Update collection
- `DELETE /api/collections/:id` - Delete collection
- `POST /api/collections/:id/recipes` - Add recipe to collection
- `DELETE /api/collections/:id/recipes/:recipeId` - Remove recipe from collection

### Health
- `GET /api/health` - Health check

## Testing

```bash
# Start PostgreSQL first
docker compose up -d

# Run backend tests
cd backend
npm test
```

70 integration tests across auth, accords, invitations, tags, stats, export/import, recipes, and collections.

## Security

- JWT refresh tokens with automatic rotation
- Short-lived access tokens (15 minutes)
- Rate limiting on auth endpoints (5 requests/15 min)
- Bcrypt password hashing (cost 10)
- SHA-256 hashed refresh tokens in database
- Token revocation (logout/logout-all)
- Complete user data isolation

**Production deployment** (https://scentora.thejoeshow.net):
- JWT_SECRET and DB_PASSWORD are generated by the setup script
- HTTPS via Let's Encrypt with auto-renewal
- Rate limiting disabled in test environment (`NODE_ENV=test`) to prevent CI failures
- Consider migrating rate limiting to Redis for multi-instance deployments

## Deployment

Scentora is deployed at **https://scentora.thejoeshow.net** on a DigitalOcean droplet with Docker, PostgreSQL, and automated CI/CD.

**CI/CD Pipeline**: GitHub Actions runs tests on PRs, then builds and deploys on push to `main`.

**Infrastructure**:
- 3 Docker containers: frontend (nginx:alpine), backend (Node.js), postgres (15-alpine)
- Host nginx reverse proxy with Let's Encrypt SSL (auto-renewal)
- Daily database backups at 2 AM (30-day retention)
- UFW firewall (ports 22, 80, 443 only)

**Quick Deploy** (fresh Ubuntu 24.04 droplet):
```bash
curl -sSL https://raw.githubusercontent.com/xupit3r/scentora/main/deploy/setup-droplet.sh | bash -s your-domain.com your@email.com
```

**See [Deployment Guide](docs/DEPLOYMENT.md) for complete instructions.**

## Documentation

- **[Deployment Guide](docs/DEPLOYMENT.md)** - Complete production deployment guide
- [API Spec](specs/api-spec.md) - Accord API specification
- [Recipe API](specs/recipe-api.md) - Recipe system API specification
- [Data Models](specs/data-models.md) - Database schema
- [Tag System](specs/tag-system.md) - Predefined tag categories

## License

MIT
