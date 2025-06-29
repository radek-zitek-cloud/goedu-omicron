/**
 * Approval Workflow Component Types
 *
 * Comprehensive type definitions for banking approval workflow components
 * with multi-step processes, role-based access, and audit trail functionality.
 */

/**
 * Workflow status enumeration
 */
export type WorkflowStatus =
  | 'draft'
  | 'in-progress'
  | 'completed'
  | 'rejected'
  | 'cancelled'
  | 'paused';

/**
 * Step status enumeration
 */
export type StepStatus = 'pending' | 'in-progress' | 'completed' | 'rejected' | 'skipped';

/**
 * Step type enumeration
 */
export type StepType = 'review' | 'approve' | 'sign' | 'validate' | 'notify' | 'custom';

/**
 * Approval step configuration
 */
export interface ApprovalStep {
  /**
   * Unique identifier for the step
   */
  id: string;

  /**
   * Step title
   */
  title: string;

  /**
   * Step description
   */
  description?: string;

  /**
   * Step type
   */
  type: StepType;

  /**
   * Current status of the step
   */
  status: StepStatus;

  /**
   * Role required to complete this step
   */
  role: string;

  /**
   * Specific user assigned to this step
   */
  assignee?: string;

  /**
   * Alternative assignees who can complete this step
   */
  alternativeAssignees?: string[];

  /**
   * Due date for completion
   */
  dueDate?: Date | string;

  /**
   * Estimated time to complete (human readable)
   */
  estimatedTime?: string;

  /**
   * Whether this step can be approved
   */
  allowApprove?: boolean;

  /**
   * Whether this step can be rejected
   */
  allowReject?: boolean;

  /**
   * Whether this step can be delegated
   */
  allowDelegate?: boolean;

  /**
   * Whether this step can be skipped
   */
  allowSkip?: boolean;

  /**
   * Conditions that must be met to complete this step
   */
  conditions?: StepCondition[];

  /**
   * Required fields or documents for this step
   */
  requirements?: StepRequirement[];

  /**
   * Step history and audit trail
   */
  history?: StepHistoryEntry[];

  /**
   * Additional metadata for the step
   */
  metadata?: Record<string, any>;
}

/**
 * Step condition interface
 */
export interface StepCondition {
  /**
   * Condition type
   */
  type: 'field_required' | 'document_uploaded' | 'approval_count' | 'custom';

  /**
   * Field or condition identifier
   */
  field?: string;

  /**
   * Required value or threshold
   */
  value?: any;

  /**
   * Human-readable description
   */
  description: string;

  /**
   * Whether condition is currently met
   */
  satisfied: boolean;
}

/**
 * Step requirement interface
 */
export interface StepRequirement {
  /**
   * Requirement type
   */
  type: 'document' | 'signature' | 'approval' | 'field' | 'custom';

  /**
   * Requirement identifier
   */
  id: string;

  /**
   * Human-readable name
   */
  name: string;

  /**
   * Whether requirement is mandatory
   */
  required: boolean;

  /**
   * Whether requirement is satisfied
   */
  satisfied: boolean;

  /**
   * Additional metadata
   */
  metadata?: Record<string, any>;
}

/**
 * Step history entry
 */
export interface StepHistoryEntry {
  /**
   * Action performed
   */
  action: string;

  /**
   * User who performed the action
   */
  user: string;

  /**
   * Timestamp of the action
   */
  timestamp: Date | string;

  /**
   * Optional comment or reason
   */
  comment?: string;

  /**
   * Additional metadata
   */
  metadata?: Record<string, any>;
}

/**
 * Workflow action event payloads
 */
export interface StepActionEvent {
  /**
   * The step that was acted upon
   */
  step: ApprovalStep;

  /**
   * Index of the step in the workflow
   */
  stepIndex: number;

  /**
   * Optional comment from the user
   */
  comment?: string;

  /**
   * Timestamp of the action
   */
  timestamp: Date;

  /**
   * Additional context
   */
  metadata?: Record<string, any>;
}

export interface WorkflowActionEvent {
  /**
   * Timestamp of the action
   */
  timestamp: Date;

  /**
   * Optional reason or comment
   */
  reason?: string;

  /**
   * Additional context
   */
  metadata?: Record<string, any>;
}

export interface ExportApprovalEvent {
  /**
   * All workflow steps
   */
  steps: ApprovalStep[];

  /**
   * Current workflow status
   */
  workflowStatus: WorkflowStatus;

  /**
   * Export timestamp
   */
  timestamp: Date;

  /**
   * Export format preference
   */
  format?: 'pdf' | 'excel' | 'json';
}

/**
 * Main component props
 */
export interface ApprovalWorkflowProps {
  /**
   * Array of workflow steps
   */
  steps: ApprovalStep[];

  /**
   * Current active step index
   */
  currentStep?: number;

  /**
   * Overall workflow status
   */
  workflowStatus?: WorkflowStatus;

  /**
   * Workflow title
   */
  title?: string;

  /**
   * Workflow description
   */
  description?: string;

  /**
   * Workflow variant for styling
   */
  variant?: 'standard' | 'high-priority' | 'low-priority' | 'urgent';

  /**
   * Whether to show the workflow header
   */
  showHeader?: boolean;

  /**
   * Whether to allow comments on actions
   */
  allowComments?: boolean;

  /**
   * Whether to show workflow-level actions
   */
  showWorkflowActions?: boolean;

  /**
   * Whether steps must be completed in strict sequence
   */
  strictSequence?: boolean;

  /**
   * Whether to automatically progress to next step after completion
   */
  autoProgress?: boolean;

  /**
   * Whether workflow can be cancelled
   */
  allowCancel?: boolean;

  /**
   * Current user information for permission checks
   */
  currentUser?: {
    id: string;
    roles: string[];
    permissions: string[];
  };

  /**
   * Workflow configuration options
   */
  config?: WorkflowConfig;

  /**
   * Custom CSS classes
   */
  class?: string | string[] | Record<string, boolean>;
}

/**
 * Workflow configuration options
 */
export interface WorkflowConfig {
  /**
   * Email notifications configuration
   */
  notifications?: {
    enabled: boolean;
    onStepAssignment: boolean;
    onStepCompletion: boolean;
    onWorkflowCompletion: boolean;
    onOverdue: boolean;
  };

  /**
   * Escalation rules
   */
  escalation?: {
    enabled: boolean;
    escalateAfterHours: number;
    escalateTo: string[];
  };

  /**
   * Approval requirements
   */
  approvalRules?: {
    requireAllSteps: boolean;
    allowParallelApprovals: boolean;
    minimumApprovals: number;
  };

  /**
   * Document requirements
   */
  documentRequirements?: {
    requireSignatures: boolean;
    allowDigitalSignatures: boolean;
    retentionPeriod: number;
  };
}

/**
 * Component events
 */
export interface ApprovalWorkflowEmits {
  /**
   * Emitted when a step is approved
   */
  (event: 'step-approved', payload: StepActionEvent): void;

  /**
   * Emitted when a step is rejected
   */
  (event: 'step-rejected', payload: StepActionEvent): void;

  /**
   * Emitted when a step is delegated
   */
  (event: 'step-delegated', payload: StepActionEvent): void;

  /**
   * Emitted when a step is skipped
   */
  (event: 'step-skipped', payload: StepActionEvent): void;

  /**
   * Emitted when workflow is started
   */
  (event: 'workflow-started', payload: WorkflowActionEvent): void;

  /**
   * Emitted when workflow is completed
   */
  (event: 'workflow-completed', payload: WorkflowActionEvent): void;

  /**
   * Emitted when workflow is cancelled
   */
  (event: 'workflow-cancelled', payload: WorkflowActionEvent): void;

  /**
   * Emitted when workflow is paused
   */
  (event: 'workflow-paused', payload: WorkflowActionEvent): void;

  /**
   * Emitted when workflow is resumed
   */
  (event: 'workflow-resumed', payload: WorkflowActionEvent): void;

  /**
   * Emitted when approval record export is requested
   */
  (event: 'export-approval', payload: ExportApprovalEvent): void;

  /**
   * Emitted when workflow configuration changes
   */
  (event: 'config-updated', config: WorkflowConfig): void;

  /**
   * Emitted when step assignment changes
   */
  (event: 'step-reassigned', payload: StepActionEvent & { newAssignee: string }): void;
}

/**
 * Workflow template interface for common banking processes
 */
export interface WorkflowTemplate {
  /**
   * Template identifier
   */
  id: string;

  /**
   * Template name
   */
  name: string;

  /**
   * Template description
   */
  description: string;

  /**
   * Banking process category
   */
  category: 'transaction' | 'loan' | 'account' | 'compliance' | 'audit' | 'risk';

  /**
   * Template steps
   */
  steps: Omit<ApprovalStep, 'id' | 'status' | 'history'>[];

  /**
   * Default configuration
   */
  defaultConfig: WorkflowConfig;

  /**
   * Required roles for this workflow type
   */
  requiredRoles: string[];

  /**
   * Estimated completion time
   */
  estimatedDuration: string;

  /**
   * Compliance requirements
   */
  complianceRequirements?: string[];
}

/**
 * Workflow analytics interface
 */
export interface WorkflowAnalytics {
  /**
   * Total workflows processed
   */
  totalWorkflows: number;

  /**
   * Average completion time
   */
  averageCompletionTime: number;

  /**
   * Success rate percentage
   */
  successRate: number;

  /**
   * Most common rejection reasons
   */
  rejectionReasons: Array<{
    reason: string;
    count: number;
    percentage: number;
  }>;

  /**
   * Step performance metrics
   */
  stepMetrics: Array<{
    stepType: StepType;
    averageTime: number;
    completionRate: number;
  }>;

  /**
   * User performance metrics
   */
  userMetrics: Array<{
    userId: string;
    completedWorkflows: number;
    averageResponseTime: number;
    approvalRate: number;
  }>;
}
