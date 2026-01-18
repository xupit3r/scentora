# Scentora Digital Ocean Deployment Guide

Complete guide to deploying Scentora on Digital Ocean with automated setup.

## Quick Start

### 1. Create Digital Ocean Droplet

1. Go to [Digital Ocean](https://cloud.digitalocean.com)
2. Create a new droplet:
   - **Image**: Ubuntu 22.04 LTS x64
   - **Size**: Basic ($12/mo - 2 GB RAM, 2 vCPUs)
   - **Datacenter**: Choose closest to your users
   - **Authentication**: SSH keys (recommended) or password
3. Wait for droplet to be created
4. Note the IP address

### 2. Point Your Domain to the Server

1. Go to your domain registrar
2. Create an A record:
   - **Name**: @ (or subdomain like `app`)
   - **Value**: Your droplet's IP address
   - **TTL**: 3600
3. Wait 5-10 minutes for DNS propagation

### 3. Deploy the Application

From your local computer:

```bash
cd /path/to/scentora

# Upload files to server
scp -r * root@YOUR_SERVER_IP:/var/www/scentora/

# SSH into server
ssh root@YOUR_SERVER_IP

# Run setup script
cd /var/www/scentora/deploy
chmod +x setup-server.sh
sudo ./setup-server.sh
```

The script will ask you for:
- Domain name (e.g., scentora.example.com)
- Email for SSL certificate
- JWT secret (auto-generated if left empty)
- CouchDB password (auto-generated if left empty)

That's it! Your application will be live at `https://yourdomain.com`

##  What Gets Installed

The setup script automatically installs and configures:

### Software
- ✅ Node.js 20.x
- ✅ npm (latest)
- ✅ Docker & Docker Compose
- ✅ Nginx (web server)
- ✅ Certbot (SSL certificates)
- ✅ PM2 (process manager)
- ✅ UFW (firewall)

### Services
- ✅ CouchDB (in Docker container)
- ✅ Backend API (via PM2, 2 instances)
- ✅ Frontend (served by Nginx)
- ✅ SSL Certificate (Let's Encrypt)

### Security
- ✅ Firewall configured (SSH, HTTP, HTTPS only)
- ✅ HTTPS enforced
- ✅ Security headers
- ✅ Rate limiting
- ✅ Auto-renewal for SSL certificates

### Automation
- ✅ Daily database backups (2 AM)
- ✅ Auto-start services on reboot
- ✅ Management commands

## Management Commands

Once deployed, use these commands on the server:

```bash
# Start all services
scentora start

# Stop all services
scentora stop

# Restart all services
scentora restart

# Check status
scentora status

# View API logs
scentora logs api

# View Database logs
scentora logs db

# Manual backup
scentora backup

# Update from git
scentora update
```

## File Locations

```
/var/www/scentora/          # Application files
├── backend/                # Backend API
├── frontend/               # Frontend app
├── docker-compose.prod.yml # CouchDB config
└── ecosystem.config.js     # PM2 config

/var/log/scentora/          # Application logs
├── api-error.log          # Backend errors
└── api-out.log            # Backend output

/var/backups/scentora/      # Database backups
└── couchdb_*.json.gz      # Daily backups (kept 7 days)

/root/scentora-credentials.txt # Saved credentials
```

## Accessing Services

### Frontend
```
https://yourdomain.com
```

### CouchDB Admin (SSH tunnel required)
```bash
# On your local computer:
ssh -L 5984:localhost:5984 root@YOUR_SERVER_IP

# Then visit: http://localhost:5984/_utils
# Username: admin
# Password: (from /root/scentora-credentials.txt)
```

### PM2 Dashboard
```bash
pm2 monit
```

## Updating the Application

### Method 1: Using Git (Recommended)

1. Push your changes to Git
2. SSH into server
3. Run: `scentora update`

### Method 2: Manual Upload

```bash
# From your local computer
rsync -avz --exclude='node_modules' ./ root@YOUR_SERVER_IP:/var/www/scentora/

# SSH into server
ssh root@YOUR_SERVER_IP

# Rebuild and restart
cd /var/www/scentora/backend && npm install && npm run build
cd /var/www/scentora/frontend && npm install && npm run build
scentora restart
```

## Monitoring

### Check Service Status
```bash
scentora status
```

### View Live Logs
```bash
# Backend
pm2 logs scentora-api

# All PM2 processes
pm2 monit
```

### Check Disk Space
```bash
df -h
```

### Check Memory Usage
```bash
free -h
```

## Backups

### Automatic Backups
- Run daily at 2 AM
- Stored in `/var/backups/scentora/`
- Kept for 7 days
- Compressed with gzip

### Manual Backup
```bash
scentora backup
```

### Restore from Backup
```bash
# List backups
ls -lh /var/backups/scentora/

# Restore (replace DATE with your backup date)
gunzip /var/backups/scentora/couchdb_DATE.json.gz
curl -X POST http://admin:PASSWORD@localhost:5984/scentora/_bulk_docs \
  -H "Content-Type: application/json" \
  -d @/var/backups/scentora/couchdb_DATE.json
```

## Troubleshooting

### Application Won't Start

```bash
# Check PM2 status
pm2 status

# Check PM2 logs
pm2 logs

# Restart PM2
pm2 restart all
```

### Database Connection Issues

```bash
# Check if CouchDB is running
docker ps

# Check CouchDB logs
docker logs scentora-couchdb

# Restart CouchDB
docker restart scentora-couchdb
```

### SSL Certificate Issues

```bash
# Check certificate status
certbot certificates

# Renew manually
certbot renew

# Renew for specific domain
certbot --nginx -d yourdomain.com
```

### Nginx Issues

```bash
# Test configuration
nginx -t

# Check status
systemctl status nginx

# Restart
systemctl restart nginx

# View error logs
tail -f /var/log/nginx/error.log
```

### Out of Disk Space

```bash
# Check disk usage
df -h

# Clean old logs
pm2 flush

# Clean old backups
find /var/backups/scentora -name "*.gz" -mtime +7 -delete

# Clean Docker
docker system prune -a
```

## Security Best Practices

### Change Default Passwords
```bash
# Edit backend .env
nano /var/www/scentora/backend/.env

# Update CouchDB password in docker-compose
nano /var/www/scentora/docker-compose.prod.yml

# Restart services
scentora restart
```

### Enable Automatic Security Updates
```bash
apt-get install unattended-upgrades
dpkg-reconfigure --priority=low unattended-upgrades
```

### Setup Fail2Ban (Brute Force Protection)
```bash
apt-get install fail2ban
systemctl enable fail2ban
systemctl start fail2ban
```

### Monitor Login Attempts
```bash
# View failed SSH attempts
grep "Failed password" /var/log/auth.log

# View successful logins
last -a
```

## Performance Optimization

### Enable Redis for Rate Limiting (Optional)
```bash
# Install Redis
apt-get install redis-server

# Update backend to use Redis
# (Requires code changes - see advanced docs)
```

### Setup CDN (Optional)
1. Sign up for Cloudflare (free tier)
2. Point DNS to Cloudflare
3. Enable caching and minification
4. Configure origin server

### Database Optimization
```bash
# Compact database
curl -X POST http://admin:PASSWORD@localhost:5984/scentora/_compact

# View database info
curl http://admin:PASSWORD@localhost:5984/scentora
```

## Scaling

### Vertical Scaling (Bigger Droplet)
1. Create snapshot of current droplet
2. Resize droplet to larger size
3. Restart services

### Horizontal Scaling (Multiple Servers)
1. Setup load balancer
2. Deploy to multiple droplets
3. Use external CouchDB cluster
4. Configure PM2 cluster mode

## Cost Estimates

### Minimum Setup
- **Droplet**: $12/month (2GB RAM)
- **Backups**: Included in droplet
- **SSL**: Free (Let's Encrypt)
- **Domain**: $10-15/year
- **Total**: ~$13/month + domain

### Recommended Setup
- **Droplet**: $24/month (4GB RAM)
- **Backups**: $2.40/month (20% of droplet)
- **Spaces (S3-like)**: $5/month (for backups)
- **Total**: ~$32/month + domain

## Support & Resources

- **Digital Ocean Docs**: https://docs.digitalocean.com
- **CouchDB Docs**: https://docs.couchdb.org
- **PM2 Docs**: https://pm2.keymetrics.io
- **Nginx Docs**: https://nginx.org/en/docs
- **Certbot Docs**: https://certbot.eff.org

## Next Steps

After deployment:
1. ✅ Register your first user account
2. ✅ Test all features work
3. ✅ Setup database backups to external storage
4. ✅ Configure monitoring/alerts
5. ✅ Document your custom configuration
6. ✅ Plan for regular updates

## Advanced Configuration

See separate guides for:
- Custom domain configuration
- Multiple environment setup (staging/prod)
- CI/CD pipeline setup
- Database replication
- Performance tuning
- Security hardening

## License

MIT
