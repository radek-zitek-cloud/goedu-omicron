/**
 * Audit Trail Component Types
 * 
 * Comprehensive type definitions for banking audit trail components
 * with detailed tracking, filtering, and compliance features.
 */

/**
 * Audit entry interface
 */
export interface AuditEntry {
  /**
   * Unique identifier for the audit entry
   */
  id: string
  
  /**
   * Action that was performed
   */
  action: string
  
  /**
   * Optional description of the action
   */
  description?: string
  
  /**
   * User who performed the action
   */
  user: string
  
  /**
   * Timestamp when the action was performed
   */
  timestamp: Date | string
  
  /**
   * Target of the action (e.g., document ID, user ID)
   */
  target?: string
  
  /**
   * Category of the action
   */
  category?: 'security' | 'financial' | 'compliance' | 'user' | 'system' | 'data'
  
  /**
   * Risk level assessment
   */
  riskLevel?: 'low' | 'medium' | 'high'
  
  /**
   * Additional details about the action
   */
  details?: Record<string, any>
  
  /**
   * Changes made (for update actions)
   */
  changes?: AuditChange[]
  
  /**
   * Additional metadata
   */
  metadata?: {
    /**
     * IP address from which the action was performed
     */
    ipAddress?: string
    
    /**
     * Geographic location
     */
    location?: string
    
    /**
     * User agent string
     */
    userAgent?: string
    
    /**
     * Session ID
     */
    sessionId?: string
    
    /**
     * Request ID for tracing
     */
    requestId?: string
    
    /**
     * Additional custom metadata
     */
    [key: string]: any
  }
}

/**
 * Audit change interface for tracking field modifications
 */
export interface AuditChange {
  /**
   * Field that was changed
   */
  field: string
  
  /**
   * Previous value
   */
  oldValue: any
  
  /**
   * New value
   */
  newValue: any
  
  /**
   * Data type of the field
   */
  dataType?: string
}

/**
 * Audit filters interface
 */
export interface AuditFilters {
  /**
   * Filter by user
   */
  user: string
  
  /**
   * Filter by action
   */
  action: string
  
  /**
   * Filter by date range - from
   */
  dateFrom: string
  
  /**
   * Filter by date range - to
   */
  dateTo: string
  
  /**
   * Filter by risk level
   */
  riskLevel: string
  
  /**
   * Filter by IP address
   */
  ipAddress: string
  
  /**
   * Filter by category
   */
  category?: string
  
  /**
   * Filter by target
   */
  target?: string
}

/**
 * Export request interface
 */
export interface AuditExportRequest {
  /**
   * Entries to export
   */
  entries: AuditEntry[]
  
  /**
   * Filters applied
   */
  filters: AuditFilters
  
  /**
   * Export timestamp
   */
  timestamp: Date
  
  /**
   * Export format
   */
  format?: 'csv' | 'excel' | 'pdf' | 'json'
  
  /**
   * Include metadata in export
   */
  includeMetadata?: boolean
}

/**
 * Main component props
 */
export interface AuditTrailProps {
  /**
   * Array of audit entries
   */
  entries?: AuditEntry[]
  
  /**
   * Component title
   */
  title?: string
  
  /**
   * Component subtitle
   */
  subtitle?: string
  
  /**
   * Component variant for styling
   */
  variant?: 'default' | 'compact' | 'detailed'
  
  /**
   * Whether to show the header
   */
  showHeader?: boolean
  
  /**
   * Whether to show summary information
   */
  showSummary?: boolean
  
  /**
   * Whether to show detailed entry information
   */
  showDetails?: boolean
  
  /**
   * Whether filtering is enabled
   */
  filterable?: boolean
  
  /**
   * Whether export functionality is enabled
   */
  exportable?: boolean
  
  /**
   * Whether refresh functionality is enabled
   */
  refreshable?: boolean
  
  /**
   * Number of items per page
   */
  itemsPerPage?: number
  
  /**
   * Whether entries are grouped by date
   */
  groupByDate?: boolean
  
  /**
   * Loading state
   */
  loading?: boolean
  
  /**
   * No data text
   */
  noDataText?: string
  
  /**
   * Custom CSS classes
   */
  class?: string | string[] | Record<string, boolean>
}

/**
 * Component events
 */
export interface AuditTrailEmits {
  /**
   * Emitted when an audit entry is selected
   */
  (event: 'entry-selected', entry: AuditEntry): void
  
  /**
   * Emitted when filters change
   */
  (event: 'filters-changed', filters: AuditFilters): void
  
  /**
   * Emitted when filters are cleared
   */
  (event: 'filters-cleared'): void
  
  /**
   * Emitted when page changes
   */
  (event: 'page-changed', page: number): void
  
  /**
   * Emitted when export is requested
   */
  (event: 'export-requested', request: AuditExportRequest): void
  
  /**
   * Emitted when refresh is requested
   */
  (event: 'refresh-requested'): void
  
  /**
   * Emitted when entry details are expanded
   */
  (event: 'entry-expanded', entry: AuditEntry): void
  
  /**
   * Emitted when entry details are collapsed
   */
  (event: 'entry-collapsed', entry: AuditEntry): void
}

/**
 * Audit trail configuration
 */
export interface AuditTrailConfig {
  /**
   * Default items per page
   */
  defaultItemsPerPage?: number
  
  /**
   * Maximum entries to display
   */
  maxEntries?: number
  
  /**
   * Default date range filter (days)
   */
  defaultDateRange?: number
  
  /**
   * Auto-refresh interval (milliseconds)
   */
  autoRefreshInterval?: number
  
  /**
   * Retention policy (days)
   */
  retentionDays?: number
  
  /**
   * High-risk action patterns
   */
  highRiskPatterns?: string[]
  
  /**
   * Medium-risk action patterns
   */
  mediumRiskPatterns?: string[]
  
  /**
   * Sensitive field patterns
   */
  sensitiveFields?: string[]
}

/**
 * Audit trail analytics
 */
export interface AuditTrailAnalytics {
  /**
   * Total entries in the time period
   */
  totalEntries: number
  
  /**
   * Entries by risk level
   */
  entriesByRisk: {
    low: number
    medium: number
    high: number
  }
  
  /**
   * Entries by category
   */
  entriesByCategory: Record<string, number>
  
  /**
   * Top users by activity
   */
  topUsers: Array<{
    user: string
    count: number
    percentage: number
  }>
  
  /**
   * Top actions performed
   */
  topActions: Array<{
    action: string
    count: number
    percentage: number
  }>
  
  /**
   * Activity timeline (hourly)
   */
  activityTimeline: Array<{
    hour: number
    count: number
  }>
  
  /**
   * Geographic distribution
   */
  geographicDistribution: Array<{
    location: string
    count: number
    percentage: number
  }>
}