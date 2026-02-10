# Documentation Directory

This directory contains user-facing and feature documentation for Scentora.

## Contents

### [AUTH_IMPLEMENTATION.md](AUTH_IMPLEMENTATION.md)
Complete guide to the authentication system:
- JWT-based authentication
- User registration and login
- Password security with bcrypt
- Protected routes and user isolation
- Frontend auth store and route guards
- Security best practices

### [REFRESH_TOKENS_RATE_LIMITING.md](REFRESH_TOKENS_RATE_LIMITING.md)
Security features documentation:
- JWT refresh token system
- Token rotation and automatic refresh
- Short-lived access tokens (15 minutes)
- Long-lived refresh tokens (7 days)
- Rate limiting on auth endpoints
- Brute force protection

### [LAUNCHER_GUIDE.md](LAUNCHER_GUIDE.md)
User guide for launcher scripts:
- Quick start commands
- Starting, stopping, and restarting services
- Checking service status
- Viewing logs
- Troubleshooting common issues
- Platform-specific instructions (Linux/Mac/Windows)

### [CLI_CONSOLE.md](CLI_CONSOLE.md)
CLI Console user guide:
- Interactive command-line interface for administration
- User management commands (create, list, delete, reset-password)
- Terminal features (auto-complete, history, colored output)
- Common tasks and troubleshooting
- Future: Accord/recipe management and database operations

## Related Documentation

- **/README.md** - Project overview and quick start
- **/QUICKSTART.md** - Getting started guide
- **/PLAN.md** - Comprehensive development plan
- **/specs/** - Technical specifications
- **/specs/api-spec.md** - API endpoint specifications
- **/specs/archives/** - Historical documentation

## For Developers

When adding new features:
1. Update relevant documentation in this directory
2. Keep docs in sync with implementation
3. Add examples and use cases
4. Include troubleshooting tips
