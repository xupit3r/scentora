<template>
  <n-card
    :bordered="true"
    class="collection-card card-hover"
  >
    <!-- Header -->
    <div class="card-header">
      <div class="collection-icon">📁</div>
      <h3 class="collection-name">{{ collection.name }}</h3>
    </div>

    <!-- Body -->
    <div class="card-body">
      <p v-if="collection.description" class="collection-description">
        {{ truncatedDescription }}
      </p>

      <!-- Recipe Count -->
      <div class="recipe-count">
        <span class="count-icon">📝</span>
        <span class="count-text">
          {{ recipeCount }} recipe{{ recipeCount !== 1 ? 's' : '' }}
        </span>
      </div>

      <!-- Metadata -->
      <div class="metadata-row">
        <span class="meta-item">
          <span class="meta-icon">📅</span>
          <span class="meta-text">{{ formattedDate }}</span>
        </span>
      </div>
    </div>

    <!-- Actions -->
    <div class="card-actions">
      <n-button
        text
        @click="$emit('view', collection)"
        class="action-btn"
        title="View collection"
      >
        <template #icon>
          <span class="action-icon">👁️</span>
        </template>
      </n-button>
      <n-button
        text
        @click="$emit('edit', collection)"
        class="action-btn"
        title="Edit collection"
      >
        <template #icon>
          <span class="action-icon">✏️</span>
        </template>
      </n-button>
      <n-button
        text
        @click="$emit('delete', collection)"
        class="action-btn action-danger"
        title="Delete collection"
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
import { NCard, NButton } from 'naive-ui';
import type { RecipeCollection } from '@/types';

interface Props {
  collection: RecipeCollection;
  recipeCount?: number;
}

const props = withDefaults(defineProps<Props>(), {
  recipeCount: 0,
});

defineEmits<{
  view: [collection: RecipeCollection];
  edit: [collection: RecipeCollection];
  delete: [collection: RecipeCollection];
}>();

// Description truncation
const truncatedDescription = computed(() => {
  if (!props.collection.description) return '';
  const maxLength = 120;
  return props.collection.description.length > maxLength
    ? props.collection.description.substring(0, maxLength) + '...'
    : props.collection.description;
});

// Date formatting
const formattedDate = computed(() => {
  const date = new Date(props.collection.createdAt);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== new Date().getFullYear() ? 'numeric' : undefined,
  });
});
</script>

<style scoped>
.collection-card {
  position: relative;
  transition: all 0.2s ease;
  border-radius: 8px;
  background: var(--card-color);
  border-left: 4px solid #f0a020;
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
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.collection-icon {
  font-size: 24px;
}

.collection-name {
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

.collection-description {
  font-size: 14px;
  color: var(--text-color-2);
  line-height: 1.6;
  margin: 0;
}

.recipe-count {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color-2);
}

.count-icon {
  font-size: 16px;
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
  .collection-name {
    font-size: 16px;
  }

  .card-actions {
    opacity: 1; /* Always show on mobile */
  }
}
</style>
