/**
 * Organization Store for GoEdu Omicron Banking Platform
 *
 * Centralized state management for organization data, including organization
 * selection, member management, and organization-specific configurations.
 *
 * Features:
 * - Organization CRUD operations
 * - Current organization selection
 * - Member management with role-based access
 * - Organization preferences and settings
 * - Multi-tenant data isolation
 * - Offline capability with persistence
 *
 * Security Considerations:
 * - Organization-based data isolation
 * - Role-based access control for organization operations
 * - Audit trail for organization changes
 * - Secure member invitation process
 *
 * @author GoEdu Development Team
 * @version 1.0.0
 * @since 2024
 */

import { defineStore } from 'pinia';
import { ref, computed, readonly } from 'vue';
import axios from 'axios';

/**
 * Organization Interface Definition
 *
 * Represents an organization within the banking platform with all necessary
 * properties for multi-tenant operations and compliance management.
 */
export interface Organization {
  id: string;
  name: string;
  code: string;
  description: string;
  industry: 'banking' | 'credit_union' | 'fintech' | 'insurance' | 'other';
  size: 'small' | 'medium' | 'large' | 'enterprise';
  region: string;
  country: string;
  timezone: string;
  currency: string;
  logoUrl?: string;
  website?: string;
  contactEmail: string;
  contactPhone: string;
  address: {
    street: string;
    city: string;
    state: string;
    postalCode: string;
    country: string;
  };
  regulatoryFrameworks: string[];
  subscriptionPlan: 'starter' | 'professional' | 'enterprise';
  subscriptionStatus: 'active' | 'suspended' | 'trial' | 'expired';
  memberCount: number;
  maxMembers: number;
  createdAt: Date;
  updatedAt: Date;
  settings: {
    allowInvitations: boolean;
    requireMfa: boolean;
    sessionTimeoutMinutes: number;
    enableAuditLog: boolean;
    dataRetentionDays: number;
    allowDataExport: boolean;
    customBranding: boolean;
    integrations: {
      sso: boolean;
      ldap: boolean;
      api: boolean;
    };
  };
}

/**
 * Organization Member Interface
 *
 * Represents a member of an organization with role and permission information.
 */
export interface OrganizationMember {
  id: string;
  userId: string;
  organizationId: string;
  email: string;
  firstName: string;
  lastName: string;
  role: 'owner' | 'admin' | 'user' | 'auditor' | 'readonly';
  status: 'active' | 'invited' | 'suspended' | 'inactive';
  permissions: string[];
  joinedAt: Date;
  lastActiveAt: Date;
  invitedBy?: string;
  invitedAt?: Date;
}

/**
 * Organization Creation Request Interface
 */
export interface CreateOrganizationRequest {
  name: string;
  code: string;
  description: string;
  industry: Organization['industry'];
  size: Organization['size'];
  region: string;
  country: string;
  timezone: string;
  currency: string;
  contactEmail: string;
  contactPhone: string;
  address: Organization['address'];
  subscriptionPlan: Organization['subscriptionPlan'];
}

/**
 * Member Invitation Request Interface
 */
export interface InviteMemberRequest {
  email: string;
  role: OrganizationMember['role'];
  permissions?: string[];
  message?: string;
}

/**
 * Organization Store Implementation
 *
 * Main store for managing organization state and operations.
 * Uses Pinia composition API pattern for better type safety and reactivity.
 */
export const useOrganizationsStore = defineStore(
  'organizations',
  () => {
    // Reactive state
    const organizations = ref<Organization[]>([]);
    const currentOrganization = ref<Organization | null>(null);
    const members = ref<OrganizationMember[]>([]);
    const isLoading = ref(false);
    const isLoadingMembers = ref(false);
    const error = ref<string | null>(null);
    const lastSyncTime = ref<Date | null>(null);

    // Configuration
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api';

    // Computed properties
    const hasCurrentOrganization = computed(() => !!currentOrganization.value);
    const currentOrganizationId = computed(() => currentOrganization.value?.id || null);
    const currentOrganizationName = computed(() => currentOrganization.value?.name || '');
    const isCurrentOrganizationOwner = computed(() => {
      if (!currentOrganization.value || !members.value.length) return false;
      // This would need to be cross-referenced with current user from auth store
      return false; // Placeholder implementation
    });
    const currentOrganizationMembers = computed(() =>
      members.value.filter(member => member.organizationId === currentOrganizationId.value)
    );
    const organizationCount = computed(() => organizations.value.length);
    const canInviteMembers = computed(
      () =>
        currentOrganization.value?.settings.allowInvitations &&
        currentOrganization.value?.memberCount < currentOrganization.value?.maxMembers
    );

    /**
     * Load Organizations
     *
     * Fetches all organizations accessible to the current user.
     *
     * @returns Promise<boolean> - Success status
     */
    async function loadOrganizations(): Promise<boolean> {
      if (isLoading.value) return false;

      try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.get(`${API_BASE_URL}/organizations`, {
          withCredentials: true,
        });

        if (response.data.success) {
          organizations.value = response.data.organizations.map((org: Record<string, any>) => ({
            ...org,
            createdAt: new Date(org.createdAt),
            updatedAt: new Date(org.updatedAt),
          }));

          lastSyncTime.value = new Date();
          console.log('✅ Organizations loaded successfully:', organizations.value.length);
          return true;
        }

        throw new Error(response.data.message || 'Failed to load organizations');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to load organizations:', errorMessage);
        return false;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Set Current Organization
     *
     * Sets the active organization for the current session.
     *
     * @param organizationId - ID of the organization to set as current
     * @returns Promise<boolean> - Success status
     */
    async function setCurrentOrganization(organizationId: string): Promise<boolean> {
      try {
        const organization = organizations.value.find(org => org.id === organizationId);

        if (!organization) {
          error.value = 'Organization not found';
          return false;
        }

        currentOrganization.value = organization;

        // Load members for the current organization
        await loadOrganizationMembers(organizationId);

        console.log('✅ Current organization set:', organization.name);
        return true;
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'Failed to set current organization';
        error.value = errorMessage;
        console.error('❌ Failed to set current organization:', errorMessage);
        return false;
      }
    }

    /**
     * Create Organization
     *
     * Creates a new organization with the provided details.
     *
     * @param request - Organization creation request data
     * @returns Promise<Organization | null> - Created organization or null on failure
     */
    async function createOrganization(
      request: CreateOrganizationRequest
    ): Promise<Organization | null> {
      try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.post(`${API_BASE_URL}/organizations`, request, {
          withCredentials: true,
        });

        if (response.data.success) {
          const newOrganization: Organization = {
            ...response.data.organization,
            createdAt: new Date(response.data.organization.createdAt),
            updatedAt: new Date(response.data.organization.updatedAt),
          };

          organizations.value.push(newOrganization);
          console.log('✅ Organization created successfully:', newOrganization.name);
          return newOrganization;
        }

        throw new Error(response.data.message || 'Failed to create organization');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to create organization:', errorMessage);
        return null;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Update Organization
     *
     * Updates an existing organization with new data.
     *
     * @param organizationId - ID of the organization to update
     * @param updates - Partial organization data to update
     * @returns Promise<boolean> - Success status
     */
    async function updateOrganization(
      organizationId: string,
      updates: Partial<Organization>
    ): Promise<boolean> {
      try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.patch(
          `${API_BASE_URL}/organizations/${organizationId}`,
          updates,
          { withCredentials: true }
        );

        if (response.data.success) {
          const updatedOrg = response.data.organization;
          const index = organizations.value.findIndex(org => org.id === organizationId);

          if (index !== -1) {
            organizations.value[index] = {
              ...updatedOrg,
              createdAt: new Date(updatedOrg.createdAt),
              updatedAt: new Date(updatedOrg.updatedAt),
            };

            // Update current organization if it's the one being updated
            if (currentOrganization.value?.id === organizationId) {
              currentOrganization.value = organizations.value[index];
            }
          }

          console.log('✅ Organization updated successfully:', organizationId);
          return true;
        }

        throw new Error(response.data.message || 'Failed to update organization');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to update organization:', errorMessage);
        return false;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Load Organization Members
     *
     * Fetches members for a specific organization.
     *
     * @param organizationId - ID of the organization
     * @returns Promise<boolean> - Success status
     */
    async function loadOrganizationMembers(organizationId: string): Promise<boolean> {
      if (isLoadingMembers.value) return false;

      try {
        isLoadingMembers.value = true;

        const response = await axios.get(
          `${API_BASE_URL}/organizations/${organizationId}/members`,
          { withCredentials: true }
        );

        if (response.data.success) {
          const orgMembers = response.data.members.map((member: Record<string, any>) => ({
            ...member,
            joinedAt: new Date(member.joinedAt),
            lastActiveAt: new Date(member.lastActiveAt),
            invitedAt: member.invitedAt ? new Date(member.invitedAt) : undefined,
          }));

          // Replace members for this organization
          members.value = members.value.filter(m => m.organizationId !== organizationId);
          members.value.push(...orgMembers);

          console.log('✅ Organization members loaded:', orgMembers.length);
          return true;
        }

        throw new Error(response.data.message || 'Failed to load members');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to load members';
        error.value = errorMessage;
        console.error('❌ Failed to load organization members:', errorMessage);
        return false;
      } finally {
        isLoadingMembers.value = false;
      }
    }

    /**
     * Invite Member
     *
     * Invites a new member to the current organization.
     *
     * @param request - Member invitation request data
     * @returns Promise<boolean> - Success status
     */
    async function inviteMember(request: InviteMemberRequest): Promise<boolean> {
      if (!currentOrganization.value) {
        error.value = 'No current organization selected';
        return false;
      }

      try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.post(
          `${API_BASE_URL}/organizations/${currentOrganization.value.id}/invite`,
          request,
          { withCredentials: true }
        );

        if (response.data.success) {
          // Reload members to get the updated list
          await loadOrganizationMembers(currentOrganization.value.id);
          console.log('✅ Member invited successfully:', request.email);
          return true;
        }

        throw new Error(response.data.message || 'Failed to invite member');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to invite member:', errorMessage);
        return false;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Remove Member
     *
     * Removes a member from the organization.
     *
     * @param memberId - ID of the member to remove
     * @returns Promise<boolean> - Success status
     */
    async function removeMember(memberId: string): Promise<boolean> {
      if (!currentOrganization.value) {
        error.value = 'No current organization selected';
        return false;
      }

      try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.delete(
          `${API_BASE_URL}/organizations/${currentOrganization.value.id}/members/${memberId}`,
          { withCredentials: true }
        );

        if (response.data.success) {
          // Remove member from local state
          members.value = members.value.filter(member => member.id !== memberId);
          console.log('✅ Member removed successfully:', memberId);
          return true;
        }

        throw new Error(response.data.message || 'Failed to remove member');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to remove member:', errorMessage);
        return false;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Clear Error
     *
     * Clears the current error state.
     */
    function clearError(): void {
      error.value = null;
    }

    /**
     * Reset Store
     *
     * Resets all store state to initial values.
     */
    function resetStore(): void {
      organizations.value = [];
      currentOrganization.value = null;
      members.value = [];
      isLoading.value = false;
      isLoadingMembers.value = false;
      error.value = null;
      lastSyncTime.value = null;
    }

    // Return store interface
    return {
      // State
      organizations: readonly(organizations),
      currentOrganization: readonly(currentOrganization),
      members: readonly(members),
      isLoading: readonly(isLoading),
      isLoadingMembers: readonly(isLoadingMembers),
      error: readonly(error),
      lastSyncTime: readonly(lastSyncTime),

      // Computed
      hasCurrentOrganization,
      currentOrganizationId,
      currentOrganizationName,
      isCurrentOrganizationOwner,
      currentOrganizationMembers,
      organizationCount,
      canInviteMembers,

      // Actions
      loadOrganizations,
      setCurrentOrganization,
      createOrganization,
      updateOrganization,
      loadOrganizationMembers,
      inviteMember,
      removeMember,
      clearError,
      resetStore,
    };
  },
  {
    // Pinia persistence configuration for offline capability
    persist: {
      key: 'goedu-organizations',
      storage: localStorage,
      paths: ['currentOrganization', 'organizations', 'lastSyncTime'],
    },
  }
);
