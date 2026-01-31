<template>
  <n-card
    :bordered="true"
    class="accord-card card-hover"
    :class="`position-indicator-${accord.pyramidPosition}`"
  >
    <!-- Header -->
    <div class="card-header">
      <h3 class="accord-name">{{ accord.name }}</h3>
      <n-tag :type="positionTagType" size="small" round>
        {{ positionLabel }}
      </n-tag>
    </div>

    <!-- Body -->
    <div class="card-body">
      <!-- Volume Info -->
      <div class="volume-section">
        <div class="volume-primary">
          <span class="volume-value">{{ accord.volumeMl }}</span>
          <span class="volume-unit">ml</span>
        </div>
        <div class="volume-secondary">
          <span class="drops-value">{{ accord.volumeDrops }}</span>
          <span class="drops-unit">drops</span>
        </div>
        <n-tag
          v-if="isLowStock"
          :type="stockLevel === 'critical' ? 'error' : 'warning'"
          size="small"
          class="stock-badge"
        >
          {{ stockWarning }}
        </n-tag>
      </div>

      <!-- Supplier -->
      <div v-if="accord.supplier" class="info-row">
        <span class="info-icon">🏢</span>
        <span class="info-text">{{ accord.supplier }}</span>
      </div>

      <!-- Dilution -->
      <div v-if="accord.dilutionPercentage" class="info-row">
        <span class="info-label">Dilution:</span>
        <span class="info-value">{{ accord.dilutionPercentage }}%</span>
      </div>

      <!-- Tags -->
      <div v-if="accord.tags.length" class="tags-section">
        <n-tag
          v-for="tag in accord.tags.slice(0, 3)"
          :key="tag"
          size="small"
          :bordered="false"
        >
          {{ tag }}
        </n-tag>
        <n-tag
          v-if="accord.tags.length > 3"
          size="small"
          :bordered="false"
          class="tag-more"
        >
          +{{ accord.tags.length - 3 }}
        </n-tag>
      </div>
    </div>

    <!-- Actions -->
    <div class="card-actions">
      <n-button
        text
        @click="$emit('view', accord)"
        class="action-btn"
      >
        <template #icon>
          <span class="action-icon">👁️</span>
        </template>
      </n-button>
      <n-button
        text
        @click="$emit('edit', accord)"
        class="action-btn"
      >
        <template #icon>
          <span class="action-icon">✏️</span>
        </template>
      </n-button>
      <n-button
        text
        @click="$emit('delete', accord)"
        class="action-btn action-danger"
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
import type { Accord } from '@/types';

const props = defineProps<{
  accord: Accord;
}>();

defineEmits<{
  view: [accord: Accord];
  edit: [accord: Accord];
  delete: [accord: Accord];
}>();

const positionLabel = computed(() => {
  const labels: Record<string, string> = {
    top: 'Top Note',
    middle: 'Middle Note',
    base: 'Base Note',
  };
  return labels[props.accord.pyramidPosition] || props.accord.pyramidPosition;
});

const positionTagType = computed(() => {
  const types: Record<string, 'warning' | 'info' | 'default'> = {
    top: 'warning',
    middle: 'info',
    base: 'default',
  };
  return types[props.accord.pyramidPosition] || 'default';
});

const isLowStock = computed(() => props.accord.volumeMl < 5);

const stockLevel = computed(() => {
  if (props.accord.volumeMl < 1) return 'critical';
  if (props.accord.volumeMl < 5) return 'low';
  return '';
});

const stockWarning = computed(() => {
  if (props.accord.volumeMl < 1) return 'Critical!';
  if (props.accord.volumeMl < 5) return 'Low Stock';
  return '';
});
</script>

<style scoped>
.accord-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* Position indicators - subtle left border */
.position-indicator-top {
  border-left: 3px solid #FDE68A;
}

.position-indicator-middle {
  border-left: 3px solid #D8B4FE;
}

.position-indicator-base {
  border-left: 3px solid #E4C9A0;
}

/* Header */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
}

.accord-name {
  font-size: 20px;
  font-weight: 600;
  color: #37352F;
  margin: 0;
  flex: 1;
  line-height: 1.3;
}

/* Body */
.card-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Volume Section */
.volume-section {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.volume-primary {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.volume-value {
  font-size: 28px;
  font-weight: 600;
  color: #37352F;
  line-height: 1;
}

.volume-unit {
  font-size: 14px;
  color: #787774;
}

.volume-secondary {
  display: flex;
  align-items: baseline;
  gap: 4px;
  color: #9B9A97;
  font-size: 14px;
}

.drops-value {
  font-weight: 500;
  color: #787774;
}

.stock-badge {
  margin-left: auto;
}

/* Info Rows */
.info-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #787774;
}

.info-icon {
  font-size: 16px;
}

.info-text {
  color: #787774;
}

.info-label {
  color: #9B9A97;
}

.info-value {
  font-weight: 500;
  color: #37352F;
}

/* Tags Section */
.tags-section {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: auto;
  padding-top: 8px;
}

.tag-more {
  background: #E9E9E7 !important;
  color: #787774 !important;
}

/* Actions */
.card-actions {
  display: flex;
  gap: 8px;
  padding-top: 16px;
  margin-top: 8px;
  border-top: 1px solid #E9E9E7;
}

.action-btn {
  opacity: 0;
  transition: opacity 200ms ease, background 200ms ease;
}

.accord-card:hover .action-btn {
  opacity: 1;
}

.action-icon {
  font-size: 16px;
}

.action-danger:hover {
  color: #991B1B !important;
}
</style>
