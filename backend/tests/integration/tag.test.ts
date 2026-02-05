import { describe, it, expect } from 'vitest';
import supertest from 'supertest';
import { createApp } from '../../src/app.js';
import { testPrisma } from '../setup.js';

const app = createApp();
const request = supertest(app.callback());

describe('Tag Routes (public)', () => {
  describe('GET /api/tags', () => {
    it('should return predefined tags', async () => {
      // Seed a test tag
      await testPrisma.predefinedTag.create({
        data: { category: 'test', tag: 'test-tag' },
      });

      const res = await request.get('/api/tags');
      expect(res.status).toBe(200);
      expect(res.body.tags).toBeDefined();
      expect(Array.isArray(res.body.tags)).toBe(true);
      expect(res.body.tags.length).toBeGreaterThanOrEqual(1);
    });
  });

  describe('GET /api/tags/search', () => {
    it('should search tags', async () => {
      await testPrisma.predefinedTag.create({
        data: { category: 'scent', tag: 'citrus-search-test' },
      });

      const res = await request.get('/api/tags/search?q=citrus-search');
      expect(res.status).toBe(200);
      expect(res.body.tags.length).toBeGreaterThanOrEqual(1);
    });

    it('should return empty for no matches', async () => {
      const res = await request.get('/api/tags/search?q=xxxxxxxxx');
      expect(res.status).toBe(200);
      expect(res.body.tags).toEqual([]);
    });
  });

  describe('GET /api/tags/categories', () => {
    it('should return categories', async () => {
      await testPrisma.predefinedTag.create({
        data: { category: 'mood-cat', tag: 'test-cat-tag' },
      });

      const res = await request.get('/api/tags/categories');
      expect(res.status).toBe(200);
      expect(res.body.categories).toBeDefined();
      expect(Array.isArray(res.body.categories)).toBe(true);
    });
  });

  describe('GET /api/tags/grouped', () => {
    it('should return grouped tags', async () => {
      await testPrisma.predefinedTag.create({
        data: { category: 'group-cat', tag: 'grouped-tag' },
      });

      const res = await request.get('/api/tags/grouped');
      expect(res.status).toBe(200);
      expect(res.body.tags).toBeDefined();
      expect(typeof res.body.tags).toBe('object');
    });
  });

  describe('GET /api/tags/category/:category', () => {
    it('should return tags for category', async () => {
      await testPrisma.predefinedTag.create({
        data: { category: 'specific-cat', tag: 'specific-tag' },
      });

      const res = await request.get('/api/tags/category/specific-cat');
      expect(res.status).toBe(200);
      expect(res.body.tags).toHaveLength(1);
      expect(res.body.tags[0].tag).toBe('specific-tag');
    });
  });
});
