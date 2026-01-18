import Router from '@koa/router';
import { exportController } from '../controllers/exportController';
import { authenticate } from '../middleware/auth';

const router = new Router({
  prefix: '/export',
});

router.use(authenticate);

router.get('/collection', exportController.exportCollection);
router.post('/import', exportController.importCollection);

export default router;
