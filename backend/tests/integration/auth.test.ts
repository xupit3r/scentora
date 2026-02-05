import { describe, it, expect } from 'vitest';
import supertest from 'supertest';
import { createApp } from '../../src/app.js';
import { createTestUser, generateAccessToken, createTestInvitation, uniqueEmail } from '../helpers/auth.helper.js';

const app = createApp();
const request = supertest(app.callback());

describe('Auth Routes', () => {
  describe('POST /api/auth/register', () => {
    it('should register a new user with valid invitation', async () => {
      const creator = await createTestUser();
      const invitation = await createTestInvitation(creator.id);

      const res = await request.post('/api/auth/register').send({
        email: uniqueEmail(),
        username: 'newuser',
        password: 'testpass',
        invitationCode: invitation.code,
      });

      expect(res.status).toBe(201);
      expect(res.body.user).toBeDefined();
      expect(res.body.user._id).toBeDefined();
      expect(res.body.accessToken).toBeDefined();
      expect(res.body.refreshToken).toBeDefined();
    });

    it('should fail with invalid invitation code', async () => {
      const res = await request.post('/api/auth/register').send({
        email: uniqueEmail(),
        username: 'newuser',
        password: 'testpass',
        invitationCode: 'invalid-code',
      });

      expect(res.status).toBe(400);
      expect(res.body.error.message).toBeDefined();
    });

    it('should fail with duplicate email', async () => {
      const creator = await createTestUser();
      const invitation = await createTestInvitation(creator.id);

      const res = await request.post('/api/auth/register').send({
        email: creator.email,
        username: 'anotheruser',
        password: 'testpass',
        invitationCode: invitation.code,
      });

      expect(res.status).toBe(400);
      expect(res.body.error.message).toContain('email already exists');
    });

    it('should fail with missing fields', async () => {
      const res = await request.post('/api/auth/register').send({
        email: 'test@example.com',
      });

      expect(res.status).toBe(400);
    });
  });

  describe('POST /api/auth/login', () => {
    it('should login with valid credentials', async () => {
      const user = await createTestUser({ password: 'mypassword' });

      const res = await request.post('/api/auth/login').send({
        email: user.email,
        password: 'mypassword',
      });

      expect(res.status).toBe(200);
      expect(res.body.user._id).toBe(user.id);
      expect(res.body.accessToken).toBeDefined();
      expect(res.body.refreshToken).toBeDefined();
    });

    it('should fail with wrong password', async () => {
      const user = await createTestUser();

      const res = await request.post('/api/auth/login').send({
        email: user.email,
        password: 'wrongpassword',
      });

      expect(res.status).toBe(401);
      expect(res.body.error.message).toContain('invalid credentials');
    });

    it('should fail with non-existent email', async () => {
      const res = await request.post('/api/auth/login').send({
        email: 'nonexistent@example.com',
        password: 'password',
      });

      expect(res.status).toBe(401);
    });
  });

  describe('POST /api/auth/refresh', () => {
    it('should refresh tokens', async () => {
      const user = await createTestUser({ password: 'mypassword' });
      const loginRes = await request.post('/api/auth/login').send({
        email: user.email,
        password: 'mypassword',
      });

      const res = await request.post('/api/auth/refresh').send({
        refreshToken: loginRes.body.refreshToken,
      });

      expect(res.status).toBe(200);
      expect(res.body.accessToken).toBeDefined();
      expect(res.body.refreshToken).toBeDefined();
    });

    it('should fail with invalid refresh token', async () => {
      const res = await request.post('/api/auth/refresh').send({
        refreshToken: 'invalid-token',
      });

      expect(res.status).toBe(401);
    });
  });

  describe('POST /api/auth/logout', () => {
    it('should logout (revoke refresh token)', async () => {
      const user = await createTestUser({ password: 'mypassword' });
      const loginRes = await request.post('/api/auth/login').send({
        email: user.email,
        password: 'mypassword',
      });

      const res = await request.post('/api/auth/logout').send({
        refreshToken: loginRes.body.refreshToken,
      });

      expect(res.status).toBe(204);

      // Old refresh token should no longer work
      const refreshRes = await request.post('/api/auth/refresh').send({
        refreshToken: loginRes.body.refreshToken,
      });
      expect(refreshRes.status).toBe(401);
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

  describe('POST /api/auth/logout-all', () => {
    it('should revoke all tokens', async () => {
      const user = await createTestUser({ password: 'mypassword' });
      const token = generateAccessToken(user.id, user.email);
      const loginRes = await request.post('/api/auth/login').send({
        email: user.email,
        password: 'mypassword',
      });

      const res = await request.post('/api/auth/logout-all')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.message).toBe('Logged out from all devices');
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
