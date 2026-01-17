import Router from '@koa/router';
import { statsController } from '../controllers/statsController';

const router = new Router({
  prefix: '/stats',
});

router.get('/', statsController.getCollectionStats);

export default router;
