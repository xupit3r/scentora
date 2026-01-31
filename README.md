# Scentora

A perfume formulation and accord management system for DIY perfumers and enthusiasts.

## 🎯 Project Status

**Current Phase**: Transitioning to Accord Inventory System  
**Latest**: Phase 8.2 Backend + Phase 8.3 Frontend Complete  
**Next**: Phase 8.9 - Notion-Inspired UI/UX Redesign

## Tech Stack

- **Backend**: Go + Echo + PostgreSQL + JWT Authentication
- **Frontend**: Vue.js 3 + TypeScript + Pinia
- **Database**: PostgreSQL
- **Planned UI**: Naive UI + Tailwind CSS (Phase 8.9)

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

### Accord Management (Current Focus)
- ✨ Manage perfume accords and essential oils
- 🏺 Track pyramid position (top, middle, base notes)
- 📦 Inventory management (volume tracking in ml and drops)
- 🏷️ Rich tagging system with 57+ predefined tags
- 🔍 Advanced filtering by position, volume, supplier, tags
- 🔎 Full-text search across names, notes, and metadata
- ⚠️ Low stock warnings
- 📊 Tag-based organization and discovery

### UI/UX (Phase 8.9 - Planned)
- 🎨 Notion-inspired clean, minimalist interface
- 📱 Responsive design (desktop, tablet, mobile)
- ⌨️ Keyboard shortcuts for power users
- ✏️ Inline editing for quick updates
- 🎯 Sidebar navigation with collapsible sections
- 🌙 Dark mode support (future)

### Analytics & Export (Future)
- 📈 Collection statistics dashboard
- 📤 Export/import collection as JSON

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

- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [AUTH_IMPLEMENTATION.md](docs/AUTH_IMPLEMENTATION.md) - Authentication details
- [REFRESH_TOKENS_RATE_LIMITING.md](docs/REFRESH_TOKENS_RATE_LIMITING.md) - Refresh tokens & rate limiting
- [LAUNCHER_GUIDE.md](docs/LAUNCHER_GUIDE.md) - Launcher script usage
- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Project overview
- [PHASE1_COMPLETE.md](PHASE1_COMPLETE.md) - Foundation phase
- [PHASE2_COMPLETE.md](PHASE2_COMPLETE.md) - CRUD operations phase
- [PHASE3_COMPLETE.md](PHASE3_COMPLETE.md) - Advanced features phase
- [PHASE4_COMPLETE.md](PHASE4_COMPLETE.md) - Analytics phase
- [PHASE5_COMPLETE.md](PHASE5_COMPLETE.md) - Authentication phase
- [PHASE5B_COMPLETE.md](PHASE5B_COMPLETE.md) - Refresh tokens & rate limiting phase

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
