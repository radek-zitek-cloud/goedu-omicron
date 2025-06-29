# Pinia Store Documentation

This document provides comprehensive information about the Pinia state management setup for the GoEdu Omicron Banking Platform.

## Overview

The application uses Pinia for centralized state management with the following key features:

- **TypeScript Integration**: Full type safety with TypeScript interfaces
- **Devtools Support**: Pinia devtools integration for development debugging
- **Persistence Layer**: Offline capability with localStorage persistence
- **Modular Design**: Separate stores for different business domains
- **HTTP Interceptor Integration**: Automatic authentication handling

## Store Architecture

### 1. Authentication Store (`useAuthStore`)

Manages user authentication, sessions, and security-related state.

#### Key Features:
- JWT token management with refresh capability
- Role-based access control (RBAC)
- Session timeout handling
- Multi-factor authentication support
- Secure logout with token invalidation

#### Usage Example:
```typescript
<script setup lang="ts">
import { useAuthStore } from '@/stores';

const authStore = useAuthStore();

// Login user
const login = async () => {
  const success = await authStore.login({
    email: 'user@example.com',
    password: 'password123',
    rememberMe: true
  });
  
  if (success) {
    console.log('Login successful');
    console.log('User:', authStore.user);
  }
};

// Check permissions
const canEditControls = authStore.canEditControls;
const hasPermission = authStore.hasPermission('EDIT_CONTROLS');

// Logout
const logout = () => authStore.logout();
</script>
```

#### Persisted State:
- User information
- Last activity timestamp

### 2. Organizations Store (`useOrganizationsStore`)

Manages organization data, member management, and organization-specific settings.

#### Key Features:
- Organization CRUD operations
- Current organization selection
- Member management with role-based access
- Organization preferences and settings
- Multi-tenant data isolation

#### Usage Example:
```typescript
<script setup lang="ts">
import { useOrganizationsStore } from '@/stores';

const organizationsStore = useOrganizationsStore();

// Load organizations
const loadOrgs = async () => {
  await organizationsStore.loadOrganizations();
  console.log('Organizations:', organizationsStore.organizations);
};

// Set current organization
const selectOrg = async (orgId: string) => {
  await organizationsStore.setCurrentOrganization(orgId);
  console.log('Current org:', organizationsStore.currentOrganization);
};

// Create new organization
const createOrg = async () => {
  const newOrg = await organizationsStore.createOrganization({
    name: 'New Bank',
    code: 'NEWBANK',
    description: 'A new banking organization',
    industry: 'banking',
    size: 'medium',
    region: 'North America',
    country: 'USA',
    timezone: 'America/New_York',
    currency: 'USD',
    contactEmail: 'admin@newbank.com',
    contactPhone: '+1-555-0123',
    address: {
      street: '123 Banking St',
      city: 'New York',
      state: 'NY',
      postalCode: '10001',
      country: 'USA'
    },
    subscriptionPlan: 'professional'
  });
  
  if (newOrg) {
    console.log('Organization created:', newOrg);
  }
};

// Invite member
const inviteMember = async () => {
  const success = await organizationsStore.inviteMember({
    email: 'newuser@example.com',
    role: 'user',
    permissions: ['VIEW_CONTROLS', 'EDIT_CONTROLS'],
    message: 'Welcome to our organization!'
  });
  
  if (success) {
    console.log('Member invited successfully');
  }
};
</script>
```

#### Persisted State:
- Current organization
- Organizations list
- Last sync timestamp

### 3. Controls Store (`useControlsStore`)

Manages banking controls, test executions, and compliance tracking.

#### Key Features:
- Control CRUD operations with full lifecycle management
- Test execution and result tracking
- Evidence management and file attachments
- Compliance reporting and analytics
- Real-time collaboration and updates
- Offline queue for sync when online

#### Usage Example:
```typescript
<script setup lang="ts">
import { useControlsStore } from '@/stores';

const controlsStore = useControlsStore();

// Load controls for organization
const loadControls = async (orgId: string) => {
  await controlsStore.loadControls(orgId);
  console.log('Controls loaded:', controlsStore.controlsCount);
};

// Create new control
const createControl = async () => {
  const newControl = await controlsStore.createControl({
    code: 'FIN-001',
    title: 'Daily Cash Reconciliation',
    description: 'Reconcile cash positions daily',
    objective: 'Ensure accurate cash reporting',
    category: 'financial',
    subcategory: 'cash_management',
    riskLevel: 'high',
    frequency: 'daily',
    ownerId: 'user123',
    reviewerId: 'manager456',
    regulatoryReferences: ['SOX', 'FDIC'],
    frameworkMappings: [{
      framework: 'COSO',
      controlId: 'CC1.1',
      requirement: 'Control Environment'
    }],
    testProcedures: [{
      step: 1,
      description: 'Retrieve cash balances from system',
      expectedResult: 'Balances retrieved successfully',
      automatable: true
    }],
    keyAttributes: [{
      name: 'threshold',
      value: '10000',
      dataType: 'number',
      required: true
    }],
    tags: ['cash', 'reconciliation', 'daily'],
    relatedControls: [],
    effectiveDate: new Date()
  });
  
  if (newControl) {
    console.log('Control created:', newControl);
  }
};

// Execute test
const executeTest = async (controlId: string) => {
  const test = await controlsStore.executeTest({
    controlId,
    testPeriod: {
      startDate: new Date('2024-01-01'),
      endDate: new Date('2024-01-31'),
      label: 'January 2024'
    },
    plannedStartDate: new Date(),
    plannedEndDate: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
    testerId: 'tester123',
    reviewerId: 'reviewer456',
    notes: 'Regular monthly test execution'
  });
  
  if (test) {
    console.log('Test execution started:', test);
  }
};

// Filter controls
const filterControls = () => {
  controlsStore.setFilters({
    status: ['active'],
    category: ['financial', 'operational'],
    riskLevel: ['high', 'critical'],
    search: 'cash'
  });
  
  console.log('Filtered controls:', controlsStore.filteredControls);
};

// Get upcoming tests
const upcomingTests = controlsStore.upcomingTests;
const overdueTests = controlsStore.overdueTests;
</script>
```

#### Persisted State:
- Controls list
- Current filters
- Last sync timestamp
- Offline queue for operations

## Store Initialization

The stores are automatically initialized when the app starts:

```typescript
// main.ts
import { initializeStores } from '@/stores';

async function initializeApp() {
  const app = createApp(App);
  
  // Register Pinia
  app.use(pinia);
  
  // Initialize stores and HTTP interceptors
  initializeStores();
  
  app.mount('#app');
}
```

## HTTP Interceptor Integration

The authentication store automatically sets up axios interceptors for:

- **Request Interceptor**: Adds authentication headers to all requests
- **Response Interceptor**: Handles 401 errors and automatic token refresh

## Persistence Configuration

All stores use localStorage persistence with the following strategy:

- **Selective Persistence**: Only critical data is persisted
- **Offline Capability**: Stores maintain functionality when offline
- **Automatic Sync**: Data syncs when connection is restored

## Devtools Integration

In development mode, Pinia devtools are automatically enabled:

- **State Inspection**: View and modify store state
- **Action Tracking**: Monitor all store actions
- **Time Travel**: Debug state changes over time
- **Performance Monitoring**: Track store performance

## TypeScript Integration

All stores are fully typed with TypeScript:

- **Interface Definitions**: Comprehensive type definitions for all data structures
- **Type Safety**: Full type checking for all store operations
- **IntelliSense Support**: Rich IDE support with autocomplete

## Best Practices

### 1. Store Usage in Components
```typescript
// ✅ Good: Use stores in setup function
<script setup lang="ts">
import { useAuthStore } from '@/stores';

const authStore = useAuthStore();
const isLoggedIn = computed(() => authStore.isAuthenticated);
</script>

// ❌ Bad: Don't use stores outside component context
const authStore = useAuthStore(); // This will cause issues
```

### 2. Error Handling
```typescript
// ✅ Good: Handle errors properly
const login = async () => {
  try {
    const success = await authStore.login(credentials);
    if (!success) {
      // Handle login failure
      console.error('Login failed:', authStore.error);
    }
  } catch (error) {
    console.error('Login error:', error);
  }
};
```

### 3. Reactive State
```typescript
// ✅ Good: Use computed for reactive derived state
const userFullName = computed(() => authStore.userFullName);

// ❌ Bad: Direct property access loses reactivity
const userFullName = authStore.userFullName;
```

### 4. Store Cleanup
```typescript
// ✅ Good: Cleanup on logout or organization change
const logout = async () => {
  await authStore.logout();
  organizationsStore.resetStore();
  controlsStore.resetStore();
};
```

## Troubleshooting

### Common Issues:

1. **Store not found**: Ensure Pinia is registered before using stores
2. **Persistence not working**: Check localStorage is available and not full
3. **HTTP interceptors not working**: Ensure `initializeStores()` is called
4. **Type errors**: Ensure TypeScript is properly configured

### Debug Information:

```typescript
// Get store health status
import { storeUtils } from '@/stores';

const healthStatus = storeUtils.getStoreHealth();
console.log('Store health:', healthStatus);
```

## Migration Guide

If migrating from Vuex or other state management:

1. **Replace Vuex modules** with Pinia stores
2. **Update component imports** to use new store functions
3. **Migrate state structure** to new interface definitions
4. **Test persistence** to ensure offline capability works

## Performance Considerations

- **Lazy Loading**: Stores are only initialized when first used
- **Selective Persistence**: Only essential data is persisted
- **Efficient Updates**: Use reactive state for optimal performance
- **Memory Management**: Reset stores when not needed

## Security Considerations

- **Token Storage**: Sensitive tokens use secure storage mechanisms
- **Data Isolation**: Organization data is properly isolated
- **Permission Checks**: Role-based access control is enforced
- **Audit Trails**: All sensitive operations are logged

This documentation provides a comprehensive guide to using the Pinia stores in the GoEdu Omicron Banking Platform. For specific implementation details, refer to the individual store files and their TypeScript interfaces.