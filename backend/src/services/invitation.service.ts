import crypto from 'crypto';
import { prisma } from '../prisma.js';

function invitationToApi(inv: any) {
  return {
    _id: inv.id,
    code: inv.code,
    ...(inv.email != null ? { email: inv.email } : {}),
    createdBy: inv.createdBy,
    expiresAt: inv.expiresAt,
    used: inv.used,
    ...(inv.usedAt != null ? { usedAt: inv.usedAt } : {}),
    ...(inv.usedBy != null ? { usedBy: inv.usedBy } : {}),
    createdAt: inv.createdAt,
  };
}

export async function createInvitation(creatorId: string, email?: string, expiresInDays = 7) {
  const code = crypto.randomBytes(16).toString('hex');
  const expiresAt = new Date(Date.now() + expiresInDays * 24 * 60 * 60 * 1000);

  const invitation = await prisma.invitation.create({
    data: {
      code,
      email: email ?? null,
      createdBy: creatorId,
      expiresAt,
    },
  });

  return invitationToApi(invitation);
}

export async function listInvitations(creatorId: string) {
  const invitations = await prisma.invitation.findMany({
    where: { createdBy: creatorId },
    orderBy: { createdAt: 'desc' },
  });
  return invitations.map(invitationToApi);
}

export async function revokeInvitation(code: string, creatorId: string) {
  const invitation = await prisma.invitation.findFirst({
    where: { code, createdBy: creatorId },
  });
  if (!invitation) throw new Error('invitation not found');

  await prisma.invitation.delete({ where: { id: invitation.id } });
}
