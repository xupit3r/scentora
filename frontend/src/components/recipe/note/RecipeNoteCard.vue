<template>
  <div class="note-card">
    <div class="note-header">
      <n-tag
        :type="noteTypeTag"
        size="small"
        round
      >
        {{ noteTypeLabel }}
      </n-tag>
      <span class="note-time">{{ formattedTime }}</span>
    </div>

    <div class="note-content">
      {{ note.content }}
    </div>

    <div v-if="showActions" class="note-actions">
      <n-button
        text
        size="small"
        @click="$emit('edit', note)"
      >
        Edit
      </n-button>
      <n-button
        text
        size="small"
        type="error"
        @click="$emit('delete', note)"
      >
        Delete
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { NTag, NButton } from 'naive-ui';
import type { RecipeNote, NoteType } from '@/types';

interface Props {
  note: RecipeNote;
  showActions?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showActions: true,
});

defineEmits<{
  edit: [note: RecipeNote];
  delete: [note: RecipeNote];
}>();

const noteTypeLabels: Record<NoteType, string> = {
  general: 'General',
  observation: 'Observation',
};

const noteTypeTags: Record<NoteType, 'default' | 'info'> = {
  general: 'default',
  observation: 'info',
};

const noteTypeLabel = computed(() => noteTypeLabels[props.note.noteType]);
const noteTypeTag = computed(() => noteTypeTags[props.note.noteType]);

// Time formatting
const formattedTime = computed(() => {
  const date = new Date(props.note.createdAt);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 1) {
    return 'Just now';
  } else if (diffMins < 60) {
    return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`;
  } else if (diffHours < 24) {
    return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
  } else if (diffDays < 7) {
    return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
  } else {
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
    });
  }
});
</script>

<style scoped>
.note-card {
  padding: 16px;
  border-radius: 8px;
  background: var(--card-color);
  border: 1px solid var(--divider-color);
  transition: all 0.15s ease;
}

.note-card:hover {
  border-color: var(--primary-color);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.note-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.note-time {
  font-size: 12px;
  color: var(--text-color-3);
}

.note-content {
  font-size: 14px;
  color: var(--text-color-1);
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

.note-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--divider-color);
}
</style>
