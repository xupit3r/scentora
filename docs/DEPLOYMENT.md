# Scentora Deployment Guide

Complete guide for deploying Scentora to DigitalOcean with Docker, PostgreSQL, and automated CI/CD.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Manual Deployment](#manual-deployment)
- [GitHub Actions Setup](#github-actions-setup)
- [Environment Variables](#environment-variables)
- [Database Management](#database-management)
- [SSL Certificate Setup](#ssl-certificate-setup)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [Backup & Restore](#backup--restore)
- [Troubleshooting](#troubleshooting)
- [Rollback Procedures](#rollback-procedures)
- [Security Hardening](#security-hardening)

## Prerequisites

### Required

- **DigitalOcean Account**: Create at https://www.digitalocean.com
- **Domain Name**: Point A record to your droplet's IP address
- **SSH Key**: Generate with `ssh-keygen` and add to DigitalOcean
- **Docker Hub Account**: For storing container images (or use GitHub Container Registry)

### Recommended Droplet Specs

- **Minimum**: $12/month (2GB RAM, 1 CPU, 50GB SSD)
- **Recommended**: $24/month (4GB RAM, 2 CPU, 80GB SSD)
- **Operating System**: Ubuntu 24.04 LTS

## Quick Start

### Option 1: Automated Setup (Recommended)

1. **Create a DigitalOcean droplet** with Ubuntu 24.04
2. **Point your domain's A record** to the droplet's IP
3. **Wait for DNS propagation** (5-10 minutes)
4. **Run the setup script** on your droplet:

```bash
ssh root@YOUR_DROPLET_IP
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/scentora/main/deploy/setup-droplet.sh | bash -s your-domain.com your@email.com
```

5. **Deploy your application files**:

```bash
# On your local machine
rsync -avz --exclude 'node_modules' --exclude '.git' . root@YOUR_DROPLET_IP:/opt/scentora/

# On the droplet
cd /opt/scentora
docker compose -f docker-compose.prod.yml up -d
```

6. **Verify deployment**:

```bash
curl https://your-domain.com/api/health
```

### Option 2: GitHub Actions CI/CD

After manual setup once, use GitHub Actions for automated deployments on every push to `main`.

See [GitHub Actions Setup](#github-actions-setup) below.

## Manual Deployment

### Step 1: Create DigitalOcean Droplet

1. Log in to DigitalOcean
2. Create > Droplets
3. Choose Ubuntu 24.04 LTS
4. Select droplet size (2GB RAM minimum)
5. Add your SSH key
6. Create Droplet

### Step 2: Configure DNS

Point your domain's A record to the droplet's IP:

```
Type: A
Name: @ (or subdomain)
Data: YOUR_DROPLET_IP
TTL: 3600
```

Wait 5-10 minutes for DNS propagation.

### Step 3: Initial Server Setup

SSH into your droplet:

```bash
ssh root@YOUR_DROPLET_IP
```

Update packages:

```bash
apt-get update && apt-get upgrade -y
```

### Step 4: Run Setup Script

```bash
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/scentora/main/deploy/setup-droplet.sh | bash -s your-domain.com your@email.com
```

This script will:
- Install Docker and Docker Compose
- Configure UFW firewall
- Install nginx and certbot
- Generate secure secrets
- Setup SSL certificates
- Configure automatic backups

### Step 5: Deploy Application

Copy your application files:

```bash
# From your local machine
rsync -avz --exclude 'node_modules' --exclude '.git' . root@YOUR_DROPLET_IP:/opt/scentora/
```

Start containers:

```bash
ssh root@YOUR_DROPLET_IP
cd /opt/scentora
docker compose -f docker-compose.prod.yml --env-file .env.production up -d
```

### Step 6: Verify Deployment

```bash
# Check container status
docker compose -f docker-compose.prod.yml ps

# Check logs
docker compose -f docker-compose.prod.yml logs -f

# Test API
curl https://your-domain.com/api/health

# Test frontend
open https://your-domain.com
```

## GitHub Actions Setup

### Step 1: Add GitHub Secrets

Go to your repository > Settings > Secrets and variables > Actions

Add these secrets:

| Secret Name | Description | Example/Command |
|------------|-------------|-----------------|
| `DOCKER_USERNAME` | Docker Hub username | `yourusername` |
| `DOCKER_PASSWORD` | Docker Hub password or token | Generate at hub.docker.com |
| `DEPLOY_HOST` | Droplet IP address | `192.168.1.100` |
| `DEPLOY_USER` | SSH user | `root` |
| `DEPLOY_KEY` | SSH private key | `cat ~/.ssh/id_rsa` |
| `DOMAIN` | Your domain name | `scentora.example.com` |
| `JWT_SECRET` | JWT signing secret | `openssl rand -hex 32` |
| `DB_PASSWORD` | PostgreSQL password | `openssl rand -base64 24` |

### Step 2: Generate SSH Key for Deployment

On your local machine:

```bash
ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/github_actions
```

Add the public key to your droplet:

```bash
ssh-copy-id -i ~/.ssh/github_actions.pub root@YOUR_DROPLET_IP
```

Copy the private key to GitHub Secrets:

```bash
cat ~/.ssh/github_actions
# Copy output to DEPLOY_KEY secret
```

### Step 3: Test Deployment

Push to main branch:

```bash
git push origin main
```

Monitor the deployment:
- GitHub > Actions tab
- Watch the workflow run
- Check for any errors

## Environment Variables

### Production Environment Variables

Located at `/opt/scentora/.env.production` on the droplet:

```bash
# Domain Configuration
DOMAIN=your-domain.com

# Database Configuration
DB_USER=admin
DB_PASSWORD=<generated>
DB_NAME=scentora

# Backend Configuration
BACKEND_PORT=3000
JWT_SECRET=<generated>
JWT_EXPIRES_IN=7d
JWT_REFRESH_EXPIRES_IN=30d

# Frontend Configuration
VITE_API_URL=/api

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://your-domain.com

# Docker Configuration
DOCKER_REGISTRY=scentora
IMAGE_TAG=latest
```

### Generating Secrets

```bash
# JWT Secret (64 characters)
openssl rand -hex 32

# Database Password (32 characters)
openssl rand -base64 24 | tr -d "=+/" | cut -c1-32
```

## Database Management

### Prisma Migrations

**Development** (create new migration):

```bash
cd backend
npx prisma migrate dev --name add_new_field
git add prisma/migrations/
git commit -m "feat: add new field"
```

**Production** (automatic on container startup):

```bash
# Migrations run automatically via docker-compose command:
npx prisma migrate deploy && node dist/index.js
```

**Manual migration**:

```bash
docker compose -f docker-compose.prod.yml exec backend npx prisma migrate deploy
```

### Check Migration Status

```bash
docker compose -f docker-compose.prod.yml exec backend npx prisma migrate status
```

### Direct Database Access

```bash
# Access PostgreSQL
docker compose -f docker-compose.prod.yml exec postgres psql -U admin -d scentora

# Run SQL query
docker compose -f docker-compose.prod.yml exec postgres psql -U admin -d scentora -c "SELECT COUNT(*) FROM users;"
```

## SSL Certificate Setup

### Automatic Setup

SSL certificates are configured automatically by the `setup-droplet.sh` script using Let's Encrypt and certbot.

### Manual Certificate Renewal

Certificates auto-renew via systemd timer. To manually renew:

```bash
certbot renew
systemctl reload nginx
```

### Check Certificate Status

```bash
certbot certificates
```

### Test SSL Configuration

Visit https://www.ssllabs.com/ssltest/ and enter your domain.

Target: **A+ rating**

## Monitoring & Maintenance

### Container Health

```bash
# Check status
docker compose -f docker-compose.prod.yml ps

# View logs
docker compose -f docker-compose.prod.yml logs -f

# View specific service logs
docker compose -f docker-compose.prod.yml logs -f backend

# Container stats
docker stats --no-stream
```

### Application Health

```bash
# API health check
curl https://your-domain.com/api/health

# Expected response:
# {"status":"ok","timestamp":"2024-02-06T12:00:00.000Z","service":"scentora-api"}
```

### System Resources

```bash
# Disk usage
df -h

# Memory usage
free -h

# Docker disk usage
docker system df
```

### Log Management

```bash
# View recent logs
docker compose -f docker-compose.prod.yml logs --tail=50 backend

# Export logs
docker compose -f docker-compose.prod.yml logs > logs_$(date +%Y%m%d).txt

# Clear old logs (if needed)
docker compose -f docker-compose.prod.yml logs --since 24h > recent.log
```

## Backup & Restore

### Automated Daily Backups

Backups run daily at 2 AM via cron (configured by setup script):

```bash
/usr/local/bin/backup-scentora-db.sh
```

Backups stored in: `/var/backups/scentora/`

Retention: 30 days

### Manual Backup

```bash
# Backup database
docker compose -f docker-compose.prod.yml exec -T postgres pg_dump -U admin scentora > backup_$(date +%Y%m%d).sql

# Backup with compression
docker compose -f docker-compose.prod.yml exec -T postgres pg_dump -U admin scentora | gzip > backup_$(date +%Y%m%d).sql.gz
```

### Restore from Backup

```bash
# Stop backend (optional, but recommended)
docker compose -f docker-compose.prod.yml stop backend

# Restore database
docker compose -f docker-compose.prod.yml exec -T postgres psql -U admin scentora < backup_20240206.sql

# Or restore from compressed backup
gunzip -c backup_20240206.sql.gz | docker compose -f docker-compose.prod.yml exec -T postgres psql -U admin scentora

# Restart backend
docker compose -f docker-compose.prod.yml start backend
```

### Off-site Backup

```bash
# Copy backups to remote storage (recommended)
rsync -avz /var/backups/scentora/ user@backup-server:/backups/scentora/

# Or use rclone for S3/Spaces/etc
rclone copy /var/backups/scentora/ spaces:scentora-backups/
```

## Troubleshooting

### Containers Won't Start

```bash
# Check logs
docker compose -f docker-compose.prod.yml logs

# Check specific service
docker compose -f docker-compose.prod.yml logs backend

# Check environment variables
docker compose -f docker-compose.prod.yml config

# Restart containers
docker compose -f docker-compose.prod.yml restart
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker compose -f docker-compose.prod.yml ps postgres

# Check PostgreSQL logs
docker compose -f docker-compose.prod.yml logs postgres

# Verify DATABASE_URL
docker compose -f docker-compose.prod.yml exec backend env | grep DATABASE_URL

# Test database connection
docker compose -f docker-compose.prod.yml exec postgres pg_isready -U admin
```

### SSL Certificate Issues

```bash
# Check certificate
certbot certificates

# Renew certificate
certbot renew

# Check nginx configuration
nginx -t

# Reload nginx
systemctl reload nginx
```

### Frontend Not Loading

```bash
# Check frontend container
docker compose -f docker-compose.prod.yml logs frontend

# Check nginx
systemctl status nginx
nginx -t

# Check nginx access logs
tail -f /var/log/nginx/access.log

# Check nginx error logs
tail -f /var/log/nginx/error.log
```

### High Memory Usage

```bash
# Check container stats
docker stats

# Restart specific container
docker compose -f docker-compose.prod.yml restart backend

# Check for memory leaks in logs
docker compose -f docker-compose.prod.yml logs backend | grep -i "memory\|heap"
```

### Disk Space Issues

```bash
# Check disk usage
df -h

# Clean up Docker
docker system prune -a

# Clean up old backups
find /var/backups/scentora -name "db_*.sql.gz" -mtime +30 -delete

# Clean up logs
journalctl --vacuum-time=7d
```

## Rollback Procedures

### Git-Based Rollback

```bash
# Revert last commit
git revert HEAD
git push origin main
# GitHub Actions will automatically deploy the previous version
```

### Docker Image Rollback

```bash
# SSH to droplet
ssh root@YOUR_DROPLET_IP
cd /opt/scentora

# List available images
docker images | grep scentora

# Edit docker-compose.prod.yml to use specific tag
nano docker-compose.prod.yml
# Change: image: username/scentora-backend:latest
# To:     image: username/scentora-backend:sha-abc1234

# Pull and restart
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d

# Verify
curl https://your-domain.com/api/health
```

### Database Rollback

```bash
# Restore from backup
docker compose -f docker-compose.prod.yml stop backend
gunzip -c /var/backups/scentora/db_20240206.sql.gz | docker compose -f docker-compose.prod.yml exec -T postgres psql -U admin scentora
docker compose -f docker-compose.prod.yml start backend
```

### Emergency Full Rollback

```bash
# Stop everything
docker compose -f docker-compose.prod.yml down

# Restore from backup
docker volume rm scentora-postgres-data
docker volume create scentora-postgres-data

# Start PostgreSQL
docker compose -f docker-compose.prod.yml up -d postgres

# Wait for PostgreSQL to be ready
sleep 10

# Restore database
gunzip -c backup_20240206.sql.gz | docker compose -f docker-compose.prod.yml exec -T postgres psql -U admin scentora

# Start all services
docker compose -f docker-compose.prod.yml up -d
```

## Security Hardening

### Firewall Status

```bash
# Check firewall
ufw status verbose

# Expected output:
# Status: active
# To                         Action      From
# --                         ------      ----
# 22/tcp                     ALLOW       Anywhere
# 80/tcp                     ALLOW       Anywhere
# 443/tcp                    ALLOW       Anywhere
```

### SSH Security

Disable password authentication (after adding SSH key):

```bash
# Edit SSH config
nano /etc/ssh/sshd_config

# Set these values:
PasswordAuthentication no
PermitRootLogin prohibit-password
PubkeyAuthentication yes

# Restart SSH
systemctl restart sshd
```

### Regular Updates

```bash
# Manual system updates
apt-get update && apt-get upgrade -y

# Automatic updates are configured by setup script
```

### Security Audit

```bash
# Check open ports
netstat -tlnp

# Check running processes
ps aux | grep docker

# Review cron jobs
crontab -l

# Check for failed login attempts
grep "Failed password" /var/log/auth.log | tail -20
```

### Secrets Management

- Never commit `.env.production` to git
- Store credentials in `/root/.scentora-credentials` (chmod 600)
- Use GitHub Secrets for CI/CD
- Rotate secrets periodically (every 90 days recommended)

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Prisma Documentation](https://www.prisma.io/docs)
- [DigitalOcean Tutorials](https://www.digitalocean.com/community/tutorials)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [nginx Documentation](https://nginx.org/en/docs/)

## Support

For issues or questions:
- Check the [Troubleshooting](#troubleshooting) section above
- Review container logs
- Check GitHub Issues

---

**Last Updated**: February 2024
**Version**: 1.0.0
