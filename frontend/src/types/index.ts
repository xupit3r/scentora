// Accord types
export interface Accord {
  _id: string;
  userId: string;
  name: string;
  pyramidPosition: 'top' | 'middle' | 'base';
  volumeMl: number;
  volumeDrops: number;
  supplier?: string;
  purchaseDate?: string;
  dilutionPercentage?: number;
  notes?: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateAccordRequest {
  name: string;
  pyramidPosition: 'top' | 'middle' | 'base';
  volumeMl: number;
  supplier?: string;
  purchaseDate?: string;
  dilutionPercentage?: number;
  notes?: string;
  tags?: string[];
}

export interface UpdateAccordRequest {
  name?: string;
  pyramidPosition?: 'top' | 'middle' | 'base';
  volumeMl?: number;
  supplier?: string;
  purchaseDate?: string;
  dilutionPercentage?: number;
  notes?: string;
  tags?: string[];
}

export interface PredefinedTag {
  _id: string;
  category: string;
  tag: string;
  createdAt: string;
}

export interface AccordFilters {
  position?: 'top' | 'middle' | 'base';
  minVolume?: number;
  maxVolume?: number;
  supplier?: string;
  search?: string;
  tags?: string[];
}

// Recipe types
export * from './recipe';
