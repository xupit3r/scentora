import { Context } from 'koa';
import { z } from 'zod';

export function validateBody(schema: z.ZodSchema) {
  return async (ctx: Context, next: () => Promise<any>) => {
    try {
      const validated = schema.parse(ctx.request.body);
      ctx.request.body = validated;
      await next();
    } catch (error) {
      if (error instanceof z.ZodError) {
        ctx.status = 400;
        ctx.body = {
          error: {
            message: 'Validation failed',
            details: error.issues,
          },
        };
        return;
      }
      throw error;
    }
  };
}
