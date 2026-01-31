# API Specification

**Last Updated**: January 31, 2026  
**Version**: 2.0 (Accord System)  
**Base URL**: `http://localhost:3000/api`

---

## Overview

This document specifies all REST API endpoints for the Scentora accord management system.

---

## Authentication

All protected endpoints require JWT authentication via the `Authorization` header:

```
Authorization: Bearer <access_token>
```

**Token Types**:
- **Access Token**: Short-lived (15 minutes), used for API requests
- **Refresh Token**: Long-lived (7 days), used to obtain new access tokens

---

## Response Format

### Success Response
```json
{
  "success": true,
  "data": { /* response payload */ }
}
```

### Error Response
```json
{
  "error": {
    "message": "Error description",
    "details": "Additional details (optional)"
  }
}
```

### HTTP Status Codes
- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no response body
- `400 Bad Request` - Validation error
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate resource
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

---

## Endpoints

### Authentication (Public)

#### POST /auth/register
Register a new user account (requires invitation code).

**Request**:
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "invitationCode": "abc123def456"
}
```

**Response (201)**:
```json
{
  "user": {
    "_id": "uuid",
    "email": "user@example.com",
    "username": "username",
    "createdAt": "2026-01-31T10:00:00Z",
    "updatedAt": "2026-01-31T10:00:00Z"
  },
  "accessToken": "eyJhbGc...",
  "refreshToken": "abc123..."
}
```

#### POST /auth/login
Authenticate with email and password.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200)**:
```json
{
  "user": {
    "_id": "uuid",
    "email": "user@example.com",
    "username": "username"
  },
  "accessToken": "eyJhbGc...",
  "refreshToken": "abc123..."
}
```

#### POST /auth/refresh
Get new access token using refresh token.

**Request**:
```json
{
  "refreshToken": "abc123..."
}
```

**Response (200)**:
```json
{
  "accessToken": "eyJhbGc...",
  "refreshToken": "xyz789..." // New refresh token (rotation)
}
```

#### POST /auth/logout
Revoke a refresh token (logout from current device).

**Request**:
```json
{
  "refreshToken": "abc123..."
}
```

**Response (200)**:
```json
{
  "message": "Logged out successfully"
}
```

#### GET /auth/me
Get current user information.

**Auth**: Required

**Response (200)**:
```json
{
  "_id": "uuid",
  "email": "user@example.com",
  "username": "username",
  "createdAt": "2026-01-31T10:00:00Z",
  "updatedAt": "2026-01-31T10:00:00Z"
}
```

#### POST /auth/logout-all
Revoke all refresh tokens (logout from all devices).

**Auth**: Required

**Response (200)**:
```json
{
  "message": "Logged out from all devices"
}
```

---

### Invitations (Protected)

#### POST /invitations
Create a new invitation code.

**Auth**: Required

**Request**:
```json
{
  "email": "friend@example.com", // Optional: email-specific invitation
  "expiresInDays": 7              // Optional: default 7
}
```

**Response (201)**:
```json
{
  "invitation": {
    "_id": "uuid",
    "code": "abc123def456",
    "email": "friend@example.com",
    "createdBy": "user-uuid",
    "expiresAt": "2026-02-07T10:00:00Z",
    "used": false,
    "createdAt": "2026-01-31T10:00:00Z"
  }
}
```

#### GET /invitations
List all invitations created by current user.

**Auth**: Required

**Response (200)**:
```json
{
  "invitations": [
    {
      "_id": "uuid",
      "code": "abc123def456",
      "email": null,
      "createdBy": "user-uuid",
      "expiresAt": "2026-02-07T10:00:00Z",
      "used": true,
      "usedAt": "2026-02-01T10:00:00Z",
      "usedBy": "another-user-uuid",
      "createdAt": "2026-01-31T10:00:00Z"
    }
  ]
}
```

#### DELETE /invitations/:code
Revoke an invitation.

**Auth**: Required (must be creator)

**Response (200)**:
```json
{
  "message": "Invitation revoked successfully"
}
```

---

### Accords (Protected)

#### GET /accords
List all accords with optional filtering and search.

**Auth**: Required

**Query Parameters**:
- `search` (string): Search by name, notes, or supplier
- `position` (string): Filter by pyramid position (`top`, `middle`, `base`)
- `minVolume` (number): Minimum volume in ml
- `maxVolume` (number): Maximum volume in ml
- `tag` (string): Filter by tag (can repeat for multiple tags)
- `supplier` (string): Filter by supplier
- `lowStock` (boolean): Show only low stock items (< 5ml)
- `sort` (string): Sort field (`name`, `createdAt`, `volumeMl`, `position`)
- `order` (string): Sort order (`asc`, `desc`)
- `page` (number): Page number (default: 1)
- `limit` (number): Items per page (default: 50)

**Example**:
```
GET /accords?position=top&tag=fresh&tag=citrus&sort=volumeMl&order=asc
```

**Response (200)**:
```json
{
  "accords": [
    {
      "_id": "uuid",
      "userId": "user-uuid",
      "name": "Citrus Fresh Accord",
      "pyramidPosition": "top",
      "volumeMl": 25.5,
      "volumeDrops": 510,
      "supplier": "Perfumer's Apprentice",
      "purchaseDate": "2025-12-15",
      "dilutionPercentage": 10,
      "notes": "Very bright and zesty.",
      "tags": ["fresh", "citrus", "energetic"],
      "createdAt": "2025-12-15T10:00:00Z",
      "updatedAt": "2026-01-20T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 42,
    "pages": 1
  }
}
```

#### POST /accords
Create a new accord.

**Auth**: Required

**Request**:
```json
{
  "name": "Citrus Fresh Accord",
  "pyramidPosition": "top",
  "volumeMl": 25.5,
  "supplier": "Perfumer's Apprentice",
  "purchaseDate": "2025-12-15",
  "dilutionPercentage": 10,
  "notes": "Very bright and zesty.",
  "tags": ["fresh", "citrus", "energetic"]
}
```

**Validation**:
- `name`: Required, 1-255 chars
- `pyramidPosition`: Required, one of: `top`, `middle`, `base`
- `volumeMl`: Required, >= 0
- `tags`: Optional array, each tag 1-50 chars

**Response (201)**:
```json
{
  "_id": "uuid",
  "userId": "user-uuid",
  "name": "Citrus Fresh Accord",
  "pyramidPosition": "top",
  "volumeMl": 25.5,
  "volumeDrops": 510,
  "supplier": "Perfumer's Apprentice",
  "purchaseDate": "2025-12-15",
  "dilutionPercentage": 10,
  "notes": "Very bright and zesty.",
  "tags": ["fresh", "citrus", "energetic"],
  "createdAt": "2026-01-31T10:00:00Z",
  "updatedAt": "2026-01-31T10:00:00Z"
}
```

**Errors**:
- `409 Conflict`: Duplicate name+position combination

#### GET /accords/:id
Get a single accord by ID.

**Auth**: Required (must be owner)

**Response (200)**:
```json
{
  "_id": "uuid",
  "userId": "user-uuid",
  "name": "Citrus Fresh Accord",
  "pyramidPosition": "top",
  "volumeMl": 25.5,
  "volumeDrops": 510,
  "supplier": "Perfumer's Apprentice",
  "purchaseDate": "2025-12-15",
  "dilutionPercentage": 10,
  "notes": "Very bright and zesty.",
  "tags": ["fresh", "citrus", "energetic"],
  "createdAt": "2025-12-15T10:00:00Z",
  "updatedAt": "2026-01-20T10:00:00Z"
}
```

**Errors**:
- `404 Not Found`: Accord doesn't exist or not owned by user

#### PUT /accords/:id
Update an accord (partial update).

**Auth**: Required (must be owner)

**Request** (all fields optional):
```json
{
  "name": "Updated Name",
  "volumeMl": 20.0,
  "notes": "Updated notes",
  "tags": ["fresh", "new-tag"]
}
```

**Response (200)**:
```json
{
  "_id": "uuid",
  "userId": "user-uuid",
  "name": "Updated Name",
  "pyramidPosition": "top",
  "volumeMl": 20.0,
  "volumeDrops": 400,
  "tags": ["fresh", "new-tag"],
  "updatedAt": "2026-01-31T10:30:00Z",
  // ... other fields
}
```

**Notes**:
- Updating `tags` replaces all tags (not additive)
- Changing `name` or `pyramidPosition` checks uniqueness

**Errors**:
- `409 Conflict`: New name+position already exists

#### DELETE /accords/:id
Delete an accord.

**Auth**: Required (must be owner)

**Response (204)**: No content

**Errors**:
- `404 Not Found`: Accord doesn't exist or not owned by user

---

### Tags (Protected)

#### GET /tags/predefined
Get all predefined tags grouped by category.

**Auth**: Required

**Response (200)**:
```json
{
  "tags": {
    "character": ["fresh", "warm", "cool", "dry"],
    "mood": ["romantic", "energetic", "calming"],
    "season": ["spring", "summer", "autumn", "winter"],
    "scent_family": ["floral", "citrus", "woody"]
    // ... all categories
  }
}
```

#### GET /tags
Get all unique tags used by current user (custom + predefined).

**Auth**: Required

**Query Parameters**:
- `search` (string): Search tags by name

**Response (200)**:
```json
{
  "tags": [
    "fresh",
    "citrus",
    "my-custom-tag",
    "energetic"
  ]
}
```

#### POST /accords/:id/tags
Add a tag to an accord.

**Auth**: Required (must be owner)

**Request**:
```json
{
  "tag": "custom-tag"
}
```

**Response (200)**:
```json
{
  "message": "Tag added successfully"
}
```

**Errors**:
- `409 Conflict`: Tag already exists on accord

#### DELETE /accords/:id/tags/:tag
Remove a tag from an accord.

**Auth**: Required (must be owner)

**Response (200)**:
```json
{
  "message": "Tag removed successfully"
}
```

---

### Statistics (Protected)

#### GET /stats
Get accord statistics for current user.

**Auth**: Required

**Response (200)**:
```json
{
  "totalAccords": 42,
  "totalVolumeMl": 1250.5,
  "byPosition": {
    "top": {
      "count": 15,
      "volumeMl": 420.5
    },
    "middle": {
      "count": 18,
      "volumeMl": 550.0
    },
    "base": {
      "count": 9,
      "volumeMl": 280.0
    }
  },
  "mostUsedTags": [
    { "tag": "fresh", "count": 12 },
    { "tag": "warm", "count": 8 },
    { "tag": "floral", "count": 7 }
  ],
  "lowStock": [
    {
      "_id": "uuid",
      "name": "Vanilla Base",
      "pyramidPosition": "base",
      "volumeMl": 2.5
    }
  ]
}
```

---

### Export/Import (Protected)

#### GET /export
Export all user data as JSON.

**Auth**: Required

**Response (200)**:
```json
{
  "version": "2.0",
  "exportDate": "2026-01-31T10:00:00Z",
  "accords": [
    {
      "name": "Citrus Fresh Accord",
      "pyramidPosition": "top",
      "volumeMl": 25.5,
      "volumeDrops": 510,
      "supplier": "Perfumer's Apprentice",
      "purchaseDate": "2025-12-15",
      "dilutionPercentage": 10,
      "notes": "Very bright and zesty.",
      "tags": ["fresh", "citrus", "energetic"],
      "createdAt": "2025-12-15T10:00:00Z",
      "updatedAt": "2026-01-20T10:00:00Z"
    }
  ]
}
```

**Headers**:
```
Content-Type: application/json
Content-Disposition: attachment; filename="scentora-export-2026-01-31.json"
```

#### POST /export/import
Import accords from JSON.

**Auth**: Required

**Request**: Same format as export response

**Query Parameters**:
- `onDuplicate` (string): How to handle duplicates (`skip`, `overwrite`, `rename`)

**Response (200)**:
```json
{
  "imported": 15,
  "skipped": 3,
  "errors": [
    {
      "accord": "Duplicate Accord",
      "reason": "Name and position already exist"
    }
  ]
}
```

---

### Health Check (Public)

#### GET /health
Check server health.

**Response (200)**:
```json
{
  "status": "ok",
  "timestamp": "2026-01-31T10:00:00Z",
  "service": "scentora-api"
}
```

---

## Rate Limiting

### Auth Endpoints
- **Limit**: 5 requests per 15 minutes per IP
- **Applies to**: `/auth/register`, `/auth/login`, `/auth/refresh`

### General Endpoints
- **Limit**: 100 requests per minute per IP

### Response Headers
```
Rate-Limit-Total: 5
Rate-Limit-Remaining: 3
Rate-Limit-Reset: 1738332000
```

### Rate Limit Error (429)
```json
{
  "error": {
    "message": "Too many requests. Please try again later."
  }
}
```

---

## Error Responses

### Validation Error (400)
```json
{
  "error": {
    "message": "Validation failed",
    "details": [
      {
        "field": "volumeMl",
        "message": "must be greater than or equal to 0"
      }
    ]
  }
}
```

### Unauthorized (401)
```json
{
  "error": {
    "message": "Unauthorized. Please provide a valid access token."
  }
}
```

### Forbidden (403)
```json
{
  "error": {
    "message": "You don't have permission to access this resource."
  }
}
```

### Not Found (404)
```json
{
  "error": {
    "message": "Accord not found"
  }
}
```

### Conflict (409)
```json
{
  "error": {
    "message": "Accord with this name and position already exists"
  }
}
```

---

## CORS

**Allowed Origins** (configurable via environment):
- `http://localhost:5173` (frontend dev)
- `http://localhost:3000` (backend)

**Allowed Methods**:
- GET, POST, PUT, DELETE, OPTIONS

**Allowed Headers**:
- Content-Type, Authorization

---

## Notes

- All timestamps in ISO 8601 format (UTC)
- UUIDs as strings
- JSON field names in camelCase
- Pagination defaults: page=1, limit=50
- Maximum limit: 100 items per page
- Search uses case-insensitive matching
- Tags are case-sensitive
- Volume calculations: 1 ml = 20 drops (approximation)
