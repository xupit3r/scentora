import Router from '@koa/router';
import { journalController } from '../controllers/journalController';
import { validateBody } from '../middleware/validation';
import { journalEntrySchema, updateJournalEntrySchema } from '../models/schemas';

const router = new Router();

// Journal entries for a specific perfume
router.get('/perfumes/:perfumeId/journal', journalController.getByPerfumeId);
router.post('/perfumes/:perfumeId/journal', validateBody(journalEntrySchema), journalController.create);

// Individual journal entry operations
router.put('/journal/:id', validateBody(updateJournalEntrySchema), journalController.update);
router.delete('/journal/:id', journalController.delete);

export default router;
