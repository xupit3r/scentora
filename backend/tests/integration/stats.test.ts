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

describe('Stats Routes', () => {
  describe('GET /api/stats', () => {
    it('should return empty stats for new user', async () => {
      const res = await request.get('/api/stats')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.overview.totalAccords).toBe(0);
      expect(res.body.pyramidStats).toBeDefined();
      expect(res.body.tagStats).toEqual([]);
      expect(res.body.supplierStats).toEqual([]);
      expect(res.body.volumeStats).toBeDefined();
      expect(res.body.lowInventory).toEqual([]);
    });

    it('should return correct stats after adding accords', async () => {
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Top1', pyramidPosition: 'top', volumeMl: 30, tags: ['fresh'] });
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Base1', pyramidPosition: 'base', volumeMl: 5, supplier: 'Supplier A' });

      const res = await request.get('/api/stats')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.overview.totalAccords).toBe(2);
      expect(res.body.overview.totalVolume).toBe(35);
      expect(res.body.pyramidStats.topCount).toBe(1);
      expect(res.body.pyramidStats.baseCount).toBe(1);
      expect(res.body.lowInventory).toHaveLength(1); // Base1 < 10ml
    });

    it('should fail without auth', async () => {
      const res = await request.get('/api/stats');
      expect(res.status).toBe(401);
    });
  });
});
