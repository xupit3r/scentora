<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <!-- Login page (no sidebar) -->
          <div v-if="!authStore.isAuthenticated" id="app" class="app-auth">
            <router-view />
          </div>

          <!-- Main app with sidebar -->
          <n-layout v-else has-sider class="app-layout">
            <app-sidebar />
            
            <n-layout-content class="main-content" :native-scrollbar="false">
              <div class="content-wrapper">
                <app-breadcrumbs />
                <router-view />
              </div>
            </n-layout-content>
          </n-layout>
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { useAuthStore } from './stores/auth';
import { naiveTheme } from './design-system/theme';
import AppSidebar from './components/layout/AppSidebar.vue';
import AppBreadcrumbs from './components/layout/AppBreadcrumbs.vue';

const authStore = useAuthStore();
const themeOverrides = naiveTheme;
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
}

.main-content {
  margin-left: 280px;
  min-height: 100vh;
  background: #FAFAFA;
  transition: margin-left 0.3s ease;
}

.content-wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px 48px;
}

/* Responsive */
@media (max-width: 1024px) {
  .content-wrapper {
    padding: 24px 32px;
  }
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }

  .content-wrapper {
    padding: 16px 20px;
  }
}
</style>
