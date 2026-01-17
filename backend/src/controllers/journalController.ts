import { Context } from 'koa';
import { db } from '../config/database';
import type { JournalEntry } from '../models/types';

export const journalController = {
  async getByPerfumeId(ctx: Context) {
    try {
      const { perfumeId } = ctx.params;
      
      const result = await db.find({
        selector: { 
          type: 'journal',
          perfumeId,
        },
        sort: [{ date: 'desc' }],
      });
      
      ctx.body = result.docs;
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async create(ctx: Context) {
    try {
      const { perfumeId } = ctx.params;
      const entryData = ctx.request.body as Omit<JournalEntry, '_id' | '_rev' | 'type' | 'createdAt' | 'updatedAt'>;
      
      // Verify perfume exists
      try {
        const perfume = await db.get(perfumeId) as any;
        if (perfume.type !== 'perfume') {
          ctx.status = 404;
          ctx.body = { error: { message: 'Perfume not found' } };
          return;
        }
      } catch (error: any) {
        if (error.statusCode === 404) {
          ctx.status = 404;
          ctx.body = { error: { message: 'Perfume not found' } };
          return;
        }
        throw error;
      }

      const entry: JournalEntry = {
        ...entryData,
        perfumeId,
        type: 'journal',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      const result = await db.insert(entry);
      const created = await db.get(result.id);
      
      ctx.status = 201;
      ctx.body = created;
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async update(ctx: Context) {
    try {
      const { id } = ctx.params;
      const updates = ctx.request.body as Partial<JournalEntry>;
      
      const existing = await db.get(id) as any;
      
      if (existing.type !== 'journal') {
        ctx.status = 404;
        ctx.body = { error: { message: 'Journal entry not found' } };
        return;
      }

      const updated: JournalEntry = {
        ...(existing as JournalEntry),
        ...updates,
        type: 'journal',
        _id: existing._id,
        _rev: existing._rev,
        perfumeId: existing.perfumeId,
        createdAt: existing.createdAt,
        updatedAt: new Date().toISOString(),
      };

      await db.insert(updated);
      const result = await db.get(id);
      
      ctx.body = result;
    } catch (error: any) {
      if (error.statusCode === 404) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Journal entry not found' } };
      } else {
        ctx.status = 500;
        ctx.body = { error: { message: error.message } };
      }
    }
  },

  async delete(ctx: Context) {
    try {
      const { id } = ctx.params;
      const entry = await db.get(id) as any;
      
      if (entry.type !== 'journal') {
        ctx.status = 404;
        ctx.body = { error: { message: 'Journal entry not found' } };
        return;
      }

      await db.destroy(id, entry._rev);
      
      ctx.status = 204;
    } catch (error: any) {
      if (error.statusCode === 404) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Journal entry not found' } };
      } else {
        ctx.status = 500;
        ctx.body = { error: { message: error.message } };
      }
    }
  },
};
