<template>
  <div class="ingredient-list">
    <div v-if="ingredients.length === 0" class="empty-ingredients">
      <span class="empty-icon">🧪</span>
      <span class="empty-text">No ingredients added yet</span>
      <n-button
        type="primary"
        size="small"
        @click="$emit('add')"
      >
        Add First Ingredient
      </n-button>
    </div>

    <div v-else class="ingredients-table">
      <n-table :bordered="false" :single-line="false" striped>
        <thead>
          <tr>
            <th>Accord</th>
            <th class="text-right">Quantity</th>
            <th class="text-right">Percentage</th>
            <th v-if="showActions" class="text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="ingredient in ingredients"
            :key="ingredient._id"
            class="ingredient-row"
          >
            <td>
              <div class="ingredient-name">
                {{ ingredient.accordName || 'Unknown Accord' }}
              </div>
              <div v-if="ingredient.notes" class="ingredient-notes">
                {{ ingredient.notes }}
              </div>
            </td>
            <td class="text-right quantity-cell">
              <div class="quantity-display">
                <span class="quantity-ml">{{ ingredient.quantityMl }} ml</span>
                <span class="quantity-drops">({{ ingredient.quantityDrops }} drops)</span>
              </div>
            </td>
            <td class="text-right percentage-cell">
              <span v-if="ingredient.percentage" class="percentage-value">
                {{ ingredient.percentage.toFixed(1) }}%
              </span>
              <span v-else class="percentage-na">-</span>
            </td>
            <td v-if="showActions" class="text-center actions-cell">
              <n-button
                text
                size="small"
                @click="$emit('edit', ingredient)"
                title="Edit ingredient"
              >
                <span class="action-icon">✏️</span>
              </n-button>
              <n-button
                text
                size="small"
                @click="$emit('remove', ingredient)"
                class="action-danger"
                title="Remove ingredient"
              >
                <span class="action-icon">🗑️</span>
              </n-button>
            </td>
          </tr>
        </tbody>
        <tfoot v-if="showTotals">
          <tr class="totals-row">
            <td><strong>Total</strong></td>
            <td class="text-right">
              <strong>{{ totalVolume.toFixed(1) }} ml</strong>
            </td>
            <td class="text-right">
              <strong>{{ totalPercentage.toFixed(1) }}%</strong>
            </td>
            <td v-if="showActions"></td>
          </tr>
        </tfoot>
      </n-table>

      <div v-if="totalPercentage !== 100 && showTotals" class="percentage-warning">
        <n-alert
          :type="Math.abs(100 - totalPercentage) > 1 ? 'warning' : 'info'"
          :show-icon="false"
        >
          <template #header>
            {{ totalPercentage > 100 ? '⚠️ Total exceeds 100%' : 'ℹ️ Total less than 100%' }}
          </template>
          Current total: {{ totalPercentage.toFixed(1) }}%
        </n-alert>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { NTable, NButton, NAlert } from 'naive-ui';
import type { RecipeIngredient } from '@/types';

interface Props {
  ingredients: RecipeIngredient[];
  showActions?: boolean;
  showTotals?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showActions: true,
  showTotals: true,
});

defineEmits<{
  add: [];
  edit: [ingredient: RecipeIngredient];
  remove: [ingredient: RecipeIngredient];
}>();

// Calculate totals
const totalVolume = computed(() => {
  return props.ingredients.reduce((sum, ing) => sum + ing.quantityMl, 0);
});

const totalPercentage = computed(() => {
  return props.ingredients.reduce((sum, ing) => sum + (ing.percentage || 0), 0);
});
</script>

<style scoped>
.ingredient-list {
  width: 100%;
}

.empty-ingredients {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 48px 24px;
  color: var(--text-color-3);
}

.empty-icon {
  font-size: 48px;
  opacity: 0.6;
}

.empty-text {
  font-size: 14px;
}

.ingredients-table {
  width: 100%;
}

.ingredient-row {
  transition: background-color 0.15s ease;
}

.ingredient-row:hover {
  background-color: var(--hover-color);
}

.ingredient-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color-1);
}

.ingredient-notes {
  font-size: 12px;
  color: var(--text-color-3);
  font-style: italic;
  margin-top: 4px;
}

.quantity-cell {
  font-size: 14px;
}

.quantity-display {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.quantity-ml {
  font-weight: 500;
  color: var(--text-color-1);
}

.quantity-drops {
  font-size: 12px;
  color: var(--text-color-3);
}

.percentage-cell {
  font-size: 14px;
}

.percentage-value {
  font-weight: 500;
  color: var(--text-color-1);
}

.percentage-na {
  color: var(--text-color-3);
}

.actions-cell {
  white-space: nowrap;
}

.action-icon {
  font-size: 16px;
}

.action-danger:hover {
  color: #d03050;
}

.totals-row {
  background: var(--table-header-color);
  font-weight: 600;
}

.percentage-warning {
  margin-top: 16px;
}

.text-right {
  text-align: right;
}

.text-center {
  text-align: center;
}

/* Responsive */
@media (max-width: 768px) {
  .ingredients-table {
    overflow-x: auto;
  }

  .quantity-drops {
    display: none;
  }
}
</style>
