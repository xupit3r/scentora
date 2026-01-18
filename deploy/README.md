# Scentora Digital Ocean Deployment

Automated deployment scripts and guides for deploying Scentora to Digital Ocean.

## 🚀 One-Click Deployment

The easiest way to deploy Scentora:

```bash
cd scentora/deploy
./one-click-deploy.sh
```

This single script will:
- ✅ Upload your application to the server
- ✅ Install all dependencies (Node.js, Docker, Nginx, etc.)
- ✅ Configure CouchDB
- ✅ Build backend and frontend
- ✅ Setup PM2 for process management
- ✅ Configure Nginx as reverse proxy
- ✅ Install SSL certificate (Let's Encrypt)
- ✅ Configure firewall
- ✅ Start all services

**Time**: ~10-15 minutes

## Prerequisites

### What You Need
1. **Digital Ocean Account** - [Sign up here](https://cloud.digitalocean.com)
2. **Ubuntu 22.04 Droplet** - Minimum 2GB RAM ($12/month)
3. **Domain Name** - Pointed to your droplet's IP
4. **SSH Access** - To your droplet

### Before You Start

1. **Create Droplet**:
   ```
   Image: Ubuntu 22.04 LTS
   Size: Basic - $12/mo (2 GB RAM, 2 vCPUs)
   Datacenter: Choose closest region
   Authentication: SSH key recommended
   ```

2. **Point DNS**:
   ```
   Type: A Record
   Name: @ (or subdomain)
   Value: YOUR_DROPLET_IP
   TTL: 3600
   ```

3. **Wait for DNS** (5-10 minutes)

## Deployment Methods

### Method 1: One-Click Script (Recommended)

```bash
cd scentora/deploy
./one-click-deploy.sh
```

You'll be asked for:
- Server IP address
- SSH username (default: root)
- Domain name
- Email for SSL certificate

**That's it!** Everything else is automatic.

### Method 2: Manual Deployment

If you prefer more control:

```bash
# 1. Upload files
scp -r * root@YOUR_SERVER_IP:/var/www/scentora/

# 2. SSH into server
ssh root@YOUR_SERVER_IP

# 3. Run setup
cd /var/www/scentora/deploy
chmod +x setup-server.sh
./setup-server.sh
```

## What Gets Deployed

### Architecture
```
┌─────────────────────────────────┐
│         Internet (HTTPS)        │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│         Nginx (Port 443)        │
│    - SSL Termination            │
│    - Static File Serving        │
│    - Reverse Proxy              │
└──────┬──────────────────┬───────┘
       │                  │
┌──────▼───────┐   ┌─────▼────────┐
│   Frontend   │   │   Backend    │
│  (Vue.js)    │   │  (Koa.js)    │
│  Static      │   │  PM2 x2      │
│  Files       │   │  Port 3000   │
└──────────────┘   └──────┬────────┘
                          │
                   ┌──────▼────────┐
                   │    CouchDB    │
                   │    Docker     │
                   │   Port 5984   │
                   └───────────────┘
```

### Installed Software
- Node.js 20.x
- Docker & Docker Compose
- Nginx
- Certbot (Let's Encrypt)
- PM2 (Process Manager)
- UFW (Firewall)

### Security Features
- ✅ HTTPS enforced
- ✅ Auto-renewing SSL certificates
- ✅ Firewall configured (SSH, HTTP, HTTPS only)
- ✅ Security headers
- ✅ Database not exposed to internet
- ✅ Environment variables secured

## Post-Deployment

### Accessing Your Application
```
https://yourdomain.com
```

### Server Management Commands

SSH into your server, then use:

```bash
# Check status
pm2 status                    # Backend status
docker ps                     # CouchDB status
systemctl status nginx        # Nginx status

# View logs
pm2 logs                      # Backend logs
docker logs scentora-couchdb  # Database logs
tail -f /var/log/nginx/error.log  # Nginx logs

# Restart services
pm2 restart all               # Restart backend
docker restart scentora-couchdb   # Restart database
systemctl restart nginx       # Restart Nginx

# View credentials
cat /root/scentora-credentials.txt
```

### Updating Your Application

```bash
# From your local machine
cd scentora/deploy
./one-click-deploy.sh
```

Or manually:
```bash
# SSH into server
ssh root@YOUR_SERVER_IP

# Pull latest code
cd /var/www/scentora
git pull

# Rebuild backend
cd backend
npm install
npm run build

# Rebuild frontend
cd ../frontend
npm install
npm run build

# Restart
pm2 restart all
```

## Monitoring

### Health Checks
```bash
# Check if backend is responding
curl http://localhost:3000/api/health

# Check if CouchDB is running
curl http://localhost:5984

# Check SSL certificate
curl -I https://yourdomain.com
```

### Performance Monitoring
```bash
# CPU and Memory
htop

# PM2 Dashboard
pm2 monit

# Disk Space
df -h

# Active Connections
netstat -an | grep :443 | wc -l
```

## Backups

### Automatic Backups
The deployment includes a daily backup script that:
- Runs at 2 AM daily
- Backs up CouchDB database
- Compresses backups
- Keeps last 7 days

Location: `/var/backups/scentora/`

### Manual Backup
```bash
# On the server
docker exec scentora-couchdb curl \
  http://admin:PASSWORD@localhost:5984/scentora/_all_docs?include_docs=true \
  > backup_$(date +%Y%m%d).json
```

### Restoring from Backup
```bash
# Upload docs back to CouchDB
curl -X POST http://admin:PASSWORD@localhost:5984/scentora/_bulk_docs \
  -H "Content-Type: application/json" \
  -d @backup_20240118.json
```

## Troubleshooting

### Application Not Accessible

```bash
# Check Nginx
systemctl status nginx
nginx -t

# Check SSL
certbot certificates

# Check firewall
ufw status
```

### Backend Not Responding

```bash
# Check PM2
pm2 status
pm2 logs

# Restart
pm2 restart all
```

### Database Connection Issues

```bash
# Check Docker
docker ps
docker logs scentora-couchdb

# Restart
docker restart scentora-couchdb
```

### SSL Certificate Issues

```bash
# Check certificate
certbot certificates

# Renew manually
certbot renew --dry-run
certbot renew
```

### Out of Disk Space

```bash
# Check space
df -h

# Clean Docker
docker system prune -a

# Clean PM2 logs
pm2 flush

# Clean old backups
find /var/backups/scentora -name "*.gz" -mtime +7 -delete
```

## Cost Breakdown

### Minimum Setup
- **Droplet**: $12/month (2GB RAM, 2 vCPUs)
- **Backups**: Included
- **SSL**: Free (Let's Encrypt)
- **Domain**: $10-15/year
- **Total**: ~$13/month

### Recommended Setup
- **Droplet**: $24/month (4GB RAM, 4 vCPUs)
- **Droplet Backups**: $4.80/month
- **Total**: ~$29/month

### Storage Needs
- Application: ~100MB
- Database: ~1GB (for thousands of perfumes)
- Logs: ~500MB/month
- Backups: ~1GB/week

## Scaling

### Vertical Scaling (More Resources)
1. Take snapshot of droplet
2. Resize to larger size in Digital Ocean
3. Restart services

### Horizontal Scaling (More Servers)
1. Setup load balancer
2. Deploy to multiple droplets
3. Use managed database
4. Configure session affinity

## Security Best Practices

### Change Default Passwords
```bash
# Generate new password
openssl rand -base64 32

# Update .env file
nano /var/www/scentora/backend/.env

# Restart services
pm2 restart all
```

### Enable Automatic Updates
```bash
apt-get install unattended-upgrades
dpkg-reconfigure --priority=low unattended-upgrades
```

### Setup Fail2Ban
```bash
apt-get install fail2ban
systemctl enable fail2ban
```

### Monitor Login Attempts
```bash
# Failed SSH attempts
grep "Failed password" /var/log/auth.log

# Successful logins
last -a
```

## Performance Tips

### Enable Gzip Compression
Already enabled in Nginx config for:
- Text files
- CSS
- JavaScript
- JSON
- XML

### Enable Browser Caching
Already configured for static assets (1 year)

### Database Optimization
```bash
# Compact database
curl -X POST http://admin:PASSWORD@localhost:5984/scentora/_compact
```

### PM2 Cluster Mode
Already enabled (2 instances)

## Files Created

```
/var/www/scentora/               # Application
├── backend/                     # Backend API
├── frontend/                    # Frontend app
├── docker-compose.prod.yml      # CouchDB config
└── ecosystem.config.js          # PM2 config

/etc/nginx/sites-available/scentora  # Nginx config
/var/log/scentora/               # App logs
/var/backups/scentora/           # Backups
/root/scentora-credentials.txt   # Credentials
```

## Support

For issues or questions:
1. Check the troubleshooting section
2. View server logs
3. Check Digital Ocean status page
4. Review application logs

## Advanced Topics

See DEPLOY_GUIDE.md for:
- Custom domain configuration
- CI/CD setup
- Database replication
- Multi-environment setup
- Performance tuning

## License

MIT
