import { Context, Next } from 'koa';
import { verifyAccessToken } from '../config/auth';
import type { AuthUser } from '../models/types';

// Extend Koa context to include user
declare module 'koa' {
  interface Context {
    user?: AuthUser;
  }
}

export async function authenticate(ctx: Context, next: Next) {
  try {
    const authHeader = ctx.headers.authorization;

    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      ctx.status = 401;
      ctx.body = { error: { message: 'Authentication required' } };
      return;
    }

    const token = authHeader.substring(7);
    const user = verifyAccessToken(token);
    
    ctx.user = user;
    await next();
  } catch (error: any) {
    if (error.name === 'JsonWebTokenError') {
      ctx.status = 401;
      ctx.body = { error: { message: 'Invalid token' } };
    } else if (error.name === 'TokenExpiredError') {
      ctx.status = 401;
      ctx.body = { error: { message: 'Token expired' } };
    } else {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  }
}

export function optionalAuth(ctx: Context, next: Next) {
  const authHeader = ctx.headers.authorization;

  if (authHeader && authHeader.startsWith('Bearer ')) {
    try {
      const token = authHeader.substring(7);
      const user = verifyAccessToken(token);
      ctx.user = user;
    } catch (error) {
      // Ignore errors for optional auth
    }
  }

  return next();
}
