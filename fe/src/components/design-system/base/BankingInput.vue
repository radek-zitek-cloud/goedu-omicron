<template>
  <VTextField
    :model-value="formattedValue"
    :variant="variant"
    :density="density"
    :label="computedLabel"
    :placeholder="placeholder"
    :hint="hint"
    :disabled="disabled"
    :readonly="readonly"
    :loading="loading"
    :rules="computedRules"
    :error="error"
    :error-messages="errorMessages"
    :clearable="clearable"
    :prepend-icon="prependIcon"
    :append-icon="appendIcon"
    :maxlength="maxLength"
    :type="computedInputType"
    :inputmode="computedInputMode"
    :autocomplete="autocomplete"
    :class="computedClasses"
    :aria-label="ariaLabel"
    :aria-describedby="ariaDescribedby"
    :validate-on="validateOnBlur ? 'blur' : 'input'"
    @update:model-value="handleInput"
    @blur="handleBlur"
    @focus="handleFocus"
    @keydown.enter="handleEnter"
    @click:clear="handleClear"
    v-bind="$attrs"
  >
    <!-- Pass through all slots -->
    <template v-for="(_, slot) of $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="slotProps" />
    </template>

    <!-- Currency symbol for currency inputs -->
    <template v-if="type === 'currency'" #prepend-inner>
      <span class="banking-input__currency-symbol">
        {{ currencySymbol }}
      </span>
    </template>

    <!-- Percentage symbol for percentage inputs -->
    <template v-if="type === 'percentage'" #append-inner>
      <span class="banking-input__percentage-symbol">%</span>
    </template>

    <!-- Security indicator for sensitive fields -->
    <template v-if="isSensitiveField" #append-inner>
      <VIcon
        icon="mdi-shield-check"
        size="small"
        class="banking-input__security-icon"
        :title="'Secure field: ' + type"
      />
    </template>
  </VTextField>
</template>

<script setup lang="ts">
import { computed, ref, watch, nextTick, type PropType } from 'vue';
import { VTextField, VIcon } from 'vuetify/components';
import type {
  BankingInputProps,
  BankingInputEmits,
  ValidationResult,
  BankingValidators,
} from './BankingInput.types';

/**
 * Banking Input Component
 *
 * A comprehensive input component designed for banking applications with
 * specialized formatting, validation, and accessibility features.
 *
 * Features:
 * - Banking-specific input types (currency, account numbers, etc.)
 * - Automatic formatting and validation
 * - WCAG 2.1 AA compliant
 * - Secure handling of sensitive data
 * - Mobile-optimized input modes
 *
 * @example
 * <BankingInput
 *   v-model="amount"
 *   type="currency"
 *   label="Transaction Amount"
 *   currency="USD"
 *   :rules="[validators.required, validators.currency]"
 *   required
 * />
 */

// Props with banking-specific defaults
const props = withDefaults(defineProps<BankingInputProps>(), {
  type: 'text',
  variant: 'outlined',
  density: 'default',
  autoFormat: true,
  validateOnBlur: false,
  immediateValidation: false,
  decimalPlaces: 2,
  currency: 'USD',
  clearable: false,
});

// Events
const emit = defineEmits<BankingInputEmits>();

// Internal state
const internalValue = ref<string | number>(props.modelValue || '');
const validationErrors = ref<string[]>([]);

/**
 * Banking validators for different field types
 */
const validators: BankingValidators = {
  required: (value: any): ValidationResult => ({
    valid: value !== null && value !== undefined && value !== '',
    message: 'This field is required',
  }),

  email: (value: string): ValidationResult => {
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return {
      valid: !value || emailPattern.test(value),
      message: 'Please enter a valid email address',
    };
  },

  currency: (value: string | number, min?: number, max?: number): ValidationResult => {
    const numValue = typeof value === 'string' ? parseFloat(value.replace(/[^0-9.-]/g, '')) : value;
    if (isNaN(numValue)) {
      return { valid: false, message: 'Please enter a valid amount' };
    }
    if (min !== undefined && numValue < min) {
      return { valid: false, message: `Amount must be at least ${formatCurrency(min)}` };
    }
    if (max !== undefined && numValue > max) {
      return { valid: false, message: `Amount cannot exceed ${formatCurrency(max)}` };
    }
    return { valid: true };
  },

  percentage: (value: string | number): ValidationResult => {
    const numValue = typeof value === 'string' ? parseFloat(value.replace(/[^0-9.-]/g, '')) : value;
    if (isNaN(numValue)) {
      return { valid: false, message: 'Please enter a valid percentage' };
    }
    if (numValue < 0 || numValue > 100) {
      return { valid: false, message: 'Percentage must be between 0 and 100' };
    }
    return { valid: true };
  },

  accountNumber: (value: string): ValidationResult => {
    const cleaned = value.replace(/\D/g, '');
    return {
      valid: cleaned.length >= 8 && cleaned.length <= 17,
      message: 'Account number must be 8-17 digits',
    };
  },

  routingNumber: (value: string): ValidationResult => {
    const cleaned = value.replace(/\D/g, '');
    if (cleaned.length !== 9) {
      return { valid: false, message: 'Routing number must be 9 digits' };
    }

    // ABA routing number checksum validation
    const digits = cleaned.split('').map(Number);
    const checksum =
      (3 * (digits[0] + digits[3] + digits[6]) +
        7 * (digits[1] + digits[4] + digits[7]) +
        (digits[2] + digits[5] + digits[8])) %
      10;

    return {
      valid: checksum === 0,
      message: 'Please enter a valid routing number',
    };
  },

  ssn: (value: string): ValidationResult => {
    const cleaned = value.replace(/\D/g, '');
    return {
      valid: cleaned.length === 9,
      message: 'SSN must be 9 digits',
    };
  },

  phone: (value: string): ValidationResult => {
    const cleaned = value.replace(/\D/g, '');
    return {
      valid: cleaned.length === 10 || cleaned.length === 11,
      message: 'Please enter a valid phone number',
    };
  },

  minLength:
    (min: number) =>
    (value: string): ValidationResult => ({
      valid: !value || value.length >= min,
      message: `Must be at least ${min} characters`,
    }),

  maxLength:
    (max: number) =>
    (value: string): ValidationResult => ({
      valid: !value || value.length <= max,
      message: `Must be no more than ${max} characters`,
    }),

  pattern:
    (pattern: RegExp, message: string) =>
    (value: string): ValidationResult => ({
      valid: !value || pattern.test(value),
      message,
    }),
};

/**
 * Computed properties
 */
const computedLabel = computed(() => {
  return props.required && props.label ? `${props.label} *` : props.label;
});

const computedInputType = computed(() => {
  switch (props.type) {
    case 'currency':
    case 'percentage':
      return 'text'; // We handle formatting manually
    case 'account-number':
    case 'routing-number':
    case 'ssn':
      return 'text'; // For security and formatting
    case 'phone':
      return 'tel';
    default:
      return props.type;
  }
});

const computedInputMode = computed(() => {
  if (props.inputmode) return props.inputmode;

  switch (props.type) {
    case 'currency':
    case 'percentage':
      return 'decimal';
    case 'account-number':
    case 'routing-number':
    case 'ssn':
    case 'phone':
      return 'numeric';
    case 'email':
      return 'email';
    default:
      return 'text';
  }
});

const computedClasses = computed(() => {
  const classes: string[] = ['banking-input'];

  // Add type-specific classes
  classes.push(`banking-input--${props.type}`);

  // Add sensitive field class
  if (isSensitiveField.value) {
    classes.push('banking-input--sensitive');
  }

  // Add user classes
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

const isSensitiveField = computed(() => {
  return ['password', 'ssn', 'account-number', 'routing-number'].includes(props.type || '');
});

const currencySymbol = computed(() => {
  try {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: props.currency,
    })
      .format(0)
      .replace(/\d/g, '')
      .replace(/[.,\s]/g, '');
  } catch {
    return '$';
  }
});

const computedRules = computed(() => {
  const rules: ((value: any) => boolean | string)[] = [];

  // Add required validation
  if (props.required) {
    rules.push((value: any) => {
      const result = validators.required(value);
      return result.valid || result.message!;
    });
  }

  // Add type-specific validation
  switch (props.type) {
    case 'email':
      rules.push((value: string) => {
        const result = validators.email(value);
        return result.valid || result.message!;
      });
      break;
    case 'currency':
      rules.push((value: string) => {
        const result = validators.currency(value, props.min, props.max);
        return result.valid || result.message!;
      });
      break;
    case 'percentage':
      rules.push((value: string) => {
        const result = validators.percentage(value);
        return result.valid || result.message!;
      });
      break;
    case 'account-number':
      rules.push((value: string) => {
        const result = validators.accountNumber(value);
        return result.valid || result.message!;
      });
      break;
    case 'routing-number':
      rules.push((value: string) => {
        const result = validators.routingNumber(value);
        return result.valid || result.message!;
      });
      break;
    case 'ssn':
      rules.push((value: string) => {
        const result = validators.ssn(value);
        return result.valid || result.message!;
      });
      break;
    case 'phone':
      rules.push((value: string) => {
        const result = validators.phone(value);
        return result.valid || result.message!;
      });
      break;
  }

  // Add custom rules
  if (props.rules) {
    rules.push(...props.rules);
  }

  return rules;
});

const formattedValue = computed({
  get: () => {
    if (!props.autoFormat) return internalValue.value;

    const value = String(internalValue.value || '');

    switch (props.type) {
      case 'currency':
        return formatCurrency(value);
      case 'percentage':
        return formatPercentage(value);
      case 'phone':
        return formatPhone(value);
      case 'ssn':
        return formatSSN(value);
      case 'account-number':
        return formatAccountNumber(value);
      case 'routing-number':
        return formatRoutingNumber(value);
      default:
        return value;
    }
  },
  set: (value: string | number) => {
    internalValue.value = value;
    emit('update:modelValue', value);
  },
});

/**
 * Formatting functions
 */
function formatCurrency(value: string | number): string {
  if (!value && value !== 0) return '';

  const cleanValue = String(value).replace(/[^0-9.-]/g, '');
  const numValue = parseFloat(cleanValue);

  if (isNaN(numValue)) return cleanValue;

  try {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: props.currency,
      minimumFractionDigits: props.decimalPlaces,
      maximumFractionDigits: props.decimalPlaces,
    }).format(numValue);
  } catch {
    return `$${numValue.toFixed(props.decimalPlaces)}`;
  }
}

function formatPercentage(value: string | number): string {
  if (!value && value !== 0) return '';

  const cleanValue = String(value).replace(/[^0-9.-]/g, '');
  const numValue = parseFloat(cleanValue);

  if (isNaN(numValue)) return cleanValue;

  return numValue.toFixed(props.decimalPlaces);
}

function formatPhone(value: string): string {
  const cleaned = value.replace(/\D/g, '');

  if (cleaned.length === 0) return '';
  if (cleaned.length <= 3) return cleaned;
  if (cleaned.length <= 6) return `(${cleaned.slice(0, 3)}) ${cleaned.slice(3)}`;

  return `(${cleaned.slice(0, 3)}) ${cleaned.slice(3, 6)}-${cleaned.slice(6, 10)}`;
}

function formatSSN(value: string): string {
  const cleaned = value.replace(/\D/g, '');

  if (cleaned.length === 0) return '';
  if (cleaned.length <= 3) return cleaned;
  if (cleaned.length <= 5) return `${cleaned.slice(0, 3)}-${cleaned.slice(3)}`;

  return `${cleaned.slice(0, 3)}-${cleaned.slice(3, 5)}-${cleaned.slice(5, 9)}`;
}

function formatAccountNumber(value: string): string {
  const cleaned = value.replace(/\D/g, '');

  // For display, show only last 4 digits for security
  if (cleaned.length > 4) {
    return '*'.repeat(cleaned.length - 4) + cleaned.slice(-4);
  }

  return cleaned;
}

function formatRoutingNumber(value: string): string {
  const cleaned = value.replace(/\D/g, '');

  if (cleaned.length === 0) return '';
  if (cleaned.length <= 3) return cleaned;
  if (cleaned.length <= 6) return `${cleaned.slice(0, 3)}-${cleaned.slice(3)}`;

  return `${cleaned.slice(0, 3)}-${cleaned.slice(3, 6)}-${cleaned.slice(6, 9)}`;
}

/**
 * Event handlers
 */
function handleInput(value: string | number): void {
  internalValue.value = value;
  emit('update:modelValue', value);
  emit('input', value);
}

function handleBlur(event: FocusEvent): void {
  emit('blur', event);
}

function handleFocus(event: FocusEvent): void {
  emit('focus', event);
}

function handleEnter(event: KeyboardEvent): void {
  emit('enter', event);
}

function handleClear(): void {
  internalValue.value = '';
  emit('update:modelValue', '');
  emit('clear');
}

/**
 * Watch for external model value changes
 */
watch(
  () => props.modelValue,
  newValue => {
    if (newValue !== internalValue.value) {
      internalValue.value = newValue || '';
    }
  },
  { immediate: true }
);
</script>

<style scoped>
/**
 * Banking Input Styles
 * 
 * Professional styling for banking input components with enhanced
 * security indicators and accessibility features.
 */

.banking-input {
  position: relative;
}

/* Currency symbol styling */
.banking-input__currency-symbol {
  color: rgb(var(--v-theme-on-surface-variant));
  font-weight: 500;
  margin-right: 4px;
  user-select: none;
}

/* Percentage symbol styling */
.banking-input__percentage-symbol {
  color: rgb(var(--v-theme-on-surface-variant));
  font-weight: 500;
  margin-left: 4px;
  user-select: none;
}

/* Security icon for sensitive fields */
.banking-input__security-icon {
  color: rgb(var(--v-theme-success));
  opacity: 0.7;
}

/* Sensitive field styling */
.banking-input--sensitive {
  position: relative;
}

.banking-input--sensitive input {
  /* Prevent text selection on sensitive fields */
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

/* Type-specific styling */
.banking-input--currency input {
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.banking-input--percentage input {
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.banking-input--account-number input,
.banking-input--routing-number input {
  font-family: monospace;
  letter-spacing: 0.5px;
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  .banking-input__currency-symbol,
  .banking-input__percentage-symbol {
    font-weight: 700;
  }

  .banking-input--sensitive {
    border: 2px solid rgb(var(--v-theme-warning));
  }
}

/* Mobile optimizations */
@media (max-width: 599px) {
  .banking-input input {
    font-size: 16px; /* Prevent zoom on iOS */
  }
}

/* Focus management for accessibility */
.banking-input :deep(.v-field--focused) {
  outline: 2px solid rgb(var(--v-theme-primary));
  outline-offset: 2px;
}

/* Error state styling */
.banking-input :deep(.v-field--error) {
  border-color: rgb(var(--v-theme-error));
}

/* Loading state styling */
.banking-input :deep(.v-field--loading) {
  opacity: 0.7;
}
</style>
