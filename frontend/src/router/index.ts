import { createRouter, createWebHistory } from 'vue-router';
import Home from '@/views/Home.vue';
import About from '@/views/About.vue';
import PerfumeDetail from '@/views/PerfumeDetail.vue';
import Statistics from '@/views/Statistics.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/perfume/:id',
      name: 'perfume-detail',
      component: PerfumeDetail,
    },
    {
      path: '/statistics',
      name: 'statistics',
      component: Statistics,
    },
    {
      path: '/about',
      name: 'about',
      component: About,
    },
  ],
});

export default router;
