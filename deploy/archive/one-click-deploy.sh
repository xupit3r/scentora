#!/bin/bash
# Scentora One-Click Deployment for Digital Ocean
# Run this script on your local machine

set -e

echo "╔════════════════════════════════════════════╗"
echo "║   🌸 Scentora One-Click Deployment 🌸     ║"
echo "╚════════════════════════════════════════════╝"
echo ""

# Get user input
read -p "Enter your Digital Ocean droplet IP: " SERVER_IP
read -p "Enter SSH username (default: root): " SSH_USER
SSH_USER=${SSH_USER:-root}
read -p "Enter your domain name: " DOMAIN
read -p "Enter your email for SSL: " EMAIL

echo ""
echo "📦 Preparing deployment..."

# Create temporary deployment script
TEMP_SCRIPT=$(mktemp)

cat > $TEMP_SCRIPT << 'SERVERSCRIPT'
#!/bin/bash
set -e

# Get configuration from arguments
DOMAIN=$1
EMAIL=$2

echo "🔧 Installing dependencies..."
apt-get update -qq
apt-get install -y curl git docker.io nginx certbot python3-certbot-nginx ufw -qq

curl -fsSL https://deb.nodesource.com/setup_20.x | bash - > /dev/null
apt-get install -y nodejs -qq

curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

npm install -g pm2 -qq
pm2 startup systemd -u root --hp /root > /dev/null

echo "✓ Dependencies installed"

# Generate secrets
JWT_SECRET=$(openssl rand -base64 32)
COUCH_PASSWORD=$(openssl rand -base64 16)

echo "🔐 Configuring application..."

# Create env file
cat > /var/www/scentora/backend/.env << ENVEOF
NODE_ENV=production
PORT=3000
COUCHDB_URL=http://localhost:5984
COUCHDB_USER=admin
COUCHDB_PASSWORD=$COUCH_PASSWORD
COUCHDB_DATABASE=scentora
JWT_SECRET=$JWT_SECRET
JWT_ACCESS_EXPIRES_IN=15m
JWT_REFRESH_EXPIRES_IN=7d
ENVEOF

# Create docker-compose for CouchDB
cat > /var/www/scentora/docker-compose.prod.yml << DOCKEREOF
version: '3.8'
services:
  couchdb:
    image: couchdb:3.3
    container_name: scentora-couchdb
    restart: always
    environment:
      - COUCHDB_USER=admin
      - COUCHDB_PASSWORD=$COUCH_PASSWORD
    ports:
      - "127.0.0.1:5984:5984"
    volumes:
      - couchdb-data:/opt/couchdb/data
volumes:
  couchdb-data:
DOCKEREOF

# Start CouchDB
docker-compose -f /var/www/scentora/docker-compose.prod.yml up -d
sleep 5

echo "🏗️  Building application..."

# Build backend
cd /var/www/scentora/backend
npm install --production -qq
npm run build

# Build frontend with production API URL
cd /var/www/scentora/frontend
echo "VITE_API_URL=https://$DOMAIN/api" > .env.production
npm install -qq
npm run build

# Create PM2 config
cat > /var/www/scentora/ecosystem.config.js << PM2EOF
module.exports = {
  apps: [{
    name: 'scentora-api',
    cwd: '/var/www/scentora/backend',
    script: './dist/index.js',
    instances: 2,
    exec_mode: 'cluster',
    env: { NODE_ENV: 'production', PORT: 3000 }
  }]
};
PM2EOF

mkdir -p /var/log/scentora
pm2 start /var/www/scentora/ecosystem.config.js
pm2 save

echo "🌐 Configuring web server..."

# Configure Nginx
cat > /etc/nginx/sites-available/scentora << NGINXEOF
server {
    listen 80;
    server_name $DOMAIN;
    root /var/www/scentora/frontend/dist;
    index index.html;

    location /api {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location / {
        try_files \$uri \$uri/ /index.html;
    }
}
NGINXEOF

ln -sf /etc/nginx/sites-available/scentora /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl reload nginx

# Setup SSL
certbot --nginx -d $DOMAIN --non-interactive --agree-tos -m $EMAIL

# Configure firewall
ufw --force enable
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow http
ufw allow https

# Save credentials
cat > /root/scentora-credentials.txt << CREDEOF
Scentora Credentials
====================
Domain: $DOMAIN
CouchDB Password: $COUCH_PASSWORD
JWT Secret: $JWT_SECRET
Generated: $(date)
CREDEOF
chmod 600 /root/scentora-credentials.txt

echo ""
echo "✅ Deployment complete!"
echo "🌐 Visit: https://$DOMAIN"
echo "🔐 Credentials saved to: /root/scentora-credentials.txt"
SERVERSCRIPT

echo "📤 Uploading files to server..."
rsync -avz --exclude='node_modules' --exclude='.git' --exclude='logs' --exclude='dist' \
  $(dirname $(dirname $0))/ ${SSH_USER}@${SERVER_IP}:/var/www/scentora/ > /dev/null

echo "🚀 Running deployment on server..."
ssh ${SSH_USER}@${SERVER_IP} "bash -s $DOMAIN $EMAIL" < $TEMP_SCRIPT

rm $TEMP_SCRIPT

echo ""
echo "╔════════════════════════════════════════════╗"
echo "║          🎉 Deployment Complete! 🎉        ║"
echo "╚════════════════════════════════════════════╝"
echo ""
echo "Your application is now live at: https://$DOMAIN"
echo ""
echo "Next steps:"
echo "1. Visit your domain and register an account"
echo "2. SSH into server to view credentials: ssh $SSH_USER@$SERVER_IP"
echo "3. View credentials: cat /root/scentora-credentials.txt"
echo ""
