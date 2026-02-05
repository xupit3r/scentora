import { Decimal } from '@prisma/client/runtime/library';
import { prisma } from '../prisma.js';
import { NotFoundError, ValidationError } from '../utils/errors.js';

function ingredientToApi(ing: any) {
  return {
    _id: ing.id,
    versionId: ing.versionId,
    accordId: ing.accordId,
    quantityMl: Number(ing.quantityMl),
    quantityDrops: ing.quantityDrops ?? Math.round(Number(ing.quantityMl) * 20),
    ...(ing.percentage != null ? { percentage: Number(ing.percentage) } : {}),
    ...(ing.notes != null ? { notes: ing.notes } : {}),
    createdAt: ing.createdAt,
  };
}

export async function addIngredient(versionId: string, userId: string, data: {
  accordId: string;
  quantityMl: number;
  percentage?: number;
  notes?: string;
}) {
  // Verify version -> recipe -> user chain
  const version = await prisma.recipeVersion.findUnique({ where: { id: versionId } });
  if (!version) throw new NotFoundError('version not found');

  const recipe = await prisma.recipe.findFirst({ where: { id: version.recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  // Verify accord belongs to user
  const accord = await prisma.accord.findFirst({ where: { id: data.accordId, userId } });
  if (!accord) throw new NotFoundError('accord not found');

  // Check for duplicate
  const existing = await prisma.recipeIngredient.findFirst({
    where: { versionId, accordId: data.accordId },
  });
  if (existing) throw new ValidationError('this accord is already in this version');

  // Volume validation if enabled
  const user = await prisma.user.findUnique({ where: { id: userId } });
  if (user?.validateRecipeVolumes) {
    if (data.quantityMl > Number(accord.volumeMl)) {
      throw new ValidationError(
        `insufficient accord volume: need ${data.quantityMl.toFixed(2)} ml, have ${Number(accord.volumeMl).toFixed(2)} ml`,
      );
    }
  }

  const ingredient = await prisma.recipeIngredient.create({
    data: {
      versionId,
      accordId: data.accordId,
      quantityMl: new Decimal(data.quantityMl),
      percentage: data.percentage != null ? new Decimal(data.percentage) : null,
      notes: data.notes ?? null,
    },
  });

  return ingredientToApi(ingredient);
}

export async function updateIngredient(ingredientId: string, userId: string, data: {
  quantityMl?: number;
  percentage?: number | null;
  notes?: string | null;
}) {
  const ingredient = await prisma.recipeIngredient.findUnique({ where: { id: ingredientId } });
  if (!ingredient) throw new NotFoundError('ingredient not found');

  const version = await prisma.recipeVersion.findUnique({ where: { id: ingredient.versionId } });
  if (!version) throw new NotFoundError('version not found');

  const recipe = await prisma.recipe.findFirst({ where: { id: version.recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  const updateData: any = {};
  if (data.quantityMl !== undefined) {
    if (data.quantityMl <= 0) throw new ValidationError('quantity must be greater than 0');

    // Volume validation if enabled
    const user = await prisma.user.findUnique({ where: { id: userId } });
    if (user?.validateRecipeVolumes) {
      const accord = await prisma.accord.findFirst({ where: { id: ingredient.accordId, userId } });
      if (accord && data.quantityMl > Number(accord.volumeMl)) {
        throw new ValidationError(
          `insufficient accord volume: need ${data.quantityMl.toFixed(2)} ml, have ${Number(accord.volumeMl).toFixed(2)} ml`,
        );
      }
    }

    updateData.quantityMl = new Decimal(data.quantityMl);
  }
  if (data.percentage !== undefined) {
    updateData.percentage = data.percentage != null ? new Decimal(data.percentage) : null;
  }
  if (data.notes !== undefined) updateData.notes = data.notes;

  const updated = await prisma.recipeIngredient.update({
    where: { id: ingredientId },
    data: updateData,
  });

  return ingredientToApi(updated);
}

export async function removeIngredient(ingredientId: string, userId: string) {
  const ingredient = await prisma.recipeIngredient.findUnique({ where: { id: ingredientId } });
  if (!ingredient) throw new NotFoundError('ingredient not found');

  const version = await prisma.recipeVersion.findUnique({ where: { id: ingredient.versionId } });
  if (!version) throw new NotFoundError('version not found');

  const recipe = await prisma.recipe.findFirst({ where: { id: version.recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  await prisma.recipeIngredient.delete({ where: { id: ingredientId } });
}
