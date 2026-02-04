<template>
  <div class="version-selector">
    <n-select
      :value="modelValue"
      :options="versionOptions"
      :loading="loading"
      placeholder="Select version..."
      @update:value="handleSelect"
      :disabled="!hasVersions || disabled"
    >
      <template #empty>
        <div class="empty-versions">
          <span class="empty-icon">📝</span>
          <span class="empty-text">No versions yet</span>
        </div>
      </template>
    </n-select>

    <n-button
      v-if="showCreateButton"
      type="primary"
      :disabled="disabled"
      @click="$emit('create')"
      class="create-version-btn"
    >
      <template #icon>
        <span>➕</span>
      </template>
      New Version
    </n-button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { NSelect, NButton } from 'naive-ui';
import type { RecipeVersion } from '@/types';

interface Props {
  versions: RecipeVersion[];
  modelValue?: string; // Selected version ID
  loading?: boolean;
  disabled?: boolean;
  showCreateButton?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  disabled: false,
  showCreateButton: true,
});

const emit = defineEmits<{
  'update:modelValue': [value: string];
  select: [version: RecipeVersion];
  create: [];
}>();

const hasVersions = computed(() => props.versions.length > 0);

const versionOptions = computed(() => {
  return props.versions.map((version) => ({
    label: formatVersionLabel(version),
    value: version._id,
    version,
  }));
});

function formatVersionLabel(version: RecipeVersion): string {
  const isActive = version.isActive ? ' ⭐' : '';
  const date = new Date(version.createdAt).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  });
  return `v${version.versionNumber} - ${version.name}${isActive} (${date})`;
}

function handleSelect(value: string) {
  const selected = props.versions.find((v) => v._id === value);
  if (selected) {
    emit('update:modelValue', value);
    emit('select', selected);
  }
}
</script>

<style scoped>
.version-selector {
  display: flex;
  gap: 12px;
  align-items: center;
}

.version-selector :deep(.n-select) {
  flex: 1;
  min-width: 200px;
}

.create-version-btn {
  flex-shrink: 0;
}

.empty-versions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
  color: var(--text-color-3);
}

.empty-icon {
  font-size: 32px;
}

.empty-text {
  font-size: 14px;
}

/* Responsive */
@media (max-width: 640px) {
  .version-selector {
    flex-direction: column;
    align-items: stretch;
  }

  .version-selector :deep(.n-select) {
    min-width: unset;
  }
}
</style>
