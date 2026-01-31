<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <div id="app">
            <header v-if="authStore.isAuthenticated">
              <nav>
                <h1>🌸 Scentora</h1>
                <div class="nav-links">
                  <router-link to="/">Collection</router-link>
                  <router-link to="/statistics">Statistics</router-link>
                  <router-link to="/about">About</router-link>
                  <div class="user-menu">
                    <button @click="showUserMenu = !showUserMenu" class="user-button">
                      {{ authStore.user?.username || 'User' }} ▾
                    </button>
                    <div v-if="showUserMenu" class="user-dropdown">
                      <div class="user-info">
                        <strong>{{ authStore.user?.username }}</strong>
                        <small>{{ authStore.user?.email }}</small>
                      </div>
                      <button @click="handleLogout" class="logout-button">Logout</button>
                    </div>
                  </div>
                </div>
              </nav>
            </header>
            <main>
              <router-view />
            </main>
          </div>
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from './stores/auth';
import { naiveTheme } from './design-system/theme';

const authStore = useAuthStore();
const router = useRouter();
const showUserMenu = ref(false);
const themeOverrides = naiveTheme;

function handleLogout() {
  authStore.logout();
  showUserMenu.value = false;
  router.push('/login');
}

// Close dropdown when clicking outside
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement;
  if (!target.closest('.user-menu')) {
    showUserMenu.value = false;
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

header {
  background: white;
  border-bottom: 1px solid #E9E9E7;
  padding: 1rem 2rem;
}

nav {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

h1 {
  font-size: 1.5rem;
  color: #37352F;
  font-weight: 600;
}

.nav-links {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.nav-links a {
  text-decoration: none;
  color: #787774;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  color: #0F766E;
}

.user-menu {
  position: relative;
}

.user-button {
  background: transparent;
  border: 1px solid #E9E9E7;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #787774;
  transition: all 0.2s;
}

.user-button:hover {
  border-color: #0F766E;
  color: #0F766E;
  background: #FAFAFA;
}

.user-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background: white;
  border: 1px solid #E9E9E7;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  min-width: 200px;
  z-index: 1000;
}

.user-info {
  padding: 16px;
  border-bottom: 1px solid #E9E9E7;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-info strong {
  color: #37352F;
  font-size: 14px;
}

.user-info small {
  color: #9B9A97;
  font-size: 12px;
}

.logout-button {
  width: 100%;
  padding: 12px 16px;
  background: transparent;
  border: none;
  color: #991B1B;
  cursor: pointer;
  text-align: left;
  font-size: 14px;
  transition: background 0.2s;
}

.logout-button:hover {
  background: #FEE2E2;
}

main {
  flex: 1;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding: 2rem;
}
</style>
