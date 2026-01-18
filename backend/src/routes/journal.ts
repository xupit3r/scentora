import Router from '@koa/router';
import { journalController } from '../controllers/journalController';
import { validateBody } from '../middleware/validation';
import { authenticate } from '../middleware/auth';
import { journalEntrySchema, updateJournalEntrySchema } from '../models/schemas';

const router = new Router();

// All journal routes require authentication
router.use(authenticate);

// Journal entries for a specific perfume
router.get('/perfumes/:perfumeId/journal', journalController.getByPerfumeId);
router.post('/perfumes/:perfumeId/journal', validateBody(journalEntrySchema), journalController.create);

// Individual journal entry operations
router.put('/journal/:id', validateBody(updateJournalEntrySchema), journalController.update);
router.delete('/journal/:id', journalController.delete);

export default router;
