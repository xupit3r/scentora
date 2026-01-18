import Router from '@koa/router';
import { perfumeController } from '../controllers/perfumeController';
import { validateBody } from '../middleware/validation';
import { authenticate } from '../middleware/auth';
import { createPerfumeSchema, updatePerfumeSchema } from '../models/schemas';

const router = new Router({
  prefix: '/perfumes',
});

// All perfume routes require authentication
router.use(authenticate);

router.get('/', perfumeController.getAll);
router.get('/:id', perfumeController.getById);
router.post('/', validateBody(createPerfumeSchema), perfumeController.create);
router.put('/:id', validateBody(updatePerfumeSchema), perfumeController.update);
router.delete('/:id', perfumeController.delete);

export default router;
