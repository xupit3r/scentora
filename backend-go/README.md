# Scentora Backend - Go Rewrite

Go/Echo backend rewrite of the Scentora perfume cataloging application with PostgreSQL database.

## Tech Stack

- **Framework**: Echo v4
- **Database**: PostgreSQL 15+
- **ORM**: sqlx
- **Authentication**: JWT (golang-jwt)
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator

## Project Structure

```
backend-go/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration and database
│   ├── models/         # Data models and DTOs
│   ├── repository/     # Data access layer
│   ├── services/       # Business logic
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Custom middleware
│   └── routes/         # Route definitions
├── migrations/         # Database migrations
└── pkg/utils/          # Reusable utilities
```

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker (optional, for PostgreSQL)

### 1. Install Dependencies

```bash
cd backend-go
go mod tidy
```

### 2. Start PostgreSQL

Using Docker:
```bash
docker run --name scentora-postgres \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=scentora \
  -p 5432:5432 \
  -d postgres:15
```

Or update docker-compose.yml in the project root.

### 3. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your configuration. **Important**: Set a strong `JWT_SECRET` for production.

### 4. Run the Server

Development:
```bash
go run cmd/server/main.go
```

Build and run:
```bash
go build -o scentora-backend cmd/server/main.go
./scentora-backend
```

The server will:
1. Connect to PostgreSQL
2. Run migrations automatically
3. Start on port 3000 (or configured PORT)

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user (protected)

### Perfumes (all protected)
- `GET /api/perfumes` - List perfumes (supports filters)
- `GET /api/perfumes/:id` - Get single perfume
- `POST /api/perfumes` - Create perfume
- `PUT /api/perfumes/:id` - Update perfume
- `DELETE /api/perfumes/:id` - Delete perfume

### Journal (all protected)
- `GET /api/perfumes/:perfumeId/journal` - List entries
- `POST /api/perfumes/:perfumeId/journal` - Create entry
- `PUT /api/journal/:id` - Update entry
- `DELETE /api/journal/:id` - Delete entry

### Other (all protected)
- `GET /api/notes` - Get all unique notes
- `GET /api/stats` - Get user statistics
- `GET /api/export` - Export user data

### Health Check
- `GET /health` - Server health check

## Database Schema

### Users
- id (UUID, PK)
- email (unique)
- username
- password_hash
- created_at, updated_at

### Perfumes
- id (UUID, PK)
- user_id (FK to users)
- name, designer
- year, concentration
- top_notes, middle_notes, base_notes (TEXT[])
- description, image_url
- created_at, updated_at

### Journal Entries
- id (UUID, PK)
- user_id, perfume_id (FKs)
- date, content
- rating (1-10), occasion, weather
- created_at, updated_at

### Refresh Tokens
- id (UUID, PK)
- user_id (FK)
- token_hash
- expires_at, revoked
- created_at

## Development

### Environment Variables

See `.env.example` for all available configuration options.

Key variables:
- `PORT` - Server port (default: 3000)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - **Required** - Secret for JWT signing
- `JWT_ACCESS_EXPIRES_IN` - Access token expiry (default: 15m)
- `JWT_REFRESH_EXPIRES_IN` - Refresh token expiry (default: 7d)
- `CORS_ALLOWED_ORIGINS` - Comma-separated allowed origins

### API Compatibility

This Go backend maintains 100% API compatibility with the TypeScript/Koa backend, allowing the Vue.js frontend to work without modifications.

Response formats match exactly:
- Field names use camelCase in JSON
- UUIDs are returned as strings
- Arrays and objects match the original structure
- Error responses follow the same format

## Testing

Basic health check:
```bash
curl http://localhost:3000/health
```

Register a user:
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

## Production Deployment

1. Build the binary:
```bash
CGO_ENABLED=0 GOOS=linux go build -o scentora-backend cmd/server/main.go
```

2. Set environment variables (especially `JWT_SECRET`)

3. Ensure PostgreSQL is running and accessible

4. Run migrations will happen automatically on startup

5. Start the server:
```bash
./scentora-backend
```

## Migration from CouchDB

A separate data migration script will be provided to migrate existing data from CouchDB to PostgreSQL.

## License

MIT
