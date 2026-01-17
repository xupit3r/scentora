import Router from '@koa/router';
import { notesController } from '../controllers/notesController';

const router = new Router({
  prefix: '/notes',
});

router.get('/', notesController.getAllNotes);

export default router;
