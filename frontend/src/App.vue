<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <!-- Login page (no sidebar) -->
          <div v-if="!authStore.isAuthenticated" id="app" class="app-auth">
            <router-view />
          </div>

          <!-- Main app with hamburger sidebar -->
          <div v-else class="app-layout">
            <div class="top-bar">
              <button class="hamburger-btn" @click="sidebarOpen = true" aria-label="Open menu">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <rect x="3" y="6" width="18" height="2" rx="1" fill="#37352F"/>
                  <rect x="3" y="11" width="18" height="2" rx="1" fill="#37352F"/>
                  <rect x="3" y="16" width="18" height="2" rx="1" fill="#37352F"/>
                </svg>
              </button>
              <span class="brand">🌸 Scentora</span>
            </div>

            <app-sidebar :open="sidebarOpen" @close="sidebarOpen = false" />

            <div class="main-content">
              <div class="content-wrapper">
                <app-breadcrumbs />
                <router-view />
              </div>
            </div>
          </div>
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from './stores/auth';
import { naiveTheme } from './design-system/theme';
import AppSidebar from './components/layout/AppSidebar.vue';
import AppBreadcrumbs from './components/layout/AppBreadcrumbs.vue';

const authStore = useAuthStore();
const themeOverrides = naiveTheme;
const sidebarOpen = ref(false);
</script>

<style scoped>
.app-auth {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #FAFAFA;
}

.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.top-bar {
  position: sticky;
  top: 0;
  z-index: 100;
  height: 56px;
  background: #FFFFFF;
  border-bottom: 1px solid #E9E9E7;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 16px;
}

.hamburger-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  transition: background 150ms ease;
  flex-shrink: 0;
}

.hamburger-btn:hover {
  background: #F7F6F3;
}

.brand {
  font-size: 18px;
  font-weight: 600;
  color: #37352F;
}

.main-content {
  flex: 1;
  background: #FAFAFA;
}

.content-wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px 48px;
}

@media (max-width: 1024px) {
  .content-wrapper {
    padding: 24px 32px;
  }
}

@media (max-width: 768px) {
  .content-wrapper {
    padding: 16px 20px;
  }
}
</style>
