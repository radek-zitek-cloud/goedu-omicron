<template>
  <div class="audit-trail" :class="computedClasses">
    <!-- Header -->
    <div v-if="showHeader" class="audit-trail__header">
      <div class="audit-trail__title-section">
        <h3 class="audit-trail__title">{{ title || 'Audit Trail' }}</h3>
        <p v-if="subtitle" class="audit-trail__subtitle">{{ subtitle }}</p>
      </div>
      
      <div class="audit-trail__actions">
        <slot name="header-actions">
          <BankingButton
            v-if="exportable"
            variant="outlined"
            icon="mdi-download"
            size="small"
            @click="handleExport"
          >
            Export
          </BankingButton>
          
          <BankingButton
            v-if="refreshable"
            variant="text"
            icon="mdi-refresh"
            size="small"
            @click="handleRefresh"
          >
            Refresh
          </BankingButton>
        </slot>
      </div>
    </div>

    <!-- Filters -->
    <div v-if="filterable" class="audit-trail__filters">
      <div class="audit-trail__filter-row">
        <BankingInput
          v-model="filters.user"
          type="text"
          label="User"
          variant="outlined"
          density="compact"
          class="audit-trail__filter-input"
          clearable
          @input="applyFilters"
        />
        
        <BankingInput
          v-model="filters.action"
          type="text"
          label="Action"
          variant="outlined"
          density="compact"
          class="audit-trail__filter-input"
          clearable
          @input="applyFilters"
        />
        
        <BankingInput
          v-model="filters.dateFrom"
          type="date"
          label="From Date"
          variant="outlined"
          density="compact"
          class="audit-trail__filter-input"
          @input="applyFilters"
        />
        
        <BankingInput
          v-model="filters.dateTo"
          type="date"
          label="To Date"
          variant="outlined"
          density="compact"
          class="audit-trail__filter-input"
          @input="applyFilters"
        />
      </div>
      
      <div class="audit-trail__filter-actions">
        <BankingButton
          variant="text"
          size="small"
          @click="clearFilters"
        >
          Clear Filters
        </BankingButton>
      </div>
    </div>

    <!-- Entry Count -->
    <div v-if="showSummary" class="audit-trail__summary">
      <div class="audit-trail__count">
        <VIcon icon="mdi-history" size="small" class="mr-2" />
        {{ filteredEntries.length }} {{ filteredEntries.length === 1 ? 'entry' : 'entries' }}
        <span v-if="totalEntries !== filteredEntries.length">
          ({{ totalEntries }} total)
        </span>
      </div>
      
      <div v-if="dateRange" class="audit-trail__date-range">
        {{ formatDateRange() }}
      </div>
    </div>

    <!-- Audit Entries -->
    <div class="audit-trail__entries">
      <div v-if="loading" class="audit-trail__loading">
        <VSkeletonLoader
          v-for="n in 5"
          :key="n"
          type="list-item-two-line"
          class="mb-2"
        />
      </div>
      
      <div v-else-if="filteredEntries.length === 0" class="audit-trail__no-data">
        <VIcon icon="mdi-history" size="48" class="mb-4" />
        <h4>No Audit Entries</h4>
        <p>{{ noDataText || 'No audit entries found matching your criteria.' }}</p>
      </div>
      
      <div v-else>
        <div
          v-for="(entry, index) in paginatedEntries"
          :key="entry.id || index"
          class="audit-trail__entry"
          :class="getEntryClasses(entry)"
        >
          <!-- Entry Icon -->
          <div class="audit-trail__entry-icon">
            <VAvatar
              :color="getEntryColor(entry)"
              :variant="getEntryVariant(entry)"
              size="32"
            >
              <VIcon
                :icon="getEntryIcon(entry)"
                size="small"
                :color="getEntryIconColor(entry)"
              />
            </VAvatar>
          </div>

          <!-- Entry Content -->
          <div class="audit-trail__entry-content">
            <div class="audit-trail__entry-header">
              <div class="audit-trail__entry-main">
                <h4 class="audit-trail__entry-action">{{ entry.action }}</h4>
                <p class="audit-trail__entry-description">
                  {{ getEntryDescription(entry) }}
                </p>
              </div>
              
              <div class="audit-trail__entry-meta">
                <div class="audit-trail__entry-user">
                  <VIcon icon="mdi-account" size="small" />
                  {{ entry.user }}
                </div>
                <div class="audit-trail__entry-timestamp">
                  <VIcon icon="mdi-clock-outline" size="small" />
                  {{ formatTimestamp(entry.timestamp) }}
                </div>
              </div>
            </div>

            <!-- Entry Details -->
            <div v-if="entry.details && showDetails" class="audit-trail__entry-details">
              <VExpansionPanels variant="accordion" flat>
                <VExpansionPanel>
                  <VExpansionPanelTitle>
                    <VIcon icon="mdi-information-outline" size="small" class="mr-2" />
                    View Details
                  </VExpansionPanelTitle>
                  <VExpansionPanelText>
                    <div class="audit-trail__entry-detail-grid">
                      <div
                        v-for="(value, key) in entry.details"
                        :key="`${entry.id}-${key}`"
                        class="audit-trail__entry-detail-item"
                      >
                        <strong class="audit-trail__entry-detail-key">{{ formatDetailKey(key) }}:</strong>
                        <span class="audit-trail__entry-detail-value">{{ formatDetailValue(value) }}</span>
                      </div>
                    </div>
                  </VExpansionPanelText>
                </VExpansionPanel>
              </VExpansionPanels>
            </div>

            <!-- Entry Changes -->
            <div v-if="entry.changes && entry.changes.length > 0" class="audit-trail__entry-changes">
              <VExpansionPanels variant="accordion" flat>
                <VExpansionPanel>
                  <VExpansionPanelTitle>
                    <VIcon icon="mdi-delta" size="small" class="mr-2" />
                    Changes ({{ entry.changes.length }})
                  </VExpansionPanelTitle>
                  <VExpansionPanelText>
                    <div class="audit-trail__changes-list">
                      <div
                        v-for="(change, changeIndex) in entry.changes"
                        :key="`${entry.id}-change-${changeIndex}`"
                        class="audit-trail__change-item"
                      >
                        <div class="audit-trail__change-field">
                          <strong>{{ formatDetailKey(change.field) }}</strong>
                        </div>
                        <div class="audit-trail__change-values">
                          <div class="audit-trail__change-old">
                            <span class="audit-trail__change-label">From:</span>
                            <code class="audit-trail__change-value">{{ formatDetailValue(change.oldValue) }}</code>
                          </div>
                          <VIcon icon="mdi-arrow-right" size="small" class="audit-trail__change-arrow" />
                          <div class="audit-trail__change-new">
                            <span class="audit-trail__change-label">To:</span>
                            <code class="audit-trail__change-value">{{ formatDetailValue(change.newValue) }}</code>
                          </div>
                        </div>
                      </div>
                    </div>
                  </VExpansionPanelText>
                </VExpansionPanel>
              </VExpansionPanels>
            </div>

            <!-- Entry Risk Assessment -->
            <div v-if="entry.riskLevel" class="audit-trail__entry-risk">
              <VChip
                :color="getRiskColor(entry.riskLevel)"
                size="small"
                variant="tonal"
              >
                <VIcon
                  :icon="getRiskIcon(entry.riskLevel)"
                  size="small"
                  start
                />
                {{ entry.riskLevel }} Risk
              </VChip>
            </div>

            <!-- Entry IP and Location -->
            <div v-if="entry.metadata" class="audit-trail__entry-metadata">
              <div v-if="entry.metadata.ipAddress" class="audit-trail__metadata-item">
                <VIcon icon="mdi-ip-network" size="small" />
                <span>{{ entry.metadata.ipAddress }}</span>
              </div>
              <div v-if="entry.metadata.location" class="audit-trail__metadata-item">
                <VIcon icon="mdi-map-marker" size="small" />
                <span>{{ entry.metadata.location }}</span>
              </div>
              <div v-if="entry.metadata.userAgent" class="audit-trail__metadata-item">
                <VIcon icon="mdi-monitor" size="small" />
                <span>{{ entry.metadata.userAgent }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="audit-trail__pagination">
      <VPagination
        v-model="currentPage"
        :length="totalPages"
        :total-visible="7"
        @update:model-value="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, type PropType } from 'vue'
import {
  VIcon,
  VAvatar,
  VSkeletonLoader,
  VExpansionPanels,
  VExpansionPanel,
  VExpansionPanelTitle,
  VExpansionPanelText,
  VChip,
  VPagination
} from 'vuetify/components'
import BankingButton from '../base/BankingButton.vue'
import BankingInput from '../base/BankingInput.vue'
import type {
  AuditTrailProps,
  AuditTrailEmits,
  AuditEntry,
  AuditFilters
} from './AuditTrail.types'

/**
 * Audit Trail Component
 * 
 * A comprehensive audit trail component for banking applications with
 * advanced filtering, detailed change tracking, and compliance features.
 * 
 * Features:
 * - Detailed audit entry display with change tracking
 * - Advanced filtering and search capabilities
 * - Risk level assessment and IP tracking
 * - Export functionality for compliance
 * - WCAG 2.1 AA compliant
 * - Responsive design
 * 
 * @example
 * <AuditTrail
 *   :entries="auditEntries"
 *   title="Transaction Audit Trail"
 *   filterable
 *   exportable
 *   @entry-selected="viewEntryDetails"
 * />
 */

// Props with defaults
const props = withDefaults(defineProps<AuditTrailProps>(), {
  showHeader: true,
  showSummary: true,
  showDetails: true,
  filterable: false,
  exportable: false,
  refreshable: false,
  itemsPerPage: 20,
  groupByDate: false
})

// Events
const emit = defineEmits<AuditTrailEmits>()

// Reactive state
const currentPage = ref(1)
const filters = ref<AuditFilters>({
  user: '',
  action: '',
  dateFrom: '',
  dateTo: '',
  riskLevel: '',
  ipAddress: ''
})

/**
 * Computed properties
 */
const computedClasses = computed(() => {
  const classes: string[] = ['audit-trail']
  
  if (props.variant) {
    classes.push(`audit-trail--${props.variant}`)
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

const filteredEntries = computed(() => {
  let entries = props.entries || []
  
  // Apply filters
  if (filters.value.user) {
    entries = entries.filter(entry => 
      entry.user.toLowerCase().includes(filters.value.user.toLowerCase())
    )
  }
  
  if (filters.value.action) {
    entries = entries.filter(entry => 
      entry.action.toLowerCase().includes(filters.value.action.toLowerCase())
    )
  }
  
  if (filters.value.dateFrom) {
    const fromDate = new Date(filters.value.dateFrom)
    entries = entries.filter(entry => 
      new Date(entry.timestamp) >= fromDate
    )
  }
  
  if (filters.value.dateTo) {
    const toDate = new Date(filters.value.dateTo)
    toDate.setHours(23, 59, 59, 999) // End of day
    entries = entries.filter(entry => 
      new Date(entry.timestamp) <= toDate
    )
  }
  
  if (filters.value.riskLevel) {
    entries = entries.filter(entry => 
      entry.riskLevel === filters.value.riskLevel
    )
  }
  
  if (filters.value.ipAddress) {
    entries = entries.filter(entry => 
      entry.metadata?.ipAddress?.includes(filters.value.ipAddress)
    )
  }
  
  // Sort by timestamp (newest first)
  return entries.sort((a, b) => 
    new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  )
})

const totalEntries = computed(() => props.entries?.length || 0)

const totalPages = computed(() => 
  Math.ceil(filteredEntries.value.length / props.itemsPerPage)
)

const paginatedEntries = computed(() => {
  const start = (currentPage.value - 1) * props.itemsPerPage
  const end = start + props.itemsPerPage
  return filteredEntries.value.slice(start, end)
})

const dateRange = computed(() => {
  if (filteredEntries.value.length === 0) return null
  
  const dates = filteredEntries.value.map(entry => new Date(entry.timestamp))
  const minDate = new Date(Math.min(...dates.map(d => d.getTime())))
  const maxDate = new Date(Math.max(...dates.map(d => d.getTime())))
  
  return { min: minDate, max: maxDate }
})

/**
 * Entry styling helpers
 */
function getEntryClasses(entry: AuditEntry): string[] {
  const classes: string[] = ['audit-trail__entry--default']
  
  if (entry.riskLevel) {
    classes.push(`audit-trail__entry--${entry.riskLevel.toLowerCase()}-risk`)
  }
  
  if (entry.category) {
    classes.push(`audit-trail__entry--${entry.category}`)
  }
  
  return classes
}

function getEntryColor(entry: AuditEntry): string {
  if (entry.riskLevel) {
    switch (entry.riskLevel.toLowerCase()) {
      case 'high': return 'error'
      case 'medium': return 'warning'
      case 'low': return 'info'
    }
  }
  
  switch (entry.action.toLowerCase()) {
    case 'create':
    case 'created': return 'success'
    case 'update':
    case 'updated': return 'info'
    case 'delete':
    case 'deleted': return 'error'
    case 'login':
    case 'logout': return 'primary'
    case 'approve':
    case 'approved': return 'success'
    case 'reject':
    case 'rejected': return 'error'
    default: return 'default'
  }
}

function getEntryVariant(entry: AuditEntry): 'flat' | 'elevated' | 'tonal' | 'outlined' | 'text' | 'plain' {
  const color = getEntryColor(entry)
  return color === 'default' ? 'outlined' : 'tonal'
}

function getEntryIcon(entry: AuditEntry): string {
  switch (entry.action.toLowerCase()) {
    case 'create':
    case 'created': return 'mdi-plus-circle'
    case 'update':
    case 'updated': return 'mdi-pencil'
    case 'delete':
    case 'deleted': return 'mdi-delete'
    case 'login': return 'mdi-login'
    case 'logout': return 'mdi-logout'
    case 'approve':
    case 'approved': return 'mdi-check-circle'
    case 'reject':
    case 'rejected': return 'mdi-cancel'
    case 'view':
    case 'viewed': return 'mdi-eye'
    case 'download':
    case 'downloaded': return 'mdi-download'
    case 'upload':
    case 'uploaded': return 'mdi-upload'
    default: return 'mdi-history'
  }
}

function getEntryIconColor(entry: AuditEntry): string {
  const color = getEntryColor(entry)
  return color === 'default' ? 'default' : 'white'
}

function getRiskColor(riskLevel: string): string {
  switch (riskLevel.toLowerCase()) {
    case 'high': return 'error'
    case 'medium': return 'warning'
    case 'low': return 'info'
    default: return 'default'
  }
}

function getRiskIcon(riskLevel: string): string {
  switch (riskLevel.toLowerCase()) {
    case 'high': return 'mdi-alert-circle'
    case 'medium': return 'mdi-alert-triangle'
    case 'low': return 'mdi-information'
    default: return 'mdi-help-circle'
  }
}

/**
 * Formatting helpers
 */
function getEntryDescription(entry: AuditEntry): string {
  if (entry.description) return entry.description
  
  // Generate description based on action and target
  if (entry.target) {
    return `${entry.action} ${entry.target}`
  }
  
  return `${entry.action} performed`
}

function formatTimestamp(timestamp: Date | string): string {
  const date = new Date(timestamp)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatDateRange(): string {
  if (!dateRange.value) return ''
  
  const { min, max } = dateRange.value
  
  if (min.toDateString() === max.toDateString()) {
    return min.toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    })
  }
  
  return `${min.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric'
  })} - ${max.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })}`
}

function formatDetailKey(key: string): string {
  return key.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase())
}

function formatDetailValue(value: any): string {
  if (value === null || value === undefined) return 'N/A'
  if (typeof value === 'boolean') return value ? 'Yes' : 'No'
  if (typeof value === 'object') return JSON.stringify(value, null, 2)
  return String(value)
}

/**
 * Event handlers
 */
function applyFilters(): void {
  currentPage.value = 1
  emit('filters-changed', filters.value)
}

function clearFilters(): void {
  filters.value = {
    user: '',
    action: '',
    dateFrom: '',
    dateTo: '',
    riskLevel: '',
    ipAddress: ''
  }
  currentPage.value = 1
  emit('filters-cleared')
}

function handlePageChange(page: number): void {
  emit('page-changed', page)
}

function handleExport(): void {
  emit('export-requested', {
    entries: filteredEntries.value,
    filters: filters.value,
    timestamp: new Date()
  })
}

function handleRefresh(): void {
  emit('refresh-requested')
}

/**
 * Watch for filter changes
 */
watch(filters, () => {
  applyFilters()
}, { deep: true })
</script>

<style scoped>
/**
 * Audit Trail Styles
 * 
 * Professional styling for audit trail components with clear
 * visual hierarchy and accessibility features.
 */

.audit-trail {
  background: rgb(var(--v-theme-surface));
  border-radius: 8px;
  overflow: hidden;
}

/* Header */
.audit-trail__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px 24px;
  border-bottom: 1px solid rgb(var(--v-theme-border));
  background: rgb(var(--v-theme-surface-bright));
}

.audit-trail__title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
  color: rgb(var(--v-theme-on-surface));
}

.audit-trail__subtitle {
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface-variant));
  margin: 4px 0 0;
}

.audit-trail__actions {
  display: flex;
  gap: 8px;
}

/* Filters */
.audit-trail__filters {
  padding: 16px 24px;
  border-bottom: 1px solid rgb(var(--v-theme-border));
  background: rgb(var(--v-theme-surface-variant));
}

.audit-trail__filter-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 12px;
}

.audit-trail__filter-actions {
  display: flex;
  justify-content: flex-end;
}

/* Summary */
.audit-trail__summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  background: rgba(var(--v-theme-surface-variant), 0.3);
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface-variant));
}

.audit-trail__count {
  display: flex;
  align-items: center;
  font-weight: 500;
}

/* Entries */
.audit-trail__entries {
  padding: 16px 24px;
}

.audit-trail__loading,
.audit-trail__no-data {
  padding: 48px 24px;
  text-align: center;
}

.audit-trail__no-data {
  color: rgb(var(--v-theme-on-surface-variant));
}

.audit-trail__entry {
  display: flex;
  margin-bottom: 24px;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid rgb(var(--v-theme-border));
  background: rgb(var(--v-theme-surface));
  transition: all 0.2s ease;
}

.audit-trail__entry:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.audit-trail__entry:last-child {
  margin-bottom: 0;
}

.audit-trail__entry-icon {
  margin-right: 16px;
  flex-shrink: 0;
}

.audit-trail__entry-content {
  flex: 1;
  min-width: 0;
}

.audit-trail__entry-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.audit-trail__entry-main {
  flex: 1;
  min-width: 0;
}

.audit-trail__entry-action {
  font-size: 1rem;
  font-weight: 600;
  margin: 0 0 4px;
  color: rgb(var(--v-theme-on-surface));
}

.audit-trail__entry-description {
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface-variant));
  margin: 0;
}

.audit-trail__entry-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  text-align: right;
  flex-shrink: 0;
  margin-left: 16px;
}

.audit-trail__entry-user,
.audit-trail__entry-timestamp {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  color: rgb(var(--v-theme-on-surface-variant));
}

/* Entry details */
.audit-trail__entry-details,
.audit-trail__entry-changes {
  margin-top: 12px;
}

.audit-trail__entry-detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 8px;
}

.audit-trail__entry-detail-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.audit-trail__entry-detail-key {
  font-size: 0.75rem;
  color: rgb(var(--v-theme-on-surface-variant));
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.audit-trail__entry-detail-value {
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface));
  word-break: break-word;
}

/* Changes */
.audit-trail__changes-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.audit-trail__change-item {
  padding: 12px;
  background: rgba(var(--v-theme-surface-variant), 0.3);
  border-radius: 6px;
}

.audit-trail__change-field {
  margin-bottom: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(var(--v-theme-on-surface));
}

.audit-trail__change-values {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.audit-trail__change-old,
.audit-trail__change-new {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 150px;
}

.audit-trail__change-label {
  font-size: 0.75rem;
  color: rgb(var(--v-theme-on-surface-variant));
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.audit-trail__change-value {
  background: rgba(var(--v-theme-surface), 0.8);
  padding: 4px 8px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.875rem;
  border: 1px solid rgb(var(--v-theme-border));
}

.audit-trail__change-arrow {
  color: rgb(var(--v-theme-on-surface-variant));
  flex-shrink: 0;
}

/* Risk and metadata */
.audit-trail__entry-risk {
  margin-top: 8px;
}

.audit-trail__entry-metadata {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 8px;
  font-size: 0.75rem;
  color: rgb(var(--v-theme-on-surface-variant));
}

.audit-trail__metadata-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* Pagination */
.audit-trail__pagination {
  display: flex;
  justify-content: center;
  padding: 16px 24px;
  border-top: 1px solid rgb(var(--v-theme-border));
  background: rgb(var(--v-theme-surface-bright));
}

/* Risk level variants */
.audit-trail__entry--high-risk {
  border-left: 4px solid rgb(var(--v-theme-error));
}

.audit-trail__entry--medium-risk {
  border-left: 4px solid rgb(var(--v-theme-warning));
}

.audit-trail__entry--low-risk {
  border-left: 4px solid rgb(var(--v-theme-info));
}

/* Responsive design */
@media (max-width: 959px) {
  .audit-trail__header {
    flex-direction: column;
    gap: 16px;
  }
  
  .audit-trail__filter-row {
    grid-template-columns: 1fr;
  }
  
  .audit-trail__entry-header {
    flex-direction: column;
    gap: 12px;
  }
  
  .audit-trail__entry-meta {
    text-align: left;
    margin-left: 0;
  }
  
  .audit-trail__change-values {
    flex-direction: column;
    align-items: stretch;
  }
  
  .audit-trail__change-arrow {
    transform: rotate(90deg);
    align-self: center;
  }
}

@media (max-width: 599px) {
  .audit-trail__entries {
    padding: 12px 16px;
  }
  
  .audit-trail__entry {
    padding: 12px;
  }
  
  .audit-trail__entry-icon {
    margin-right: 12px;
  }
  
  .audit-trail__entry-metadata {
    flex-direction: column;
    gap: 8px;
  }
}

/* High contrast mode */
@media (prefers-contrast: high) {
  .audit-trail__entry {
    border-width: 2px;
  }
  
  .audit-trail__entry--high-risk,
  .audit-trail__entry--medium-risk,
  .audit-trail__entry--low-risk {
    border-left-width: 6px;
  }
}
</style>