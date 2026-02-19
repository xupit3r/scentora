<template>
  <n-drawer :show="open" :width="280" placement="left" @update:show="(val: boolean) => { if (!val) emit('close') }">
    <div class="sidebar-content">
      <!-- Logo / Header -->
      <div class="sidebar-header">
        <div class="logo">
          <span class="logo-icon">🌸</span>
          <span class="logo-text">Scentora</span>
        </div>
        <button class="close-btn" @click="emit('close')" aria-label="Close menu">✕</button>
      </div>

      <!-- Navigation Menu -->
      <n-menu
        v-model:value="activeKey"
        :options="menuOptions"
        :indent="24"
        @update:value="handleMenuSelect"
      />

      <!-- User Profile at Bottom -->
      <div class="sidebar-footer">
        <div class="user-profile" @click="toggleUserMenu">
          <div class="user-avatar">
            {{ userInitial }}
          </div>
          <div class="user-info">
            <div class="user-name">{{ authStore.user?.username }}</div>
            <div class="user-email">{{ authStore.user?.email }}</div>
          </div>
        </div>

        <!-- User Menu Dropdown -->
        <transition name="slide-up">
          <div v-if="showUserMenu" class="user-menu-dropdown">
            <div class="menu-item" @click="handleLogout">
              <span class="menu-icon">🚪</span>
              <span>Logout</span>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </n-drawer>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { MenuOption } from 'naive-ui';
import { useAuthStore } from '@/stores/auth';

const props = defineProps<{ open: boolean }>();
const emit = defineEmits<{ (e: 'close'): void }>();

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const showUserMenu = ref(false);
const activeKey = ref<string>(route.path);

// Update active key when route changes, and close drawer on navigation
router.afterEach((to) => {
  activeKey.value = to.path;
  emit('close');
});

const userInitial = computed(() => {
  return authStore.user?.username?.charAt(0).toUpperCase() || 'U';
});

function renderIcon(icon: string) {
  return () => h('span', { style: 'font-size: 18px;' }, icon);
}

const menuOptions = computed<MenuOption[]>(() => [
  {
    label: 'Collection',
    key: '/',
    icon: renderIcon('🏺'),
  },
]);

function handleMenuSelect(key: string) {
  router.push(key);
}

function toggleUserMenu() {
  showUserMenu.value = !showUserMenu.value;
}

function handleLogout() {
  authStore.logout();
  showUserMenu.value = false;
  router.push('/login');
}

// Close user menu when clicking outside
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement;
  if (!target.closest('.sidebar-footer')) {
    showUserMenu.value = false;
  }
}

if (typeof window !== 'undefined') {
  document.addEventListener('click', handleClickOutside);
}
</script>

<style scoped>
.sidebar-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #FFFFFF;
}

/* Logo / Header */
.sidebar-header {
  padding: 16px 20px;
  border-bottom: 1px solid #E9E9E7;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 20px;
  font-weight: 600;
  color: #37352F;
}

.logo-icon {
  font-size: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  white-space: nowrap;
}

.close-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 18px;
  color: #787774;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background 150ms ease;
  line-height: 1;
}

.close-btn:hover {
  background: #F7F6F3;
  color: #37352F;
}

/* User Profile Footer */
.sidebar-footer {
  margin-top: auto;
  padding: 16px;
  border-top: 1px solid #E9E9E7;
  position: relative;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 200ms ease;
}

.user-profile:hover {
  background: #F7F6F3;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: linear-gradient(135deg, #0F766E 0%, #14B8A6 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
}

.user-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #37352F;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-email {
  font-size: 12px;
  color: #9B9A97;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* User Menu Dropdown */
.user-menu-dropdown {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 16px;
  right: 16px;
  background: white;
  border: 1px solid #E9E9E7;
  border-radius: 8px;
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  z-index: 1000;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  font-size: 14px;
  color: #991B1B;
  cursor: pointer;
  transition: background 200ms ease;
}

.menu-item:hover {
  background: #FEE2E2;
}

.menu-icon {
  font-size: 16px;
}

/* Transitions */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 200ms ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
