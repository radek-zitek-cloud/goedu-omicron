/**
 * Banking Table Component Types
 *
 * Comprehensive type definitions for banking-specific data table components
 * with advanced features for financial data display and manipulation.
 */

/**
 * Banking table header configuration
 */
export interface BankingTableHeader {
  /**
   * Unique key for the column
   */
  key: string;

  /**
   * Display title for the column
   */
  title: string;

  /**
   * Column data type for specialized rendering
   */
  type?:
    | 'text'
    | 'number'
    | 'currency'
    | 'percentage'
    | 'date'
    | 'datetime'
    | 'status'
    | 'boolean'
    | 'actions'
    | 'audit-trail';

  /**
   * Column alignment
   */
  align?: 'start' | 'center' | 'end';

  /**
   * Column width
   */
  width?: string | number;

  /**
   * Minimum column width
   */
  minWidth?: string | number;

  /**
   * Maximum column width
   */
  maxWidth?: string | number;

  /**
   * Whether column is sortable
   */
  sortable?: boolean;

  /**
   * Whether column is filterable
   */
  filterable?: boolean;

  /**
   * Whether column is resizable
   */
  resizable?: boolean;

  /**
   * Whether column is hidden
   */
  hidden?: boolean;

  /**
   * Column format options
   */
  format?: BankingColumnFormat;

  /**
   * Custom CSS classes for the column
   */
  class?: string | string[];

  /**
   * Custom CSS classes for the header cell
   */
  headerClass?: string | string[];

  /**
   * Custom CSS classes for data cells
   */
  cellClass?: string | string[];

  /**
   * Tooltip text for the header
   */
  tooltip?: string;

  /**
   * Whether this column contains sensitive data
   */
  sensitive?: boolean;

  /**
   * Column validation rules
   */
  rules?: ((value: any) => boolean | string)[];
}

/**
 * Column formatting options
 */
export interface BankingColumnFormat {
  /**
   * Currency formatting options
   */
  currency?: {
    code: string;
    locale?: string;
    minimumFractionDigits?: number;
    maximumFractionDigits?: number;
  };

  /**
   * Percentage formatting options
   */
  percentage?: {
    minimumFractionDigits?: number;
    maximumFractionDigits?: number;
    showSymbol?: boolean;
  };

  /**
   * Date formatting options
   */
  date?: {
    format: string;
    locale?: string;
    timezone?: string;
  };

  /**
   * Number formatting options
   */
  number?: {
    minimumFractionDigits?: number;
    maximumFractionDigits?: number;
    useGrouping?: boolean;
    locale?: string;
  };

  /**
   * Status formatting options
   */
  status?: {
    colorMap?: Record<string, string>;
    iconMap?: Record<string, string>;
    textMap?: Record<string, string>;
  };

  /**
   * Boolean formatting options
   */
  boolean?: {
    trueText?: string;
    falseText?: string;
    trueIcon?: string;
    falseIcon?: string;
    trueColor?: string;
    falseColor?: string;
  };
}

/**
 * Banking table item interface
 */
export interface BankingTableItem {
  /**
   * Unique identifier for the row
   */
  id: string | number;

  /**
   * Row data - flexible to accommodate different data structures
   */
  [key: string]: any;

  /**
   * Optional row metadata
   */
  _meta?: {
    /**
     * Row variant for styling
     */
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info';

    /**
     * Whether row is selectable
     */
    selectable?: boolean;

    /**
     * Whether row is expandable
     */
    expandable?: boolean;

    /**
     * Custom CSS classes for the row
     */
    class?: string | string[];

    /**
     * Row tooltip
     */
    tooltip?: string;

    /**
     * Audit information
     */
    audit?: {
      createdAt: Date;
      createdBy: string;
      updatedAt?: Date;
      updatedBy?: string;
      version?: number;
    };
  };
}

/**
 * Table options for server-side operations
 */
export interface TableOptions {
  /**
   * Current page
   */
  page: number;

  /**
   * Items per page
   */
  itemsPerPage: number;

  /**
   * Sort configuration
   */
  sortBy: {
    key: string;
    order: 'asc' | 'desc';
  }[];

  /**
   * Search query
   */
  search?: string;

  /**
   * Filter configuration
   */
  filters?: Record<string, any>;
}

/**
 * Main component props
 */
export interface BankingTableProps {
  /**
   * Table headers configuration
   */
  headers: BankingTableHeader[];

  /**
   * Table data items
   */
  items?: BankingTableItem[];

  /**
   * Table loading state
   */
  loading?: boolean;

  /**
   * Table title
   */
  title?: string;

  /**
   * Table subtitle
   */
  subtitle?: string;

  /**
   * Table variant for styling
   */
  variant?: 'default' | 'audit' | 'compliance' | 'alert' | 'critical';

  /**
   * Items per page
   */
  itemsPerPage?: number;

  /**
   * Items per page options
   */
  itemsPerPageOptions?: number[];

  /**
   * Current page
   */
  page?: number;

  /**
   * Total items (for server-side pagination)
   */
  totalItems?: number;

  /**
   * Default sort configuration
   */
  sortBy?: { key: string; order: 'asc' | 'desc' }[];

  /**
   * Table density
   */
  density?: 'default' | 'comfortable' | 'compact';

  /**
   * Whether to show header
   */
  showHeader?: boolean;

  /**
   * Whether table is searchable
   */
  searchable?: boolean;

  /**
   * Whether table is filterable
   */
  filterable?: boolean;

  /**
   * Whether table is sortable
   */
  sortable?: boolean;

  /**
   * Whether rows are selectable
   */
  selectable?: boolean;

  /**
   * Whether rows are expandable
   */
  expandable?: boolean;

  /**
   * Whether to expand on click
   */
  expandOnClick?: boolean;

  /**
   * Whether table is exportable
   */
  exportable?: boolean;

  /**
   * Whether table is refreshable
   */
  refreshable?: boolean;

  /**
   * Whether sorting is required
   */
  mustSort?: boolean;

  /**
   * Whether multi-column sorting is enabled
   */
  multiSort?: boolean;

  /**
   * Mobile breakpoint
   */
  mobileBreakpoint?: number;

  /**
   * Whether to show striped rows
   */
  striped?: boolean;

  /**
   * Whether to show hover effects
   */
  hover?: boolean;

  /**
   * Pagination total visible pages
   */
  totalVisible?: number;

  /**
   * Server-side mode
   */
  serverMode?: boolean;

  /**
   * No data text
   */
  noDataText?: string;

  /**
   * Custom CSS classes
   */
  class?: string | string[] | Record<string, boolean>;
}

/**
 * Component events
 */
export interface BankingTableEmits {
  /**
   * Emitted when a row is clicked
   */
  (event: 'row-click', item: BankingTableItem, originalEvent: Event): void;

  /**
   * Emitted when table options change
   */
  (event: 'options-update', options: TableOptions): void;

  /**
   * Emitted when items per page changes
   */
  (event: 'items-per-page-update', itemsPerPage: number): void;

  /**
   * Emitted when sort configuration changes
   */
  (event: 'sort-update', sortBy: any[]): void;

  /**
   * Emitted when page changes
   */
  (event: 'page-update', page: number): void;

  /**
   * Emitted when search query changes
   */
  (event: 'search', query: string): void;

  /**
   * Emitted when export is requested
   */
  (event: 'export', items: BankingTableItem[]): void;

  /**
   * Emitted when refresh is requested
   */
  (event: 'refresh'): void;

  /**
   * Emitted when selection changes
   */
  (event: 'selection-change', selected: BankingTableItem[]): void;

  /**
   * Emitted when rows are expanded/collapsed
   */
  (event: 'expand-change', expanded: BankingTableItem[]): void;
}

/**
 * Table cell component props
 */
export interface BankingTableCellProps {
  /**
   * Cell value
   */
  value: any;

  /**
   * Cell type
   */
  type?: BankingTableHeader['type'];

  /**
   * Formatting options
   */
  format?: BankingColumnFormat;

  /**
   * Full row item
   */
  item?: BankingTableItem;

  /**
   * Header configuration
   */
  header?: BankingTableHeader;

  /**
   * Whether cell is editable
   */
  editable?: boolean;

  /**
   * Custom CSS classes
   */
  class?: string | string[];
}

/**
 * Export options
 */
export interface ExportOptions {
  /**
   * Export format
   */
  format: 'csv' | 'excel' | 'pdf';

  /**
   * Export filename
   */
  filename?: string;

  /**
   * Columns to include
   */
  columns?: string[];

  /**
   * Whether to include headers
   */
  includeHeaders?: boolean;

  /**
   * Date range filter
   */
  dateRange?: {
    start: Date;
    end: Date;
  };

  /**
   * Additional filters
   */
  filters?: Record<string, any>;
}
