#!/bin/bash

# Scentora Application Launcher
# Manages the full stack: CouchDB, Backend API, and Frontend

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Process tracking
BACKEND_PID=""
FRONTEND_PID=""
COUCHDB_RUNNING=false

# Directories
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"

# Log files
LOG_DIR="$SCRIPT_DIR/logs"
BACKEND_LOG="$LOG_DIR/backend.log"
FRONTEND_LOG="$LOG_DIR/frontend.log"

# Create logs directory
mkdir -p "$LOG_DIR"

# Print with color
print_header() {
    echo -e "${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}  ${MAGENTA}🌸 Scentora Application Manager 🌸${NC}       ${CYAN}║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if dependencies are installed
check_dependencies() {
    local missing=()
    
    if ! command -v node &> /dev/null; then
        missing+=("node")
    fi
    
    if ! command -v npm &> /dev/null; then
        missing+=("npm")
    fi
    
    if ! command -v docker &> /dev/null; then
        missing+=("docker")
    fi
    
    if [ ${#missing[@]} -ne 0 ]; then
        print_error "Missing dependencies: ${missing[*]}"
        echo ""
        echo "Please install:"
        for dep in "${missing[@]}"; do
            echo "  - $dep"
        done
        exit 1
    fi
}

# Check if CouchDB is running
check_couchdb() {
    if curl -s http://localhost:5984 > /dev/null 2>&1; then
        COUCHDB_RUNNING=true
        return 0
    else
        COUCHDB_RUNNING=false
        return 1
    fi
}

# Start CouchDB
start_couchdb() {
    print_info "Starting CouchDB..."
    
    if check_couchdb; then
        print_success "CouchDB already running"
        return 0
    fi
    
    cd "$SCRIPT_DIR"
    
    # Try with sudo if permission denied
    if docker-compose up -d 2>&1 | grep -q "permission denied"; then
        print_warning "Docker requires elevated privileges"
        if command -v sudo &> /dev/null; then
            sudo docker-compose up -d > /dev/null 2>&1 || true
        fi
    else
        docker-compose up -d > /dev/null 2>&1 || true
    fi
    
    # Wait for CouchDB to be ready
    local retries=0
    local max_retries=30
    
    while [ $retries -lt $max_retries ]; do
        if check_couchdb; then
            print_success "CouchDB is ready (http://localhost:5984)"
            return 0
        fi
        sleep 1
        retries=$((retries + 1))
    done
    
    print_warning "CouchDB may not be running (will continue anyway)"
    return 0
}

# Stop CouchDB
stop_couchdb() {
    print_info "Stopping CouchDB..."
    cd "$SCRIPT_DIR"
    
    if docker-compose down > /dev/null 2>&1; then
        print_success "CouchDB stopped"
    elif sudo docker-compose down > /dev/null 2>&1; then
        print_success "CouchDB stopped"
    else
        print_warning "Could not stop CouchDB (may not be running)"
    fi
}

# Check if NPM dependencies are installed
check_npm_deps() {
    local dir=$1
    local name=$2
    
    if [ ! -d "$dir/node_modules" ]; then
        print_warning "$name dependencies not installed"
        print_info "Installing $name dependencies..."
        cd "$dir"
        npm install > /dev/null 2>&1
        print_success "$name dependencies installed"
    fi
}

# Start Backend
start_backend() {
    print_info "Starting Backend API..."
    
    check_npm_deps "$BACKEND_DIR" "Backend"
    
    cd "$BACKEND_DIR"
    
    # Start backend in background
    npm run dev > "$BACKEND_LOG" 2>&1 &
    BACKEND_PID=$!
    
    # Wait for backend to be ready
    local retries=0
    local max_retries=30
    
    while [ $retries -lt $max_retries ]; do
        if curl -s http://localhost:3000/api/health > /dev/null 2>&1; then
            print_success "Backend API ready (http://localhost:3000)"
            return 0
        fi
        sleep 1
        retries=$((retries + 1))
    done
    
    print_error "Backend failed to start. Check logs: $BACKEND_LOG"
    return 1
}

# Start Frontend
start_frontend() {
    print_info "Starting Frontend..."
    
    check_npm_deps "$FRONTEND_DIR" "Frontend"
    
    cd "$FRONTEND_DIR"
    
    # Start frontend in background
    npm run dev > "$FRONTEND_LOG" 2>&1 &
    FRONTEND_PID=$!
    
    # Wait for frontend to be ready
    local retries=0
    local max_retries=30
    
    while [ $retries -lt $max_retries ]; do
        if curl -s http://localhost:5173 > /dev/null 2>&1; then
            print_success "Frontend ready (http://localhost:5173)"
            return 0
        fi
        sleep 1
        retries=$((retries + 1))
    done
    
    print_error "Frontend failed to start. Check logs: $FRONTEND_LOG"
    return 1
}

# Stop all services
stop_all() {
    echo ""
    print_info "Stopping all services..."
    
    if [ ! -z "$BACKEND_PID" ]; then
        print_info "Stopping Backend..."
        kill $BACKEND_PID 2>/dev/null || true
        print_success "Backend stopped"
    fi
    
    if [ ! -z "$FRONTEND_PID" ]; then
        print_info "Stopping Frontend..."
        kill $FRONTEND_PID 2>/dev/null || true
        print_success "Frontend stopped"
    fi
    
    stop_couchdb
    
    echo ""
    print_success "All services stopped"
}

# Cleanup on exit
cleanup() {
    echo ""
    print_info "Shutting down gracefully..."
    stop_all
    exit 0
}

# Trap signals for graceful shutdown
trap cleanup SIGINT SIGTERM

# Show status
show_status() {
    print_header
    echo -e "${CYAN}Service Status:${NC}"
    echo ""
    
    # CouchDB
    if check_couchdb; then
        echo -e "  ${GREEN}●${NC} CouchDB    - http://localhost:5984"
    else
        echo -e "  ${RED}●${NC} CouchDB    - Not running"
    fi
    
    # Backend
    if curl -s http://localhost:3000/api/health > /dev/null 2>&1; then
        echo -e "  ${GREEN}●${NC} Backend    - http://localhost:3000"
    else
        echo -e "  ${RED}●${NC} Backend    - Not running"
    fi
    
    # Frontend
    if curl -s http://localhost:5173 > /dev/null 2>&1; then
        echo -e "  ${GREEN}●${NC} Frontend   - http://localhost:5173"
    else
        echo -e "  ${RED}●${NC} Frontend   - Not running"
    fi
    
    echo ""
}

# Show logs
show_logs() {
    local service=$1
    
    case $service in
        backend)
            if [ -f "$BACKEND_LOG" ]; then
                tail -f "$BACKEND_LOG"
            else
                print_error "Backend log not found"
            fi
            ;;
        frontend)
            if [ -f "$FRONTEND_LOG" ]; then
                tail -f "$FRONTEND_LOG"
            else
                print_error "Frontend log not found"
            fi
            ;;
        *)
            print_error "Unknown service: $service"
            echo "Available services: backend, frontend"
            ;;
    esac
}

# Start all services
start_all() {
    print_header
    
    check_dependencies
    
    echo -e "${YELLOW}Starting all services...${NC}"
    echo ""
    
    start_couchdb
    start_backend
    start_frontend
    
    echo ""
    print_success "All services started!"
    echo ""
    echo -e "${CYAN}URLs:${NC}"
    echo -e "  Frontend:  ${GREEN}http://localhost:5173${NC}"
    echo -e "  Backend:   ${GREEN}http://localhost:3000${NC}"
    echo -e "  CouchDB:   ${GREEN}http://localhost:5984/_utils${NC}"
    echo ""
    echo -e "${CYAN}Logs:${NC}"
    echo -e "  Backend:   $BACKEND_LOG"
    echo -e "  Frontend:  $FRONTEND_LOG"
    echo ""
    print_info "Press Ctrl+C to stop all services"
    echo ""
    
    # Keep script running and monitor services
    while true; do
        # Check if services are still running
        if [ ! -z "$BACKEND_PID" ]; then
            if ! kill -0 $BACKEND_PID 2>/dev/null; then
                print_error "Backend crashed! Check logs: $BACKEND_LOG"
                stop_all
                exit 1
            fi
        fi
        
        if [ ! -z "$FRONTEND_PID" ]; then
            if ! kill -0 $FRONTEND_PID 2>/dev/null; then
                print_error "Frontend crashed! Check logs: $FRONTEND_LOG"
                stop_all
                exit 1
            fi
        fi
        
        sleep 5
    done
}

# Show help
show_help() {
    print_header
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start       Start all services (default)"
    echo "  stop        Stop all services"
    echo "  restart     Restart all services"
    echo "  status      Show service status"
    echo "  logs        Show logs (specify: backend or frontend)"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                      # Start all services"
    echo "  $0 start                # Start all services"
    echo "  $0 status               # Check status"
    echo "  $0 logs backend         # View backend logs"
    echo "  $0 stop                 # Stop all services"
    echo ""
}

# Main command handler
case "${1:-start}" in
    start)
        start_all
        ;;
    stop)
        print_header
        stop_all
        ;;
    restart)
        print_header
        stop_all
        echo ""
        sleep 2
        start_all
        ;;
    status)
        show_status
        ;;
    logs)
        if [ -z "$2" ]; then
            print_error "Please specify service: backend or frontend"
            exit 1
        fi
        show_logs "$2"
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
