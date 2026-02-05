import type Koa from 'koa';
import jwt from 'jsonwebtoken';
import { config } from '../config.js';

interface JWTPayload {
  userId: string;
  email: string;
}

export function jwtAuth(): Koa.Middleware {
  return async (ctx, next) => {
    const authHeader = ctx.get('Authorization');
    if (!authHeader) {
      ctx.status = 401;
      ctx.body = { error: { message: 'Missing authorization header' } };
      return;
    }

    const parts = authHeader.split(' ');
    if (parts.length !== 2 || parts[0] !== 'Bearer') {
      ctx.status = 401;
      ctx.body = { error: { message: 'Invalid authorization header format' } };
      return;
    }

    try {
      const payload = jwt.verify(parts[1], config.jwtSecret) as JWTPayload;
      ctx.state.userId = payload.userId;
      ctx.state.email = payload.email;
    } catch {
      ctx.status = 401;
      ctx.body = { error: { message: 'Invalid or expired token' } };
      return;
    }

    await next();
  };
}
