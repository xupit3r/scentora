import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import axios from 'axios';

interface User {
  id: string;
  email: string;
  username: string;
}

interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

interface RefreshResponse {
  accessToken: string;
  refreshToken: string;
}

const ACCESS_TOKEN_KEY = 'scentora_access_token';
const REFRESH_TOKEN_KEY = 'scentora_refresh_token';
const API_BASE = 'http://localhost:3000/api';

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(localStorage.getItem(ACCESS_TOKEN_KEY));
  const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_TOKEN_KEY));
  const user = ref<User | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const isAuthenticated = computed(() => !!accessToken.value && !!user.value);

  // Set auth token in axios headers
  function setAuthHeader(token: string | null) {
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
      delete axios.defaults.headers.common['Authorization'];
    }
  }

  // Store tokens
  function storeTokens(access: string, refresh: string) {
    accessToken.value = access;
    refreshToken.value = refresh;
    localStorage.setItem(ACCESS_TOKEN_KEY, access);
    localStorage.setItem(REFRESH_TOKEN_KEY, refresh);
    setAuthHeader(access);
  }

  // Clear tokens
  function clearTokens() {
    accessToken.value = null;
    refreshToken.value = null;
    user.value = null;
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    setAuthHeader(null);
  }

  // Initialize auth state on store creation
  if (accessToken.value) {
    setAuthHeader(accessToken.value);
  }

  // Setup axios interceptor for automatic token refresh
  axios.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config;

      // If error is 401 and we haven't tried refreshing yet
      if (error.response?.status === 401 && !originalRequest._retry && refreshToken.value) {
        originalRequest._retry = true;

        try {
          // Try to refresh the token
          const response = await axios.post<RefreshResponse>(`${API_BASE}/auth/refresh`, {
            refreshToken: refreshToken.value,
          });

          // Store new tokens
          storeTokens(response.data.accessToken, response.data.refreshToken);

          // Retry the original request with new token
          originalRequest.headers['Authorization'] = `Bearer ${response.data.accessToken}`;
          return axios(originalRequest);
        } catch (refreshError) {
          // Refresh failed, log out user
          clearTokens();
          window.location.href = '/login';
          return Promise.reject(refreshError);
        }
      }

      return Promise.reject(error);
    }
  );

  async function register(email: string, username: string, password: string) {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await axios.post<AuthResponse>(`${API_BASE}/auth/register`, {
        email,
        username,
        password,
      });

      storeTokens(response.data.accessToken, response.data.refreshToken);
      user.value = response.data.user;

      return true;
    } catch (err: any) {
      error.value = err.response?.data?.error?.message || 'Registration failed';
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function login(email: string, password: string) {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await axios.post<AuthResponse>(`${API_BASE}/auth/login`, {
        email,
        password,
      });

      storeTokens(response.data.accessToken, response.data.refreshToken);
      user.value = response.data.user;

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
      // Token is invalid, clear auth
      clearTokens();
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function logout() {
    try {
      if (refreshToken.value) {
        await axios.post(`${API_BASE}/auth/logout`, {
          refreshToken: refreshToken.value,
        });
      }
    } catch (err) {
      // Ignore errors on logout
    } finally {
      clearTokens();
    }
  }

  async function logoutAll() {
    try {
      await axios.post(`${API_BASE}/auth/logout-all`);
    } catch (err) {
      // Ignore errors
    } finally {
      clearTokens();
    }
  }

  return {
    accessToken,
    refreshToken,
    user,
    isLoading,
    error,
    isAuthenticated,
    register,
    login,
    fetchCurrentUser,
    logout,
    logoutAll,
  };
});

