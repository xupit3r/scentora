<template>
  <div class="perfume-pyramid">
    <div class="pyramid-level top">
      <h4>Top Notes</h4>
      <div class="notes">
        <span v-for="note in pyramid.top" :key="note" class="note">{{ note }}</span>
        <span v-if="pyramid.top.length === 0" class="empty">No notes</span>
      </div>
    </div>
    <div class="pyramid-level middle">
      <h4>Middle Notes (Heart)</h4>
      <div class="notes">
        <span v-for="note in pyramid.middle" :key="note" class="note">{{ note }}</span>
        <span v-if="pyramid.middle.length === 0" class="empty">No notes</span>
      </div>
    </div>
    <div class="pyramid-level base">
      <h4>Base Notes</h4>
      <div class="notes">
        <span v-for="note in pyramid.base" :key="note" class="note">{{ note }}</span>
        <span v-if="pyramid.base.length === 0" class="empty">No notes</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PerfumePyramid } from '@/types';

interface Props {
  pyramid: PerfumePyramid;
}

defineProps<Props>();
</script>

<style scoped>
.perfume-pyramid {
  width: 100%;
}

.pyramid-level {
  padding: 1rem;
  margin-bottom: 0.5rem;
  border-radius: 8px;
}

.pyramid-level.top {
  background: linear-gradient(135deg, #fff9e6 0%, #fff3cc 100%);
}

.pyramid-level.middle {
  background: linear-gradient(135deg, #ffe6f0 0%, #ffcce0 100%);
}

.pyramid-level.base {
  background: linear-gradient(135deg, #e6f0ff 0%, #cce0ff 100%);
}

h4 {
  font-size: 0.9rem;
  font-weight: 600;
  color: #555;
  margin-bottom: 0.5rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.notes {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.note {
  display: inline-block;
  padding: 0.3rem 0.8rem;
  background: white;
  border-radius: 15px;
  font-size: 0.85rem;
  color: #666;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.empty {
  color: #999;
  font-style: italic;
  font-size: 0.85rem;
}
</style>
