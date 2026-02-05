import Koa from 'koa';
import cors from '@koa/cors';
import bodyParser from 'koa-bodyparser';
import { config } from './config.js';
import { errorHandler } from './middleware/errorHandler.js';
import { generalRateLimiter } from './middleware/rateLimiter.js';
import { mountRoutes } from './routes/index.js';

export function createApp() {
  const app = new Koa();

  // Global error handler (first middleware)
  app.use(errorHandler());

  // CORS
  app.use(cors({
    origin: (ctx) => {
      const origin = ctx.get('Origin');
      if (config.corsAllowedOrigins.includes(origin)) return origin;
      return '';
    },
    credentials: true,
  }));

  // Body parser
  app.use(bodyParser());

  // General rate limiter
  app.use(generalRateLimiter());

  // Health check
  app.use(async (ctx, next) => {
    if (ctx.path === '/api/health' && ctx.method === 'GET') {
      ctx.body = {
        status: 'ok',
        timestamp: new Date().toISOString(),
        service: 'scentora-api',
      };
      return;
    }
    await next();
  });

  // Mount all route groups
  mountRoutes(app);

  return app;
}
