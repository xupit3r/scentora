import Router from '@koa/router';
import { statsController } from '../controllers/statsController';
import { authenticate } from '../middleware/auth';

const router = new Router({
  prefix: '/stats',
});

router.use(authenticate);

router.get('/', statsController.getCollectionStats);

export default router;
