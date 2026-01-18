import { Context } from 'koa';
import { db } from '../config/database';

export const statsController = {
  async getCollectionStats(ctx: Context) {
    try {
      const userId = ctx.user!.id;
      
      const perfumesResult = await db.find({
        selector: { type: 'perfume', userId },
      });

      const journalResult = await db.find({
        selector: { type: 'journal', userId },
      });

      const perfumes = perfumesResult.docs as any[];
      const journals = journalResult.docs as any[];

      // Calculate statistics
      const totalPerfumes = perfumes.length;
      const totalJournalEntries = journals.length;

      // Designer stats
      const designerCounts: Record<string, number> = {};
      perfumes.forEach(p => {
        designerCounts[p.designer] = (designerCounts[p.designer] || 0) + 1;
      });
      const topDesigners = Object.entries(designerCounts)
        .sort(([, a], [, b]) => b - a)
        .slice(0, 10)
        .map(([designer, count]) => ({ designer, count }));

      // Note frequency
      const noteCounts: Record<string, number> = {};
      perfumes.forEach(p => {
        [...p.pyramid.top, ...p.pyramid.middle, ...p.pyramid.base].forEach(note => {
          noteCounts[note] = (noteCounts[note] || 0) + 1;
        });
      });
      const topNotes = Object.entries(noteCounts)
        .sort(([, a], [, b]) => b - a)
        .slice(0, 20)
        .map(([note, count]) => ({ note, count }));

      // Concentration distribution
      const concentrationCounts: Record<string, number> = {};
      perfumes.forEach(p => {
        if (p.concentration) {
          concentrationCounts[p.concentration] = (concentrationCounts[p.concentration] || 0) + 1;
        }
      });

      // Rating stats
      const ratings = journals.filter(j => j.rating).map(j => j.rating);
      const avgRating = ratings.length > 0 
        ? ratings.reduce((sum, r) => sum + r, 0) / ratings.length 
        : 0;

      // Year distribution
      const yearCounts: Record<string, number> = {};
      perfumes.forEach(p => {
        if (p.year) {
          yearCounts[p.year] = (yearCounts[p.year] || 0) + 1;
        }
      });
      const yearDistribution = Object.entries(yearCounts)
        .sort(([a], [b]) => parseInt(a) - parseInt(b))
        .map(([year, count]) => ({ year: parseInt(year), count }));

      // Pyramid level stats
      const topNotesCount = perfumes.reduce((sum, p) => sum + p.pyramid.top.length, 0);
      const middleNotesCount = perfumes.reduce((sum, p) => sum + p.pyramid.middle.length, 0);
      const baseNotesCount = perfumes.reduce((sum, p) => sum + p.pyramid.base.length, 0);

      ctx.body = {
        overview: {
          totalPerfumes,
          totalJournalEntries,
          averageRating: avgRating.toFixed(1),
          uniqueNotes: Object.keys(noteCounts).length,
          uniqueDesigners: Object.keys(designerCounts).length,
        },
        topDesigners,
        topNotes,
        concentrationDistribution: concentrationCounts,
        yearDistribution,
        pyramidStats: {
          topNotes: topNotesCount,
          middleNotes: middleNotesCount,
          baseNotes: baseNotesCount,
          avgNotesPerPerfume: totalPerfumes > 0 
            ? ((topNotesCount + middleNotesCount + baseNotesCount) / totalPerfumes).toFixed(1)
            : 0,
        },
      };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },
};
