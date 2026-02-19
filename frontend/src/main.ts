import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import { useAuthStore } from './stores/auth';

// Import styles
import './style.css';

// Naive UI setup (global configuration)
import {
  create,
  NButton,
  NInput,
  NCard,
  NSpace,
  NGrid,
  NGridItem,
  NTag,
  NSelect,
  NModal,
  NForm,
  NFormItem,
  NCheckbox,
  NMenu,
  NIcon,
  NLayout,
  NLayoutSider,
  NLayoutHeader,
  NLayoutContent,
  NDrawer,
  NConfigProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
} from 'naive-ui';

const naive = create({
  components: [
    NButton,
    NInput,
    NCard,
    NSpace,
    NGrid,
    NGridItem,
    NTag,
    NSelect,
    NModal,
    NForm,
    NFormItem,
    NCheckbox,
    NMenu,
    NIcon,
    NLayout,
    NLayoutSider,
    NLayoutHeader,
    NLayoutContent,
    NDrawer,
    NConfigProvider,
    NMessageProvider,
    NNotificationProvider,
    NDialogProvider,
  ],
});

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);
app.use(naive);

// Initialize auth store early to set up axios interceptors
useAuthStore();

app.mount('#app');
