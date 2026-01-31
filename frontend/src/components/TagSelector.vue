<template>
  <div class="tag-selector">
    <label v-if="label" class="label">{{ label }}</label>
    
    <div class="selected-tags">
      <span
        v-for="tag in modelValue"
        :key="tag"
        class="selected-tag"
      >
        {{ tag }}
        <button @click="removeTag(tag)" type="button" class="remove-btn">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </span>
    </div>

    <div class="input-wrapper">
      <input
        v-model="searchQuery"
        @input="handleSearch"
        @focus="showDropdown = true"
        @blur="handleBlur"
        @keydown.enter.prevent="handleEnter"
        @keydown.down.prevent="navigateDown"
        @keydown.up.prevent="navigateUp"
        type="text"
        :placeholder="placeholder"
        class="tag-input"
      />
      <div v-if="searchQuery && showDropdown" class="dropdown">
        <div
          v-for="(tag, index) in filteredTags"
          :key="tag.tag"
          @mousedown.prevent="selectTag(tag.tag)"
          @mouseover="highlightedIndex = index"
          class="dropdown-item"
          :class="{ highlighted: highlightedIndex === index }"
        >
          <span class="tag-name">{{ tag.tag }}</span>
          <span class="tag-category">{{ tag.category }}</span>
        </div>
        <div v-if="filteredTags.length === 0 && searchQuery.length >= 2" class="dropdown-item empty">
          No tags found. Press Enter to add "{{ searchQuery }}"
        </div>
      </div>
    </div>

    <div v-if="showSuggestions && !searchQuery" class="suggestions">
      <button
        v-for="tag in popularTags"
        :key="tag"
        @click="selectTag(tag)"
        type="button"
        class="suggestion-tag"
      >
        {{ tag }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { tagService } from '@/services/api';
import type { PredefinedTag } from '@/types';

const props = withDefaults(defineProps<{
  modelValue: string[];
  label?: string;
  placeholder?: string;
  showSuggestions?: boolean;
}>(), {
  placeholder: 'Search and add tags...',
  showSuggestions: true,
});

const emit = defineEmits<{
  'update:modelValue': [value: string[]];
}>();

const searchQuery = ref('');
const showDropdown = ref(false);
const allTags = ref<PredefinedTag[]>([]);
const highlightedIndex = ref(0);

const filteredTags = computed(() => {
  if (!searchQuery.value || searchQuery.value.length < 2) return [];
  
  const query = searchQuery.value.toLowerCase();
  return allTags.value
    .filter(tag => 
      tag.tag.toLowerCase().includes(query) &&
      !props.modelValue.includes(tag.tag)
    )
    .slice(0, 10);
});

const popularTags = computed(() => {
  // Return some popular tags from different categories
  const popular = ['fresh', 'warm', 'woody', 'floral', 'spicy', 'citrus'];
  return popular.filter(tag => !props.modelValue.includes(tag));
});

onMounted(async () => {
  try {
    allTags.value = await tagService.getAll();
  } catch (error) {
    console.error('Failed to load tags:', error);
  }
});

watch(searchQuery, () => {
  highlightedIndex.value = 0;
});

function handleSearch() {
  showDropdown.value = true;
}

function handleBlur() {
  setTimeout(() => {
    showDropdown.value = false;
  }, 200);
}

function handleEnter() {
  if (highlightedIndex.value >= 0 && filteredTags.value[highlightedIndex.value]) {
    selectTag(filteredTags.value[highlightedIndex.value].tag);
  } else if (searchQuery.value.trim().length >= 2) {
    // Add custom tag
    selectTag(searchQuery.value.trim().toLowerCase());
  }
}

function navigateDown() {
  if (highlightedIndex.value < filteredTags.value.length - 1) {
    highlightedIndex.value++;
  }
}

function navigateUp() {
  if (highlightedIndex.value > 0) {
    highlightedIndex.value--;
  }
}

function selectTag(tag: string) {
  if (!props.modelValue.includes(tag)) {
    emit('update:modelValue', [...props.modelValue, tag]);
  }
  searchQuery.value = '';
  showDropdown.value = false;
}

function removeTag(tag: string) {
  emit('update:modelValue', props.modelValue.filter(t => t !== tag));
}
</script>

<style scoped>
.tag-selector {
  width: 100%;
}

.label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #37352F;
  margin-bottom: 8px;
}

.selected-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
  min-height: 32px;
}

.selected-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  background: #0F766E;
  color: white;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
}

.remove-btn {
  padding: 0;
  background: none;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  opacity: 0.8;
  transition: opacity 200ms;
}

.remove-btn:hover {
  opacity: 1;
}

.remove-btn svg {
  width: 14px;
  height: 14px;
}

.input-wrapper {
  position: relative;
}

.tag-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #E9E9E7;
  border-radius: 8px;
  font-size: 14px;
  transition: all 200ms;
  color: #37352F;
}

.tag-input:focus {
  outline: none;
  border-color: #0F766E;
  box-shadow: 0 0 0 3px rgba(15, 118, 110, 0.1);
}

.tag-input::placeholder {
  color: #9B9A97;
}

.dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: white;
  border: 1px solid #E9E9E7;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.06);
  max-height: 256px;
  overflow-y: auto;
  z-index: 50;
}

.dropdown-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  cursor: pointer;
  transition: background-color 200ms;
}

.dropdown-item:hover,
.dropdown-item.highlighted {
  background-color: #F7F6F3;
}

.dropdown-item.empty {
  color: #787774;
  cursor: default;
}

.dropdown-item.empty:hover {
  background-color: white;
}

.tag-name {
  font-weight: 500;
  color: #37352F;
  font-size: 14px;
}

.tag-category {
  font-size: 12px;
  color: #9B9A97;
  text-transform: capitalize;
}

.suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.suggestion-tag {
  padding: 6px 12px;
  background: #FAFAFA;
  border: 1px solid #E9E9E7;
  border-radius: 6px;
  font-size: 14px;
  color: #787774;
  cursor: pointer;
  transition: all 200ms;
}

.suggestion-tag:hover {
  background: #F7F6F3;
  border-color: #D9D9D7;
  color: #37352F;
}
</style>
