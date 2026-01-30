export interface PerfumePyramid {
  top: string[];
  middle: string[];
  base: string[];
}

export interface Perfume {
  _id?: string;
  _rev?: string;
  type: 'perfume';
  userId: string;
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
  userId: string;
  perfumeId: string;
  date: string;
  content: string;
  rating?: number;
  occasion?: string;
  weather?: string;
  createdAt: string;
  updatedAt: string;
}

export interface User {
  _id?: string;
  _rev?: string;
  type: 'user';
  email: string;
  username: string;
  password: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthUser {
  id: string;
  email: string;
  username: string;
}

export interface RefreshToken {
  _id?: string;
  _rev?: string;
  type: 'refresh_token';
  userId: string;
  token: string;
  expiresAt: string;
  createdAt: string;
  revoked: boolean;
}

export interface Invitation {
  _id?: string;
  _rev?: string;
  type: 'invitation';
  code: string;
  email?: string;
  createdBy: string;
  expiresAt: string;
  used: boolean;
  usedAt?: string;
  usedBy?: string;
  createdAt: string;
}
