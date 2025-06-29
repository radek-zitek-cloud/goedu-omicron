/**
 * Controls Store for GoEdu Omicron Banking Platform
 *
 * Centralized state management for banking controls, including control
 * definitions, test executions, results, and compliance tracking.
 *
 * Features:
 * - Control CRUD operations with full lifecycle management
 * - Test execution and result tracking
 * - Evidence management and file attachments
 * - Compliance reporting and analytics
 * - Real-time collaboration and updates
 * - Offline capability with sync when online
 * - Audit trail for all control activities
 *
 * Banking Control Categories:
 * - Operational Controls: Daily operations and procedures
 * - Financial Controls: Financial reporting and accuracy
 * - Compliance Controls: Regulatory compliance requirements
 * - Security Controls: Information security and access
 * - Risk Controls: Risk management and mitigation
 *
 * @author GoEdu Development Team
 * @version 1.0.0
 * @since 2024
 */

import { defineStore } from 'pinia';
import { ref, computed, readonly } from 'vue';
import axios from 'axios';

/**
 * Control Status Enumeration
 */
export type ControlStatus =
  | 'draft' // Control is being created/edited
  | 'active' // Control is active and being tested
  | 'inactive' // Control is temporarily disabled
  | 'retired' // Control is no longer used
  | 'under_review' // Control is being reviewed/updated
  | 'deprecated'; // Control is deprecated but kept for history

/**
 * Test Status Enumeration
 */
export type TestStatus =
  | 'not_started' // Test not yet started
  | 'in_progress' // Test is currently being executed
  | 'completed' // Test completed successfully
  | 'failed' // Test failed or found issues
  | 'cancelled' // Test was cancelled
  | 'on_hold' // Test is temporarily on hold
  | 'needs_review'; // Test results need review

/**
 * Test Result Enumeration
 */
export type TestResult =
  | 'pass' // Control is operating effectively
  | 'fail' // Control has deficiencies
  | 'exception' // Control has exceptions that need addressing
  | 'not_tested' // Control was not tested in this period
  | 'not_applicable'; // Control is not applicable for this period

/**
 * Control Frequency Enumeration
 */
export type ControlFrequency =
  | 'daily'
  | 'weekly'
  | 'monthly'
  | 'quarterly'
  | 'semi_annual'
  | 'annual'
  | 'ad_hoc'
  | 'continuous';

/**
 * Control Interface Definition
 *
 * Represents a banking control with all necessary properties for
 * compliance testing and management.
 */
export interface Control {
  id: string;
  organizationId: string;
  code: string;
  title: string;
  description: string;
  objective: string;
  category: 'operational' | 'financial' | 'compliance' | 'security' | 'risk';
  subcategory: string;
  riskLevel: 'low' | 'medium' | 'high' | 'critical';
  status: ControlStatus;
  frequency: ControlFrequency;
  owner: {
    userId: string;
    name: string;
    email: string;
    role: string;
  };
  reviewer: {
    userId: string;
    name: string;
    email: string;
    role: string;
  };
  regulatoryReferences: string[];
  frameworkMappings: {
    framework: string;
    controlId: string;
    requirement: string;
  }[];
  testProcedures: {
    id: string;
    step: number;
    description: string;
    expectedResult: string;
    automatable: boolean;
  }[];
  keyAttributes: {
    name: string;
    value: string;
    dataType: 'text' | 'number' | 'date' | 'boolean' | 'list';
    required: boolean;
  }[];
  tags: string[];
  relatedControls: string[];
  effectiveDate: Date;
  nextTestDate: Date;
  lastTestDate: Date | null;
  testHistory: ControlTest[];
  metrics: {
    totalTests: number;
    passRate: number;
    avgTestDuration: number;
    lastFailureDate: Date | null;
    consecutivePasses: number;
  };
  createdAt: Date;
  updatedAt: Date;
  createdBy: string;
  updatedBy: string;
}

/**
 * Control Test Interface
 *
 * Represents a test execution for a specific control.
 */
export interface ControlTest {
  id: string;
  controlId: string;
  organizationId: string;
  testPeriod: {
    startDate: Date;
    endDate: Date;
    label: string;
  };
  status: TestStatus;
  result: TestResult;
  tester: {
    userId: string;
    name: string;
    email: string;
  };
  reviewer: {
    userId: string;
    name: string;
    email: string;
  } | null;
  plannedStartDate: Date;
  actualStartDate: Date | null;
  plannedEndDate: Date;
  actualEndDate: Date | null;
  testSteps: {
    stepId: string;
    description: string;
    status: 'pending' | 'in_progress' | 'completed' | 'failed' | 'skipped';
    result: TestResult | null;
    notes: string;
    evidence: TestEvidence[];
    completedAt: Date | null;
    completedBy: string | null;
  }[];
  findings: {
    id: string;
    severity: 'low' | 'medium' | 'high' | 'critical';
    category: string;
    description: string;
    impact: string;
    recommendation: string;
    status: 'open' | 'in_progress' | 'resolved' | 'accepted';
    assignee: string;
    dueDate: Date;
    evidence: TestEvidence[];
  }[];
  summary: {
    overview: string;
    conclusion: string;
    recommendations: string[];
    nextSteps: string[];
  };
  evidence: TestEvidence[];
  comments: {
    id: string;
    userId: string;
    userName: string;
    message: string;
    timestamp: Date;
    edited: boolean;
  }[];
  createdAt: Date;
  updatedAt: Date;
  createdBy: string;
  updatedBy: string;
}

/**
 * Test Evidence Interface
 *
 * Represents evidence/documentation attached to tests.
 */
export interface TestEvidence {
  id: string;
  fileName: string;
  originalName: string;
  fileSize: number;
  mimeType: string;
  description: string;
  category: 'screenshot' | 'document' | 'data_export' | 'email' | 'other';
  uploadedAt: Date;
  uploadedBy: string;
  tags: string[];
  url: string;
  downloadUrl: string;
}

/**
 * Control Filter Interface
 */
export interface ControlFilter {
  status?: ControlStatus[];
  category?: string[];
  riskLevel?: string[];
  owner?: string[];
  tags?: string[];
  search?: string;
  dateRange?: {
    field: 'nextTestDate' | 'lastTestDate' | 'createdAt' | 'effectiveDate';
    start: Date;
    end: Date;
  };
}

/**
 * Control Creation Request Interface
 */
export interface CreateControlRequest {
  code: string;
  title: string;
  description: string;
  objective: string;
  category: Control['category'];
  subcategory: string;
  riskLevel: Control['riskLevel'];
  frequency: ControlFrequency;
  ownerId: string;
  reviewerId: string;
  regulatoryReferences: string[];
  frameworkMappings: Control['frameworkMappings'];
  testProcedures: Omit<Control['testProcedures'][0], 'id'>[];
  keyAttributes: Control['keyAttributes'];
  tags: string[];
  relatedControls: string[];
  effectiveDate: Date;
}

/**
 * Test Execution Request Interface
 */
export interface ExecuteTestRequest {
  controlId: string;
  testPeriod: {
    startDate: Date;
    endDate: Date;
    label: string;
  };
  plannedStartDate: Date;
  plannedEndDate: Date;
  testerId: string;
  reviewerId?: string;
  notes?: string;
}

/**
 * Controls Store Implementation
 *
 * Main store for managing controls and test execution state.
 * Uses Pinia composition API pattern for better type safety and reactivity.
 */
export const useControlsStore = defineStore(
  'controls',
  () => {
    // Reactive state
    const controls = ref<Control[]>([]);
    const currentControl = ref<Control | null>(null);
    const tests = ref<ControlTest[]>([]);
    const currentTest = ref<ControlTest | null>(null);
    const isLoading = ref(false);
    const isLoadingTests = ref(false);
    const isExecutingTest = ref(false);
    const error = ref<string | null>(null);
    const filters = ref<ControlFilter>({});
    const lastSyncTime = ref<Date | null>(null);
    const offlineQueue = ref<any[]>([]); // Queue for offline operations

    // Configuration
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api';

    // Computed properties
    const filteredControls = computed(() => {
      let filtered = [...controls.value];
      const filter = filters.value;

      // Apply status filter
      if (filter.status?.length) {
        filtered = filtered.filter(control => filter.status!.includes(control.status));
      }

      // Apply category filter
      if (filter.category?.length) {
        filtered = filtered.filter(control => filter.category!.includes(control.category));
      }

      // Apply risk level filter
      if (filter.riskLevel?.length) {
        filtered = filtered.filter(control => filter.riskLevel!.includes(control.riskLevel));
      }

      // Apply owner filter
      if (filter.owner?.length) {
        filtered = filtered.filter(control => filter.owner!.includes(control.owner.userId));
      }

      // Apply tags filter
      if (filter.tags?.length) {
        filtered = filtered.filter(control => filter.tags!.some(tag => control.tags.includes(tag)));
      }

      // Apply search filter
      if (filter.search?.trim()) {
        const searchTerm = filter.search.toLowerCase();
        filtered = filtered.filter(
          control =>
            control.title.toLowerCase().includes(searchTerm) ||
            control.description.toLowerCase().includes(searchTerm) ||
            control.code.toLowerCase().includes(searchTerm) ||
            control.tags.some(tag => tag.toLowerCase().includes(searchTerm))
        );
      }

      // Apply date range filter
      if (filter.dateRange) {
        const { field, start, end } = filter.dateRange;
        filtered = filtered.filter(control => {
          const date = control[field] as Date;
          return date >= start && date <= end;
        });
      }

      return filtered;
    });

    const controlsByCategory = computed(() => {
      const grouped: Record<string, Control[]> = {};
      filteredControls.value.forEach(control => {
        if (!grouped[control.category]) {
          grouped[control.category] = [];
        }
        grouped[control.category].push(control);
      });
      return grouped;
    });

    const controlsByStatus = computed(() => {
      const grouped: Record<ControlStatus, Control[]> = {
        draft: [],
        active: [],
        inactive: [],
        retired: [],
        under_review: [],
        deprecated: [],
      };
      filteredControls.value.forEach(control => {
        grouped[control.status].push(control);
      });
      return grouped;
    });

    const upcomingTests = computed(() => {
      const now = new Date();
      const thirtyDaysFromNow = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000);

      return filteredControls.value
        .filter(control => control.nextTestDate <= thirtyDaysFromNow && control.status === 'active')
        .sort((a, b) => a.nextTestDate.getTime() - b.nextTestDate.getTime());
    });

    const overdueTests = computed(() => {
      const now = new Date();
      return filteredControls.value
        .filter(control => control.nextTestDate < now && control.status === 'active')
        .sort((a, b) => a.nextTestDate.getTime() - b.nextTestDate.getTime());
    });

    const controlsCount = computed(() => ({
      total: controls.value.length,
      active: controls.value.filter(c => c.status === 'active').length,
      draft: controls.value.filter(c => c.status === 'draft').length,
      overdue: overdueTests.value.length,
      upcoming: upcomingTests.value.length,
    }));

    /**
     * Load Controls
     *
     * Fetches controls for the current organization.
     *
     * @param organizationId - ID of the organization
     * @returns Promise<boolean> - Success status
     */
    async function loadControls(organizationId: string): Promise<boolean> {
      if (isLoading.value) return false;

      try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.get(
          `${API_BASE_URL}/organizations/${organizationId}/controls`,
          {
            withCredentials: true,
          }
        );

        if (response.data.success) {
          controls.value = response.data.controls.map((control: Record<string, any>) => ({
            ...control,
            effectiveDate: new Date(control.effectiveDate),
            nextTestDate: new Date(control.nextTestDate),
            lastTestDate: control.lastTestDate ? new Date(control.lastTestDate) : null,
            createdAt: new Date(control.createdAt),
            updatedAt: new Date(control.updatedAt),
            testHistory: control.testHistory.map((test: Record<string, any>) => ({
              ...test,
              testPeriod: {
                ...test.testPeriod,
                startDate: new Date(test.testPeriod.startDate),
                endDate: new Date(test.testPeriod.endDate),
              },
              plannedStartDate: new Date(test.plannedStartDate),
              actualStartDate: test.actualStartDate ? new Date(test.actualStartDate) : null,
              plannedEndDate: new Date(test.plannedEndDate),
              actualEndDate: test.actualEndDate ? new Date(test.actualEndDate) : null,
              createdAt: new Date(test.createdAt),
              updatedAt: new Date(test.updatedAt),
            })),
          }));

          lastSyncTime.value = new Date();
          console.log('✅ Controls loaded successfully:', controls.value.length);
          return true;
        }

        throw new Error(response.data.message || 'Failed to load controls');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to load controls:', errorMessage);
        return false;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Create Control
     *
     * Creates a new control with the provided details.
     *
     * @param request - Control creation request data
     * @returns Promise<Control | null> - Created control or null on failure
     */
    async function createControl(request: CreateControlRequest): Promise<Control | null> {
      try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.post(`${API_BASE_URL}/controls`, request, {
          withCredentials: true,
        });

        if (response.data.success) {
          const newControl: Control = {
            ...response.data.control,
            effectiveDate: new Date(response.data.control.effectiveDate),
            nextTestDate: new Date(response.data.control.nextTestDate),
            lastTestDate: response.data.control.lastTestDate
              ? new Date(response.data.control.lastTestDate)
              : null,
            createdAt: new Date(response.data.control.createdAt),
            updatedAt: new Date(response.data.control.updatedAt),
            testHistory: [],
          };

          controls.value.push(newControl);
          console.log('✅ Control created successfully:', newControl.code);
          return newControl;
        }

        throw new Error(response.data.message || 'Failed to create control');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to create control:', errorMessage);
        return null;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Update Control
     *
     * Updates an existing control with new data.
     *
     * @param controlId - ID of the control to update
     * @param updates - Partial control data to update
     * @returns Promise<boolean> - Success status
     */
    async function updateControl(controlId: string, updates: Partial<Control>): Promise<boolean> {
      try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.patch(`${API_BASE_URL}/controls/${controlId}`, updates, {
          withCredentials: true,
        });

        if (response.data.success) {
          const updatedControl = response.data.control;
          const index = controls.value.findIndex(control => control.id === controlId);

          if (index !== -1) {
            controls.value[index] = {
              ...updatedControl,
              effectiveDate: new Date(updatedControl.effectiveDate),
              nextTestDate: new Date(updatedControl.nextTestDate),
              lastTestDate: updatedControl.lastTestDate
                ? new Date(updatedControl.lastTestDate)
                : null,
              createdAt: new Date(updatedControl.createdAt),
              updatedAt: new Date(updatedControl.updatedAt),
            };

            // Update current control if it's the one being updated
            if (currentControl.value?.id === controlId) {
              currentControl.value = controls.value[index];
            }
          }

          console.log('✅ Control updated successfully:', controlId);
          return true;
        }

        throw new Error(response.data.message || 'Failed to update control');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to update control:', errorMessage);
        return false;
      } finally {
        isLoading.value = false;
      }
    }

    /**
     * Execute Test
     *
     * Starts a new test execution for a control.
     *
     * @param request - Test execution request data
     * @returns Promise<ControlTest | null> - Created test or null on failure
     */
    async function executeTest(request: ExecuteTestRequest): Promise<ControlTest | null> {
      try {
        isExecutingTest.value = true;
        error.value = null;

        const response = await axios.post(`${API_BASE_URL}/controls/tests`, request, {
          withCredentials: true,
        });

        if (response.data.success) {
          const newTest: ControlTest = {
            ...response.data.test,
            testPeriod: {
              ...response.data.test.testPeriod,
              startDate: new Date(response.data.test.testPeriod.startDate),
              endDate: new Date(response.data.test.testPeriod.endDate),
            },
            plannedStartDate: new Date(response.data.test.plannedStartDate),
            actualStartDate: response.data.test.actualStartDate
              ? new Date(response.data.test.actualStartDate)
              : null,
            plannedEndDate: new Date(response.data.test.plannedEndDate),
            actualEndDate: response.data.test.actualEndDate
              ? new Date(response.data.test.actualEndDate)
              : null,
            createdAt: new Date(response.data.test.createdAt),
            updatedAt: new Date(response.data.test.updatedAt),
          };

          tests.value.push(newTest);
          currentTest.value = newTest;

          // Update the control's test history
          const controlIndex = controls.value.findIndex(c => c.id === request.controlId);
          if (controlIndex !== -1) {
            controls.value[controlIndex].testHistory.push(newTest);
          }

          console.log('✅ Test execution started successfully:', newTest.id);
          return newTest;
        }

        throw new Error(response.data.message || 'Failed to start test execution');
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        error.value = errorMessage;
        console.error('❌ Failed to start test execution:', errorMessage);
        return null;
      } finally {
        isExecutingTest.value = false;
      }
    }

    /**
     * Set Current Control
     *
     * Sets the active control for detailed view/editing.
     *
     * @param controlId - ID of the control to set as current
     */
    function setCurrentControl(controlId: string): void {
      const control = controls.value.find(c => c.id === controlId);
      currentControl.value = control || null;
    }

    /**
     * Set Filters
     *
     * Updates the control filters.
     *
     * @param newFilters - Filter criteria to apply
     */
    function setFilters(newFilters: Partial<ControlFilter>): void {
      filters.value = { ...filters.value, ...newFilters };
    }

    /**
     * Clear Filters
     *
     * Resets all filters to default state.
     */
    function clearFilters(): void {
      filters.value = {};
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
      controls.value = [];
      currentControl.value = null;
      tests.value = [];
      currentTest.value = null;
      isLoading.value = false;
      isLoadingTests.value = false;
      isExecutingTest.value = false;
      error.value = null;
      filters.value = {};
      lastSyncTime.value = null;
      offlineQueue.value = [];
    }

    // Return store interface
    return {
      // State
      controls: readonly(controls),
      currentControl: readonly(currentControl),
      tests: readonly(tests),
      currentTest: readonly(currentTest),
      isLoading: readonly(isLoading),
      isLoadingTests: readonly(isLoadingTests),
      isExecutingTest: readonly(isExecutingTest),
      error: readonly(error),
      filters: readonly(filters),
      lastSyncTime: readonly(lastSyncTime),
      offlineQueue: readonly(offlineQueue),

      // Computed
      filteredControls,
      controlsByCategory,
      controlsByStatus,
      upcomingTests,
      overdueTests,
      controlsCount,

      // Actions
      loadControls,
      createControl,
      updateControl,
      executeTest,
      setCurrentControl,
      setFilters,
      clearFilters,
      clearError,
      resetStore,
    };
  },
  {
    // Pinia persistence configuration for offline capability
    persist: {
      key: 'goedu-controls',
      storage: localStorage,
      paths: ['controls', 'filters', 'lastSyncTime', 'offlineQueue'],
    },
  }
);
