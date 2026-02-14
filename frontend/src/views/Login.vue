<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1>🌸</h1>
        <h2>Welcome to Scentora</h2>
        <p class="subtitle">Enter the password to continue</p>
      </div>

      <n-form @submit.prevent="handleLogin" class="auth-form">
        <n-form-item label="Password">
          <n-input
            v-model:value="password"
            type="password"
            placeholder="••••••••"
            size="large"
            show-password-on="click"
          />
        </n-form-item>

        <n-button
          type="primary"
          attr-type="submit"
          :loading="authStore.isLoading"
          block
          size="large"
        >
          {{ authStore.isLoading ? 'Logging in...' : 'Login' }}
        </n-button>

        <div v-if="authStore.error" class="error-message">
          {{ authStore.error }}
        </div>
      </n-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const password = ref('');

async function handleLogin() {
  const success = await authStore.login(password.value);
  if (success) {
    router.push('/');
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #FAFAFA;
  padding: 20px;
}

.auth-card {
  background: white;
  border-radius: 16px;
  padding: 48px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.06);
  border: 1px solid #E9E9E7;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  font-size: 48px;
  margin: 0 0 16px 0;
}

.auth-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #787774;
  margin: 0;
  font-size: 14px;
}

.auth-form {
  margin-bottom: 24px;
}

.error-message {
  margin-top: 16px;
  padding: 12px;
  background: #FEE2E2;
  color: #991B1B;
  border-radius: 8px;
  font-size: 14px;
  text-align: center;
}
</style>
