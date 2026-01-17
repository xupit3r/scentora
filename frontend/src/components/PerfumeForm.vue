<template>
  <div class="perfume-form">
    <h2>{{ isEditing ? 'Edit Perfume' : 'Add New Perfume' }}</h2>
    
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="name">Name *</label>
        <input
          id="name"
          v-model="form.name"
          type="text"
          required
          placeholder="e.g., Bleu de Chanel"
        />
      </div>

      <div class="form-group">
        <label for="designer">Designer *</label>
        <input
          id="designer"
          v-model="form.designer"
          type="text"
          required
          placeholder="e.g., Chanel"
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="year">Year</label>
          <input
            id="year"
            v-model.number="form.year"
            type="number"
            min="1800"
            :max="new Date().getFullYear() + 1"
            placeholder="e.g., 2010"
          />
        </div>

        <div class="form-group">
          <label for="concentration">Concentration</label>
          <select id="concentration" v-model="form.concentration">
            <option value="">Select...</option>
            <option value="Parfum">Parfum</option>
            <option value="EDP">Eau de Parfum (EDP)</option>
            <option value="EDT">Eau de Toilette (EDT)</option>
            <option value="EDC">Eau de Cologne (EDC)</option>
          </select>
        </div>
      </div>

      <div class="form-group">
        <label for="imageUrl">Image URL</label>
        <input
          id="imageUrl"
          v-model="form.imageUrl"
          type="url"
          placeholder="https://example.com/image.jpg"
        />
      </div>

      <div class="pyramid-section">
        <h3>Perfume Pyramid</h3>
        
        <div class="form-group">
          <label for="top-notes">Top Notes</label>
          <input
            id="top-notes"
            v-model="topNotesInput"
            type="text"
            placeholder="e.g., Bergamot, Lemon, Grapefruit (comma-separated)"
            list="notes-datalist"
          />
          <div class="note-chips">
            <span v-for="(note, index) in form.pyramid.top" :key="index" class="chip">
              {{ note }}
              <button type="button" @click="removeNote('top', index)">×</button>
            </span>
          </div>
        </div>

        <div class="form-group">
          <label for="middle-notes">Middle Notes (Heart)</label>
          <input
            id="middle-notes"
            v-model="middleNotesInput"
            type="text"
            placeholder="e.g., Jasmine, Rose, Ginger"
            list="notes-datalist"
          />
          <div class="note-chips">
            <span v-for="(note, index) in form.pyramid.middle" :key="index" class="chip">
              {{ note }}
              <button type="button" @click="removeNote('middle', index)">×</button>
            </span>
          </div>
        </div>

        <div class="form-group">
          <label for="base-notes">Base Notes</label>
          <input
            id="base-notes"
            v-model="baseNotesInput"
            type="text"
            placeholder="e.g., Sandalwood, Vetiver, Cedar"
            list="notes-datalist"
          />
          <datalist id="notes-datalist">
            <option v-for="note in availableNotes" :key="note" :value="note" />
          </datalist>
          <div class="note-chips">
            <span v-for="(note, index) in form.pyramid.base" :key="index" class="chip">
              {{ note }}
              <button type="button" @click="removeNote('base', index)">×</button>
            </span>
          </div>
        </div>
      </div>

      <div class="form-group">
        <label for="description">Description</label>
        <textarea
          id="description"
          v-model="form.description"
          rows="4"
          placeholder="Notes about this perfume..."
        ></textarea>
      </div>

      <div class="form-actions">
        <button type="button" class="btn-secondary" @click="$emit('cancel')">
          Cancel
        </button>
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? 'Saving...' : 'Save Perfume' }}
        </button>
      </div>

      <div v-if="error" class="error-message">{{ error }}</div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { notesService } from '@/services/api';
import type { Perfume } from '@/types';

interface Props {
  perfume?: Perfume;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  submit: [perfume: Partial<Perfume>];
  cancel: [];
}>();

const isEditing = ref(!!props.perfume);
const loading = ref(false);
const error = ref<string | null>(null);
const availableNotes = ref<string[]>([]);

const form = ref({
  name: props.perfume?.name || '',
  designer: props.perfume?.designer || '',
  year: props.perfume?.year,
  concentration: props.perfume?.concentration || '',
  imageUrl: props.perfume?.imageUrl || '',
  pyramid: {
    top: props.perfume?.pyramid?.top || [],
    middle: props.perfume?.pyramid?.middle || [],
    base: props.perfume?.pyramid?.base || [],
  },
  description: props.perfume?.description || '',
});

const topNotesInput = ref('');
const middleNotesInput = ref('');
const baseNotesInput = ref('');

// Load available notes for autocomplete
onMounted(async () => {
  try {
    availableNotes.value = await notesService.getAll();
  } catch (err) {
    console.error('Failed to load notes:', err);
  }
});

// Watch for comma-separated input and convert to array
watch(topNotesInput, (value) => {
  if (value.includes(',')) {
    const notes = value.split(',').map(n => n.trim()).filter(n => n);
    form.value.pyramid.top = [...new Set([...form.value.pyramid.top, ...notes])];
    topNotesInput.value = '';
  }
});

watch(middleNotesInput, (value) => {
  if (value.includes(',')) {
    const notes = value.split(',').map(n => n.trim()).filter(n => n);
    form.value.pyramid.middle = [...new Set([...form.value.pyramid.middle, ...notes])];
    middleNotesInput.value = '';
  }
});

watch(baseNotesInput, (value) => {
  if (value.includes(',')) {
    const notes = value.split(',').map(n => n.trim()).filter(n => n);
    form.value.pyramid.base = [...new Set([...form.value.pyramid.base, ...notes])];
    baseNotesInput.value = '';
  }
});

function removeNote(type: 'top' | 'middle' | 'base', index: number) {
  form.value.pyramid[type].splice(index, 1);
}

async function handleSubmit() {
  loading.value = true;
  error.value = null;

  try {
    // Add any remaining notes
    if (topNotesInput.value.trim()) {
      form.value.pyramid.top.push(topNotesInput.value.trim());
    }
    if (middleNotesInput.value.trim()) {
      form.value.pyramid.middle.push(middleNotesInput.value.trim());
    }
    if (baseNotesInput.value.trim()) {
      form.value.pyramid.base.push(baseNotesInput.value.trim());
    }

    emit('submit', form.value);
  } catch (err: any) {
    error.value = err.message || 'Failed to save perfume';
    loading.value = false;
  }
}
</script>

<style scoped>
.perfume-form {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  max-width: 800px;
  margin: 0 auto;
}

h2 {
  margin-bottom: 2rem;
  color: #333;
}

h3 {
  margin-bottom: 1rem;
  color: #6b4f9e;
  font-size: 1.2rem;
}

.form-group {
  margin-bottom: 1.5rem;
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
select,
textarea {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

input:focus,
select:focus,
textarea:focus {
  outline: none;
  border-color: #6b4f9e;
}

textarea {
  resize: vertical;
  font-family: inherit;
}

.pyramid-section {
  background: #f9f9f9;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.note-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: #6b4f9e;
  color: white;
  padding: 0.4rem 0.8rem;
  border-radius: 20px;
  font-size: 0.9rem;
}

.chip button {
  background: none;
  border: none;
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.btn-primary,
.btn-secondary {
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

.btn-primary:hover:not(:disabled) {
  background: #5a3f8d;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: #e0e0e0;
  color: #666;
}

.btn-secondary:hover {
  background: #d0d0d0;
}

.error-message {
  color: #d32f2f;
  margin-top: 1rem;
  padding: 0.75rem;
  background: #ffebee;
  border-radius: 8px;
}
</style>
