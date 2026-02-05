import { Decimal } from '@prisma/client/runtime/library';
import { prisma } from '../prisma.js';
import { NotFoundError, ValidationError } from '../utils/errors.js';

function accordToApi(accord: any) {
  return {
    _id: accord.id,
    userId: accord.userId,
    name: accord.name,
    pyramidPosition: accord.pyramidPosition,
    volumeMl: Number(accord.volumeMl),
    volumeDrops: accord.volumeDrops ?? Math.round(Number(accord.volumeMl) * 20),
    ...(accord.supplier != null ? { supplier: accord.supplier } : {}),
    ...(accord.purchaseDate != null ? { purchaseDate: accord.purchaseDate } : {}),
    ...(accord.dilutionPercentage != null
      ? { dilutionPercentage: Number(accord.dilutionPercentage) }
      : {}),
    ...(accord.notes != null ? { notes: accord.notes } : {}),
    tags: accord.tags?.map((t: any) => t.tag) ?? [],
    createdAt: accord.createdAt,
    updatedAt: accord.updatedAt,
  };
}

export async function createAccord(userId: string, data: {
  name: string;
  pyramidPosition: string;
  volumeMl: number;
  supplier?: string;
  purchaseDate?: string;
  dilutionPercentage?: number;
  notes?: string;
  tags?: string[];
}) {
  const accord = await prisma.accord.create({
    data: {
      userId,
      name: data.name,
      pyramidPosition: data.pyramidPosition,
      volumeMl: new Decimal(data.volumeMl),
      supplier: data.supplier ?? null,
      purchaseDate: data.purchaseDate ? new Date(data.purchaseDate) : null,
      dilutionPercentage: data.dilutionPercentage != null ? new Decimal(data.dilutionPercentage) : null,
      notes: data.notes ?? null,
      ...(data.tags && data.tags.length > 0
        ? { tags: { create: data.tags.map((tag) => ({ tag })) } }
        : {}),
    },
    include: { tags: true },
  });

  return accordToApi(accord);
}

export async function getAccord(accordId: string, userId: string) {
  const accord = await prisma.accord.findFirst({
    where: { id: accordId, userId },
    include: { tags: true },
  });
  if (!accord) throw new NotFoundError('Accord not found');
  return accordToApi(accord);
}

export async function listAccords(userId: string) {
  const accords = await prisma.accord.findMany({
    where: { userId },
    include: { tags: true },
    orderBy: { createdAt: 'desc' },
  });
  return accords.map(accordToApi);
}

export async function searchAccords(
  userId: string,
  filters: {
    position?: string;
    minVolume?: number;
    maxVolume?: number;
    supplier?: string;
    search?: string;
    tags?: string[];
  },
) {
  // If tags filter, do tag-based search (ALL tags must match)
  if (filters.tags && filters.tags.length > 0) {
    const accords = await prisma.$queryRaw<any[]>`
      SELECT a.* FROM accords a
      WHERE a.user_id = ${userId}::uuid
      AND (
        SELECT COUNT(DISTINCT at.tag)
        FROM accord_tags at
        WHERE at.accord_id = a.id AND at.tag = ANY(${filters.tags})
      ) = ${filters.tags.length}
      ORDER BY a.created_at DESC
    `;

    // Load tags for each accord
    const result = [];
    for (const accord of accords) {
      const tags = await prisma.accordTag.findMany({ where: { accordId: accord.id } });
      result.push(accordToApi({
        ...accord,
        volumeMl: accord.volume_ml,
        pyramidPosition: accord.pyramid_position,
        dilutionPercentage: accord.dilution_percentage,
        purchaseDate: accord.purchase_date,
        createdAt: accord.created_at,
        updatedAt: accord.updated_at,
        userId: accord.user_id,
        volumeDrops: accord.volume_drops,
        tags,
      }));
    }
    return result;
  }

  // Otherwise use Prisma filters
  const where: any = { userId };

  if (filters.position) where.pyramidPosition = filters.position;
  if (filters.supplier) where.supplier = { contains: filters.supplier, mode: 'insensitive' };
  if (filters.minVolume != null || filters.maxVolume != null) {
    where.volumeMl = {};
    if (filters.minVolume != null) where.volumeMl.gte = new Decimal(filters.minVolume);
    if (filters.maxVolume != null) where.volumeMl.lte = new Decimal(filters.maxVolume);
  }
  if (filters.search) {
    where.OR = [
      { name: { contains: filters.search, mode: 'insensitive' } },
      { notes: { contains: filters.search, mode: 'insensitive' } },
    ];
  }

  const accords = await prisma.accord.findMany({
    where,
    include: { tags: true },
    orderBy: { createdAt: 'desc' },
  });
  return accords.map(accordToApi);
}

export async function updateAccord(accordId: string, userId: string, data: {
  name?: string;
  pyramidPosition?: string;
  volumeMl?: number;
  supplier?: string | null;
  purchaseDate?: string | null;
  dilutionPercentage?: number | null;
  notes?: string | null;
  tags?: string[];
}) {
  const accord = await prisma.accord.findFirst({ where: { id: accordId, userId } });
  if (!accord) throw new NotFoundError('Accord not found');

  const updateData: any = {};
  if (data.name !== undefined) updateData.name = data.name;
  if (data.pyramidPosition !== undefined) updateData.pyramidPosition = data.pyramidPosition;
  if (data.volumeMl !== undefined) updateData.volumeMl = new Decimal(data.volumeMl);
  if (data.supplier !== undefined) updateData.supplier = data.supplier;
  if (data.purchaseDate !== undefined) {
    updateData.purchaseDate = data.purchaseDate ? new Date(data.purchaseDate) : null;
  }
  if (data.dilutionPercentage !== undefined) {
    updateData.dilutionPercentage = data.dilutionPercentage != null ? new Decimal(data.dilutionPercentage) : null;
  }
  if (data.notes !== undefined) updateData.notes = data.notes;

  await prisma.accord.update({ where: { id: accordId }, data: updateData });

  // Update tags if provided
  if (data.tags !== undefined) {
    await prisma.accordTag.deleteMany({ where: { accordId } });
    if (data.tags.length > 0) {
      await prisma.accordTag.createMany({
        data: data.tags.map((tag) => ({ accordId, tag })),
      });
    }
  }

  return getAccord(accordId, userId);
}

export async function deleteAccord(accordId: string, userId: string) {
  const accord = await prisma.accord.findFirst({ where: { id: accordId, userId } });
  if (!accord) throw new NotFoundError('Accord not found');

  // Check if used in recipes
  const usedInRecipes = await prisma.recipeIngredient.findFirst({
    where: { accordId },
    include: {
      version: {
        include: { recipe: true },
      },
    },
  });

  if (usedInRecipes) {
    // Get recipe names
    const ingredients = await prisma.recipeIngredient.findMany({
      where: { accordId },
      include: { version: { include: { recipe: { select: { name: true } } } } },
      take: 6,
    });
    const names = [...new Set(ingredients.map((i) => i.version.recipe.name))];
    const nameList = names.length > 5
      ? `${names.slice(0, 5).join(', ')}, and ${names.length - 5} more`
      : names.join(', ');
    throw new ValidationError(`cannot delete accord: used in recipes (${nameList})`);
  }

  await prisma.accord.delete({ where: { id: accordId } });
}

export async function addTag(accordId: string, userId: string, tag: string) {
  const accord = await prisma.accord.findFirst({ where: { id: accordId, userId } });
  if (!accord) throw new NotFoundError('Accord not found');

  // Upsert-like: ignore if exists
  try {
    await prisma.accordTag.create({ data: { accordId, tag } });
  } catch (err: any) {
    if (err.code !== 'P2002') throw err;
    // Already exists, ignore
  }
}

export async function removeTag(accordId: string, userId: string, tag: string) {
  const accord = await prisma.accord.findFirst({ where: { id: accordId, userId } });
  if (!accord) throw new NotFoundError('Accord not found');

  await prisma.accordTag.deleteMany({ where: { accordId, tag } });
}
