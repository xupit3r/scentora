<template>
  <div class="recipe-list">
    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <skeleton-card v-for="i in 3" :key="i" />
    </div>

    <!-- Empty State -->
    <empty-state
      v-else-if="!loading && recipes.length === 0"
      icon="📝"
      :title="emptyTitle"
      :description="emptyDescription"
    >
      <template #action>
        <n-button type="primary" @click="$emit('create')">
          Create First Recipe
        </n-button>
      </template>
    </empty-state>

    <!-- Recipe Grid -->
    <div v-else class="recipe-grid">
      <recipe-card
        v-for="recipe in recipes"
        :key="recipe._id"
        :recipe="recipe"
        :version-count="getVersionCount(recipe._id)"
        :active-version-number="getActiveVersionNumber(recipe._id)"
        @view="$emit('view', $event)"
        @edit="$emit('edit', $event)"
        @duplicate="$emit('duplicate', $event)"
        @delete="$emit('delete', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { NButton } from 'naive-ui';
import RecipeCard from './RecipeCard.vue';
import EmptyState from '@/components/ui/EmptyState.vue';
import SkeletonCard from '@/components/ui/SkeletonCard.vue';
import type { Recipe } from '@/types';

interface Props {
  recipes: Recipe[];
  loading?: boolean;
  emptyTitle?: string;
  emptyDescription?: string;
  versionCounts?: Map<string, number>;
  activeVersionNumbers?: Map<string, number>;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyTitle: 'No recipes yet',
  emptyDescription: 'Create your first perfume recipe to get started',
  versionCounts: () => new Map(),
  activeVersionNumbers: () => new Map(),
});

defineEmits<{
  view: [recipe: Recipe];
  edit: [recipe: Recipe];
  duplicate: [recipe: Recipe];
  delete: [recipe: Recipe];
  create: [];
}>();

const getVersionCount = (recipeId: string): number => {
  return props.versionCounts.get(recipeId) || 0;
};

const getActiveVersionNumber = (recipeId: string): number => {
  return props.activeVersionNumbers.get(recipeId) || 0;
};
</script>

<style scoped>
.recipe-list {
  width: 100%;
}

.loading-state {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
  width: 100%;
}

.recipe-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
  width: 100%;
}

/* Responsive */
@media (max-width: 768px) {
  .recipe-grid,
  .loading-state {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}

@media (min-width: 1920px) {
  .recipe-grid,
  .loading-state {
    grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  }
}
</style>
