<template>
  <div class="home">
    <div class="header-section">
      <h2>My Collection</h2>
      <button class="btn-primary">+ Add Perfume</button>
    </div>
    
    <div v-if="loading" class="loading">Loading your collection...</div>
    
    <div v-else-if="error" class="error">{{ error }}</div>
    
    <div v-else-if="perfumes.length === 0" class="empty-state">
      <p>📦 Your collection is empty</p>
      <p>Start by adding your first perfume!</p>
    </div>
    
    <div v-else class="perfume-grid">
      <div v-for="perfume in perfumes" :key="perfume._id" class="perfume-card">
        <div class="perfume-image">
          <img v-if="perfume.imageUrl" :src="perfume.imageUrl" :alt="perfume.name" />
          <div v-else class="placeholder">🌸</div>
        </div>
        <div class="perfume-info">
          <h3>{{ perfume.name }}</h3>
          <p class="designer">{{ perfume.designer }}</p>
          <p class="year" v-if="perfume.year">{{ perfume.year }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { perfumeService } from '@/services/api';
import type { Perfume } from '@/types';

const perfumes = ref<Perfume[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

onMounted(async () => {
  try {
    perfumes.value = await perfumeService.getAll();
  } catch (err: any) {
    error.value = err.message || 'Failed to load perfumes';
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.home {
  width: 100%;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

h2 {
  font-size: 2rem;
  color: #333;
}

.btn-primary {
  background: #6b4f9e;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #5a3f8d;
}

.loading,
.error,
.empty-state {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.error {
  color: #d32f2f;
}

.empty-state p:first-child {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.perfume-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.perfume-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.perfume-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.perfume-image {
  width: 100%;
  height: 200px;
  background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.perfume-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder {
  font-size: 4rem;
}

.perfume-info {
  padding: 1rem;
}

.perfume-info h3 {
  font-size: 1.1rem;
  margin-bottom: 0.5rem;
  color: #333;
}

.designer {
  color: #666;
  font-weight: 500;
}

.year {
  color: #999;
  font-size: 0.9rem;
  margin-top: 0.25rem;
}
</style>
