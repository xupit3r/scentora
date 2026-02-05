import Router from '@koa/router';
import { jwtAuth } from '../middleware/auth.js';
import { authRateLimiter } from '../middleware/rateLimiter.js';
import { registerSchema, loginSchema, refreshSchema } from '../schemas/auth.schema.js';
import * as authService from '../services/auth.service.js';
import { ValidationError } from '../utils/errors.js';

const router = new Router({ prefix: '/api/auth' });
const authLimiter = authRateLimiter();

// POST /api/auth/register
router.post('/register', authLimiter, async (ctx) => {
  const parsed = registerSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  const { email, username, password, invitationCode } = parsed.data;
  const result = await authService.register(email, username, password, invitationCode);
  ctx.status = 201;
  ctx.body = result;
});

// POST /api/auth/login
router.post('/login', authLimiter, async (ctx) => {
  const parsed = loginSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  const result = await authService.login(parsed.data.email, parsed.data.password);
  ctx.status = 200;
  ctx.body = result;
});

// POST /api/auth/refresh
router.post('/refresh', authLimiter, async (ctx) => {
  const parsed = refreshSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  const result = await authService.refresh(parsed.data.refreshToken);
  ctx.status = 200;
  ctx.body = result;
});

// POST /api/auth/logout
router.post('/logout', authLimiter, async (ctx) => {
  const parsed = refreshSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  await authService.logout(parsed.data.refreshToken);
  ctx.status = 204;
  ctx.body = null;
});

// GET /api/auth/me (protected)
router.get('/me', jwtAuth(), async (ctx) => {
  const user = await authService.getMe(ctx.state.userId);
  ctx.body = user;
});

// POST /api/auth/logout-all (protected)
router.post('/logout-all', jwtAuth(), async (ctx) => {
  await authService.logoutAll(ctx.state.userId);
  ctx.body = { message: 'Logged out from all devices' };
});

export default router;
