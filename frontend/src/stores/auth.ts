import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from 'axios';

interface User {
  id: string;
  email: string;
  username: string;
}

interface LoginResponse {
  accessToken: string;
}

const ACCESS_TOKEN_KEY = 'scentora_access_token';
const API_BASE = 'http://localhost:3000/api';

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(localStorage.getItem(ACCESS_TOKEN_KEY));
  const user = ref<User | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const isAuthenticated = computed(() => !!accessToken.value && !!user.value);

  function setAuthHeader(token: string | null) {
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
      delete axios.defaults.headers.common['Authorization'];
    }
  }

  function storeToken(token: string) {
    accessToken.value = token;
    localStorage.setItem(ACCESS_TOKEN_KEY, token);
    setAuthHeader(token);
  }

  function clearTokens() {
    accessToken.value = null;
    user.value = null;
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    setAuthHeader(null);
  }

  // Initialize auth header if token exists
  if (accessToken.value) {
    setAuthHeader(accessToken.value);
  }

  // Redirect to login on 401
  axios.interceptors.response.use(
    (response) => response,
    async (err) => {
      if (err.response?.status === 401 && !err.config.url?.includes('/auth/login')) {
        clearTokens();
        window.location.href = '/login';
      }
      return Promise.reject(err);
    }
  );

  async function login(password: string) {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await axios.post<LoginResponse>(`${API_BASE}/auth/login`, {
        password,
      });

      storeToken(response.data.accessToken);

      // Fetch user info after login
      await fetchCurrentUser();

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Login failed';
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchCurrentUser() {
    if (!accessToken.value) return false;

    isLoading.value = true;
    error.value = null;
    try {
      const response = await axios.get<{ user: User }>(`${API_BASE}/auth/me`);
      user.value = response.data.user;
      return true;
    } catch (err: any) {
      clearTokens();
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function logout() {
    clearTokens();
  }

  return {
    accessToken,
    user,
    isLoading,
    error,
    isAuthenticated,
    login,
    fetchCurrentUser,
    logout,
  };
});
