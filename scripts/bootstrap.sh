#!/bin/bash
set -euo pipefail

###############################################################################
# Scentora Bootstrap - Create first user and invitation on the production server
# Must be run from a machine with SSH access to the droplet.
###############################################################################

HOST="${SCENTORA_HOST:-root@scentora.thejoeshow.net}"
APP_DIR="/opt/scentora"

echo "=== Scentora Bootstrap ==="
echo "This creates the first user and invitation directly on the server."
echo ""

read -rp "Email: " EMAIL
read -rp "Username: " USERNAME
read -rsp "Password: " PASSWORD
echo ""

if [ -z "$EMAIL" ] || [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
  echo "All fields are required."
  exit 1
fi

echo ""
echo "Creating user and invitation on $HOST..."

# Base64-encode credentials to safely pass through shell layers
B64_EMAIL=$(printf '%s' "$EMAIL" | base64)
B64_USERNAME=$(printf '%s' "$USERNAME" | base64)
B64_PASSWORD=$(printf '%s' "$PASSWORD" | base64)

# SSH into the droplet and run a Node.js script inside the backend container.
# Using Prisma + bcrypt avoids all shell escaping issues with passwords/SQL.
ssh "$HOST" "cd $APP_DIR && docker compose -f docker-compose.prod.yml --env-file .env.production exec -T \
  -e BOOT_EMAIL=$B64_EMAIL \
  -e BOOT_USERNAME=$B64_USERNAME \
  -e BOOT_PASSWORD=$B64_PASSWORD \
  backend node -e '
const bcrypt = await import(\"bcrypt\");
const crypto = await import(\"crypto\");
const { PrismaClient } = await import(\"@prisma/client\");

const decode = (b64) => Buffer.from(b64, \"base64\").toString();
const email = decode(process.env.BOOT_EMAIL);
const username = decode(process.env.BOOT_USERNAME);
const password = decode(process.env.BOOT_PASSWORD);

const prisma = new PrismaClient();
try {
  const hash = await bcrypt.hash(password, 10);
  const user = await prisma.user.create({
    data: { email, username, passwordHash: hash },
  });

  const code = crypto.randomBytes(16).toString(\"hex\");
  await prisma.invitation.create({
    data: {
      code,
      createdBy: user.id,
      expiresAt: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
    },
  });

  console.log(\"\");
  console.log(\"=== Bootstrap Complete ===\");
  console.log(\"  User:       \" + email + \" (\" + username + \")\");
  console.log(\"  Invitation: \" + code);
  console.log(\"  Expires:    30 days\");
  console.log(\"\");
  console.log(\"Log in with the user above, or register a new account\");
  console.log(\"at https://scentora.thejoeshow.net using the invitation code.\");
} finally {
  await prisma.\$disconnect();
}
'"
