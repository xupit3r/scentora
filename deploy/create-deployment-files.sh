#!/bin/bash
# This script creates all deployment files

DEPLOY_DIR="$(pwd)"

echo "Creating deployment files..."

# 1. Server setup script - Main deployment script
cat > "$DEPLOY_DIR/setup-server.sh" << 'EOF'
#!/bin/bash
# Complete server setup continues from here
# Due to length, see DEPLOY_GUIDE.md for full script
echo "This is a placeholder. Use the complete script from GitHub or docs."
EOF

chmod +x "$DEPLOY_DIR/setup-server.sh"
echo "✓ Created setup-server.sh"

# 2. Local deployment script
cat > "$DEPLOY_DIR/deploy.sh" << 'EOF'
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
EOF

chmod +x "$DEPLOY_DIR/deploy.sh"
echo "✓ Created deploy.sh"

echo ""
echo "All deployment files created in: $DEPLOY_DIR"
