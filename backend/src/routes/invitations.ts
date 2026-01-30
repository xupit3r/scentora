import Router from '@koa/router';
import { invitationController } from '../controllers/invitationController';
import { validateBody } from '../middleware/validation';
import { createInvitationSchema } from '../models/schemas';
import { authenticate } from '../middleware/auth';

const router = new Router({
  prefix: '/invitations',
});

// All invitation routes require authentication
router.post('/', authenticate, validateBody(createInvitationSchema), invitationController.create);
router.get('/', authenticate, invitationController.list);
router.delete('/:code', authenticate, invitationController.revoke);

export default router;
