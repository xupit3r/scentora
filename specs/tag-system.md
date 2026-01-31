# Tag System Specification

**Last Updated**: January 31, 2026  
**Version**: 2.0

---

## Overview

The tag system allows users to categorize and describe their accords using predefined descriptive tags or custom tags. Tags enable flexible organization, powerful filtering, and help users understand the characteristics of their accords.

---

## Tag Categories

### 1. Character (10 tags)
Describes the overall feel or personality of the accord.

- `fresh` - Clean, invigorating, bright
- `warm` - Cozy, comforting, enveloping
- `cool` - Refreshing, crisp, airy
- `dry` - Not sweet, austere, mineral
- `powdery` - Soft, talc-like, intimate
- `creamy` - Rich, smooth, indulgent
- `sharp` - Piercing, distinct, attention-grabbing
- `soft` - Gentle, subtle, delicate
- `rich` - Opulent, full-bodied, luxurious
- `light` - Airy, weightless, translucent

### 2. Mood (8 tags)
The emotional impact or feeling evoked by the accord.

- `romantic` - Intimate, seductive, loving
- `sensual` - Provocative, alluring, sexy
- `energetic` - Lively, invigorating, dynamic
- `calming` - Soothing, relaxing, peaceful
- `mysterious` - Enigmatic, intriguing, complex
- `playful` - Fun, whimsical, lighthearted
- `sophisticated` - Elegant, refined, mature
- `innocent` - Pure, youthful, naive

### 3. Season (4 tags)
Best suited season for the accord.

- `spring` - Fresh, blooming, renewal
- `summer` - Bright, light, sunny
- `autumn` - Warm, cozy, transitional
- `winter` - Deep, rich, comforting

### 4. Time (4 tags)
Best time of day for the accord.

- `morning` - Fresh start, energizing
- `afternoon` - Versatile, balanced
- `evening` - Sophisticated, transitional
- `night` - Deep, sensual, bold

### 5. Intensity (5 tags)
The strength and projection of the accord.

- `subtle` - Barely there, skin scent
- `moderate` - Noticeable but not overpowering
- `strong` - Prominent, clear presence
- `intense` - Powerful, commanding
- `bold` - Dramatic, unforgettable

### 6. Quality (7 tags)
The character and nature of the ingredients.

- `clean` - Fresh, scrubbed, soapy
- `dirty` - Earthy, raw, unrefined
- `animalic` - Musky, leathery, primal
- `synthetic` - Modern, chemical, sharp
- `natural` - Organic, raw, authentic
- `modern` - Contemporary, innovative
- `vintage` - Classic, nostalgic, traditional

### 7. Scent Family (8 tags)
Traditional perfume classification families.

- `floral` - Flowers, petals, blooms
- `fruity` - Fruits, berries, juicy
- `woody` - Woods, bark, resinous
- `oriental` - Spices, amber, incense
- `fresh` - Citrus, aquatic, green
- `aromatic` - Herbs, lavender, sage
- `spicy` - Pepper, cinnamon, clove
- `gourmand` - Edible, sweet, dessert-like

### 8. Texture (6 tags)
The tactile or visceral quality of the scent.

- `smooth` - Even, seamless, polished
- `rough` - Textured, uneven, rustic
- `silky` - Luxurious, flowing, elegant
- `velvety` - Plush, soft, enveloping
- `airy` - Light, floating, ethereal
- `dense` - Heavy, thick, substantial

### 9. Style (5 tags)
The overall aesthetic or vibe.

- `casual` - Everyday, approachable, relaxed
- `formal` - Professional, polished, refined
- `sporty` - Athletic, dynamic, active
- `elegant` - Graceful, sophisticated, timeless
- `edgy` - Unconventional, daring, bold

---

## Total Predefined Tags: 57

---

## Database Seeding

### SQL Insert Script

```sql
-- Character tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('character', 'fresh'),
    ('character', 'warm'),
    ('character', 'cool'),
    ('character', 'dry'),
    ('character', 'powdery'),
    ('character', 'creamy'),
    ('character', 'sharp'),
    ('character', 'soft'),
    ('character', 'rich'),
    ('character', 'light');

-- Mood tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('mood', 'romantic'),
    ('mood', 'sensual'),
    ('mood', 'energetic'),
    ('mood', 'calming'),
    ('mood', 'mysterious'),
    ('mood', 'playful'),
    ('mood', 'sophisticated'),
    ('mood', 'innocent');

-- Season tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('season', 'spring'),
    ('season', 'summer'),
    ('season', 'autumn'),
    ('season', 'winter');

-- Time tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('time', 'morning'),
    ('time', 'afternoon'),
    ('time', 'evening'),
    ('time', 'night');

-- Intensity tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('intensity', 'subtle'),
    ('intensity', 'moderate'),
    ('intensity', 'strong'),
    ('intensity', 'intense'),
    ('intensity', 'bold');

-- Quality tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('quality', 'clean'),
    ('quality', 'dirty'),
    ('quality', 'animalic'),
    ('quality', 'synthetic'),
    ('quality', 'natural'),
    ('quality', 'modern'),
    ('quality', 'vintage');

-- Scent family tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('scent_family', 'floral'),
    ('scent_family', 'fruity'),
    ('scent_family', 'woody'),
    ('scent_family', 'oriental'),
    ('scent_family', 'fresh'),
    ('scent_family', 'aromatic'),
    ('scent_family', 'spicy'),
    ('scent_family', 'gourmand');

-- Texture tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('texture', 'smooth'),
    ('texture', 'rough'),
    ('texture', 'silky'),
    ('texture', 'velvety'),
    ('texture', 'airy'),
    ('texture', 'dense');

-- Style tags
INSERT INTO predefined_tags (category, tag) VALUES
    ('style', 'casual'),
    ('style', 'formal'),
    ('style', 'sporty'),
    ('style', 'elegant'),
    ('style', 'edgy');
```

---

## Custom Tags

Users can create custom tags beyond the predefined set.

**Rules**:
- 1-50 characters
- Case-sensitive
- No duplicates per accord
- No category required
- Stored in `accord_tags` table
- Not stored in `predefined_tags` (user-specific)

**Examples of Custom Tags**:
- `my-signature-blend`
- `needs-testing`
- `favorite`
- `discontinued`
- `expensive`
- `gift-idea`

---

## Tag Management

### Adding Tags to Accord

**During Creation**:
```json
POST /api/accords
{
  "name": "Citrus Fresh",
  "pyramidPosition": "top",
  "volumeMl": 25.0,
  "tags": ["fresh", "citrus", "energetic", "my-favorite"]
}
```

**After Creation**:
```json
POST /api/accords/:id/tags
{
  "tag": "summer"
}
```

### Removing Tags

```json
DELETE /api/accords/:id/tags/summer
```

### Updating All Tags at Once

```json
PUT /api/accords/:id
{
  "tags": ["fresh", "warm", "new-tag"]
}
```

*Note: This replaces all tags, not additive*

---

## Tag Display

### UI Presentation

**Predefined Tags**:
- Displayed with category grouping in autocomplete
- Color-coded by category (optional)
- Alphabetically sorted within category

**Custom Tags**:
- Displayed separately as "Your Tags"
- Alphabetically sorted
- Option to mark as custom (e.g., with icon)

**Selected Tags**:
- Shown as colored chips/badges
- Removable (X button)
- Maximum 50 tags per accord (practical limit)

---

## Tag Filtering

### Single Tag Filter
```
GET /api/accords?tag=fresh
```

### Multiple Tag Filter (AND logic)
```
GET /api/accords?tag=fresh&tag=citrus&tag=summer
```

Returns accords that have ALL specified tags.

### Tag Search (Autocomplete)
```
GET /api/tags?search=fre
```

Returns: `["fresh", "freeform-custom-tag"]`

---

## Tag Statistics

### Most Used Tags
```json
GET /api/stats

{
  "mostUsedTags": [
    { "tag": "fresh", "count": 12 },
    { "tag": "warm", "count": 8 },
    { "tag": "floral", "count": 7 }
  ]
}
```

### User's Unique Tags
```json
GET /api/tags

{
  "tags": [
    "fresh",
    "warm",
    "my-custom-tag",
    "another-custom"
  ]
}
```

---

## Tag Color Coding (Frontend)

### Category Colors (Optional)

- **Character**: Blue (#3B82F6)
- **Mood**: Pink (#EC4899)
- **Season**: Green (#10B981)
- **Time**: Purple (#8B5CF6)
- **Intensity**: Orange (#F59E0B)
- **Quality**: Gray (#6B7280)
- **Scent Family**: Teal (#14B8A6)
- **Texture**: Indigo (#6366F1)
- **Style**: Amber (#F59E0B)
- **Custom**: Dark Gray (#374151)

---

## Tag Validation

### Backend Validation

```go
type TagValidator struct {
    Tag string `validate:"required,min=1,max=50"`
}
```

### Rules
- Required: Must not be empty
- Length: 1-50 characters
- Format: No restrictions (alphanumeric, hyphens, spaces allowed)
- Case: Preserved as entered
- Uniqueness: Per accord (cannot add same tag twice)

### Error Messages

- Empty tag: `"Tag cannot be empty"`
- Too long: `"Tag must be 50 characters or less"`
- Duplicate: `"This tag already exists on this accord"`

---

## Tag Search Algorithm

### Fuzzy Matching
- Case-insensitive search
- Partial match from start: `"fre"` matches `"fresh"`
- Not substring: `"res"` does NOT match `"fresh"`

### Priority
1. Predefined tags (exact match)
2. Predefined tags (starts with)
3. User's custom tags (exact match)
4. User's custom tags (starts with)

---

## Best Practices

### For Users

**DO**:
- Use multiple tags for better organization
- Combine predefined and custom tags
- Be consistent with custom tag naming
- Use lowercase for custom tags (easier to search)
- Tag immediately after creating accord

**DON'T**:
- Create near-duplicate tags (`fresh` vs `freshness`)
- Use overly long tag names
- Use tags as notes (use notes field instead)
- Create too many categories of custom tags

### For Development

**DO**:
- Validate tag length on frontend and backend
- Trim whitespace from tags
- Index `tag` field for fast filtering
- Cache predefined tags in frontend
- Show tag suggestions as user types

**DON'T**:
- Allow unlimited tags per accord
- Allow empty strings as tags
- Case-insensitive duplicate checking
- Store custom tags in predefined_tags table

---

## Future Enhancements

### Potential Features
- Tag synonyms (map "summery" → "summer")
- Tag popularity ranking
- Suggested tags based on accord name
- Tag combinations ("summer" + "morning" = "beach")
- Tag clouds visualization
- Collaborative tagging (community suggestions)
- Tag hierarchies (parent-child relationships)
- Bulk tag operations (add tag to multiple accords)

---

## Notes

- Predefined tags are shared across all users (read-only)
- Custom tags are user-specific (not shared)
- Tags are not normalized (case-preserved)
- No limit on number of predefined tags
- Practical limit of 50 tags per accord (UI constraint)
- Tags can be filtered in combination (AND logic)
- Tag deletion from accord doesn't delete custom tag from system
- Unused custom tags remain in database (no auto-cleanup)

---

## References

- Fragrance Wheel (Michael Edwards)
- Perfume families and classifications
- Common perfumery terminology
- User research on fragrance descriptions
