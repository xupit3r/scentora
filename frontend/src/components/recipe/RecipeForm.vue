<template>
  <n-modal
    :show="true"
    @update:show="(show: boolean) => !show && $emit('close')"
    preset="card"
    :title="isEdit ? 'Edit Recipe' : 'New Recipe'"
    class="recipe-modal"
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

        <n-form-item label="Recipe Name" path="name" required>
          <n-input
            v-model:value="form.name"
            placeholder="e.g., Summer Citrus Cologne"
            maxlength="255"
            show-count
          />
        </n-form-item>

        <n-form-item label="Description" path="description">
          <n-input
            v-model:value="form.description"
            type="textarea"
            placeholder="Describe your perfume recipe..."
            :rows="3"
            maxlength="500"
            show-count
          />
        </n-form-item>

        <n-form-item label="Target Volume (ml)" path="targetVolumeMl" required>
          <n-input
            v-model:value="form.targetVolumeMl"
            type="number"
            placeholder="100"
            min="0"
            step="0.1"
          >
            <template #suffix>
              <span class="input-hint">ml</span>
            </template>
          </n-input>
          <template #feedback>
            <span class="form-hint">
              This is the total volume you plan to make
            </span>
          </template>
        </n-form-item>

        <n-form-item v-if="isEdit" label="Status" path="status">
          <n-select
            v-model:value="form.status"
            :options="statusOptions"
            placeholder="Select status..."
          />
        </n-form-item>
      </div>
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
          {{ isEdit ? 'Save Changes' : 'Create Recipe' }}
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
import type { Recipe, RecipeStatus, CreateRecipeRequest, UpdateRecipeRequest } from '@/types';

interface Props {
  recipe?: Recipe;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

const emit = defineEmits<{
  close: [];
  submit: [data: CreateRecipeRequest | UpdateRecipeRequest];
}>();

const formRef = ref<FormInst | null>(null);

const isEdit = computed(() => !!props.recipe);

// Form data
const form = ref<{
  name: string;
  description: string;
  targetVolumeMl: string;
  status?: RecipeStatus;
}>({
  name: '',
  description: '',
  targetVolumeMl: '',
  status: 'draft',
});

// Initialize form with recipe data if editing
watch(
  () => props.recipe,
  (recipe) => {
    if (recipe) {
      form.value = {
        name: recipe.name,
        description: recipe.description || '',
        targetVolumeMl: recipe.targetVolumeMl.toString(),
        status: recipe.status,
      };
    }
  },
  { immediate: true }
);

// Status options
const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Tested', value: 'tested' },
  { label: 'Finalized', value: 'finalized' },
  { label: 'Archived', value: 'archived' },
];

// Form validation rules
const rules: FormRules = {
  name: [
    { required: true, message: 'Recipe name is required', trigger: ['blur', 'input'] },
    { min: 3, message: 'Name must be at least 3 characters', trigger: ['blur', 'input'] },
    { max: 255, message: 'Name must be less than 255 characters', trigger: ['blur', 'input'] },
  ],
  targetVolumeMl: [
    { required: true, message: 'Target volume is required', trigger: ['blur', 'input'] },
    {
      validator: (rule, value) => {
        const num = parseFloat(value);
        if (isNaN(num) || num <= 0) {
          return new Error('Volume must be greater than 0');
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

    const data = {
      name: form.value.name.trim(),
      description: form.value.description.trim() || undefined,
      targetVolumeMl: parseFloat(form.value.targetVolumeMl),
      ...(isEdit.value && form.value.status ? { status: form.value.status } : {}),
    };

    emit('submit', data);
  } catch (error) {
    console.error('Form validation failed:', error);
  }
}
</script>

<style scoped>
.recipe-modal :deep(.n-card__content) {
  padding: 0;
}

.form-section {
  padding: 24px;
  border-bottom: 1px solid var(--divider-color);
}

.form-section:last-child {
  border-bottom: none;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color-1);
  margin: 0 0 16px 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
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
  padding: 16px 24px;
}

/* Responsive */
@media (max-width: 640px) {
  .form-section {
    padding: 16px;
  }

  .modal-footer {
    padding: 12px 16px;
  }
}
</style>
