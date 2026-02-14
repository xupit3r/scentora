# Quick Start Guide

## Prerequisites

- Node.js 20+
- Docker & Docker Compose (for PostgreSQL)

## Starting the Application

### 1. Start PostgreSQL (requires Docker)
```bash
docker compose up -d
```

PostgreSQL will be available at `localhost:5435` (mapped from container port 5432).

### 2. Start the Backend
```bash
cd backend
npm install
cp .env.example .env    # First time only
npx prisma generate     # First time only
npx prisma db push      # First time only (creates tables)
npm run dev
```

The API will be available at http://localhost:3000
- Health check: http://localhost:3000/api/health

### 3. Start the Frontend
```bash
cd frontend
npm install
npm run dev
```

The app will be available at http://localhost:5173

## Backend Stack

- **Koa.js** - HTTP framework
- **TypeScript** - Type safety
- **Prisma** - Database ORM (14 models)
- **Zod** - Request validation
- **jsonwebtoken** - JWT authentication
- **bcryptjs** - Password hashing
- **Vitest + Supertest** - Testing

## Architecture

```
Routes → Services → Prisma (no repository layer)
```

- **Routes**: HTTP handlers, request validation, response formatting
- **Services**: Business logic, authorization checks
- **Prisma**: Database queries via generated client

## Running Tests

```bash
# Ensure PostgreSQL is running
docker compose up -d

# Run all tests
cd backend
npm test

# Watch mode
npm run test:watch
```

## Key Scripts

```bash
# Backend
npm run dev          # Start dev server with hot reload
npm run build        # Compile TypeScript
npm run start        # Run compiled JS
npm test             # Run tests
npx prisma studio    # Visual database browser

# Frontend
npm run dev          # Start dev server
npm run build        # Production build
```

## API Overview

| Domain       | Endpoints | Auth     |
|-------------|-----------|----------|
| Auth         | 6         | Mixed    |
| Invitations  | 3         | Required |
| Accords      | 7         | Required |
| Tags         | 5         | Public   |
| Stats        | 1         | Required |
| Export       | 2         | Required |
| Recipes      | 21        | Required |
| Collections  | 7         | Required |
| Health       | 1         | Public   |

**Total: 53 endpoints**
