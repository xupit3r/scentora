import Router from '@koa/router';
import { notesController } from '../controllers/notesController';
import { authenticate } from '../middleware/auth';

const router = new Router({
  prefix: '/notes',
});

router.use(authenticate);

router.get('/', notesController.getAllNotes);

export default router;
