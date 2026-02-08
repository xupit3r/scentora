#!/bin/bash
set -euo pipefail

###############################################################################
# Scentora Bootstrap - Create first user and invitation on the production server
# Must be run from a machine with SSH access to the droplet.
# Uses ssh -t so the container's interactive password prompt works.
###############################################################################

HOST="${SCENTORA_HOST:-root@scentora.thejoeshow.net}"
APP_DIR="/opt/scentora"

read -rp "Email: " EMAIL
read -rp "Username: " USERNAME

if [ -z "$EMAIL" ] || [ -z "$USERNAME" ]; then
  echo "Email and username are required."
  exit 1
fi

# ssh -t allocates a TTY so the interactive password prompt works inside the container
ssh -t "$HOST" "cd $APP_DIR && docker compose -f docker-compose.prod.yml --env-file .env.production \
  exec backend node dist/cli.js bootstrap --email '$EMAIL' --username '$USERNAME'"
