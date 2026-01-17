# Scentora

A perfume cataloging application for tracking scent profiles, notes, and personal journal entries.

## Tech Stack

- **Backend**: Koa.js + TypeScript + CouchDB
- **Frontend**: Vue.js 3 + TypeScript
- **Database**: CouchDB

## Getting Started

### Prerequisites

- Node.js 18+ 
- Docker & Docker Compose
- npm or yarn

### Setup

1. **Start CouchDB**:
   ```bash
   docker-compose up -d
   ```
   Access CouchDB admin at http://localhost:5984/_utils (admin/password)

2. **Backend Setup**:
   ```bash
   cd backend
   npm install
   npm run dev
   ```

3. **Frontend Setup**:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## Features

- Catalog perfumes with detailed information
- Track perfume pyramid (top, middle, base notes)
- Personal journal entries for each perfume
- Search and filter by notes, designer, etc.

## Project Structure

```
scentora/
├── backend/              # Koa.js API server
├── frontend/             # Vue.js SPA
└── docker-compose.yml    # CouchDB setup
```
