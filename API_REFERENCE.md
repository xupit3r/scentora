# Scentora API Reference

## Base URL
```
http://localhost:3000/api
```

## Endpoints

### Health Check
```http
GET /health
```
Returns server status.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2026-01-17T20:00:00.000Z",
  "service": "scentora-api"
}
```

---

### Perfumes

#### List All Perfumes
```http
GET /perfumes
```
Returns all perfumes sorted by creation date (newest first).

**Response:**
```json
[
  {
    "_id": "abc123",
    "_rev": "1-xyz",
    "type": "perfume",
    "name": "Bleu de Chanel",
    "designer": "Chanel",
    "year": 2010,
    "concentration": "EDP",
    "pyramid": {
      "top": ["Grapefruit", "Lemon", "Mint"],
      "middle": ["Ginger", "Nutmeg", "Jasmine"],
      "base": ["Incense", "Vetiver", "Cedar"]
    },
    "description": "A woody aromatic fragrance",
    "imageUrl": "https://...",
    "createdAt": "2026-01-17T20:00:00.000Z",
    "updatedAt": "2026-01-17T20:00:00.000Z"
  }
]
```

#### Get Perfume by ID
```http
GET /perfumes/:id
```

**Response:** Single perfume object (same structure as above)

#### Create Perfume
```http
POST /perfumes
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Bleu de Chanel",
  "designer": "Chanel",
  "year": 2010,
  "concentration": "EDP",
  "pyramid": {
    "top": ["Grapefruit", "Lemon", "Mint"],
    "middle": ["Ginger", "Nutmeg", "Jasmine"],
    "base": ["Incense", "Vetiver", "Cedar"]
  },
  "description": "A woody aromatic fragrance",
  "imageUrl": "https://..."
}
```

**Required Fields:** `name`, `designer`, `pyramid` (with top/middle/base arrays)

**Optional Fields:** `year`, `concentration`, `description`, `imageUrl`

**Response:** Created perfume object with `_id`, `_rev`, `createdAt`, `updatedAt`

#### Update Perfume
```http
PUT /perfumes/:id
Content-Type: application/json
```

**Request Body:** Same as create, but all fields optional (partial update)

**Response:** Updated perfume object

#### Delete Perfume
```http
DELETE /perfumes/:id
```

**Response:** 204 No Content

---

### Journal Entries

#### Get Journal Entries for Perfume
```http
GET /perfumes/:perfumeId/journal
```
Returns all journal entries for a perfume, sorted by date (newest first).

**Response:**
```json
[
  {
    "_id": "entry123",
    "_rev": "1-xyz",
    "type": "journal",
    "perfumeId": "abc123",
    "date": "2026-01-17",
    "content": "Wore this to work today. Got many compliments!",
    "rating": 9,
    "occasion": "Work",
    "weather": "Sunny",
    "createdAt": "2026-01-17T20:00:00.000Z",
    "updatedAt": "2026-01-17T20:00:00.000Z"
  }
]
```

#### Create Journal Entry
```http
POST /perfumes/:perfumeId/journal
Content-Type: application/json
```

**Request Body:**
```json
{
  "perfumeId": "abc123",
  "date": "2026-01-17",
  "content": "Wore this to work today. Got many compliments!",
  "rating": 9,
  "occasion": "Work",
  "weather": "Sunny"
}
```

**Required Fields:** `perfumeId`, `date`, `content`

**Optional Fields:** `rating` (1-10), `occasion`, `weather`

**Response:** Created journal entry object

#### Update Journal Entry
```http
PUT /journal/:id
Content-Type: application/json
```

**Request Body:** Same as create (excluding `perfumeId`), all fields optional

**Response:** Updated journal entry object

#### Delete Journal Entry
```http
DELETE /journal/:id
```

**Response:** 204 No Content

---

## Error Responses

### Validation Error (400)
```json
{
  "error": {
    "message": "Validation failed",
    "details": [
      {
        "code": "too_small",
        "minimum": 1,
        "type": "string",
        "path": ["name"],
        "message": "Name is required"
      }
    ]
  }
}
```

### Not Found (404)
```json
{
  "error": {
    "message": "Perfume not found"
  }
}
```

### Server Error (500)
```json
{
  "error": {
    "message": "Internal server error"
  }
}
```

---

## Data Models

### Perfume Pyramid
```typescript
{
  top: string[];      // Top notes (volatile, first impression)
  middle: string[];   // Middle/heart notes (core character)
  base: string[];     // Base notes (lasting, foundation)
}
```

### Concentrations
Common values: `Parfum`, `EDP`, `EDT`, `EDC`

### Rating Scale
1-10, where:
- 1-3: Poor
- 4-6: Average
- 7-8: Good
- 9-10: Excellent

---

## cURL Examples

### Create a Perfume
```bash
curl -X POST http://localhost:3000/api/perfumes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bleu de Chanel",
    "designer": "Chanel",
    "year": 2010,
    "concentration": "EDP",
    "pyramid": {
      "top": ["Grapefruit", "Lemon", "Mint"],
      "middle": ["Ginger", "Nutmeg", "Jasmine"],
      "base": ["Incense", "Vetiver", "Cedar"]
    },
    "description": "A woody aromatic fragrance"
  }'
```

### Get All Perfumes
```bash
curl http://localhost:3000/api/perfumes
```

### Add Journal Entry
```bash
curl -X POST http://localhost:3000/api/perfumes/abc123/journal \
  -H "Content-Type: application/json" \
  -d '{
    "perfumeId": "abc123",
    "date": "2026-01-17",
    "content": "Great for spring weather!",
    "rating": 8,
    "occasion": "Casual",
    "weather": "Sunny"
  }'
```
