/**
 * GoEdu Design System - Component Library Index
 *
 * Central export file for all design system components, types, and utilities.
 * This provides a clean import interface for consuming applications.
 *
 * @example
 * import { BankingButton, BankingInput, ApprovalWorkflow } from '@/components/design-system'
 */

// Base Components
export { default as BankingButton } from './base/BankingButton.vue';
export { default as BankingInput } from './base/BankingInput.vue';
export { default as BankingTable } from './base/BankingTable.vue';
export { default as BankingTableCell } from './base/BankingTableCell.vue';

// Banking-Specific Components
export { default as ApprovalWorkflow } from './banking/ApprovalWorkflow.vue';
export { default as AuditTrail } from './banking/AuditTrail.vue';

// Showcase Component
export { default as DesignSystemShowcase } from './showcase/DesignSystemShowcase.vue';

// Type Exports - Base Components
export type {
  BankingButtonProps,
  BankingButtonEmits,
  BankingButtonClickEvent,
  ConfirmationDialogOptions,
  BankingButtonSlots,
} from './base/BankingButton.types';

export type {
  BankingInputProps,
  BankingInputEmits,
  ValidationResult,
  BankingValidators,
  BankingInputSlots,
  FormattingOptions,
} from './base/BankingInput.types';

export type {
  BankingTableProps,
  BankingTableEmits,
  BankingTableHeader,
  BankingTableItem,
  BankingColumnFormat,
  TableOptions,
  ExportOptions,
  BankingTableCellProps,
} from './base/BankingTable.types';

// Type Exports - Banking Components
export type {
  ApprovalWorkflowProps,
  ApprovalWorkflowEmits,
  ApprovalStep,
  WorkflowStatus,
  StepStatus,
  StepType,
  StepCondition,
  StepRequirement,
  StepHistoryEntry,
  StepActionEvent,
  WorkflowActionEvent,
  ExportApprovalEvent,
  WorkflowConfig,
  WorkflowTemplate,
  WorkflowAnalytics,
} from './banking/ApprovalWorkflow.types';

export type {
  AuditTrailProps,
  AuditTrailEmits,
  AuditEntry,
  AuditChange,
  AuditFilters,
  AuditExportRequest,
  AuditTrailConfig,
  AuditTrailAnalytics,
} from './banking/AuditTrail.types';

/**
 * Design System Configuration
 */
export interface DesignSystemConfig {
  /**
   * Default theme (light/dark)
   */
  theme: 'light' | 'dark';

  /**
   * Default currency for currency inputs
   */
  defaultCurrency: string;

  /**
   * Default locale for formatting
   */
  defaultLocale: string;

  /**
   * Accessibility features
   */
  accessibility: {
    /**
     * Enable high contrast mode
     */
    highContrast: boolean;

    /**
     * Enable reduced motion
     */
    reducedMotion: boolean;

    /**
     * Enable screen reader optimizations
     */
    screenReader: boolean;
  };

  /**
   * Banking-specific settings
   */
  banking: {
    /**
     * Default decimal places for currency
     */
    currencyDecimals: number;

    /**
     * Default decimal places for percentages
     */
    percentageDecimals: number;

    /**
     * Enable audit trail by default
     */
    enableAuditTrail: boolean;

    /**
     * Default approval workflow timeout (hours)
     */
    approvalTimeout: number;
  };
}

/**
 * Default design system configuration
 */
export const defaultConfig: DesignSystemConfig = {
  theme: 'light',
  defaultCurrency: 'USD',
  defaultLocale: 'en-US',
  accessibility: {
    highContrast: false,
    reducedMotion: false,
    screenReader: true,
  },
  banking: {
    currencyDecimals: 2,
    percentageDecimals: 2,
    enableAuditTrail: true,
    approvalTimeout: 48,
  },
};

/**
 * Component categories for organization
 */
export const componentCategories = {
  base: ['BankingButton', 'BankingInput', 'BankingTable', 'BankingTableCell'],
  banking: ['ApprovalWorkflow', 'AuditTrail'],
  layout: [
    // Future layout components
  ],
  utility: ['DesignSystemShowcase'],
} as const;

/**
 * Accessibility compliance information
 */
export const accessibilityCompliance = {
  standard: 'WCAG 2.1 AA',
  features: [
    'Keyboard navigation support',
    'Screen reader compatibility',
    'High contrast mode support',
    'Focus management',
    'Semantic HTML structure',
    'ARIA labels and roles',
    'Color contrast compliance',
    'Text scaling support (up to 200%)',
    'Reduced motion support',
  ],
  testing: [
    'Automated accessibility testing',
    'Manual keyboard testing',
    'Screen reader testing (NVDA, JAWS, VoiceOver)',
    'Color contrast validation',
    'Mobile accessibility testing',
  ],
} as const;

/**
 * Banking industry compliance information
 */
export const bankingCompliance = {
  standards: [
    'SOX Section 404 (Internal Controls)',
    'FFIEC Guidelines (IT Risk Management)',
    'Basel III (Operational Risk)',
    'GDPR Article 32 (Security of Processing)',
    'PCI DSS (Payment Card Industry)',
  ],
  features: [
    'Audit trail logging',
    'Multi-step approval workflows',
    'Role-based access control',
    'Data encryption support',
    'Session management',
    'Secure input handling',
    'Change tracking',
    'Compliance reporting',
  ],
} as const;

/**
 * Theme configuration
 */
export interface ThemeConfig {
  colors: {
    primary: string;
    secondary: string;
    success: string;
    warning: string;
    error: string;
    info: string;
    surface: string;
    background: string;
  };
  typography: {
    fontFamily: string;
    fontSize: {
      xs: string;
      sm: string;
      base: string;
      lg: string;
      xl: string;
    };
  };
  spacing: {
    xs: string;
    sm: string;
    md: string;
    lg: string;
    xl: string;
  };
  borderRadius: {
    sm: string;
    md: string;
    lg: string;
  };
}

/**
 * Component usage statistics for optimization
 */
export interface UsageStats {
  component: string;
  usageCount: number;
  lastUsed: Date;
  averageRenderTime: number;
  bundleSize: number;
}

/**
 * Performance metrics
 */
export interface PerformanceMetrics {
  totalComponents: number;
  bundleSize: string;
  treeShakeable: boolean;
  lazyLoadable: boolean;
  averageRenderTime: number;
  memoryUsage: string;
}

/**
 * Design system metadata
 */
export const designSystemMetadata = {
  name: 'GoEdu Banking Design System',
  version: '1.0.0',
  description: 'Professional banking UI components with WCAG 2.1 AA compliance',
  author: 'GoEdu Team',
  license: 'Proprietary',
  dependencies: {
    vue: '^3.5.17',
    vuetify: '^3.7.9',
    typescript: '~5.8.0',
  },
  peerDependencies: {
    '@mdi/font': '^7.4.47',
  },
  features: [
    'Banking-specific components',
    'WCAG 2.1 AA compliance',
    'TypeScript support',
    'Responsive design',
    'Dark mode support',
    'Comprehensive documentation',
    'Interactive component showcase',
    'Automated testing',
    'Performance optimized',
  ],
} as const;
