<template>
  <div class="home">
    <div class="top-bar">
      <div class="header-section">
        <h1>Accord Inventory</h1>
        <p class="subtitle">Manage your essential oils and perfume accords</p>
      </div>
      <div class="header-actions">
        <button @click="toggleFilters" class="btn-filter" :class="{ active: showFilters }">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
          </svg>
          <span>Filters</span>
        </button>
        <button @click="showForm = true" class="btn-primary">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span>New Accord</span>
        </button>
      </div>
    </div>

    <div class="main-content">
      <aside v-if="showFilters" class="filters-sidebar">
        <AccordFilters v-model="filters" />
      </aside>

      <div class="content-area">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Loading your accords...</p>
        </div>

        <div v-else-if="error" class="error-state">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p>{{ error }}</p>
          <button @click="loadAccords" class="btn-secondary">Try Again</button>
        </div>

        <div v-else-if="accords.length === 0 && !hasActiveFilters" class="empty-state">
          <div class="empty-icon">📦</div>
          <h3>Your inventory is empty</h3>
          <p>Start by adding your first essential oil or accord</p>
          <button @click="showForm = true" class="btn-primary">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            <span>Add First Accord</span>
          </button>
        </div>

        <div v-else-if="accords.length === 0 && hasActiveFilters" class="empty-state">
          <div class="empty-icon">🔍</div>
          <h3>No accords match your filters</h3>
          <p>Try adjusting your search criteria</p>
          <button @click="clearFilters" class="btn-secondary">Clear Filters</button>
        </div>

        <div v-else class="results-section">
          <div class="results-header">
            <div class="results-count">
              <strong>{{ accords.length }}</strong> {{ accords.length === 1 ? 'accord' : 'accords' }}
            </div>
          </div>

          <div class="accords-grid">
            <AccordCard
              v-for="accord in accords"
              :key="accord._id"
              :accord="accord"
              @view="viewAccord"
              @edit="editAccord"
              @delete="confirmDelete"
            />
          </div>
        </div>
      </div>
    </div>

    <AccordForm
      v-if="showForm"
      :accord="editingAccord"
      :loading="submitting"
      @close="closeForm"
      @submit="handleSubmit"
    />

    <div v-if="deletingAccord" class="modal-overlay" @click.self="deletingAccord = null">
      <div class="delete-modal">
        <div class="delete-icon">⚠️</div>
        <h3>Delete Accord?</h3>
        <p>Are you sure you want to delete <strong>{{ deletingAccord.name }}</strong>?</p>
        <p class="warning-text">This action cannot be undone.</p>
        <div class="modal-actions">
          <button @click="deletingAccord = null" class="btn-secondary">Cancel</button>
          <button @click="handleDelete" class="btn-danger" :disabled="submitting">
            {{ submitting ? 'Deleting...' : 'Delete' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import AccordCard from '@/components/AccordCard.vue';
import AccordForm from '@/components/AccordForm.vue';
import AccordFilters from '@/components/AccordFilters.vue';
import { accordService } from '@/services/api';
import type { Accord, AccordFilters as Filters, CreateAccordRequest, UpdateAccordRequest } from '@/types';

const router = useRouter();

const accords = ref<Accord[]>([]);
const filters = ref<Filters>({});
const loading = ref(false);
const submitting = ref(false);
const error = ref('');
const showForm = ref(false);
const showFilters = ref(true);
const editingAccord = ref<Accord | undefined>();
const deletingAccord = ref<Accord | null>(null);

const hasActiveFilters = computed(() => {
  return Object.keys(filters.value).length > 0;
});

onMounted(() => {
  loadAccords();
});

async function loadAccords() {
  loading.value = true;
  error.value = '';
  
  try {
    accords.value = await accordService.getAll(filters.value);
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'Failed to load accords';
    console.error('Load accords error:', err);
  } finally {
    loading.value = false;
  }
}

function toggleFilters() {
  showFilters.value = !showFilters.value;
}

function clearFilters() {
  filters.value = {};
  loadAccords();
}

function viewAccord(accord: Accord) {
  router.push(`/accords/${accord._id}`);
}

function editAccord(accord: Accord) {
  editingAccord.value = accord;
  showForm.value = true;
}

function closeForm() {
  showForm.value = false;
  editingAccord.value = undefined;
}

function confirmDelete(accord: Accord) {
  deletingAccord.value = accord;
}

async function handleSubmit(data: CreateAccordRequest | UpdateAccordRequest) {
  submitting.value = true;
  error.value = '';

  try {
    if (editingAccord.value) {
      await accordService.update(editingAccord.value._id, data as UpdateAccordRequest);
    } else {
      await accordService.create(data as CreateAccordRequest);
    }
    
    closeForm();
    await loadAccords();
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'Failed to save accord';
    console.error('Save accord error:', err);
  } finally {
    submitting.value = false;
  }
}

async function handleDelete() {
  if (!deletingAccord.value) return;

  submitting.value = true;
  error.value = '';

  try {
    await accordService.delete(deletingAccord.value._id);
    deletingAccord.value = null;
    await loadAccords();
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'Failed to delete accord';
    console.error('Delete accord error:', err);
  } finally {
    submitting.value = false;
  }
}

watch(filters, () => {
  loadAccords();
}, { deep: true });
</script>

<style scoped>
.home {
  min-height: 100vh;
  background: #F5F5F5;
}

.top-bar {
  background: white;
  border-bottom: 1px solid #E5E7EB;
  padding: 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
}

.header-section h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #1F2937;
  margin: 0 0 0.5rem 0;
}

.subtitle {
  color: #6B7280;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 1rem;
  flex-shrink: 0;
}

.btn-filter {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: white;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-filter svg {
  width: 1.25rem;
  height: 1.25rem;
}

.btn-filter:hover,
.btn-filter.active {
  background: #F3F4F6;
  border-color: #9CA3AF;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #14B8A6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary svg {
  width: 1.25rem;
  height: 1.25rem;
}

.btn-primary:hover {
  background: #0D9488;
}

.main-content {
  display: flex;
  gap: 2rem;
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.filters-sidebar {
  flex-shrink: 0;
  width: 280px;
}

.content-area {
  flex: 1;
  min-width: 0;
}

.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
}

.spinner {
  width: 3rem;
  height: 3rem;
  border: 4px solid #E5E7EB;
  border-top-color: #14B8A6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-state svg,
.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.error-state svg {
  width: 4rem;
  height: 4rem;
  color: #DC2626;
}

.empty-state h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1F2937;
  margin: 0 0 0.5rem 0;
}

.empty-state p {
  color: #6B7280;
  margin: 0 0 1.5rem 0;
}

.results-section {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #E5E7EB;
}

.results-count {
  font-size: 0.875rem;
  color: #6B7280;
}

.results-count strong {
  color: #1F2937;
  font-size: 1rem;
}

.accords-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
}

.btn-secondary {
  padding: 0.75rem 1.5rem;
  background: white;
  color: #374151;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #F9FAFB;
}

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

.delete-modal {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  max-width: 400px;
  text-align: center;
}

.delete-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.delete-modal h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1F2937;
  margin: 0 0 1rem 0;
}

.delete-modal p {
  color: #6B7280;
  margin: 0 0 0.5rem 0;
}

.warning-text {
  color: #DC2626;
  font-size: 0.875rem;
  margin-bottom: 1.5rem !important;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: center;
}

.btn-danger {
  padding: 0.75rem 1.5rem;
  background: #DC2626;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-danger:hover {
  background: #B91C1C;
}

.btn-danger:disabled,
.btn-primary:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
}

@media (max-width: 1024px) {
  .filters-sidebar {
    display: none;
  }
}

@media (max-width: 768px) {
  .top-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .header-actions {
    width: 100%;
  }

  .accords-grid {
    grid-template-columns: 1fr;
  }
}
</style>
