#!/bin/bash
# Phase 10.5 - Manual Integration Testing
# Comprehensive test of all Recipe API endpoints

BASE_URL="http://localhost:3000/api"
LOG_FILE="/tmp/phase10-5-test.log"

echo "🧪 Phase 10.5 - Recipe API Integration Testing" | tee $LOG_FILE
echo "=================================================" | tee -a $LOG_FILE
echo "" | tee -a $LOG_FILE

# Login
echo "1️⃣  Logging in..." | tee -a $LOG_FILE
LOGIN_RESP=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"demo@scentora.com","password":"demo1234"}')

TOKEN=$(echo "$LOGIN_RESP" | grep -oP '"accessToken":"\K[^"]+')
if [ -z "$TOKEN" ]; then
    echo "❌ Login failed!" | tee -a $LOG_FILE
    echo "$LOGIN_RESP" | tee -a $LOG_FILE
    exit 1
fi
echo "✅ Login successful (token: ${TOKEN:0:30}...)" | tee -a $LOG_FILE
echo "" | tee -a $LOG_FILE

# Create Accord (needed for ingredients)
echo "2️⃣  Creating test accord..." | tee -a $LOG_FILE
ACCORD_RESP=$(curl -s -X POST "$BASE_URL/accords" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Phase105 Test Accord","pyramidPosition":"top","volumeMl":50}')

ACCORD_ID=$(echo "$ACCORD_RESP" | grep -oP '"_id":"\K[^"]+' | head -1)
if [ -z "$ACCORD_ID" ]; then
    echo "❌ Accord creation failed!" | tee -a $LOG_FILE
    echo "$ACCORD_RESP" | tee -a $LOG_FILE
    exit 1
fi
echo "✅ Accord created: $ACCORD_ID" | tee -a $LOG_FILE
echo "" | tee -a $LOG_FILE

# CREATE RECIPE
echo "3️⃣  Creating recipe..." | tee -a $LOG_FILE
RECIPE_RESP=$(curl -s -X POST "$BASE_URL/recipes" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Phase 10.5 Integration Test","description":"Full API test","targetVolumeMl":100,"status":"draft"}')

RECIPE_ID=$(echo "$RECIPE_RESP" | grep -oP '"_id":"\K[^"]+' | head -1)
if [ -z "$RECIPE_ID" ]; then
    echo "❌ Recipe creation failed!" | tee -a $LOG_FILE
    echo "$RECIPE_RESP" | tee -a $LOG_FILE
    exit 1
fi
echo "✅ Recipe created: $RECIPE_ID" | tee -a $LOG_FILE
echo "" | tee -a $LOG_FILE

# GET RECIPE
echo "4️⃣  Getting recipe details..." | tee -a $LOG_FILE
GET_RESP=$(curl -s -X GET "$BASE_URL/recipes/$RECIPE_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$GET_RESP" | grep -q "$RECIPE_ID"; then
    echo "✅ Recipe retrieved successfully" | tee -a $LOG_FILE
else
    echo "❌ Recipe retrieval failed!" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# LIST RECIPES
echo "5️⃣  Listing all recipes..." | tee -a $LOG_FILE
LIST_RESP=$(curl -s -X GET "$BASE_URL/recipes" \
    -H "Authorization: Bearer $TOKEN")

if echo "$LIST_RESP" | grep -q "$RECIPE_ID"; then
    echo "✅ Recipe list includes our recipe" | tee -a $LOG_FILE
else
    echo "❌ Recipe not in list!" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# UPDATE RECIPE
echo "6️⃣  Updating recipe..." | tee -a $LOG_FILE
UPDATE_RESP=$(curl -s -X PUT "$BASE_URL/recipes/$RECIPE_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Phase 10.5 UPDATED","description":"Updated desc","targetVolumeMl":150,"status":"in_progress"}')

if echo "$UPDATE_RESP" | grep -q "UPDATED"; then
    echo "✅ Recipe updated successfully" | tee -a $LOG_FILE
else
    echo "❌ Recipe update failed!" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# SEARCH RECIPES
echo "7️⃣  Searching recipes..." | tee -a $LOG_FILE
SEARCH_RESP=$(curl -s -X GET "$BASE_URL/recipes/search?query=Phase" \
    -H "Authorization: Bearer $TOKEN")

if echo "$SEARCH_RESP" | grep -q "$RECIPE_ID"; then
    echo "✅ Recipe search works" | tee -a $LOG_FILE
else
    echo "⚠️  Recipe search returned no results (might be OK)" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# CREATE VERSION
echo "8️⃣  Creating recipe version..." | tee -a $LOG_FILE
VERSION_RESP=$(curl -s -X POST "$BASE_URL/recipes/$RECIPE_ID/versions" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"notes":"Version 1 - Initial formulation"}')

VERSION_ID=$(echo "$VERSION_RESP" | grep -oP '"_id":"\K[^"]+' | head -1)
if [ -z "$VERSION_ID" ]; then
    echo "❌ Version creation failed!" | tee -a $LOG_FILE
    echo "$VERSION_RESP" | tee -a $LOG_FILE
else
    echo "✅ Version created: $VERSION_ID" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# LIST VERSIONS
echo "9️⃣  Listing versions..." | tee -a $LOG_FILE
VERSION_LIST=$(curl -s -X GET "$BASE_URL/recipes/$RECIPE_ID/versions" \
    -H "Authorization: Bearer $TOKEN")

if echo "$VERSION_LIST" | grep -q "$VERSION_ID"; then
    echo "✅ Version list retrieved" | tee -a $LOG_FILE
else
    echo "❌ Version not in list!" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# ADD INGREDIENT
echo "🔟 Adding ingredient to version..." | tee -a $LOG_FILE
INGREDIENT_RESP=$(curl -s -X POST "$BASE_URL/recipes/$RECIPE_ID/versions/$VERSION_ID/ingredients" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"accordId\":\"$ACCORD_ID\",\"volumeMl\":10,\"percentage\":10}")

INGREDIENT_ID=$(echo "$INGREDIENT_RESP" | grep -oP '"_id":"\K[^"]+' | head -1)
if [ -z "$INGREDIENT_ID" ]; then
    echo "❌ Ingredient creation failed!" | tee -a $LOG_FILE
    echo "$INGREDIENT_RESP" | tee -a $LOG_FILE
else
    echo "✅ Ingredient added: $INGREDIENT_ID" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# UPDATE INGREDIENT
echo "1️⃣1️⃣ Updating ingredient..." | tee -a $LOG_FILE
INGREDIENT_UPDATE=$(curl -s -X PUT "$BASE_URL/recipes/$RECIPE_ID/versions/$VERSION_ID/ingredients/$INGREDIENT_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"volumeMl":15,"percentage":15}')

if echo "$INGREDIENT_UPDATE" | grep -q "15"; then
    echo "✅ Ingredient updated" | tee -a $LOG_FILE
else
    echo "⚠️  Ingredient update response: $INGREDIENT_UPDATE" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# CREATE NOTE
echo "1️⃣2️⃣ Creating recipe note..." | tee -a $LOG_FILE
NOTE_RESP=$(curl -s -X POST "$BASE_URL/recipes/$RECIPE_ID/notes" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"content":"Testing notes feature","noteType":"general"}')

NOTE_ID=$(echo "$NOTE_RESP" | grep -oP '"_id":"\K[^"]+' | head -1)
if [ -z "$NOTE_ID" ]; then
    echo "❌ Note creation failed!" | tee -a $LOG_FILE
    echo "$NOTE_RESP" | tee -a $LOG_FILE
else
    echo "✅ Note created: $NOTE_ID" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# LIST NOTES
echo "1️⃣3️⃣ Listing notes..." | tee -a $LOG_FILE
NOTE_LIST=$(curl -s -X GET "$BASE_URL/recipes/$RECIPE_ID/notes" \
    -H "Authorization: Bearer $TOKEN")

if echo "$NOTE_LIST" | grep -q "$NOTE_ID"; then
    echo "✅ Note list retrieved" | tee -a $LOG_FILE
else
    echo "❌ Note not in list!" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# UPDATE NOTE
echo "1️⃣4️⃣ Updating note..." | tee -a $LOG_FILE
NOTE_UPDATE=$(curl -s -X PUT "$BASE_URL/recipes/$RECIPE_ID/notes/$NOTE_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"content":"Updated note content!"}')

if echo "$NOTE_UPDATE" | grep -q "Updated note"; then
    echo "✅ Note updated" | tee -a $LOG_FILE
else
    echo "⚠️  Note update: check manually" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# ADD TAG
echo "1️⃣5️⃣ Adding tag to recipe..." | tee -a $LOG_FILE
TAG_RESP=$(curl -s -X POST "$BASE_URL/recipes/$RECIPE_ID/tags" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"tag":"citrus"}')

if echo "$TAG_RESP" | grep -q '"tags"'; then
    echo "✅ Tag added" | tee -a $LOG_FILE
else
    echo "⚠️  Tag add: check response" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# GET POPULAR TAGS
echo "1️⃣6️⃣ Getting popular tags..." | tee -a $LOG_FILE
POPULAR_TAGS=$(curl -s -X GET "$BASE_URL/recipes/tags/popular" \
    -H "Authorization: Bearer $TOKEN")

if echo "$POPULAR_TAGS" | grep -q '\['; then
    echo "✅ Popular tags retrieved" | tee -a $LOG_FILE
else
    echo "⚠️  Popular tags: no data yet" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# CREATE COLLECTION
echo "1️⃣7️⃣ Creating collection..." | tee -a $LOG_FILE
COLLECTION_RESP=$(curl -s -X POST "$BASE_URL/collections" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Phase 10.5 Test Collection","description":"Integration test collection"}')

COLLECTION_ID=$(echo "$COLLECTION_RESP" | grep -oP '"_id":"\K[^"]+' | head -1)
if [ -z "$COLLECTION_ID" ]; then
    echo "❌ Collection creation failed!" | tee -a $LOG_FILE
    echo "$COLLECTION_RESP" | tee -a $LOG_FILE
else
    echo "✅ Collection created: $COLLECTION_ID" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# ADD RECIPE TO COLLECTION
echo "1️⃣8️⃣ Adding recipe to collection..." | tee -a $LOG_FILE
ADD_TO_COLL=$(curl -s -X POST "$BASE_URL/collections/$COLLECTION_ID/recipes" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"recipeId\":\"$RECIPE_ID\"}")

if echo "$ADD_TO_COLL" | grep -q "$RECIPE_ID"; then
    echo "✅ Recipe added to collection" | tee -a $LOG_FILE
else
    echo "⚠️  Add to collection: check response" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# GET COLLECTION
echo "1️⃣9️⃣ Getting collection details..." | tee -a $LOG_FILE
COLL_GET=$(curl -s -X GET "$BASE_URL/collections/$COLLECTION_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$COLL_GET" | grep -q "$COLLECTION_ID"; then
    echo "✅ Collection retrieved" | tee -a $LOG_FILE
else
    echo "❌ Collection retrieval failed!" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# LIST COLLECTIONS
echo "2️⃣0️⃣ Listing collections..." | tee -a $LOG_FILE
COLL_LIST=$(curl -s -X GET "$BASE_URL/collections" \
    -H "Authorization: Bearer $TOKEN")

if echo "$COLL_LIST" | grep -q "$COLLECTION_ID"; then
    echo "✅ Collection list retrieved" | tee -a $LOG_FILE
else
    echo "❌ Collection not in list!" | tee -a $LOG_FILE
fi
echo "" | tee -a $LOG_FILE

# CLEANUP
echo "" | tee -a $LOG_FILE
echo "🧹 Cleanup (optional - leaving test data for manual verification)" | tee -a $LOG_FILE
echo "   To clean up manually:" | tee -a $LOG_FILE
echo "   - Recipe ID: $RECIPE_ID" | tee -a $LOG_FILE
echo "   - Collection ID: $COLLECTION_ID" | tee -a $LOG_FILE
echo "   - Accord ID: $ACCORD_ID" | tee -a $LOG_FILE
echo "" | tee -a $LOG_FILE

echo "=================================================" | tee -a $LOG_FILE
echo "✅ Phase 10.5 Integration Testing Complete!" | tee -a $LOG_FILE
echo "=================================================" | tee -a $LOG_FILE
echo "" | tee -a $LOG_FILE
echo "Full log saved to: $LOG_FILE" | tee -a $LOG_FILE
echo "" | tee -a $LOG_FILE

echo "Summary:" | tee -a $LOG_FILE
echo "- Created Recipe: $RECIPE_ID" | tee -a $LOG_FILE
echo "- Created Version: $VERSION_ID" | tee -a $LOG_FILE
echo "- Created Ingredient: $INGREDIENT_ID" | tee -a $LOG_FILE
echo "- Created Note: $NOTE_ID" | tee -a $LOG_FILE
echo "- Created Collection: $COLLECTION_ID" | tee -a $LOG_FILE
echo "- Created Accord: $ACCORD_ID" | tee -a $LOG_FILE
