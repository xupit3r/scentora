<template>
  <div class="empty-state">
    <div class="empty-icon">{{ icon }}</div>
    <h3 class="empty-title">{{ title }}</h3>
    <p class="empty-description">{{ description }}</p>
    <n-button
      v-if="actionLabel"
      type="primary"
      @click="$emit('action')"
      class="empty-action"
    >
      {{ actionLabel }}
    </n-button>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  icon?: string;
  title: string;
  description: string;
  actionLabel?: string;
}>(), {
  icon: '📦',
});

defineEmits<{
  action: [];
}>();
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  text-align: center;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 24px;
  opacity: 0.6;
}

.empty-title {
  font-size: 20px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 12px 0;
}

.empty-description {
  font-size: 14px;
  color: #787774;
  max-width: 400px;
  margin: 0 0 24px 0;
  line-height: 1.6;
}

.empty-action {
  margin-top: 8px;
}
</style>
