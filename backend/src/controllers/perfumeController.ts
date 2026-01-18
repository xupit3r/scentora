import { Context } from 'koa';
import { db } from '../config/database';
import type { Perfume } from '../models/types';

export const perfumeController = {
  async getAll(ctx: Context) {
    try {
      const { search, concentration, year, note } = ctx.query;
      const userId = ctx.user!.id;
      
      const result = await db.find({
        selector: { type: 'perfume', userId },
        sort: [{ createdAt: 'desc' }],
      });

      let perfumes = result.docs as Perfume[];

      // Apply filters
      if (search && typeof search === 'string') {
        const searchLower = search.toLowerCase();
        perfumes = perfumes.filter(p => 
          p.name.toLowerCase().includes(searchLower) ||
          p.designer.toLowerCase().includes(searchLower) ||
          p.description?.toLowerCase().includes(searchLower)
        );
      }

      if (concentration && typeof concentration === 'string') {
        perfumes = perfumes.filter(p => p.concentration === concentration);
      }

      if (year && typeof year === 'string') {
        perfumes = perfumes.filter(p => p.year === parseInt(year));
      }

      if (note && typeof note === 'string') {
        const noteLower = note.toLowerCase();
        perfumes = perfumes.filter(p => {
          const allNotes = [
            ...p.pyramid.top,
            ...p.pyramid.middle,
            ...p.pyramid.base
          ].map(n => n.toLowerCase());
          return allNotes.some(n => n.includes(noteLower));
        });
      }

      ctx.body = perfumes;
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async getById(ctx: Context) {
    try {
      const { id } = ctx.params;
      const userId = ctx.user!.id;
      const perfume = await db.get(id) as any;
      
      if (perfume.type !== 'perfume' || perfume.userId !== userId) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Perfume not found' } };
        return;
      }
      
      ctx.body = perfume;
    } catch (error: any) {
      if (error.statusCode === 404) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Perfume not found' } };
      } else {
        ctx.status = 500;
        ctx.body = { error: { message: error.message } };
      }
    }
  },

  async create(ctx: Context) {
    try {
      const userId = ctx.user!.id;
      const perfumeData = ctx.request.body as Omit<Perfume, '_id' | '_rev' | 'type' | 'userId' | 'createdAt' | 'updatedAt'>;
      
      const perfume: Perfume = {
        ...perfumeData,
        userId,
        type: 'perfume',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      const result = await db.insert(perfume);
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
      const userId = ctx.user!.id;
      const updates = ctx.request.body as Partial<Perfume>;
      
      const existing = await db.get(id) as any;
      
      if (existing.type !== 'perfume' || existing.userId !== userId) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Perfume not found' } };
        return;
      }

      const updated: Perfume = {
        ...(existing as Perfume),
        ...updates,
        type: 'perfume',
        userId,
        _id: existing._id,
        _rev: existing._rev,
        createdAt: existing.createdAt,
        updatedAt: new Date().toISOString(),
      };

      await db.insert(updated);
      const result = await db.get(id);
      
      ctx.body = result;
    } catch (error: any) {
      if (error.statusCode === 404) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Perfume not found' } };
      } else {
        ctx.status = 500;
        ctx.body = { error: { message: error.message } };
      }
    }
  },

  async delete(ctx: Context) {
    try {
      const { id } = ctx.params;
      const userId = ctx.user!.id;
      const perfume = await db.get(id) as any;
      
      if (perfume.type !== 'perfume' || perfume.userId !== userId) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Perfume not found' } };
        return;
      }

      await db.destroy(id, perfume._rev);
      
      ctx.status = 204;
    } catch (error: any) {
      if (error.statusCode === 404) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Perfume not found' } };
      } else {
        ctx.status = 500;
        ctx.body = { error: { message: error.message } };
      }
    }
  },
};
