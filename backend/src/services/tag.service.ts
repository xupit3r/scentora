import { prisma } from '../prisma.js';

function tagToApi(tag: { id: string; category: string; tag: string; createdAt: Date }) {
  return {
    _id: tag.id,
    category: tag.category,
    tag: tag.tag,
    createdAt: tag.createdAt,
  };
}

export async function getAllTags() {
  const tags = await prisma.predefinedTag.findMany({ orderBy: [{ category: 'asc' }, { tag: 'asc' }] });
  return tags.map(tagToApi);
}

export async function searchTags(query: string) {
  if (!query) return [];
  const tags = await prisma.predefinedTag.findMany({
    where: { tag: { contains: query, mode: 'insensitive' } },
    orderBy: { tag: 'asc' },
  });
  return tags.map(tagToApi);
}

export async function getCategories() {
  const result = await prisma.predefinedTag.findMany({
    select: { category: true },
    distinct: ['category'],
    orderBy: { category: 'asc' },
  });
  return result.map((r) => r.category);
}

export async function getTagsGrouped() {
  const tags = await prisma.predefinedTag.findMany({ orderBy: [{ category: 'asc' }, { tag: 'asc' }] });
  const grouped: Record<string, string[]> = {};
  for (const tag of tags) {
    if (!grouped[tag.category]) grouped[tag.category] = [];
    grouped[tag.category].push(tag.tag);
  }
  return grouped;
}

export async function getByCategory(category: string) {
  const tags = await prisma.predefinedTag.findMany({
    where: { category },
    orderBy: { tag: 'asc' },
  });
  return tags.map(tagToApi);
}
