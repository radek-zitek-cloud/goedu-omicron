<template>
  <div class="banking-table" :class="computedClasses">
    <!-- Table Header with Actions -->
    <div v-if="showHeader" class="banking-table__header">
      <div class="banking-table__title-section">
        <h3 v-if="title" class="banking-table__title">{{ title }}</h3>
        <p v-if="subtitle" class="banking-table__subtitle">{{ subtitle }}</p>
      </div>
      
      <div class="banking-table__actions">
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

    <!-- Search and Filters -->
    <div v-if="searchable || filterable" class="banking-table__controls">
      <BankingInput
        v-if="searchable"
        v-model="searchQuery"
        type="text"
        label="Search"
        prepend-icon="mdi-magnify"
        variant="outlined"
        density="compact"
        class="banking-table__search"
        clearable
        @input="handleSearch"
      />
      
      <div v-if="filterable" class="banking-table__filters">
        <slot name="filters" />
      </div>
    </div>

    <!-- Data Table -->
    <VDataTable
      :headers="computedHeaders"
      :items="computedItems"
      :loading="loading"
      :items-per-page="itemsPerPage"
      :items-per-page-options="itemsPerPageOptions"
      :sort-by="sortBy"
      :show-select="selectable"
      :show-expand="expandable"
      :expand-on-click="expandOnClick"
      :must-sort="mustSort"
      :multi-sort="multiSort"
      :mobile-breakpoint="mobileBreakpoint"
      :density="density"
      :class="tableClasses"
      :items-length="totalItems"
      :server-items-length="serverMode ? totalItems : undefined"
      @update:options="handleOptionsUpdate"
      @update:items-per-page="handleItemsPerPageUpdate"
      @update:sort-by="handleSortUpdate"
      @click:row="handleRowClick"
      v-model:expanded="expandedItems"
      v-model="selectedItems"
    >
      <!-- Pass through all custom column slots -->
      <template
        v-for="header in headers"
        :key="header.key"
        #[`item.${header.key}`]="{ item, value }"
      >
        <slot
          :name="`item.${header.key}`"
          :item="item"
          :value="value"
          :header="header"
        >
          <!-- Banking-specific column rendering -->
          <BankingTableCell
            :value="value"
            :type="header.type"
            :format="header.format"
            :item="item"
            :header="header"
          />
        </slot>
      </template>

      <!-- Expanded row content -->
      <template v-if="expandable" #expanded-row="{ item }">
        <slot name="expanded-row" :item="item">
          <tr>
            <td :colspan="computedHeaders.length + (selectable ? 1 : 0)">
              <div class="banking-table__expanded-content">
                <p>No expanded content provided</p>
              </div>
            </td>
          </tr>
        </slot>
      </template>

      <!-- Loading state -->
      <template #loading>
        <slot name="loading">
          <VSkeletonLoader
            type="table-row-divider@6"
            class="banking-table__loading"
          />
        </slot>
      </template>

      <!-- No data state -->
      <template #no-data>
        <slot name="no-data">
          <div class="banking-table__no-data">
            <VIcon icon="mdi-database-off" size="48" class="mb-4" />
            <h4>No Data Available</h4>
            <p>{{ noDataText || 'No records found matching your criteria.' }}</p>
          </div>
        </slot>
      </template>

      <!-- Bottom pagination -->
      <template #bottom>
        <div class="banking-table__footer">
          <div class="banking-table__info">
            <span class="text-caption">
              Showing {{ startIndex + 1 }}-{{ endIndex }} of {{ totalItems }} items
            </span>
          </div>
          
          <VPagination
            v-if="totalPages > 1"
            :model-value="currentPage"
            :length="totalPages"
            :total-visible="totalVisible"
            @update:model-value="handlePageUpdate"
          />
        </div>
      </template>
    </VDataTable>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, type PropType } from 'vue'
import {
  VDataTable,
  VIcon,
  VSkeletonLoader,
  VPagination
} from 'vuetify/components'
import BankingButton from './BankingButton.vue'
import BankingInput from './BankingInput.vue'
import BankingTableCell from './BankingTableCell.vue'
import type {
  BankingTableProps,
  BankingTableEmits,
  BankingTableHeader,
  BankingTableItem,
  TableOptions
} from './BankingTable.types'

/**
 * Banking Data Table Component
 * 
 * A comprehensive data table component designed for banking applications
 * with advanced filtering, sorting, pagination, and accessibility features.
 * 
 * Features:
 * - Banking-specific column types (currency, percentage, status)
 * - Advanced search and filtering
 * - Export functionality
 * - Responsive design
 * - WCAG 2.1 AA compliant
 * - Audit trail integration
 * 
 * @example
 * <BankingTable
 *   :headers="transactionHeaders"
 *   :items="transactions"
 *   title="Transaction History"
 *   searchable
 *   exportable
 *   @row-click="viewTransaction"
 * />
 */

// Props with banking-specific defaults
const props = withDefaults(defineProps<BankingTableProps>(), {
  density: 'default',
  itemsPerPage: 25,
  itemsPerPageOptions: () => [10, 25, 50, 100],
  showHeader: true,
  searchable: false,
  filterable: false,
  sortable: true,
  selectable: false,
  expandable: false,
  expandOnClick: false,
  exportable: false,
  refreshable: false,
  mustSort: false,
  multiSort: false,
  mobileBreakpoint: 960,
  totalVisible: 7,
  serverMode: false
})

// Events
const emit = defineEmits<BankingTableEmits>()

// Reactive state
const searchQuery = ref('')
const selectedItems = ref<BankingTableItem[]>([])
const expandedItems = ref<string[]>([])
const currentPage = ref(1)

/**
 * Computed properties
 */
const computedHeaders = computed(() => {
  return props.headers.map(header => ({
    ...header,
    sortable: header.sortable !== false && props.sortable,
    align: header.align || getDefaultAlign(header.type)
  }))
})

const computedItems = computed(() => {
  let items = props.items || []
  
  // Apply search filter
  if (searchQuery.value && props.searchable) {
    const query = searchQuery.value.toLowerCase()
    items = items.filter(item => 
      Object.values(item).some(value => 
        String(value).toLowerCase().includes(query)
      )
    )
  }
  
  return items
})

const computedClasses = computed(() => {
  const classes: string[] = ['banking-table']
  
  if (props.variant) {
    classes.push(`banking-table--${props.variant}`)
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

const tableClasses = computed(() => [
  'banking-table__data-table',
  { 'banking-table__data-table--striped': props.striped },
  { 'banking-table__data-table--hover': props.hover }
])

const totalItems = computed(() => 
  props.serverMode ? (props.totalItems || 0) : computedItems.value.length
)

const totalPages = computed(() => 
  Math.ceil(totalItems.value / props.itemsPerPage)
)

const startIndex = computed(() => 
  (currentPage.value - 1) * props.itemsPerPage
)

const endIndex = computed(() => 
  Math.min(startIndex.value + props.itemsPerPage, totalItems.value)
)

/**
 * Get default column alignment based on type
 */
function getDefaultAlign(type?: string): 'start' | 'center' | 'end' {
  switch (type) {
    case 'currency':
    case 'percentage':
    case 'number':
      return 'end'
    case 'status':
    case 'boolean':
    case 'actions':
      return 'center'
    default:
      return 'start'
  }
}

/**
 * Event handlers
 */
function handleSearch(value: string | number): void {
  const searchValue = String(value)
  searchQuery.value = searchValue
  currentPage.value = 1
  emit('search', searchValue)
}

function handleRowClick(event: Event, { item }: { item: BankingTableItem }): void {
  emit('row-click', item, event)
}

function handleOptionsUpdate(options: TableOptions): void {
  emit('options-update', options)
}

function handleItemsPerPageUpdate(itemsPerPage: number): void {
  emit('items-per-page-update', itemsPerPage)
}

function handleSortUpdate(sortBy: any[]): void {
  emit('sort-update', sortBy)
}

function handlePageUpdate(page: number): void {
  currentPage.value = page
  emit('page-update', page)
}

function handleExport(): void {
  emit('export', computedItems.value)
}

function handleRefresh(): void {
  emit('refresh')
}

/**
 * Watch for external changes
 */
watch(() => props.page, (newPage) => {
  if (newPage !== undefined && newPage !== currentPage.value) {
    currentPage.value = newPage
  }
}, { immediate: true })
</script>

<style scoped>
/**
 * Banking Table Styles
 * 
 * Professional styling for banking data tables with enhanced
 * accessibility and responsive design.
 */

.banking-table {
  background: rgb(var(--v-theme-surface));
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* Header section */
.banking-table__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px 24px 16px;
  border-bottom: 1px solid rgb(var(--v-theme-border));
  background: rgb(var(--v-theme-surface-bright));
}

.banking-table__title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
  color: rgb(var(--v-theme-on-surface));
}

.banking-table__subtitle {
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface-variant));
  margin: 4px 0 0;
}

.banking-table__actions {
  display: flex;
  gap: 8px;
}

/* Controls section */
.banking-table__controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  gap: 16px;
  border-bottom: 1px solid rgb(var(--v-theme-border));
  background: rgb(var(--v-theme-surface-variant));
}

.banking-table__search {
  min-width: 300px;
  max-width: 400px;
}

.banking-table__filters {
  display: flex;
  gap: 12px;
}

/* Data table styling */
.banking-table__data-table {
  background: transparent;
}

.banking-table__data-table--striped :deep(tbody tr:nth-child(even)) {
  background: rgba(var(--v-theme-on-surface), 0.02);
}

.banking-table__data-table--hover :deep(tbody tr:hover) {
  background: rgba(var(--v-theme-primary), 0.04);
  cursor: pointer;
}

/* Expanded content */
.banking-table__expanded-content {
  padding: 16px 24px;
  background: rgba(var(--v-theme-surface-variant), 0.5);
  border-top: 1px solid rgb(var(--v-theme-border));
}

/* Loading state */
.banking-table__loading {
  padding: 24px;
}

/* No data state */
.banking-table__no-data {
  text-align: center;
  padding: 48px 24px;
  color: rgb(var(--v-theme-on-surface-variant));
}

/* Footer */
.banking-table__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  border-top: 1px solid rgb(var(--v-theme-border));
  background: rgb(var(--v-theme-surface-bright));
}

.banking-table__info {
  font-size: 0.875rem;
  color: rgb(var(--v-theme-on-surface-variant));
}

/* Responsive design */
@media (max-width: 959px) {
  .banking-table__header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .banking-table__controls {
    flex-direction: column;
    align-items: stretch;
  }
  
  .banking-table__search {
    min-width: unset;
    max-width: unset;
  }
  
  .banking-table__footer {
    flex-direction: column;
    gap: 12px;
  }
}

@media (max-width: 599px) {
  .banking-table__header,
  .banking-table__controls,
  .banking-table__footer {
    padding: 12px 16px;
  }
  
  .banking-table__expanded-content {
    padding: 12px 16px;
  }
}

/* High contrast mode */
@media (prefers-contrast: high) {
  .banking-table {
    border: 2px solid rgb(var(--v-theme-outline));
  }
  
  .banking-table__header,
  .banking-table__controls,
  .banking-table__footer {
    border-bottom-width: 2px;
    border-top-width: 2px;
  }
}

/* Focus management */
.banking-table :deep(.v-data-table__tr--clickable:focus) {
  outline: 2px solid rgb(var(--v-theme-primary));
  outline-offset: -2px;
}

/* Banking variants */
.banking-table--audit {
  border-left: 4px solid rgb(var(--v-theme-info));
}

.banking-table--compliance {
  border-left: 4px solid rgb(var(--v-theme-success));
}

.banking-table--alert {
  border-left: 4px solid rgb(var(--v-theme-warning));
}

.banking-table--critical {
  border-left: 4px solid rgb(var(--v-theme-error));
}
</style>