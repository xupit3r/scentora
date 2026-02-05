import { PrismaClient } from '@prisma/client';
import { afterAll, beforeEach } from 'vitest';

export const testPrisma = new PrismaClient();

export async function cleanDatabase() {
  // Delete in FK-safe order
  await testPrisma.recipeCollectionMember.deleteMany();
  await testPrisma.recipeCollection.deleteMany();
  await testPrisma.recipeTag.deleteMany();
  await testPrisma.recipeNote.deleteMany();
  await testPrisma.recipeIngredient.deleteMany();
  await testPrisma.recipeVersion.deleteMany();
  await testPrisma.recipe.deleteMany();
  await testPrisma.accordTag.deleteMany();
  await testPrisma.accord.deleteMany();
  await testPrisma.refreshToken.deleteMany();
  await testPrisma.invitation.deleteMany();
  await testPrisma.user.deleteMany();
  await testPrisma.predefinedTag.deleteMany();
}

beforeEach(async () => {
  await cleanDatabase();
});

afterAll(async () => {
  await cleanDatabase();
  await testPrisma.$disconnect();
});
