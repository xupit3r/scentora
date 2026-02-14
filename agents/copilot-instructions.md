# Scentora - GitHub Copilot Instructions

**Project**: Perfume formulation and accord management system
**Stack**: Koa.js/TypeScript backend + Vue 3/TypeScript frontend + PostgreSQL
**Phase**: Phase 13.1 Complete (CLI Console) - February 10, 2026

---

## Quick Context

### What is Scentora?
A web app for DIY perfumers to:
- Manage accord inventory (scent building blocks)
- Create recipes by combining accords
- Track recipe versions (v1, v2, v3...)
- Organize with tags and collections

### Current Status
- ✅ Backend 100% ready (Koa.js 3.1.1, 70 integration tests passing)
- ✅ Frontend complete (Vue 3, Naive UI, Notion-inspired design)
- ✅ Recipe system complete (versions, ingredients, notes, collections)
- ✅ Deployed to https://scentora.thejoeshow.net (Docker + CI/CD + SSL)
- ✅ CLI Console (interactive REPL for user management)
- 🎯 Next: CLI tests and future phases (accord/recipe management, DB operations)

---

## Tech Stack

**Backend** (Port 3000)
- Node.js 20+ with Koa.js framework
- TypeScript 5.x (strict mode enabled)
- Prisma ORM (PostgreSQL 15)
- Zod for validation
- JWT auth (jsonwebtoken, bcryptjs)
- Vitest + Supertest for testing

**Frontend** (Port 5173)
- Vue 3.5 (Composition API with `<script setup>`)
- TypeScript 5.x
- Pinia for state management
- Naive UI components (minimalist design)
- Tailwind CSS v4
- Vite 6.x build tool
- Vue Router 4.x

**Database**
- PostgreSQL 15 (Docker on port 5435)
- Prisma migrations
- 14 models (users, accords, recipes, versions, ingredients, notes, tags, collections)

---

## Project Structure

```
scentora/
├── backend/                    # Koa.js API
│   ├── prisma/
│   │   └── schema.prisma      # Database schema (14 models)
│   ├── cli/                   # Interactive CLI console (Phase 13)
│   │   ├── commands/          # User, accord, recipe commands
│   │   ├── utils/             # Output, validation utilities
│   │   ├── types/             # Context types
│   │   ├── context.ts         # App context loader
│   │   └── index.ts           # CLI entry point
│   ├── src/
│   │   ├── middleware/        # Auth, CORS, rate limiting, error handling
│   │   ├── routes/            # HTTP route handlers
│   │   ├── services/          # Business logic
│   │   ├── schemas/           # Zod validation schemas
│   │   ├── utils/             # Errors, response helpers
│   │   ├── types/             # TypeScript types
│   │   ├── app.ts             # Koa app setup
│   │   └── index.ts           # Entry point
│   └── tests/                 # Vitest integration tests
├── frontend/                   # Vue.js SPA
│   ├── src/
│   │   ├── components/        # Vue components
│   │   │   ├── layout/       # AppSidebar, AppBreadcrumbs
│   │   │   ├── recipe/       # Recipe, version, ingredient, note, collection components
│   │   │   ├── ui/           # EmptyState, SkeletonCard
│   │   │   └── *.vue         # Accord components
│   │   ├── views/             # Page components
│   │   ├── stores/            # Pinia stores (auth, recipe, recipeCollection)
│   │   ├── services/          # API client (axios)
│   │   ├── composables/       # useKeyboard, useConfirmDialog
│   │   ├── utils/             # volume.ts (conversion utilities)
│   │   ├── types/             # TypeScript interfaces
│   │   ├── router/            # Vue Router config
│   │   └── design-system/     # Design tokens
│   └── ...
├── specs/                      # API and design specifications
│   ├── cli-console.md         # CLI console specification
│   └── ...
├── docs/                       # Documentation
│   ├── CLI_CONSOLE.md         # CLI user guide
│   └── phases/                # Phase completion documents
├── docker-compose.yml          # PostgreSQL container (dev)
├── docker-compose.prod.yml     # Production containers (postgres, backend, frontend)
├── .github/workflows/          # CI/CD (test.yml, deploy.yml)
└── plan.md                     # Master development plan
```

---

## Architecture Pattern

```
Routes → Services → Prisma (no repository layer)
```

**Backend Layers:**
1. **Routes** (`src/routes/`): HTTP handlers, request validation, response formatting
2. **Services** (`src/services/`): Business logic, authorization, data orchestration
3. **Prisma**: Database queries (no manual repository layer)

**Frontend Layers:**
1. **Components**: Reusable UI elements
2. **Views**: Page-level components
3. **Stores**: Pinia state management
4. **Services**: API client (axios wrapper)

---

## Code Style

### TypeScript Standards

```typescript
// ✅ Good: Explicit types, clear names
interface CreateRecipeRequest {
  name: string;
  description?: string;
  targetVolumeMl: number;
}

async function createRecipe(userId: string, data: CreateRecipeRequest): Promise<Recipe> {
  // Implementation
}

// ❌ Bad: any, unclear names
async function cr(uid: any, d: any): Promise<any> {
  // Implementation
}
```

### Vue Component Structure

```vue
<script setup lang="ts">
// 1. Imports
import { ref, computed } from 'vue';
import type { Recipe } from '@/types/recipe';

// 2. Props with TypeScript interface
interface Props {
  recipe: Recipe;
  compact?: boolean;
}
const props = withDefaults(defineProps<Props>(), {
  compact: false
});

// 3. Emits with types
interface Emits {
  (e: 'click', recipe: Recipe): void;
  (e: 'delete', id: string): void;
}
const emit = defineEmits<Emits>();

// 4. Reactive state
const isHovered = ref(false);

// 5. Computed properties
const displayName = computed(() => {
  return props.compact ? props.recipe.name : `${props.recipe.name} v${props.recipe.activeVersionNumber}`;
});

// 6. Methods
function handleClick() {
  emit('click', props.recipe);
}
</script>

<template>
  <!-- Use semantic HTML, clear class names -->
  <div class="recipe-card" @click="handleClick">
    <h3>{{ displayName }}</h3>
    <p v-if="recipe.description">{{ recipe.description }}</p>
  </div>
</template>

<style scoped>
/* Use design system variables, 8px grid spacing */
.recipe-card {
  padding: 1rem;
  border-radius: 8px;
  background: var(--bg-secondary);
  cursor: pointer;
  transition: background 150ms ease;
}
</style>
```

### Naming Conventions

**Files:**
- Components: PascalCase (`RecipeCard.vue`, `IngredientList.vue`)
- Composables: camelCase with `use` prefix (`useKeyboard.ts`)
- Utilities: kebab-case (`volume.ts`)
- Types: PascalCase (`recipe.ts`)

**Variables/Functions:**
- Variables: camelCase (`isHovered`, `displayName`)
- Functions: camelCase (`createRecipe`, `handleClick`)
- Constants: UPPER_SNAKE_CASE (`MAX_RECIPES`, `DEFAULT_PAGE_SIZE`)
- Classes/Interfaces: PascalCase (`Recipe`, `CreateRecipeRequest`)

**Database:**
- Tables: snake_case (`recipe_ingredients`, `recipe_versions`)
- Columns: snake_case (`target_volume_ml`, `is_active`)

---

## API Design Patterns

### Backend Route Handler

```typescript
// src/routes/recipes.ts
import Router from '@koa/router';
import { authMiddleware } from '../middleware/auth';
import { createRecipeSchema } from '../schemas/recipe';
import * as recipeService from '../services/recipe';

const router = new Router();

router.post('/api/recipes', authMiddleware, async (ctx) => {
  // 1. Get user from auth middleware
  const userId = ctx.state.user.userId;

  // 2. Validate request body
  const body = createRecipeSchema.parse(ctx.request.body);

  // 3. Call service
  const recipe = await recipeService.createRecipe(userId, body);

  // 4. Return response
  ctx.status = 201;
  ctx.body = recipe;
});

export default router;
```

### Backend Service

```typescript
// src/services/recipe.ts
import { prisma } from '../utils/prisma';
import type { CreateRecipeRequest } from '../types/recipe';

export async function createRecipe(userId: string, data: CreateRecipeRequest) {
  // 1. Validation
  const exists = await prisma.recipe.findFirst({
    where: { userId, name: data.name }
  });
  if (exists) {
    throw new Error('Recipe name already exists');
  }

  // 2. Create with user isolation
  const recipe = await prisma.recipe.create({
    data: {
      userId,
      ...data,
      status: 'draft'
    }
  });

  // 3. Return data
  return recipe;
}
```

### Frontend API Service

```typescript
// src/services/recipe.ts
import axios from 'axios';
import type { Recipe, CreateRecipeRequest } from '@/types/recipe';

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000';

export async function createRecipe(data: CreateRecipeRequest): Promise<Recipe> {
  const response = await axios.post<Recipe>(`${API_BASE}/api/recipes`, data);
  return response.data;
}

export async function getRecipes(): Promise<Recipe[]> {
  const response = await axios.get<Recipe[]>(`${API_BASE}/api/recipes`);
  return response.data;
}
```

### Frontend Pinia Store

```typescript
// src/stores/recipe.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import * as recipeApi from '@/services/recipe';
import type { Recipe } from '@/types/recipe';

export const useRecipeStore = defineStore('recipe', () => {
  // State
  const recipes = ref<Recipe[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Getters
  const draftRecipes = computed(() =>
    recipes.value.filter(r => r.status === 'draft')
  );

  // Actions
  async function fetchRecipes() {
    loading.value = true;
    error.value = null;
    try {
      recipes.value = await recipeApi.getRecipes();
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch recipes';
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function createRecipe(data: CreateRecipeRequest) {
    const recipe = await recipeApi.createRecipe(data);
    recipes.value.push(recipe);
    return recipe;
  }

  return {
    // State
    recipes,
    loading,
    error,
    // Getters
    draftRecipes,
    // Actions
    fetchRecipes,
    createRecipe
  };
});
```

---

## Database Schema (Key Models)

### User & Auth
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
```

### Accord System
```prisma
model Accord {
  id                  String   @id @default(uuid())
  userId              String
  name                String
  pyramidPosition     String   // 'top' | 'middle' | 'base'
  volumeMl            Float
  supplier            String?
  notes               String?
  createdAt           DateTime @default(now())
  updatedAt           DateTime @updatedAt

  @@unique([userId, name, pyramidPosition])
}
```

### Recipe System
```prisma
model Recipe {
  id              String    @id @default(uuid())
  userId          String
  name            String
  description     String?
  targetVolumeMl  Float
  status          String    @default("draft")
  activeVersionId String?
  createdAt       DateTime  @default(now())
  updatedAt       DateTime  @updatedAt

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

  @@unique([recipeId, versionNumber])
}

model RecipeIngredient {
  id         String   @id @default(uuid())
  versionId  String
  accordId   String   // ON DELETE RESTRICT (protect accords)
  quantityMl Float
  percentage Float?
  notes      String?
  createdAt  DateTime @default(now())

  @@unique([versionId, accordId])
}
```

---

## Testing Patterns

### Integration Test (Backend)

```typescript
// tests/recipes.test.ts
import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import request from 'supertest';
import { app } from '../src/app';
import { prisma } from '../src/utils/prisma';

describe('POST /api/recipes', () => {
  let accessToken: string;
  let userId: string;

  beforeAll(async () => {
    // Create test user and login
    const res = await request(app)
      .post('/api/auth/register')
      .send({
        email: 'test@example.com',
        username: 'testuser',
        password: 'password123',
        invitationCode: 'test-code'
      });

    accessToken = res.body.accessToken;
    userId = res.body.user._id;
  });

  afterAll(async () => {
    // Cleanup
    await prisma.recipe.deleteMany({ where: { userId } });
    await prisma.user.delete({ where: { id: userId } });
  });

  it('creates a recipe', async () => {
    const res = await request(app)
      .post('/api/recipes')
      .set('Authorization', `Bearer ${accessToken}`)
      .send({
        name: 'Test Recipe',
        description: 'A test recipe',
        targetVolumeMl: 50
      })
      .expect(201);

    expect(res.body.name).toBe('Test Recipe');
    expect(res.body.userId).toBe(userId);
    expect(res.body.status).toBe('draft');
  });

  it('rejects duplicate names', async () => {
    // Create first recipe
    await request(app)
      .post('/api/recipes')
      .set('Authorization', `Bearer ${accessToken}`)
      .send({ name: 'Duplicate', targetVolumeMl: 50 });

    // Try to create duplicate
    const res = await request(app)
      .post('/api/recipes')
      .set('Authorization', `Bearer ${accessToken}`)
      .send({ name: 'Duplicate', targetVolumeMl: 50 })
      .expect(409);

    expect(res.body.error).toBeDefined();
  });
});
```

---

## Security Best Practices

### Authentication

```typescript
// ✅ Always verify user ownership
async function getRecipe(userId: string, recipeId: string) {
  const recipe = await prisma.recipe.findFirst({
    where: { id: recipeId, userId } // ← User isolation
  });
  if (!recipe) throw new Error('Recipe not found');
  return recipe;
}

// ❌ Never trust client-provided user IDs
async function getRecipeBad(recipeId: string) {
  // Missing userId check - security vulnerability!
  return await prisma.recipe.findUnique({ where: { id: recipeId } });
}
```

### Input Validation

```typescript
// ✅ Validate with Zod
import { z } from 'zod';

const createRecipeSchema = z.object({
  name: z.string().min(1).max(255),
  description: z.string().optional(),
  targetVolumeMl: z.number().positive()
});

router.post('/api/recipes', async (ctx) => {
  const body = createRecipeSchema.parse(ctx.request.body); // ← Validates
  // ...
});

// ❌ Never trust unvalidated input
router.post('/api/recipes', async (ctx) => {
  const body = ctx.request.body; // ← Dangerous!
  await createRecipe(body);
});
```

### Password Handling

```typescript
// ✅ Hash passwords
import bcrypt from 'bcryptjs';

const passwordHash = await bcrypt.hash(password, 10);

// ❌ Never store plain passwords
user.password = password; // ← Never do this!
```

---

## Design System (Notion-Inspired)

### Color Palette

```css
/* Text colors */
--text-primary: #37352F;
--text-secondary: #787774;
--text-tertiary: #9B9A97;

/* Background colors */
--bg-primary: #FFFFFF;
--bg-secondary: #FAFAFA;
--bg-tertiary: #F7F6F3;

/* Border */
--border-default: #E9E9E7;

/* Accent */
--accent-primary: #0F766E;  /* Teal for actions */

/* Pyramid positions */
--pyramid-top: linear-gradient(135deg, #FFD93D, #FFA800);
--pyramid-middle: linear-gradient(135deg, #B565D8, #8B5CF6);
--pyramid-base: linear-gradient(135deg, #A0826D, #6B4423);
```

### Spacing (8px Grid)

```css
/* Use multiples of 8px for consistency */
--spacing-xs: 0.5rem;   /* 8px */
--spacing-sm: 1rem;     /* 16px */
--spacing-md: 1.5rem;   /* 24px */
--spacing-lg: 2rem;     /* 32px */
--spacing-xl: 3rem;     /* 48px */
--spacing-2xl: 4rem;    /* 64px */
```

### Component Patterns

```vue
<!-- Card with hover effect -->
<div class="card">
  <h3>Title</h3>
  <p>Description</p>
</div>

<style scoped>
.card {
  padding: 1rem;
  border-radius: 8px;
  background: var(--bg-secondary);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12);
  transition: background 150ms ease;
  cursor: pointer;
}

.card:hover {
  background: var(--bg-tertiary);
}
</style>
```

---

## Important Rules

### Do's ✅

1. **Testing**: Always write tests before committing
2. **Type Safety**: Use TypeScript strictly, avoid `any`
3. **User Isolation**: Filter all queries by `userId`
4. **Error Handling**: Provide clear, actionable error messages
5. **Documentation**: Keep all documentation up-to-date. When you change code, update any affected docs (README, agent instructions, specs, QUICKSTART, etc.) in the same commit or immediately after.
6. **Commit Frequently**: Make small, focused commits as you work. Don't accumulate large uncommitted changesets. Each logical change should be its own commit.
7. **Security**: Validate all inputs, never expose secrets
8. **Patterns**: Follow established code patterns
9. **Simplicity**: Keep solutions simple, avoid over-engineering

### Don'ts ❌

1. **No Over-Engineering**: Don't add unrequested features
2. **No `any` Type**: Use explicit types always
3. **No Skipping Tests**: Never commit untested code
4. **No Commented Code**: Remove code, don't comment out
5. **No Breaking Changes**: Without approval first
6. **No Hardcoded Secrets**: Use environment variables
7. **No Production Data**: Use test database
8. **No Force Push**: To main/master branch

### Recipe System Rules

1. **Volume Validation**: Disabled by default (user opt-in)
2. **Version Auto-Activation**: New versions automatically active
3. **Accord Protection**: Prevent deletion if used in recipes (409 error)
4. **Immutable Versions**: Once created, cannot edit (duplicate instead)
5. **Recipe Status**: draft | in_progress | tested | finalized | archived

---

## Common Commands

### Development

```bash
# Start PostgreSQL
docker compose up -d

# Backend (terminal 1)
cd backend
npm run dev

# Frontend (terminal 2)
cd frontend
npm run dev

# Run tests
cd backend
npm test

# Database operations
npx prisma studio            # Visual database browser
npx prisma generate          # Update Prisma client
npx prisma db push           # Apply schema changes
```

### API Testing

```bash
# Health check
curl http://localhost:3000/api/health

# Login (get access token)
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Create recipe
curl -X POST http://localhost:3000/api/recipes \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Recipe","targetVolumeMl":50}'
```

---

## Useful Links

**Documentation:**
- Master Plan: `plan.md` (41KB, comprehensive roadmap)
- API Spec: `specs/api-spec.md` (Accord endpoints)
- Recipe Spec: `specs/recipe-system.md` (Phase 10 details)
- Database: `specs/data-models.md`

**Tech Docs:**
- Koa.js: https://koajs.com/
- Prisma: https://www.prisma.io/docs
- Vue 3: https://vuejs.org/guide/
- Pinia: https://pinia.vuejs.org/
- Naive UI: https://www.naiveui.com/
- Zod: https://zod.dev/

**Project:**
- Repository: https://github.com/xupit3r/scentora
- Issues: https://github.com/xupit3r/scentora/issues

---

## Current Focus (February 7, 2026)

### Phase 10.9: Recipe Views & Pages

**Goal**: Create page-level views and integrate routing

**Tasks**:
- [ ] RecipesView.vue - Main recipes list page
- [ ] RecipeDetailView.vue - Single recipe detail page
- [ ] RecipeCreateView.vue - Create recipe flow
- [ ] RecipeEditView.vue - Edit recipe flow
- [ ] CollectionsView.vue - Collections list
- [ ] CollectionDetailView.vue - Collection with recipes
- [ ] Update router/index.ts with routes
- [ ] Add route guards for auth
- [ ] Implement breadcrumbs
- [ ] Test navigation

**Next Phases**:
- Phase 10.10: Features & polish (search, filtering, export/import)
- Phase 10.11: Testing (integration, E2E)

---

**Last Updated**: February 7, 2026
**Test Status**: 70/70 passing (100%)
**Deployment**: Live at https://scentora.thejoeshow.net
**Current Phase**: 10.9 (Recipe Views & Pages)
