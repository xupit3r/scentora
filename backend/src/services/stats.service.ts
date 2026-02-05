import { prisma } from '../prisma.js';

export async function getStatistics(userId: string) {
  const accords = await prisma.accord.findMany({
    where: { userId },
    include: { tags: true },
  });

  // Overview
  const totalVolume = accords.reduce((sum, a) => sum + Number(a.volumeMl), 0);
  const suppliers = new Set<string>();
  const tagSet = new Set<string>();

  for (const accord of accords) {
    if (accord.supplier) suppliers.add(accord.supplier);
    for (const t of accord.tags) tagSet.add(t.tag);
  }

  const overview = {
    totalAccords: accords.length,
    totalVolume,
    totalSuppliers: suppliers.size,
    totalTags: tagSet.size,
  };

  // Pyramid stats
  const pyramidStats = { topCount: 0, middleCount: 0, baseCount: 0, topVolume: 0, middleVolume: 0, baseVolume: 0 };
  for (const accord of accords) {
    const vol = Number(accord.volumeMl);
    switch (accord.pyramidPosition) {
      case 'top': pyramidStats.topCount++; pyramidStats.topVolume += vol; break;
      case 'middle': pyramidStats.middleCount++; pyramidStats.middleVolume += vol; break;
      case 'base': pyramidStats.baseCount++; pyramidStats.baseVolume += vol; break;
    }
  }

  // Tag stats
  const tagCounts: Record<string, number> = {};
  for (const accord of accords) {
    for (const t of accord.tags) {
      tagCounts[t.tag] = (tagCounts[t.tag] || 0) + 1;
    }
  }
  const tagStats = Object.entries(tagCounts).map(([tag, count]) => ({ tag, count }));

  // Supplier stats
  const supplierCounts: Record<string, number> = {};
  for (const accord of accords) {
    if (accord.supplier) supplierCounts[accord.supplier] = (supplierCounts[accord.supplier] || 0) + 1;
  }
  const supplierStats = Object.entries(supplierCounts).map(([supplier, count]) => ({ supplier, count }));

  // Volume stats
  let volumeStats = { averageVolume: 0, minVolume: 0, maxVolume: 0 };
  if (accords.length > 0) {
    const volumes = accords.map((a) => Number(a.volumeMl));
    volumeStats = {
      averageVolume: totalVolume / accords.length,
      minVolume: Math.min(...volumes),
      maxVolume: Math.max(...volumes),
    };
  }

  // Low inventory (< 10ml)
  const lowInventory = accords
    .filter((a) => Number(a.volumeMl) < 10)
    .map((a) => ({
      accordId: a.id,
      name: a.name,
      volumeMl: Number(a.volumeMl),
      ...(a.supplier != null ? { supplier: a.supplier } : {}),
    }));

  return { overview, pyramidStats, tagStats, supplierStats, volumeStats, lowInventory };
}
