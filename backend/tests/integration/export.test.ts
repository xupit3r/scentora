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

describe('Export/Import Routes', () => {
  describe('GET /api/export', () => {
    it('should export accords as JSON', async () => {
      await request.post('/api/accords')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Export1', pyramidPosition: 'top', volumeMl: 10 });

      const res = await request.get('/api/export')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.version).toBe('1.0');
      expect(res.body.accords).toHaveLength(1);
      expect(res.body.accords[0].name).toBe('Export1');
      expect(res.headers['content-disposition']).toContain('scentora-export.json');
    });
  });

  describe('POST /api/export/import', () => {
    it('should import accords', async () => {
      const res = await request.post('/api/export/import')
        .set('Authorization', `Bearer ${token}`)
        .send({
          version: '1.0',
          accords: [
            { name: 'Imported1', pyramidPosition: 'top', volumeMl: 10, tags: ['fresh'] },
            { name: 'Imported2', pyramidPosition: 'base', volumeMl: 20 },
          ],
        });

      expect(res.status).toBe(200);
      expect(res.body.totalRecords).toBe(2);
      expect(res.body.importedRecords).toBe(2);
      expect(res.body.failedRecords).toBe(0);

      // Verify imported
      const listRes = await request.get('/api/accords')
        .set('Authorization', `Bearer ${token}`);
      expect(listRes.body.accords).toHaveLength(2);
    });
  });
});
