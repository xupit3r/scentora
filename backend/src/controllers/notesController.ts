import { Context } from 'koa';
import { db } from '../config/database';

export const notesController = {
  async getAllNotes(ctx: Context) {
    try {
      const result = await db.find({
        selector: { type: 'perfume' },
      });

      const notesSet = new Set<string>();
      
      result.docs.forEach((doc: any) => {
        if (doc.pyramid) {
          doc.pyramid.top?.forEach((note: string) => notesSet.add(note));
          doc.pyramid.middle?.forEach((note: string) => notesSet.add(note));
          doc.pyramid.base?.forEach((note: string) => notesSet.add(note));
        }
      });

      const notes = Array.from(notesSet).sort();
      
      ctx.body = { notes, count: notes.length };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },
};
