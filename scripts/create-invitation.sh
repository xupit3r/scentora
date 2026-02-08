#!/bin/bash
set -euo pipefail

###############################################################################
# Scentora Create Invitation - Create an invitation on the production server
# Must be run from a machine with SSH access to the droplet.
# No longer needs curl/jq/API auth — works directly against the DB.
###############################################################################

HOST="${SCENTORA_HOST:-root@scentora.thejoeshow.net}"
APP_DIR="/opt/scentora"

usage() {
  echo "Usage: $0 --email <creator-email> [--for-email <email>] [-d <days>]"
  echo ""
  echo "Options:"
  echo "  --email      Creator's email address (required)"
  echo "  --for-email  Restrict invitation to a specific email"
  echo "  -d, --days   Expiration in days (default: 7)"
  exit 1
}

# Pass all arguments through to the CLI tool
if [ $# -eq 0 ]; then
  usage
fi

ssh -t "$HOST" "cd $APP_DIR && docker compose -f docker-compose.prod.yml --env-file .env.production \
  exec -T backend node dist/cli.js create-invitation $*"
