import { prisma } from '../prisma.js';

const SYSTEM_EMAIL = 'system@scentora.local';
const SYSTEM_USERNAME = 'scentora';

export async function ensureSystemUser() {
  const user = await prisma.user.upsert({
    where: { email: SYSTEM_EMAIL },
    update: {},
    create: {
      email: SYSTEM_EMAIL,
      username: SYSTEM_USERNAME,
      passwordHash: 'not-used',
    },
  });
  return user;
}
