# Migration from Node.js to Go Backend

## Date: January 30, 2026

Scentora has successfully migrated from a Node.js/Koa + CouchDB backend to a Go/Echo + PostgreSQL backend.

## Why the Migration?

1. **Performance**: Go's compiled nature and efficient concurrency model provide better performance
2. **Type Safety**: Static typing reduces runtime errors
3. **Lower Resource Usage**: Smaller memory footprint and faster startup times
4. **Better Tooling**: Built-in testing, profiling, and cross-compilation
5. **Production Ready**: PostgreSQL is more suitable for production than CouchDB for this use case

## What Changed

### Backend
- **Language**: TypeScript/Koa → Go/Echo
- **Database**: CouchDB → PostgreSQL
- **Location**: `backend/` → `backend-go/`

### API Compatibility
✅ **100% API Compatible** - No frontend changes required
- All endpoints remain the same
- Request/response formats unchanged
- Authentication flow identical
- Error handling matches

### New Features Added
- Logout from all devices
- Import collection data
- Enhanced export format (version, exportDate)
- Rate limiting middleware
- Optional authentication middleware

## Migration Steps (If you have existing data)

1. **Export from Node.js backend** (if you were using it):
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:3000/api/export > collection.json
   ```

2. **Start Go backend**:
   ```bash
   cd backend-go
   ./scentora.sh start
   ```

3. **Import to Go backend**:
   ```bash
   curl -X POST -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d @collection.json \
     http://localhost:3000/api/export/import
   ```

## Architecture Changes

### Old Stack (Removed)
- Node.js 18+
- Koa.js web framework
- CouchDB document database
- npm for package management

### New Stack (Current)
- Go 1.21+
- Echo web framework
- PostgreSQL relational database
- Go modules for dependency management

## Files Removed

The following have been removed:
- `backend/` directory (entire Node.js backend)
- CouchDB container from docker-compose.yml
- Node.js backend dependencies from package.json

## Files Updated

- `scentora.sh` - Now starts Go backend
- `scentora.bat` - Now starts Go backend
- `package.json` - Updated scripts for Go backend
- `docker-compose.yml` - PostgreSQL only
- `README.md` - Updated tech stack and instructions

## Rollback (Not Recommended)

If you need to rollback to the Node.js backend:
```bash
git checkout <commit-before-migration> -- backend/
git checkout <commit-before-migration> -- scentora.sh scentora.bat package.json docker-compose.yml
```

However, the Go backend is production-ready and recommended.

## Documentation

For the old Node.js backend implementation details, see:
- Git history: commits before this migration
- Documentation: PHASE*.md files (archived for reference)

For the Go backend:
- See `backend-go/README.md`
- See `backend-go/IMPLEMENTATION_COMPLETE.md`
- See `backend-go/BACKEND_PARITY_COMPLETE.md`

## Support

The Node.js backend is no longer maintained. All future development will be on the Go backend.
