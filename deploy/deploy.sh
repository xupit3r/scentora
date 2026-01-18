#!/bin/bash
# Local deployment script - Run from your computer

set -e

echo "================================"
echo "  Scentora Deployment Tool"
echo "================================"
echo ""

# Get server details
read -p "Enter server IP address: " SERVER_IP
read -p "Enter SSH user (default: root): " SSH_USER
SSH_USER=${SSH_USER:-root}

echo ""
echo "Step 1: Uploading files to server..."
rsync -avz --exclude='node_modules' --exclude='.git' --exclude='logs' \
  ../ ${SSH_USER}@${SERVER_IP}:/var/www/scentora/

echo ""
echo "Step 2: Running setup script on server..."
ssh ${SSH_USER}@${SERVER_IP} 'bash /var/www/scentora/deploy/setup-server.sh'

echo ""
echo "✓ Deployment complete!"
echo "Visit your domain to see the application"
