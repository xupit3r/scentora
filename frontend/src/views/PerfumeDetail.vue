<template>
  <div class="perfume-detail">
    <div v-if="loading" class="loading">Loading perfume details...</div>
    
    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button class="btn-secondary" @click="router.push('/')">Back to Collection</button>
    </div>

    <div v-else-if="editMode && perfume" class="edit-mode">
      <PerfumeForm
        :perfume="perfume"
        @submit="handleUpdatePerfume"
        @cancel="editMode = false"
      />
    </div>
    
    <div v-else-if="perfume" class="content">
      <div class="header">
        <button class="back-btn" @click="router.push('/')">← Back</button>
        <div class="actions">
          <button class="btn-secondary" @click="startEdit">Edit</button>
          <button class="btn-danger" @click="handleDelete">Delete</button>
        </div>
      </div>

      <div class="perfume-main">
        <div class="perfume-image-large">
          <img v-if="perfume.imageUrl" :src="perfume.imageUrl" :alt="perfume.name" />
          <div v-else class="placeholder">🌸</div>
        </div>

        <div class="perfume-details">
          <h1>{{ perfume.name }}</h1>
          <p class="designer">by {{ perfume.designer }}</p>
          <div class="meta">
            <span v-if="perfume.year" class="meta-item">{{ perfume.year }}</span>
            <span v-if="perfume.concentration" class="meta-item">{{ perfume.concentration }}</span>
          </div>
          
          <p v-if="perfume.description" class="description">{{ perfume.description }}</p>

          <div class="pyramid-section">
            <h2>Perfume Pyramid</h2>
            <PerfumePyramid :pyramid="perfume.pyramid" />
          </div>
        </div>
      </div>

      <div class="journal-section">
        <div class="journal-header">
          <h2>Journal Entries</h2>
          <button class="btn-primary" @click="showJournalForm ? cancelEditEntry() : (showJournalForm = true)">
            {{ showJournalForm ? 'Cancel' : '+ Add Entry' }}
          </button>
        </div>

        <div v-if="showJournalForm" class="journal-form">
          <h3>{{ editingEntry ? 'Edit Entry' : 'New Entry' }}</h3>
          <form @submit.prevent="handleAddJournalEntry">
            <div class="form-group">
              <label for="date">Date</label>
              <input
                id="date"
                v-model="journalForm.date"
                type="date"
                required
              />
            </div>

            <div class="form-group">
              <label for="rating">Rating (1-10)</label>
              <input
                id="rating"
                v-model.number="journalForm.rating"
                type="number"
                min="1"
                max="10"
              />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="occasion">Occasion</label>
                <input
                  id="occasion"
                  v-model="journalForm.occasion"
                  type="text"
                  placeholder="e.g., Work, Date night"
                />
              </div>

              <div class="form-group">
                <label for="weather">Weather</label>
                <input
                  id="weather"
                  v-model="journalForm.weather"
                  type="text"
                  placeholder="e.g., Sunny, Cold"
                />
              </div>
            </div>

            <div class="form-group">
              <label for="content">Notes *</label>
              <textarea
                id="content"
                v-model="journalForm.content"
                rows="4"
                required
                placeholder="Your thoughts about this wearing..."
              ></textarea>
            </div>

            <button type="submit" class="btn-primary">
              {{ editingEntry ? 'Update Entry' : 'Save Entry' }}
            </button>
          </form>
        </div>

        <div v-if="journalEntries.length === 0 && !showJournalForm" class="empty-journal">
          <p>No journal entries yet. Start documenting your experiences!</p>
        </div>

        <div v-else class="journal-entries">
          <div v-for="entry in journalEntries" :key="entry._id" class="journal-entry">
            <div class="entry-header">
              <span class="entry-date">{{ formatDate(entry.date) }}</span>
              <div class="entry-actions">
                <span v-if="entry.rating" class="entry-rating">⭐ {{ entry.rating }}/10</span>
                <button class="action-btn-small" @click="startEditEntry(entry)" title="Edit">
                  ✏️
                </button>
                <button class="action-btn-small danger" @click="handleDeleteEntry(entry._id!)" title="Delete">
                  🗑️
                </button>
              </div>
            </div>
            <div v-if="entry.occasion || entry.weather" class="entry-meta">
              <span v-if="entry.occasion" class="tag">{{ entry.occasion }}</span>
              <span v-if="entry.weather" class="tag">{{ entry.weather }}</span>
            </div>
            <p class="entry-content">{{ entry.content }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { perfumeService, journalService } from '@/services/api';
import PerfumePyramid from '@/components/PerfumePyramid.vue';
import PerfumeForm from '@/components/PerfumeForm.vue';
import type { Perfume, JournalEntry } from '@/types';

const route = useRoute();
const router = useRouter();

const perfume = ref<Perfume | null>(null);
const journalEntries = ref<JournalEntry[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const showJournalForm = ref(false);
const editMode = ref(false);
const editingEntry = ref<JournalEntry | null>(null);

const journalForm = ref({
  date: new Date().toISOString().split('T')[0],
  content: '',
  rating: undefined as number | undefined,
  occasion: '',
  weather: '',
});

async function loadPerfume() {
  try {
    const id = route.params.id as string;
    perfume.value = await perfumeService.getById(id);
    journalEntries.value = await journalService.getByPerfumeId(id);
  } catch (err: any) {
    error.value = err.message || 'Failed to load perfume';
  } finally {
    loading.value = false;
  }
}

async function handleUpdatePerfume(perfumeData: Partial<Perfume>) {
  if (!perfume.value?._id) return;

  try {
    perfume.value = await perfumeService.update(perfume.value._id, perfumeData);
    editMode.value = false;
  } catch (err: any) {
    error.value = err.message || 'Failed to update perfume';
  }
}

async function handleAddJournalEntry() {
  if (!perfume.value?._id) return;

  try {
    if (editingEntry.value?._id) {
      // Update existing entry
      await journalService.update(editingEntry.value._id, journalForm.value);
      editingEntry.value = null;
    } else {
      // Create new entry
      await journalService.create({
        perfumeId: perfume.value._id,
        ...journalForm.value,
      } as any);
    }

    // Reset form
    journalForm.value = {
      date: new Date().toISOString().split('T')[0],
      content: '',
      rating: undefined,
      occasion: '',
      weather: '',
    };

    showJournalForm.value = false;
    
    // Reload journal entries
    journalEntries.value = await journalService.getByPerfumeId(perfume.value._id);
  } catch (err: any) {
    error.value = err.message || 'Failed to save journal entry';
  }
}

function startEditEntry(entry: JournalEntry) {
  editingEntry.value = entry;
  journalForm.value = {
    date: entry.date,
    content: entry.content,
    rating: entry.rating,
    occasion: entry.occasion || '',
    weather: entry.weather || '',
  };
  showJournalForm.value = true;
}

function cancelEditEntry() {
  editingEntry.value = null;
  journalForm.value = {
    date: new Date().toISOString().split('T')[0],
    content: '',
    rating: undefined,
    occasion: '',
    weather: '',
  };
  showJournalForm.value = false;
}

async function handleDeleteEntry(entryId: string) {
  if (!confirm('Are you sure you want to delete this journal entry?')) return;

  try {
    await journalService.delete(entryId);
    journalEntries.value = journalEntries.value.filter(e => e._id !== entryId);
  } catch (err: any) {
    error.value = err.message || 'Failed to delete journal entry';
  }
}

function startEdit() {
  editMode.value = true;
}

async function handleDelete() {
  if (!perfume.value?._id) return;

  if (confirm(`Are you sure you want to delete "${perfume.value.name}"?`)) {
    try {
      await perfumeService.delete(perfume.value._id);
      router.push('/');
    } catch (err: any) {
      error.value = err.message || 'Failed to delete perfume';
    }
  }
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  });
}

onMounted(() => {
  loadPerfume();
});
</script>

<style scoped>
.perfume-detail {
  width: 100%;
}

.loading,
.error {
  text-align: center;
  padding: 3rem;
}

.error {
  color: #d32f2f;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.back-btn {
  background: none;
  border: none;
  color: #6b4f9e;
  font-size: 1rem;
  cursor: pointer;
  padding: 0.5rem 1rem;
  transition: opacity 0.2s;
}

.back-btn:hover {
  opacity: 0.7;
}

.actions {
  display: flex;
  gap: 1rem;
}

.perfume-main {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 2rem;
  background: white;
  padding: 2rem;
  border-radius: 12px;
  margin-bottom: 2rem;
}

.perfume-image-large {
  width: 100%;
  aspect-ratio: 1;
  background: linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.perfume-image-large img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder {
  font-size: 6rem;
}

.perfume-details h1 {
  font-size: 2.5rem;
  color: #333;
  margin-bottom: 0.5rem;
}

.designer {
  font-size: 1.3rem;
  color: #666;
  margin-bottom: 1rem;
}

.meta {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.meta-item {
  padding: 0.5rem 1rem;
  background: #f5f5f5;
  border-radius: 20px;
  font-size: 0.9rem;
  color: #666;
}

.description {
  color: #555;
  line-height: 1.6;
  margin-bottom: 2rem;
}

.pyramid-section h2 {
  font-size: 1.3rem;
  color: #6b4f9e;
  margin-bottom: 1rem;
}

.journal-section {
  background: white;
  padding: 2rem;
  border-radius: 12px;
}

.journal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.journal-header h2 {
  font-size: 1.5rem;
  color: #333;
}

.journal-form {
  background: #f9f9f9;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #555;
}

input,
textarea {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
}

input:focus,
textarea:focus {
  outline: none;
  border-color: #6b4f9e;
}

textarea {
  resize: vertical;
  font-family: inherit;
}

.btn-primary,
.btn-secondary,
.btn-danger {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #6b4f9e;
  color: white;
}

.btn-primary:hover {
  background: #5a3f8d;
}

.btn-secondary {
  background: #e0e0e0;
  color: #666;
}

.btn-secondary:hover {
  background: #d0d0d0;
}

.btn-danger {
  background: #d32f2f;
  color: white;
}

.btn-danger:hover {
  background: #b71c1c;
}

.empty-journal {
  text-align: center;
  padding: 2rem;
  color: #999;
}

.journal-entries {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.journal-entry {
  padding: 1.5rem;
  background: #f9f9f9;
  border-radius: 8px;
  border-left: 4px solid #6b4f9e;
}

.entry-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.entry-actions {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.action-btn-small {
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.action-btn-small:hover {
  background: #f5f5f5;
  transform: scale(1.05);
}

.action-btn-small.danger:hover {
  background: #ffebee;
  border-color: #d32f2f;
}

.edit-mode {
  width: 100%;
}

.journal-form h3 {
  margin-bottom: 1rem;
  color: #6b4f9e;
}

.entry-date {
  font-weight: 600;
  color: #333;
}

.entry-rating {
  color: #6b4f9e;
  font-weight: 600;
}

.entry-meta {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.tag {
  padding: 0.25rem 0.75rem;
  background: white;
  border-radius: 15px;
  font-size: 0.85rem;
  color: #666;
}

.entry-content {
  color: #555;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .perfume-main {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
