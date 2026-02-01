# Recipe System API Documentation

**Version**: 1.0  
**Last Updated**: February 1, 2026  
**Base URL**: `http://localhost:3000/api`

---

## Table of Contents

1. [Authentication](#authentication)
2. [Recipe Endpoints](#recipe-endpoints)
3. [Recipe Version Endpoints](#recipe-version-endpoints)
4. [Recipe Ingredient Endpoints](#recipe-ingredient-endpoints)
5. [Recipe Note Endpoints](#recipe-note-endpoints)
6. [Recipe Tag Endpoints](#recipe-tag-endpoints)
7. [Recipe Collection Endpoints](#recipe-collection-endpoints)
8. [Error Responses](#error-responses)
9. [Examples](#examples)

---

## Authentication

All recipe endpoints require JWT authentication.

**Header**:
```
Authorization: Bearer <access_token>
```

**Getting a Token**:
```bash
POST /api/auth/register
POST /api/auth/login
```

---

## Recipe Endpoints

### Create Recipe

**Endpoint**: `POST /api/recipes`  
**Auth**: Required  
**Description**: Create a new recipe

**Request Body**:
```json
{
  "name": "Summer Citrus Blend",
  "description": "A bright and refreshing citrus-forward fragrance",
  "notes": "Top: Bergamot, Lemon\nHeart: Neroli\nBase: Amber",
  "targetVolumeMl": 100
}
```

**Validation**:
- `name`: Required, max 255 chars, must be unique per user
- `description`: Optional
- `notes`: Optional
- `targetVolumeMl`: Required, must be > 0

**Response**: `201 Created`
```json
{
  "_id": "uuid",
  "userId": "uuid",
  "name": "Summer Citrus Blend",
  "description": "A bright and refreshing citrus-forward fragrance",
  "targetVolumeMl": 100,
  "status": "draft",
  "createdAt": "2026-02-01T12:00:00Z",
  "updatedAt": "2026-02-01T12:00:00Z"
}
```

---

### List Recipes

**Endpoint**: `GET /api/recipes`  
**Auth**: Required  
**Description**: List all recipes for the authenticated user

**Query Parameters**:
- `status` (optional): Filter by status (draft, in_progress, tested, finalized, archived)
- `limit` (optional): Number of results (default: 50)
- `offset` (optional): Pagination offset (default: 0)

**Example**: `GET /api/recipes?status=draft&limit=10`

**Response**: `200 OK`
```json
[
  {
    "_id": "uuid",
    "userId": "uuid",
    "name": "Summer Citrus Blend",
    "description": "A bright and refreshing citrus-forward fragrance",
    "targetVolumeMl": 100,
    "status": "draft",
    "createdAt": "2026-02-01T12:00:00Z",
    "updatedAt": "2026-02-01T12:00:00Z"
  }
]
```

---

### Get Recipe

**Endpoint**: `GET /api/recipes/:id`  
**Auth**: Required  
**Description**: Get a specific recipe by ID

**Response**: `200 OK`
```json
{
  "_id": "uuid",
  "userId": "uuid",
  "name": "Summer Citrus Blend",
  "description": "A bright and refreshing citrus-forward fragrance",
  "targetVolumeMl": 100,
  "status": "draft",
  "createdAt": "2026-02-01T12:00:00Z",
  "updatedAt": "2026-02-01T12:00:00Z"
}
```

**Errors**:
- `404`: Recipe not found or doesn't belong to user

---

### Update Recipe

**Endpoint**: `PUT /api/recipes/:id`  
**Auth**: Required  
**Description**: Update recipe details

**Request Body** (all fields optional):
```json
{
  "name": "Summer Citrus Blend v2",
  "description": "Updated description",
  "notes": "Updated notes",
  "status": "in_progress"
}
```

**Valid Status Values**:
- `draft`
- `in_progress`
- `tested`
- `finalized`
- `archived`

**Response**: `200 OK`
```json
{
  "_id": "uuid",
  "userId": "uuid",
  "name": "Summer Citrus Blend v2",
  "description": "Updated description",
  "targetVolumeMl": 100,
  "status": "in_progress",
  "createdAt": "2026-02-01T12:00:00Z",
  "updatedAt": "2026-02-01T12:05:00Z"
}
```

---

### Delete Recipe

**Endpoint**: `DELETE /api/recipes/:id`  
**Auth**: Required  
**Description**: Delete a recipe (cascades to versions, ingredients, notes, tags)

**Response**: `200 OK` or `204 No Content`

---

### Search Recipes

**Endpoint**: `GET /api/recipes/search`  
**Auth**: Required  
**Description**: Full-text search across recipe names, descriptions, and notes

**Query Parameters**:
- `q`: Search query (required)

**Example**: `GET /api/recipes/search?q=citrus`

**Response**: `200 OK`
```json
[
  {
    "_id": "uuid",
    "userId": "uuid",
    "name": "Summer Citrus Blend",
    "description": "A bright and refreshing citrus-forward fragrance",
    "targetVolumeMl": 100,
    "status": "draft",
    "createdAt": "2026-02-01T12:00:00Z",
    "updatedAt": "2026-02-01T12:00:00Z"
  }
]
```

---

## Recipe Version Endpoints

### Create Version

**Endpoint**: `POST /api/recipes/:id/versions`  
**Auth**: Required  
**Description**: Create a new version of a recipe (auto-activates as active version)

**Request Body**:
```json
{
  "notes": "Adjusted bergamot ratio for better balance",
  "changes": "Increased bergamot from 15% to 20%"
}
```

**Response**: `201 Created`
```json
{
  "_id": "uuid",
  "recipeId": "uuid",
  "versionNumber": 2,
  "name": "v2",
  "notes": "Adjusted bergamot ratio for better balance",
  "isActive": true,
  "createdAt": "2026-02-01T12:10:00Z"
}
```

**Note**: New versions automatically become the active version

---

### List Versions

**Endpoint**: `GET /api/recipes/:id/versions`  
**Auth**: Required  
**Description**: List all versions of a recipe

**Response**: `200 OK`
```json
[
  {
    "_id": "uuid",
    "recipeId": "uuid",
    "versionNumber": 2,
    "name": "v2",
    "notes": "Adjusted bergamot ratio",
    "isActive": true,
    "createdAt": "2026-02-01T12:10:00Z"
  },
  {
    "_id": "uuid",
    "recipeId": "uuid",
    "versionNumber": 1,
    "name": "v1",
    "notes": "Initial version",
    "isActive": false,
    "createdAt": "2026-02-01T12:00:00Z"
  }
]
```

---

### Get Version

**Endpoint**: `GET /api/recipes/:id/versions/:versionId`  
**Auth**: Required  
**Description**: Get a specific version by ID

**Response**: `200 OK`
```json
{
  "_id": "uuid",
  "recipeId": "uuid",
  "versionNumber": 2,
  "name": "v2",
  "notes": "Adjusted bergamot ratio",
  "isActive": true,
  "createdAt": "2026-02-01T12:10:00Z"
}
```

---

### Activate Version

**Endpoint**: `POST /api/recipes/:id/versions/:versionNumber/activate`  
**Auth**: Required  
**Description**: Set a specific version as the active version

**Response**: `200 OK`
```json
{
  "_id": "uuid",
  "recipeId": "uuid",
  "versionNumber": 1,
  "name": "v1",
  "notes": "Initial version",
  "isActive": true,
  "createdAt": "2026-02-01T12:00:00Z"
}
```

---

### Duplicate Version

**Endpoint**: `POST /api/recipes/:id/versions/:versionId/duplicate`  
**Auth**: Required  
**Description**: Create a copy of an existing version

**Response**: `201 Created`
```json
{
  "_id": "uuid",
  "recipeId": "uuid",
  "versionNumber": 3,
  "name": "v3",
  "notes": "Copy of v2",
  "isActive": true,
  "createdAt": "2026-02-01T12:15:00Z"
}
```

---

## Recipe Ingredient Endpoints

### Add Ingredient

**Endpoint**: `POST /api/recipes/:id/versions/:versionId/ingredients`  
**Auth**: Required  
**Description**: Add an accord to a recipe version

**Request Body**:
```json
{
  "accordId": "uuid",
  "quantityMl": 15.5
}
```

**Validation**:
- `accordId`: Required, must be a valid accord belonging to the user
- `quantityMl`: Required, must be > 0

**Response**: `201 Created`
```json
{
  "_id": "uuid",
  "versionId": "uuid",
  "accordId": "uuid",
  "quantityMl": 15.5,
  "createdAt": "2026-02-01T12:20:00Z"
}
```

**Note**: If user has `validate_recipe_volumes` enabled, validates against accord inventory

---

### Update Ingredient

**Endpoint**: `PUT /api/recipes/:id/versions/:versionId/ingredients/:ingredientId`  
**Auth**: Required  
**Description**: Update ingredient quantity

**Request Body**:
```json
{
  "quantityMl": 20.0
}
```

**Response**: `200 OK`
```json
{
  "_id": "uuid",
  "versionId": "uuid",
  "accordId": "uuid",
  "quantityMl": 20.0,
  "createdAt": "2026-02-01T12:20:00Z"
}
```

---

### Remove Ingredient

**Endpoint**: `DELETE /api/recipes/:id/versions/:versionId/ingredients/:ingredientId`  
**Auth**: Required  
**Description**: Remove an ingredient from a version

**Response**: `200 OK` or `204 No Content`

---

## Recipe Note Endpoints

### Create Note

**Endpoint**: `POST /api/recipes/:id/notes`  
**Auth**: Required  
**Description**: Add a note to a recipe

**Request Body**:
```json
{
  "content": "This formula performs exceptionally well in an alcohol base. Consider testing in oil base next.",
  "type": "observation"
}
```

**Note Types**:
- `general` (default)
- `observation`
- `adjustment`
- `reminder`

**Response**: `201 Created`
```json
{
  "_id": "uuid",
  "recipeId": "uuid",
  "content": "This formula performs exceptionally well in an alcohol base.",
  "noteType": "observation",
  "createdAt": "2026-02-01T12:25:00Z",
  "updatedAt": "2026-02-01T12:25:00Z"
}
```

---

### List Notes

**Endpoint**: `GET /api/recipes/:id/notes`  
**Auth**: Required  
**Description**: Get all notes for a recipe

**Response**: `200 OK`
```json
[
  {
    "_id": "uuid",
    "recipeId": "uuid",
    "content": "This formula performs exceptionally well in an alcohol base.",
    "noteType": "observation",
    "createdAt": "2026-02-01T12:25:00Z",
    "updatedAt": "2026-02-01T12:25:00Z"
  }
]
```

---

### Update Note

**Endpoint**: `PUT /api/recipes/:id/notes/:noteId`  
**Auth**: Required  
**Description**: Update note content

**Request Body**:
```json
{
  "content": "Updated note content"
}
```

**Response**: `200 OK`

---

### Delete Note

**Endpoint**: `DELETE /api/recipes/:id/notes/:noteId`  
**Auth**: Required  
**Description**: Delete a note

**Response**: `200 OK` or `204 No Content`

---

## Recipe Tag Endpoints

### Add Tag

**Endpoint**: `POST /api/recipes/:id/tags`  
**Auth**: Required  
**Description**: Add a tag to a recipe

**Request Body**:
```json
{
  "tag": "citrus"
}
```

**Response**: `201 Created`
```json
{
  "tag": "citrus"
}
```

---

### Get Tags

**Endpoint**: `GET /api/recipes/:id/tags`  
**Auth**: Required  
**Description**: Get all tags for a recipe

**Response**: `200 OK`
```json
[
  {"tag": "citrus"},
  {"tag": "fresh"},
  {"tag": "summer"}
]
```

---

### Remove Tag

**Endpoint**: `DELETE /api/recipes/:id/tags/:tag`  
**Auth**: Required  
**Description**: Remove a tag from a recipe

**Response**: `200 OK` or `204 No Content`

---

## Recipe Collection Endpoints

### Create Collection

**Endpoint**: `POST /api/collections`  
**Auth**: Required  
**Description**: Create a new recipe collection

**Request Body**:
```json
{
  "name": "Summer Scents",
  "description": "Light and fresh recipes perfect for warm weather"
}
```

**Response**: `201 Created`
```json
{
  "_id": "uuid",
  "userId": "uuid",
  "name": "Summer Scents",
  "description": "Light and fresh recipes perfect for warm weather",
  "createdAt": "2026-02-01T12:30:00Z",
  "updatedAt": "2026-02-01T12:30:00Z"
}
```

---

### List Collections

**Endpoint**: `GET /api/collections`  
**Auth**: Required  
**Description**: List all collections

**Query Parameters**:
- `limit` (optional): Number of results (default: 50)
- `offset` (optional): Pagination offset (default: 0)

**Response**: `200 OK`
```json
[
  {
    "_id": "uuid",
    "userId": "uuid",
    "name": "Summer Scents",
    "description": "Light and fresh recipes",
    "createdAt": "2026-02-01T12:30:00Z",
    "updatedAt": "2026-02-01T12:30:00Z"
  }
]
```

---

### Get Collection

**Endpoint**: `GET /api/collections/:id`  
**Auth**: Required  
**Description**: Get a specific collection

**Response**: `200 OK`
```json
{
  "_id": "uuid",
  "userId": "uuid",
  "name": "Summer Scents",
  "description": "Light and fresh recipes",
  "createdAt": "2026-02-01T12:30:00Z",
  "updatedAt": "2026-02-01T12:30:00Z"
}
```

---

### Update Collection

**Endpoint**: `PUT /api/collections/:id`  
**Auth**: Required  
**Description**: Update collection details

**Request Body**:
```json
{
  "name": "Summer Scents 2026",
  "description": "Updated description"
}
```

**Response**: `200 OK`

---

### Delete Collection

**Endpoint**: `DELETE /api/collections/:id`  
**Auth**: Required  
**Description**: Delete a collection (does not delete recipes)

**Response**: `200 OK` or `204 No Content`

---

### Add Recipe to Collection

**Endpoint**: `POST /api/collections/:id/recipes`  
**Auth**: Required  
**Description**: Add a recipe to a collection

**Request Body**:
```json
{
  "recipeId": "uuid"
}
```

**Response**: `200 OK`

---

### Remove Recipe from Collection

**Endpoint**: `DELETE /api/collections/:id/recipes/:recipeId`  
**Auth**: Required  
**Description**: Remove a recipe from a collection

**Response**: `200 OK` or `204 No Content`

---

## Error Responses

### Standard Error Format

```json
{
  "error": "Error message"
}
```

Or with details:

```json
{
  "error": {
    "message": "Validation failed",
    "details": "Field 'name' is required"
  }
}
```

### HTTP Status Codes

- `200 OK`: Success
- `201 Created`: Resource created successfully
- `204 No Content`: Success with no response body
- `400 Bad Request`: Validation error or invalid input
- `401 Unauthorized`: Missing or invalid JWT token
- `404 Not Found`: Resource not found or doesn't belong to user
- `409 Conflict`: Constraint violation (e.g., duplicate name)
- `500 Internal Server Error`: Server error

---

## Examples

### Complete Workflow Example

```bash
# 1. Register/Login
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.accessToken')

# 2. Create a recipe
RECIPE=$(curl -s -X POST http://localhost:3000/api/recipes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Lavender Dreams",
    "description":"Calming lavender blend",
    "targetVolumeMl":100
  }')

RECIPE_ID=$(echo $RECIPE | jq -r '._id')

# 3. Create a version
VERSION=$(curl -s -X POST http://localhost:3000/api/recipes/$RECIPE_ID/versions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "notes":"Initial formula",
    "changes":"First version"
  }')

VERSION_ID=$(echo $VERSION | jq -r '._id')

# 4. Add note
curl -s -X POST http://localhost:3000/api/recipes/$RECIPE_ID/notes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content":"Great performance in tests",
    "type":"observation"
  }'

# 5. Add tag
curl -s -X POST http://localhost:3000/api/recipes/$RECIPE_ID/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"tag":"floral"}'

# 6. Create collection
COLLECTION=$(curl -s -X POST http://localhost:3000/api/collections \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Favorites",
    "description":"My best recipes"
  }')

COLL_ID=$(echo $COLLECTION | jq -r '._id')

# 7. Add recipe to collection
curl -s -X POST http://localhost:3000/api/collections/$COLL_ID/recipes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"recipeId\":\"$RECIPE_ID\"}"

# 8. Search recipes
curl -s -X GET "http://localhost:3000/api/recipes/search?q=lavender" \
  -H "Authorization: Bearer $TOKEN"

# 9. Update recipe status
curl -s -X PUT http://localhost:3000/api/recipes/$RECIPE_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"tested"}'

# 10. List all recipes
curl -s -X GET http://localhost:3000/api/recipes \
  -H "Authorization: Bearer $TOKEN"
```

---

## Notes

- All timestamps are in ISO 8601 format (UTC)
- UUIDs are used for all resource IDs
- Pagination uses limit/offset pattern
- All endpoints enforce user isolation (recipes only visible to owner)
- Cascade deletes: Deleting a recipe deletes all versions, ingredients, notes, and tags

---

**End of API Documentation**
