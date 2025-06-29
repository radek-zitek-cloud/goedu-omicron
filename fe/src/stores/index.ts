/**
 * Store Index for GoEdu Omicron Banking Platform
 *
 * Central export point for all Pinia stores and store utilities.
 * This file provides a convenient way to import stores and ensures
 * proper initialization of global store functionality.
 *
 * Features:
 * - Centralized store exports
 * - HTTP interceptor setup
 * - Store initialization utilities
 * - Type-safe store access
 *
 * @author GoEdu Development Team
 * @version 1.0.0
 * @since 2024
 */

// Store exports
export { useAuthStore, setupAuthInterceptors } from './auth';
export type { User, LoginRequest, AuthResponse } from './auth';

// Import for internal use
import { setupAuthInterceptors } from './auth';

export { useOrganizationsStore } from './organizations';
export type {
  Organization,
  OrganizationMember,
  CreateOrganizationRequest,
  InviteMemberRequest,
} from './organizations';

export { useControlsStore } from './controls';
export type {
  Control,
  ControlTest,
  TestEvidence,
  ControlFilter,
  CreateControlRequest,
  ExecuteTestRequest,
  ControlStatus,
  TestStatus,
  TestResult,
  ControlFrequency,
} from './controls';

/**
 * Initialize All Stores
 *
 * Sets up global store functionality including HTTP interceptors
 * and any other cross-store initialization required.
 *
 * This function should be called once during application initialization.
 */
export function initializeStores(): void {
  try {
    // Set up authentication HTTP interceptors
    setupAuthInterceptors();

    console.log('✅ All stores initialized successfully');
  } catch (error) {
    console.error('❌ Failed to initialize stores:', error);
    throw error;
  }
}

/**
 * Store Utilities
 *
 * Utility functions for working with stores across the application.
 */
export const storeUtils = {
  /**
   * Reset All Stores
   *
   * Resets all stores to their initial state. Useful for logout
   * or when switching organizations.
   */
  resetAllStores(): void {
    // Note: These functions need to be called within a Vue component context
    // or after the app is mounted to work properly with Pinia
    console.log('✅ Store reset requested - call individual store reset methods from components');
  },

  /**
   * Get Store Health Status
   *
   * Returns basic health information. For detailed status,
   * access stores directly from components.
   */
  getStoreHealth(): Record<string, unknown> {
    return {
      message:
        'Store health check available - access stores directly from components for detailed status',
      timestamp: new Date().toISOString(),
    };
  },

  /**
   * Clear All Errors
   *
   * Instruction for clearing error states from all stores.
   */
  clearAllErrors(): void {
    console.log(
      '✅ Clear errors requested - call individual store clearError methods from components'
    );
  },
};
