import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import Home from '@/views/Home.vue';
import About from '@/views/About.vue';
import Statistics from '@/views/Statistics.vue';
import Login from '@/views/Login.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { public: true },
    },
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: { requiresAuth: true },
    },
    {
      path: '/statistics',
      name: 'statistics',
      component: Statistics,
      meta: { requiresAuth: true },
    },
    {
      path: '/about',
      name: 'about',
      component: About,
      meta: { requiresAuth: true },
    },
  ],
});

// Navigation guard to protect routes
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore();

  // Check if route requires auth
  if (to.meta.requiresAuth) {
    // If no token, redirect to login
    if (!authStore.accessToken) {
      next({ name: 'login', query: { redirect: to.fullPath } });
      return;
    }

    // If we have a token but no user, fetch user data
    if (!authStore.user) {
      const success = await authStore.fetchCurrentUser();
      if (!success) {
        next({ name: 'login', query: { redirect: to.fullPath } });
        return;
      }
    }
  }

  // If logged in and trying to access auth pages, redirect to home
  if (to.meta.public && authStore.isAuthenticated) {
    next({ name: 'home' });
    return;
  }

  next();
});

export default router;
