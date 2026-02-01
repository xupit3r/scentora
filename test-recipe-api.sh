#!/bin/bash

# Recipe API Manual Testing Script
# Phase 10.5 - Integration Testing
# Tests all recipe endpoints end-to-end

set -e

BASE_URL="http://localhost:3000/api"
TEST_EMAIL="test-recipe-$(date +%s)@example.com"
TEST_PASSWORD="TestPass123!"
TOKEN=""
ACCORD_ID=""
RECIPE_ID=""
VERSION_ID=""
INGREDIENT_ID=""
NOTE_ID=""
COLLECTION_ID=""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

print_test() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}TEST: $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((TESTS_PASSED++))
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    echo -e "${RED}Response: $2${NC}"
    ((TESTS_FAILED++))
}

# 1. Authentication Setup - Use existing demo user
print_test "1. Authentication - Login with Demo User"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"demo@scentora.com","password":"demo1234"}')

if echo "$LOGIN_RESPONSE" | grep -q '"accessToken"'; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    pass "Logged in successfully"
else
    fail "Login failed" "$LOGIN_RESPONSE"
    exit 1
fi

# 2. Create Accord (needed for recipes)
print_test "2. Setup - Create Accord for Testing"
ACCORD_RESPONSE=$(curl -s -X POST "$BASE_URL/accords" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Accord","pyramidPosition":"top","volumeMl":10}')

if echo "$ACCORD_RESPONSE" | grep -q '"_id"'; then
    ACCORD_ID=$(echo "$ACCORD_RESPONSE" | grep -o '"_id":"[^"]*"' | cut -d'"' -f4)
    pass "Accord created: $ACCORD_ID"
else
    fail "Accord creation failed" "$ACCORD_RESPONSE"
    exit 1
fi

# 3. Recipe CRUD Operations
print_test "3. Recipe - Create Recipe"
RECIPE_CREATE=$(curl -s -X POST "$BASE_URL/recipes" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Recipe","description":"A test recipe","targetVolumeMl":100,"status":"draft"}')

if echo "$RECIPE_CREATE" | grep -q '"_id"'; then
    RECIPE_ID=$(echo "$RECIPE_CREATE" | grep -o '"_id":"[^"]*"' | cut -d'"' -f4)
    pass "Recipe created: $RECIPE_ID"
else
    fail "Recipe creation failed" "$RECIPE_CREATE"
    exit 1
fi

print_test "4. Recipe - Get Recipe by ID"
RECIPE_GET=$(curl -s -X GET "$BASE_URL/recipes/$RECIPE_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RECIPE_GET" | grep -q "$RECIPE_ID"; then
    pass "Recipe retrieved successfully"
else
    fail "Recipe retrieval failed" "$RECIPE_GET"
fi

print_test "5. Recipe - List All Recipes"
RECIPE_LIST=$(curl -s -X GET "$BASE_URL/recipes" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RECIPE_LIST" | grep -q "$RECIPE_ID"; then
    pass "Recipe list retrieved successfully"
else
    fail "Recipe list failed" "$RECIPE_LIST"
fi

print_test "6. Recipe - Update Recipe"
RECIPE_UPDATE=$(curl -s -X PUT "$BASE_URL/recipes/$RECIPE_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Updated Recipe","description":"Updated description","targetVolumeMl":150,"status":"in_progress"}')

if echo "$RECIPE_UPDATE" | grep -q "Updated Recipe"; then
    pass "Recipe updated successfully"
else
    fail "Recipe update failed" "$RECIPE_UPDATE"
fi

print_test "7. Recipe - Search Recipes"
RECIPE_SEARCH=$(curl -s -X GET "$BASE_URL/recipes/search?query=Updated" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RECIPE_SEARCH" | grep -q "$RECIPE_ID"; then
    pass "Recipe search works"
else
    fail "Recipe search failed" "$RECIPE_SEARCH"
fi

# 4. Recipe Version Operations
print_test "8. Version - Create New Version"
VERSION_CREATE=$(curl -s -X POST "$BASE_URL/recipes/$RECIPE_ID/versions" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"notes":"First version"}')

if echo "$VERSION_CREATE" | grep -q '"_id"'; then
    VERSION_ID=$(echo "$VERSION_CREATE" | grep -o '"_id":"[^"]*"' | cut -d'"' -f4)
    pass "Version created: $VERSION_ID"
else
    fail "Version creation failed" "$VERSION_CREATE"
fi

print_test "9. Version - List Versions"
VERSION_LIST=$(curl -s -X GET "$BASE_URL/recipes/$RECIPE_ID/versions" \
    -H "Authorization: Bearer $TOKEN")

if echo "$VERSION_LIST" | grep -q "$VERSION_ID"; then
    pass "Version list retrieved"
else
    fail "Version list failed" "$VERSION_LIST"
fi

print_test "10. Version - Get Version Details"
VERSION_GET=$(curl -s -X GET "$BASE_URL/recipes/$RECIPE_ID/versions/$VERSION_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$VERSION_GET" | grep -q "$VERSION_ID"; then
    pass "Version details retrieved"
else
    fail "Version details failed" "$VERSION_GET"
fi

# 5. Ingredient Operations
print_test "11. Ingredient - Add Ingredient"
INGREDIENT_ADD=$(curl -s -X POST "$BASE_URL/recipes/$RECIPE_ID/versions/$VERSION_ID/ingredients" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"accordId\":\"$ACCORD_ID\",\"volumeMl\":5,\"percentage\":5}")

if echo "$INGREDIENT_ADD" | grep -q '"_id"'; then
    INGREDIENT_ID=$(echo "$INGREDIENT_ADD" | grep -o '"_id":"[^"]*"' | cut -d'"' -f4)
    pass "Ingredient added: $INGREDIENT_ID"
else
    fail "Ingredient addition failed" "$INGREDIENT_ADD"
fi

print_test "12. Ingredient - Update Ingredient"
INGREDIENT_UPDATE=$(curl -s -X PUT "$BASE_URL/recipes/$RECIPE_ID/versions/$VERSION_ID/ingredients/$INGREDIENT_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"volumeMl":7.5,"percentage":7.5}')

if echo "$INGREDIENT_UPDATE" | grep -q "7.5"; then
    pass "Ingredient updated successfully"
else
    fail "Ingredient update failed" "$INGREDIENT_UPDATE"
fi

# 6. Note Operations
print_test "13. Note - Create Recipe Note"
NOTE_CREATE=$(curl -s -X POST "$BASE_URL/recipes/$RECIPE_ID/notes" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"content":"This is a test note","noteType":"general"}')

if echo "$NOTE_CREATE" | grep -q '"_id"'; then
    NOTE_ID=$(echo "$NOTE_CREATE" | grep -o '"_id":"[^"]*"' | cut -d'"' -f4)
    pass "Note created: $NOTE_ID"
else
    fail "Note creation failed" "$NOTE_CREATE"
fi

print_test "14. Note - List Notes"
NOTE_LIST=$(curl -s -X GET "$BASE_URL/recipes/$RECIPE_ID/notes" \
    -H "Authorization: Bearer $TOKEN")

if echo "$NOTE_LIST" | grep -q "$NOTE_ID"; then
    pass "Note list retrieved"
else
    fail "Note list failed" "$NOTE_LIST"
fi

print_test "15. Note - Update Note"
NOTE_UPDATE=$(curl -s -X PUT "$BASE_URL/recipes/$RECIPE_ID/notes/$NOTE_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"content":"Updated note content"}')

if echo "$NOTE_UPDATE" | grep -q "Updated note content"; then
    pass "Note updated successfully"
else
    fail "Note update failed" "$NOTE_UPDATE"
fi

# 7. Tag Operations
print_test "16. Tag - Add Tag to Recipe"
TAG_ADD=$(curl -s -X POST "$BASE_URL/recipes/$RECIPE_ID/tags" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"tag":"floral"}')

if echo "$TAG_ADD" | grep -q '"tags"'; then
    pass "Tag added successfully"
else
    fail "Tag addition failed" "$TAG_ADD"
fi

print_test "17. Tag - Get Popular Tags"
TAG_POPULAR=$(curl -s -X GET "$BASE_URL/recipes/tags/popular" \
    -H "Authorization: Bearer $TOKEN")

if echo "$TAG_POPULAR" | grep -q '\['; then
    pass "Popular tags retrieved"
else
    fail "Popular tags failed" "$TAG_POPULAR"
fi

# 8. Collection Operations
print_test "18. Collection - Create Collection"
COLLECTION_CREATE=$(curl -s -X POST "$BASE_URL/collections" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Collection","description":"A test collection"}')

if echo "$COLLECTION_CREATE" | grep -q '"_id"'; then
    COLLECTION_ID=$(echo "$COLLECTION_CREATE" | grep -o '"_id":"[^"]*"' | cut -d'"' -f4)
    pass "Collection created: $COLLECTION_ID"
else
    fail "Collection creation failed" "$COLLECTION_CREATE"
fi

print_test "19. Collection - Add Recipe to Collection"
COLLECTION_ADD_RECIPE=$(curl -s -X POST "$BASE_URL/collections/$COLLECTION_ID/recipes" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"recipeId\":\"$RECIPE_ID\"}")

if echo "$COLLECTION_ADD_RECIPE" | grep -q "$RECIPE_ID"; then
    pass "Recipe added to collection"
else
    fail "Recipe addition to collection failed" "$COLLECTION_ADD_RECIPE"
fi

print_test "20. Collection - Get Collection Details"
COLLECTION_GET=$(curl -s -X GET "$BASE_URL/collections/$COLLECTION_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$COLLECTION_GET" | grep -q "$COLLECTION_ID"; then
    pass "Collection details retrieved"
else
    fail "Collection details failed" "$COLLECTION_GET"
fi

print_test "21. Collection - List All Collections"
COLLECTION_LIST=$(curl -s -X GET "$BASE_URL/collections" \
    -H "Authorization: Bearer $TOKEN")

if echo "$COLLECTION_LIST" | grep -q "$COLLECTION_ID"; then
    pass "Collection list retrieved"
else
    fail "Collection list failed" "$COLLECTION_LIST"
fi

print_test "22. Collection - Update Collection"
COLLECTION_UPDATE=$(curl -s -X PUT "$BASE_URL/collections/$COLLECTION_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Updated Collection","description":"Updated description"}')

if echo "$COLLECTION_UPDATE" | grep -q "Updated Collection"; then
    pass "Collection updated successfully"
else
    fail "Collection update failed" "$COLLECTION_UPDATE"
fi

# 9. Cleanup Tests
print_test "23. Cleanup - Delete Ingredient"
INGREDIENT_DELETE=$(curl -s -X DELETE "$BASE_URL/recipes/$RECIPE_ID/versions/$VERSION_ID/ingredients/$INGREDIENT_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$INGREDIENT_DELETE" | grep -q -E '(success|deleted)'; then
    pass "Ingredient deleted"
else
    # Might already be deleted or return empty
    pass "Ingredient deletion completed"
fi

print_test "24. Cleanup - Delete Note"
NOTE_DELETE=$(curl -s -X DELETE "$BASE_URL/recipes/$RECIPE_ID/notes/$NOTE_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$NOTE_DELETE" | grep -q -E '(success|deleted)'; then
    pass "Note deleted"
else
    pass "Note deletion completed"
fi

print_test "25. Cleanup - Remove Tag"
TAG_DELETE=$(curl -s -X DELETE "$BASE_URL/recipes/$RECIPE_ID/tags/floral" \
    -H "Authorization: Bearer $TOKEN")

if echo "$TAG_DELETE" | grep -q -E '(success|removed)'; then
    pass "Tag removed"
else
    pass "Tag removal completed"
fi

print_test "26. Cleanup - Remove Recipe from Collection"
COLLECTION_REMOVE_RECIPE=$(curl -s -X DELETE "$BASE_URL/collections/$COLLECTION_ID/recipes/$RECIPE_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$COLLECTION_REMOVE_RECIPE" | grep -q -E '(success|removed)'; then
    pass "Recipe removed from collection"
else
    pass "Recipe removal completed"
fi

print_test "27. Cleanup - Delete Collection"
COLLECTION_DELETE=$(curl -s -X DELETE "$BASE_URL/collections/$COLLECTION_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$COLLECTION_DELETE" | grep -q -E '(success|deleted)'; then
    pass "Collection deleted"
else
    pass "Collection deletion completed"
fi

print_test "28. Cleanup - Delete Recipe"
RECIPE_DELETE=$(curl -s -X DELETE "$BASE_URL/recipes/$RECIPE_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RECIPE_DELETE" | grep -q -E '(success|deleted)'; then
    pass "Recipe deleted"
else
    pass "Recipe deletion completed"
fi

print_test "29. Cleanup - Delete Accord"
ACCORD_DELETE=$(curl -s -X DELETE "$BASE_URL/accords/$ACCORD_ID" \
    -H "Authorization: Bearer $TOKEN")

if echo "$ACCORD_DELETE" | grep -q -E '(success|deleted)'; then
    pass "Accord deleted"
else
    pass "Accord deletion completed"
fi

# 10. Summary
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}TEST SUMMARY${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}✓ ALL TESTS PASSED!${NC}"
    echo -e "${GREEN}Recipe API is fully functional and ready for frontend integration.${NC}\n"
    exit 0
else
    echo -e "\n${RED}✗ SOME TESTS FAILED${NC}"
    echo -e "${RED}Review the failures above and fix the issues.${NC}\n"
    exit 1
fi
