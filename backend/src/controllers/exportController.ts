import { Context } from 'koa';
import { db } from '../config/database';

export const exportController = {
  async exportCollection(ctx: Context) {
    try {
      const userId = ctx.user!.id;
      
      const perfumesResult = await db.find({
        selector: { type: 'perfume', userId },
      });

      const journalResult = await db.find({
        selector: { type: 'journal', userId },
      });

      const exportData = {
        version: '1.0',
        exportDate: new Date().toISOString(),
        perfumes: perfumesResult.docs,
        journalEntries: journalResult.docs,
      };

      ctx.set('Content-Disposition', 'attachment; filename="scentora-collection.json"');
      ctx.set('Content-Type', 'application/json');
      ctx.body = exportData;
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async importCollection(ctx: Context) {
    try {
      const userId = ctx.user!.id;
      const data = ctx.request.body as any;

      if (!data.perfumes || !Array.isArray(data.perfumes)) {
        ctx.status = 400;
        ctx.body = { error: { message: 'Invalid import data format' } };
        return;
      }

      const results = {
        perfumesImported: 0,
        journalEntriesImported: 0,
        errors: [] as string[],
      };

      // Import perfumes
      for (const perfume of data.perfumes) {
        try {
          // Remove _id and _rev to let CouchDB assign new ones
          const { _id, _rev, ...perfumeData } = perfume;
          await db.insert({
            ...perfumeData,
            userId,
            type: 'perfume',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          });
          results.perfumesImported++;
        } catch (err: any) {
          results.errors.push(`Failed to import perfume: ${perfume.name}`);
        }
      }

      // Import journal entries (if present)
      if (data.journalEntries && Array.isArray(data.journalEntries)) {
        for (const entry of data.journalEntries) {
          try {
            const { _id, _rev, ...entryData } = entry;
            await db.insert({
              ...entryData,
              userId,
              type: 'journal',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            });
            results.journalEntriesImported++;
          } catch (err: any) {
            results.errors.push(`Failed to import journal entry`);
          }
        }
      }

      ctx.body = results;
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },
};
