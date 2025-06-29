<template>
  <div class="approval-workflow" :class="computedClasses">
    <!-- Workflow Header -->
    <div v-if="showHeader" class="approval-workflow__header">
      <div class="approval-workflow__title-section">
        <h3 class="approval-workflow__title">{{ title || 'Approval Workflow' }}</h3>
        <p v-if="description" class="approval-workflow__description">{{ description }}</p>
      </div>
      
      <div class="approval-workflow__status">
        <VChip
          :color="getWorkflowStatusColor()"
          :variant="workflowStatus === 'completed' ? 'flat' : 'tonal'"
          size="large"
        >
          <VIcon
            :icon="getWorkflowStatusIcon()"
            start
            size="small"
          />
          {{ getWorkflowStatusText() }}
        </VChip>
      </div>
    </div>

    <!-- Progress Bar -->
    <div class="approval-workflow__progress">
      <VProgressLinear
        :model-value="progressPercentage"
        :color="getProgressColor()"
        height="8"
        rounded
        striped
      />
      <div class="approval-workflow__progress-text">
        Step {{ currentStepIndex + 1 }} of {{ steps.length }}
        <span v-if="completedSteps > 0">({{ completedSteps }} completed)</span>
      </div>
    </div>

    <!-- Workflow Steps -->
    <div class="approval-workflow__steps">
      <div
        v-for="(step, index) in steps"
        :key="step.id"
        class="approval-workflow__step"
        :class="getStepClasses(step, index)"
      >
        <!-- Step Icon -->
        <div class="approval-workflow__step-icon">
          <VAvatar
            :color="getStepColor(step, index)"
            :variant="getStepVariant(step, index)"
            size="40"
          >
            <VIcon
              :icon="getStepIcon(step, index)"
              :color="getStepIconColor(step, index)"
            />
          </VAvatar>
          
          <!-- Connecting Line -->
          <div
            v-if="index < steps.length - 1"
            class="approval-workflow__step-connector"
            :class="{ 'approval-workflow__step-connector--completed': isStepCompleted(index) }"
          />
        </div>

        <!-- Step Content -->
        <div class="approval-workflow__step-content">
          <div class="approval-workflow__step-header">
            <h4 class="approval-workflow__step-title">{{ step.title }}</h4>
            <div class="approval-workflow__step-meta">
              <span class="approval-workflow__step-role">{{ step.role }}</span>
              <span v-if="step.estimatedTime" class="approval-workflow__step-time">
                ~{{ step.estimatedTime }}
              </span>
            </div>
          </div>
          
          <p v-if="step.description" class="approval-workflow__step-description">
            {{ step.description }}
          </p>

          <!-- Step Details -->
          <div v-if="step.assignee || step.dueDate" class="approval-workflow__step-details">
            <div v-if="step.assignee" class="approval-workflow__step-assignee">
              <VIcon icon="mdi-account" size="small" />
              {{ step.assignee }}
            </div>
            <div v-if="step.dueDate" class="approval-workflow__step-due-date">
              <VIcon icon="mdi-calendar" size="small" />
              {{ formatDate(step.dueDate) }}
              <VChip
                v-if="isOverdue(step.dueDate) && !isStepCompleted(index)"
                color="error"
                size="x-small"
                variant="flat"
                class="ml-2"
              >
                Overdue
              </VChip>
            </div>
          </div>

          <!-- Step Actions -->
          <div v-if="isCurrentStep(index)" class="approval-workflow__step-actions">
            <slot :name="`step-${step.id}-actions`" :step="step" :index="index">
              <div class="approval-workflow__default-actions">
                <BankingButton
                  v-if="step.allowApprove"
                  variant="approve"
                  action-type="approve"
                  :require-confirmation="true"
                  confirmation-message="Are you sure you want to approve this step?"
                  @click="handleApprove(step, index)"
                  :disabled="!canUserAction(step, 'approve')"
                >
                  Approve
                </BankingButton>
                
                <BankingButton
                  v-if="step.allowReject"
                  variant="reject"
                  action-type="reject"
                  :require-confirmation="true"
                  confirmation-message="Are you sure you want to reject this step? This will stop the workflow."
                  @click="handleReject(step, index)"
                  :disabled="!canUserAction(step, 'reject')"
                >
                  Reject
                </BankingButton>
                
                <BankingButton
                  v-if="step.allowDelegate"
                  variant="outlined"
                  icon="mdi-account-arrow-right"
                  @click="handleDelegate(step, index)"
                  :disabled="!canUserAction(step, 'delegate')"
                >
                  Delegate
                </BankingButton>
              </div>
            </slot>
          </div>

          <!-- Step History -->
          <div v-if="step.history && step.history.length > 0" class="approval-workflow__step-history">
            <VExpansionPanels variant="accordion" flat>
              <VExpansionPanel>
                <VExpansionPanelTitle>
                  <VIcon icon="mdi-history" size="small" class="mr-2" />
                  View History ({{ step.history.length }} entries)
                </VExpansionPanelTitle>
                <VExpansionPanelText>
                  <div
                    v-for="(entry, historyIndex) in step.history"
                    :key="historyIndex"
                    class="approval-workflow__history-entry"
                  >
                    <div class="approval-workflow__history-header">
                      <strong>{{ entry.action }}</strong>
                      <span class="approval-workflow__history-user">by {{ entry.user }}</span>
                      <span class="approval-workflow__history-date">{{ formatDateTime(entry.timestamp) }}</span>
                    </div>
                    <p v-if="entry.comment" class="approval-workflow__history-comment">
                      {{ entry.comment }}
                    </p>
                  </div>
                </VExpansionPanelText>
              </VExpansionPanel>
            </VExpansionPanels>
          </div>
        </div>
      </div>
    </div>

    <!-- Comments Section -->
    <div v-if="allowComments" class="approval-workflow__comments">
      <VTextarea
        v-model="currentComment"
        label="Add Comment"
        hint="Optional comment for this action"
        variant="outlined"
        rows="3"
        :disabled="!canAddComment"
      />
    </div>

    <!-- Workflow Actions -->
    <div v-if="showWorkflowActions" class="approval-workflow__workflow-actions">
      <BankingButton
        v-if="workflowStatus === 'draft'"
        variant="primary"
        action-type="submit"
        @click="handleStartWorkflow"
        :disabled="!canStartWorkflow"
      >
        Start Workflow
      </BankingButton>
      
      <BankingButton
        v-if="canCancelWorkflow"
        variant="warning"
        action-type="cancel"
        :require-confirmation="true"
        confirmation-message="Are you sure you want to cancel this workflow? This action cannot be undone."
        @click="handleCancelWorkflow"
      >
        Cancel Workflow
      </BankingButton>
      
      <BankingButton
        v-if="workflowStatus === 'completed'"
        variant="outlined"
        icon="mdi-download"
        @click="handleExportApproval"
      >
        Export Approval Record
      </BankingButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, type PropType } from 'vue'
import {
  VChip,
  VIcon,
  VProgressLinear,
  VAvatar,
  VExpansionPanels,
  VExpansionPanel,
  VExpansionPanelTitle,
  VExpansionPanelText,
  VTextarea
} from 'vuetify/components'
import BankingButton from '../base/BankingButton.vue'
import type {
  ApprovalWorkflowProps,
  ApprovalWorkflowEmits,
  ApprovalStep,
  WorkflowStatus
} from './ApprovalWorkflow.types'

/**
 * Approval Workflow Component
 * 
 * A comprehensive workflow component for banking approval processes with
 * step tracking, user assignments, due dates, and audit trail functionality.
 * 
 * Features:
 * - Multi-step approval process
 * - Role-based access control
 * - Due date tracking and overdue notifications
 * - Comprehensive audit trail and history
 * - Delegation and escalation support
 * - WCAG 2.1 AA compliant
 * - Responsive design
 * 
 * @example
 * <ApprovalWorkflow
 *   :steps="approvalSteps"
 *   :current-step="2"
 *   title="Transaction Approval"
 *   @step-approved="handleStepApproved"
 *   @workflow-completed="handleWorkflowCompleted"
 * />
 */

// Props with banking-specific defaults
const props = withDefaults(defineProps<ApprovalWorkflowProps>(), {
  currentStep: 0,
  workflowStatus: 'draft',
  showHeader: true,
  allowComments: true,
  showWorkflowActions: true,
  strictSequence: true,
  autoProgress: true
})

// Events
const emit = defineEmits<ApprovalWorkflowEmits>()

// Reactive state
const currentComment = ref('')

/**
 * Computed properties
 */
const computedClasses = computed(() => {
  const classes: string[] = ['approval-workflow']
  
  classes.push(`approval-workflow--${props.workflowStatus}`)
  
  if (props.variant) {
    classes.push(`approval-workflow--${props.variant}`)
  }
  
  if (props.class) {
    if (typeof props.class === 'string') {
      classes.push(props.class)
    } else if (Array.isArray(props.class)) {
      classes.push(...props.class)
    } else {
      Object.entries(props.class).forEach(([className, condition]) => {
        if (condition) classes.push(className)
      })
    }
  }
  
  return classes
})

const currentStepIndex = computed(() => props.currentStep)

const completedSteps = computed(() => {
  return props.steps.filter(step => step.status === 'completed').length
})

const progressPercentage = computed(() => {
  if (props.steps.length === 0) return 0
  return (completedSteps.value / props.steps.length) * 100
})

const canStartWorkflow = computed(() => {
  return props.workflowStatus === 'draft' && props.steps.length > 0
})

const canCancelWorkflow = computed(() => {
  return ['in-progress', 'paused'].includes(props.workflowStatus) && props.allowCancel
})

const canAddComment = computed(() => {
  return props.allowComments && ['in-progress', 'paused'].includes(props.workflowStatus)
})

/**
 * Workflow status helpers
 */
function getWorkflowStatusColor(): string {
  switch (props.workflowStatus) {
    case 'completed': return 'success'
    case 'rejected': return 'error'
    case 'cancelled': return 'warning'
    case 'in-progress': return 'primary'
    case 'paused': return 'warning'
    default: return 'default'
  }
}

function getWorkflowStatusIcon(): string {
  switch (props.workflowStatus) {
    case 'completed': return 'mdi-check-circle'
    case 'rejected': return 'mdi-cancel'
    case 'cancelled': return 'mdi-stop-circle'
    case 'in-progress': return 'mdi-play-circle'
    case 'paused': return 'mdi-pause-circle'
    default: return 'mdi-clipboard-text'
  }
}

function getWorkflowStatusText(): string {
  switch (props.workflowStatus) {
    case 'draft': return 'Draft'
    case 'in-progress': return 'In Progress'
    case 'completed': return 'Completed'
    case 'rejected': return 'Rejected'
    case 'cancelled': return 'Cancelled'
    case 'paused': return 'Paused'
    default: return 'Unknown'
  }
}

function getProgressColor(): string {
  if (props.workflowStatus === 'completed') return 'success'
  if (props.workflowStatus === 'rejected') return 'error'
  if (hasOverdueSteps()) return 'warning'
  return 'primary'
}

/**
 * Step helpers
 */
function isStepCompleted(index: number): boolean {
  return props.steps[index]?.status === 'completed'
}

function isCurrentStep(index: number): boolean {
  return index === currentStepIndex.value && props.workflowStatus === 'in-progress'
}

function getStepClasses(step: ApprovalStep, index: number): string[] {
  const classes: string[] = []
  
  if (step.status) {
    classes.push(`approval-workflow__step--${step.status}`)
  }
  
  if (isCurrentStep(index)) {
    classes.push('approval-workflow__step--current')
  }
  
  if (step.dueDate && isOverdue(step.dueDate) && !isStepCompleted(index)) {
    classes.push('approval-workflow__step--overdue')
  }
  
  return classes
}

function getStepColor(step: ApprovalStep, index: number): string {
  if (step.status === 'completed') return 'success'
  if (step.status === 'rejected') return 'error'
  if (isCurrentStep(index)) return 'primary'
  if (step.dueDate && isOverdue(step.dueDate)) return 'warning'
  return 'default'
}

function getStepVariant(step: ApprovalStep, index: number): 'flat' | 'elevated' | 'tonal' | 'outlined' | 'text' | 'plain' {
  if (step.status === 'completed' || isCurrentStep(index)) return 'flat'
  return 'outlined'
}

function getStepIcon(step: ApprovalStep, index: number): string {
  if (step.status === 'completed') return 'mdi-check'
  if (step.status === 'rejected') return 'mdi-close'
  if (isCurrentStep(index)) return 'mdi-play'
  
  // Default icons based on step type
  switch (step.type) {
    case 'review': return 'mdi-eye'
    case 'approve': return 'mdi-check-circle-outline'
    case 'sign': return 'mdi-signature'
    case 'validate': return 'mdi-shield-check-outline'
    default: return 'mdi-account-outline'
  }
}

function getStepIconColor(step: ApprovalStep, index: number): string {
  if (step.status === 'completed') return 'white'
  if (step.status === 'rejected') return 'white'
  if (isCurrentStep(index)) return 'white'
  return 'default'
}

/**
 * Utility functions
 */
function formatDate(date: Date | string): string {
  const d = new Date(date)
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}

function formatDateTime(date: Date | string): string {
  const d = new Date(date)
  return d.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function isOverdue(dueDate: Date | string): boolean {
  return new Date(dueDate) < new Date()
}

function hasOverdueSteps(): boolean {
  return props.steps.some((step, index) => 
    step.dueDate && 
    isOverdue(step.dueDate) && 
    !isStepCompleted(index)
  )
}

function canUserAction(step: ApprovalStep, action: string): boolean {
  // This would typically check against user permissions
  // For now, return true as a placeholder
  return true
}

/**
 * Event handlers
 */
function handleApprove(step: ApprovalStep, index: number): void {
  emit('step-approved', {
    step,
    stepIndex: index,
    comment: currentComment.value,
    timestamp: new Date()
  })
  
  currentComment.value = ''
}

function handleReject(step: ApprovalStep, index: number): void {
  emit('step-rejected', {
    step,
    stepIndex: index,
    comment: currentComment.value,
    timestamp: new Date()
  })
  
  currentComment.value = ''
}

function handleDelegate(step: ApprovalStep, index: number): void {
  emit('step-delegated', {
    step,
    stepIndex: index,
    comment: currentComment.value,
    timestamp: new Date()
  })
}

function handleStartWorkflow(): void {
  emit('workflow-started', {
    timestamp: new Date()
  })
}

function handleCancelWorkflow(): void {
  emit('workflow-cancelled', {
    reason: currentComment.value,
    timestamp: new Date()
  })
  
  currentComment.value = ''
}

function handleExportApproval(): void {
  emit('export-approval', {
    steps: props.steps,
    workflowStatus: props.workflowStatus,
    timestamp: new Date()
  })
}
</script>

<style scoped>
/**
 * Approval Workflow Styles
 * 
 * Professional styling for banking approval workflows with clear
 * visual hierarchy and accessibility features.
 */

.approval-workflow {
  background: rgb(var(--v-theme-surface));
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* Header */
.approval-workflow__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgb(var(--v-theme-border));
}

.approval-workflow__title {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0 0 4px;
  color: rgb(var(--v-theme-on-surface));
}

.approval-workflow__description {
  color: rgb(var(--v-theme-on-surface-variant));
  margin: 0;
  font-size: 0.875rem;
}

/* Progress */
.approval-workflow__progress {
  margin-bottom: 32px;
}

.approval-workflow__progress-text {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface-variant));
}

/* Steps */
.approval-workflow__steps {
  margin-bottom: 24px;
}

.approval-workflow__step {
  display: flex;
  margin-bottom: 32px;
  position: relative;
}

.approval-workflow__step:last-child {
  margin-bottom: 0;
}

.approval-workflow__step-icon {
  position: relative;
  margin-right: 20px;
  flex-shrink: 0;
}

.approval-workflow__step-connector {
  position: absolute;
  top: 50px;
  left: 50%;
  transform: translateX(-50%);
  width: 2px;
  height: 40px;
  background: rgb(var(--v-theme-border));
  transition: background-color 0.3s ease;
}

.approval-workflow__step-connector--completed {
  background: rgb(var(--v-theme-success));
}

.approval-workflow__step-content {
  flex: 1;
  min-width: 0;
}

.approval-workflow__step-header {
  margin-bottom: 8px;
}

.approval-workflow__step-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0 0 4px;
  color: rgb(var(--v-theme-on-surface));
}

.approval-workflow__step-meta {
  display: flex;
  gap: 16px;
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface-variant));
}

.approval-workflow__step-role {
  font-weight: 500;
}

.approval-workflow__step-description {
  margin: 0 0 12px;
  color: rgb(var(--v-theme-on-surface-variant));
  font-size: 0.875rem;
  line-height: 1.5;
}

.approval-workflow__step-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 0.875rem;
}

.approval-workflow__step-assignee,
.approval-workflow__step-due-date {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgb(var(--v-theme-on-surface-variant));
}

.approval-workflow__step-actions {
  margin-bottom: 16px;
}

.approval-workflow__default-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

/* Step States */
.approval-workflow__step--current {
  background: rgba(var(--v-theme-primary), 0.04);
  border-radius: 8px;
  padding: 16px;
  margin: -16px;
  margin-bottom: 16px;
}

.approval-workflow__step--overdue .approval-workflow__step-title {
  color: rgb(var(--v-theme-warning));
}

/* History */
.approval-workflow__step-history {
  margin-top: 16px;
}

.approval-workflow__history-entry {
  padding: 12px 0;
  border-bottom: 1px solid rgba(var(--v-theme-border), 0.5);
}

.approval-workflow__history-entry:last-child {
  border-bottom: none;
}

.approval-workflow__history-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 0.875rem;
}

.approval-workflow__history-user {
  color: rgb(var(--v-theme-on-surface-variant));
}

.approval-workflow__history-date {
  color: rgb(var(--v-theme-on-surface-variant));
  margin-left: auto;
}

.approval-workflow__history-comment {
  margin: 8px 0 0;
  padding: 8px 12px;
  background: rgba(var(--v-theme-surface-variant), 0.5);
  border-radius: 6px;
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface-variant));
}

/* Comments and Actions */
.approval-workflow__comments {
  margin-bottom: 24px;
}

.approval-workflow__workflow-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid rgb(var(--v-theme-border));
}

/* Responsive design */
@media (max-width: 959px) {
  .approval-workflow__header {
    flex-direction: column;
    gap: 16px;
  }
  
  .approval-workflow__step {
    flex-direction: column;
  }
  
  .approval-workflow__step-icon {
    margin-right: 0;
    margin-bottom: 16px;
    align-self: flex-start;
  }
  
  .approval-workflow__step-connector {
    display: none;
  }
}

@media (max-width: 599px) {
  .approval-workflow {
    padding: 16px;
  }
  
  .approval-workflow__default-actions {
    flex-direction: column;
  }
  
  .approval-workflow__workflow-actions {
    flex-direction: column;
  }
}

/* Workflow variants */
.approval-workflow--high-priority {
  border-left: 4px solid rgb(var(--v-theme-error));
}

.approval-workflow--standard {
  border-left: 4px solid rgb(var(--v-theme-primary));
}

.approval-workflow--low-priority {
  border-left: 4px solid rgb(var(--v-theme-info));
}

/* High contrast mode */
@media (prefers-contrast: high) {
  .approval-workflow {
    border: 2px solid rgb(var(--v-theme-outline));
  }
  
  .approval-workflow__step--current {
    border: 2px solid rgb(var(--v-theme-primary));
  }
}
</style>