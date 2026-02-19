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
          :model-value="localFilters.tags || []"
          @update:model-value="(tags) => { localFilters.tags = tags; emitFilters(); }"
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

const localFilters = ref<AccordFilters>({ 
  ...props.modelValue,
  tags: props.modelValue.tags || []
});
const showLowStock = ref(false);

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { 
    ...newVal,
    tags: newVal.tags || []
  };
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
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid #E9E9E7;
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
  transition: transform 300ms ease;
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
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #E9E9E7;
}

.panel-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #37352F;
  margin: 0;
}

.close-btn {
  padding: 8px;
  background: none;
  border: none;
  cursor: pointer;
  color: #787774;
  transition: color 200ms;
}

.close-btn:hover {
  color: #37352F;
}

.close-btn svg {
  width: 20px;
  height: 20px;
}

.filters-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-label {
  font-size: 14px;
  font-weight: 600;
  color: #37352F;
}

.filter-input {
  padding: 10px 12px;
  border: 1px solid #E9E9E7;
  border-radius: 8px;
  font-size: 14px;
  transition: all 200ms;
  color: #37352F;
}

.filter-input:focus {
  outline: none;
  border-color: #0F766E;
  box-shadow: 0 0 0 3px rgba(15, 118, 110, 0.1);
}

.filter-input::placeholder {
  color: #9B9A97;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #787774;
  transition: color 200ms;
}

.radio-label:hover {
  color: #37352F;
}

.radio-label input[type="radio"] {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: #0F766E;
}

.position-badge {
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.badge-top {
  background: #FEF3C7;
  color: #92400E;
}

.badge-middle {
  background: #E9D5FF;
  color: #6B21A8;
}

.badge-base {
  background: #F5E6D3;
  color: #78350F;
}

.range-inputs {
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 8px;
}

.range-inputs .filter-input {
  flex: 1;
}

.range-separator {
  color: #9B9A97;
  font-size: 14px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #787774;
  transition: color 200ms;
}

.checkbox-label:hover {
  color: #37352F;
}

.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: #0F766E;
}

.filter-actions {
  padding-top: 12px;
  border-top: 1px solid #E9E9E7;
}

.btn-clear {
  width: 100%;
  padding: 10px 16px;
  background: white;
  color: #787774;
  border: 1px solid #E9E9E7;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 200ms;
}

.btn-clear:hover {
  background: #FAFAFA;
  color: #37352F;
  border-color: #D9D9D7;
}
</style>
