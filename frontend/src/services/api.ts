import axios from 'axios';
import type { Accord, CreateAccordRequest, UpdateAccordRequest, AccordFilters, PredefinedTag } from '@/types';

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

// Statistics API Service
export const statsService = {
  async getStats(): Promise<any> {
    const { data } = await api.get('/stats');
    return data;
  },
};

// Export/Import API Service
export const exportService = {
  async exportCollection(): Promise<Blob> {
    const { data } = await api.get('/export', {
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

// Re-export recipe service
export { recipeService } from './recipe';

export default api;

