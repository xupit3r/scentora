<template>
  <div class="filters-panel" :class="{ mobile: isMobile, open: isOpen }">
    <div v-if="isMobile" class="panel-header">
      <h3>Filters</h3>
      <button @click="$emit('close')" class="close-btn">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <div class="filters-content">
      <div class="filter-group">
        <label class="filter-label">Search</label>
        <input
          v-model="localFilters.search"
          @input="emitFilters"
          type="text"
          placeholder="Search by name..."
          class="filter-input"
        />
      </div>

      <div class="filter-group">
        <label class="filter-label">Pyramid Position</label>
        <div class="radio-group">
          <label class="radio-label">
            <input
              type="radio"
              :value="undefined"
              v-model="localFilters.position"
              @change="emitFilters"
            />
            <span>All</span>
          </label>
          <label class="radio-label">
            <input
              type="radio"
              value="top"
              v-model="localFilters.position"
              @change="emitFilters"
            />
            <span class="position-badge badge-top">Top</span>
          </label>
          <label class="radio-label">
            <input
              type="radio"
              value="middle"
              v-model="localFilters.position"
              @change="emitFilters"
            />
            <span class="position-badge badge-middle">Middle</span>
          </label>
          <label class="radio-label">
            <input
              type="radio"
              value="base"
              v-model="localFilters.position"
              @change="emitFilters"
            />
            <span class="position-badge badge-base">Base</span>
          </label>
        </div>
      </div>

      <div class="filter-group">
        <label class="filter-label">Volume Range</label>
        <div class="range-inputs">
          <input
            v-model.number="localFilters.minVolume"
            @input="emitFilters"
            type="number"
            step="0.1"
            min="0"
            placeholder="Min (ml)"
            class="filter-input"
          />
          <span class="range-separator">to</span>
          <input
            v-model.number="localFilters.maxVolume"
            @input="emitFilters"
            type="number"
            step="0.1"
            min="0"
            placeholder="Max (ml)"
            class="filter-input"
          />
        </div>
      </div>

      <div class="filter-group">
        <label class="filter-label">Supplier</label>
        <input
          v-model="localFilters.supplier"
          @input="emitFilters"
          type="text"
          placeholder="Filter by supplier..."
          class="filter-input"
        />
      </div>

      <div class="filter-group">
        <label class="filter-label">Tags</label>
        <TagSelector
          v-model="localFilters.tags"
          @update:modelValue="emitFilters"
          placeholder="Filter by tags..."
          :show-suggestions="false"
        />
      </div>

      <div class="filter-group">
        <label class="checkbox-label">
          <input
            type="checkbox"
            v-model="showLowStock"
            @change="handleLowStockToggle"
          />
          <span>Show only low stock (< 5ml)</span>
        </label>
      </div>

      <div class="filter-actions">
        <button @click="clearFilters" class="btn-clear">
          Clear All
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import TagSelector from './TagSelector.vue';
import type { AccordFilters } from '@/types';

const props = withDefaults(defineProps<{
  modelValue: AccordFilters;
  isMobile?: boolean;
  isOpen?: boolean;
}>(), {
  isMobile: false,
  isOpen: false,
});

const emit = defineEmits<{
  'update:modelValue': [filters: AccordFilters];
  close: [];
}>();

const localFilters = ref<AccordFilters>({ ...props.modelValue });
const showLowStock = ref(false);

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...newVal };
  // Check if low stock filter is active
  showLowStock.value = newVal.maxVolume === 5;
}, { deep: true });

function emitFilters() {
  // Remove empty values
  const filters: AccordFilters = {};
  
  if (localFilters.value.position) filters.position = localFilters.value.position;
  if (localFilters.value.minVolume !== undefined && localFilters.value.minVolume !== null) {
    filters.minVolume = localFilters.value.minVolume;
  }
  if (localFilters.value.maxVolume !== undefined && localFilters.value.maxVolume !== null) {
    filters.maxVolume = localFilters.value.maxVolume;
  }
  if (localFilters.value.supplier) filters.supplier = localFilters.value.supplier;
  if (localFilters.value.search) filters.search = localFilters.value.search;
  if (localFilters.value.tags && localFilters.value.tags.length > 0) {
    filters.tags = localFilters.value.tags;
  }

  emit('update:modelValue', filters);
}

function handleLowStockToggle() {
  if (showLowStock.value) {
    localFilters.value.maxVolume = 5;
  } else {
    localFilters.value.maxVolume = undefined;
  }
  emitFilters();
}

function clearFilters() {
  localFilters.value = {};
  showLowStock.value = false;
  emitFilters();
}
</script>

<style scoped>
.filters-panel {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.filters-panel.mobile {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 85%;
  max-width: 400px;
  border-radius: 0;
  transform: translateX(100%);
  transition: transform 0.3s ease;
  z-index: 1000;
  overflow-y: auto;
}

.filters-panel.mobile.open {
  transform: translateX(0);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #E5E7EB;
}

.panel-header h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1F2937;
  margin: 0;
}

.close-btn {
  padding: 0.5rem;
  background: none;
  border: none;
  cursor: pointer;
  color: #6B7280;
}

.close-btn:hover {
  color: #1F2937;
}

.close-btn svg {
  width: 1.5rem;
  height: 1.5rem;
}

.filters-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.filter-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: #374151;
}

.filter-input {
  padding: 0.625rem 0.875rem;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-size: 0.875rem;
  transition: all 0.2s;
}

.filter-input:focus {
  outline: none;
  border-color: #14B8A6;
  box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.1);
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  color: #4B5563;
}

.radio-label input[type="radio"] {
  width: 1rem;
  height: 1rem;
  cursor: pointer;
}

.position-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.badge-top {
  background: linear-gradient(135deg, #FFD93D, #FFA800);
  color: #4a2800;
}

.badge-middle {
  background: linear-gradient(135deg, #B565D8, #8B5CF6);
  color: white;
}

.badge-base {
  background: linear-gradient(135deg, #A0826D, #6B4423);
  color: white;
}

.range-inputs {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.range-inputs .filter-input {
  flex: 1;
}

.range-separator {
  color: #6B7280;
  font-size: 0.875rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  color: #4B5563;
}

.checkbox-label input[type="checkbox"] {
  width: 1rem;
  height: 1rem;
  cursor: pointer;
}

.filter-actions {
  padding-top: 0.5rem;
  border-top: 1px solid #E5E7EB;
}

.btn-clear {
  width: 100%;
  padding: 0.625rem 1rem;
  background: white;
  color: #6B7280;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-clear:hover {
  background: #F9FAFB;
  color: #1F2937;
}
</style>
