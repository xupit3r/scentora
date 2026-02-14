<template>
  <n-modal
    :show="true"
    @update:show="(show: boolean) => !show && $emit('close')"
    preset="card"
    :title="isEdit ? 'Edit Ingredient' : 'Add Ingredient'"
    class="ingredient-modal"
    style="width: 500px; max-width: 90vw;"
  >
    <n-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-placement="top"
      size="medium"
    >
      <n-form-item label="Select Accord" path="accordId" required>
        <n-select
          v-model:value="form.accordId"
          :options="accordOptions"
          :loading="accordsLoading"
          filterable
          placeholder="Search for an accord..."
          :disabled="isEdit"
        >
          <template #empty>
            <div class="empty-accords">
              <span>No accords found</span>
            </div>
          </template>
        </n-select>
      </n-form-item>

      <n-form-item label="Quantity (ml)" path="quantityMl" required>
        <n-input
          v-model:value="form.quantityMl"
          placeholder="15.0"
        >
          <template #suffix>
            <span class="input-hint">≈ {{ estimatedDrops }} drops</span>
          </template>
        </n-input>
      </n-form-item>

      <n-form-item label="Percentage" path="percentage">
        <n-input
          v-model:value="form.percentage"
          placeholder="15.0"
        >
          <template #suffix>
            <span class="input-hint">%</span>
          </template>
        </n-input>
        <template #feedback>
          <span class="form-hint">
            Optional: Percentage of total formula
          </span>
        </template>
      </n-form-item>

      <n-form-item label="Notes" path="notes">
        <n-input
          v-model:value="form.notes"
          type="textarea"
          placeholder="Optional notes about this ingredient..."
          :rows="2"
          maxlength="200"
          show-count
        />
      </n-form-item>
    </n-form>

    <template #footer>
      <div class="modal-footer">
        <n-button @click="$emit('close')">
          Cancel
        </n-button>
        <n-button
          type="primary"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ isEdit ? 'Save Changes' : 'Add Ingredient' }}
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import {
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NButton,
  type FormInst,
  type FormRules,
} from 'naive-ui';
import type {
  Accord,
  RecipeIngredient,
  CreateRecipeIngredientRequest,
  UpdateRecipeIngredientRequest,
} from '@/types';

interface Props {
  accords: Accord[];
  ingredient?: RecipeIngredient;
  loading?: boolean;
  accordsLoading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  accordsLoading: false,
});

const emit = defineEmits<{
  close: [];
  submit: [data: CreateRecipeIngredientRequest | UpdateRecipeIngredientRequest];
}>();

const formRef = ref<FormInst | null>(null);

const isEdit = computed(() => !!props.ingredient);

// Form data
const form = ref<{
  accordId: string;
  quantityMl: string;
  percentage: string;
  notes: string;
}>({
  accordId: '',
  quantityMl: '',
  percentage: '',
  notes: '',
});

// Initialize form if editing
watch(
  () => props.ingredient,
  (ingredient) => {
    if (ingredient) {
      form.value = {
        accordId: ingredient.accordId,
        quantityMl: ingredient.quantityMl.toString(),
        percentage: ingredient.percentage?.toString() || '',
        notes: ingredient.notes || '',
      };
    }
  },
  { immediate: true }
);

// Accord options
const accordOptions = computed(() => {
  return props.accords.map((accord) => ({
    label: `${accord.name} (${accord.pyramidPosition})`,
    value: accord._id,
  }));
});

// Estimated drops
const estimatedDrops = computed(() => {
  const ml = parseFloat(form.value.quantityMl);
  if (isNaN(ml) || ml <= 0) return 0;
  return Math.round(ml * 20);
});

// Form validation rules
const rules: FormRules = {
  accordId: [
    { required: true, message: 'Please select an accord', trigger: 'change' },
  ],
  quantityMl: [
    { required: true, message: 'Quantity is required', trigger: ['blur', 'input'] },
    {
      validator: (_rule, value) => {
        const num = parseFloat(value);
        if (isNaN(num) || num <= 0) {
          return new Error('Quantity must be greater than 0');
        }
        return true;
      },
      trigger: ['blur', 'input'],
    },
  ],
  percentage: [
    {
      validator: (_rule, value) => {
        if (!value) return true; // Optional field
        const num = parseFloat(value);
        if (isNaN(num) || num < 0 || num > 100) {
          return new Error('Percentage must be between 0 and 100');
        }
        return true;
      },
      trigger: ['blur', 'input'],
    },
  ],
};

// Submit handler
async function handleSubmit() {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();

    const data: any = {
      quantityMl: parseFloat(form.value.quantityMl),
    };

    if (!isEdit.value) {
      data.accordId = form.value.accordId;
    }

    if (form.value.percentage) {
      data.percentage = parseFloat(form.value.percentage);
    }

    if (form.value.notes.trim()) {
      data.notes = form.value.notes.trim();
    }

    emit('submit', data);
  } catch (error) {
    console.error('Form validation failed:', error);
  }
}
</script>

<style scoped>
.ingredient-modal :deep(.n-card__content) {
  padding: 24px;
}

.empty-accords {
  padding: 24px;
  text-align: center;
  color: var(--text-color-3);
  font-size: 14px;
}

.input-hint {
  color: var(--text-color-3);
  font-size: 13px;
}

.form-hint {
  color: var(--text-color-3);
  font-size: 12px;
  font-style: italic;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Responsive */
@media (max-width: 640px) {
  .ingredient-modal :deep(.n-card__content) {
    padding: 16px;
  }
}
</style>
