import { prisma } from '../prisma.js';
import { NotFoundError, ValidationError } from '../utils/errors.js';

function versionToApi(version: any) {
  return {
    _id: version.id,
    recipeId: version.recipeId,
    versionNumber: version.versionNumber,
    name: version.name,
    ...(version.notes != null ? { notes: version.notes } : {}),
    isActive: version.isActive,
    createdAt: version.createdAt,
  };
}

function ingredientToApi(ing: any) {
  return {
    _id: ing.id,
    versionId: ing.versionId,
    accordId: ing.accordId,
    ...(ing.accord?.name ? { accordName: ing.accord.name } : {}),
    quantityMl: Number(ing.quantityMl),
    quantityDrops: ing.quantityDrops ?? Math.round(Number(ing.quantityMl) * 20),
    ...(ing.percentage != null ? { percentage: Number(ing.percentage) } : {}),
    ...(ing.notes != null ? { notes: ing.notes } : {}),
    createdAt: ing.createdAt,
  };
}

function versionWithIngredientsToApi(version: any, ingredients: any[]) {
  return {
    ...versionToApi(version),
    ingredients: ingredients.map(ingredientToApi),
  };
}

export async function createVersion(recipeId: string, userId: string, data: { name: string; notes?: string }) {
  // Verify recipe ownership
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  // Get next version number
  const count = await prisma.recipeVersion.count({ where: { recipeId } });

  const version = await prisma.recipeVersion.create({
    data: {
      recipeId,
      versionNumber: count + 1,
      name: data.name,
      notes: data.notes ?? null,
      isActive: false,
    },
  });

  // Auto-activate if first version
  if (count === 0) {
    await prisma.recipeVersion.update({ where: { id: version.id }, data: { isActive: true } });
    await prisma.recipe.update({ where: { id: recipeId }, data: { activeVersionId: version.id } });
    version.isActive = true;
  }

  return versionToApi(version);
}

export async function getVersion(versionId: string) {
  const version = await prisma.recipeVersion.findUnique({ where: { id: versionId } });
  if (!version) throw new NotFoundError('version not found');

  const ingredients = await prisma.recipeIngredient.findMany({
    where: { versionId },
    include: { accord: { select: { name: true } } },
    orderBy: { createdAt: 'asc' },
  });

  return versionWithIngredientsToApi(version, ingredients);
}

export async function listVersions(recipeId: string, userId: string) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  const versions = await prisma.recipeVersion.findMany({
    where: { recipeId },
    orderBy: { versionNumber: 'asc' },
  });

  return versions.map(versionToApi);
}

export async function activateVersion(recipeId: string, versionNumber: number, userId: string) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  const version = await prisma.recipeVersion.findFirst({
    where: { recipeId, versionNumber },
  });
  if (!version) throw new NotFoundError('Version not found');

  // Deactivate all, activate target
  await prisma.recipeVersion.updateMany({ where: { recipeId }, data: { isActive: false } });
  await prisma.recipeVersion.update({ where: { id: version.id }, data: { isActive: true } });
  await prisma.recipe.update({ where: { id: recipeId }, data: { activeVersionId: version.id } });

  return { message: 'Version activated successfully' };
}

export async function duplicateVersion(versionId: string, userId: string) {
  const version = await prisma.recipeVersion.findUnique({ where: { id: versionId } });
  if (!version) throw new NotFoundError('version not found');

  const recipe = await prisma.recipe.findFirst({ where: { id: version.recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  // Get ingredients
  const ingredients = await prisma.recipeIngredient.findMany({ where: { versionId } });

  // Get next version number
  const count = await prisma.recipeVersion.count({ where: { recipeId: version.recipeId } });

  const newVersion = await prisma.recipeVersion.create({
    data: {
      recipeId: version.recipeId,
      versionNumber: count + 1,
      name: `v${count + 1}`,
      notes: version.notes,
      isActive: false,
    },
  });

  // Copy ingredients
  if (ingredients.length > 0) {
    await prisma.recipeIngredient.createMany({
      data: ingredients.map((ing) => ({
        versionId: newVersion.id,
        accordId: ing.accordId,
        quantityMl: ing.quantityMl,
        percentage: ing.percentage,
        notes: ing.notes,
      })),
    });
  }

  return versionToApi(newVersion);
}
