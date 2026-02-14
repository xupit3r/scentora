import Router from '@koa/router';
import { jwtAuth } from '../middleware/auth.js';
import { loginSchema } from '../schemas/auth.schema.js';
import * as authService from '../services/auth.service.js';
import { ValidationError } from '../utils/errors.js';

const router = new Router({ prefix: '/api/auth' });

// POST /api/auth/login
router.post('/login', async (ctx) => {
  const parsed = loginSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  const result = await authService.login(parsed.data.password);
  ctx.status = 200;
  ctx.body = result;
});

// GET /api/auth/me (protected)
router.get('/me', jwtAuth(), async (ctx) => {
  const user = await authService.getMe(ctx.state.userId);
  ctx.body = { user };
});

export default router;
