<template>
  <div class="banking-table-cell" :class="computedClasses">
    <!-- Currency display -->
    <template v-if="type === 'currency'">
      <span class="banking-table-cell__currency">
        {{ formatCurrency(value) }}
      </span>
    </template>

    <!-- Percentage display -->
    <template v-else-if="type === 'percentage'">
      <span class="banking-table-cell__percentage">
        {{ formatPercentage(value) }}
      </span>
    </template>

    <!-- Number display -->
    <template v-else-if="type === 'number'">
      <span class="banking-table-cell__number">
        {{ formatNumber(value) }}
      </span>
    </template>

    <!-- Date display -->
    <template v-else-if="type === 'date' || type === 'datetime'">
      <span class="banking-table-cell__date" :title="formatDateFull(value)">
        {{ formatDate(value) }}
      </span>
    </template>

    <!-- Status display -->
    <template v-else-if="type === 'status'">
      <VChip
        :color="getStatusColor(value)"
        :variant="getStatusVariant(value)"
        size="small"
        class="banking-table-cell__status"
      >
        <VIcon v-if="getStatusIcon(value)" :icon="getStatusIcon(value)" size="small" start />
        {{ getStatusText(value) }}
      </VChip>
    </template>

    <!-- Boolean display -->
    <template v-else-if="type === 'boolean'">
      <div class="banking-table-cell__boolean">
        <VIcon
          :icon="value ? getBooleanIcon(true) : getBooleanIcon(false)"
          :color="value ? getBooleanColor(true) : getBooleanColor(false)"
          size="small"
        />
        <span class="ml-1">{{ getBooleanText(value) }}</span>
      </div>
    </template>

    <!-- Actions display -->
    <template v-else-if="type === 'actions'">
      <div class="banking-table-cell__actions">
        <slot name="actions" :item="item" :value="value">
          <VBtn icon="mdi-dots-vertical" variant="text" size="small" @click="handleActionClick" />
        </slot>
      </div>
    </template>

    <!-- Audit trail display -->
    <template v-else-if="type === 'audit-trail'">
      <div class="banking-table-cell__audit-trail">
        <div class="banking-table-cell__audit-user">
          {{ value?.user || 'Unknown' }}
        </div>
        <div class="banking-table-cell__audit-date">
          {{ formatDate(value?.timestamp) }}
        </div>
        <VIcon
          v-if="value?.action"
          :icon="getAuditIcon(value.action)"
          size="small"
          :color="getAuditColor(value.action)"
          class="ml-2"
        />
      </div>
    </template>

    <!-- Default text display -->
    <template v-else>
      <span class="banking-table-cell__text">
        {{ displayValue }}
      </span>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, type PropType } from 'vue';
import { VChip, VIcon, VBtn } from 'vuetify/components';
import type { BankingTableCellProps } from './BankingTable.types';

/**
 * Banking Table Cell Component
 *
 * Specialized cell component for rendering different data types in banking tables
 * with proper formatting, accessibility, and visual indicators.
 */

// Props
const props = withDefaults(defineProps<BankingTableCellProps>(), {
  type: 'text',
  editable: false,
});

// Events
const emit = defineEmits<{
  (event: 'action-click', item: any): void;
}>();

/**
 * Computed properties
 */
const computedClasses = computed(() => {
  const classes: string[] = ['banking-table-cell'];

  if (props.type) {
    classes.push(`banking-table-cell--${props.type}`);
  }

  if (props.editable) {
    classes.push('banking-table-cell--editable');
  }

  if (props.class) {
    if (typeof props.class === 'string') {
      classes.push(props.class);
    } else if (Array.isArray(props.class)) {
      classes.push(...props.class);
    }
  }

  return classes;
});

const displayValue = computed(() => {
  if (props.value === null || props.value === undefined) {
    return '—';
  }
  return String(props.value);
});

/**
 * Currency formatting
 */
function formatCurrency(value: any): string {
  if (value === null || value === undefined || value === '') return '—';

  const numValue = Number(value);
  if (isNaN(numValue)) return String(value);

  const options = props.format?.currency || { code: 'USD' };

  try {
    return new Intl.NumberFormat(options.locale || 'en-US', {
      style: 'currency',
      currency: options.code,
      minimumFractionDigits: options.minimumFractionDigits ?? 2,
      maximumFractionDigits: options.maximumFractionDigits ?? 2,
    }).format(numValue);
  } catch {
    return `$${numValue.toFixed(2)}`;
  }
}

/**
 * Percentage formatting
 */
function formatPercentage(value: any): string {
  if (value === null || value === undefined || value === '') return '—';

  const numValue = Number(value);
  if (isNaN(numValue)) return String(value);

  const options = props.format?.percentage || {};
  const decimals = options.minimumFractionDigits ?? 2;
  const showSymbol = options.showSymbol ?? true;

  return `${numValue.toFixed(decimals)}${showSymbol ? '%' : ''}`;
}

/**
 * Number formatting
 */
function formatNumber(value: any): string {
  if (value === null || value === undefined || value === '') return '—';

  const numValue = Number(value);
  if (isNaN(numValue)) return String(value);

  const options = props.format?.number || {};

  try {
    return new Intl.NumberFormat(options.locale || 'en-US', {
      minimumFractionDigits: options.minimumFractionDigits ?? 0,
      maximumFractionDigits: options.maximumFractionDigits ?? 2,
      useGrouping: options.useGrouping ?? true,
    }).format(numValue);
  } catch {
    return numValue.toLocaleString();
  }
}

/**
 * Date formatting
 */
function formatDate(value: any): string {
  if (!value) return '—';

  const date = new Date(value);
  if (isNaN(date.getTime())) return String(value);

  const options = props.format?.date;

  if (options?.format) {
    // Custom format would require a date formatting library
    // For now, use built-in formatting
    return date.toLocaleDateString(options.locale || 'en-US');
  }

  return props.type === 'datetime' ? date.toLocaleString() : date.toLocaleDateString();
}

/**
 * Full date formatting for tooltips
 */
function formatDateFull(value: any): string {
  if (!value) return '';

  const date = new Date(value);
  if (isNaN(date.getTime())) return String(value);

  return date.toLocaleString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/**
 * Status formatting
 */
function getStatusColor(value: any): string {
  const colorMap = props.format?.status?.colorMap || {};
  return colorMap[String(value)] || getDefaultStatusColor(value);
}

function getDefaultStatusColor(value: any): string {
  const val = String(value).toLowerCase();

  if (['approved', 'complete', 'success', 'active'].includes(val)) return 'success';
  if (['pending', 'in-progress', 'warning'].includes(val)) return 'warning';
  if (['rejected', 'failed', 'error', 'inactive'].includes(val)) return 'error';
  if (['under-review', 'draft', 'info'].includes(val)) return 'info';

  return 'default';
}

function getStatusVariant(
  value: any
): 'flat' | 'elevated' | 'tonal' | 'outlined' | 'text' | 'plain' {
  const color = getStatusColor(value);
  return color === 'default' ? 'outlined' : 'tonal';
}

function getStatusIcon(value: any): string | undefined {
  const iconMap = props.format?.status?.iconMap || {};
  const customIcon = iconMap[String(value)];

  if (customIcon) return customIcon;

  const val = String(value).toLowerCase();

  if (['approved', 'complete', 'success', 'active'].includes(val)) return 'mdi-check-circle';
  if (['pending', 'in-progress'].includes(val)) return 'mdi-clock-outline';
  if (['rejected', 'failed', 'error'].includes(val)) return 'mdi-alert-circle';
  if (['under-review', 'draft'].includes(val)) return 'mdi-eye-outline';
  if (['warning'].includes(val)) return 'mdi-alert-triangle';

  return undefined;
}

function getStatusText(value: any): string {
  const textMap = props.format?.status?.textMap || {};
  return textMap[String(value)] || String(value);
}

/**
 * Boolean formatting
 */
function getBooleanText(value: any): string {
  const options = props.format?.boolean || {};
  return value ? options.trueText || 'Yes' : options.falseText || 'No';
}

function getBooleanIcon(value: boolean): string {
  const options = props.format?.boolean || {};
  return value ? options.trueIcon || 'mdi-check-circle' : options.falseIcon || 'mdi-close-circle';
}

function getBooleanColor(value: boolean): string {
  const options = props.format?.boolean || {};
  return value ? options.trueColor || 'success' : options.falseColor || 'error';
}

/**
 * Audit trail formatting
 */
function getAuditIcon(action: string): string {
  switch (action?.toLowerCase()) {
    case 'create':
    case 'created':
      return 'mdi-plus-circle';
    case 'update':
    case 'updated':
      return 'mdi-pencil-circle';
    case 'delete':
    case 'deleted':
      return 'mdi-delete-circle';
    case 'approve':
    case 'approved':
      return 'mdi-check-circle';
    case 'reject':
    case 'rejected':
      return 'mdi-cancel';
    case 'review':
    case 'reviewed':
      return 'mdi-eye-circle';
    default:
      return 'mdi-history';
  }
}

function getAuditColor(action: string): string {
  switch (action?.toLowerCase()) {
    case 'create':
    case 'created':
      return 'success';
    case 'update':
    case 'updated':
      return 'info';
    case 'delete':
    case 'deleted':
      return 'error';
    case 'approve':
    case 'approved':
      return 'success';
    case 'reject':
    case 'rejected':
      return 'error';
    case 'review':
    case 'reviewed':
      return 'warning';
    default:
      return 'default';
  }
}

/**
 * Event handlers
 */
function handleActionClick(): void {
  emit('action-click', props.item);
}
</script>

<style scoped>
/**
 * Banking Table Cell Styles
 */

.banking-table-cell {
  position: relative;
  font-variant-numeric: tabular-nums;
}

/* Currency cells */
.banking-table-cell--currency {
  font-family: monospace;
  text-align: right;
}

.banking-table-cell__currency {
  font-weight: 500;
  color: rgb(var(--v-theme-on-surface));
}

/* Percentage cells */
.banking-table-cell--percentage {
  text-align: right;
}

.banking-table-cell__percentage {
  font-weight: 500;
  color: rgb(var(--v-theme-on-surface));
}

/* Number cells */
.banking-table-cell--number {
  text-align: right;
}

.banking-table-cell__number {
  font-weight: 500;
  color: rgb(var(--v-theme-on-surface));
}

/* Date cells */
.banking-table-cell__date {
  color: rgb(var(--v-theme-on-surface-variant));
  font-size: 0.875rem;
}

/* Status cells */
.banking-table-cell__status {
  font-weight: 500;
  text-transform: capitalize;
}

/* Boolean cells */
.banking-table-cell__boolean {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
}

/* Actions cells */
.banking-table-cell__actions {
  display: flex;
  justify-content: center;
  align-items: center;
}

/* Audit trail cells */
.banking-table-cell__audit-trail {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 0.875rem;
}

.banking-table-cell__audit-user {
  font-weight: 500;
  color: rgb(var(--v-theme-on-surface));
}

.banking-table-cell__audit-date {
  color: rgb(var(--v-theme-on-surface-variant));
  font-size: 0.75rem;
}

/* Text cells */
.banking-table-cell__text {
  color: rgb(var(--v-theme-on-surface));
}

/* Editable cells */
.banking-table-cell--editable {
  cursor: pointer;
  border-radius: 4px;
  padding: 4px 8px;
  margin: -4px -8px;
  transition: background-color 0.2s ease;
}

.banking-table-cell--editable:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
}

.banking-table-cell--editable:focus {
  outline: 2px solid rgb(var(--v-theme-primary));
  outline-offset: 2px;
}

/* Responsive adjustments */
@media (max-width: 599px) {
  .banking-table-cell__audit-trail {
    align-items: center;
  }

  .banking-table-cell__boolean {
    justify-content: flex-start;
  }
}

/* High contrast mode */
@media (prefers-contrast: high) {
  .banking-table-cell__currency,
  .banking-table-cell__percentage,
  .banking-table-cell__number {
    font-weight: 700;
  }

  .banking-table-cell--editable {
    border: 1px solid transparent;
  }

  .banking-table-cell--editable:hover {
    border-color: rgb(var(--v-theme-primary));
  }
}
</style>
