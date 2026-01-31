<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>{{ isEdit ? 'Edit Accord' : 'New Accord' }}</h2>
        <button @click="$emit('close')" class="close-btn">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-section">
          <h3 class="section-title">Basic Information</h3>
          
          <div class="form-group">
            <label for="name" class="required">Name</label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              class="form-input"
              placeholder="e.g., Bergamot Essential Oil"
            />
          </div>

          <div class="form-group">
            <label for="position" class="required">Pyramid Position</label>
            <select
              id="position"
              v-model="form.pyramidPosition"
              required
              class="form-select"
            >
              <option value="">Select position...</option>
              <option value="top">Top Note</option>
              <option value="middle">Middle Note</option>
              <option value="base">Base Note</option>
            </select>
          </div>
        </div>

        <div class="form-section">
          <h3 class="section-title">Inventory</h3>
          
          <div class="form-group">
            <label for="volume" class="required">Volume (ml)</label>
            <input
              id="volume"
              v-model.number="form.volumeMl"
              type="number"
              step="0.01"
              min="0"
              required
              class="form-input"
              placeholder="25.0"
            />
            <span class="helper-text">
              ≈ {{ estimatedDrops }} drops
            </span>
          </div>

          <div class="form-group">
            <label for="supplier">Supplier</label>
            <input
              id="supplier"
              v-model="form.supplier"
              type="text"
              class="form-input"
              placeholder="e.g., Perfumer's Apprentice"
            />
          </div>

          <div class="form-group">
            <label for="purchaseDate">Purchase Date</label>
            <input
              id="purchaseDate"
              v-model="form.purchaseDate"
              type="date"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="dilution">Dilution Percentage</label>
            <input
              id="dilution"
              v-model.number="form.dilutionPercentage"
              type="number"
              step="0.1"
              min="0"
              max="100"
              class="form-input"
              placeholder="10.0"
            />
            <span class="helper-text">
              Optional: concentration if diluted
            </span>
          </div>
        </div>

        <div class="form-section">
          <h3 class="section-title">Tags</h3>
          <TagSelector
            v-model="form.tags"
            placeholder="Add tags to categorize this accord..."
          />
        </div>

        <div class="form-section">
          <h3 class="section-title">Notes</h3>
          <div class="form-group">
            <label for="notes">Additional Notes</label>
            <textarea
              id="notes"
              v-model="form.notes"
              rows="4"
              class="form-textarea"
              placeholder="Any additional information, usage notes, or observations..."
            ></textarea>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" @click="$emit('close')" class="btn btn-secondary">
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? 'Saving...' : (isEdit ? 'Update' : 'Create') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import TagSelector from './TagSelector.vue';
import type { Accord, CreateAccordRequest, UpdateAccordRequest } from '@/types';

const props = defineProps<{
  accord?: Accord;
  loading?: boolean;
}>();

const emit = defineEmits<{
  close: [];
  submit: [data: CreateAccordRequest | UpdateAccordRequest];
}>();

const isEdit = computed(() => !!props.accord);

const form = ref<CreateAccordRequest>({
  name: '',
  pyramidPosition: 'top',
  volumeMl: 0,
  supplier: '',
  purchaseDate: '',
  dilutionPercentage: undefined,
  notes: '',
  tags: [],
});

// Initialize form with accord data if editing
watch(() => props.accord, (accord) => {
  if (accord) {
    form.value = {
      name: accord.name,
      pyramidPosition: accord.pyramidPosition,
      volumeMl: accord.volumeMl,
      supplier: accord.supplier || '',
      purchaseDate: accord.purchaseDate || '',
      dilutionPercentage: accord.dilutionPercentage,
      notes: accord.notes || '',
      tags: [...accord.tags],
    };
  }
}, { immediate: true });

const estimatedDrops = computed(() => {
  return Math.round((form.value.volumeMl || 0) * 20);
});

function handleSubmit() {
  const data: CreateAccordRequest | UpdateAccordRequest = {
    name: form.value.name,
    pyramidPosition: form.value.pyramidPosition,
    volumeMl: form.value.volumeMl,
    tags: form.value.tags,
  };

  if (form.value.supplier) data.supplier = form.value.supplier;
  if (form.value.purchaseDate) data.purchaseDate = form.value.purchaseDate;
  if (form.value.dilutionPercentage !== undefined) data.dilutionPercentage = form.value.dilutionPercentage;
  if (form.value.notes) data.notes = form.value.notes;

  emit('submit', data);
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #E5E7EB;
}

.modal-header h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1F2937;
  margin: 0;
}

.close-btn {
  padding: 0.5rem;
  background: none;
  border: none;
  cursor: pointer;
  color: #6B7280;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #1F2937;
}

.close-btn svg {
  width: 1.5rem;
  height: 1.5rem;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
}

.form-section {
  margin-bottom: 2rem;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1F2937;
  margin: 0 0 1rem 0;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #F3F4F6;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group:last-child {
  margin-bottom: 0;
}

label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.5rem;
}

label.required::after {
  content: ' *';
  color: #DC2626;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-size: 0.875rem;
  transition: all 0.2s;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: #14B8A6;
  box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.helper-text {
  display: block;
  margin-top: 0.375rem;
  font-size: 0.75rem;
  color: #6B7280;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.5rem;
  border-top: 1px solid #E5E7EB;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-secondary {
  background: white;
  color: #374151;
  border: 1px solid #D1D5DB;
}

.btn-secondary:hover {
  background: #F9FAFB;
}

.btn-primary {
  background: #14B8A6;
  color: white;
}

.btn-primary:hover {
  background: #0D9488;
}

.btn-primary:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
}
</style>
