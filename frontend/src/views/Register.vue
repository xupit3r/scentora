<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1>🌸</h1>
        <h2>Join Scentora</h2>
        <p class="subtitle">Create your accord inventory account</p>
      </div>

      <n-form @submit.prevent="handleRegister" class="auth-form">
        <n-form-item label="Invitation Code" required>
          <n-input
            v-model:value="invitationCode"
            type="text"
            placeholder="Enter your invitation code"
            size="large"
          />
        </n-form-item>

        <n-form-item label="Username" required>
          <n-input
            v-model:value="username"
            type="text"
            placeholder="Your display name"
            size="large"
          />
        </n-form-item>

        <n-form-item label="Email" required>
          <n-input
            v-model:value="email"
            type="email"
            placeholder="your@email.com"
            size="large"
          />
        </n-form-item>

        <n-form-item label="Password" required>
          <n-input
            v-model:value="password"
            type="password"
            placeholder="••••••••"
            size="large"
            show-password-on="click"
          />
        </n-form-item>

        <n-form-item label="Confirm Password" required>
          <n-input
            v-model:value="confirmPassword"
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
          {{ authStore.isLoading ? 'Creating Account...' : 'Create Account' }}
        </n-button>

        <div v-if="authStore.error" class="error-message">
          {{ authStore.error }}
        </div>
      </n-form>

      <div class="auth-footer">
        <p class="info-text">
          Already have an account?
          <router-link to="/login" class="link">Login</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const invitationCode = ref('');
const username = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');

async function handleRegister() {
  if (password.value !== confirmPassword.value) {
    authStore.error = 'Passwords do not match';
    return;
  }

  if (password.value.length < 4) {
    authStore.error = 'Password must be at least 4 characters';
    return;
  }

  const success = await authStore.register(
    email.value,
    username.value,
    password.value,
    invitationCode.value
  );

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

.auth-footer {
  text-align: center;
  padding-top: 24px;
  border-top: 1px solid #E9E9E7;
}

.info-text {
  color: #787774;
  font-size: 14px;
  margin: 0;
}

.link {
  color: #0F766E;
  text-decoration: none;
  font-weight: 600;
  margin-left: 4px;
}

.link:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .auth-card {
    padding: 32px 24px;
  }

  .auth-header h1 {
    font-size: 40px;
  }

  .auth-header h2 {
    font-size: 20px;
  }
}
</style>
