<template>
  <VBtn
    :variant="computedVariant"
    :size="size"
    :color="computedColor"
    :disabled="disabled || loading"
    :loading="loading"
    :block="block"
    :elevation="elevated ? 2 : 0"
    :class="computedClasses"
    :style="style"
    :aria-label="computedAriaLabel"
    :aria-describedby="ariaDescribedby"
    :aria-disabled="disabled"
    @click="handleClick"
    v-bind="$attrs"
  >
    <!-- Icon slot -->
    <template v-if="icon && iconPosition === 'start'" #prepend>
      <slot name="icon">
        <VIcon :icon="icon" />
      </slot>
    </template>

    <!-- Default content slot -->
    <slot>
      {{ defaultButtonText }}
    </slot>

    <!-- Loading slot -->
    <template v-if="loading" #loader>
      <slot name="loading">
        <VProgressCircular indeterminate size="16" width="2" :color="loadingColor" />
      </slot>
    </template>

    <!-- End icon slot -->
    <template v-if="icon && iconPosition === 'end'" #append>
      <slot name="icon">
        <VIcon :icon="icon" />
      </slot>
    </template>
  </VBtn>

  <!-- Confirmation Dialog -->
  <VDialog v-model="showConfirmation" max-width="400" persistent>
    <VCard>
      <VCardTitle class="text-h6">
        {{ confirmationOptions.title || 'Confirm Action' }}
      </VCardTitle>

      <VCardText>
        <div class="mb-4">
          {{ confirmationOptions.message }}
        </div>

        <VAlert v-if="confirmationOptions.destructive" type="warning" variant="tonal" class="mb-2">
          This action cannot be undone.
        </VAlert>
      </VCardText>

      <VCardActions>
        <VSpacer />

        <VBtn variant="text" @click="cancelConfirmation">
          {{ confirmationOptions.cancelText || 'Cancel' }}
        </VBtn>

        <VBtn
          :variant="confirmationOptions.destructive ? 'flat' : 'elevated'"
          :color="
            confirmationOptions.severity === 'error'
              ? 'error'
              : confirmationOptions.severity === 'warning'
                ? 'warning'
                : 'primary'
          "
          @click="acceptConfirmation"
        >
          {{ confirmationOptions.confirmText || 'Confirm' }}
        </VBtn>
      </VCardActions>
    </VCard>
  </VDialog>
</template>

<script setup lang="ts">
import { computed, ref, type PropType } from 'vue';
import {
  VBtn,
  VIcon,
  VProgressCircular,
  VDialog,
  VCard,
  VCardTitle,
  VCardText,
  VCardActions,
  VSpacer,
  VAlert,
} from 'vuetify/components';
import type {
  BankingButtonProps,
  BankingButtonClickEvent,
  ConfirmationDialogOptions,
  BankingButtonEmits,
} from './BankingButton.types';

/**
 * Banking Button Component
 *
 * A professional button component designed for banking applications with enhanced
 * accessibility, confirmation dialogs, and banking-specific styling.
 *
 * Features:
 * - WCAG 2.1 AA compliant
 * - Banking-specific variants and colors
 * - Confirmation dialogs for critical actions
 * - Loading states and disabled states with reasons
 * - Full keyboard navigation support
 * - Responsive design
 *
 * @example
 * <BankingButton
 *   variant="approve"
 *   action-type="approve"
 *   :require-confirmation="true"
 *   confirmation-message="Are you sure you want to approve this transaction?"
 *   @click="handleApproval"
 * >
 *   Approve Transaction
 * </BankingButton>
 */

// Props definition with comprehensive banking-specific options
const props = withDefaults(defineProps<BankingButtonProps>(), {
  variant: 'primary',
  size: 'default',
  iconPosition: 'start',
  disabled: false,
  loading: false,
  requireConfirmation: false,
  block: false,
  elevated: false,
});

// Events definition
const emit = defineEmits<BankingButtonEmits>();

// Reactive state for confirmation dialog
const showConfirmation = ref(false);
const pendingClickEvent = ref<Event | null>(null);

// Computed confirmation options
const confirmationOptions = computed<ConfirmationDialogOptions>(() => ({
  message: props.confirmationMessage || 'Are you sure you want to perform this action?',
  severity: props.actionType === 'reject' ? 'warning' : 'info',
  destructive: ['reject', 'cancel'].includes(props.actionType || ''),
  confirmText: getConfirmText(),
  cancelText: 'Cancel',
}));

/**
 * Get appropriate confirmation text based on action type
 */
function getConfirmText(): string {
  switch (props.actionType) {
    case 'approve':
      return 'Approve';
    case 'reject':
      return 'Reject';
    case 'submit':
      return 'Submit';
    case 'save-draft':
      return 'Save Draft';
    case 'send-reminder':
      return 'Send Reminder';
    case 'export':
      return 'Export';
    default:
      return 'Confirm';
  }
}

/**
 * Compute Vuetify variant based on banking variant
 */
const computedVariant = computed(() => {
  switch (props.variant) {
    case 'primary':
    case 'approve':
      return 'elevated';
    case 'secondary':
      return 'flat';
    case 'reject':
    case 'warning':
      return 'flat';
    case 'outlined':
      return 'outlined';
    case 'text':
      return 'text';
    default:
      return 'elevated';
  }
});

/**
 * Compute color based on banking variant and action type
 */
const computedColor = computed(() => {
  if (props.color) return props.color;

  switch (props.variant) {
    case 'primary':
      return 'primary';
    case 'secondary':
      return 'secondary';
    case 'approve':
      return 'success';
    case 'reject':
      return 'error';
    case 'warning':
      return 'warning';
    default:
      return 'primary';
  }
});

/**
 * Compute CSS classes for banking-specific styling
 */
const computedClasses = computed(() => {
  const classes: string[] = ['banking-button'];

  // Add variant-specific classes
  if (props.variant) {
    classes.push(`banking-button--${props.variant}`);
  }

  // Add action-type specific classes
  if (props.actionType) {
    classes.push(`banking-button--${props.actionType}`);
  }

  // Add disabled reason class for styling
  if (props.disabled && props.disabledReason) {
    classes.push('banking-button--disabled-with-reason');
  }

  // Add user-provided classes
  if (props.class) {
    if (typeof props.class === 'string') {
      classes.push(props.class);
    } else if (Array.isArray(props.class)) {
      classes.push(...props.class);
    } else {
      Object.entries(props.class).forEach(([className, condition]) => {
        if (condition) classes.push(className);
      });
    }
  }

  return classes;
});

/**
 * Compute accessibility label
 */
const computedAriaLabel = computed(() => {
  if (props.ariaLabel) return props.ariaLabel;

  // If disabled with reason, include that in the label
  if (props.disabled && props.disabledReason) {
    return `Button disabled: ${props.disabledReason}`;
  }

  // Generate descriptive label based on action type
  if (props.actionType) {
    return `${getConfirmText()} action button`;
  }

  return undefined;
});

/**
 * Default button text based on action type
 */
const defaultButtonText = computed(() => {
  switch (props.actionType) {
    case 'approve':
      return 'Approve';
    case 'reject':
      return 'Reject';
    case 'submit':
      return 'Submit';
    case 'save-draft':
      return 'Save Draft';
    case 'send-reminder':
      return 'Send Reminder';
    case 'export':
      return 'Export';
    case 'cancel':
      return 'Cancel';
    default:
      return 'Button';
  }
});

/**
 * Loading indicator color based on variant
 */
const loadingColor = computed(() => {
  return props.variant === 'text' || props.variant === 'outlined'
    ? computedColor.value
    : 'on-primary';
});

/**
 * Handle button click with confirmation support
 */
function handleClick(event: Event): void {
  if (props.requireConfirmation) {
    pendingClickEvent.value = event;
    showConfirmation.value = true;
    emit('confirmation-shown');
  } else {
    executeClick(event);
  }
}

/**
 * Execute the actual click action
 */
function executeClick(event: Event, confirmed = false): void {
  const clickEvent: BankingButtonClickEvent = {
    originalEvent: event,
    actionType: props.actionType,
    timestamp: new Date(),
    confirmed,
  };

  emit('click', clickEvent);
}

/**
 * Accept confirmation and execute click
 */
function acceptConfirmation(): void {
  showConfirmation.value = false;
  emit('confirmation-accepted');

  if (pendingClickEvent.value) {
    executeClick(pendingClickEvent.value, true);
    pendingClickEvent.value = null;
  }
}

/**
 * Cancel confirmation dialog
 */
function cancelConfirmation(): void {
  showConfirmation.value = false;
  emit('confirmation-rejected');
  pendingClickEvent.value = null;
}
</script>

<style scoped>
/**
 * Banking Button Styles
 * 
 * Professional styling for banking applications with enhanced
 * accessibility and responsive design.
 */

.banking-button {
  /* Enhanced focus indicators for accessibility */
  position: relative;
  font-weight: 500;
  letter-spacing: 0.0892857143em;
  text-transform: none; /* Override Vuetify default uppercase */
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.banking-button:focus-visible {
  outline: 2px solid rgb(var(--v-theme-primary));
  outline-offset: 2px;
}

/* Banking-specific variant styles */
.banking-button--approve {
  /* Success/approve specific styling */
  font-weight: 600;
}

.banking-button--reject {
  /* Error/reject specific styling */
  font-weight: 600;
}

.banking-button--warning {
  /* Warning specific styling */
  font-weight: 500;
}

.banking-button--save-draft {
  /* Draft action styling */
  font-style: italic;
}

/* Disabled state with reason tooltip support */
.banking-button--disabled-with-reason {
  cursor: help;
}

/* Responsive adjustments */
@media (max-width: 599px) {
  .banking-button {
    /* Increase touch target size on mobile */
    min-height: 44px;
    padding: 0 16px;
  }

  .banking-button.v-btn--size-small {
    min-height: 40px;
    padding: 0 12px;
  }
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  .banking-button {
    border: 2px solid;
    font-weight: 700;
  }
}

/* Reduced motion support */
@media (prefers-reduced-motion: reduce) {
  .banking-button {
    transition: none;
  }
}
</style>
