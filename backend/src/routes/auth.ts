import Router from '@koa/router';
import { authController } from '../controllers/authController';
import { validateBody } from '../middleware/validation';
import { registerSchema, loginSchema } from '../models/schemas';
import { authenticate } from '../middleware/auth';
import { authRateLimiter } from '../middleware/rateLimit';

const router = new Router({
  prefix: '/auth',
});

// Public routes with rate limiting
router.post('/register', authRateLimiter, validateBody(registerSchema), authController.register);
router.post('/login', authRateLimiter, validateBody(loginSchema), authController.login);
router.post('/refresh', authRateLimiter, authController.refresh);
router.post('/logout', authController.logout);

// Protected routes
router.post('/logout-all', authenticate, authController.logoutAll);
router.get('/me', authenticate, authController.me);

export default router;
