import jwt from 'jsonwebtoken';
import { prisma } from '../prisma.js';
import { config } from '../config.js';
import { parseDuration } from '../utils/duration.js';
import { AppError, UnauthorizedError } from '../utils/errors.js';

const SYSTEM_EMAIL = 'system@scentora.local';

function generateAccessToken(userId: string, email: string): string {
  const expiresInMs = parseDuration(config.jwtAccessExpiresIn);
  return jwt.sign({ userId, email }, config.jwtSecret, {
    expiresIn: Math.floor(expiresInMs / 1000),
  });
}

function userToApi(user: { id: string; email: string; username: string; validateRecipeVolumes: boolean; createdAt: Date; updatedAt: Date }) {
  return {
    _id: user.id,
    email: user.email,
    username: user.username,
    validateRecipeVolumes: user.validateRecipeVolumes,
    createdAt: user.createdAt,
    updatedAt: user.updatedAt,
  };
}

export async function login(password: string) {
  if (password !== config.masterPassword) {
    throw new UnauthorizedError('invalid credentials');
  }

  const user = await prisma.user.findUnique({ where: { email: SYSTEM_EMAIL } });
  if (!user) throw new AppError(500, 'System user not found');

  const accessToken = generateAccessToken(user.id, user.email);

  return { accessToken };
}

export async function getMe(userId: string) {
  const user = await prisma.user.findUnique({ where: { id: userId } });
  if (!user) throw new AppError(404, 'User not found');
  return userToApi(user);
}
