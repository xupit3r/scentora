import { z } from 'zod';

export const createInvitationSchema = z.object({
  email: z.string().email().optional(),
  expiresInDays: z.number().int().min(1).max(365).optional(),
});
