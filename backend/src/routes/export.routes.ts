import Router from '@koa/router';
import { jwtAuth } from '../middleware/auth.js';
import * as exportService from '../services/export.service.js';

const router = new Router({ prefix: '/api/export' });
router.use(jwtAuth());

// GET /api/export
router.get('/', async (ctx) => {
  const data = await exportService.exportAccords(ctx.state.userId);
  ctx.set('Content-Type', 'application/json');
  ctx.set('Content-Disposition', 'attachment; filename=scentora-export.json');
  ctx.body = data;
});

// POST /api/export/import
router.post('/import', async (ctx) => {
  const result = await exportService.importAccords(ctx.state.userId, ctx.request.body as any);
  ctx.body = result;
});

export default router;
