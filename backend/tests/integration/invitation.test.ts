import { describe, it, expect, beforeEach } from 'vitest';
import supertest from 'supertest';
import { createApp } from '../../src/app.js';
import { createTestUser, generateAccessToken } from '../helpers/auth.helper.js';

const app = createApp();
const request = supertest(app.callback());

let token: string;

beforeEach(async () => {
  const user = await createTestUser();
  token = generateAccessToken(user.id, user.email);
});

describe('Invitation Routes', () => {
  describe('POST /api/invitations', () => {
    it('should create an invitation', async () => {
      const res = await request.post('/api/invitations')
        .set('Authorization', `Bearer ${token}`)
        .send({});

      expect(res.status).toBe(201);
      expect(res.body.invitation._id).toBeDefined();
      expect(res.body.invitation.code).toBeDefined();
      expect(res.body.invitation.code.length).toBe(32); // 16 bytes hex
      expect(res.body.invitation.used).toBe(false);
    });

    it('should create an email-specific invitation', async () => {
      const res = await request.post('/api/invitations')
        .set('Authorization', `Bearer ${token}`)
        .send({ email: 'friend@example.com' });

      expect(res.status).toBe(201);
      expect(res.body.invitation.email).toBe('friend@example.com');
    });
  });

  describe('GET /api/invitations', () => {
    it('should list invitations for creator', async () => {
      await request.post('/api/invitations')
        .set('Authorization', `Bearer ${token}`)
        .send({});
      await request.post('/api/invitations')
        .set('Authorization', `Bearer ${token}`)
        .send({});

      const res = await request.get('/api/invitations')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.invitations).toHaveLength(2);
    });
  });

  describe('DELETE /api/invitations/:code', () => {
    it('should revoke an invitation', async () => {
      const createRes = await request.post('/api/invitations')
        .set('Authorization', `Bearer ${token}`)
        .send({});

      const code = createRes.body.invitation.code;
      const res = await request.delete(`/api/invitations/${code}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.message).toContain('revoked');
    });
  });
});
