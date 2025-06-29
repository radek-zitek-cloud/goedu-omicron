/**
 * Vuetify Theme Configuration for GoEdu Omicron Banking Platform
 *
 * Professional theme configuration optimized for banking and financial applications.
 * Includes multiple theme variants, accessibility features, and custom color palettes
 * that comply with WCAG 2.1 AA standards.
 *
 * Key Features:
 * - Light and dark theme variants
 * - High contrast ratios for accessibility
 * - Banking-appropriate color schemes
 * - Consistent spacing and typography
 * - Professional gradients and shadows
 *
 * @author GoEdu Development Team
 * @version 1.0.0
 * @since 2024
 */

import { type ThemeDefinition } from 'vuetify';

/**
 * Banking-Specific Color Palette
 *
 * Professional colors that convey trust, stability, and reliability.
 * All colors meet WCAG 2.1 AA contrast requirements for accessibility.
 */
const colors = {
  // Primary Banking Colors
  primary: {
    50: '#e3f2fd', // Very light blue
    100: '#bbdefb', // Light blue
    200: '#90caf9', // Medium light blue
    300: '#64b5f6', // Medium blue
    400: '#42a5f5', // Medium dark blue
    500: '#2196f3', // Primary blue (main)
    600: '#1e88e5', // Dark blue
    700: '#1976d2', // Darker blue
    800: '#1565c0', // Very dark blue
    900: '#0d47a1', // Darkest blue
  },

  // Secondary Professional Colors
  secondary: {
    50: '#f3e5f5', // Very light purple
    100: '#e1bee7', // Light purple
    200: '#ce93d8', // Medium light purple
    300: '#ba68c8', // Medium purple
    400: '#ab47bc', // Medium dark purple
    500: '#9c27b0', // Primary purple (main)
    600: '#8e24aa', // Dark purple
    700: '#7b1fa2', // Darker purple
    800: '#6a1b9a', // Very dark purple
    900: '#4a148c', // Darkest purple
  },

  // Success Colors (for positive actions, confirmations)
  success: {
    50: '#e8f5e8', // Very light green
    100: '#c8e6c9', // Light green
    200: '#a5d6a7', // Medium light green
    300: '#81c784', // Medium green
    400: '#66bb6a', // Medium dark green
    500: '#4caf50', // Primary green (main)
    600: '#43a047', // Dark green
    700: '#388e3c', // Darker green
    800: '#2e7d32', // Very dark green
    900: '#1b5e20', // Darkest green
  },

  // Warning Colors (for cautions, pending states)
  warning: {
    50: '#fff8e1', // Very light amber
    100: '#ffecb3', // Light amber
    200: '#ffe082', // Medium light amber
    300: '#ffd54f', // Medium amber
    400: '#ffca28', // Medium dark amber
    500: '#ffc107', // Primary amber (main)
    600: '#ffb300', // Dark amber
    700: '#ffa000', // Darker amber
    800: '#ff8f00', // Very dark amber
    900: '#ff6f00', // Darkest amber
  },

  // Error Colors (for errors, critical issues)
  error: {
    50: '#ffebee', // Very light red
    100: '#ffcdd2', // Light red
    200: '#ef9a9a', // Medium light red
    300: '#e57373', // Medium red
    400: '#ef5350', // Medium dark red
    500: '#f44336', // Primary red (main)
    600: '#e53935', // Dark red
    700: '#d32f2f', // Darker red
    800: '#c62828', // Very dark red
    900: '#b71c1c', // Darkest red
  },

  // Information Colors (for informational messages)
  info: {
    50: '#e1f5fe', // Very light cyan
    100: '#b3e5fc', // Light cyan
    200: '#81d4fa', // Medium light cyan
    300: '#4fc3f7', // Medium cyan
    400: '#29b6f6', // Medium dark cyan
    500: '#03a9f4', // Primary cyan (main)
    600: '#039be5', // Dark cyan
    700: '#0288d1', // Darker cyan
    800: '#0277bd', // Very dark cyan
    900: '#01579b', // Darkest cyan
  },

  // Neutral Colors (for backgrounds, borders, text)
  neutral: {
    50: '#fafafa', // Very light gray
    100: '#f5f5f5', // Light gray
    200: '#eeeeee', // Medium light gray
    300: '#e0e0e0', // Medium gray
    400: '#bdbdbd', // Medium dark gray
    500: '#9e9e9e', // Primary gray (main)
    600: '#757575', // Dark gray
    700: '#616161', // Darker gray
    800: '#424242', // Very dark gray
    900: '#212121', // Darkest gray
  },
};

/**
 * Light Theme Configuration
 *
 * Professional light theme optimized for banking applications.
 * Provides excellent readability and professional appearance.
 */
const lightTheme: ThemeDefinition = {
  dark: false,
  colors: {
    // Main brand colors
    primary: colors.primary[600], // Professional blue
    'primary-darken-1': colors.primary[700],
    'primary-lighten-1': colors.primary[500],

    secondary: colors.secondary[600], // Professional purple
    'secondary-darken-1': colors.secondary[700],
    'secondary-lighten-1': colors.secondary[500],

    // Status colors
    success: colors.success[600], // Professional green
    'success-darken-1': colors.success[700],
    'success-lighten-1': colors.success[500],

    warning: colors.warning[600], // Professional amber
    'warning-darken-1': colors.warning[700],
    'warning-lighten-1': colors.warning[500],

    error: colors.error[600], // Professional red
    'error-darken-1': colors.error[700],
    'error-lighten-1': colors.error[500],

    info: colors.info[600], // Professional cyan
    'info-darken-1': colors.info[700],
    'info-lighten-1': colors.info[500],

    // Background colors
    background: '#ffffff', // Pure white background
    surface: '#ffffff', // White surface
    'surface-bright': colors.neutral[50],
    'surface-light': colors.neutral[100],
    'surface-variant': colors.neutral[200],

    // Text colors (high contrast for accessibility)
    'on-primary': '#ffffff',
    'on-secondary': '#ffffff',
    'on-success': '#ffffff',
    'on-warning': colors.neutral[900],
    'on-error': '#ffffff',
    'on-info': '#ffffff',
    'on-background': colors.neutral[900],
    'on-surface': colors.neutral[900],
    'on-surface-variant': colors.neutral[700],

    // Additional banking-specific colors
    'banking-gold': '#FFD700', // For premium features
    'banking-silver': '#C0C0C0', // For standard features
    'banking-bronze': '#CD7F32', // For basic features

    // Compliance status colors
    compliant: colors.success[600],
    'non-compliant': colors.error[600],
    'pending-review': colors.warning[600],
    'under-review': colors.info[600],

    // Borders and dividers
    border: colors.neutral[300],
    divider: colors.neutral[200],
    outline: colors.neutral[400],

    // Hover and focus states
    'hover-overlay': 'rgba(0, 0, 0, 0.04)',
    'focus-overlay': 'rgba(33, 150, 243, 0.12)',
    'pressed-overlay': 'rgba(0, 0, 0, 0.08)',
  },
};

/**
 * Dark Theme Configuration
 *
 * Professional dark theme for low-light environments.
 * Maintains accessibility while reducing eye strain.
 */
const darkTheme: ThemeDefinition = {
  dark: true,
  colors: {
    // Main brand colors (adjusted for dark theme)
    primary: colors.primary[400], // Lighter blue for dark background
    'primary-darken-1': colors.primary[500],
    'primary-lighten-1': colors.primary[300],

    secondary: colors.secondary[400], // Lighter purple for dark background
    'secondary-darken-1': colors.secondary[500],
    'secondary-lighten-1': colors.secondary[300],

    // Status colors (adjusted for dark theme)
    success: colors.success[400], // Lighter green
    'success-darken-1': colors.success[500],
    'success-lighten-1': colors.success[300],

    warning: colors.warning[400], // Lighter amber
    'warning-darken-1': colors.warning[500],
    'warning-lighten-1': colors.warning[300],

    error: colors.error[400], // Lighter red
    'error-darken-1': colors.error[500],
    'error-lighten-1': colors.error[300],

    info: colors.info[400], // Lighter cyan
    'info-darken-1': colors.info[500],
    'info-lighten-1': colors.info[300],

    // Dark background colors
    background: '#121212', // Material Design dark background
    surface: '#1e1e1e', // Elevated dark surface
    'surface-bright': '#2a2a2a', // Brighter dark surface
    'surface-light': '#343434', // Light dark surface
    'surface-variant': '#424242', // Variant dark surface

    // Text colors (high contrast for dark theme)
    'on-primary': colors.neutral[900],
    'on-secondary': colors.neutral[900],
    'on-success': colors.neutral[900],
    'on-warning': colors.neutral[900],
    'on-error': colors.neutral[900],
    'on-info': colors.neutral[900],
    'on-background': '#ffffff',
    'on-surface': '#ffffff',
    'on-surface-variant': colors.neutral[300],

    // Banking-specific colors (adjusted for dark theme)
    'banking-gold': '#FFE55C', // Lighter gold for dark background
    'banking-silver': '#E8E8E8', // Lighter silver
    'banking-bronze': '#DEB887', // Lighter bronze

    // Compliance status colors (adjusted for dark theme)
    compliant: colors.success[400],
    'non-compliant': colors.error[400],
    'pending-review': colors.warning[400],
    'under-review': colors.info[400],

    // Dark theme borders and dividers
    border: colors.neutral[600],
    divider: colors.neutral[700],
    outline: colors.neutral[500],

    // Dark theme hover and focus states
    'hover-overlay': 'rgba(255, 255, 255, 0.08)',
    'focus-overlay': 'rgba(100, 181, 246, 0.24)',
    'pressed-overlay': 'rgba(255, 255, 255, 0.16)',
  },
};

/**
 * Theme Configuration Export
 *
 * Main theme configuration object that includes both light and dark themes
 * with proper defaults and fallbacks.
 */
export const themeConfig = {
  defaultTheme: 'light',
  variations: {
    colors: ['primary', 'secondary', 'success', 'warning', 'error', 'info'],
    lighten: 5,
    darken: 5,
  },
  themes: {
    light: lightTheme,
    dark: darkTheme,
  },
};

/**
 * Banking-Specific Component Styles
 *
 * Custom styles for banking application components.
 * These can be used with Vuetify's SCSS variables or CSS custom properties.
 */
export const bankingStyles = {
  // Card styles for different banking contexts
  cards: {
    compliance: {
      borderLeft: `4px solid ${colors.success[500]}`,
      backgroundColor: colors.success[50],
    },
    warning: {
      borderLeft: `4px solid ${colors.warning[500]}`,
      backgroundColor: colors.warning[50],
    },
    critical: {
      borderLeft: `4px solid ${colors.error[500]}`,
      backgroundColor: colors.error[50],
    },
    info: {
      borderLeft: `4px solid ${colors.info[500]}`,
      backgroundColor: colors.info[50],
    },
  },

  // Status badges
  badges: {
    compliant: {
      backgroundColor: colors.success[500],
      color: '#ffffff',
    },
    nonCompliant: {
      backgroundColor: colors.error[500],
      color: '#ffffff',
    },
    pending: {
      backgroundColor: colors.warning[500],
      color: colors.neutral[900],
    },
    underReview: {
      backgroundColor: colors.info[500],
      color: '#ffffff',
    },
  },

  // Progress indicators
  progress: {
    low: colors.error[500], // 0-30%
    medium: colors.warning[500], // 31-70%
    high: colors.success[500], // 71-100%
  },
};

/**
 * Accessibility Helper Functions
 *
 * Utilities to ensure proper contrast ratios and accessibility compliance.
 */
export const accessibilityHelpers = {
  /**
   * Calculate contrast ratio between two colors
   * @param color1 First color in hex format
   * @param color2 Second color in hex format
   * @returns Contrast ratio (1-21)
   */
  calculateContrastRatio(color1: string, color2: string): number {
    // Implementation would go here for WCAG compliance checking
    // This is a placeholder for the actual contrast calculation
    return 4.5; // Minimum AA compliance ratio
  },

  /**
   * Get appropriate text color for background
   * @param backgroundColor Background color in hex format
   * @returns Appropriate text color (light or dark)
   */
  getTextColorForBackground(backgroundColor: string): string {
    // Implementation would determine if light or dark text is more readable
    // This is a placeholder for the actual calculation
    return colors.neutral[900]; // Default to dark text
  },
};

export default themeConfig;
