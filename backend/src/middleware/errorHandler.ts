import type Koa from 'koa';
import { AppError } from '../utils/errors.js';

export function errorHandler(): Koa.Middleware {
  return async (ctx, next) => {
    try {
      await next();
    } catch (err: unknown) {
      if (err instanceof AppError) {
        ctx.status = err.statusCode;
        ctx.body = {
          error: {
            message: err.message,
            ...(err.details !== undefined ? { details: err.details } : {}),
          },
        };
        return;
      }

      // Prisma unique constraint violation
      if (
        err &&
        typeof err === 'object' &&
        'code' in err &&
        (err as { code: string }).code === 'P2002'
      ) {
        ctx.status = 409;
        ctx.body = {
          error: { message: 'A record with that value already exists' },
        };
        return;
      }

      // Prisma foreign key constraint violation
      if (
        err &&
        typeof err === 'object' &&
        'code' in err &&
        (err as { code: string }).code === 'P2003'
      ) {
        ctx.status = 409;
        ctx.body = {
          error: { message: 'Cannot delete: record is referenced by other records' },
        };
        return;
      }

      console.error('Unhandled error:', err);
      ctx.status = 500;
      ctx.body = {
        error: { message: 'Internal server error' },
      };
    }
  };
}
