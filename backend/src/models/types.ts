export interface PerfumePyramid {
  top: string[];
  middle: string[];
  base: string[];
}

export interface Perfume {
  _id?: string;
  _rev?: string;
  type: 'perfume';
  name: string;
  designer: string;
  year?: number;
  concentration?: string;
  pyramid: PerfumePyramid;
  description?: string;
  imageUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface JournalEntry {
  _id?: string;
  _rev?: string;
  type: 'journal';
  perfumeId: string;
  date: string;
  content: string;
  rating?: number;
  occasion?: string;
  weather?: string;
  createdAt: string;
  updatedAt: string;
}
