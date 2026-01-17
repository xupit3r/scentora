<template>
  <div class="search-filters">
    <div class="search-bar">
      <input
        v-model="localFilters.search"
        type="text"
        placeholder="Search perfumes by name, designer, or description..."
        @input="handleFilterChange"
      />
      <button v-if="localFilters.search" class="clear-btn" @click="clearSearch">
        ×
      </button>
    </div>

    <div class="filters-toggle">
      <button class="filter-btn" @click="showFilters = !showFilters">
        🔍 Filters {{ hasActiveFilters ? `(${activeFilterCount})` : '' }}
      </button>
      <button v-if="hasActiveFilters" class="clear-all-btn" @click="clearAllFilters">
        Clear All
      </button>
    </div>

    <div v-if="showFilters" class="filters-panel">
      <div class="filter-group">
        <label for="concentration">Concentration</label>
        <select
          id="concentration"
          v-model="localFilters.concentration"
          @change="handleFilterChange"
        >
          <option value="">All</option>
          <option value="Parfum">Parfum</option>
          <option value="EDP">Eau de Parfum (EDP)</option>
          <option value="EDT">Eau de Toilette (EDT)</option>
          <option value="EDC">Eau de Cologne (EDC)</option>
        </select>
      </div>

      <div class="filter-group">
        <label for="year">Year</label>
        <input
          id="year"
          v-model.number="localFilters.year"
          type="number"
          placeholder="e.g., 2010"
          @input="handleFilterChange"
        />
      </div>

      <div class="filter-group">
        <label for="note">Note</label>
        <input
          id="note"
          v-model="localFilters.note"
          type="text"
          placeholder="e.g., Bergamot"
          list="notes-list"
          @input="handleFilterChange"
        />
        <datalist id="notes-list">
          <option v-for="note in availableNotes" :key="note" :value="note" />
        </datalist>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { notesService, type PerfumeFilters } from '@/services/api';

interface Props {
  modelValue: PerfumeFilters;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  'update:modelValue': [filters: PerfumeFilters];
}>();

const localFilters = ref<PerfumeFilters>({ ...props.modelValue });
const showFilters = ref(false);
const availableNotes = ref<string[]>([]);

const hasActiveFilters = computed(() => {
  return !!(
    localFilters.value.search ||
    localFilters.value.concentration ||
    localFilters.value.year ||
    localFilters.value.note
  );
});

const activeFilterCount = computed(() => {
  let count = 0;
  if (localFilters.value.search) count++;
  if (localFilters.value.concentration) count++;
  if (localFilters.value.year) count++;
  if (localFilters.value.note) count++;
  return count;
});

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal };
}, { deep: true });

function handleFilterChange() {
  // Clean up empty values
  const filters: PerfumeFilters = {};
  if (localFilters.value.search?.trim()) filters.search = localFilters.value.search.trim();
  if (localFilters.value.concentration) filters.concentration = localFilters.value.concentration;
  if (localFilters.value.year) filters.year = localFilters.value.year;
  if (localFilters.value.note?.trim()) filters.note = localFilters.value.note.trim();
  
  emit('update:modelValue', filters);
}

function clearSearch() {
  localFilters.value.search = '';
  handleFilterChange();
}

function clearAllFilters() {
  localFilters.value = {};
  showFilters.value = false;
  emit('update:modelValue', {});
}

onMounted(async () => {
  try {
    availableNotes.value = await notesService.getAll();
  } catch (error) {
    console.error('Failed to load notes:', error);
  }
});
</script>

<style scoped>
.search-filters {
  margin-bottom: 2rem;
}

.search-bar {
  position: relative;
  margin-bottom: 1rem;
}

.search-bar input {
  width: 100%;
  padding: 1rem 3rem 1rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.search-bar input:focus {
  outline: none;
  border-color: #6b4f9e;
}

.clear-btn {
  position: absolute;
  right: 0.5rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #999;
  cursor: pointer;
  padding: 0.5rem;
  line-height: 1;
}

.clear-btn:hover {
  color: #333;
}

.filters-toggle {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.filter-btn,
.clear-all-btn {
  padding: 0.5rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  background: white;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-btn:hover {
  border-color: #6b4f9e;
  background: #f9f9f9;
}

.clear-all-btn {
  border-color: #d32f2f;
  color: #d32f2f;
}

.clear-all-btn:hover {
  background: #ffebee;
}

.filters-panel {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  padding: 1.5rem;
  background: white;
  border-radius: 8px;
  border: 2px solid #e0e0e0;
}

.filter-group {
  display: flex;
  flex-direction: column;
}

.filter-group label {
  margin-bottom: 0.5rem;
  font-weight: 600;
  font-size: 0.9rem;
  color: #555;
}

.filter-group input,
.filter-group select {
  padding: 0.5rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 0.9rem;
}

.filter-group input:focus,
.filter-group select:focus {
  outline: none;
  border-color: #6b4f9e;
}
</style>
