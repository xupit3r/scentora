import jwt from 'jsonwebtoken';
import bcrypt from 'bcryptjs';
import crypto from 'crypto';
import { config } from '../config';
import { db } from '../config/database';
import type { AuthUser, RefreshToken } from '../models/types';

export function generateAccessToken(user: AuthUser): string {
  return jwt.sign(
    { id: user.id, email: user.email, username: user.username },
    config.jwtSecret as string,
    { expiresIn: config.jwtAccessExpiresIn as any }
  );
}

export async function generateRefreshToken(userId: string): Promise<string> {
  const token = crypto.randomBytes(64).toString('hex');
  const expiresAt = new Date();
  expiresAt.setDate(expiresAt.getDate() + 7); // 7 days from now

  const refreshToken: RefreshToken = {
    type: 'refresh_token',
    userId,
    token,
    expiresAt: expiresAt.toISOString(),
    createdAt: new Date().toISOString(),
    revoked: false,
  };

  await db.insert(refreshToken);
  return token;
}

export async function verifyRefreshToken(token: string): Promise<string | null> {
  try {
    const result = await db.find({
      selector: {
        type: 'refresh_token',
        token,
        revoked: false,
      },
    });

    if (result.docs.length === 0) {
      return null;
    }

    const refreshToken = result.docs[0] as any as RefreshToken;

    // Check if expired
    if (new Date(refreshToken.expiresAt) < new Date()) {
      return null;
    }

    return refreshToken.userId;
  } catch (error) {
    return null;
  }
}

export async function revokeRefreshToken(token: string): Promise<boolean> {
  try {
    const result = await db.find({
      selector: {
        type: 'refresh_token',
        token,
      },
    });

    if (result.docs.length === 0) {
      return false;
    }

    const refreshToken = result.docs[0] as any;
    const updated = {
      ...refreshToken,
      revoked: true,
    };
    await db.insert(updated as any);

    return true;
  } catch (error) {
    return false;
  }
}

export async function revokeAllUserRefreshTokens(userId: string): Promise<void> {
  try {
    const result = await db.find({
      selector: {
        type: 'refresh_token',
        userId,
        revoked: false,
      },
    });

    for (const doc of result.docs) {
      const updated = {
        ...(doc as any),
        revoked: true,
      };
      await db.insert(updated as any);
    }
  } catch (error) {
    // Silently fail
  }
}

export function verifyAccessToken(token: string): AuthUser {
  return jwt.verify(token, config.jwtSecret as string) as AuthUser;
}

export async function hashPassword(password: string): Promise<string> {
  return bcrypt.hash(password, 10);
}

export async function comparePassword(password: string, hashedPassword: string): Promise<boolean> {
  return bcrypt.compare(password, hashedPassword);
}
