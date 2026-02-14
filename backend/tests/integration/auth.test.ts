import { describe, it, expect, beforeEach } from 'vitest';
import supertest from 'supertest';
import { createApp } from '../../src/app.js';
import { createTestUser, generateAccessToken } from '../helpers/auth.helper.js';
import { testPrisma } from '../setup.js';
import { ensureSystemUser } from '../../src/services/systemUser.service.js';

const app = createApp();
const request = supertest(app.callback());

describe('Auth Routes', () => {
  beforeEach(async () => {
    await ensureSystemUser();
  });

  describe('POST /api/auth/login', () => {
    it('should login with correct master password', async () => {
      const res = await request.post('/api/auth/login').send({
        password: 'test-master-password',
      });

      expect(res.status).toBe(200);
      expect(res.body.accessToken).toBeDefined();
    });

    it('should fail with wrong password', async () => {
      const res = await request.post('/api/auth/login').send({
        password: 'wrong-password',
      });

      expect(res.status).toBe(401);
      expect(res.body.error.message).toContain('invalid credentials');
    });

    it('should fail with missing password', async () => {
      const res = await request.post('/api/auth/login').send({});

      expect(res.status).toBe(400);
    });
  });

  describe('GET /api/auth/me', () => {
    it('should return current user', async () => {
      const user = await createTestUser();
      const token = generateAccessToken(user.id, user.email);

      const res = await request.get('/api/auth/me')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.user._id).toBe(user.id);
      expect(res.body.user.email).toBe(user.email);
    });

    it('should fail without auth token', async () => {
      const res = await request.get('/api/auth/me');
      expect(res.status).toBe(401);
    });
  });
});

describe('Health Check', () => {
  it('GET /api/health should return ok', async () => {
    const res = await request.get('/api/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
    expect(res.body.service).toBe('scentora-api');
  });
});
