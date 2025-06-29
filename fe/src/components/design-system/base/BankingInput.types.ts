/**
 * Banking Input Component Types
 * 
 * Comprehensive type definitions for banking-specific input components
 * with validation, formatting, and accessibility features.
 */

export interface BankingInputProps {
  /**
   * Input value (v-model)
   */
  modelValue?: string | number
  
  /**
   * Input type with banking-specific options
   */
  type?: 'text' | 'email' | 'password' | 'number' | 'currency' | 'percentage' | 'date' | 'phone' | 'account-number' | 'routing-number' | 'ssn'
  
  /**
   * Input label
   */
  label?: string
  
  /**
   * Placeholder text
   */
  placeholder?: string
  
  /**
   * Helper text shown below input
   */
  hint?: string
  
  /**
   * Required field indicator
   */
  required?: boolean
  
  /**
   * Disabled state
   */
  disabled?: boolean
  
  /**
   * Readonly state
   */
  readonly?: boolean
  
  /**
   * Loading state
   */
  loading?: boolean
  
  /**
   * Validation rules
   */
  rules?: ((value: any) => boolean | string)[]
  
  /**
   * Error state
   */
  error?: boolean
  
  /**
   * Error messages
   */
  errorMessages?: string[]
  
  /**
   * Input variant
   */
  variant?: 'filled' | 'outlined' | 'underlined' | 'solo' | 'solo-inverted'
  
  /**
   * Input density
   */
  density?: 'default' | 'comfortable' | 'compact'
  
  /**
   * Maximum length
   */
  maxLength?: number
  
  /**
   * Minimum value (for numeric inputs)
   */
  min?: number
  
  /**
   * Maximum value (for numeric inputs)
   */
  max?: number
  
  /**
   * Step value (for numeric inputs)
   */
  step?: number
  
  /**
   * Decimal places for currency/percentage
   */
  decimalPlaces?: number
  
  /**
   * Currency code for currency inputs
   */
  currency?: string
  
  /**
   * Show clear button
   */
  clearable?: boolean
  
  /**
   * Prepend icon
   */
  prependIcon?: string
  
  /**
   * Append icon
   */
  appendIcon?: string
  
  /**
   * Auto-format input (currency, phone, etc.)
   */
  autoFormat?: boolean
  
  /**
   * Mask pattern for input formatting
   */
  mask?: string
  
  /**
   * Validation on blur
   */
  validateOnBlur?: boolean
  
  /**
   * Show validation immediately
   */
  immediateValidation?: boolean
  
  /**
   * Custom CSS classes
   */
  class?: string | string[] | Record<string, boolean>
  
  /**
   * Accessibility label when label is not descriptive enough
   */
  ariaLabel?: string
  
  /**
   * ARIA describedby for additional context
   */
  ariaDescribedby?: string
  
  /**
   * Autocomplete attribute
   */
  autocomplete?: string
  
  /**
   * Input mode for mobile keyboards
   */
  inputmode?: 'none' | 'text' | 'decimal' | 'numeric' | 'tel' | 'search' | 'email' | 'url'
}

/**
 * Banking input validation rules
 */
export interface ValidationResult {
  valid: boolean
  message?: string
}

/**
 * Banking-specific validation functions
 */
export interface BankingValidators {
  email: (value: string) => ValidationResult
  currency: (value: string | number, min?: number, max?: number) => ValidationResult
  percentage: (value: string | number) => ValidationResult
  accountNumber: (value: string) => ValidationResult
  routingNumber: (value: string) => ValidationResult
  ssn: (value: string) => ValidationResult
  phone: (value: string) => ValidationResult
  required: (value: any) => ValidationResult
  minLength: (min: number) => (value: string) => ValidationResult
  maxLength: (max: number) => (value: string) => ValidationResult
  pattern: (pattern: RegExp, message: string) => (value: string) => ValidationResult
}

/**
 * Component emits interface
 */
export interface BankingInputEmits {
  /**
   * Emitted when input value changes
   */
  (event: 'update:modelValue', value: string | number): void
  
  /**
   * Emitted on input event
   */
  (event: 'input', value: string | number): void
  
  /**
   * Emitted on blur event
   */
  (event: 'blur', focusEvent: FocusEvent): void
  
  /**
   * Emitted on focus event
   */
  (event: 'focus', focusEvent: FocusEvent): void
  
  /**
   * Emitted when validation state changes
   */
  (event: 'validation-change', valid: boolean, errors: string[]): void
  
  /**
   * Emitted when input is cleared
   */
  (event: 'clear'): void
  
  /**
   * Emitted on Enter key press
   */
  (event: 'enter', keyboardEvent: KeyboardEvent): void
}

/**
 * Component slots interface
 */
export interface BankingInputSlots {
  /**
   * Prepend slot
   */
  prepend?: () => any
  
  /**
   * Prepend inner slot
   */
  'prepend-inner'?: () => any
  
  /**
   * Append slot
   */
  append?: () => any
  
  /**
   * Append inner slot
   */
  'append-inner'?: () => any
  
  /**
   * Details slot (for custom error/hint display)
   */
  details?: () => any
  
  /**
   * Label slot
   */
  label?: () => any
}

/**
 * Formatting options for different input types
 */
export interface FormattingOptions {
  /**
   * Currency formatting options
   */
  currency: {
    locale?: string
    currency: string
    minimumFractionDigits?: number
    maximumFractionDigits?: number
  }
  
  /**
   * Percentage formatting options
   */
  percentage: {
    minimumFractionDigits?: number
    maximumFractionDigits?: number
    showSymbol?: boolean
  }
  
  /**
   * Phone formatting options
   */
  phone: {
    country?: string
    format?: 'national' | 'international'
  }
  
  /**
   * Account number formatting options
   */
  accountNumber: {
    showLastFour?: boolean
    maskCharacter?: string
  }
}