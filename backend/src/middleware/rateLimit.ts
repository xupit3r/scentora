import ratelimit from 'koa-ratelimit';

// In-memory store for rate limiting
// For production, consider using Redis
const db = new Map();

export const authRateLimiter = ratelimit({
  driver: 'memory',
  db,
  duration: 15 * 60 * 1000, // 15 minutes
  errorMessage: 'Too many requests, please try again later',
  id: (ctx) => ctx.ip,
  headers: {
    remaining: 'Rate-Limit-Remaining',
    reset: 'Rate-Limit-Reset',
    total: 'Rate-Limit-Total',
  },
  max: 5, // 5 requests per 15 minutes for auth endpoints
  disableHeader: false,
});

export const generalRateLimiter = ratelimit({
  driver: 'memory',
  db,
  duration: 60 * 1000, // 1 minute
  errorMessage: 'Too many requests, please try again later',
  id: (ctx) => ctx.ip,
  headers: {
    remaining: 'Rate-Limit-Remaining',
    reset: 'Rate-Limit-Reset',
    total: 'Rate-Limit-Total',
  },
  max: 100, // 100 requests per minute for general endpoints
  disableHeader: false,
});
