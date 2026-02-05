import { prisma } from '../prisma.js';
import { NotFoundError } from '../utils/errors.js';

export async function addTag(recipeId: string, userId: string, tag: string) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  try {
    await prisma.recipeTag.create({ data: { recipeId, tag } });
  } catch (err: any) {
    if (err.code !== 'P2002') throw err;
    // Already exists, no-op
  }
}

export async function removeTag(recipeId: string, userId: string, tag: string) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  await prisma.recipeTag.deleteMany({ where: { recipeId, tag } });
}

export async function getPopularTags(userId: string, limit = 20) {
  const result = await prisma.$queryRaw<{ tag: string; count: bigint }[]>`
    SELECT rt.tag, COUNT(*) as count
    FROM recipe_tags rt
    JOIN recipes r ON rt.recipe_id = r.id
    WHERE r.user_id = ${userId}::uuid
    GROUP BY rt.tag
    ORDER BY count DESC
    LIMIT ${limit}
  `;

  return result.map((r) => ({ tag: r.tag, count: Number(r.count) }));
}
