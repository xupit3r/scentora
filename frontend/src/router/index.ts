import { createRouter, createWebHistory } from 'vue-router';
import Home from '@/views/Home.vue';
import About from '@/views/About.vue';
import PerfumeDetail from '@/views/PerfumeDetail.vue';

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
      path: '/about',
      name: 'about',
      component: About,
    },
  ],
});

export default router;
