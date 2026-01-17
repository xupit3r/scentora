import Router from '@koa/router';

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

export default router;
