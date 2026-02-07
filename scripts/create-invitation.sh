#!/bin/bash
set -euo pipefail

API_URL="${SCENTORA_API_URL:-https://scentora.thejoeshow.net/api}"

usage() {
  echo "Usage: $0 [-e email] [-d days]"
  echo ""
  echo "Options:"
  echo "  -e EMAIL   Restrict invitation to a specific email"
  echo "  -d DAYS    Expiration in days (default: 7, max: 365)"
  echo ""
  echo "Environment:"
  echo "  SCENTORA_API_URL   API base URL (default: https://scentora.thejoeshow.net/api)"
  exit 1
}

EMAIL=""
DAYS=""

while getopts "e:d:h" opt; do
  case $opt in
    e) EMAIL="$OPTARG" ;;
    d) DAYS="$OPTARG" ;;
    h) usage ;;
    *) usage ;;
  esac
done

# Prompt for credentials
read -rp "Email: " LOGIN_EMAIL
read -rsp "Password: " LOGIN_PASSWORD
echo ""

# Login
LOGIN_RESPONSE=$(curl -sf "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$LOGIN_EMAIL\", \"password\": \"$LOGIN_PASSWORD\"}" 2>&1) || {
  echo "Login failed. Check your credentials."
  exit 1
}

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.accessToken')
if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
  echo "Login failed: could not extract token."
  exit 1
fi

# Build invitation request body
BODY="{}"
if [ -n "$EMAIL" ] && [ -n "$DAYS" ]; then
  BODY="{\"email\": \"$EMAIL\", \"expiresInDays\": $DAYS}"
elif [ -n "$EMAIL" ]; then
  BODY="{\"email\": \"$EMAIL\"}"
elif [ -n "$DAYS" ]; then
  BODY="{\"expiresInDays\": $DAYS}"
fi

# Create invitation
INVITE_RESPONSE=$(curl -sf "$API_URL/invitations" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d "$BODY" 2>&1) || {
  echo "Failed to create invitation."
  exit 1
}

CODE=$(echo "$INVITE_RESPONSE" | jq -r '.invitation.code')
EXPIRES=$(echo "$INVITE_RESPONSE" | jq -r '.invitation.expiresAt')

echo ""
echo "Invitation created!"
echo "  Code:    $CODE"
[ -n "$EMAIL" ] && echo "  For:     $EMAIL"
echo "  Expires: $EXPIRES"
