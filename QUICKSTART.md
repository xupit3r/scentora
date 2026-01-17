# Quick Start Guide

## Phase 1 - Complete! ✅

The project structure is now set up with:
- Backend: Koa.js + TypeScript + CouchDB client
- Frontend: Vue.js 3 + TypeScript + Vite
- Docker Compose for CouchDB

## Starting the Application

### 1. Start CouchDB (requires Docker)
```bash
docker compose up -d
```

Visit http://localhost:5984/_utils to access Fauxton (CouchDB admin UI)
- Username: `admin`
- Password: `password`

### 2. Start the Backend
```bash
cd backend
npm run dev
```

The API will be available at http://localhost:3000
- Health check: http://localhost:3000/api/health

### 3. Start the Frontend
```bash
cd frontend
npm run dev
```

The app will be available at http://localhost:5173

## What's Included

### Backend Structure
- ✅ Koa.js server with TypeScript
- ✅ CouchDB connection & initialization
- ✅ Error handling middleware
- ✅ Health check endpoint
- ✅ Type definitions for Perfume & JournalEntry
- ✅ Environment configuration

### Frontend Structure
- ✅ Vue 3 with Composition API
- ✅ TypeScript configuration
- ✅ Vue Router setup
- ✅ Pinia for state management
- ✅ Axios API client with service layer
- ✅ Home page (collection view)
- ✅ About page
- ✅ Basic styling

## Next Steps (Phase 2)

1. Implement perfume CRUD API endpoints
2. Create CouchDB views for querying
3. Add validation middleware
4. Build add/edit perfume forms
5. Implement perfume detail view

## Directory Structure
```
scentora/
├── backend/
│   ├── src/
│   │   ├── config/         # Database & env config
│   │   ├── middleware/     # Error handling
│   │   ├── models/         # TypeScript types
│   │   ├── routes/         # API routes
│   │   └── index.ts        # Main server
│   ├── .env                # Environment variables
│   ├── package.json
│   └── tsconfig.json
├── frontend/
│   ├── src/
│   │   ├── components/     # Vue components (empty)
│   │   ├── views/          # Home, About pages
│   │   ├── router/         # Vue Router config
│   │   ├── services/       # API client
│   │   ├── stores/         # Pinia stores (empty)
│   │   ├── types/          # TypeScript interfaces
│   │   └── main.ts         # App entry point
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
├── docker-compose.yml      # CouchDB container
└── README.md
```

## API Endpoints (to be implemented in Phase 2)
- `GET /api/health` - Health check ✅
- `GET /api/perfumes` - List all perfumes
- `POST /api/perfumes` - Create new perfume
- `GET /api/perfumes/:id` - Get perfume details
- `PUT /api/perfumes/:id` - Update perfume
- `DELETE /api/perfumes/:id` - Delete perfume
- `GET /api/perfumes/:id/journal` - Get journal entries
- `POST /api/perfumes/:id/journal` - Create journal entry
- `PUT /api/journal/:id` - Update journal entry
- `DELETE /api/journal/:id` - Delete journal entry
