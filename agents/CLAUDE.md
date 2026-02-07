# Scentora - Claude AI Agent Instructions

**Last Updated**: February 7, 2026
**Project Phase**: Phase 11 Complete (Deployment) | Phase 10.9 (Views & Pages) Next
**Current State**: Backend 100% Ready (70 integration tests, 53 endpoints) | Deployed to scentora.thejoeshow.net | Frontend UI Components Complete

---

## 📋 Quick Reference

### Project Identity
- **Name**: Scentora
- **Purpose**: Perfume formulation and accord management system for DIY perfumers
- **Architecture**: Koa.js/TypeScript backend + Vue 3/TypeScript frontend + PostgreSQL
- **Repository**: github.com/xupit3r/scentora
- **Maintainer**: Joe (xupit3r)

### Current Priorities
1. **Phase 10.9**: Recipe Views & Pages (RecipesView, RecipeDetailView, routing)
2. **Testing**: Maintain 100% test passing rate (70/70 integration tests)
3. **Documentation**: Keep plan.md and phase docs updated
4. **Code Quality**: Follow established patterns, no over-engineering

---

## 🎯 Project Overview

### What is Scentora?

Scentora enables DIY perfumers to:
- **Manage accord inventory** (scent building blocks with volumes, suppliers, tags)
- **Create perfume recipes/formulas** by combining accords with version control
- **Track iterations** with immutable versions (v1, v2, v3...)
- **Organize work** with tags, collections, and rich notes
- **Export/share** recipes as JSON

### Evolution Path
1. ✅ **Phases 1-7**: Foundation (auth, basic CRUD, migration to Koa.js/TypeScript)
2. ✅ **Phase 8**: Accord system (inventory management)
3. ✅ **Phase 9**: Testing infrastructure (70 integration tests)
4. ✅ **Phase 10.1-10.5**: Recipe backend (models, repos, services, handlers, routes)
5. ✅ **Phase 10.6-10.8**: Recipe frontend types, stores, UI components
6. ✅ **Phase 11**: Deployment (Docker, CI/CD, DigitalOcean, SSL)
7. 🎯 **Phase 10.9 (NEXT)**: Recipe views & pages
8. 📋 **Phase 10.10-10.12**: Features, testing, polish

---

## 🏗️ Architecture

### Tech Stack

**Backend** (`backend/`)
- **Runtime**: Node.js 18+
- **Framework**: Koa.js (HTTP server)
- **Language**: TypeScript 5.x (strict mode)
- **Database ORM**: Prisma
- **Validation**: Zod schemas
- **Auth**: JWT (jsonwebtoken, bcryptjs)
- **Testing**: Vitest + Supertest
- **Port**: 3000

**Frontend** (`frontend/`)
- **Framework**: Vue 3.5 (Composition API with `<script setup>`)
- **Language**: TypeScript 5.x
- **State**: Pinia stores
- **UI Library**: Naive UI (minimalist components)
- **Styling**: Tailwind CSS v4
- **Build**: Vite 6.x
- **Router**: Vue Router 4.x
- **Port**: 5173

**Database**
- **Engine**: PostgreSQL 15
- **Host**: localhost:5435 (Docker)
- **Credentials**: admin/password (dev)
- **Migrations**: Prisma migrations

**Infrastructure**
- **Container (dev)**: Docker Compose (PostgreSQL on port 5435)
- **Container (prod)**: `docker-compose.prod.yml` (postgres, backend, frontend)
- **CI/CD**: GitHub Actions (`test.yml` for PRs, `deploy.yml` for main)
- **Deployment**: DigitalOcean droplet, nginx reverse proxy, Let's Encrypt SSL
- **Live URL**: https://scentora.thejoeshow.net

### Architecture Pattern

```
Routes → Services → Prisma (no repository layer)
```

**Backend Layers:**
- **Routes** (`src/routes/`): HTTP handlers, request validation, response formatting
- **Services** (`src/services/`): Business logic, authorization checks, data orchestration
- **Prisma**: Database queries via generated client (no manual repo layer)
- **Middleware** (`src/middleware/`): Auth, CORS, rate limiting, error handling
- **Schemas** (`src/schemas/`): Zod validation schemas

**Frontend Layers:**
- **Components** (`src/components/`): Reusable UI elements
- **Views** (`src/views/`): Page-level components
- **Stores** (`src/stores/`): Pinia state management
- **Services** (`src/services/`): API client (axios wrapper)
- **Router** (`src/router/`): Vue Router config
- **Types** (`src/types/`): TypeScript interfaces

---

## 📊 Current Phase Status

### ✅ Phase 11 Complete: Deployment (February 7, 2026)

**Delivered**: Full production deployment pipeline
- Docker multi-stage builds (backend: Node Alpine, frontend: nginx Alpine)
- `docker-compose.prod.yml` with health checks and dependency ordering
- GitHub Actions CI/CD: `test.yml` (PRs + reusable) → `deploy.yml` (push to main)
- Automated Prisma migrations on container startup
- DigitalOcean droplet with nginx reverse proxy + Let's Encrypt SSL
- Rate limiting disabled in test environment (`NODE_ENV=test`)
- Container health checks use `127.0.0.1` (not `localhost`) to avoid IPv6 issues

**Live**: https://scentora.thejoeshow.net

### ✅ Phase 10.8 Complete: UI Components (February 6, 2026)

**Delivered**: 9 recipe components (1,933 lines of code)
- RecipeCard.vue (display recipes in grid)
- RecipeForm.vue (create/edit recipes)
- RecipeList.vue (filtered list with empty states)
- VersionSelector.vue (dropdown to switch versions)
- VersionTimeline.vue (visual version history)
- IngredientPicker.vue (select accords)
- IngredientList.vue (ingredients with quantities)
- RecipeNoteCard.vue (display notes)
- CollectionCard.vue (display collections)

**Achievement**: Notion-inspired design, fully responsive, type-safe

### 🎯 Phase 10.9 Next: Views & Pages

**Goal**: Create page-level views and integrate routing

**Tasks**:
- [ ] RecipesView.vue - Main recipes list page
- [ ] RecipeDetailView.vue - Single recipe detail page
- [ ] RecipeCreateView.vue - Create new recipe flow
- [ ] RecipeEditView.vue - Edit existing recipe
- [ ] CollectionsView.vue - Collections list
- [ ] CollectionDetailView.vue - Collection with recipes
- [ ] Update router/index.ts with recipe routes
- [ ] Add route guards for authentication
- [ ] Implement breadcrumb navigation
- [ ] Test full navigation flow

**Estimated Duration**: 2-3 days

### Backend Status (100% Ready)

**Completion**: February 4, 2026
- ✅ 14 Prisma models (recipes, versions, ingredients, notes, tags, collections)
- ✅ Migrations applied (8 new tables)
- ✅ 6 service modules (recipe, version, ingredient, note, tag, collection)
- ✅ 53 API endpoints (fully functional)
- ✅ 70 integration tests passing (100% success rate)
- ✅ Manual testing complete (phase10-5-manual-test.sh)

**API Endpoints Available**:
- Recipes: 6 endpoints (CRUD + search)
- Versions: 5 endpoints (create, list, get, activate, duplicate)
- Ingredients: 3 endpoints (add, update, remove)
- Notes: 4 endpoints (CRUD)
- Tags: 3 endpoints (add, remove, popular)
- Collections: 7 endpoints (CRUD + membership)

---

## 🎨 Design System (Notion-Inspired)

### Design Principles

1. **Clean Typography**: System font stack, clear hierarchy, 1.5-1.7 line height
2. **Minimalist Palette**: Restrained colors, subtle accents
3. **Generous Spacing**: 8px grid system, ample whitespace
4. **Subtle Interactions**: Background changes on hover, smooth 150ms transitions
5. **Inline Actions**: Actions appear on hover only
6. **Responsive**: Mobile-first, collapsible sidebar

### Color Palette

```typescript
// Text colors
text: {
  primary: '#37352F',
  secondary: '#787774',
  tertiary: '#9B9A97',
}

// Background colors
background: {
  primary: '#FFFFFF',
  secondary: '#FAFAFA',
  tertiary: '#F7F6F3',
}

// Borders
border: {
  default: '#E9E9E7',
}

// Accents
accent: {
  primary: '#0F766E', // Teal for actions
}

// Pyramid positions (accords)
pyramid: {
  top: '#FFD93D → #FFA800',     // Yellow gradient
  middle: '#B565D8 → #8B5CF6',  // Purple gradient
  base: '#A0826D → #6B4423',    // Brown gradient
}
```

### Component Patterns

- **Cards**: Subtle shadows (0 1px 3px rgba(0,0,0,0.12)), 6-8px rounded corners
- **Hover**: Background changes, not borders. Smooth transitions (150ms)
- **Buttons**: Ghost buttons (transparent until hover)
- **Modals**: Backdrop blur, centered, max-width constraints
- **Forms**: Single-column layout, inline labels, clear focus states
- **Empty States**: Centered, icon + text + action button

---

## 📐 Database Schema

### Core Models (14 total)

**Authentication & Users**
```prisma
model User {
  id                    String   @id @default(uuid())
  email                 String   @unique
  username              String   @unique
  passwordHash          String
  validateRecipeVolumes Boolean  @default(false)
  createdAt             DateTime @default(now())
  updatedAt             DateTime @updatedAt
}

model RefreshToken {
  id           String   @id @default(uuid())
  userId       String
  tokenHash    String   @unique
  expiresAt    DateTime
  createdAt    DateTime @default(now())
  user         User     @relation(...)
}

model Invitation {
  id        String    @id @default(uuid())
  code      String    @unique
  email     String?
  createdBy String
  expiresAt DateTime
  used      Boolean   @default(false)
  usedAt    DateTime?
  usedBy    String?
  createdAt DateTime  @default(now())
  creator   User      @relation(...)
}
```

**Accord System**
```prisma
model Accord {
  id                  String   @id @default(uuid())
  userId              String
  name                String
  pyramidPosition     String   // 'top' | 'middle' | 'base'
  volumeMl            Float
  supplier            String?
  purchaseDate        DateTime?
  dilutionPercentage  Float?
  notes               String?
  createdAt           DateTime @default(now())
  updatedAt           DateTime @updatedAt

  user                User     @relation(...)
  tags                AccordTag[]
  recipeIngredients   RecipeIngredient[]

  @@unique([userId, name, pyramidPosition])
}

model AccordTag {
  id        String   @id @default(uuid())
  accordId  String
  tag       String
  createdAt DateTime @default(now())
  accord    Accord   @relation(...)

  @@unique([accordId, tag])
}

model PredefinedTag {
  id        String   @id @default(uuid())
  category  String
  tag       String   @unique
  createdAt DateTime @default(now())
}
```

**Recipe System** (Phase 10)
```prisma
model Recipe {
  id              String    @id @default(uuid())
  userId          String
  name            String
  description     String?
  targetVolumeMl  Float
  status          String    @default("draft") // draft | in_progress | tested | finalized | archived
  activeVersionId String?
  createdAt       DateTime  @default(now())
  updatedAt       DateTime  @updatedAt

  user            User      @relation(...)
  versions        RecipeVersion[]
  notes           RecipeNote[]
  tags            RecipeTag[]
  collections     RecipeCollectionMember[]

  @@unique([userId, name])
}

model RecipeVersion {
  id            String   @id @default(uuid())
  recipeId      String
  versionNumber Int
  name          String
  notes         String?
  isActive      Boolean  @default(false)
  createdAt     DateTime @default(now())

  recipe        Recipe   @relation(...)
  ingredients   RecipeIngredient[]
  versionNotes  RecipeNote[]

  @@unique([recipeId, versionNumber])
}

model RecipeIngredient {
  id             String   @id @default(uuid())
  versionId      String
  accordId       String
  quantityMl     Float
  percentage     Float?
  notes          String?
  createdAt      DateTime @default(now())

  version        RecipeVersion @relation(...)
  accord         Accord        @relation(...)

  @@unique([versionId, accordId])
}

model RecipeNote {
  id        String    @id @default(uuid())
  recipeId  String
  versionId String?
  content   String
  noteType  String    @default("general") // general | testing | observation | adjustment | reminder
  createdAt DateTime  @default(now())
  updatedAt DateTime  @updatedAt

  recipe    Recipe    @relation(...)
  version   RecipeVersion? @relation(...)
}

model RecipeTag {
  id        String   @id @default(uuid())
  recipeId  String
  tag       String
  createdAt DateTime @default(now())

  recipe    Recipe   @relation(...)

  @@unique([recipeId, tag])
}

model RecipeCollection {
  id          String   @id @default(uuid())
  userId      String
  name        String
  description String?
  createdAt   DateTime @default(now())
  updatedAt   DateTime @updatedAt

  user        User     @relation(...)
  members     RecipeCollectionMember[]

  @@unique([userId, name])
}

model RecipeCollectionMember {
  id           String   @id @default(uuid())
  collectionId String
  recipeId     String
  addedAt      DateTime @default(now())

  collection   RecipeCollection @relation(...)
  recipe       Recipe           @relation(...)

  @@unique([collectionId, recipeId])
}
```

### Key Constraints

**Foreign Key Behavior:**
- `recipe_ingredients.accord_id` → **ON DELETE RESTRICT** (prevent accord deletion if used)
- All other recipe FKs → **ON DELETE CASCADE** (clean removal)
- `users` → **ON DELETE CASCADE** (remove all user data)

**Unique Constraints:**
- `(userId, name)` on recipes (unique names per user)
- `(userId, name, pyramidPosition)` on accords
- `(recipeId, versionNumber)` on versions
- `(versionId, accordId)` on ingredients (no duplicate ingredients per version)

---

## 🔒 Security & Authentication

### JWT Token System

**Access Tokens:**
- Lifetime: 15 minutes
- Usage: API authentication via `Authorization: Bearer <token>`
- Storage: Memory only (not localStorage)

**Refresh Tokens:**
- Lifetime: 7 days
- Usage: Obtain new access tokens via `/api/auth/refresh`
- Storage: httpOnly cookie (secure in production)
- Rotation: New refresh token issued on each refresh

**Password Security:**
- Hashing: bcryptjs (cost factor 10)
- Validation: Zod schemas enforce minimum requirements
- Never log or expose passwords

### Rate Limiting

**Auth Endpoints** (`/api/auth/*`):
- Limit: 5 requests per 15 minutes per IP
- Purpose: Brute force protection

**General Endpoints**:
- Limit: 100 requests per minute per IP
- Implementation: In-memory (consider Redis for production)

### User Data Isolation

**Enforcement Points:**
1. **Service Layer**: All queries filtered by `userId` from JWT
2. **Route Handlers**: Extract `userId` from `ctx.state.user`
3. **Database**: User-specific queries via Prisma `where: { userId }`
4. **Authorization**: Verify ownership before mutations

**Never:**
- Expose other users' data
- Allow cross-user operations
- Trust client-provided user IDs

---

## ⚙️ Development Standards

### Code Style

**TypeScript:**
- Strict mode enabled
- Explicit types (avoid `any`)
- Interfaces for object shapes, types for unions
- Descriptive names (no abbreviations unless standard)

**Vue Components:**
- Use `<script setup>` syntax
- Composition API (not Options API)
- Props with TypeScript interfaces
- Emit events with types

**Naming Conventions:**
- Files: PascalCase for components (`RecipeCard.vue`), kebab-case for utilities
- Functions: camelCase (`createRecipe`)
- Classes/Interfaces: PascalCase (`Recipe`, `RecipeRequest`)
- Constants: UPPER_SNAKE_CASE (`MAX_RECIPES`)

### Testing Philosophy (CRITICAL)

**Testing is mandatory and non-negotiable.**

**When to Test:**
- ✅ Before committing code
- ✅ Before completing a phase
- ✅ Before creating pull requests
- ✅ After fixing bugs
- ✅ Before making breaking changes

**Coverage Targets:**
- Overall: 80%+
- Services: 85%+
- Routes: 80%+
- Critical paths: 100%

**Test Types:**
1. **Unit Tests**: Services, utilities (mock dependencies)
2. **Integration Tests**: API endpoints (use test database)
3. **Component Tests**: Vue components (deferred to future)

**Running Tests:**
```bash
cd backend
npm test                 # Run all tests
npm run test:watch       # Watch mode
```

**Current Status**: 70/70 integration tests passing (100%)

### Error Handling

**Backend Errors:**
```typescript
// Always use structured error responses
return ctx.body = {
  error: {
    message: 'User-friendly message',
    details: 'Technical details (optional)'
  }
};

// Status codes
200: Success
201: Created
400: Validation error
401: Unauthorized (missing/invalid token)
403: Forbidden (insufficient permissions)
404: Not found
409: Conflict (duplicate resource)
429: Rate limit exceeded
500: Internal server error
```

**Frontend Errors:**
```typescript
// Use try-catch with toast notifications
try {
  await recipeStore.createRecipe(data);
  message.success('Recipe created successfully');
  router.push(`/recipes/${newRecipe._id}`);
} catch (error) {
  message.error(error.message || 'Failed to create recipe');
}
```

### Documentation Standards

**File Headers:**
- Add brief comment at top of new files explaining purpose
- No need for authorship or dates (git tracks that)

**Functions:**
- JSDoc for public APIs
- Inline comments for complex logic (explain WHY, not WHAT)
- Self-documenting code preferred

**Phase Completion:**
- Document in `docs/phases/PHASE{X}_{Y}_COMPLETE.md`
- Include: Summary, changes, test results, next steps
- Update `plan.md` with completion date

---

## 🚀 Common Workflows

### Starting Development

```bash
# 1. Start PostgreSQL
docker compose up -d

# 2. Start backend (terminal 1)
cd backend
npm install          # First time only
npx prisma generate  # First time only
npx prisma db push   # First time only (applies schema)
npm run dev

# 3. Start frontend (terminal 2)
cd frontend
npm install          # First time only
npm run dev

# 4. Access app
# Frontend: http://localhost:5173
# Backend: http://localhost:3000
# Health: http://localhost:3000/api/health
```

### Creating API Endpoints

**Backend Flow:**
1. Define Zod schema in `src/schemas/`
2. Add service method in `src/services/`
3. Add route handler in `src/routes/`
4. Register route in `src/app.ts`
5. Write tests in `tests/`
6. Update `specs/recipe-api.md`

**Example:**
```typescript
// 1. Schema (src/schemas/recipe.ts)
export const createRecipeSchema = z.object({
  name: z.string().min(1).max(255),
  description: z.string().optional(),
  targetVolumeMl: z.number().positive(),
});

// 2. Service (src/services/recipe.ts)
export async function createRecipe(userId: string, data: CreateRecipeData) {
  // Validate unique name
  const exists = await prisma.recipe.findFirst({
    where: { userId, name: data.name }
  });
  if (exists) throw new Error('Recipe name already exists');

  // Create recipe
  return await prisma.recipe.create({
    data: { userId, ...data }
  });
}

// 3. Route (src/routes/recipes.ts)
router.post('/api/recipes', authMiddleware, async (ctx) => {
  const userId = ctx.state.user.userId;
  const body = createRecipeSchema.parse(ctx.request.body);

  const recipe = await createRecipe(userId, body);

  ctx.status = 201;
  ctx.body = recipe;
});

// 4. Test (tests/recipes.test.ts)
describe('POST /api/recipes', () => {
  it('creates a recipe', async () => {
    const res = await request(app)
      .post('/api/recipes')
      .set('Authorization', `Bearer ${accessToken}`)
      .send({ name: 'Test Recipe', targetVolumeMl: 50 })
      .expect(201);

    expect(res.body.name).toBe('Test Recipe');
  });
});
```

### Creating Vue Components

**Component Structure:**
```vue
<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Recipe } from '@/types/recipe';

// Props
interface Props {
  recipe: Recipe;
  compact?: boolean;
}
const props = withDefaults(defineProps<Props>(), {
  compact: false
});

// Emits
interface Emits {
  (e: 'click', recipe: Recipe): void;
  (e: 'delete', id: string): void;
}
const emit = defineEmits<Emits>();

// State
const isHovered = ref(false);

// Computed
const displayName = computed(() => {
  return props.compact ? props.recipe.name : `${props.recipe.name} v${props.recipe.activeVersionNumber}`;
});

// Methods
function handleClick() {
  emit('click', props.recipe);
}
</script>

<template>
  <div
    class="recipe-card"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false"
    @click="handleClick"
  >
    <h3>{{ displayName }}</h3>
    <p v-if="recipe.description">{{ recipe.description }}</p>

    <div v-if="isHovered" class="actions">
      <button @click.stop="$emit('delete', recipe._id)">Delete</button>
    </div>
  </div>
</template>

<style scoped>
.recipe-card {
  padding: 1rem;
  border-radius: 8px;
  background: var(--bg-secondary);
  cursor: pointer;
  transition: background 150ms ease;
}

.recipe-card:hover {
  background: var(--bg-tertiary);
}

.actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.75rem;
}
</style>
```

### Adding Database Models

**Using Prisma:**
1. Edit `backend/prisma/schema.prisma`
2. Run `npx prisma db push` (applies changes)
3. Run `npx prisma generate` (updates Prisma client)
4. Update TypeScript types if needed
5. Create service methods
6. Add API endpoints
7. Write tests

**Example Migration (manual if needed):**
```bash
# Generate migration
npx prisma migrate dev --name add_recipe_collections

# Apply migration
npx prisma migrate deploy
```

---

## 🎯 Important Constraints & Rules

### Do's ✅

1. **Testing**: Always write tests before committing
2. **Type Safety**: Use TypeScript strictly, no `any`
3. **User Isolation**: Filter all queries by `userId`
4. **Error Messages**: Provide clear, actionable error messages
5. **Documentation**: Update plan.md and phase docs
6. **Security**: Validate all inputs, never expose secrets
7. **Patterns**: Follow established code patterns
8. **Simplicity**: Keep solutions simple, avoid over-engineering

### Don'ts ❌

1. **No Over-Engineering**: Don't add features not requested
2. **No Unnecessary Abstractions**: Three similar lines > premature abstraction
3. **No Skipping Tests**: Never commit untested code
4. **No Commented Code**: Remove code, don't comment it out (git tracks history)
5. **No Breaking Changes**: Without user approval first
6. **No Hardcoded Secrets**: Use environment variables
7. **No Production Data**: Use test database for development
8. **No Force Push**: To main/master branch

### Recipe System Specific Rules

1. **Volume Validation**: Disabled by default (user opt-in via preferences)
2. **Version Auto-Activation**: New versions automatically become active
3. **Accord Protection**: Prevent deletion if used in recipes (return 409 with recipe list)
4. **Immutable Versions**: Once created, versions cannot be edited (duplicate instead)
5. **Recipe Status**: Five states (draft, in_progress, tested, finalized, archived)

---

## 📝 API Quick Reference

### Authentication
```
POST   /api/auth/register       Register (requires invitation)
POST   /api/auth/login          Login
POST   /api/auth/refresh        Refresh tokens
POST   /api/auth/logout         Logout (revoke token)
POST   /api/auth/logout-all     Logout all devices
GET    /api/auth/me             Get current user
```

### Accords (Protected)
```
GET    /api/accords             List accords (with filters)
POST   /api/accords             Create accord
GET    /api/accords/:id         Get accord
PUT    /api/accords/:id         Update accord
DELETE /api/accords/:id         Delete accord
POST   /api/accords/:id/tags    Add tag
DELETE /api/accords/:id/tags/:tag  Remove tag
```

### Recipes (Protected)
```
GET    /api/recipes             List recipes
POST   /api/recipes             Create recipe
GET    /api/recipes/:id         Get recipe
PUT    /api/recipes/:id         Update recipe
DELETE /api/recipes/:id         Delete recipe
GET    /api/recipes/search?q=   Search recipes
```

### Recipe Versions (Protected)
```
GET    /api/recipes/:id/versions                      List versions
POST   /api/recipes/:id/versions                      Create version
GET    /api/recipes/:id/versions/:versionId           Get version
POST   /api/recipes/:id/versions/:versionId/activate  Set active version
POST   /api/recipes/:id/versions/:versionId/duplicate Duplicate version
```

### Recipe Ingredients (Protected)
```
POST   /api/recipes/:id/versions/:versionId/ingredients                Add ingredient
PUT    /api/recipes/:id/versions/:versionId/ingredients/:ingredientId  Update ingredient
DELETE /api/recipes/:id/versions/:versionId/ingredients/:ingredientId  Remove ingredient
```

### Recipe Notes (Protected)
```
GET    /api/recipes/:id/notes         List notes
POST   /api/recipes/:id/notes         Create note
PUT    /api/recipes/:id/notes/:noteId Update note
DELETE /api/recipes/:id/notes/:noteId Delete note
```

### Recipe Tags (Protected)
```
POST   /api/recipes/:id/tags      Add tag
DELETE /api/recipes/:id/tags/:tag Remove tag
GET    /api/recipes/tags/popular  Popular tags
```

### Collections (Protected)
```
GET    /api/collections                                List collections
POST   /api/collections                                Create collection
GET    /api/collections/:id                            Get collection
PUT    /api/collections/:id                            Update collection
DELETE /api/collections/:id                            Delete collection
POST   /api/collections/:id/recipes                    Add recipe
DELETE /api/collections/:id/recipes/:recipeId          Remove recipe
```

### Tags (Public)
```
GET    /api/tags                  All predefined tags (grouped)
GET    /api/tags/search?q=        Search tags
GET    /api/tags/categories       List categories
```

### Stats & Export (Protected)
```
GET    /api/stats                 Accord statistics
GET    /api/export                Export accords JSON
POST   /api/export/import         Import accords JSON
```

### Health (Public)
```
GET    /api/health                Health check
```

**Total Endpoints**: 53

---

## 🔍 Debugging Tips

### Common Issues

**"Token expired" errors:**
- Access tokens expire in 15 minutes
- Implement automatic refresh in frontend auth interceptor
- Check `ctx.state.user` exists in protected routes

**"Accord is used in recipes" error:**
- Expected behavior (ON DELETE RESTRICT)
- User must remove accord from recipes first
- Error includes list of recipes using the accord

**Test failures:**
- Run tests sequentially: `npm test`
- Check test database is clean (tests should clean up)
- Verify Prisma schema is up to date

**CORS errors:**
- Check `CORS_ALLOWED_ORIGINS` in `.env`
- Ensure frontend origin (http://localhost:5173) is allowed

**Database connection errors:**
- Verify PostgreSQL container is running: `docker compose ps`
- Check `DATABASE_URL` in `.env`
- Run `npx prisma db push` to sync schema

### Logging

**Backend:**
```typescript
// Use console.log for development debugging
console.log('Creating recipe:', data);

// For production, consider structured logging
// (e.g., winston, pino)
```

**Frontend:**
```typescript
// Use console.debug for development
console.debug('Recipe store state:', recipeStore.recipes);

// Use Vue DevTools for component inspection
// Install: https://devtools.vuejs.org/
```

---

## 📚 Key Documentation Files

**Planning & Tracking:**
- `plan.md` - Master development plan (41KB, comprehensive)
- `docs/phases/` - Phase completion documents
- `specs/recipe-system.md` - Phase 10 specification

**Technical Specs:**
- `specs/api-spec.md` - Accord API documentation
- `specs/recipe-api.md` - Recipe API documentation
- `specs/data-models.md` - Database schema
- `specs/tag-system.md` - Tag categories (57+ tags)

**Guides:**
- `README.md` - Project overview, quick start
- `QUICKSTART.md` - Quick start guide
- `docs/AUTH_IMPLEMENTATION.md` - Authentication details
- `docs/TESTING_IMPLEMENTATION.md` - Testing framework

**Frontend:**
- `frontend/src/design-system/README.md` - Design system tokens

---

## 🎓 Learning Resources

**Project-Specific:**
- Study `plan.md` for project history and roadmap
- Read `docs/CODEBASE_STUDY_SUMMARY.md` for architecture overview
- Review test files for usage examples

**Technology:**
- Koa.js: https://koajs.com/
- Prisma: https://www.prisma.io/docs
- Vue 3: https://vuejs.org/guide/
- Pinia: https://pinia.vuejs.org/
- Naive UI: https://www.naiveui.com/
- Zod: https://zod.dev/

---

## 🚨 Emergency Contacts

**Repository Issues**: https://github.com/xupit3r/scentora/issues
**Maintainer**: Joe (xupit3r)
**License**: MIT

---

## 🎯 Current Priorities (February 7, 2026)

1. **Phase 10.9**: Recipe Views & Pages
   - RecipesView.vue (main list page)
   - RecipeDetailView.vue (single recipe)
   - Router integration
   - Navigation testing

2. **Maintain Quality**:
   - Keep 100% test pass rate (70/70)
   - Update documentation as you go
   - Follow established patterns

3. **No Over-Engineering**:
   - Implement only what's specified
   - Keep solutions simple
   - Build on existing components

---

**Last Updated**: February 7, 2026
**Status**: Phase 11 (Deployment) Complete | Ready for Phase 10.9 - Views & Pages
**Test Status**: 70/70 passing (100%)
**Deployment**: Live at https://scentora.thejoeshow.net
**Backend**: 100% complete and functional
**Frontend**: UI components complete, views next
