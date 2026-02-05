import Router from '@koa/router';
import { jwtAuth } from '../middleware/auth.js';
import {
  createRecipeSchema, updateRecipeSchema,
  createVersionSchema, addIngredientSchema, updateIngredientSchema,
  createNoteSchema, updateNoteSchema, addRecipeTagSchema,
} from '../schemas/recipe.schema.js';
import * as recipeService from '../services/recipe.service.js';
import * as versionService from '../services/recipeVersion.service.js';
import * as ingredientService from '../services/recipeIngredient.service.js';
import * as noteService from '../services/recipeNote.service.js';
import * as tagService from '../services/recipeTag.service.js';
import { ValidationError } from '../utils/errors.js';

const router = new Router({ prefix: '/api/recipes' });
router.use(jwtAuth());

// POST /api/recipes
router.post('/', async (ctx) => {
  const parsed = createRecipeSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  const recipe = await recipeService.createRecipe(ctx.state.userId, parsed.data);
  ctx.status = 201;
  ctx.body = recipe;
});

// GET /api/recipes
router.get('/', async (ctx) => {
  const status = ctx.query.status as string | undefined;
  const limit = ctx.query.limit ? parseInt(ctx.query.limit as string, 10) : 100;
  const offset = ctx.query.offset ? parseInt(ctx.query.offset as string, 10) : 0;
  ctx.body = await recipeService.listRecipes(ctx.state.userId, status, limit, offset);
});

// GET /api/recipes/search
router.get('/search', async (ctx) => {
  const q = ctx.query.q as string;
  if (!q) {
    ctx.status = 400;
    ctx.body = { error: { message: 'Search query required' } };
    return;
  }
  const limit = ctx.query.limit ? parseInt(ctx.query.limit as string, 10) : 100;
  const offset = ctx.query.offset ? parseInt(ctx.query.offset as string, 10) : 0;
  ctx.body = await recipeService.searchRecipes(ctx.state.userId, q, limit, offset);
});

// GET /api/recipes/tags/popular
router.get('/tags/popular', async (ctx) => {
  const limit = ctx.query.limit ? parseInt(ctx.query.limit as string, 10) : 20;
  ctx.body = await tagService.getPopularTags(ctx.state.userId, limit);
});

// GET /api/recipes/:id
router.get('/:id', async (ctx) => {
  ctx.body = await recipeService.getRecipe(ctx.params.id, ctx.state.userId);
});

// PUT /api/recipes/:id
router.put('/:id', async (ctx) => {
  const parsed = updateRecipeSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  ctx.body = await recipeService.updateRecipe(ctx.params.id, ctx.state.userId, parsed.data);
});

// DELETE /api/recipes/:id
router.delete('/:id', async (ctx) => {
  await recipeService.deleteRecipe(ctx.params.id, ctx.state.userId);
  ctx.body = { message: 'Recipe deleted successfully' };
});

// --- Versions ---

// POST /api/recipes/:id/versions
router.post('/:id/versions', async (ctx) => {
  const parsed = createVersionSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  const version = await versionService.createVersion(ctx.params.id, ctx.state.userId, parsed.data);
  ctx.status = 201;
  ctx.body = version;
});

// GET /api/recipes/:id/versions
router.get('/:id/versions', async (ctx) => {
  ctx.body = await versionService.listVersions(ctx.params.id, ctx.state.userId);
});

// GET /api/recipes/:id/versions/:versionId
router.get('/:id/versions/:versionId', async (ctx) => {
  ctx.body = await versionService.getVersion(ctx.params.versionId);
});

// POST /api/recipes/:id/versions/:versionNumber/activate
router.post('/:id/versions/:versionNumber/activate', async (ctx) => {
  const versionNumber = parseInt(ctx.params.versionNumber, 10);
  if (isNaN(versionNumber)) {
    ctx.status = 400;
    ctx.body = { error: { message: 'Invalid version number' } };
    return;
  }
  ctx.body = await versionService.activateVersion(ctx.params.id, versionNumber, ctx.state.userId);
});

// POST /api/recipes/:id/versions/:versionId/duplicate
router.post('/:id/versions/:versionId/duplicate', async (ctx) => {
  const version = await versionService.duplicateVersion(ctx.params.versionId, ctx.state.userId);
  ctx.status = 201;
  ctx.body = version;
});

// --- Ingredients ---

// POST /api/recipes/:id/versions/:versionId/ingredients
router.post('/:id/versions/:versionId/ingredients', async (ctx) => {
  const parsed = addIngredientSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  const ingredient = await ingredientService.addIngredient(
    ctx.params.versionId, ctx.state.userId, parsed.data,
  );
  ctx.status = 201;
  ctx.body = ingredient;
});

// PUT /api/recipes/:id/versions/:versionId/ingredients/:ingredientId
router.put('/:id/versions/:versionId/ingredients/:ingredientId', async (ctx) => {
  const parsed = updateIngredientSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  ctx.body = await ingredientService.updateIngredient(
    ctx.params.ingredientId, ctx.state.userId, parsed.data,
  );
});

// DELETE /api/recipes/:id/versions/:versionId/ingredients/:ingredientId
router.delete('/:id/versions/:versionId/ingredients/:ingredientId', async (ctx) => {
  await ingredientService.removeIngredient(ctx.params.ingredientId, ctx.state.userId);
  ctx.body = { message: 'Ingredient removed successfully' };
});

// --- Notes ---

// POST /api/recipes/:id/notes
router.post('/:id/notes', async (ctx) => {
  const parsed = createNoteSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  const note = await noteService.createNote(ctx.params.id, ctx.state.userId, parsed.data);
  ctx.status = 201;
  ctx.body = note;
});

// GET /api/recipes/:id/notes
router.get('/:id/notes', async (ctx) => {
  ctx.body = await noteService.listNotes(ctx.params.id, ctx.state.userId);
});

// PUT /api/recipes/:id/notes/:noteId
router.put('/:id/notes/:noteId', async (ctx) => {
  const parsed = updateNoteSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  ctx.body = await noteService.updateNote(ctx.params.noteId, ctx.state.userId, parsed.data);
});

// DELETE /api/recipes/:id/notes/:noteId
router.delete('/:id/notes/:noteId', async (ctx) => {
  await noteService.deleteNote(ctx.params.noteId, ctx.state.userId);
  ctx.body = { message: 'Note deleted successfully' };
});

// --- Tags ---

// POST /api/recipes/:id/tags
router.post('/:id/tags', async (ctx) => {
  const parsed = addRecipeTagSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  await tagService.addTag(ctx.params.id, ctx.state.userId, parsed.data.tag);
  ctx.status = 201;
  ctx.body = { message: 'Tag added successfully', tag: parsed.data.tag };
});

// DELETE /api/recipes/:id/tags/:tag
router.delete('/:id/tags/:tag', async (ctx) => {
  await tagService.removeTag(ctx.params.id, ctx.state.userId, ctx.params.tag);
  ctx.body = { message: 'Tag removed successfully' };
});

export default router;
