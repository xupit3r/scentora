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
        <!-- Skeleton Loading -->
        <div v-if="loading" class="accords-grid">
          <SkeletonCard v-for="i in 6" :key="i" />
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="error-state">
          <EmptyState
            icon="⚠️"
            title="Something went wrong"
            :description="error"
            action-label="Try Again"
            @action="loadAccords"
          />
        </div>

        <!-- Empty State - No Accords -->
        <EmptyState
          v-else-if="accords.length === 0 && !hasActiveFilters"
          icon="🌸"
          title="Your inventory is empty"
          description="Start by adding your first essential oil or accord"
          action-label="Add First Accord"
          @action="showForm = true"
        />

        <!-- Empty State - No Results -->
        <EmptyState
          v-else-if="accords.length === 0 && hasActiveFilters"
          icon="🔍"
          title="No accords match your filters"
          description="Try adjusting your search criteria"
          action-label="Clear Filters"
          @action="clearFilters"
        />

        <!-- Results -->
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
import { useMessage } from 'naive-ui';
import AccordCard from '@/components/AccordCard.vue';
import AccordForm from '@/components/AccordForm.vue';
import AccordFilters from '@/components/AccordFilters.vue';
import SkeletonCard from '@/components/ui/SkeletonCard.vue';
import EmptyState from '@/components/ui/EmptyState.vue';
import { useSimpleKeyboard } from '@/composables/useKeyboard';
import { accordService } from '@/services/api';
import type { Accord, AccordFilters as Filters, CreateAccordRequest, UpdateAccordRequest } from '@/types';

const router = useRouter();
const message = useMessage();

const accords = ref<Accord[]>([]);
const filters = ref<Filters>({});
const loading = ref(false);
const submitting = ref(false);
const error = ref('');
const showForm = ref(false);
const showFilters = ref(true);
const editingAccord = ref<Accord | undefined>();
const deletingAccord = ref<Accord | null>(null);

// Keyboard shortcuts
useSimpleKeyboard('n', () => {
  if (!showForm.value) {
    showForm.value = true;
    message.info('Press Esc to close');
  }
});

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
      message.success('Accord updated successfully');
    } else {
      await accordService.create(data as CreateAccordRequest);
      message.success('Accord created successfully');
    }
    
    closeForm();
    await loadAccords();
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'Failed to save accord';
    message.error(error.value);
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
    message.success('Accord deleted successfully');
    deletingAccord.value = null;
    await loadAccords();
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'Failed to delete accord';
    message.error(error.value);
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
}

.top-bar {
  background: white;
  border-bottom: 1px solid #E9E9E7;
  padding: 32px 48px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 32px;
}

.header-section h1 {
  font-size: 32px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #787774;
  margin: 0;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

.btn-filter {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: white;
  border: 1px solid #E9E9E7;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 200ms;
  font-size: 14px;
  color: #37352F;
}

.btn-filter svg {
  width: 18px;
  height: 18px;
}

.btn-filter:hover,
.btn-filter.active {
  background: #FAFAFA;
  border-color: #D9D9D7;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #0F766E;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 200ms;
  font-size: 14px;
}

.btn-primary svg {
  width: 18px;
  height: 18px;
}

.btn-primary:hover {
  background: #0D6460;
}

.main-content {
  display: flex;
  gap: 32px;
  padding: 32px 48px;
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

.error-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 64px 32px;
}

.results-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  border: 1px solid #E9E9E7;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #E9E9E7;
}

.results-count {
  font-size: 14px;
  color: #787774;
}

.results-count strong {
  color: #37352F;
  font-size: 16px;
}

.accords-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.btn-secondary {
  padding: 10px 16px;
  background: white;
  color: #37352F;
  border: 1px solid #E9E9E7;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 200ms;
  font-size: 14px;
}

.btn-secondary:hover {
  background: #FAFAFA;
  border-color: #D9D9D7;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
}

.delete-modal {
  background: white;
  border-radius: 16px;
  padding: 32px;
  max-width: 400px;
  text-align: center;
}

.delete-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.delete-modal h3 {
  font-size: 20px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 16px 0;
}

.delete-modal p {
  color: #787774;
  margin: 0 0 8px 0;
  font-size: 14px;
}

.warning-text {
  color: #991B1B;
  font-size: 14px;
  margin-bottom: 24px !important;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.btn-danger {
  padding: 10px 16px;
  background: #991B1B;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 200ms;
  font-size: 14px;
}

.btn-danger:hover {
  background: #7F1D1D;
}

.btn-danger:disabled,
.btn-primary:disabled {
  background: #9B9A97;
  cursor: not-allowed;
}

@media (max-width: 1024px) {
  .filters-sidebar {
    display: none;
  }

  .top-bar,
  .main-content {
    padding: 24px 32px;
  }
}

@media (max-width: 768px) {
  .top-bar {
    flex-direction: column;
    align-items: stretch;
    padding: 20px;
  }

  .header-actions {
    width: 100%;
  }

  .main-content {
    padding: 20px;
  }

  .accords-grid {
    grid-template-columns: 1fr;
  }
}
</style>
