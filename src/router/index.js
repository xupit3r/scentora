import { createRouter, createWebHistory } from 'vue-router'
import PerfumeListView from '../views/PerfumeListView.vue';
import PerfumeView from '../views/PerfumeView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: PerfumeListView,
    },
    {
      path: '/perfume',
      name: 'perfume',
      component: PerfumeView,
    },
  ],
})

export default router
