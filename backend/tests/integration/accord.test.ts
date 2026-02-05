import { describe, it, expect, beforeEach } from 'vitest';
import supertest from 'supertest';
import { createApp } from '../../src/app.js';
import { createTestUser, generateAccessToken } from '../helpers/auth.helper.js';

const app = createApp();
const request = supertest(app.callback());

let userId: string;
let token: string;

beforeEach(async () => {
  const user = await createTestUser();
  userId = user.id;
  token = generateAccessToken(user.id, user.email);
});

describe('Accord Routes', () => {
  describe('POST /api/accords', () => {
    it('should create an accord', async () => {
      const res = await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({
          name: 'Lavender',
          pyramidPosition: 'top',
          volumeMl: 50,
          tags: ['fresh', 'floral'],
        });

      expect(res.status).toBe(201);
      expect(res.body.accord._id).toBeDefined();
      expect(res.body.accord.name).toBe('Lavender');
      expect(res.body.accord.pyramidPosition).toBe('top');
      expect(res.body.accord.volumeMl).toBe(50);
      expect(res.body.accord.volumeDrops).toBe(1000);
      expect(res.body.accord.tags).toEqual(['fresh', 'floral']);
    });

    it('should fail without auth', async () => {
      const res = await request.post('/api/accords').send({
        name: 'Test',
        pyramidPosition: 'top',
        volumeMl: 10,
      });
      expect(res.status).toBe(401);
    });

    it('should fail with invalid pyramid position', async () => {
      const res = await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({
          name: 'Test',
          pyramidPosition: 'invalid',
          volumeMl: 10,
        });
      expect(res.status).toBe(400);
    });
  });

  describe('GET /api/accords', () => {
    it('should list user accords', async () => {
      // Create two accords
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Accord 1', pyramidPosition: 'top', volumeMl: 10 });
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Accord 2', pyramidPosition: 'base', volumeMl: 20 });

      const res = await request.get('/api/accords')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.accords).toHaveLength(2);
    });

    it('should filter by position', async () => {
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Top', pyramidPosition: 'top', volumeMl: 10 });
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Base', pyramidPosition: 'base', volumeMl: 20 });

      const res = await request.get('/api/accords?position=top')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.accords).toHaveLength(1);
      expect(res.body.accords[0].pyramidPosition).toBe('top');
    });

    it('should filter by tags', async () => {
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Tagged', pyramidPosition: 'top', volumeMl: 10, tags: ['floral'] });
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Untagged', pyramidPosition: 'base', volumeMl: 20 });

      const res = await request.get('/api/accords?tags=floral')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.accords).toHaveLength(1);
      expect(res.body.accords[0].name).toBe('Tagged');
    });
  });

  describe('GET /api/accords/:id', () => {
    it('should get a single accord', async () => {
      const createRes = await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Rose', pyramidPosition: 'middle', volumeMl: 30 });

      const res = await request.get(`/api/accords/${createRes.body.accord._id}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.accord.name).toBe('Rose');
    });

    it('should return 404 for non-existent accord', async () => {
      const res = await request.get('/api/accords/00000000-0000-0000-0000-000000000000')
        .set('Authorization', `Bearer ${token}`);
      expect(res.status).toBe(404);
    });
  });

  describe('PUT /api/accords/:id', () => {
    it('should update an accord', async () => {
      const createRes = await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Vetiver', pyramidPosition: 'base', volumeMl: 25 });

      const res = await request.put(`/api/accords/${createRes.body.accord._id}`)
        .set('Authorization', `Bearer ${token}`)
        .send({ volumeMl: 35, supplier: 'Supplier A' });

      expect(res.status).toBe(200);
      expect(res.body.accord.volumeMl).toBe(35);
      expect(res.body.accord.supplier).toBe('Supplier A');
    });
  });

  describe('DELETE /api/accords/:id', () => {
    it('should delete an accord', async () => {
      const createRes = await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'ToDelete', pyramidPosition: 'top', volumeMl: 5 });

      const res = await request.delete(`/api/accords/${createRes.body.accord._id}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.message).toContain('deleted');

      // Verify deleted
      const getRes = await request.get(`/api/accords/${createRes.body.accord._id}`)
        .set('Authorization', `Bearer ${token}`);
      expect(getRes.status).toBe(404);
    });
  });

  describe('POST /api/accords/:id/tags', () => {
    it('should add a tag to an accord', async () => {
      const createRes = await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Jasmine', pyramidPosition: 'middle', volumeMl: 15 });

      const res = await request.post(`/api/accords/${createRes.body.accord._id}/tags`)
        .set('Authorization', `Bearer ${token}`)
        .send({ tag: 'floral' });

      expect(res.status).toBe(200);
      expect(res.body.message).toContain('Tag added');
    });
  });

  describe('DELETE /api/accords/:id/tags/:tag', () => {
    it('should remove a tag from an accord', async () => {
      const createRes = await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Bergamot', pyramidPosition: 'top', volumeMl: 10, tags: ['citrus'] });

      const res = await request.delete(`/api/accords/${createRes.body.accord._id}/tags/citrus`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);

      // Verify tag removed
      const getRes = await request.get(`/api/accords/${createRes.body.accord._id}`)
        .set('Authorization', `Bearer ${token}`);
      expect(getRes.body.accord.tags).toEqual([]);
    });
  });
});
