<template>
  <div class="home">
    <div v-if="!showForm && !editingPerfume" class="header-section">
      <h2>My Collection</h2>
      <button class="btn-primary" @click="showForm = true">+ Add Perfume</button>
    </div>
    
    <PerfumeForm
      v-if="showForm && !editingPerfume"
      @submit="handleCreatePerfume"
      @cancel="showForm = false"
    />

    <PerfumeForm
      v-if="editingPerfume"
      :perfume="editingPerfume"
      @submit="handleUpdatePerfume"
      @cancel="cancelEdit"
    />
    
    <div v-if="!showForm && !editingPerfume">
      <SearchFilters v-model="filters" />

      <div v-if="loading" class="loading">Loading your collection...</div>
      
      <div v-else-if="error" class="error">{{ error }}</div>
      
      <div v-else-if="filteredPerfumes.length === 0 && !hasActiveFilters" class="empty-state">
        <p>📦 Your collection is empty</p>
        <p>Start by adding your first perfume!</p>
      </div>

      <div v-else-if="filteredPerfumes.length === 0 && hasActiveFilters" class="empty-state">
        <p>🔍 No perfumes match your filters</p>
        <p>Try adjusting your search criteria</p>
      </div>
      
      <div v-else>
        <div class="results-count">
          Showing {{ filteredPerfumes.length }} {{ filteredPerfumes.length === 1 ? 'perfume' : 'perfumes' }}
        </div>
        <div class="perfume-grid">
          <div
            v-for="perfume in filteredPerfumes"
            :key="perfume._id"
            class="perfume-card"
          >
            <div class="card-actions">
              <button class="action-btn" @click.stop="startEdit(perfume)" title="Edit">
                ✏️
              </button>
              <button class="action-btn" @click.stop="viewPerfume(perfume._id!)" title="View">
                👁️
              </button>
            </div>
            <div @click="viewPerfume(perfume._id!)">
              <div class="perfume-image">
                <img v-if="perfume.imageUrl" :src="perfume.imageUrl" :alt="perfume.name" />
                <div v-else class="placeholder">🌸</div>
              </div>
              <div class="perfume-info">
                <h3>{{ perfume.name }}</h3>
                <p class="designer">{{ perfume.designer }}</p>
                <p class="year" v-if="perfume.year">{{ perfume.year }}</p>
                <div class="notes-preview">
                  <span v-if="perfume.pyramid.top.length" class="note-count">
                    {{ perfume.pyramid.top.length + perfume.pyramid.middle.length + perfume.pyramid.base.length }} notes
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { perfumeService, type PerfumeFilters } from '@/services/api';
import PerfumeForm from '@/components/PerfumeForm.vue';
import SearchFilters from '@/components/SearchFilters.vue';
import type { Perfume } from '@/types';

const router = useRouter();
const perfumes = ref<Perfume[]>([]);
const filteredPerfumes = ref<Perfume[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const showForm = ref(false);
const editingPerfume = ref<Perfume | null>(null);
const filters = ref<PerfumeFilters>({});

const hasActiveFilters = computed(() => {
  return !!(filters.value.search || filters.value.concentration || filters.value.year || filters.value.note);
});

watch(filters, async () => {
  await loadPerfumes();
}, { deep: true });

async function loadPerfumes() {
  try {
    loading.value = true;
    const result = await perfumeService.getAll(filters.value);
    perfumes.value = result;
    filteredPerfumes.value = result;
  } catch (err: any) {
    error.value = err.message || 'Failed to load perfumes';
  } finally {
    loading.value = false;
  }
}

async function handleCreatePerfume(perfumeData: Partial<Perfume>) {
  try {
    await perfumeService.create(perfumeData as any);
    showForm.value = false;
    await loadPerfumes();
  } catch (err: any) {
    error.value = err.message || 'Failed to create perfume';
  }
}

async function handleUpdatePerfume(perfumeData: Partial<Perfume>) {
  if (!editingPerfume.value?._id) return;

  try {
    await perfumeService.update(editingPerfume.value._id, perfumeData);
    editingPerfume.value = null;
    await loadPerfumes();
  } catch (err: any) {
    error.value = err.message || 'Failed to update perfume';
  }
}

function startEdit(perfume: Perfume) {
  editingPerfume.value = perfume;
  showForm.value = false;
}

function cancelEdit() {
  editingPerfume.value = null;
}

function viewPerfume(id: string) {
  router.push(`/perfume/${id}`);
}

// Initial load
loadPerfumes();
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
  position: relative;
}

.perfume-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.card-actions {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  display: flex;
  gap: 0.5rem;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.2s;
}

.perfume-card:hover .card-actions {
  opacity: 1;
}

.action-btn {
  background: white;
  border: none;
  border-radius: 50%;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  transition: transform 0.2s;
  font-size: 1rem;
}

.action-btn:hover {
  transform: scale(1.1);
}

.results-count {
  margin-bottom: 1rem;
  color: #666;
  font-size: 0.9rem;
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

.notes-preview {
  margin-top: 0.5rem;
}

.note-count {
  font-size: 0.85rem;
  color: #6b4f9e;
  font-weight: 500;
}
</style>
