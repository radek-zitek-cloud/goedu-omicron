/**
 * Banking Button Component Types
 *
 * Comprehensive type definitions for the banking-specific button component.
 * Supports various banking workflows and accessibility requirements.
 */

export interface BankingButtonProps {
  /**
   * Button variant determining the visual style and semantic meaning
   * @default 'primary'
   */
  variant?: 'primary' | 'secondary' | 'approve' | 'reject' | 'warning' | 'text' | 'outlined';

  /**
   * Button size for different contexts and device types
   * @default 'default'
   */
  size?: 'small' | 'default' | 'large' | 'x-large';

  /**
   * Banking-specific action types for enhanced semantics
   */
  actionType?:
    | 'submit'
    | 'approve'
    | 'reject'
    | 'cancel'
    | 'save-draft'
    | 'send-reminder'
    | 'export';

  /**
   * Disable state with optional reason for accessibility
   */
  disabled?: boolean;

  /**
   * Reason for disabled state (used for accessibility)
   */
  disabledReason?: string;

  /**
   * Loading state for async operations
   */
  loading?: boolean;

  /**
   * Show confirmation dialog before action
   */
  requireConfirmation?: boolean;

  /**
   * Confirmation dialog message
   */
  confirmationMessage?: string;

  /**
   * Icon to display (Material Design Icons name)
   */
  icon?: string;

  /**
   * Position of the icon relative to text
   * @default 'start'
   */
  iconPosition?: 'start' | 'end';

  /**
   * Full width button (useful for mobile)
   */
  block?: boolean;

  /**
   * Elevated appearance with shadow
   */
  elevated?: boolean;

  /**
   * Custom color (overrides variant colors)
   */
  color?: string;

  /**
   * Accessibility label when button text is not descriptive enough
   */
  ariaLabel?: string;

  /**
   * ARIA describedby for additional context
   */
  ariaDescribedby?: string;

  /**
   * Custom CSS classes
   */
  class?: string | string[] | Record<string, boolean>;

  /**
   * Custom styles
   */
  style?: string | Record<string, string | number>;
}

/**
 * Button click event with enhanced context for banking operations
 */
export interface BankingButtonClickEvent {
  /**
   * Original mouse/keyboard event
   */
  originalEvent: Event;

  /**
   * Button action type
   */
  actionType?: string;

  /**
   * Timestamp of the action
   */
  timestamp: Date;

  /**
   * Whether confirmation was required and accepted
   */
  confirmed?: boolean;
}

/**
 * Confirmation dialog options
 */
export interface ConfirmationDialogOptions {
  /**
   * Dialog title
   */
  title?: string;

  /**
   * Confirmation message
   */
  message: string;

  /**
   * Confirm button text
   * @default 'Confirm'
   */
  confirmText?: string;

  /**
   * Cancel button text
   * @default 'Cancel'
   */
  cancelText?: string;

  /**
   * Dialog severity level
   * @default 'info'
   */
  severity?: 'info' | 'warning' | 'error';

  /**
   * Whether to show additional warning for destructive actions
   */
  destructive?: boolean;
}

/**
 * Component emits interface
 */
export interface BankingButtonEmits {
  /**
   * Emitted when button is clicked (after confirmation if required)
   */
  (event: 'click', payload: BankingButtonClickEvent): void;

  /**
   * Emitted when confirmation dialog is shown
   */
  (event: 'confirmation-shown'): void;

  /**
   * Emitted when confirmation is accepted
   */
  (event: 'confirmation-accepted'): void;

  /**
   * Emitted when confirmation is rejected
   */
  (event: 'confirmation-rejected'): void;
}

/**
 * Component slots interface
 */
export interface BankingButtonSlots {
  /**
   * Default slot for button content
   */
  default: () => any;

  /**
   * Icon slot for custom icons
   */
  icon?: () => any;

  /**
   * Loading slot for custom loading indicators
   */
  loading?: () => any;
}
