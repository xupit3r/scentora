import { prisma } from '../prisma.js';
import { NotFoundError } from '../utils/errors.js';

function noteToApi(note: any) {
  return {
    _id: note.id,
    recipeId: note.recipeId,
    ...(note.versionId != null ? { versionId: note.versionId } : {}),
    content: note.content,
    noteType: note.noteType,
    createdAt: note.createdAt,
    updatedAt: note.updatedAt,
  };
}

export async function createNote(recipeId: string, userId: string, data: {
  content: string;
  versionId?: string;
  noteType?: string;
}) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  const note = await prisma.recipeNote.create({
    data: {
      recipeId,
      versionId: data.versionId ?? null,
      content: data.content,
      noteType: data.noteType || 'general',
    },
  });

  return noteToApi(note);
}

export async function listNotes(recipeId: string, userId: string) {
  const recipe = await prisma.recipe.findFirst({ where: { id: recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  const notes = await prisma.recipeNote.findMany({
    where: { recipeId },
    orderBy: { createdAt: 'desc' },
  });

  return notes.map(noteToApi);
}

export async function updateNote(noteId: string, userId: string, data: { content: string }) {
  const note = await prisma.recipeNote.findUnique({ where: { id: noteId } });
  if (!note) throw new NotFoundError('note not found');

  const recipe = await prisma.recipe.findFirst({ where: { id: note.recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  const updated = await prisma.recipeNote.update({
    where: { id: noteId },
    data: { content: data.content },
  });

  return noteToApi(updated);
}

export async function deleteNote(noteId: string, userId: string) {
  const note = await prisma.recipeNote.findUnique({ where: { id: noteId } });
  if (!note) throw new NotFoundError('note not found');

  const recipe = await prisma.recipe.findFirst({ where: { id: note.recipeId, userId } });
  if (!recipe) throw new NotFoundError('recipe not found');

  await prisma.recipeNote.delete({ where: { id: noteId } });
}
