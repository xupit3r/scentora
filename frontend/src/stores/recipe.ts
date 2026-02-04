import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { recipeService } from '@/services/recipe';
import type {
  Recipe,
  CreateRecipeRequest,
  UpdateRecipeRequest,
  RecipeFilters,
  RecipeStatus,
  RecipeVersion,
  CreateRecipeVersionRequest,
  RecipeIngredient,
  CreateRecipeIngredientRequest,
  UpdateRecipeIngredientRequest,
  RecipeNote,
  CreateRecipeNoteRequest,
  UpdateRecipeNoteRequest,
} from '@/types';

export const useRecipeStore = defineStore('recipe', () => {
  // ============================================================================
  // State
  // ============================================================================

  const recipes = ref<Recipe[]>([]);
  const currentRecipe = ref<Recipe | null>(null);
  const currentVersions = ref<RecipeVersion[]>([]);
  const currentIngredients = ref<RecipeIngredient[]>([]);
  const currentNotes = ref<RecipeNote[]>([]);
  const currentTags = ref<string[]>([]);

  const isLoading = ref(false);
  const isLoadingVersions = ref(false);
  const isLoadingIngredients = ref(false);
  const isLoadingNotes = ref(false);
  const error = ref<string | null>(null);

  // Cache for recipes (indexed by ID)
  const recipeCache = ref<Map<string, Recipe>>(new Map());
  const versionCache = ref<Map<string, RecipeVersion[]>>(new Map());

  // ============================================================================
  // Getters
  // ============================================================================

  const hasRecipes = computed(() => recipes.value.length > 0);

  const recipesByStatus = computed(() => (status: RecipeStatus) => {
    return recipes.value.filter((r) => r.status === status);
  });

  const draftRecipes = computed(() => recipesByStatus.value('draft'));
  const inProgressRecipes = computed(() => recipesByStatus.value('in_progress'));
  const testedRecipes = computed(() => recipesByStatus.value('tested'));
  const finalizedRecipes = computed(() => recipesByStatus.value('finalized'));
  const archivedRecipes = computed(() => recipesByStatus.value('archived'));

  const recipeCount = computed(() => recipes.value.length);

  const activeVersion = computed(() => {
    return currentVersions.value.find((v) => v.isActive) || null;
  });

  // ============================================================================
  // Actions - Recipe CRUD
  // ============================================================================

  /**
   * Fetch all recipes with optional filters
   */
  async function fetchRecipes(filters?: RecipeFilters): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      const data = await recipeService.getAll(filters);
      recipes.value = data;

      // Update cache
      data.forEach((recipe) => {
        recipeCache.value.set(recipe._id, recipe);
      });

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch recipes';
      console.error('Error fetching recipes:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Fetch a specific recipe by ID
   */
  async function fetchRecipe(id: string): Promise<boolean> {
    // Check cache first
    if (recipeCache.value.has(id)) {
      currentRecipe.value = recipeCache.value.get(id)!;
    }

    isLoading.value = true;
    error.value = null;

    try {
      const data = await recipeService.getById(id);
      currentRecipe.value = data;
      recipeCache.value.set(id, data);

      // Also fetch versions and notes
      await Promise.all([
        fetchVersions(id),
        fetchNotes(id),
        fetchTags(id),
      ]);

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch recipe';
      console.error('Error fetching recipe:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Create a new recipe
   */
  async function createRecipe(request: CreateRecipeRequest): Promise<Recipe | null> {
    isLoading.value = true;
    error.value = null;

    try {
      const newRecipe = await recipeService.create(request);

      // Add to list and cache
      recipes.value.unshift(newRecipe);
      recipeCache.value.set(newRecipe._id, newRecipe);
      currentRecipe.value = newRecipe;

      return newRecipe;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to create recipe';
      console.error('Error creating recipe:', err);
      return null;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Update an existing recipe
   */
  async function updateRecipe(id: string, request: UpdateRecipeRequest): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      const updated = await recipeService.update(id, request);

      // Update in list
      const index = recipes.value.findIndex((r) => r._id === id);
      if (index !== -1) {
        recipes.value[index] = updated;
      }

      // Update cache and current
      recipeCache.value.set(id, updated);
      if (currentRecipe.value?._id === id) {
        currentRecipe.value = updated;
      }

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update recipe';
      console.error('Error updating recipe:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Delete a recipe
   */
  async function deleteRecipe(id: string): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      await recipeService.delete(id);

      // Remove from list and cache
      recipes.value = recipes.value.filter((r) => r._id !== id);
      recipeCache.value.delete(id);
      versionCache.value.delete(id);

      if (currentRecipe.value?._id === id) {
        currentRecipe.value = null;
        currentVersions.value = [];
        currentIngredients.value = [];
        currentNotes.value = [];
        currentTags.value = [];
      }

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete recipe';
      console.error('Error deleting recipe:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Search recipes
   */
  async function searchRecipes(query: string): Promise<boolean> {
    if (!query.trim()) {
      return fetchRecipes();
    }

    isLoading.value = true;
    error.value = null;

    try {
      const data = await recipeService.search(query);
      recipes.value = data;
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Search failed';
      console.error('Error searching recipes:', err);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  // ============================================================================
  // Actions - Recipe Versions
  // ============================================================================

  /**
   * Fetch all versions for a recipe
   */
  async function fetchVersions(recipeId: string): Promise<boolean> {
    // Check cache
    if (versionCache.value.has(recipeId)) {
      currentVersions.value = versionCache.value.get(recipeId)!;
    }

    isLoadingVersions.value = true;
    error.value = null;

    try {
      const data = await recipeService.getVersions(recipeId);
      currentVersions.value = data;
      versionCache.value.set(recipeId, data);
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch versions';
      console.error('Error fetching versions:', err);
      return false;
    } finally {
      isLoadingVersions.value = false;
    }
  }

  /**
   * Create a new version
   */
  async function createVersion(
    recipeId: string,
    request: CreateRecipeVersionRequest
  ): Promise<RecipeVersion | null> {
    isLoadingVersions.value = true;
    error.value = null;

    try {
      const newVersion = await recipeService.createVersion(recipeId, request);

      // Add to current versions
      currentVersions.value.unshift(newVersion);

      // Update cache
      versionCache.value.set(recipeId, currentVersions.value);

      return newVersion;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to create version';
      console.error('Error creating version:', err);
      return null;
    } finally {
      isLoadingVersions.value = false;
    }
  }

  /**
   * Activate a version
   */
  async function activateVersion(recipeId: string, versionNumber: number): Promise<boolean> {
    isLoadingVersions.value = true;
    error.value = null;

    try {
      await recipeService.activateVersion(recipeId, versionNumber);

      // Update local state
      currentVersions.value = currentVersions.value.map((v) => ({
        ...v,
        isActive: v.versionNumber === versionNumber,
      }));

      // Update cache
      versionCache.value.set(recipeId, currentVersions.value);

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to activate version';
      console.error('Error activating version:', err);
      return false;
    } finally {
      isLoadingVersions.value = false;
    }
  }

  /**
   * Duplicate a version
   */
  async function duplicateVersion(recipeId: string, versionId: string): Promise<RecipeVersion | null> {
    isLoadingVersions.value = true;
    error.value = null;

    try {
      const newVersion = await recipeService.duplicateVersion(recipeId, versionId);

      // Refresh versions to get updated list
      await fetchVersions(recipeId);

      return newVersion;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to duplicate version';
      console.error('Error duplicating version:', err);
      return null;
    } finally {
      isLoadingVersions.value = false;
    }
  }

  // ============================================================================
  // Actions - Recipe Ingredients
  // ============================================================================

  /**
   * Add an ingredient to a version
   */
  async function addIngredient(
    recipeId: string,
    versionId: string,
    request: CreateRecipeIngredientRequest
  ): Promise<RecipeIngredient | null> {
    isLoadingIngredients.value = true;
    error.value = null;

    try {
      const newIngredient = await recipeService.addIngredient(recipeId, versionId, request);
      currentIngredients.value.push(newIngredient);
      return newIngredient;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to add ingredient';
      console.error('Error adding ingredient:', err);
      return null;
    } finally {
      isLoadingIngredients.value = false;
    }
  }

  /**
   * Update an ingredient
   */
  async function updateIngredient(
    recipeId: string,
    versionId: string,
    ingredientId: string,
    request: UpdateRecipeIngredientRequest
  ): Promise<boolean> {
    isLoadingIngredients.value = true;
    error.value = null;

    try {
      const updated = await recipeService.updateIngredient(recipeId, versionId, ingredientId, request);

      // Update in list
      const index = currentIngredients.value.findIndex((i) => i._id === ingredientId);
      if (index !== -1) {
        currentIngredients.value[index] = updated;
      }

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update ingredient';
      console.error('Error updating ingredient:', err);
      return false;
    } finally {
      isLoadingIngredients.value = false;
    }
  }

  /**
   * Remove an ingredient
   */
  async function removeIngredient(
    recipeId: string,
    versionId: string,
    ingredientId: string
  ): Promise<boolean> {
    isLoadingIngredients.value = true;
    error.value = null;

    try {
      await recipeService.removeIngredient(recipeId, versionId, ingredientId);
      currentIngredients.value = currentIngredients.value.filter((i) => i._id !== ingredientId);
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to remove ingredient';
      console.error('Error removing ingredient:', err);
      return false;
    } finally {
      isLoadingIngredients.value = false;
    }
  }

  // ============================================================================
  // Actions - Recipe Notes
  // ============================================================================

  /**
   * Fetch all notes for a recipe
   */
  async function fetchNotes(recipeId: string): Promise<boolean> {
    isLoadingNotes.value = true;
    error.value = null;

    try {
      const data = await recipeService.getNotes(recipeId);
      currentNotes.value = data;
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch notes';
      console.error('Error fetching notes:', err);
      return false;
    } finally {
      isLoadingNotes.value = false;
    }
  }

  /**
   * Create a new note
   */
  async function createNote(recipeId: string, request: CreateRecipeNoteRequest): Promise<RecipeNote | null> {
    isLoadingNotes.value = true;
    error.value = null;

    try {
      const newNote = await recipeService.createNote(recipeId, request);
      currentNotes.value.unshift(newNote);
      return newNote;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to create note';
      console.error('Error creating note:', err);
      return null;
    } finally {
      isLoadingNotes.value = false;
    }
  }

  /**
   * Update a note
   */
  async function updateNote(
    recipeId: string,
    noteId: string,
    request: UpdateRecipeNoteRequest
  ): Promise<boolean> {
    isLoadingNotes.value = true;
    error.value = null;

    try {
      const updated = await recipeService.updateNote(recipeId, noteId, request);

      // Update in list
      const index = currentNotes.value.findIndex((n) => n._id === noteId);
      if (index !== -1) {
        currentNotes.value[index] = updated;
      }

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update note';
      console.error('Error updating note:', err);
      return false;
    } finally {
      isLoadingNotes.value = false;
    }
  }

  /**
   * Delete a note
   */
  async function deleteNote(recipeId: string, noteId: string): Promise<boolean> {
    isLoadingNotes.value = true;
    error.value = null;

    try {
      await recipeService.deleteNote(recipeId, noteId);
      currentNotes.value = currentNotes.value.filter((n) => n._id !== noteId);
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete note';
      console.error('Error deleting note:', err);
      return false;
    } finally {
      isLoadingNotes.value = false;
    }
  }

  // ============================================================================
  // Actions - Recipe Tags
  // ============================================================================

  /**
   * Fetch all tags for a recipe
   */
  async function fetchTags(recipeId: string): Promise<boolean> {
    error.value = null;

    try {
      const data = await recipeService.getTags(recipeId);
      currentTags.value = data.map((t) => t.tag);
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch tags';
      console.error('Error fetching tags:', err);
      return false;
    }
  }

  /**
   * Add a tag
   */
  async function addTag(recipeId: string, tag: string): Promise<boolean> {
    error.value = null;

    try {
      await recipeService.addTag(recipeId, { tag });
      if (!currentTags.value.includes(tag)) {
        currentTags.value.push(tag);
      }
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to add tag';
      console.error('Error adding tag:', err);
      return false;
    }
  }

  /**
   * Remove a tag
   */
  async function removeTag(recipeId: string, tag: string): Promise<boolean> {
    error.value = null;

    try {
      await recipeService.removeTag(recipeId, tag);
      currentTags.value = currentTags.value.filter((t) => t !== tag);
      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to remove tag';
      console.error('Error removing tag:', err);
      return false;
    }
  }

  // ============================================================================
  // Utility Actions
  // ============================================================================

  /**
   * Clear current recipe and related state
   */
  function clearCurrent() {
    currentRecipe.value = null;
    currentVersions.value = [];
    currentIngredients.value = [];
    currentNotes.value = [];
    currentTags.value = [];
  }

  /**
   * Clear all state (for logout, etc.)
   */
  function clearAll() {
    recipes.value = [];
    currentRecipe.value = null;
    currentVersions.value = [];
    currentIngredients.value = [];
    currentNotes.value = [];
    currentTags.value = [];
    recipeCache.value.clear();
    versionCache.value.clear();
    error.value = null;
  }

  // ============================================================================
  // Return
  // ============================================================================

  return {
    // State
    recipes,
    currentRecipe,
    currentVersions,
    currentIngredients,
    currentNotes,
    currentTags,
    isLoading,
    isLoadingVersions,
    isLoadingIngredients,
    isLoadingNotes,
    error,

    // Getters
    hasRecipes,
    recipesByStatus,
    draftRecipes,
    inProgressRecipes,
    testedRecipes,
    finalizedRecipes,
    archivedRecipes,
    recipeCount,
    activeVersion,

    // Actions - Recipes
    fetchRecipes,
    fetchRecipe,
    createRecipe,
    updateRecipe,
    deleteRecipe,
    searchRecipes,

    // Actions - Versions
    fetchVersions,
    createVersion,
    activateVersion,
    duplicateVersion,

    // Actions - Ingredients
    addIngredient,
    updateIngredient,
    removeIngredient,

    // Actions - Notes
    fetchNotes,
    createNote,
    updateNote,
    deleteNote,

    // Actions - Tags
    fetchTags,
    addTag,
    removeTag,

    // Utility
    clearCurrent,
    clearAll,
  };
});
