import { z } from 'zod';

export const createRecipeSchema = z.object({
  name: z.string().min(1),
  description: z.string().optional(),
  targetVolumeMl: z.number().gt(0),
  status: z.enum(['draft', 'in_progress', 'tested', 'finalized', 'archived']).optional(),
});

export const updateRecipeSchema = z.object({
  name: z.string().min(1).optional(),
  description: z.string().nullable().optional(),
  targetVolumeMl: z.number().gt(0).optional(),
  status: z.enum(['draft', 'in_progress', 'tested', 'finalized', 'archived']).optional(),
});

export const createVersionSchema = z.object({
  name: z.string().min(1).optional(),
  notes: z.string().optional(),
});

export const addIngredientSchema = z.object({
  accordId: z.string().uuid(),
  quantityMl: z.number().gt(0),
  percentage: z.number().gte(0).lte(100).optional(),
  notes: z.string().optional(),
});

export const updateIngredientSchema = z.object({
  quantityMl: z.number().gt(0).optional(),
  percentage: z.number().gte(0).lte(100).nullable().optional(),
  notes: z.string().nullable().optional(),
});

export const createNoteSchema = z.object({
  versionId: z.string().uuid().optional(),
  content: z.string().min(1),
  noteType: z.enum(['general', 'testing', 'observation', 'adjustment', 'reminder']).optional(),
});

export const updateNoteSchema = z.object({
  content: z.string().min(1),
});

export const addRecipeTagSchema = z.object({
  tag: z.string().min(1).max(50),
});

export const createCollectionSchema = z.object({
  name: z.string().min(1),
  description: z.string().optional(),
});

export const updateCollectionSchema = z.object({
  name: z.string().min(1).optional(),
  description: z.string().nullable().optional(),
});

export const addRecipeToCollectionSchema = z.object({
  recipeId: z.string().uuid(),
});
