#!/bin/bash

###############################################################################
# Scentora - DigitalOcean Droplet Setup Script
# Automated deployment setup for Ubuntu 24.04 with PostgreSQL, Docker, SSL
###############################################################################

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script configuration
DOMAIN="${1:-}"
EMAIL="${2:-}"
APP_DIR="/opt/scentora"
CREDENTIALS_FILE="/root/.scentora-credentials"

###############################################################################
# Helper Functions
###############################################################################

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

check_args() {
    if [[ -z "$DOMAIN" ]]; then
        log_error "Domain name is required"
        echo "Usage: $0 <domain> <email>"
        echo "Example: $0 scentora.example.com admin@example.com"
        exit 1
    fi

    if [[ -z "$EMAIL" ]]; then
        log_error "Email address is required (for Let's Encrypt)"
        echo "Usage: $0 <domain> <email>"
        exit 1
    fi
}

generate_secret() {
    openssl rand -hex 32
}

generate_password() {
    openssl rand -base64 24 | tr -d "=+/" | cut -c1-32
}

###############################################################################
# Installation Functions
###############################################################################

install_docker() {
    log_info "Installing Docker and Docker Compose..."

    # Update package index
    apt-get update

    # Install prerequisites
    apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release

    # Add Docker's official GPG key
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc

    # Set up Docker repository
    echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
        $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
        tee /etc/apt/sources.list.d/docker.list > /dev/null

    # Install Docker Engine
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # Enable and start Docker
    systemctl enable docker
    systemctl start docker

    log_success "Docker installed successfully"
}

configure_firewall() {
    log_info "Configuring UFW firewall..."

    # Install ufw if not present
    apt-get install -y ufw

    # Configure firewall rules
    ufw --force reset
    ufw default deny incoming
    ufw default allow outgoing
    ufw allow 22/tcp comment 'SSH'
    ufw allow 80/tcp comment 'HTTP'
    ufw allow 443/tcp comment 'HTTPS'

    # Enable firewall
    ufw --force enable

    log_success "Firewall configured"
}

install_nginx() {
    log_info "Installing nginx..."

    apt-get install -y nginx
    systemctl enable nginx
    systemctl start nginx

    log_success "nginx installed"
}

install_certbot() {
    log_info "Installing certbot..."

    apt-get install -y certbot python3-certbot-nginx

    log_success "certbot installed"
}

setup_application_directory() {
    log_info "Setting up application directory..."

    mkdir -p "$APP_DIR"
    cd "$APP_DIR"

    log_success "Application directory created"
}

generate_secrets() {
    log_info "Generating secure secrets..."

    JWT_SECRET=$(generate_secret)
    DB_PASSWORD=$(generate_password)

    log_success "Secrets generated"
}

create_env_file() {
    log_info "Creating production environment file..."

    cat > "$APP_DIR/.env.production" <<EOF
# Production Environment Variables
# Generated on $(date)

# Domain Configuration
DOMAIN=$DOMAIN

# Database Configuration
DB_USER=admin
DB_PASSWORD=$DB_PASSWORD
DB_NAME=scentora

# Backend Configuration
BACKEND_PORT=3000
JWT_SECRET=$JWT_SECRET
JWT_EXPIRES_IN=7d
JWT_REFRESH_EXPIRES_IN=30d

# Frontend Configuration
VITE_API_URL=/api

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://$DOMAIN

# Docker Configuration
DOCKER_REGISTRY=scentora
IMAGE_TAG=latest
EOF

    chmod 600 "$APP_DIR/.env.production"

    log_success "Environment file created"
}

save_credentials() {
    log_info "Saving credentials securely..."

    cat > "$CREDENTIALS_FILE" <<EOF
# Scentora Credentials
# Generated on $(date)
# Keep this file secure!

Domain: $DOMAIN
Database Password: $DB_PASSWORD
JWT Secret: $JWT_SECRET

# Application location: $APP_DIR
# Environment file: $APP_DIR/.env.production
EOF

    chmod 600 "$CREDENTIALS_FILE"

    log_success "Credentials saved to $CREDENTIALS_FILE"
}

configure_nginx_proxy() {
    log_info "Configuring nginx reverse proxy..."

    # Create nginx config (quoted heredoc preserves $ for nginx variables)
    cat > /etc/nginx/sites-available/scentora <<'NGINXCONF'
server {
    listen 80;
    server_name __DOMAIN__;

    # Let's Encrypt challenge
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    # Redirect all other HTTP traffic to HTTPS
    location / {
        return 301 https://$server_name$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name __DOMAIN__;

    # SSL certificates (will be configured by certbot)
    ssl_certificate /etc/letsencrypt/live/__DOMAIN__/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/__DOMAIN__/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Proxy to frontend container
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Increase timeouts for API calls
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;

    # Client max body size (for file uploads if needed in future)
    client_max_body_size 10M;
}
NGINXCONF

    # Substitute domain placeholder
    sed -i "s/__DOMAIN__/$DOMAIN/g" /etc/nginx/sites-available/scentora

    # Enable site
    ln -sf /etc/nginx/sites-available/scentora /etc/nginx/sites-enabled/

    # Remove default site
    rm -f /etc/nginx/sites-enabled/default

    # Test nginx configuration
    nginx -t

    log_success "nginx configured"
}

obtain_ssl_certificate() {
    log_info "Obtaining SSL certificate from Let's Encrypt..."

    # Temporarily comment out SSL directives for initial certificate
    sed -i 's/ssl_certificate/#ssl_certificate/' /etc/nginx/sites-available/scentora
    systemctl reload nginx

    # Obtain certificate
    certbot certonly --nginx -d "$DOMAIN" \
        --non-interactive \
        --agree-tos \
        --email "$EMAIL" \
        --no-eff-email

    # Uncomment SSL directives
    sed -i 's/#ssl_certificate/ssl_certificate/' /etc/nginx/sites-available/scentora
    systemctl reload nginx

    # Setup auto-renewal
    systemctl enable certbot.timer
    systemctl start certbot.timer

    log_success "SSL certificate obtained and auto-renewal configured"
}

clone_or_create_compose() {
    log_info "Setting up docker-compose configuration..."

    # For now, we'll assume the files are already on the server
    # In a real deployment, you'd clone from git or copy files
    log_warning "Make sure to copy docker-compose.prod.yml and .env.production to $APP_DIR"

    # Create a minimal docker-compose if it doesn't exist
    if [[ ! -f "$APP_DIR/docker-compose.prod.yml" ]]; then
        log_warning "docker-compose.prod.yml not found. You'll need to deploy it separately."
    fi
}

start_containers() {
    log_info "Starting Docker containers..."

    cd "$APP_DIR"

    if [[ -f "docker-compose.prod.yml" ]]; then
        docker compose -f docker-compose.prod.yml --env-file .env.production up -d
        log_success "Containers started"
    else
        log_warning "docker-compose.prod.yml not found. Skipping container startup."
        log_warning "Deploy your application files and run:"
        log_warning "  cd $APP_DIR && docker compose -f docker-compose.prod.yml up -d"
    fi
}

setup_automatic_updates() {
    log_info "Configuring automatic security updates..."

    apt-get install -y unattended-upgrades

    # Configure unattended-upgrades
    cat > /etc/apt/apt.conf.d/50unattended-upgrades <<EOF
Unattended-Upgrade::Allowed-Origins {
    "\${distro_id}:\${distro_codename}";
    "\${distro_id}:\${distro_codename}-security";
    "\${distro_id}:\${distro_codename}-updates";
};
Unattended-Upgrade::AutoFixInterruptedDpkg "true";
Unattended-Upgrade::MinimalSteps "true";
Unattended-Upgrade::Remove-Unused-Kernel-Packages "true";
Unattended-Upgrade::Remove-Unused-Dependencies "true";
Unattended-Upgrade::Automatic-Reboot "false";
EOF

    # Enable automatic updates
    cat > /etc/apt/apt.conf.d/20auto-upgrades <<EOF
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Download-Upgradeable-Packages "1";
APT::Periodic::AutocleanInterval "7";
APT::Periodic::Unattended-Upgrade "1";
EOF

    log_success "Automatic updates configured"
}

setup_backup_cron() {
    log_info "Setting up daily database backup..."

    # Create backup directory
    mkdir -p /var/backups/scentora

    # Create backup script
    cat > /usr/local/bin/backup-scentora-db.sh <<'EOF'
#!/bin/bash
BACKUP_DIR="/var/backups/scentora"
DATE=$(date +%Y%m%d)
cd /opt/scentora
docker compose -f docker-compose.prod.yml exec -T postgres pg_dump -U admin scentora | gzip > "$BACKUP_DIR/db_$DATE.sql.gz"
# Keep last 30 days of backups
find "$BACKUP_DIR" -name "db_*.sql.gz" -mtime +30 -delete
EOF

    chmod +x /usr/local/bin/backup-scentora-db.sh

    # Add to crontab (2 AM daily)
    (crontab -l 2>/dev/null; echo "0 2 * * * /usr/local/bin/backup-scentora-db.sh") | crontab -

    log_success "Daily backups configured (2 AM)"
}

print_summary() {
    echo ""
    echo "========================================================================"
    log_success "Scentora deployment setup complete!"
    echo "========================================================================"
    echo ""
    echo "Domain: https://$DOMAIN"
    echo "Application Directory: $APP_DIR"
    echo "Credentials File: $CREDENTIALS_FILE"
    echo ""
    echo "Next steps:"
    echo "  1. Copy your docker-compose.prod.yml to $APP_DIR"
    echo "  2. Deploy your application:"
    echo "     cd $APP_DIR"
    echo "     docker compose -f docker-compose.prod.yml pull"
    echo "     docker compose -f docker-compose.prod.yml up -d"
    echo ""
    echo "  3. Check container status:"
    echo "     docker compose -f docker-compose.prod.yml ps"
    echo ""
    echo "  4. View logs:"
    echo "     docker compose -f docker-compose.prod.yml logs -f"
    echo ""
    echo "  5. Test the application:"
    echo "     curl https://$DOMAIN/api/health"
    echo ""
    echo "========================================================================"
    echo ""

    cat "$CREDENTIALS_FILE"
}

###############################################################################
# Main Installation Flow
###############################################################################

main() {
    log_info "Starting Scentora droplet setup..."

    check_root
    check_args

    # System setup
    install_docker
    configure_firewall
    install_nginx
    install_certbot
    setup_automatic_updates

    # Application setup
    setup_application_directory
    generate_secrets
    create_env_file
    save_credentials
    configure_nginx_proxy
    obtain_ssl_certificate

    # Docker setup
    clone_or_create_compose
    start_containers

    # Maintenance setup
    setup_backup_cron

    # Summary
    print_summary
}

# Run main function
main "$@"
