import api from './api';
import type {
  Recipe,
  CreateRecipeRequest,
  UpdateRecipeRequest,
  RecipeFilters,
  RecipeVersion,
  CreateRecipeVersionRequest,
  RecipeIngredient,
  CreateRecipeIngredientRequest,
  UpdateRecipeIngredientRequest,
  RecipeNote,
  CreateRecipeNoteRequest,
  UpdateRecipeNoteRequest,
  RecipeTag,
  AddRecipeTagRequest,
  RecipeCollection,
  CreateRecipeCollectionRequest,
  UpdateRecipeCollectionRequest,
  RecipeCollectionFilters,
  AddRecipeToCollectionRequest,
} from '@/types';

/**
 * Recipe API Service
 * Handles all recipe-related API calls
 */
export const recipeService = {
  // ============================================================================
  // Recipe CRUD
  // ============================================================================

  /**
   * Get all recipes for the current user
   * @param filters Optional filters (status, search, pagination)
   * @returns Array of recipes
   */
  async getAll(filters?: RecipeFilters): Promise<Recipe[]> {
    const params: Record<string, any> = {};
    if (filters?.status) params.status = filters.status;
    if (filters?.search) params.search = filters.search;
    if (filters?.limit !== undefined) params.limit = filters.limit;
    if (filters?.offset !== undefined) params.offset = filters.offset;

    const { data } = await api.get('/recipes', { params });
    return Array.isArray(data) ? data : [];
  },

  /**
   * Get a specific recipe by ID
   * @param id Recipe ID
   * @returns Recipe details
   */
  async getById(id: string): Promise<Recipe> {
    const { data } = await api.get(`/recipes/${id}`);
    return data;
  },

  /**
   * Create a new recipe
   * @param recipe Recipe data
   * @returns Created recipe
   */
  async create(recipe: CreateRecipeRequest): Promise<Recipe> {
    const { data } = await api.post('/recipes', recipe);
    return data;
  },

  /**
   * Update an existing recipe
   * @param id Recipe ID
   * @param recipe Updated recipe data
   * @returns Updated recipe
   */
  async update(id: string, recipe: UpdateRecipeRequest): Promise<Recipe> {
    const { data } = await api.put(`/recipes/${id}`, recipe);
    return data;
  },

  /**
   * Delete a recipe
   * @param id Recipe ID
   */
  async delete(id: string): Promise<void> {
    await api.delete(`/recipes/${id}`);
  },

  /**
   * Search recipes by query string
   * @param query Search query
   * @returns Array of matching recipes
   */
  async search(query: string): Promise<Recipe[]> {
    const { data } = await api.get('/recipes/search', {
      params: { q: query },
    });
    return Array.isArray(data) ? data : [];
  },

  // ============================================================================
  // Recipe Versions
  // ============================================================================

  /**
   * Get all versions for a recipe
   * @param recipeId Recipe ID
   * @returns Array of versions
   */
  async getVersions(recipeId: string): Promise<RecipeVersion[]> {
    const { data } = await api.get(`/recipes/${recipeId}/versions`);
    return Array.isArray(data) ? data : [];
  },

  /**
   * Get a specific version by ID
   * @param recipeId Recipe ID
   * @param versionId Version ID
   * @returns Version details
   */
  async getVersion(recipeId: string, versionId: string): Promise<RecipeVersion> {
    const { data } = await api.get(`/recipes/${recipeId}/versions/${versionId}`);
    return data;
  },

  /**
   * Create a new version for a recipe
   * @param recipeId Recipe ID
   * @param version Version data
   * @returns Created version
   */
  async createVersion(
    recipeId: string,
    version: CreateRecipeVersionRequest
  ): Promise<RecipeVersion> {
    const { data } = await api.post(`/recipes/${recipeId}/versions`, version);
    return data;
  },

  /**
   * Activate a specific version
   * @param recipeId Recipe ID
   * @param versionNumber Version number to activate
   * @returns Activated version
   */
  async activateVersion(recipeId: string, versionNumber: number): Promise<RecipeVersion> {
    const { data } = await api.post(
      `/recipes/${recipeId}/versions/${versionNumber}/activate`
    );
    return data;
  },

  /**
   * Duplicate a version
   * @param recipeId Recipe ID
   * @param versionId Version ID to duplicate
   * @returns New version
   */
  async duplicateVersion(recipeId: string, versionId: string): Promise<RecipeVersion> {
    const { data } = await api.post(`/recipes/${recipeId}/versions/${versionId}/duplicate`);
    return data;
  },

  // ============================================================================
  // Recipe Ingredients
  // ============================================================================

  /**
   * Add an ingredient to a version
   * @param recipeId Recipe ID
   * @param versionId Version ID
   * @param ingredient Ingredient data
   * @returns Created ingredient
   */
  async addIngredient(
    recipeId: string,
    versionId: string,
    ingredient: CreateRecipeIngredientRequest
  ): Promise<RecipeIngredient> {
    const { data } = await api.post(
      `/recipes/${recipeId}/versions/${versionId}/ingredients`,
      ingredient
    );
    return data;
  },

  /**
   * Update an ingredient quantity
   * @param recipeId Recipe ID
   * @param versionId Version ID
   * @param ingredientId Ingredient ID
   * @param ingredient Updated ingredient data
   * @returns Updated ingredient
   */
  async updateIngredient(
    recipeId: string,
    versionId: string,
    ingredientId: string,
    ingredient: UpdateRecipeIngredientRequest
  ): Promise<RecipeIngredient> {
    const { data } = await api.put(
      `/recipes/${recipeId}/versions/${versionId}/ingredients/${ingredientId}`,
      ingredient
    );
    return data;
  },

  /**
   * Remove an ingredient from a version
   * @param recipeId Recipe ID
   * @param versionId Version ID
   * @param ingredientId Ingredient ID
   */
  async removeIngredient(
    recipeId: string,
    versionId: string,
    ingredientId: string
  ): Promise<void> {
    await api.delete(
      `/recipes/${recipeId}/versions/${versionId}/ingredients/${ingredientId}`
    );
  },

  // ============================================================================
  // Recipe Notes
  // ============================================================================

  /**
   * Get all notes for a recipe
   * @param recipeId Recipe ID
   * @returns Array of notes
   */
  async getNotes(recipeId: string): Promise<RecipeNote[]> {
    const { data } = await api.get(`/recipes/${recipeId}/notes`);
    return Array.isArray(data) ? data : [];
  },

  /**
   * Create a new note for a recipe
   * @param recipeId Recipe ID
   * @param note Note data
   * @returns Created note
   */
  async createNote(recipeId: string, note: CreateRecipeNoteRequest): Promise<RecipeNote> {
    const { data } = await api.post(`/recipes/${recipeId}/notes`, note);
    return data;
  },

  /**
   * Update a note
   * @param recipeId Recipe ID
   * @param noteId Note ID
   * @param note Updated note data
   * @returns Updated note
   */
  async updateNote(
    recipeId: string,
    noteId: string,
    note: UpdateRecipeNoteRequest
  ): Promise<RecipeNote> {
    const { data } = await api.put(`/recipes/${recipeId}/notes/${noteId}`, note);
    return data;
  },

  /**
   * Delete a note
   * @param recipeId Recipe ID
   * @param noteId Note ID
   */
  async deleteNote(recipeId: string, noteId: string): Promise<void> {
    await api.delete(`/recipes/${recipeId}/notes/${noteId}`);
  },

  // ============================================================================
  // Recipe Tags
  // ============================================================================

  /**
   * Get all tags for a recipe
   * @param recipeId Recipe ID
   * @returns Array of tags
   */
  async getTags(recipeId: string): Promise<RecipeTag[]> {
    const { data } = await api.get(`/recipes/${recipeId}/tags`);
    return Array.isArray(data) ? data : [];
  },

  /**
   * Add a tag to a recipe
   * @param recipeId Recipe ID
   * @param tag Tag data
   * @returns Created tag
   */
  async addTag(recipeId: string, tag: AddRecipeTagRequest): Promise<RecipeTag> {
    const { data } = await api.post(`/recipes/${recipeId}/tags`, tag);
    return data;
  },

  /**
   * Remove a tag from a recipe
   * @param recipeId Recipe ID
   * @param tag Tag name to remove
   */
  async removeTag(recipeId: string, tag: string): Promise<void> {
    await api.delete(`/recipes/${recipeId}/tags/${tag}`);
  },

  // ============================================================================
  // Recipe Collections
  // ============================================================================

  /**
   * Get all collections for the current user
   * @param filters Optional filters (pagination)
   * @returns Array of collections
   */
  async getCollections(filters?: RecipeCollectionFilters): Promise<RecipeCollection[]> {
    const params: Record<string, any> = {};
    if (filters?.limit !== undefined) params.limit = filters.limit;
    if (filters?.offset !== undefined) params.offset = filters.offset;

    const { data } = await api.get('/collections', { params });
    return Array.isArray(data) ? data : [];
  },

  /**
   * Get a specific collection by ID
   * @param id Collection ID
   * @returns Collection details
   */
  async getCollection(id: string): Promise<RecipeCollection> {
    const { data } = await api.get(`/collections/${id}`);
    return data;
  },

  /**
   * Create a new collection
   * @param collection Collection data
   * @returns Created collection
   */
  async createCollection(
    collection: CreateRecipeCollectionRequest
  ): Promise<RecipeCollection> {
    const { data } = await api.post('/collections', collection);
    return data;
  },

  /**
   * Update a collection
   * @param id Collection ID
   * @param collection Updated collection data
   * @returns Updated collection
   */
  async updateCollection(
    id: string,
    collection: UpdateRecipeCollectionRequest
  ): Promise<RecipeCollection> {
    const { data } = await api.put(`/collections/${id}`, collection);
    return data;
  },

  /**
   * Delete a collection
   * @param id Collection ID
   */
  async deleteCollection(id: string): Promise<void> {
    await api.delete(`/collections/${id}`);
  },

  /**
   * Add a recipe to a collection
   * @param collectionId Collection ID
   * @param request Request with recipe ID
   */
  async addRecipeToCollection(
    collectionId: string,
    request: AddRecipeToCollectionRequest
  ): Promise<void> {
    await api.post(`/collections/${collectionId}/recipes`, request);
  },

  /**
   * Remove a recipe from a collection
   * @param collectionId Collection ID
   * @param recipeId Recipe ID
   */
  async removeRecipeFromCollection(collectionId: string, recipeId: string): Promise<void> {
    await api.delete(`/collections/${collectionId}/recipes/${recipeId}`);
  },
};

export default recipeService;
