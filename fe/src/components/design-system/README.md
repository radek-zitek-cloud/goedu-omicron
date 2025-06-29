# GoEdu Omicron Design System

A comprehensive design system and component library for banking applications, built with Vue 3 + Vuetify 3, providing consistent, accessible, and professional UI components specifically designed for financial services.

## 🏦 Banking-First Design

This design system is purpose-built for banking and financial applications, incorporating:

- **Banking UX Standards**: Components follow established patterns from successful banking platforms
- **Regulatory Compliance**: Built-in support for SOX, FFIEC, Basel III, and other financial regulations
- **Security by Design**: Secure input handling, audit trails, and permission-based interactions
- **Professional Aesthetics**: Clean, trustworthy design language that instills confidence

## ✨ Key Features

### 🎯 Comprehensive Component Library

#### Base Components
- **BankingButton**: Professional buttons with confirmation dialogs and action types
- **BankingInput**: Specialized inputs for currency, percentages, account numbers, and more
- **BankingTable**: Advanced data tables with banking-specific column types and formatting
- **BankingTableCell**: Smart cell component with automatic formatting for financial data

#### Banking-Specific Components
- **ApprovalWorkflow**: Multi-step approval processes with role-based assignments
- **AuditTrail**: Comprehensive audit logging with change tracking and risk assessment

### ♿ WCAG 2.1 AA Compliance

- **Keyboard Navigation**: Full keyboard support with logical tab order
- **Screen Reader Support**: Proper ARIA labels and semantic HTML structure
- **High Contrast Mode**: Components work seamlessly with high contrast themes
- **Focus Management**: Clear focus indicators and logical focus flow
- **Text Scaling**: Interface scales appropriately up to 200% text size

### 📱 Responsive Design

Components are designed to work across all device sizes:
- **Desktop**: Optimized for large screens with full feature sets
- **Tablet**: Adapted layouts for medium screens (iPad, Android tablets)
- **Mobile**: Touch-friendly interfaces with simplified layouts

### 🎨 Professional Theming

- **Light/Dark Mode**: Seamless switching between themes
- **Banking Color Palette**: Professional colors (banking gold, silver, bronze)
- **Status Colors**: Consistent compliance status indicators
- **Customizable**: Easy theme customization while maintaining consistency

## 🚀 Quick Start

### Installation

The design system is already included in the GoEdu project. To use components:

```vue
<template>
  <div>
    <!-- Banking Button with confirmation -->
    <BankingButton 
      variant="approve" 
      action-type="approve"
      :require-confirmation="true"
      confirmation-message="Are you sure you want to approve this transaction?"
      @click="handleApproval"
    >
      Approve Transaction
    </BankingButton>
    
    <!-- Currency Input with validation -->
    <BankingInput
      v-model="amount"
      type="currency"
      label="Transaction Amount"
      currency="USD"
      :rules="[required, minAmount]"
      required
    />
    
    <!-- Advanced Data Table -->
    <BankingTable
      :headers="transactionHeaders"
      :items="transactions"
      title="Transaction History"
      searchable
      exportable
      @row-click="viewTransaction"
    />
  </div>
</template>

<script setup lang="ts">
import { 
  BankingButton, 
  BankingInput, 
  BankingTable 
} from '@/components/design-system'

// Your component logic...
</script>
```

### Import Components

```typescript
// Import individual components
import { BankingButton, BankingInput } from '@/components/design-system'

// Import types
import type { 
  BankingButtonProps, 
  BankingInputProps, 
  ApprovalStep 
} from '@/components/design-system'
```

## 📖 Component Documentation

### Base Components

#### BankingButton

Professional button component with banking-specific variants and confirmation dialogs.

**Features:**
- Action types: approve, reject, submit, save-draft, send-reminder, export
- Confirmation dialogs for critical actions
- Loading states and disabled states with reasons
- Icon support and positioning
- Full keyboard navigation

**Example:**
```vue
<BankingButton 
  variant="approve" 
  action-type="approve"
  :require-confirmation="true"
  confirmation-message="Approve this transaction?"
  icon="mdi-check-circle"
  @click="handleApproval"
>
  Approve Transaction
</BankingButton>
```

#### BankingInput

Comprehensive input component with banking-specific formatting and validation.

**Supported Types:**
- `currency`: Automatic currency formatting with locale support
- `percentage`: Percentage inputs with validation
- `account-number`: Secure account number handling
- `routing-number`: Bank routing number with checksum validation
- `phone`: Automatic phone number formatting
- `ssn`: Secure SSN input with masking
- `email`: Email validation
- Standard HTML5 input types

**Example:**
```vue
<BankingInput
  v-model="transactionAmount"
  type="currency"
  label="Amount"
  currency="USD"
  :rules="[required, minAmount(100)]"
  :decimal-places="2"
  required
/>
```

#### BankingTable

Advanced data table with banking-specific column types and features.

**Features:**
- Currency, percentage, date, status, and audit trail columns
- Advanced search and filtering
- Export functionality (CSV, Excel, PDF)
- Responsive design with mobile optimization
- Row selection and expansion
- Server-side pagination support

**Example:**
```vue
<BankingTable
  :headers="[
    { key: 'id', title: 'Transaction ID', type: 'text' },
    { key: 'amount', title: 'Amount', type: 'currency' },
    { key: 'status', title: 'Status', type: 'status' },
    { key: 'date', title: 'Date', type: 'date' }
  ]"
  :items="transactions"
  title="Recent Transactions"
  searchable
  exportable
  striped
/>
```

### Banking Components

#### ApprovalWorkflow

Multi-step approval workflow component for banking processes.

**Features:**
- Step-by-step progress tracking
- Role-based assignments and permissions
- Due date tracking with overdue notifications
- Delegation and escalation support
- Comprehensive audit trail
- Confirmation dialogs for critical actions

**Example:**
```vue
<ApprovalWorkflow
  :steps="loanApprovalSteps"
  :current-step="2"
  workflow-status="in-progress"
  title="Loan Approval Process"
  @step-approved="handleStepApproved"
  @workflow-completed="handleWorkflowCompleted"
/>
```

#### AuditTrail

Comprehensive audit trail component with advanced filtering and export capabilities.

**Features:**
- Detailed change tracking
- Risk level assessment
- IP address and location tracking
- Advanced filtering by user, action, date, risk level
- Export functionality for compliance
- Real-time updates

**Example:**
```vue
<AuditTrail
  :entries="systemAuditEntries"
  title="System Audit Trail"
  filterable
  exportable
  @filters-changed="updateFilters"
  @export-requested="exportAuditLog"
/>
```

## 🎨 Design Tokens

The design system uses consistent design tokens defined in `src/plugins/vuetify.ts`:

### Colors

```typescript
// Banking-specific colors
'banking-gold': '#FFD700'      // Premium features
'banking-silver': '#C0C0C0'    // Standard features  
'banking-bronze': '#CD7F32'    // Basic features

// Compliance status colors
'compliant': colors.success[600]
'non-compliant': colors.error[600]
'pending-review': colors.warning[600]
'under-review': colors.info[600]
```

### Typography

- **Font Family**: Inter, system fonts for optimal readability
- **Font Scales**: Consistent sizing from 0.75rem to 2.5rem
- **Font Weights**: 400 (normal), 500 (medium), 600 (semibold), 700 (bold)

### Spacing

- **Base Unit**: 4px
- **Scale**: 4px, 8px, 12px, 16px, 24px, 32px, 48px, 64px

## 🛠️ Development

### Testing Components

Visit `/design-system` in the application to see the interactive component showcase with:
- All component variants and states
- Live examples with editable props
- Accessibility testing tools
- Code examples and documentation

### Adding New Components

1. Create component in appropriate directory (`base/` or `banking/`)
2. Create corresponding `.types.ts` file with TypeScript interfaces
3. Add comprehensive documentation and examples
4. Include accessibility features (ARIA labels, keyboard navigation)
5. Add to showcase for testing
6. Export from `index.ts`

### Accessibility Testing

Components should be tested with:
- **Keyboard navigation**: Tab through all interactive elements
- **Screen readers**: Test with NVDA, JAWS, or VoiceOver
- **High contrast mode**: Ensure visibility with high contrast themes
- **Color blindness**: Test with color blindness simulators
- **Mobile accessibility**: Test touch targets and mobile screen readers

## 📋 Compliance Standards

### Banking Regulations
- **SOX Section 404**: Internal control assessment support
- **FFIEC Guidelines**: IT risk management compliance
- **Basel III**: Operational risk data aggregation
- **GDPR Article 32**: Security of processing requirements
- **PCI DSS**: Payment card industry data security

### Accessibility Standards
- **WCAG 2.1 AA**: Full compliance with web accessibility guidelines
- **Section 508**: US federal accessibility requirements
- **EN 301 549**: European accessibility standard

## 🏗️ Architecture

### Component Structure
```
design-system/
├── base/                   # Core UI components
│   ├── BankingButton.vue
│   ├── BankingButton.types.ts
│   ├── BankingInput.vue
│   ├── BankingInput.types.ts
│   ├── BankingTable.vue
│   ├── BankingTable.types.ts
│   └── BankingTableCell.vue
├── banking/                # Banking-specific components
│   ├── ApprovalWorkflow.vue
│   ├── ApprovalWorkflow.types.ts
│   ├── AuditTrail.vue
│   └── AuditTrail.types.ts
├── layout/                 # Layout components (future)
├── showcase/               # Component showcase
│   └── DesignSystemShowcase.vue
├── index.ts               # Main export file
└── README.md             # This documentation
```

### Built With
- **Vue 3**: Modern reactive framework
- **Vuetify 3**: Material Design component library
- **TypeScript**: Type-safe development
- **Vite**: Fast build tool
- **Material Design Icons**: Comprehensive icon set

## 🤝 Contributing

When contributing to the design system:

1. Follow existing patterns and conventions
2. Include comprehensive TypeScript types
3. Add accessibility features (ARIA, keyboard navigation)
4. Write documentation and examples
5. Test across browsers and devices
6. Ensure WCAG 2.1 AA compliance

## 📄 License

This design system is proprietary to GoEdu and is intended for use within GoEdu banking applications only.

---

Built with ❤️ for the banking industry by the GoEdu team.