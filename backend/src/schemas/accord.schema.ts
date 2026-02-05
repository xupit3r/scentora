import { z } from 'zod';

export const createAccordSchema = z.object({
  name: z.string().min(1),
  pyramidPosition: z.enum(['top', 'middle', 'base']),
  volumeMl: z.number().gte(0),
  supplier: z.string().optional(),
  purchaseDate: z.string().optional(), // ISO date string
  dilutionPercentage: z.number().gte(0).lte(100).optional(),
  notes: z.string().optional(),
  tags: z.array(z.string()).optional(),
});

export const updateAccordSchema = z.object({
  name: z.string().min(1).optional(),
  pyramidPosition: z.enum(['top', 'middle', 'base']).optional(),
  volumeMl: z.number().gte(0).optional(),
  supplier: z.string().nullable().optional(),
  purchaseDate: z.string().nullable().optional(),
  dilutionPercentage: z.number().gte(0).lte(100).nullable().optional(),
  notes: z.string().nullable().optional(),
  tags: z.array(z.string()).optional(),
});

export const addTagSchema = z.object({
  tag: z.string().min(1).max(50),
});
