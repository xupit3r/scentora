import Router from '@koa/router';
import { exportController } from '../controllers/exportController';

const router = new Router({
  prefix: '/export',
});

router.get('/collection', exportController.exportCollection);
router.post('/import', exportController.importCollection);

export default router;
