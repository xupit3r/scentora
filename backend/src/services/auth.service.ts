import crypto from 'crypto';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import { prisma } from '../prisma.js';
import { config } from '../config.js';
import { parseDuration } from '../utils/duration.js';
import { AppError, UnauthorizedError, ValidationError } from '../utils/errors.js';

function hashToken(token: string): string {
  return crypto.createHash('sha256').update(token).digest('hex');
}

function generateAccessToken(userId: string, email: string): string {
  const expiresInMs = parseDuration(config.jwtAccessExpiresIn);
  return jwt.sign({ userId, email }, config.jwtSecret, {
    expiresIn: Math.floor(expiresInMs / 1000),
  });
}

async function generateRefreshToken(userId: string): Promise<string> {
  const expiresInMs = parseDuration(config.jwtRefreshExpiresIn);
  const tokenString = `${userId}:${Date.now()}`;
  const tokenHash = hashToken(tokenString);

  await prisma.refreshToken.create({
    data: {
      userId,
      tokenHash,
      expiresAt: new Date(Date.now() + expiresInMs),
    },
  });

  return tokenString;
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

export async function register(email: string, username: string, password: string, invitationCode: string) {
  // Validate invitation code
  const invitation = await prisma.invitation.findUnique({ where: { code: invitationCode } });
  if (!invitation) throw new ValidationError('invalid invitation code');
  if (invitation.used) throw new ValidationError('invitation code has already been used');
  if (new Date() > invitation.expiresAt) throw new ValidationError('invitation code has expired');
  if (invitation.email && invitation.email !== email) {
    throw new ValidationError('this invitation is for a different email address');
  }

  // Check uniqueness
  const emailExists = await prisma.user.findUnique({ where: { email } });
  if (emailExists) throw new ValidationError('email already exists');

  const usernameExists = await prisma.user.findFirst({ where: { username } });
  if (usernameExists) throw new ValidationError('username already exists');

  // Create user
  const passwordHash = await bcrypt.hash(password, 10);
  const user = await prisma.user.create({
    data: { email, username, passwordHash },
  });

  // Mark invitation used
  await prisma.invitation.update({
    where: { id: invitation.id },
    data: { used: true, usedAt: new Date(), usedBy: user.id },
  });

  const accessToken = generateAccessToken(user.id, user.email);
  const refreshToken = await generateRefreshToken(user.id);

  return {
    user: userToApi(user),
    accessToken,
    refreshToken,
  };
}

export async function login(email: string, password: string) {
  const user = await prisma.user.findUnique({ where: { email } });
  if (!user) throw new UnauthorizedError('invalid credentials');

  const valid = await bcrypt.compare(password, user.passwordHash);
  if (!valid) throw new UnauthorizedError('invalid credentials');

  const accessToken = generateAccessToken(user.id, user.email);
  const refreshToken = await generateRefreshToken(user.id);

  return {
    user: userToApi(user),
    accessToken,
    refreshToken,
  };
}

export async function refresh(refreshTokenString: string) {
  const tokenHash = hashToken(refreshTokenString);

  const token = await prisma.refreshToken.findUnique({ where: { tokenHash } });
  if (!token || token.revoked) throw new UnauthorizedError('invalid refresh token');
  if (new Date() > token.expiresAt) throw new UnauthorizedError('refresh token expired');

  const user = await prisma.user.findUnique({ where: { id: token.userId } });
  if (!user) throw new UnauthorizedError('user not found');

  const accessToken = generateAccessToken(user.id, user.email);
  const newRefreshToken = await generateRefreshToken(user.id);

  // Revoke old
  await prisma.refreshToken.update({ where: { id: token.id }, data: { revoked: true } });

  return {
    accessToken,
    refreshToken: newRefreshToken,
  };
}

export async function logout(refreshTokenString: string) {
  const tokenHash = hashToken(refreshTokenString);
  await prisma.refreshToken.updateMany({
    where: { tokenHash },
    data: { revoked: true },
  });
}

export async function logoutAll(userId: string) {
  await prisma.refreshToken.updateMany({
    where: { userId },
    data: { revoked: true },
  });
}

export async function getMe(userId: string) {
  const user = await prisma.user.findUnique({ where: { id: userId } });
  if (!user) throw new AppError(404, 'User not found');
  return userToApi(user);
}
