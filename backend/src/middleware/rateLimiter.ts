import type Koa from 'koa';

interface RateLimitEntry {
  count: number;
  resetAt: number;
}

function createRateLimiter(maxRequests: number, windowMs: number): Koa.Middleware {
  // Disable rate limiting in test environment
  if (process.env.NODE_ENV === 'test') {
    return async (_ctx, next) => next();
  }

  const store = new Map<string, RateLimitEntry>();

  // Cleanup expired entries every minute
  setInterval(() => {
    const now = Date.now();
    for (const [key, entry] of store) {
      if (now > entry.resetAt) store.delete(key);
    }
  }, 60_000).unref();

  return async (ctx, next) => {
    const key = ctx.ip;
    const now = Date.now();

    let entry = store.get(key);
    if (!entry || now > entry.resetAt) {
      entry = { count: 0, resetAt: now + windowMs };
      store.set(key, entry);
    }

    entry.count++;

    if (entry.count > maxRequests) {
      ctx.status = 429;
      ctx.body = {
        error: { message: 'Too many requests, please try again later' },
      };
      return;
    }

    await next();
  };
}

/** General rate limiter: 100 requests per minute */
export function generalRateLimiter(): Koa.Middleware {
  return createRateLimiter(100, 60 * 1000);
}
