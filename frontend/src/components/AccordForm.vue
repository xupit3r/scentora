<template>
  <n-modal
    :show="true"
    @update:show="(show: boolean) => !show && $emit('close')"
    preset="card"
    :title="isEdit ? 'Edit Accord' : 'New Accord'"
    class="accord-modal"
    style="width: 600px; max-width: 90vw;"
    :segmented="{ content: true }"
  >
    <n-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-placement="top"
      label-width="auto"
      require-mark-placement="right-hanging"
      size="medium"
    >
      <!-- Basic Information -->
      <div class="form-section">
        <h3 class="section-title">Basic Information</h3>
        
        <n-form-item label="Name" path="name" required>
          <n-input
            v-model:value="form.name"
            placeholder="e.g., Bergamot Essential Oil"
          />
        </n-form-item>

        <n-form-item label="Pyramid Position" path="pyramidPosition" required>
          <n-select
            v-model:value="form.pyramidPosition"
            :options="positionOptions"
            placeholder="Select position..."
          />
        </n-form-item>
      </div>

      <!-- Inventory -->
      <div class="form-section">
        <h3 class="section-title">Inventory</h3>
        
        <n-form-item label="Volume (ml)" path="volumeMl" required>
          <n-input
            v-model:value="form.volumeMl"
            type="number"
            placeholder="25.0"
          >
            <template #suffix>
              <span class="input-hint">≈ {{ estimatedDrops }} drops</span>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="Supplier" path="supplier">
          <n-input
            v-model:value="form.supplier"
            placeholder="e.g., Perfumer's Apprentice"
          />
        </n-form-item>

        <n-form-item label="Purchase Date" path="purchaseDate">
          <n-input
            v-model:value="form.purchaseDate"
            type="date"
          />
        </n-form-item>

        <n-form-item label="Dilution Percentage" path="dilutionPercentage">
          <n-input
            v-model:value="form.dilutionPercentage"
            type="number"
            placeholder="10.0"
          >
            <template #suffix>
              <span class="input-hint">%</span>
            </template>
          </n-input>
        </n-form-item>
      </div>

      <!-- Tags -->
      <div class="form-section">
        <h3 class="section-title">Tags</h3>
        <n-form-item path="tags">
          <TagSelector
            :model-value="form.tags || []"
            @update:model-value="(tags) => form.tags = tags"
            placeholder="Add tags to categorize this accord..."
          />
        </n-form-item>
      </div>

      <!-- Notes -->
      <div class="form-section">
        <h3 class="section-title">Notes</h3>
        <n-form-item path="notes">
          <n-input
            v-model:value="form.notes"
            type="textarea"
            :rows="4"
            placeholder="Any additional information, usage notes, or observations..."
          />
        </n-form-item>
      </div>
    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="$emit('close')">
          Cancel
        </n-button>
        <n-button
          type="primary"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ isEdit ? 'Update' : 'Create' }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { FormInst, FormRules } from 'naive-ui';
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

const formRef = ref<FormInst | null>(null);
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

const positionOptions = [
  { label: 'Top Note', value: 'top' },
  { label: 'Middle Note', value: 'middle' },
  { label: 'Base Note', value: 'base' },
];

const rules: FormRules = {
  name: {
    required: true,
    message: 'Please enter a name',
    trigger: 'blur',
  },
  pyramidPosition: {
    required: true,
    message: 'Please select a position',
    trigger: 'change',
  },
  volumeMl: {
    required: true,
    type: 'number',
    message: 'Please enter volume in ml',
    trigger: 'blur',
  },
};

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
      tags: accord.tags ? [...accord.tags] : [],
    };
  }
}, { immediate: true });

const estimatedDrops = computed(() => {
  return Math.round((Number(form.value.volumeMl) || 0) * 20);
});

function handleSubmit() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      const data: CreateAccordRequest | UpdateAccordRequest = {
        name: form.value.name,
        pyramidPosition: form.value.pyramidPosition,
        volumeMl: Number(form.value.volumeMl),
        tags: form.value.tags,
      };

      if (form.value.supplier) data.supplier = form.value.supplier;
      if (form.value.purchaseDate) data.purchaseDate = form.value.purchaseDate;
      if (form.value.dilutionPercentage !== undefined) {
        data.dilutionPercentage = Number(form.value.dilutionPercentage);
      }
      if (form.value.notes) data.notes = form.value.notes;

      emit('submit', data);
    }
  });
}
</script>

<style scoped>
.accord-modal {
  max-height: 90vh;
}

.form-section {
  margin-bottom: 24px;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid #E9E9E7;
}

.input-hint {
  font-size: 12px;
  color: #9B9A97;
}
</style>
