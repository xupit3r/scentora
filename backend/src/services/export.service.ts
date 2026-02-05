import { Decimal } from '@prisma/client/runtime/library';
import { prisma } from '../prisma.js';

export async function exportAccords(userId: string) {
  const accords = await prisma.accord.findMany({
    where: { userId },
    include: { tags: true },
    orderBy: { createdAt: 'desc' },
  });

  const accordResponses = accords.map((a) => ({
    _id: a.id,
    userId: a.userId,
    name: a.name,
    pyramidPosition: a.pyramidPosition,
    volumeMl: Number(a.volumeMl),
    volumeDrops: a.volumeDrops ?? Math.round(Number(a.volumeMl) * 20),
    ...(a.supplier != null ? { supplier: a.supplier } : {}),
    ...(a.purchaseDate != null ? { purchaseDate: a.purchaseDate } : {}),
    ...(a.dilutionPercentage != null ? { dilutionPercentage: Number(a.dilutionPercentage) } : {}),
    ...(a.notes != null ? { notes: a.notes } : {}),
    tags: a.tags.map((t) => t.tag),
    createdAt: a.createdAt,
    updatedAt: a.updatedAt,
  }));

  return {
    version: '1.0',
    exportedAt: new Date(),
    accords: accordResponses,
  };
}

export async function importAccords(userId: string, data: { accords: any[] }) {
  let imported = 0;
  let failed = 0;
  const errors: string[] = [];

  for (const accord of data.accords) {
    try {
      await prisma.accord.create({
        data: {
          userId,
          name: accord.name,
          pyramidPosition: accord.pyramidPosition,
          volumeMl: new Decimal(accord.volumeMl),
          supplier: accord.supplier ?? null,
          purchaseDate: accord.purchaseDate ? new Date(accord.purchaseDate) : null,
          dilutionPercentage: accord.dilutionPercentage != null ? new Decimal(accord.dilutionPercentage) : null,
          notes: accord.notes ?? null,
          ...(accord.tags && accord.tags.length > 0
            ? { tags: { create: accord.tags.map((tag: string) => ({ tag })) } }
            : {}),
        },
      });
      imported++;
    } catch (err: any) {
      failed++;
      errors.push(err.message || String(err));
    }
  }

  return {
    totalRecords: data.accords.length,
    importedRecords: imported,
    failedRecords: failed,
    ...(errors.length > 0 ? { errors } : {}),
  };
}
