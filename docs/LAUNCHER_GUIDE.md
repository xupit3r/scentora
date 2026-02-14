# Scentora Launcher Scripts

This directory contains launcher scripts to easily start, stop, and manage the entire Scentora application stack.

## Quick Start

### Linux/Mac
```bash
./scentora.sh start
```

### Windows
```cmd
scentora.bat start
```

### Using npm
```bash
npm start
```

## Available Scripts

### Main Launcher (scentora.sh / scentora.bat)

#### Start All Services
```bash
./scentora.sh start
# or just
./scentora.sh
# or
npm start
```

This will:
1. Start CouchDB (via Docker)
2. Install dependencies if needed
3. Start the Backend API (port 3000)
4. Start the Frontend (port 5173)
5. Monitor services and restart if they crash
6. Show service URLs and log locations

#### Stop All Services
```bash
./scentora.sh stop
# or
npm stop
```

Gracefully stops all services:
- Backend API
- Frontend dev server
- CouchDB Docker container

#### Restart All Services
```bash
./scentora.sh restart
# or
npm restart
```

Stops and starts all services with a 2-second delay.

#### Check Status
```bash
./scentora.sh status
# or
npm run status
```

Shows the current status of all services:
- ● CouchDB - http://localhost:5984
- ● Backend - http://localhost:3000
- ● Frontend - http://localhost:5173

Green ● means running, Red ● means stopped.

#### View Logs
```bash
# Backend logs
./scentora.sh logs backend
# or
npm run logs:backend

# Frontend logs
./scentora.sh logs frontend
# or
npm run logs:frontend
```

Tails the log files in real-time. Press Ctrl+C to exit.

Log files are stored in:
- `logs/backend.log`
- `logs/frontend.log`

## NPM Scripts Reference

### Development
```bash
npm start              # Start all services
npm stop               # Stop all services
npm restart            # Restart all services
npm run status         # Check service status
```

### Setup
```bash
npm run setup          # Install all dependencies
npm run setup:backend  # Install backend dependencies only
npm run setup:frontend # Install frontend dependencies only
```

### Building
```bash
npm run build          # Build both backend and frontend
npm run build:backend  # Build backend only
npm run build:frontend # Build frontend only
```

### Testing
```bash
npm test               # Run all tests
npm run test:backend   # Run backend tests only
npm run test:frontend  # Run frontend tests only
```

### Docker
```bash
npm run docker:up      # Start CouchDB only
npm run docker:down    # Stop CouchDB only
npm run docker:logs    # View CouchDB logs
```

### Maintenance
```bash
npm run clean          # Remove all node_modules and build artifacts
```

## Features

### Automatic Dependency Installation
The launcher checks if `node_modules` exist and automatically runs `npm install` if needed.

### Health Checks
The launcher waits for each service to be ready before proceeding:
- CouchDB: Checks http://localhost:5984
- Backend: Checks http://localhost:3000/api/health
- Frontend: Checks http://localhost:5173

### Process Management
- Services run in background with PID tracking
- Graceful shutdown with Ctrl+C
- Automatic cleanup on exit
- Service monitoring (restarts on crash)

### Colored Output
The launcher uses color-coded output for better readability:
- 🟢 Green: Success messages
- 🔴 Red: Error messages
- 🟡 Yellow: Warning messages
- 🔵 Blue: Info messages

## System Requirements

### Required
- Node.js 20+
- npm or yarn
- Docker (for CouchDB)

### Optional
- curl (for health checks)
- sudo (if Docker requires elevated privileges)

## Troubleshooting

### "Docker permission denied"
If you get permission denied errors with Docker:

**Linux:**
```bash
# Add your user to docker group
sudo usermod -aG docker $USER
# Log out and log back in
```

**Or run with sudo:**
```bash
sudo ./scentora.sh start
```

### "Port already in use"
If you get port conflicts:

1. Check what's using the ports:
```bash
# Linux/Mac
lsof -i :3000  # Backend
lsof -i :5173  # Frontend
lsof -i :5984  # CouchDB

# Windows
netstat -ano | findstr :3000
```

2. Stop the conflicting service or change ports in config

### Services won't start
1. Check the log files:
```bash
cat logs/backend.log
cat logs/frontend.log
```

2. Ensure dependencies are installed:
```bash
npm run setup
```

3. Check if CouchDB is running:
```bash
curl http://localhost:5984
```

### Script not executable (Linux/Mac)
```bash
chmod +x scentora.sh
```

## Advanced Usage

### Custom Docker Compose
If you need to customize CouchDB settings:

1. Edit `docker-compose.yml`
2. Restart services:
```bash
./scentora.sh restart
```

### Environment Variables
To use custom environment variables:

1. Edit `backend/.env`
2. Restart services:
```bash
./scentora.sh restart
```

### Running in Production
The launcher is designed for development. For production:

1. Build the applications:
```bash
npm run build
```

2. Use a process manager like PM2:
```bash
pm2 start backend/dist/index.js --name scentora-api
```

3. Serve frontend with nginx or similar

### Background Mode (Linux/Mac)
To run in detached mode:

```bash
nohup ./scentora.sh start > scentora.out 2>&1 &
```

To stop:
```bash
./scentora.sh stop
```

## Examples

### First Time Setup
```bash
# Clone the repository
git clone <repo-url>
cd scentora

# Make launcher executable (Linux/Mac)
chmod +x scentora.sh

# Install dependencies
npm run setup

# Start everything
./scentora.sh start
```

### Daily Development Workflow
```bash
# Start in the morning
./scentora.sh start

# Check status anytime
./scentora.sh status

# View logs while debugging
./scentora.sh logs backend

# Stop at end of day
./scentora.sh stop
```

### Fixing a Crashed Service
```bash
# If backend crashes, check logs
./scentora.sh logs backend

# Fix the issue, then restart
./scentora.sh restart
```

## Platform-Specific Notes

### Linux
- Use `scentora.sh`
- May need sudo for Docker
- All features fully supported

### macOS
- Use `scentora.sh`
- Docker Desktop recommended
- All features fully supported

### Windows
- Use `scentora.bat`
- Docker Desktop required
- Services open in separate windows
- Manual window closure needed for stop

### Windows Subsystem for Linux (WSL)
- Use `scentora.sh`
- Docker Desktop with WSL2 backend
- All features fully supported

## File Structure

```
scentora/
├── scentora.sh         # Main launcher (Linux/Mac)
├── scentora.bat        # Windows launcher
├── package.json        # NPM scripts
├── docker-compose.yml  # CouchDB config
├── logs/               # Runtime logs (created automatically)
│   ├── backend.log
│   └── frontend.log
├── backend/            # Backend API
│   ├── src/
│   └── package.json
└── frontend/           # Frontend app
    ├── src/
    └── package.json
```

## Support

For issues or questions:
1. Check the troubleshooting section above
2. View logs: `./scentora.sh logs backend|frontend`
3. Check service status: `./scentora.sh status`
4. Review documentation in the docs folder

## Contributing

When adding new services or features:
1. Update `scentora.sh` with new service logic
2. Update `scentora.bat` for Windows support
3. Add NPM scripts to root `package.json`
4. Update this README

## License

MIT
