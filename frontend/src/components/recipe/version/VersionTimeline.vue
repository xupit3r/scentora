<template>
  <div class="version-timeline">
    <div v-if="sortedVersions.length === 0" class="empty-timeline">
      <span class="empty-icon">🕒</span>
      <span class="empty-text">No version history yet</span>
    </div>

    <n-timeline v-else>
      <n-timeline-item
        v-for="version in sortedVersions"
        :key="version._id"
        :type="version.isActive ? 'success' : 'default'"
        :title="versionTitle(version)"
        :time="versionTime(version)"
      >
        <div class="timeline-content">
          <div class="version-info">
            <n-tag
              v-if="version.isActive"
              type="success"
              size="small"
              round
            >
              Active
            </n-tag>
            <span class="version-name">{{ version.name }}</span>
          </div>

          <p v-if="version.notes" class="version-notes">
            {{ version.notes }}
          </p>

          <div class="version-actions">
            <n-button
              text
              size="small"
              @click="$emit('view', version)"
            >
              View
            </n-button>
            <n-button
              v-if="!version.isActive"
              text
              size="small"
              type="primary"
              @click="$emit('activate', version)"
            >
              Activate
            </n-button>
            <n-button
              text
              size="small"
              @click="$emit('duplicate', version)"
            >
              Duplicate
            </n-button>
          </div>
        </div>
      </n-timeline-item>
    </n-timeline>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { NTimeline, NTimelineItem, NTag, NButton } from 'naive-ui';
import type { RecipeVersion } from '@/types';

interface Props {
  versions: RecipeVersion[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  view: [version: RecipeVersion];
  activate: [version: RecipeVersion];
  duplicate: [version: RecipeVersion];
}>();

// Sort versions by version number (newest first)
const sortedVersions = computed(() => {
  return [...props.versions].sort((a, b) => b.versionNumber - a.versionNumber);
});

function versionTitle(version: RecipeVersion): string {
  return `Version ${version.versionNumber}`;
}

function versionTime(version: RecipeVersion): string {
  const date = new Date(version.createdAt);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays === 0) {
    return 'Today';
  } else if (diffDays === 1) {
    return 'Yesterday';
  } else if (diffDays < 7) {
    return `${diffDays} days ago`;
  } else if (diffDays < 30) {
    const weeks = Math.floor(diffDays / 7);
    return `${weeks} week${weeks > 1 ? 's' : ''} ago`;
  } else {
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
    });
  }
}
</script>

<style scoped>
.version-timeline {
  padding: 16px 0;
}

.empty-timeline {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
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

.timeline-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-top: 4px;
}

.version-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.version-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color-2);
}

.version-notes {
  font-size: 13px;
  color: var(--text-color-3);
  line-height: 1.6;
  margin: 0;
  padding: 8px 12px;
  background: var(--code-color);
  border-radius: 6px;
}

.version-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

/* Responsive */
@media (max-width: 640px) {
  .version-timeline {
    padding: 8px 0;
  }

  .version-actions {
    gap: 4px;
  }
}
</style>
