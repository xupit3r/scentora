import { Context } from 'koa';
import { db } from '../config/database';
import { 
  hashPassword, 
  comparePassword, 
  generateAccessToken,
  generateRefreshToken,
  verifyRefreshToken,
  revokeRefreshToken,
  revokeAllUserRefreshTokens
} from '../config/auth';
import type { User, AuthUser } from '../models/types';

export const authController = {
  async register(ctx: Context) {
    try {
      const { email, username, password } = ctx.request.body as { email: string; username: string; password: string };

      // Check if user already exists
      const existingUsers = await db.find({
        selector: {
          type: 'user',
          $or: [
            { email },
            { username }
          ]
        }
      });

      if (existingUsers.docs.length > 0) {
        ctx.status = 400;
        ctx.body = { error: { message: 'Email or username already exists' } };
        return;
      }

      // Hash password
      const hashedPassword = await hashPassword(password);

      // Create user
      const user: User = {
        type: 'user',
        email,
        username,
        password: hashedPassword,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      const result = await db.insert(user);
      const createdUser = await db.get(result.id) as any;

      // Generate tokens
      const authUser: AuthUser = {
        id: createdUser._id,
        email: createdUser.email,
        username: createdUser.username,
      };

      const accessToken = generateAccessToken(authUser);
      const refreshToken = await generateRefreshToken(createdUser._id);

      ctx.status = 201;
      ctx.body = {
        user: authUser,
        accessToken,
        refreshToken,
      };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async login(ctx: Context) {
    try {
      const { email, password } = ctx.request.body as { email: string; password: string };

      // Find user by email
      const result = await db.find({
        selector: {
          type: 'user',
          email,
        },
      });

      if (result.docs.length === 0) {
        ctx.status = 401;
        ctx.body = { error: { message: 'Invalid email or password' } };
        return;
      }

      const user = result.docs[0] as any;

      // Verify password
      const isValid = await comparePassword(password, user.password);

      if (!isValid) {
        ctx.status = 401;
        ctx.body = { error: { message: 'Invalid email or password' } };
        return;
      }

      // Generate tokens
      const authUser: AuthUser = {
        id: user._id,
        email: user.email,
        username: user.username,
      };

      const accessToken = generateAccessToken(authUser);
      const refreshToken = await generateRefreshToken(user._id);

      ctx.body = {
        user: authUser,
        accessToken,
        refreshToken,
      };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async refresh(ctx: Context) {
    try {
      const { refreshToken } = ctx.request.body as { refreshToken: string };

      if (!refreshToken) {
        ctx.status = 400;
        ctx.body = { error: { message: 'Refresh token required' } };
        return;
      }

      // Verify refresh token
      const userId = await verifyRefreshToken(refreshToken);

      if (!userId) {
        ctx.status = 401;
        ctx.body = { error: { message: 'Invalid or expired refresh token' } };
        return;
      }

      // Get user
      const user = await db.get(userId) as any;

      if (!user || user.type !== 'user') {
        ctx.status = 401;
        ctx.body = { error: { message: 'User not found' } };
        return;
      }

      // Revoke old refresh token (rotation)
      await revokeRefreshToken(refreshToken);

      // Generate new tokens
      const authUser: AuthUser = {
        id: user._id,
        email: user.email,
        username: user.username,
      };

      const accessToken = generateAccessToken(authUser);
      const newRefreshToken = await generateRefreshToken(user._id);

      ctx.body = {
        accessToken,
        refreshToken: newRefreshToken,
      };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async logout(ctx: Context) {
    try {
      const { refreshToken } = ctx.request.body as { refreshToken?: string };

      if (refreshToken) {
        await revokeRefreshToken(refreshToken);
      }

      ctx.body = { message: 'Logged out successfully' };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async logoutAll(ctx: Context) {
    try {
      const userId = ctx.user!.id;
      await revokeAllUserRefreshTokens(userId);

      ctx.body = { message: 'Logged out from all devices' };
    } catch (error: any) {
      ctx.status = 500;
      ctx.body = { error: { message: error.message } };
    }
  },

  async me(ctx: Context) {
    // User is already attached by authenticate middleware
    ctx.body = { user: ctx.user };
  },
};
