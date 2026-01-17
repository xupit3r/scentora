import axios from 'axios';
import type { Perfume, JournalEntry } from '@/types';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface PerfumeFilters {
  search?: string;
  concentration?: string;
  year?: number;
  note?: string;
}

export const perfumeService = {
  async getAll(filters?: PerfumeFilters): Promise<Perfume[]> {
    const { data } = await api.get('/perfumes', { params: filters });
    return data;
  },

  async getById(id: string): Promise<Perfume> {
    const { data } = await api.get(`/perfumes/${id}`);
    return data;
  },

  async create(perfume: Omit<Perfume, '_id' | '_rev' | 'createdAt' | 'updatedAt'>): Promise<Perfume> {
    const { data } = await api.post('/perfumes', perfume);
    return data;
  },

  async update(id: string, perfume: Partial<Perfume>): Promise<Perfume> {
    const { data } = await api.put(`/perfumes/${id}`, perfume);
    return data;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/perfumes/${id}`);
  },
};

export const journalService = {
  async getByPerfumeId(perfumeId: string): Promise<JournalEntry[]> {
    const { data } = await api.get(`/perfumes/${perfumeId}/journal`);
    return data;
  },

  async create(entry: Omit<JournalEntry, '_id' | '_rev' | 'createdAt' | 'updatedAt'>): Promise<JournalEntry> {
    const { data } = await api.post(`/perfumes/${entry.perfumeId}/journal`, entry);
    return data;
  },

  async update(id: string, entry: Partial<JournalEntry>): Promise<JournalEntry> {
    const { data } = await api.put(`/journal/${id}`, entry);
    return data;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/journal/${id}`);
  },
};

export const notesService = {
  async getAll(): Promise<string[]> {
    const { data } = await api.get('/notes');
    return data.notes;
  },
};

export interface CollectionStats {
  overview: {
    totalPerfumes: number;
    totalJournalEntries: number;
    averageRating: string;
    uniqueNotes: number;
    uniqueDesigners: number;
  };
  topDesigners: Array<{ designer: string; count: number }>;
  topNotes: Array<{ note: string; count: number }>;
  concentrationDistribution: Record<string, number>;
  yearDistribution: Array<{ year: number; count: number }>;
  pyramidStats: {
    topNotes: number;
    middleNotes: number;
    baseNotes: number;
    avgNotesPerPerfume: string;
  };
}

export const statsService = {
  async getCollectionStats(): Promise<CollectionStats> {
    const { data } = await api.get('/stats');
    return data;
  },
};

export const exportService = {
  async exportCollection(): Promise<Blob> {
    const { data } = await api.get('/export/collection', {
      responseType: 'blob',
    });
    return data;
  },

  async importCollection(file: File): Promise<any> {
    const text = await file.text();
    const data = JSON.parse(text);
    const response = await api.post('/export/import', data);
    return response.data;
  },
};

export default api;
