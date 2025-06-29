<!--
  Store Demo Component
  
  This component demonstrates how to use the Pinia stores in Vue components.
  It shows basic usage patterns for all three stores with proper TypeScript integration.
-->
<template>
  <div class="store-demo">
    <v-container>
      <v-row>
        <v-col cols="12">
          <h1>Pinia Store Demo</h1>
          <p>This demo shows how to interact with all the Pinia stores.</p>
        </v-col>
      </v-row>

      <!-- Authentication Store Demo -->
      <v-row>
        <v-col cols="12" md="4">
          <v-card>
            <v-card-title>Authentication Store</v-card-title>
            <v-card-text>
              <div class="mb-4">
                <strong>Status:</strong>
                <v-chip
                  :color="authStore.isAuthenticated ? 'success' : 'error'"
                  size="small"
                  class="ml-2"
                >
                  {{ authStore.isAuthenticated ? 'Authenticated' : 'Not Authenticated' }}
                </v-chip>
              </div>
              
              <div v-if="authStore.user" class="mb-4">
                <strong>User:</strong> {{ authStore.userFullName }}<br>
                <strong>Role:</strong> {{ authStore.user.role }}<br>
                <strong>Organization:</strong> {{ authStore.user.organizationName }}
              </div>

              <div class="mb-4">
                <strong>Permissions:</strong>
                <v-chip-group>
                  <v-chip size="small" color="primary">
                    Edit Controls: {{ authStore.canEditControls ? '✓' : '✗' }}
                  </v-chip>
                  <v-chip size="small" color="primary">
                    View Reports: {{ authStore.canViewReports ? '✓' : '✗' }}
                  </v-chip>
                </v-chip-group>
              </div>

              <v-btn
                v-if="!authStore.isAuthenticated"
                @click="simulateLogin"
                :loading="authStore.isLoading"
                color="primary"
              >
                Simulate Login
              </v-btn>
              
              <v-btn
                v-else
                @click="logout"
                :loading="authStore.isLoading"
                color="error"
              >
                Logout
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Organizations Store Demo -->
        <v-col cols="12" md="4">
          <v-card>
            <v-card-title>Organizations Store</v-card-title>
            <v-card-text>
              <div class="mb-4">
                <strong>Organizations Count:</strong> {{ organizationsStore.organizationCount }}
              </div>
              
              <div v-if="organizationsStore.currentOrganization" class="mb-4">
                <strong>Current Org:</strong> {{ organizationsStore.currentOrganizationName }}<br>
                <strong>Industry:</strong> {{ organizationsStore.currentOrganization.industry }}<br>
                <strong>Members:</strong> {{ organizationsStore.currentOrganizationMembers.length }}
              </div>

              <div class="mb-4">
                <strong>Status:</strong>
                <v-chip
                  :color="organizationsStore.hasCurrentOrganization ? 'success' : 'warning'"
                  size="small"
                  class="ml-2"
                >
                  {{ organizationsStore.hasCurrentOrganization ? 'Organization Selected' : 'No Organization' }}
                </v-chip>
              </div>

              <v-btn
                @click="loadOrganizations"
                :loading="organizationsStore.isLoading"
                color="primary"
                class="mb-2"
                block
              >
                Load Organizations
              </v-btn>
              
              <v-btn
                @click="simulateCreateOrg"
                :loading="organizationsStore.isLoading"
                color="secondary"
                block
              >
                Simulate Create Org
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Controls Store Demo -->
        <v-col cols="12" md="4">
          <v-card>
            <v-card-title>Controls Store</v-card-title>
            <v-card-text>
              <div class="mb-4">
                <strong>Controls:</strong>
                <div class="d-flex flex-wrap ga-1">
                  <v-chip size="small" color="primary">
                    Total: {{ controlsStore.controlsCount.total }}
                  </v-chip>
                  <v-chip size="small" color="success">
                    Active: {{ controlsStore.controlsCount.active }}
                  </v-chip>
                  <v-chip size="small" color="warning">
                    Draft: {{ controlsStore.controlsCount.draft }}
                  </v-chip>
                  <v-chip size="small" color="error">
                    Overdue: {{ controlsStore.controlsCount.overdue }}
                  </v-chip>
                </div>
              </div>

              <div class="mb-4">
                <strong>Upcoming Tests:</strong> {{ controlsStore.upcomingTests.length }}<br>
                <strong>Overdue Tests:</strong> {{ controlsStore.overdueTests.length }}
              </div>

              <div class="mb-4">
                <strong>Filters Applied:</strong>
                <v-chip
                  v-if="Object.keys(controlsStore.filters).length === 0"
                  size="small"
                  color="grey"
                >
                  None
                </v-chip>
                <v-chip
                  v-else
                  size="small"
                  color="info"
                >
                  {{ Object.keys(controlsStore.filters).length }} filters
                </v-chip>
              </div>

              <v-btn
                @click="loadControls"
                :loading="controlsStore.isLoading"
                color="primary"
                class="mb-2"
                block
              >
                Load Controls
              </v-btn>
              
              <v-btn
                @click="simulateCreateControl"
                :loading="controlsStore.isLoading"
                color="secondary"
                class="mb-2"
                block
              >
                Simulate Create Control
              </v-btn>

              <v-btn
                @click="toggleFilters"
                color="info"
                block
              >
                Toggle Filters
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Store Health Status -->
      <v-row>
        <v-col cols="12">
          <v-card>
            <v-card-title>Store Health Status</v-card-title>
            <v-card-text>
              <pre>{{ JSON.stringify(storeHealth, null, 2) }}</pre>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { 
  useAuthStore, 
  useOrganizationsStore, 
  useControlsStore, 
  storeUtils 
} from '@/stores';

// Initialize stores
const authStore = useAuthStore();
const organizationsStore = useOrganizationsStore();
const controlsStore = useControlsStore();

// Computed store health status
const storeHealth = computed(() => storeUtils.getStoreHealth());

// Authentication actions
const simulateLogin = async () => {
  // In a real app, this would call authStore.login() with actual credentials
  console.log('This would normally call authStore.login() with credentials');
  
  // For demo purposes, we'll show what the login process would look like
  alert('In a real app, this would authenticate with the backend API');
};

const logout = async () => {
  await authStore.logout();
  organizationsStore.resetStore();
  controlsStore.resetStore();
};

// Organization actions
const loadOrganizations = async () => {
  const success = await organizationsStore.loadOrganizations();
  if (success) {
    console.log('Organizations loaded successfully');
  } else {
    console.error('Failed to load organizations:', organizationsStore.error);
    alert('This would normally load organizations from the API');
  }
};

const simulateCreateOrg = async () => {
  // In a real app, this would call organizationsStore.createOrganization()
  console.log('This would create a new organization');
  alert('In a real app, this would create an organization via API');
};

// Controls actions
const loadControls = async () => {
  if (!organizationsStore.currentOrganizationId) {
    alert('Please select an organization first');
    return;
  }
  
  const success = await controlsStore.loadControls(organizationsStore.currentOrganizationId);
  if (success) {
    console.log('Controls loaded successfully');
  } else {
    console.error('Failed to load controls:', controlsStore.error);
    alert('This would normally load controls from the API');
  }
};

const simulateCreateControl = async () => {
  console.log('This would create a new control');
  alert('In a real app, this would create a control via API');
};

const toggleFilters = () => {
  if (Object.keys(controlsStore.filters).length === 0) {
    // Apply demo filters
    controlsStore.setFilters({
      status: ['active'],
      category: ['financial', 'operational'],
      riskLevel: ['high'],
      search: 'demo'
    });
    console.log('Demo filters applied');
  } else {
    // Clear filters
    controlsStore.clearFilters();
    console.log('Filters cleared');
  }
};

// Initialize on mount
onMounted(async () => {
  // Initialize auth store
  await authStore.initialize();
  
  console.log('Store demo component initialized');
  console.log('Store health:', storeHealth.value);
});
</script>

<style scoped>
.store-demo {
  padding: 20px;
}

pre {
  background-color: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
  font-size: 12px;
  max-height: 300px;
  overflow-y: auto;
}

.v-chip-group {
  flex-wrap: wrap;
}
</style>