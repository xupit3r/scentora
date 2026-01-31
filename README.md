# Scentora

A perfume formulation and accord management system for DIY perfumers and enthusiasts.

## 🎯 Project Status

**Current Phase**: Testing & Quality Assurance (Phase 9)  
**Latest**: Phase 9.2 (Backend Service Tests) Complete - 116 tests passing ✅  
**Coverage**: Repository 59.9% | Services 59.6%  
**Next**: Phase 9.3 (Backend Handler/Integration Tests)

## Tech Stack

- **Backend**: Go + Echo + PostgreSQL + JWT Authentication
- **Frontend**: Vue.js 3 + TypeScript + Pinia + Naive UI + Tailwind CSS v4
- **Database**: PostgreSQL
- **Design System**: Notion-inspired minimalist UI

## Features

### Authentication & Security
- 🔐 JWT-based authentication with refresh tokens
- 🔄 Automatic token refresh (15-min access, 7-day refresh)
- 🛡️ Rate limiting on auth endpoints (brute force protection)
- 👥 Multi-user support with complete data isolation
- 🔒 Secure password hashing (bcrypt)
- 📱 Persistent sessions with logout-all capability
- 🚫 Token rotation for enhanced security
- 🎫 Invitation-only registration system

### Accord Management (Core Features)
- ✨ Manage perfume accords and essential oils
- 🏺 Track pyramid position (top, middle, base notes)
- 📦 Inventory management (volume tracking in ml and drops)
- 🏷️ Rich tagging system with 57+ predefined tags
- 🔍 Advanced filtering by position, volume, supplier, tags
- 🔎 Full-text search across names, notes, and metadata
- ⚠️ Low stock warnings and inventory alerts
- 📊 Tag-based organization and discovery

### Analytics & Export (Phase 8.4 Complete)
- 📈 Collection statistics dashboard
- 📊 Pyramid position distribution analysis
- 🏷️ Tag usage statistics
- 🏢 Supplier breakdown
- 📉 Volume analytics (min/max/average)
- ⚠️ Low inventory alerts (< 10ml threshold)
- 📤 Export collection as JSON
- 📥 Import accords from JSON backup

### UI/UX (Phase 8.9 Complete)
- 🎨 Notion-inspired clean, minimalist interface
- 📱 Fully responsive design (desktop, tablet, mobile)
- ⌨️ Keyboard shortcuts (N for new accord)
- 🎯 Collapsible sidebar navigation (280px ↔ 64px)
- 🎭 Skeleton loading states with shimmer animation
- 💫 Smooth transitions (200ms standard)
- 🎪 Empty state designs with helpful CTAs
- 🌈 Subtle pastel color accents
- 📏 8px spacing grid for visual consistency
- ✨ Hover-reveal action buttons
- 🔔 Toast notifications for user feedback
- 🖼️ **[View Design Mockups](mockups/)** - Interactive gallery of wireframes and high-fidelity designs

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+ 
- Docker & Docker Compose (for PostgreSQL)
- npm or yarn

### Quick Start (Automated)

**The easiest way to run Scentora:**

```bash
# Linux/Mac
./scentora.sh start

# Windows
scentora.bat start

# Or using npm
npm start
```

That's it! The launcher will:
- Start PostgreSQL
- Build Go backend if needed
- Install frontend dependencies if needed
- Start backend and frontend
- Show you the URLs

Access the app at http://localhost:5173

See [LAUNCHER_GUIDE.md](docs/LAUNCHER_GUIDE.md) for more launcher options.

### Manual Setup

If you prefer to start services manually:

1. **Clone and navigate to the project**:
   ```bash
   git clone <repo-url>
   cd scentora
   ```

2. **Start PostgreSQL**:
   ```bash
   docker compose up -d postgres
   ```

3. **Backend Setup**:
   ```bash
   cd backend-go
   
   # Configure environment (optional - defaults work locally)
   cp .env.example .env
   # Edit .env to set JWT_SECRET in production
   
   # Build and run
   go build -o scentora-backend cmd/server/main.go
   ./scentora-backend
   ```
   Backend runs at http://localhost:3000

4. **Frontend Setup** (in a new terminal):
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
   Frontend runs at http://localhost:5173

5. **Access the app**:
   - Open http://localhost:5173
   - You'll need an invitation code to register
   - See [Creating Invitations](#creating-invitations) below
   - Start cataloging your perfumes!

## Environment Variables

Create `backend/.env` (or use defaults):
```env
# Server
PORT=3000
NODE_ENV=development

# CouchDB
COUCHDB_URL=http://localhost:5984
COUCHDB_USER=admin
COUCHDB_PASSWORD=password
COUCHDB_DATABASE=scentora

# JWT (CHANGE IN PRODUCTION!)
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_EXPIRES_IN=15m   # Access token expiry
JWT_REFRESH_EXPIRES_IN=7d   # Refresh token expiry
```

## Creating Invitations

Scentora uses an invitation-only registration system. To create your first user and generate invitations:

### Method 1: Direct Database Creation (First User)

For the very first user, create an invitation directly in CouchDB:

```bash
# Generate a random invitation code
INVITATION_CODE=$(openssl rand -hex 16)
echo "Your invitation code: $INVITATION_CODE"

# Create invitation in database
curl -X POST "http://admin:password@localhost:5984/scentora" \
  -H "Content-Type: application/json" \
  -d "{
    \"type\": \"invitation\",
    \"code\": \"$INVITATION_CODE\",
    \"createdBy\": \"system\",
    \"expiresAt\": \"$(date -u -d '+7 days' +%Y-%m-%dT%H:%M:%S.000Z)\",
    \"used\": false,
    \"createdAt\": \"$(date -u +%Y-%m-%dT%H:%M:%S.000Z)\"
  }"
```

Then use this code to register at http://localhost:5173/register

### Method 2: API (Authenticated Users)

Once you have an account, you can create invitations through the API:

```bash
# Login and get access token
ACCESS_TOKEN="your-access-token"

# Create invitation
curl -X POST "http://localhost:3000/api/invitations" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "email": "friend@example.com",
    "expiresInDays": 7
  }'
```

The response will include the invitation code to share with the new user.

## Project Structure

```
scentora/
├── backend/              # Koa.js API server
├── frontend/             # Vue.js SPA
└── docker-compose.yml    # CouchDB setup
```

## API Endpoints

### Authentication (Public)
- `POST /api/auth/register` - Create new account (requires invitation code)
- `POST /api/auth/login` - Login with email/password
- `GET /api/auth/me` - Get current user info

### Invitations (Protected)
- `POST /api/invitations` - Create a new invitation
- `GET /api/invitations` - List your created invitations
- `DELETE /api/invitations/:code` - Revoke an invitation

### Perfumes (Protected)
- `GET /api/perfumes` - List all perfumes (with filters)
- `GET /api/perfumes/:id` - Get perfume details
- `POST /api/perfumes` - Create new perfume
- `PUT /api/perfumes/:id` - Update perfume
- `DELETE /api/perfumes/:id` - Delete perfume

### Journal Entries (Protected)
- `GET /api/perfumes/:id/journal` - Get journal entries for perfume
- `POST /api/perfumes/:id/journal` - Create journal entry
- `PUT /api/journal/:id` - Update journal entry
- `DELETE /api/journal/:id` - Delete journal entry

### Other (Protected)
- `GET /api/notes` - Get all unique notes in collection
- `GET /api/stats` - Get collection statistics
- `GET /api/export/collection` - Export collection as JSON
- `POST /api/export/import` - Import collection from JSON

## Documentation

### Design & Planning
- **[Design Mockups](mockups/)** - Interactive gallery of wireframes and high-fidelity designs for Phase 8.9
- [PLAN.md](PLAN.md) - Complete project roadmap and phase planning
- [TESTING_PLAN.md](docs/TESTING_PLAN.md) - Comprehensive testing strategy and timeline

### Technical Documentation
- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [AUTH_IMPLEMENTATION.md](docs/AUTH_IMPLEMENTATION.md) - Authentication details
- [REFRESH_TOKENS_RATE_LIMITING.md](docs/REFRESH_TOKENS_RATE_LIMITING.md) - Refresh tokens & rate limiting
- [LAUNCHER_GUIDE.md](docs/LAUNCHER_GUIDE.md) - Launcher script usage
- [TESTING_IMPLEMENTATION.md](TESTING_IMPLEMENTATION.md) - Testing framework and guidelines

### Testing
- [TESTING_PLAN.md](docs/TESTING_PLAN.md) - Comprehensive testing plan (Phase 9)
- **Backend Tests:** 116 tests passing ✅
  - Repository layer: 59.9% coverage (40 tests)
  - Service layer: 59.6% coverage (70 tests)
    - AuthService: 20 tests (registration, login, tokens)
    - InvitationService: 16 tests (creation, validation)
    - TagService: 11 tests (search, categorization)
    - AccordService: 23 tests (CRUD, validation)
  - Run: `cd backend && go test ./...`
- **Frontend Tests:** Setup in progress

### Phase History
All phase completion documents are in [docs/phases/](docs/phases/):
- [PHASE9_2_COMPLETE.md](docs/phases/PHASE9_2_COMPLETE.md) - Backend service tests (70 tests) ✅
- [PHASE9_1_COMPLETE.md](docs/phases/PHASE9_1_COMPLETE.md) - Backend repository tests (40 tests) ✅
- [PHASE8_9_COMPLETE.md](docs/phases/PHASE8_9_COMPLETE.md) - Notion-inspired UI redesign
- [PHASE8_8_COMPLETE.md](docs/phases/PHASE8_8_COMPLETE.md) - Features & Polish
- [PHASE8_5_COMPLETE.md](docs/phases/PHASE8_5_COMPLETE.md) - Frontend cleanup
- [PHASE8_4_COMPLETE.md](docs/phases/PHASE8_4_COMPLETE.md) - Statistics & Export
- See [docs/phases/README.md](docs/phases/README.md) for complete list

## Security Features

✅ **Implemented:**
- JWT refresh tokens with automatic rotation
- Short-lived access tokens (15 minutes)
- Rate limiting on auth endpoints (5 requests/15 min)
- Bcrypt password hashing
- Token revocation (logout/logout-all)
- Complete user data isolation

⚠️ **Before deploying to production:**
- Change `JWT_SECRET` to a strong random value
- Enable HTTPS
- Consider migrating rate limiting to Redis
- Review all security recommendations in documentation
- Add rate limiting for auth endpoints

## License

MIT

## Contributing

Contributions welcome! Please feel free to submit a Pull Request.
