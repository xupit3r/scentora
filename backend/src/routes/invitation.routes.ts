import Router from '@koa/router';
import { jwtAuth } from '../middleware/auth.js';
import { createInvitationSchema } from '../schemas/invitation.schema.js';
import * as invitationService from '../services/invitation.service.js';
import { ValidationError } from '../utils/errors.js';

const router = new Router({ prefix: '/api/invitations' });
router.use(jwtAuth());

// POST /api/invitations
router.post('/', async (ctx) => {
  const parsed = createInvitationSchema.safeParse(ctx.request.body);
  if (!parsed.success) {
    throw new ValidationError('Validation failed', parsed.error.format());
  }
  const invitation = await invitationService.createInvitation(
    ctx.state.userId,
    parsed.data.email,
    parsed.data.expiresInDays,
  );
  ctx.status = 201;
  ctx.body = { invitation };
});

// GET /api/invitations
router.get('/', async (ctx) => {
  const invitations = await invitationService.listInvitations(ctx.state.userId);
  ctx.body = { invitations };
});

// DELETE /api/invitations/:code
router.delete('/:code', async (ctx) => {
  await invitationService.revokeInvitation(ctx.params.code, ctx.state.userId);
  ctx.body = { message: 'Invitation revoked successfully' };
});

export default router;
