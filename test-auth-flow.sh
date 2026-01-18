#!/bin/bash

echo "=== Scentora Auth Flow Test ==="
echo ""

API_BASE="http://localhost:3000/api"

echo "1. Testing health endpoint..."
curl -s $API_BASE/health | jq . || echo "Backend not running"
echo ""

echo "2. Creating test user (register)..."
REGISTER_RESPONSE=$(curl -s -X POST $API_BASE/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "username": "testuser",
    "password": "password123"
  }')

if echo "$REGISTER_RESPONSE" | jq -e '.accessToken' > /dev/null 2>&1; then
  echo "✓ Registration successful"
  ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.accessToken')
  REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.refreshToken')
  echo "Access Token: ${ACCESS_TOKEN:0:50}..."
  echo "Refresh Token: ${REFRESH_TOKEN:0:50}..."
else
  echo "✗ Registration failed (user might already exist)"
  echo "$REGISTER_RESPONSE" | jq .
  
  echo ""
  echo "3. Trying to login instead..."
  LOGIN_RESPONSE=$(curl -s -X POST $API_BASE/auth/login \
    -H "Content-Type: application/json" \
    -d '{
      "email": "testuser@example.com",
      "password": "password123"
    }')
  
  ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.accessToken')
  REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refreshToken')
  echo "Access Token: ${ACCESS_TOKEN:0:50}..."
  echo "Refresh Token: ${REFRESH_TOKEN:0:50}..."
fi

echo ""
echo "4. Testing authenticated endpoint (GET /perfumes)..."
PERFUMES_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" $API_BASE/perfumes \
  -H "Authorization: Bearer $ACCESS_TOKEN")

HTTP_CODE=$(echo "$PERFUMES_RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$PERFUMES_RESPONSE" | sed '/HTTP_CODE:/d')

if [ "$HTTP_CODE" = "200" ]; then
  echo "✓ Authenticated request successful (HTTP $HTTP_CODE)"
  echo "$BODY" | jq .
else
  echo "✗ Authenticated request failed (HTTP $HTTP_CODE)"
  echo "$BODY" | jq .
fi

echo ""
echo "5. Testing /me endpoint..."
ME_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" $API_BASE/auth/me \
  -H "Authorization: Bearer $ACCESS_TOKEN")

HTTP_CODE=$(echo "$ME_RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$ME_RESPONSE" | sed '/HTTP_CODE:/d')

if [ "$HTTP_CODE" = "200" ]; then
  echo "✓ /me request successful (HTTP $HTTP_CODE)"
  echo "$BODY" | jq .
else
  echo "✗ /me request failed (HTTP $HTTP_CODE)"
  echo "$BODY" | jq .
fi

echo ""
echo "6. Testing token refresh..."
REFRESH_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST $API_BASE/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refreshToken\": \"$REFRESH_TOKEN\"}")

HTTP_CODE=$(echo "$REFRESH_RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$REFRESH_RESPONSE" | sed '/HTTP_CODE:/d')

if [ "$HTTP_CODE" = "200" ]; then
  echo "✓ Token refresh successful (HTTP $HTTP_CODE)"
  NEW_ACCESS_TOKEN=$(echo "$BODY" | jq -r '.accessToken')
  echo "New Access Token: ${NEW_ACCESS_TOKEN:0:50}..."
else
  echo "✗ Token refresh failed (HTTP $HTTP_CODE)"
  echo "$BODY" | jq .
fi

echo ""
echo "=== Test Complete ==="
