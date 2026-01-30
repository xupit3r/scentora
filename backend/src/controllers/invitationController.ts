import { Context } from 'koa';
import { db } from '../config/database';
import { randomBytes } from 'crypto';
import type { Invitation } from '../models/types';

export const invitationController = {
  async create(ctx: Context) {
    try {
      const { email, expiresInDays = 7 } = ctx.request.body as { email?: string; expiresInDays?: number };
      const userId = ctx.user!.id;

      // Generate a unique invitation code
      const code = randomBytes(16).toString('hex');

      // Calculate expiration date
      const expiresAt = new Date();
      expiresAt.setDate(expiresAt.getDate() + expiresInDays);

      const invitation: Invitation = {
        type: 'invitation',
        code,
        email,
        createdBy: userId,
        expiresAt: expiresAt.toISOString(),
        used: false,
        createdAt: new Date().toISOString(),
      };

      const result = await db.insert(invitation);
      const createdInvitation = await db.get(result.id);

      ctx.status = 201;
      ctx.body = { invitation: createdInvitation };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async list(ctx: Context) {
    try {
      const userId = ctx.user!.id;

      // Get all invitations created by this user
      const result = await db.find({
        selector: {
          type: 'invitation',
          createdBy: userId,
        },
        sort: [{ createdAt: 'desc' }],
      });

      ctx.body = { invitations: result.docs };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async revoke(ctx: Context) {
    try {
      const { code } = ctx.params;
      const userId = ctx.user!.id;

      // Find the invitation
      const result = await db.find({
        selector: {
          type: 'invitation',
          code,
        },
      });

      if (result.docs.length === 0) {
        ctx.status = 404;
        ctx.body = { error: { message: 'Invitation not found' } };
        return;
      }

      const invitation = result.docs[0] as any;

      // Check if the user owns this invitation
      if (invitation.createdBy !== userId) {
        ctx.status = 403;
        ctx.body = { error: { message: 'Not authorized to revoke this invitation' } };
        return;
      }

      // Mark as used (effectively revoking it)
      invitation.used = true;
      invitation.usedAt = new Date().toISOString();
      await db.insert(invitation);

      ctx.body = { message: 'Invitation revoked successfully' };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },
};
