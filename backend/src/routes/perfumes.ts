import Router from '@koa/router';
import { perfumeController } from '../controllers/perfumeController';
import { validateBody } from '../middleware/validation';
import { createPerfumeSchema, updatePerfumeSchema } from '../models/schemas';

const router = new Router({
  prefix: '/perfumes',
});

router.get('/', perfumeController.getAll);
router.get('/:id', perfumeController.getById);
router.post('/', validateBody(createPerfumeSchema), perfumeController.create);
router.put('/:id', validateBody(updatePerfumeSchema), perfumeController.update);
router.delete('/:id', perfumeController.delete);

export default router;
