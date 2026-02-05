import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import { testPrisma } from '../setup.js';
import { config } from '../../src/config.js';

let counter = 0;

export function uniqueEmail() {
  counter++;
  return `test${counter}-${Date.now()}@example.com`;
}

export async function createTestUser(overrides: { email?: string; username?: string; password?: string } = {}) {
  const email = overrides.email || uniqueEmail();
  const username = overrides.username || `user${counter}`;
  const password = overrides.password || 'testpass';
  const passwordHash = await bcrypt.hash(password, 10);

  const user = await testPrisma.user.create({
    data: { email, username, passwordHash },
  });

  return { ...user, password };
}

export function generateAccessToken(userId: string, email: string) {
  return jwt.sign({ userId, email }, config.jwtSecret, { expiresIn: '15m' });
}

export async function createTestInvitation(creatorId: string) {
  const code = `test-${Date.now()}-${Math.random().toString(36).slice(2)}`;
  return testPrisma.invitation.create({
    data: {
      code,
      createdBy: creatorId,
      expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
    },
  });
}
