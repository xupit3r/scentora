import axios from 'axios';
import type { Perfume, JournalEntry, Accord, CreateAccordRequest, UpdateAccordRequest, AccordFilters, PredefinedTag } from '@/types';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Accord API Service
export const accordService = {
  async getAll(filters?: AccordFilters): Promise<Accord[]> {
    const params: Record<string, any> = {};
    if (filters?.position) params.position = filters.position;
    if (filters?.minVolume !== undefined) params.minVolume = filters.minVolume;
    if (filters?.maxVolume !== undefined) params.maxVolume = filters.maxVolume;
    if (filters?.supplier) params.supplier = filters.supplier;
    if (filters?.search) params.search = filters.search;
    if (filters?.tags?.length) params.tags = filters.tags;
    
    const { data } = await api.get('/accords', { params });
    return data.accords || [];
  },

  async getById(id: string): Promise<Accord> {
    const { data } = await api.get(`/accords/${id}`);
    return data.accord;
  },

  async create(accord: CreateAccordRequest): Promise<Accord> {
    const { data } = await api.post('/accords', accord);
    return data.accord;
  },

  async update(id: string, accord: UpdateAccordRequest): Promise<Accord> {
    const { data } = await api.put(`/accords/${id}`, accord);
    return data.accord;
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/accords/${id}`);
  },

  async addTag(id: string, tag: string): Promise<void> {
    await api.post(`/accords/${id}/tags`, { tag });
  },

  async removeTag(id: string, tag: string): Promise<void> {
    await api.delete(`/accords/${id}/tags/${tag}`);
  },
};

// Tag API Service
export const tagService = {
  async getAll(): Promise<PredefinedTag[]> {
    const { data } = await api.get('/tags');
    return data.tags || [];
  },

  async search(query: string): Promise<PredefinedTag[]> {
    const { data } = await api.get('/tags/search', { params: { q: query } });
    return data.tags || [];
  },

  async getCategories(): Promise<string[]> {
    const { data } = await api.get('/tags/categories');
    return data.categories || [];
  },

  async getGrouped(): Promise<Record<string, string[]>> {
    const { data } = await api.get('/tags/grouped');
    return data.tags || {};
  },

  async getByCategory(category: string): Promise<PredefinedTag[]> {
    const { data } = await api.get(`/tags/category/${category}`);
    return data.tags || [];
  },
};

// Legacy services (to be removed later)
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
