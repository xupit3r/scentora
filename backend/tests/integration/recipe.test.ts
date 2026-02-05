import { describe, it, expect, beforeEach } from 'vitest';
import supertest from 'supertest';
import { createApp } from '../../src/app.js';
import { createTestUser, generateAccessToken } from '../helpers/auth.helper.js';

const app = createApp();
const request = supertest(app.callback());

let token: string;
let accordId: string;

beforeEach(async () => {
  const user = await createTestUser();
  token = generateAccessToken(user.id, user.email);

  // Create an accord for recipe ingredients
  const accordRes = await request.post('/api/accords')
    .set('Authorization', `Bearer ${token}`)
    .send({ name: 'Lavender', pyramidPosition: 'top', volumeMl: 50 });
  accordId = accordRes.body.accord._id;
});

describe('Recipe Routes', () => {
  describe('POST /api/recipes', () => {
    it('should create a recipe', async () => {
      const res = await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'My Perfume', targetVolumeMl: 100 });

      expect(res.status).toBe(201);
      expect(res.body._id).toBeDefined();
      expect(res.body.name).toBe('My Perfume');
      expect(res.body.status).toBe('draft');
      expect(res.body.targetVolumeMl).toBe(100);
    });

    it('should fail with duplicate name', async () => {
      await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Duplicate', targetVolumeMl: 100 });

      const res = await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Duplicate', targetVolumeMl: 50 });

      expect(res.status).toBe(400);
    });
  });

  describe('GET /api/recipes', () => {
    it('should list recipes', async () => {
      await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Recipe 1', targetVolumeMl: 100 });
      await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Recipe 2', targetVolumeMl: 50 });

      const res = await request.get('/api/recipes')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body).toHaveLength(2);
    });

    it('should filter by status', async () => {
      await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Draft', targetVolumeMl: 100 });
      await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'In Progress', targetVolumeMl: 50, status: 'in_progress' });

      const res = await request.get('/api/recipes?status=draft')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body).toHaveLength(1);
      expect(res.body[0].name).toBe('Draft');
    });
  });

  describe('GET /api/recipes/search', () => {
    it('should search recipes', async () => {
      await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Summer Breeze', targetVolumeMl: 100 });
      await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Winter Night', targetVolumeMl: 50 });

      const res = await request.get('/api/recipes/search?q=summer')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body).toHaveLength(1);
      expect(res.body[0].name).toBe('Summer Breeze');
    });
  });

  describe('GET /api/recipes/:id', () => {
    it('should get recipe with active version', async () => {
      const createRes = await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Detailed', targetVolumeMl: 100 });

      // Add a version
      await request.post(`/api/recipes/${createRes.body._id}/versions`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'v1' });

      const res = await request.get(`/api/recipes/${createRes.body._id}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body._id).toBe(createRes.body._id);
      expect(res.body.activeVersion).toBeDefined();
      expect(res.body.activeVersion.name).toBe('v1');
    });
  });

  describe('PUT /api/recipes/:id', () => {
    it('should update a recipe', async () => {
      const createRes = await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Original', targetVolumeMl: 100 });

      const res = await request.put(`/api/recipes/${createRes.body._id}`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Updated', status: 'in_progress' });

      expect(res.status).toBe(200);
      expect(res.body.name).toBe('Updated');
      expect(res.body.status).toBe('in_progress');
    });
  });

  describe('DELETE /api/recipes/:id', () => {
    it('should delete a recipe', async () => {
      const createRes = await request.post('/api/recipes')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'ToDelete', targetVolumeMl: 100 });

      const res = await request.delete(`/api/recipes/${createRes.body._id}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body.message).toContain('deleted');
    });
  });
});

describe('Recipe Version Routes', () => {
  let recipeId: string;

  beforeEach(async () => {
    const res = await request.post('/api/recipes')
      .set('Authorization', `Bearer ${token}`)
      .send({ name: 'Version Test', targetVolumeMl: 100 });
    recipeId = res.body._id;
  });

  describe('POST /api/recipes/:id/versions', () => {
    it('should create a version and auto-activate first', async () => {
      const res = await request.post(`/api/recipes/${recipeId}/versions`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'v1' });

      expect(res.status).toBe(201);
      expect(res.body._id).toBeDefined();
      expect(res.body.versionNumber).toBe(1);
      expect(res.body.isActive).toBe(true);
    });
  });

  describe('GET /api/recipes/:id/versions', () => {
    it('should list versions', async () => {
      await request.post(`/api/recipes/${recipeId}/versions`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'v1' });
      await request.post(`/api/recipes/${recipeId}/versions`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'v2' });

      const res = await request.get(`/api/recipes/${recipeId}/versions`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body).toHaveLength(2);
    });
  });

  describe('POST /api/recipes/:id/versions/:versionNumber/activate', () => {
    it('should activate a version', async () => {
      await request.post(`/api/recipes/${recipeId}/versions`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'v1' });
      await request.post(`/api/recipes/${recipeId}/versions`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'v2' });

      const res = await request.post(`/api/recipes/${recipeId}/versions/2/activate`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body._id).toBeDefined();
      expect(res.body.isActive).toBe(true);
      expect(res.body.versionNumber).toBe(2);
    });
  });

  describe('POST /api/recipes/:id/versions/:versionId/duplicate', () => {
    it('should duplicate a version', async () => {
      const v1Res = await request.post(`/api/recipes/${recipeId}/versions`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'v1' });

      const res = await request.post(`/api/recipes/${recipeId}/versions/${v1Res.body._id}/duplicate`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(201);
      expect(res.body.versionNumber).toBe(2);
    });
  });
});

describe('Recipe Ingredient Routes', () => {
  let recipeId: string;
  let versionId: string;

  beforeEach(async () => {
    const recipeRes = await request.post('/api/recipes')
      .set('Authorization', `Bearer ${token}`)
      .send({ name: 'Ingredient Test', targetVolumeMl: 100 });
    recipeId = recipeRes.body._id;

    const versionRes = await request.post(`/api/recipes/${recipeId}/versions`)
      .set('Authorization', `Bearer ${token}`)
      .send({ name: 'v1' });
    versionId = versionRes.body._id;
  });

  describe('POST /api/recipes/:id/versions/:versionId/ingredients', () => {
    it('should add an ingredient', async () => {
      const res = await request.post(`/api/recipes/${recipeId}/versions/${versionId}/ingredients`)
        .set('Authorization', `Bearer ${token}`)
        .send({ accordId, quantityMl: 10 });

      expect(res.status).toBe(201);
      expect(res.body._id).toBeDefined();
      expect(res.body.accordId).toBe(accordId);
      expect(res.body.quantityMl).toBe(10);
      expect(res.body.quantityDrops).toBe(200);
    });

    it('should fail with duplicate accord', async () => {
      await request.post(`/api/recipes/${recipeId}/versions/${versionId}/ingredients`)
        .set('Authorization', `Bearer ${token}`)
        .send({ accordId, quantityMl: 10 });

      const res = await request.post(`/api/recipes/${recipeId}/versions/${versionId}/ingredients`)
        .set('Authorization', `Bearer ${token}`)
        .send({ accordId, quantityMl: 5 });

      expect(res.status).toBe(400);
    });
  });

  describe('PUT /api/recipes/:id/versions/:versionId/ingredients/:ingredientId', () => {
    it('should update an ingredient', async () => {
      const addRes = await request.post(`/api/recipes/${recipeId}/versions/${versionId}/ingredients`)
        .set('Authorization', `Bearer ${token}`)
        .send({ accordId, quantityMl: 10 });

      const res = await request.put(`/api/recipes/${recipeId}/versions/${versionId}/ingredients/${addRes.body._id}`)
        .set('Authorization', `Bearer ${token}`)
        .send({ quantityMl: 20 });

      expect(res.status).toBe(200);
      expect(res.body.quantityMl).toBe(20);
    });
  });

  describe('DELETE /api/recipes/:id/versions/:versionId/ingredients/:ingredientId', () => {
    it('should remove an ingredient', async () => {
      const addRes = await request.post(`/api/recipes/${recipeId}/versions/${versionId}/ingredients`)
        .set('Authorization', `Bearer ${token}`)
        .send({ accordId, quantityMl: 10 });

      const res = await request.delete(`/api/recipes/${recipeId}/versions/${versionId}/ingredients/${addRes.body._id}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
    });
  });
});

describe('Recipe Note Routes', () => {
  let recipeId: string;

  beforeEach(async () => {
    const res = await request.post('/api/recipes')
      .set('Authorization', `Bearer ${token}`)
      .send({ name: 'Note Test', targetVolumeMl: 100 });
    recipeId = res.body._id;
  });

  describe('POST /api/recipes/:id/notes', () => {
    it('should create a note', async () => {
      const res = await request.post(`/api/recipes/${recipeId}/notes`)
        .set('Authorization', `Bearer ${token}`)
        .send({ content: 'Test note', noteType: 'general' });

      expect(res.status).toBe(201);
      expect(res.body._id).toBeDefined();
      expect(res.body.content).toBe('Test note');
      expect(res.body.noteType).toBe('general');
    });
  });

  describe('GET /api/recipes/:id/notes', () => {
    it('should list notes', async () => {
      await request.post(`/api/recipes/${recipeId}/notes`)
        .set('Authorization', `Bearer ${token}`)
        .send({ content: 'Note 1' });
      await request.post(`/api/recipes/${recipeId}/notes`)
        .set('Authorization', `Bearer ${token}`)
        .send({ content: 'Note 2' });

      const res = await request.get(`/api/recipes/${recipeId}/notes`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body).toHaveLength(2);
    });
  });

  describe('PUT /api/recipes/:id/notes/:noteId', () => {
    it('should update a note', async () => {
      const createRes = await request.post(`/api/recipes/${recipeId}/notes`)
        .set('Authorization', `Bearer ${token}`)
        .send({ content: 'Original' });

      const res = await request.put(`/api/recipes/${recipeId}/notes/${createRes.body._id}`)
        .set('Authorization', `Bearer ${token}`)
        .send({ content: 'Updated' });

      expect(res.status).toBe(200);
      expect(res.body.content).toBe('Updated');
    });
  });

  describe('DELETE /api/recipes/:id/notes/:noteId', () => {
    it('should delete a note', async () => {
      const createRes = await request.post(`/api/recipes/${recipeId}/notes`)
        .set('Authorization', `Bearer ${token}`)
        .send({ content: 'To delete' });

      const res = await request.delete(`/api/recipes/${recipeId}/notes/${createRes.body._id}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
    });
  });
});

describe('Recipe Tag Routes', () => {
  let recipeId: string;

  beforeEach(async () => {
    const res = await request.post('/api/recipes')
      .set('Authorization', `Bearer ${token}`)
      .send({ name: 'Tag Test', targetVolumeMl: 100 });
    recipeId = res.body._id;
  });

  describe('POST /api/recipes/:id/tags', () => {
    it('should add a tag', async () => {
      const res = await request.post(`/api/recipes/${recipeId}/tags`)
        .set('Authorization', `Bearer ${token}`)
        .send({ tag: 'floral' });

      expect(res.status).toBe(201);
      expect(res.body.tag).toBe('floral');
    });
  });

  describe('DELETE /api/recipes/:id/tags/:tag', () => {
    it('should remove a tag', async () => {
      await request.post(`/api/recipes/${recipeId}/tags`)
        .set('Authorization', `Bearer ${token}`)
        .send({ tag: 'woody' });

      const res = await request.delete(`/api/recipes/${recipeId}/tags/woody`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
    });
  });

  describe('GET /api/recipes/tags/popular', () => {
    it('should get popular tags', async () => {
      await request.post(`/api/recipes/${recipeId}/tags`)
        .set('Authorization', `Bearer ${token}`)
        .send({ tag: 'popular-tag' });

      const res = await request.get('/api/recipes/tags/popular')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
    });
  });
});

describe('Collection Routes', () => {
  let recipeId: string;

  beforeEach(async () => {
    const res = await request.post('/api/recipes')
      .set('Authorization', `Bearer ${token}`)
      .send({ name: 'Collection Recipe', targetVolumeMl: 100 });
    recipeId = res.body._id;
  });

  describe('POST /api/collections', () => {
    it('should create a collection', async () => {
      const res = await request.post('/api/collections')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'My Collection', description: 'A test collection' });

      expect(res.status).toBe(201);
      expect(res.body._id).toBeDefined();
      expect(res.body.name).toBe('My Collection');
    });
  });

  describe('GET /api/collections', () => {
    it('should list collections', async () => {
      await request.post('/api/collections')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Collection 1' });
      await request.post('/api/collections')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Collection 2' });

      const res = await request.get('/api/collections')
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
      expect(res.body).toHaveLength(2);
    });
  });

  describe('PUT /api/collections/:id', () => {
    it('should update a collection', async () => {
      const createRes = await request.post('/api/collections')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Original' });

      const res = await request.put(`/api/collections/${createRes.body._id}`)
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Renamed' });

      expect(res.status).toBe(200);
      expect(res.body.name).toBe('Renamed');
    });
  });

  describe('DELETE /api/collections/:id', () => {
    it('should delete a collection', async () => {
      const createRes = await request.post('/api/collections')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'To Delete' });

      const res = await request.delete(`/api/collections/${createRes.body._id}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
    });
  });

  describe('POST /api/collections/:id/recipes', () => {
    it('should add a recipe to a collection', async () => {
      const collRes = await request.post('/api/collections')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'My Coll' });

      const res = await request.post(`/api/collections/${collRes.body._id}/recipes`)
        .set('Authorization', `Bearer ${token}`)
        .send({ recipeId });

      expect(res.status).toBe(200);
      expect(res.body.message).toContain('added');
    });
  });

  describe('DELETE /api/collections/:id/recipes/:recipeId', () => {
    it('should remove a recipe from a collection', async () => {
      const collRes = await request.post('/api/collections')
        .set('Authorization', `Bearer ${token}`)
        .send({ name: 'Remove Coll' });

      await request.post(`/api/collections/${collRes.body._id}/recipes`)
        .set('Authorization', `Bearer ${token}`)
        .send({ recipeId });

      const res = await request.delete(`/api/collections/${collRes.body._id}/recipes/${recipeId}`)
        .set('Authorization', `Bearer ${token}`);

      expect(res.status).toBe(200);
    });
  });
});
