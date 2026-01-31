<template>
  <div class="breadcrumbs">
    <span
      v-for="(crumb, index) in crumbs"
      :key="index"
      class="breadcrumb-item"
    >
      <router-link
        v-if="crumb.path && index < crumbs.length - 1"
        :to="crumb.path"
        class="breadcrumb-link"
      >
        {{ crumb.label }}
      </router-link>
      <span v-else class="breadcrumb-current">
        {{ crumb.label }}
      </span>
      <span v-if="index < crumbs.length - 1" class="breadcrumb-separator">/</span>
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

interface Breadcrumb {
  label: string;
  path?: string;
}

const route = useRoute();

const crumbs = computed<Breadcrumb[]>(() => {
  const matched = route.matched;
  const breadcrumbs: Breadcrumb[] = [];

  // Map routes to breadcrumbs
  for (const record of matched) {
    if (record.meta?.breadcrumb) {
      breadcrumbs.push({
        label: record.meta.breadcrumb as string,
        path: record.path,
      });
    } else if (record.name) {
      // Fallback to route name
      const label = String(record.name)
        .split('-')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');
      breadcrumbs.push({
        label,
        path: record.path,
      });
    }
  }

  // If no breadcrumbs from routes, create from path
  if (breadcrumbs.length === 0) {
    const pathSegments = route.path.split('/').filter(Boolean);
    
    if (pathSegments.length === 0) {
      breadcrumbs.push({ label: 'Collection', path: '/' });
    } else {
      breadcrumbs.push({ label: 'Home', path: '/' });
      
      let currentPath = '';
      pathSegments.forEach((segment, index) => {
        currentPath += `/${segment}`;
        const label = segment
          .split('-')
          .map(word => word.charAt(0).toUpperCase() + word.slice(1))
          .join(' ');
        
        breadcrumbs.push({
          label,
          path: index < pathSegments.length - 1 ? currentPath : undefined,
        });
      });
    }
  }

  return breadcrumbs;
});
</script>

<style scoped>
.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #787774;
  padding: 16px 0;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.breadcrumb-link {
  color: #787774;
  text-decoration: none;
  transition: color 200ms ease;
}

.breadcrumb-link:hover {
  color: #0F766E;
}

.breadcrumb-current {
  color: #37352F;
  font-weight: 500;
}

.breadcrumb-separator {
  color: #D9D9D7;
}
</style>
