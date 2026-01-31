<template>
  <div class="accord-card" :class="`position-${accord.pyramidPosition}`">
    <div class="card-header">
      <h3 class="accord-name">{{ accord.name }}</h3>
      <span class="position-badge" :class="`badge-${accord.pyramidPosition}`">
        {{ positionLabel }}
      </span>
    </div>

    <div class="card-body">
      <div class="volume-info">
        <div class="volume-primary">
          <span class="volume-value">{{ accord.volumeMl }}</span>
          <span class="volume-unit">ml</span>
        </div>
        <div class="volume-secondary">
          <span class="drops-value">{{ accord.volumeDrops }}</span>
          <span class="drops-unit">drops</span>
        </div>
        <div v-if="isLowStock" class="low-stock-badge" :class="stockLevel">
          {{ stockWarning }}
        </div>
      </div>

      <div v-if="accord.supplier" class="supplier-info">
        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
        </svg>
        <span>{{ accord.supplier }}</span>
      </div>

      <div v-if="accord.dilutionPercentage" class="dilution-info">
        <span class="dilution-label">Dilution:</span>
        <span class="dilution-value">{{ accord.dilutionPercentage }}%</span>
      </div>

      <div v-if="accord.tags.length" class="tags">
        <span v-for="tag in accord.tags.slice(0, 3)" :key="tag" class="tag">
          {{ tag }}
        </span>
        <span v-if="accord.tags.length > 3" class="tag-more">
          +{{ accord.tags.length - 3 }}
        </span>
      </div>
    </div>

    <div class="card-actions">
      <button @click="$emit('view', accord)" class="btn-icon" title="View Details">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
      </button>
      <button @click="$emit('edit', accord)" class="btn-icon" title="Edit">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
        </svg>
      </button>
      <button @click="$emit('delete', accord)" class="btn-icon btn-danger" title="Delete">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </button>
    </div>
  </div>
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
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border-left: 4px solid;
}

.accord-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.accord-card.position-top {
  border-left-color: #FFD93D;
}

.accord-card.position-middle {
  border-left-color: #B565D8;
}

.accord-card.position-base {
  border-left-color: #A0826D;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 1rem;
  gap: 0.5rem;
}

.accord-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  flex: 1;
}

.position-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  white-space: nowrap;
}

.badge-top {
  background: linear-gradient(135deg, #FFD93D, #FFA800);
  color: #4a2800;
}

.badge-middle {
  background: linear-gradient(135deg, #B565D8, #8B5CF6);
  color: white;
}

.badge-base {
  background: linear-gradient(135deg, #A0826D, #6B4423);
  color: white;
}

.card-body {
  margin-bottom: 1rem;
}

.volume-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 0.75rem;
}

.volume-primary {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
}

.volume-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
}

.volume-unit {
  font-size: 0.875rem;
  color: #6b7280;
}

.volume-secondary {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
  color: #6b7280;
  font-size: 0.875rem;
}

.drops-value {
  font-weight: 600;
}

.low-stock-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  margin-left: auto;
}

.low-stock-badge.low {
  background: #FEF3C7;
  color: #92400E;
}

.low-stock-badge.critical {
  background: #FEE2E2;
  color: #991B1B;
}

.supplier-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #6b7280;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
}

.supplier-info .icon {
  width: 1rem;
  height: 1rem;
}

.dilution-info {
  display: flex;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: #6b7280;
  margin-bottom: 0.5rem;
}

.dilution-value {
  font-weight: 600;
  color: #1f2937;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.75rem;
}

.tag {
  padding: 0.25rem 0.625rem;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 0.75rem;
  color: #4b5563;
  font-weight: 500;
}

.tag-more {
  padding: 0.25rem 0.625rem;
  background: #e5e7eb;
  border-radius: 6px;
  font-size: 0.75rem;
  color: #6b7280;
  font-weight: 600;
}

.card-actions {
  display: flex;
  gap: 0.5rem;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
}

.btn-icon {
  padding: 0.5rem;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon svg {
  width: 1.25rem;
  height: 1.25rem;
  color: #6b7280;
}

.btn-icon:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
}

.btn-icon:hover svg {
  color: #1f2937;
}

.btn-icon.btn-danger:hover {
  background: #FEE2E2;
  border-color: #FCA5A5;
}

.btn-icon.btn-danger:hover svg {
  color: #DC2626;
}
</style>
