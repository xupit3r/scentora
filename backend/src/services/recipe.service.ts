import { Decimal } from '@prisma/client/runtime/library';
import { prisma } from '../prisma.js';
import { NotFoundError, ValidationError } from '../utils/errors.js';

function recipeToApi(recipe: any) {
  return {
    _id: recipe.id,
    userId: recipe.userId,
    name: recipe.name,
    ...(recipe.description != null ? { description: recipe.description } : {}),
    targetVolumeMl: Number(recipe.targetVolumeMl),
    status: recipe.status,
    ...(recipe.activeVersionId != null ? { activeVersionId: recipe.activeVersionId } : {}),
    createdAt: recipe.createdAt,
    updatedAt: recipe.updatedAt,
  };
}

function recipeResponseToApi(recipe: any, tags: string[], activeVersion: any | null) {
  return {
    _id: recipe.id,
    userId: recipe.userId,
    name: recipe.name,
    ...(recipe.description != null ? { description: recipe.description } : {}),
    targetVolumeMl: Number(recipe.targetVolumeMl),
    status: recipe.status,
    ...(recipe.activeVersionId != null ? { activeVersionId: recipe.activeVersionId } : {}),
    ...(activeVersion != null ? { activeVersion } : {}),
    ...(tags.length > 0 ? { tags } : {}),
    createdAt: recipe.createdAt,
    updatedAt: recipe.updatedAt,
  };
}

export async function createRecipe(userId: string, data: {
  name: string;
  description?: string;
  targetVolumeMl: number;
  status?: string;
}) {
  // Check name uniqueness
  const exists = await prisma.recipe.findFirst({ where: { userId, name: data.name } });
  if (exists) throw new ValidationError('recipe with this name already exists');

  const recipe = await prisma.recipe.create({
    data: {
      userId,
      name: data.name,
      description: data.description ?? null,
      targetVolumeMl: new Decimal(data.targetVolumeMl),
      status: data.status || 'draft',
    },
  });

  return recipeToApi(recipe);
}

export async function getRecipe(recipeId: string, userId: string) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  // Get tags
  const tags = await prisma.recipeTag.findMany({ where: { recipeId } });
  const tagNames = tags.map((t) => t.tag);

  // Get active version
  let activeVersion = null;
  if (recipe.activeVersionId) {
    const version = await prisma.recipeVersion.findUnique({ where: { id: recipe.activeVersionId } });
    if (version) {
      activeVersion = {
        _id: version.id,
        recipeId: version.recipeId,
        versionNumber: version.versionNumber,
        name: version.name,
        ...(version.notes != null ? { notes: version.notes } : {}),
        isActive: version.isActive,
        createdAt: version.createdAt,
      };
    }
  }

  return recipeResponseToApi(recipe, tagNames, activeVersion);
}

export async function listRecipes(userId: string, status?: string, limit = 100, offset = 0) {
  const where: any = { userId };
  if (status) where.status = status;

  const recipes = await prisma.recipe.findMany({
    where,
    orderBy: { createdAt: 'desc' },
    take: limit,
    skip: offset,
  });

  return recipes.map(recipeToApi);
}

export async function searchRecipes(userId: string, query: string, limit = 100, offset = 0) {
  const recipes = await prisma.recipe.findMany({
    where: {
      userId,
      OR: [
        { name: { contains: query, mode: 'insensitive' } },
        { description: { contains: query, mode: 'insensitive' } },
      ],
    },
    orderBy: { createdAt: 'desc' },
    take: limit,
    skip: offset,
  });

  return recipes.map(recipeToApi);
}

export async function updateRecipe(recipeId: string, userId: string, data: {
  name?: string;
  description?: string | null;
  targetVolumeMl?: number;
  status?: string;
}) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  // Check name uniqueness on rename
  if (data.name && data.name !== recipe.name) {
    const exists = await prisma.recipe.findFirst({
      where: { userId, name: data.name, id: { not: recipeId } },
    });
    if (exists) throw new ValidationError('recipe with this name already exists');
  }

  const updateData: any = {};
  if (data.name !== undefined) updateData.name = data.name;
  if (data.description !== undefined) updateData.description = data.description;
  if (data.targetVolumeMl !== undefined) updateData.targetVolumeMl = new Decimal(data.targetVolumeMl);
  if (data.status !== undefined) updateData.status = data.status;

  const updated = await prisma.recipe.update({ where: { id: recipeId }, data: updateData });
  return recipeToApi(updated);
}

export async function deleteRecipe(recipeId: string, userId: string) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');
  await prisma.recipe.delete({ where: { id: recipeId } });
}
