import axios from 'axios';
import type { Perfume, JournalEntry } from '@/types';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const perfumeService = {
  async getAll(): Promise<Perfume[]> {
    const { data } = await api.get('/perfumes');
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

export default api;
