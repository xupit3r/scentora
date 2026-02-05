import Router from '@koa/router';
import { jwtAuth } from '../middleware/auth.js';
import { createAccordSchema, updateAccordSchema, addTagSchema } from '../schemas/accord.schema.js';
import * as accordService from '../services/accord.service.js';
import { ValidationError } from '../utils/errors.js';

const router = new Router({ prefix: '/api/accords' });
router.use(jwtAuth());

// POST /api/accords
router.post('/', async (ctx) => {
  const parsed = createAccordSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  const accord = await accordService.createAccord(ctx.state.userId, parsed.data);
  ctx.status = 201;
  ctx.body = { accord };
});

// GET /api/accords
router.get('/', async (ctx) => {
  const { position, minVolume, maxVolume, supplier, search } = ctx.query;
  const tags = ctx.query.tags;

  const hasFilters = position || minVolume || maxVolume || supplier || search || tags;

  if (hasFilters) {
    const tagArray = tags
      ? Array.isArray(tags) ? tags as string[] : [tags as string]
      : undefined;
    const accords = await accordService.searchAccords(ctx.state.userId, {
      position: position as string | undefined,
      minVolume: minVolume ? parseFloat(minVolume as string) : undefined,
      maxVolume: maxVolume ? parseFloat(maxVolume as string) : undefined,
      supplier: supplier as string | undefined,
      search: search as string | undefined,
      tags: tagArray,
    });
    ctx.body = { accords };
  } else {
    const accords = await accordService.listAccords(ctx.state.userId);
    ctx.body = { accords };
  }
});

// GET /api/accords/:id
router.get('/:id', async (ctx) => {
  const accord = await accordService.getAccord(ctx.params.id, ctx.state.userId);
  ctx.body = { accord };
});

// PUT /api/accords/:id
router.put('/:id', async (ctx) => {
  const parsed = updateAccordSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  const accord = await accordService.updateAccord(ctx.params.id, ctx.state.userId, parsed.data);
  ctx.body = { accord };
});

// DELETE /api/accords/:id
router.delete('/:id', async (ctx) => {
  await accordService.deleteAccord(ctx.params.id, ctx.state.userId);
  ctx.body = { message: 'Accord deleted successfully' };
});

// POST /api/accords/:id/tags
router.post('/:id/tags', async (ctx) => {
  const parsed = addTagSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  await accordService.addTag(ctx.params.id, ctx.state.userId, parsed.data.tag);
  ctx.body = { message: 'Tag added successfully' };
});

// DELETE /api/accords/:id/tags/:tag
router.delete('/:id/tags/:tag', async (ctx) => {
  await accordService.removeTag(ctx.params.id, ctx.state.userId, ctx.params.tag);
  ctx.body = { message: 'Tag removed successfully' };
});

export default router;
