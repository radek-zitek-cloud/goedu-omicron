<template>
  <div class="design-system-showcase">
    <VContainer fluid>
      <!-- Header -->
      <div class="showcase-header">
        <h1 class="text-h3 mb-2">GoEdu Design System</h1>
        <p class="text-h6 text-medium-emphasis mb-6">
          Banking-specific UI components with WCAG 2.1 AA compliance
        </p>
      </div>

      <!-- Navigation -->
      <VTabs v-model="activeTab" class="mb-6">
        <VTab value="buttons">Buttons</VTab>
        <VTab value="inputs">Inputs</VTab>
        <VTab value="tables">Tables</VTab>
        <VTab value="workflows">Workflows</VTab>
        <VTab value="audit">Audit Trail</VTab>
      </VTabs>

      <!-- Content -->
      <VWindow v-model="activeTab">
        <!-- Buttons Showcase -->
        <VWindowItem value="buttons">
          <VCard>
            <VCardTitle>Banking Buttons</VCardTitle>
            <VCardText>
              <div class="showcase-section">
                <h3 class="text-h5 mb-4">Button Variants</h3>
                <div class="button-grid">
                  <BankingButton variant="primary">Primary Button</BankingButton>
                  <BankingButton variant="secondary">Secondary Button</BankingButton>
                  <BankingButton variant="approve" action-type="approve">Approve</BankingButton>
                  <BankingButton variant="reject" action-type="reject">Reject</BankingButton>
                  <BankingButton variant="warning">Warning</BankingButton>
                  <BankingButton variant="outlined">Outlined</BankingButton>
                  <BankingButton variant="text">Text Button</BankingButton>
                </div>
              </div>

              <div class="showcase-section">
                <h3 class="text-h5 mb-4">Button Sizes</h3>
                <div class="button-grid">
                  <BankingButton size="small">Small</BankingButton>
                  <BankingButton size="default">Default</BankingButton>
                  <BankingButton size="large">Large</BankingButton>
                  <BankingButton size="x-large">X-Large</BankingButton>
                </div>
              </div>

              <div class="showcase-section">
                <h3 class="text-h5 mb-4">Banking Actions</h3>
                <div class="button-grid">
                  <BankingButton
                    action-type="approve"
                    variant="approve"
                    :require-confirmation="true"
                    confirmation-message="Approve this transaction?"
                  >
                    Approve Transaction
                  </BankingButton>
                  
                  <BankingButton
                    action-type="reject"
                    variant="reject"
                    :require-confirmation="true"
                    confirmation-message="Reject this transaction?"
                  >
                    Reject Transaction
                  </BankingButton>
                  
                  <BankingButton
                    action-type="send-reminder"
                    variant="outlined"
                    icon="mdi-email"
                  >
                    Send Reminder
                  </BankingButton>
                  
                  <BankingButton
                    action-type="export"
                    variant="outlined"
                    icon="mdi-download"
                  >
                    Export Data
                  </BankingButton>
                </div>
              </div>

              <div class="showcase-section">
                <h3 class="text-h5 mb-4">Button States</h3>
                <div class="button-grid">
                  <BankingButton>Normal</BankingButton>
                  <BankingButton :loading="true">Loading</BankingButton>
                  <BankingButton :disabled="true">Disabled</BankingButton>
                  <BankingButton 
                    :disabled="true" 
                    disabled-reason="Insufficient permissions"
                  >
                    Disabled with Reason
                  </BankingButton>
                </div>
              </div>
            </VCardText>
          </VCard>
        </VWindowItem>

        <!-- Inputs Showcase -->
        <VWindowItem value="inputs">
          <VCard>
            <VCardTitle>Banking Inputs</VCardTitle>
            <VCardText>
              <div class="showcase-section">
                <h3 class="text-h5 mb-4">Input Types</h3>
                <VRow>
                  <VCol cols="12" md="6">
                    <BankingInput
                      v-model="sampleData.text"
                      type="text"
                      label="Text Input"
                      placeholder="Enter text"
                    />
                  </VCol>
                  <VCol cols="12" md="6">
                    <BankingInput
                      v-model="sampleData.email"
                      type="email"
                      label="Email Address"
                      :rules="[emailRule]"
                    />
                  </VCol>
                  <VCol cols="12" md="6">
                    <BankingInput
                      v-model="sampleData.currency"
                      type="currency"
                      label="Transaction Amount"
                      currency="USD"
                      :rules="[requiredRule, currencyRule]"
                      required
                    />
                  </VCol>
                  <VCol cols="12" md="6">
                    <BankingInput
                      v-model="sampleData.percentage"
                      type="percentage"
                      label="Interest Rate"
                      :rules="[percentageRule]"
                    />
                  </VCol>
                  <VCol cols="12" md="6">
                    <BankingInput
                      v-model="sampleData.phone"
                      type="phone"
                      label="Phone Number"
                      :auto-format="true"
                    />
                  </VCol>
                  <VCol cols="12" md="6">
                    <BankingInput
                      v-model="sampleData.account"
                      type="account-number"
                      label="Account Number"
                      :auto-format="true"
                    />
                  </VCol>
                </VRow>
              </div>

              <div class="showcase-section">
                <h3 class="text-h5 mb-4">Input Variants</h3>
                <VRow>
                  <VCol cols="12" md="4">
                    <BankingInput
                      v-model="sampleData.outlined"
                      variant="outlined"
                      label="Outlined"
                    />
                  </VCol>
                  <VCol cols="12" md="4">
                    <BankingInput
                      v-model="sampleData.filled"
                      variant="filled"
                      label="Filled"
                    />
                  </VCol>
                  <VCol cols="12" md="4">
                    <BankingInput
                      v-model="sampleData.underlined"
                      variant="underlined"
                      label="Underlined"
                    />
                  </VCol>
                </VRow>
              </div>
            </VCardText>
          </VCard>
        </VWindowItem>

        <!-- Tables Showcase -->
        <VWindowItem value="tables">
          <VCard>
            <VCardTitle>Banking Tables</VCardTitle>
            <VCardText>
              <BankingTable
                :headers="tableHeaders"
                :items="tableItems"
                title="Transaction History"
                subtitle="Recent banking transactions"
                searchable
                exportable
                striped
                hover
                @row-click="handleRowClick"
                @export="handleExport"
              />
            </VCardText>
          </VCard>
        </VWindowItem>

        <!-- Approval Workflow Showcase -->
        <VWindowItem value="workflows">
          <VCard>
            <VCardTitle>Approval Workflows</VCardTitle>
            <VCardText>
              <ApprovalWorkflow
                :steps="workflowSteps"
                :current-step="1"
                workflow-status="in-progress"
                title="Loan Approval Process"
                description="Multi-step approval process for loan applications"
                @step-approved="handleStepApproved"
                @step-rejected="handleStepRejected"
              />
            </VCardText>
          </VCard>
        </VWindowItem>

        <!-- Audit Trail Showcase -->
        <VWindowItem value="audit">
          <VCard>
            <VCardTitle>Audit Trail</VCardTitle>
            <VCardText>
              <AuditTrail
                :entries="auditEntries"
                title="System Audit Trail"
                subtitle="Recent system activities and changes"
                filterable
                exportable
                @filters-changed="handleFiltersChanged"
                @export-requested="handleAuditExport"
              />
            </VCardText>
          </VCard>
        </VWindowItem>
      </VWindow>
    </VContainer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import {
  VContainer,
  VCard,
  VCardTitle,
  VCardText,
  VTabs,
  VTab,
  VWindow,
  VWindowItem,
  VRow,
  VCol
} from 'vuetify/components'

// Import design system components
import BankingButton from '../base/BankingButton.vue'
import BankingInput from '../base/BankingInput.vue'
import BankingTable from '../base/BankingTable.vue'
import ApprovalWorkflow from '../banking/ApprovalWorkflow.vue'
import AuditTrail from '../banking/AuditTrail.vue'

// Import types
import type { BankingTableHeader, BankingTableItem } from '../base/BankingTable.types'
import type { ApprovalStep } from '../banking/ApprovalWorkflow.types'
import type { AuditEntry, AuditFilters } from '../banking/AuditTrail.types'

/**
 * Design System Showcase Component
 * 
 * Interactive showcase of all banking design system components
 * with examples, documentation, and testing capabilities.
 */

// Active tab state
const activeTab = ref('buttons')

// Sample data for inputs
const sampleData = reactive({
  text: '',
  email: '',
  currency: '',
  percentage: '',
  phone: '',
  account: '',
  outlined: '',
  filled: '',
  underlined: ''
})

// Validation rules
const requiredRule = (value: any) => !!value || 'Field is required'
const emailRule = (value: string) => {
  const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return !value || pattern.test(value) || 'Invalid email address'
}
const currencyRule = (value: string) => {
  const numValue = parseFloat(value.replace(/[^0-9.-]/g, ''))
  return !value || !isNaN(numValue) || 'Invalid currency amount'
}
const percentageRule = (value: string) => {
  const numValue = parseFloat(value)
  return !value || (!isNaN(numValue) && numValue >= 0 && numValue <= 100) || 'Invalid percentage'
}

// Table configuration
const tableHeaders: BankingTableHeader[] = [
  {
    key: 'id',
    title: 'Transaction ID',
    type: 'text',
    sortable: true
  },
  {
    key: 'date',
    title: 'Date',
    type: 'date',
    sortable: true
  },
  {
    key: 'amount',
    title: 'Amount',
    type: 'currency',
    format: {
      currency: {
        code: 'USD'
      }
    },
    sortable: true
  },
  {
    key: 'status',
    title: 'Status',
    type: 'status',
    format: {
      status: {
        colorMap: {
          'completed': 'success',
          'pending': 'warning',
          'failed': 'error'
        }
      }
    }
  },
  {
    key: 'type',
    title: 'Type',
    type: 'text'
  }
]

const tableItems: BankingTableItem[] = [
  {
    id: 'TXN-001',
    date: '2024-01-15',
    amount: 1250.00,
    status: 'completed',
    type: 'Transfer'
  },
  {
    id: 'TXN-002',
    date: '2024-01-14',
    amount: 850.50,
    status: 'pending',
    type: 'Payment'
  },
  {
    id: 'TXN-003',
    date: '2024-01-13',
    amount: 2100.75,
    status: 'completed',
    type: 'Deposit'
  },
  {
    id: 'TXN-004',
    date: '2024-01-12',
    amount: 500.00,
    status: 'failed',
    type: 'Withdrawal'
  }
]

// Workflow steps
const workflowSteps: ApprovalStep[] = [
  {
    id: 'step-1',
    title: 'Initial Review',
    description: 'Review application for completeness',
    type: 'review',
    status: 'completed',
    role: 'Loan Officer',
    assignee: 'John Smith',
    dueDate: '2024-01-10',
    allowApprove: true,
    allowReject: true,
    history: [
      {
        action: 'Approved',
        user: 'John Smith',
        timestamp: '2024-01-10T10:30:00Z',
        comment: 'Application complete and meets initial criteria'
      }
    ]
  },
  {
    id: 'step-2',
    title: 'Credit Check',
    description: 'Verify credit history and score',
    type: 'validate',
    status: 'in-progress',
    role: 'Credit Analyst',
    assignee: 'Sarah Johnson',
    dueDate: '2024-01-16',
    allowApprove: true,
    allowReject: true,
    allowDelegate: true
  },
  {
    id: 'step-3',
    title: 'Manager Approval',
    description: 'Final approval from branch manager',
    type: 'approve',
    status: 'pending',
    role: 'Branch Manager',
    assignee: 'Michael Davis',
    dueDate: '2024-01-18',
    allowApprove: true,
    allowReject: true
  }
]

// Audit entries
const auditEntries: AuditEntry[] = [
  {
    id: 'audit-001',
    action: 'Login',
    description: 'User logged into the system',
    user: 'john.smith@bank.com',
    timestamp: '2024-01-15T14:30:00Z',
    category: 'security',
    riskLevel: 'low',
    metadata: {
      ipAddress: '192.168.1.100',
      location: 'New York, NY',
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
    }
  },
  {
    id: 'audit-002',
    action: 'Transaction Approved',
    description: 'Large transaction approved',
    user: 'sarah.johnson@bank.com',
    timestamp: '2024-01-15T14:25:00Z',
    target: 'TXN-12345',
    category: 'financial',
    riskLevel: 'medium',
    details: {
      transactionId: 'TXN-12345',
      amount: 50000,
      approver: 'sarah.johnson@bank.com'
    },
    changes: [
      {
        field: 'status',
        oldValue: 'pending',
        newValue: 'approved'
      },
      {
        field: 'approvedBy',
        oldValue: null,
        newValue: 'sarah.johnson@bank.com'
      }
    ],
    metadata: {
      ipAddress: '192.168.1.105',
      location: 'Chicago, IL'
    }
  },
  {
    id: 'audit-003',
    action: 'Account Modified',
    description: 'Customer account information updated',
    user: 'admin@bank.com',
    timestamp: '2024-01-15T14:20:00Z',
    target: 'ACC-67890',
    category: 'data',
    riskLevel: 'high',
    changes: [
      {
        field: 'address',
        oldValue: '123 Old Street',
        newValue: '456 New Avenue'
      },
      {
        field: 'phone',
        oldValue: '555-0123',
        newValue: '555-0456'
      }
    ],
    metadata: {
      ipAddress: '192.168.1.10',
      location: 'Head Office'
    }
  }
]

// Event handlers
function handleRowClick(item: BankingTableItem): void {
  console.log('Row clicked:', item)
}

function handleExport(items: BankingTableItem[]): void {
  console.log('Export requested:', items)
}

function handleStepApproved(event: any): void {
  console.log('Step approved:', event)
}

function handleStepRejected(event: any): void {
  console.log('Step rejected:', event)
}

function handleFiltersChanged(filters: AuditFilters): void {
  console.log('Filters changed:', filters)
}

function handleAuditExport(request: any): void {
  console.log('Audit export requested:', request)
}
</script>

<style scoped>
.design-system-showcase {
  min-height: 100vh;
  background: rgb(var(--v-theme-background));
}

.showcase-header {
  text-align: center;
  padding: 48px 0;
  border-bottom: 1px solid rgb(var(--v-theme-border));
  margin-bottom: 32px;
}

.showcase-section {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid rgba(var(--v-theme-border), 0.3);
}

.showcase-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.button-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: center;
}

@media (max-width: 599px) {
  .button-grid {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>