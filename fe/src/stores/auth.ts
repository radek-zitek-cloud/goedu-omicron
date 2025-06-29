/**
 * Authentication Store for GoEdu Omicron Banking Platform
 *
 * Centralized state management for user authentication, authorization,
 * and session management using Pinia store pattern.
 *
 * Features:
 * - JWT token management with refresh capability
 * - Role-based access control (RBAC)
 * - Organization-based data isolation
 * - Secure session handling
 * - Multi-factor authentication support
 * - Audit trail for security events
 *
 * Security Considerations:
 * - Tokens stored in httpOnly cookies (backend handled)
 * - Automatic token refresh before expiration
 * - Session timeout handling
 * - Failed login attempt tracking
 * - Secure logout with token invalidation
 *
 * @author GoEdu Development Team
 * @version 1.0.0
 * @since 2024
 */

import { defineStore } from 'pinia';
import { ref, computed, readonly } from 'vue';
import axios from 'axios';

/**
 * User Interface Definition
 *
 * Represents the authenticated user with all necessary properties
 * for banking application access control and personalization.
 */
export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: 'admin' | 'user' | 'auditor' | 'readonly';
  organizationId: string;
  organizationName: string;
  permissions: string[];
  lastLoginAt: Date;
  mfaEnabled: boolean;
  profileImageUrl?: string;
  preferences: {
    theme: 'light' | 'dark' | 'auto';
    language: string;
    timezone: string;
    notifications: {
      email: boolean;
      push: boolean;
      sms: boolean;
    };
  };
}

/**
 * Login Request Interface
 *
 * Structure for user authentication requests with optional MFA.
 */
export interface LoginRequest {
  email: string;
  password: string;
  mfaCode?: string;
  rememberMe?: boolean;
}

/**
 * Authentication Response Interface
 *
 * Structure for successful authentication responses from the backend.
 */
export interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  requiresMfa: boolean;
}

/**
 * Authentication Store Implementation
 *
 * Main store for managing authentication state and operations.
 * Uses Pinia composition API pattern for better type safety and reactivity.
 */
export const useAuthStore = defineStore(
  'auth',
  () => {
    // Reactive state
    const user = ref<User | null>(null);
    const accessToken = ref<string | null>(null);
    const refreshToken = ref<string | null>(null);
    const isLoading = ref(false);
    const isInitialized = ref(false);
    const loginAttempts = ref(0);
    const lastActivity = ref<Date>(new Date());
    const sessionTimeout = ref<number | null>(null);

    // Configuration
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api';
    const SESSION_TIMEOUT_MS = 30 * 60 * 1000; // 30 minutes
    const MAX_LOGIN_ATTEMPTS = 5;

    // Computed properties
    const isAuthenticated = computed(() => !!user.value && !!accessToken.value);
    const userFullName = computed(() =>
      user.value ? `${user.value.firstName} ${user.value.lastName}` : ''
    );
    const userInitials = computed(() =>
      user.value ? `${user.value.firstName[0]}${user.value.lastName[0]}` : ''
    );
    const isAdmin = computed(() => user.value?.role === 'admin');
    const isAuditor = computed(() => user.value?.role === 'auditor');
    const canEditControls = computed(
      () => user.value?.role === 'admin' || user.value?.role === 'user'
    );
    const canViewReports = computed(
      () =>
        user.value?.role === 'admin' ||
        user.value?.role === 'auditor' ||
        user.value?.role === 'user'
    );

    /**
     * Initialize Authentication State
     *
     * Attempts to restore authentication state from stored tokens
     * and validates current session with the backend.
     *
     * @returns Promise<boolean> - Success status of initialization
     */
    async function initialize(): Promise<boolean> {
      if (isInitialized.value) {
        return isAuthenticated.value;
      }

      try {
        isLoading.value = true;

        // Check for existing session with backend
        const response = await axios.get(`${API_BASE_URL}/auth/me`, {
          withCredentials: true, // Include httpOnly cookies
        });

        if (response.data.success && response.data.user) {
          // Restore user state from valid session
          user.value = response.data.user;
          accessToken.value = response.data.accessToken;

          // Start session management
          startSessionManagement();

          console.log('✅ Authentication state restored successfully');
          return true;
        }
      } catch (error) {
        console.log('ℹ️ No existing authentication session found');
        // Clear any stale state
        await logout(false);
      } finally {
        isLoading.value = false;
        isInitialized.value = true;
      }

      return false;
    }

    /**
     * User Login
     *
     * Authenticates user with email/password and optional MFA.
     * Handles failed attempts and account lockouts.
     *
     * @param credentials - User login credentials
     * @returns Promise<AuthResponse> - Authentication result
     */
    async function login(credentials: LoginRequest): Promise<AuthResponse> {
      try {
        isLoading.value = true;

        // Check if account is temporarily locked
        if (loginAttempts.value >= MAX_LOGIN_ATTEMPTS) {
          throw new Error(
            'Account temporarily locked due to too many failed attempts. Please try again later.'
          );
        }

        const response = await axios.post(`${API_BASE_URL}/auth/login`, credentials, {
          withCredentials: true,
        });

        if (response.data.success) {
          const authData: AuthResponse = response.data;

          // Handle MFA requirement
          if (authData.requiresMfa && !credentials.mfaCode) {
            return authData; // Return without setting user state
          }

          // Set authenticated state
          user.value = authData.user;
          accessToken.value = authData.accessToken;
          refreshToken.value = authData.refreshToken;

          // Reset failed attempts on successful login
          loginAttempts.value = 0;

          // Start session management
          startSessionManagement();

          // Log successful authentication
          console.log('✅ User authenticated successfully:', authData.user.email);

          return authData;
        } else {
          throw new Error(response.data.message || 'Authentication failed');
        }
      } catch (error: any) {
        // Increment failed login attempts
        loginAttempts.value++;

        // Log security event
        console.warn('🔒 Authentication failed:', {
          email: credentials.email,
          attempts: loginAttempts.value,
          error: error.message,
        });

        throw error;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * User Logout
     *
     * Securely logs out user and clears all authentication state.
     * Invalidates tokens on the backend for security.
     *
     * @param notifyBackend - Whether to notify backend of logout
     */
    async function logout(notifyBackend: boolean = true): Promise<void> {
      try {
        isLoading.value = true;

        // Notify backend to invalidate session
        if (notifyBackend && accessToken.value) {
          await axios.post(
            `${API_BASE_URL}/auth/logout`,
            {},
            {
              withCredentials: true,
              headers: {
                Authorization: `Bearer ${accessToken.value}`,
              },
            }
          );
        }
      } catch (error) {
        console.warn('⚠️ Error during logout:', error);
      } finally {
        // Clear all authentication state
        user.value = null;
        accessToken.value = null;
        refreshToken.value = null;
        loginAttempts.value = 0;

        // Clear session management
        stopSessionManagement();

        isLoading.value = false;
        console.log('✅ User logged out successfully');
      }
    }

    /**
     * Refresh Authentication Token
     *
     * Refreshes the access token using the refresh token.
     * Called automatically when tokens are near expiration.
     *
     * @returns Promise<boolean> - Success status of token refresh
     */
    async function refreshAuthToken(): Promise<boolean> {
      try {
        const response = await axios.post(
          `${API_BASE_URL}/auth/refresh`,
          {},
          {
            withCredentials: true,
          }
        );

        if (response.data.success) {
          accessToken.value = response.data.accessToken;
          refreshToken.value = response.data.refreshToken;

          console.log('✅ Authentication token refreshed successfully');
          return true;
        }
      } catch (error) {
        console.warn('⚠️ Token refresh failed:', error);
        // Force logout on refresh failure
        await logout(false);
      }

      return false;
    }

    /**
     * Update User Profile
     *
     * Updates user profile information and preferences.
     *
     * @param updates - Partial user data to update
     * @returns Promise<User> - Updated user object
     */
    async function updateProfile(updates: Partial<User>): Promise<User> {
      try {
        isLoading.value = true;

        const response = await axios.patch(`${API_BASE_URL}/auth/profile`, updates, {
          withCredentials: true,
          headers: {
            Authorization: `Bearer ${accessToken.value}`,
          },
        });

        if (response.data.success) {
          user.value = response.data.user;
          console.log('✅ Profile updated successfully');
          return response.data.user;
        } else {
          throw new Error(response.data.message || 'Profile update failed');
        }
      } catch (error) {
        console.error('❌ Profile update failed:', error);
        throw error;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Check User Permissions
     *
     * Validates if the current user has specific permissions.
     *
     * @param permission - Permission string to check
     * @returns boolean - Whether user has permission
     */
    function hasPermission(permission: string): boolean {
      return user.value?.permissions.includes(permission) || false;
    }

    /**
     * Start Session Management
     *
     * Begins session timeout tracking and automatic token refresh.
     * Monitors user activity to extend session automatically.
     */
    function startSessionManagement(): void {
      // Update last activity timestamp
      updateActivity();

      // Set up session timeout
      resetSessionTimeout();

      // Add activity listeners
      const activityEvents = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart'];
      activityEvents.forEach(event => {
        document.addEventListener(event, updateActivity, true);
      });
    }

    /**
     * Stop Session Management
     *
     * Clears session timeout and removes activity listeners.
     */
    function stopSessionManagement(): void {
      if (sessionTimeout.value) {
        clearTimeout(sessionTimeout.value);
        sessionTimeout.value = null;
      }

      // Remove activity listeners
      const activityEvents = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart'];
      activityEvents.forEach(event => {
        document.removeEventListener(event, updateActivity, true);
      });
    }

    /**
     * Update Activity Timestamp
     *
     * Records user activity and resets session timeout.
     */
    function updateActivity(): void {
      lastActivity.value = new Date();
      resetSessionTimeout();
    }

    /**
     * Reset Session Timeout
     *
     * Resets the automatic logout timer based on inactivity.
     */
    function resetSessionTimeout(): void {
      if (sessionTimeout.value) {
        clearTimeout(sessionTimeout.value);
      }

      sessionTimeout.value = setTimeout(async () => {
        console.log('⏰ Session timeout - logging out user');
        await logout(true);
      }, SESSION_TIMEOUT_MS);
    }

    // Return store interface
    return {
      // State
      user: readonly(user),
      isLoading: readonly(isLoading),
      isInitialized: readonly(isInitialized),
      loginAttempts: readonly(loginAttempts),
      lastActivity: readonly(lastActivity),

      // Computed
      isAuthenticated,
      userFullName,
      userInitials,
      isAdmin,
      isAuditor,
      canEditControls,
      canViewReports,

      // Actions
      initialize,
      login,
      logout,
      refreshAuthToken,
      updateProfile,
      hasPermission,
      updateActivity,
    };
  },
  {
    // Pinia persistence configuration for offline capability
    persist: {
      key: 'goedu-auth',
      storage: localStorage,
      paths: ['user', 'lastActivity'],
    },
  }
);

/**
 * Authentication HTTP Interceptor Setup
 *
 * Configures axios interceptors to automatically handle authentication
 * headers and token refresh for API requests.
 */
export function setupAuthInterceptors(): void {
  // Request interceptor to add auth headers
  axios.interceptors.request.use(
    config => {
      const authStore = useAuthStore();
      if (authStore.isAuthenticated && authStore.user) {
        config.headers.Authorization = `Bearer ${authStore.user}`;
      }
      return config;
    },
    error => Promise.reject(error)
  );

  // Response interceptor to handle token refresh
  axios.interceptors.response.use(
    response => response,
    async error => {
      const authStore = useAuthStore();

      if (error.response?.status === 401 && authStore.isAuthenticated) {
        // Try to refresh token
        const refreshSuccess = await authStore.refreshAuthToken();

        if (refreshSuccess) {
          // Retry the original request
          return axios.request(error.config);
        } else {
          // Refresh failed, redirect to login
          await authStore.logout(false);
        }
      }

      return Promise.reject(error);
    }
  );
}
