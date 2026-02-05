import Router from '@koa/router';
import * as tagService from '../services/tag.service.js';

const router = new Router({ prefix: '/api/tags' });

// GET /api/tags
router.get('/', async (ctx) => {
  const tags = await tagService.getAllTags();
  ctx.body = { tags };
});

// GET /api/tags/search
router.get('/search', async (ctx) => {
  const q = (ctx.query.q as string) || '';
  const tags = await tagService.searchTags(q);
  ctx.body = { tags };
});

// GET /api/tags/categories
router.get('/categories', async (ctx) => {
  const categories = await tagService.getCategories();
  ctx.body = { categories };
});

// GET /api/tags/grouped
router.get('/grouped', async (ctx) => {
  const tags = await tagService.getTagsGrouped();
  ctx.body = { tags };
});

// GET /api/tags/category/:category
router.get('/category/:category', async (ctx) => {
  const tags = await tagService.getByCategory(ctx.params.category);
  ctx.body = { tags };
});

export default router;
