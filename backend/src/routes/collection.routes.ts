import Router from '@koa/router';
import { jwtAuth } from '../middleware/auth.js';
import { createCollectionSchema, updateCollectionSchema, addRecipeToCollectionSchema } from '../schemas/recipe.schema.js';
import * as collectionService from '../services/recipeCollection.service.js';
import { ValidationError } from '../utils/errors.js';

const router = new Router({ prefix: '/api/collections' });
router.use(jwtAuth());

// POST /api/collections
router.post('/', async (ctx) => {
  const parsed = createCollectionSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  const collection = await collectionService.createCollection(ctx.state.userId, parsed.data);
  ctx.status = 201;
  ctx.body = collection;
});

// GET /api/collections
router.get('/', async (ctx) => {
  const limit = ctx.query.limit ? parseInt(ctx.query.limit as string, 10) : 100;
  const offset = ctx.query.offset ? parseInt(ctx.query.offset as string, 10) : 0;
  ctx.body = await collectionService.listCollections(ctx.state.userId, limit, offset);
});

// GET /api/collections/:id
router.get('/:id', async (ctx) => {
  ctx.body = await collectionService.getCollection(ctx.params.id, ctx.state.userId);
});

// PUT /api/collections/:id
router.put('/:id', async (ctx) => {
  const parsed = updateCollectionSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  ctx.body = await collectionService.updateCollection(ctx.params.id, ctx.state.userId, parsed.data);
});

// DELETE /api/collections/:id
router.delete('/:id', async (ctx) => {
  await collectionService.deleteCollection(ctx.params.id, ctx.state.userId);
  ctx.body = { message: 'Collection deleted successfully' };
});

// POST /api/collections/:id/recipes
router.post('/:id/recipes', async (ctx) => {
  const parsed = addRecipeToCollectionSchema.safeParse(ctx.request.body);
  if (!parsed.success) throw new ValidationError('Validation failed', parsed.error.format());
  await collectionService.addRecipe(ctx.params.id, parsed.data.recipeId, ctx.state.userId);
  ctx.body = { message: 'Recipe added to collection' };
});

// DELETE /api/collections/:id/recipes/:recipeId
router.delete('/:id/recipes/:recipeId', async (ctx) => {
  await collectionService.removeRecipe(ctx.params.id, ctx.params.recipeId, ctx.state.userId);
  ctx.body = { message: 'Recipe removed from collection' };
});

export default router;
