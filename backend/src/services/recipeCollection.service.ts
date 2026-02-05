import { prisma } from '../prisma.js';
import { NotFoundError, ValidationError } from '../utils/errors.js';

function collectionToApi(collection: any) {
  return {
    _id: collection.id,
    userId: collection.userId,
    name: collection.name,
    ...(collection.description != null ? { description: collection.description } : {}),
    createdAt: collection.createdAt,
    updatedAt: collection.updatedAt,
  };
}

export async function createCollection(userId: string, data: { name: string; description?: string }) {
  const exists = await prisma.recipeCollection.findFirst({ where: { userId, name: data.name } });
  if (exists) throw new ValidationError('collection with this name already exists');

  const collection = await prisma.recipeCollection.create({
    data: {
      userId,
      name: data.name,
      description: data.description ?? null,
    },
  });

  return collectionToApi(collection);
}

export async function listCollections(userId: string, limit = 100, offset = 0) {
  const collections = await prisma.recipeCollection.findMany({
    where: { userId },
    orderBy: { createdAt: 'desc' },
    take: limit,
    skip: offset,
  });
  return collections.map(collectionToApi);
}

export async function getCollection(collectionId: string, userId: string) {
  const collection = await prisma.recipeCollection.findFirst({
    where: { id: collectionId, userId },
  });
  if (!collection) throw new NotFoundError('collection not found');
  return collectionToApi(collection);
}

export async function updateCollection(collectionId: string, userId: string, data: {
  name?: string;
  description?: string | null;
}) {
  const collection = await prisma.recipeCollection.findFirst({
    where: { id: collectionId, userId },
  });
  if (!collection) throw new NotFoundError('collection not found');

  // Check name uniqueness on rename
  if (data.name && data.name !== collection.name) {
    const exists = await prisma.recipeCollection.findFirst({
      where: { userId, name: data.name, id: { not: collectionId } },
    });
    if (exists) throw new ValidationError('collection with this name already exists');
  }

  const updateData: any = {};
  if (data.name !== undefined) updateData.name = data.name;
  if (data.description !== undefined) updateData.description = data.description;

  const updated = await prisma.recipeCollection.update({
    where: { id: collectionId },
    data: updateData,
  });

  return collectionToApi(updated);
}

export async function deleteCollection(collectionId: string, userId: string) {
  const collection = await prisma.recipeCollection.findFirst({
    where: { id: collectionId, userId },
  });
  if (!collection) throw new NotFoundError('collection not found');

  await prisma.recipeCollection.delete({ where: { id: collectionId } });
}

export async function addRecipe(collectionId: string, recipeId: string, userId: string) {
  const collection = await prisma.recipeCollection.findFirst({ where: { id: collectionId, userId } });
  if (!collection) throw new NotFoundError('collection not found');

  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  try {
    await prisma.recipeCollectionMember.create({
      data: { collectionId, recipeId },
    });
  } catch (err: any) {
    if (err.code !== 'P2002') throw err;
    // Already in collection, no-op
  }
}

export async function removeRecipe(collectionId: string, recipeId: string, userId: string) {
  const collection = await prisma.recipeCollection.findFirst({ where: { id: collectionId, userId } });
  if (!collection) throw new NotFoundError('collection not found');

  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  await prisma.recipeCollectionMember.deleteMany({ where: { collectionId, recipeId } });
}
