import Router from '@koa/router';
import { jwtAuth } from '../middleware/auth.js';
import * as statsService from '../services/stats.service.js';

const router = new Router();

// GET /api/stats
router.get('/api/stats', jwtAuth(), async (ctx) => {
  const stats = await statsService.getStatistics(ctx.state.userId);
  ctx.body = stats;
});

export default router;
