// Recipe System Types

export type RecipeStatus = 'draft' | 'in_progress' | 'tested' | 'finalized' | 'archived';
export type NoteType = 'general' | 'observation' | 'adjustment' | 'reminder';

// Recipe
export interface Recipe {
  _id: string;
  userId: string;
  name: string;
  description?: string;
  targetVolumeMl: number;
  status: RecipeStatus;
  activeVersionId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateRecipeRequest {
  name: string;
  description?: string;
  notes?: string;
  targetVolumeMl: number;
}

export interface UpdateRecipeRequest {
  name?: string;
  description?: string;
  notes?: string;
  status?: RecipeStatus;
}

export interface RecipeFilters {
  status?: RecipeStatus;
  search?: string;
  limit?: number;
  offset?: number;
}

// Recipe Version
export interface RecipeVersion {
  _id: string;
  recipeId: string;
  versionNumber: number;
  name: string;
  notes?: string;
  isActive: boolean;
  createdAt: string;
  updatedAt?: string;
}

export interface CreateRecipeVersionRequest {
  notes?: string;
  changes?: string;
}

// Recipe Ingredient
export interface RecipeIngredient {
  _id: string;
  versionId: string;
  accordId: string;
  quantityMl: number;
  createdAt: string;
  updatedAt?: string;
}

export interface CreateRecipeIngredientRequest {
  accordId: string;
  quantityMl: number;
}

export interface UpdateRecipeIngredientRequest {
  quantityMl: number;
}

// Recipe Note
export interface RecipeNote {
  _id: string;
  recipeId: string;
  content: string;
  noteType: NoteType;
  createdAt: string;
  updatedAt: string;
}

export interface CreateRecipeNoteRequest {
  content: string;
  type?: NoteType;
}

export interface UpdateRecipeNoteRequest {
  content: string;
}

// Recipe Tag
export interface RecipeTag {
  tag: string;
}

export interface AddRecipeTagRequest {
  tag: string;
}

// Recipe Collection
export interface RecipeCollection {
  _id: string;
  userId: string;
  name: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateRecipeCollectionRequest {
  name: string;
  description?: string;
}

export interface UpdateRecipeCollectionRequest {
  name?: string;
  description?: string;
}

export interface RecipeCollectionFilters {
  limit?: number;
  offset?: number;
}

export interface AddRecipeToCollectionRequest {
  recipeId: string;
}

// Extended Recipe with relationships (for detailed views)
export interface RecipeWithDetails extends Recipe {
  versions?: RecipeVersion[];
  notes?: RecipeNote[];
  tags?: RecipeTag[];
  activeVersion?: RecipeVersion;
}

// Version with ingredients (for version details)
export interface RecipeVersionWithIngredients extends RecipeVersion {
  ingredients?: RecipeIngredient[];
}

// Collection with recipes (for collection views)
export interface RecipeCollectionWithRecipes extends RecipeCollection {
  recipes?: Recipe[];
  recipeCount?: number;
}
