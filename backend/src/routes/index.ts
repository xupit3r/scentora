import Router from '@koa/router';
import perfumeRoutes from './perfumes';
import journalRoutes from './journal';
import notesRoutes from './notes';
import statsRoutes from './stats';
import exportRoutes from './export';

const router = new Router({
  prefix: '/api',
});

router.get('/health', (ctx) => {
  ctx.body = { 
    status: 'ok',
    timestamp: new Date().toISOString(),
    service: 'scentora-api'
  };
});

// Mount sub-routes
router.use(perfumeRoutes.routes());
router.use(perfumeRoutes.allowedMethods());
router.use(journalRoutes.routes());
router.use(journalRoutes.allowedMethods());
router.use(notesRoutes.routes());
router.use(notesRoutes.allowedMethods());
router.use(statsRoutes.routes());
router.use(statsRoutes.allowedMethods());
router.use(exportRoutes.routes());
router.use(exportRoutes.allowedMethods());

export default router;
