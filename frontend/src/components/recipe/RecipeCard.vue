<template>
  <n-card
    :bordered="true"
    class="recipe-card card-hover"
    :class="`status-indicator-${recipe.status}`"
  >
    <!-- Header -->
    <div class="card-header">
      <h3 class="recipe-name">{{ recipe.name }}</h3>
      <n-tag :type="statusTagType" size="small" round>
        {{ statusLabel }}
      </n-tag>
    </div>

    <!-- Body -->
    <div class="card-body">
      <!-- Description -->
      <p v-if="recipe.description" class="recipe-description">
        {{ truncatedDescription }}
      </p>

      <!-- Volume Info -->
      <div class="volume-section">
        <div class="volume-primary">
          <span class="volume-label">Target:</span>
          <span class="volume-value">{{ recipe.targetVolumeMl }}</span>
          <span class="volume-unit">ml</span>
        </div>
        <div v-if="showVersionCount" class="version-info">
          <span class="version-icon">🔄</span>
          <span class="version-count">{{ versionCount }} version{{ versionCount !== 1 ? 's' : '' }}</span>
        </div>
      </div>

      <!-- Metadata -->
      <div class="metadata-row">
        <span class="meta-item">
          <span class="meta-icon">📅</span>
          <span class="meta-text">{{ formattedDate }}</span>
        </span>
        <span v-if="hasActiveVersion" class="meta-item">
          <span class="meta-icon">✨</span>
          <span class="meta-text">v{{ activeVersionNumber }}</span>
        </span>
      </div>
    </div>

    <!-- Actions (hover-reveal) -->
    <div class="card-actions">
      <n-button
        text
        @click="$emit('view', recipe)"
        class="action-btn"
        title="View recipe"
      >
        <template #icon>
          <span class="action-icon">👁️</span>
        </template>
      </n-button>
      <n-button
        text
        @click="$emit('edit', recipe)"
        class="action-btn"
        title="Edit recipe"
      >
        <template #icon>
          <span class="action-icon">✏️</span>
        </template>
      </n-button>
      <n-button
        text
        @click="$emit('duplicate', recipe)"
        class="action-btn"
        title="Duplicate recipe"
      >
        <template #icon>
          <span class="action-icon">📋</span>
        </template>
      </n-button>
      <n-button
        text
        @click="$emit('delete', recipe)"
        class="action-btn action-danger"
        title="Delete recipe"
      >
        <template #icon>
          <span class="action-icon">🗑️</span>
        </template>
      </n-button>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { NCard, NTag, NButton } from 'naive-ui';
import type { Recipe, RecipeStatus } from '@/types';

interface Props {
  recipe: Recipe;
  versionCount?: number;
  activeVersionNumber?: number;
}

const props = withDefaults(defineProps<Props>(), {
  versionCount: 0,
  activeVersionNumber: 1,
});

defineEmits<{
  view: [recipe: Recipe];
  edit: [recipe: Recipe];
  duplicate: [recipe: Recipe];
  delete: [recipe: Recipe];
}>();

// Status display
const statusLabels: Record<RecipeStatus, string> = {
  draft: 'Draft',
  in_progress: 'In Progress',
  tested: 'Tested',
  finalized: 'Finalized',
  archived: 'Archived',
};

const statusTagTypes: Record<RecipeStatus, 'default' | 'info' | 'warning' | 'success' | 'error'> = {
  draft: 'default',
  in_progress: 'info',
  tested: 'warning',
  finalized: 'success',
  archived: 'default',
};

const statusLabel = computed(() => statusLabels[props.recipe.status]);
const statusTagType = computed(() => statusTagTypes[props.recipe.status]);

// Description truncation
const truncatedDescription = computed(() => {
  if (!props.recipe.description) return '';
  const maxLength = 120;
  return props.recipe.description.length > maxLength
    ? props.recipe.description.substring(0, maxLength) + '...'
    : props.recipe.description;
});

// Version info
const showVersionCount = computed(() => props.versionCount > 0);
const hasActiveVersion = computed(() => props.activeVersionNumber > 0);

// Date formatting
const formattedDate = computed(() => {
  const date = new Date(props.recipe.createdAt);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== new Date().getFullYear() ? 'numeric' : undefined,
  });
});
</script>

<style scoped>
.recipe-card {
  position: relative;
  transition: all 0.2s ease;
  border-radius: 8px;
  background: var(--card-color);
  border-left: 4px solid var(--border-color);
}

/* Status indicators */
.recipe-card.status-indicator-draft {
  border-left-color: #e0e0e0;
}

.recipe-card.status-indicator-in_progress {
  border-left-color: #2080f0;
}

.recipe-card.status-indicator-tested {
  border-left-color: #f0a020;
}

.recipe-card.status-indicator-finalized {
  border-left-color: #18a058;
}

.recipe-card.status-indicator-archived {
  border-left-color: #909399;
  opacity: 0.7;
}

.card-hover:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-hover:hover .card-actions {
  opacity: 1;
}

/* Header */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
}

.recipe-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color-1);
  margin: 0;
  line-height: 1.4;
  flex: 1;
}

/* Body */
.card-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recipe-description {
  font-size: 14px;
  color: var(--text-color-2);
  line-height: 1.6;
  margin: 0;
}

.volume-section {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 0;
}

.volume-primary {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.volume-label {
  font-size: 12px;
  color: var(--text-color-3);
  margin-right: 4px;
}

.volume-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-color-1);
}

.volume-unit {
  font-size: 13px;
  color: var(--text-color-3);
}

.version-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-color-2);
}

.version-icon {
  font-size: 14px;
}

.version-count {
  font-weight: 500;
}

/* Metadata */
.metadata-row {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-color-3);
}

.meta-icon {
  font-size: 14px;
}

/* Actions */
.card-actions {
  display: flex;
  gap: 4px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--divider-color);
  opacity: 0;
  transition: opacity 0.2s ease;
}

.action-btn {
  color: var(--text-color-3);
  transition: all 0.15s ease;
}

.action-btn:hover {
  color: var(--primary-color);
  transform: scale(1.1);
}

.action-btn.action-danger:hover {
  color: #d03050;
}

.action-icon {
  font-size: 16px;
}

/* Responsive */
@media (max-width: 640px) {
  .recipe-name {
    font-size: 16px;
  }

  .volume-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .card-actions {
    opacity: 1; /* Always show on mobile */
  }
}
</style>
