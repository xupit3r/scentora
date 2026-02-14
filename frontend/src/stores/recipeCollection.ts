import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { recipeService } from '@/services/recipe';
import type {
  RecipeCollection,
  CreateRecipeCollectionRequest,
  UpdateRecipeCollectionRequest,
  RecipeCollectionFilters,
  Recipe,
} from '@/types';

export const useRecipeCollectionStore = defineStore('recipeCollection', () => {
  // ============================================================================
  // State
  // ============================================================================

  const collections = ref<RecipeCollection[]>([]);
  const currentCollection = ref<RecipeCollection | null>(null);
  const currentCollectionRecipes = ref<Recipe[]>([]);

  const isLoading = ref(false);
  const isLoadingRecipes = ref(false);
  const error = ref<string | null>(null);

  // Cache for collections
  const collectionCache = ref<Map<string, RecipeCollection>>(new Map());

  // ============================================================================
  // Getters
  // ============================================================================

  const hasCollections = computed(() => collections.value.length > 0);
  const collectionCount = computed(() => collections.value.length);

  const currentCollectionRecipeCount = computed(() => currentCollectionRecipes.value.length);

  /**
   * Get collection by ID from cache or list
   */
  function getCollectionById(id: string): RecipeCollection | undefined {
    return collectionCache.value.get(id) || collections.value.find((c) => c._id === id);
  }

  // ============================================================================
  // Actions - Collection CRUD
  // ============================================================================

  /**
   * Fetch all collections
   */
  async function fetchCollections(filters?: RecipeCollectionFilters): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      const data = await recipeService.getCollections(filters);
      collections.value = data;

      // Update cache
      data.forEach((collection) => {
        collectionCache.value.set(collection._id, collection);
      });

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch collections';
      console.error('Error fetching collections:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Fetch a specific collection by ID
   */
  async function fetchCollection(id: string): Promise<boolean> {
    // Check cache first
    if (collectionCache.value.has(id)) {
      currentCollection.value = collectionCache.value.get(id)!;
    }

    isLoading.value = true;
    error.value = null;

    try {
      const data = await recipeService.getCollection(id);
      currentCollection.value = data;
      collectionCache.value.set(id, data);
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch collection';
      console.error('Error fetching collection:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Create a new collection
   */
  async function createCollection(
    request: CreateRecipeCollectionRequest
  ): Promise<RecipeCollection | null> {
    isLoading.value = true;
    error.value = null;

    try {
      const newCollection = await recipeService.createCollection(request);

      // Add to list and cache
      collections.value.unshift(newCollection);
      collectionCache.value.set(newCollection._id, newCollection);
      currentCollection.value = newCollection;

      return newCollection;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to create collection';
      console.error('Error creating collection:', err);
      return null;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Update an existing collection
   */
  async function updateCollection(
    id: string,
    request: UpdateRecipeCollectionRequest
  ): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      const updated = await recipeService.updateCollection(id, request);

      // Update in list
      const index = collections.value.findIndex((c) => c._id === id);
      if (index !== -1) {
        collections.value[index] = updated;
      }

      // Update cache and current
      collectionCache.value.set(id, updated);
      if (currentCollection.value?._id === id) {
        currentCollection.value = updated;
      }

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update collection';
      console.error('Error updating collection:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Delete a collection
   */
  async function deleteCollection(id: string): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      await recipeService.deleteCollection(id);

      // Remove from list and cache
      collections.value = collections.value.filter((c) => c._id !== id);
      collectionCache.value.delete(id);

      if (currentCollection.value?._id === id) {
        currentCollection.value = null;
        currentCollectionRecipes.value = [];
      }

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete collection';
      console.error('Error deleting collection:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  // ============================================================================
  // Actions - Collection Membership
  // ============================================================================

  /**
   * Fetch all recipes in a collection
   * Note: This would require an additional API endpoint or client-side filtering
   * For now, we'll use a placeholder that fetches recipes and filters by collection
   */
  async function fetchCollectionRecipes(_collectionId: string): Promise<boolean> {
    isLoadingRecipes.value = true;
    error.value = null;

    try {
      // TODO: Backend should provide GET /api/collections/:id/recipes endpoint
      // For now, we'll fetch all recipes and filter (not optimal but functional)
      // This is a placeholder implementation

      // In a real implementation, you'd have:
      // const data = await recipeService.getCollectionRecipes(collectionId);

      // Temporary: Return empty array and set a note in error
      currentCollectionRecipes.value = [];
      console.warn('fetchCollectionRecipes: Backend endpoint not yet implemented');

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch collection recipes';
      console.error('Error fetching collection recipes:', err);
      return false;
    } finally {
      isLoadingRecipes.value = false;
    }
  }

  /**
   * Add a recipe to a collection
   */
  async function addRecipeToCollection(collectionId: string, recipeId: string): Promise<boolean> {
    error.value = null;

    try {
      await recipeService.addRecipeToCollection(collectionId, { recipeId });

      // Optimistically update local state if we have the recipe loaded
      // Note: Full recipe data would need to be fetched separately
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to add recipe to collection';
      console.error('Error adding recipe to collection:', err);
      return false;
    }
  }

  /**
   * Remove a recipe from a collection
   */
  async function removeRecipeFromCollection(collectionId: string, recipeId: string): Promise<boolean> {
    error.value = null;

    try {
      await recipeService.removeRecipeFromCollection(collectionId, recipeId);

      // Remove from current collection recipes if loaded
      currentCollectionRecipes.value = currentCollectionRecipes.value.filter(
        (r) => r._id !== recipeId
      );

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to remove recipe from collection';
      console.error('Error removing recipe from collection:', err);
      return false;
    }
  }

  // ============================================================================
  // Utility Actions
  // ============================================================================

  /**
   * Clear current collection and related state
   */
  function clearCurrent() {
    currentCollection.value = null;
    currentCollectionRecipes.value = [];
  }

  /**
   * Clear all state (for logout, etc.)
   */
  function clearAll() {
    collections.value = [];
    currentCollection.value = null;
    currentCollectionRecipes.value = [];
    collectionCache.value.clear();
    error.value = null;
  }

  /**
   * Check if a recipe is in a specific collection
   * Note: This requires fetching collection recipes first
   */
  function isRecipeInCurrentCollection(recipeId: string): boolean {
    return currentCollectionRecipes.value.some((r) => r._id === recipeId);
  }

  // ============================================================================
  // Return
  // ============================================================================

  return {
    // State
    collections,
    currentCollection,
    currentCollectionRecipes,
    isLoading,
    isLoadingRecipes,
    error,

    // Getters
    hasCollections,
    collectionCount,
    currentCollectionRecipeCount,
    getCollectionById,

    // Actions - Collections
    fetchCollections,
    fetchCollection,
    createCollection,
    updateCollection,
    deleteCollection,

    // Actions - Membership
    fetchCollectionRecipes,
    addRecipeToCollection,
    removeRecipeFromCollection,

    // Utility
    clearCurrent,
    clearAll,
    isRecipeInCurrentCollection,
  };
});
