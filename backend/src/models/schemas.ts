import { z } from 'zod';

export const perfumePyramidSchema = z.object({
  top: z.array(z.string()).min(0),
  middle: z.array(z.string()).min(0),
  base: z.array(z.string()).min(0),
});

export const createPerfumeSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  designer: z.string().min(1, 'Designer is required'),
  year: z.number().int().min(1800).max(new Date().getFullYear() + 1).optional(),
  concentration: z.string().optional(),
  pyramid: perfumePyramidSchema,
  description: z.string().optional(),
  imageUrl: z.string().url().optional().or(z.literal('')),
});

export const updatePerfumeSchema = createPerfumeSchema.partial();

export const journalEntrySchema = z.object({
  perfumeId: z.string().min(1),
  date: z.string(),
  content: z.string().min(1, 'Content is required'),
  rating: z.number().min(1).max(10).optional(),
  occasion: z.string().optional(),
  weather: z.string().optional(),
});

export const updateJournalEntrySchema = journalEntrySchema.omit({ perfumeId: true }).partial();
