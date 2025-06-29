/**
 * Main application entry point for GoEdu Omicron Banking Control Testing Platform
 *
 * This file initializes the Vue.js application with all necessary plugins and configurations
 * for a production-ready banking compliance and control testing system.
 *
 * Key Features:
 * - Vue 3 with Composition API
 * - Vuetify 3 Material Design components
 * - Pinia state management
 * - Vue Router for navigation
 * - Apollo GraphQL client
 * - VueUse utilities
 *
 * @author GoEdu Development Team
 * @version 1.0.0
 * @since 2024
 */

import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { createPinia } from 'pinia';
import { createVuetify } from 'vuetify';
import { DefaultApolloClient } from '@vue/apollo-composable';
import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';

// Import Vuetify styles and icons
import 'vuetify/styles';
import '@mdi/font/css/materialdesignicons.css';
import 'material-design-icons-iconfont/dist/material-design-icons.css';

// Import main components and styles
import App from './App.vue';
import './assets/main.css';
import router from './router/index';

// Import Vuetify theme configuration
// import { themeConfig } from './plugins/vuetify' // Temporarily disabled for initial setup

/**
 * Apollo GraphQL Client Configuration
 *
 * Configures the GraphQL client for API communication with the backend.
 * Includes authentication headers and caching strategies.
 *
 * Environment Variables Required:
 * - VITE_GRAPHQL_URI: Backend GraphQL endpoint URL
 * - VITE_API_BASE_URL: Base API URL for REST endpoints
 */
const httpLink = createHttpLink({
  uri: import.meta.env.VITE_GRAPHQL_URI || 'http://localhost:4000/graphql',
  credentials: 'include', // Include cookies for authentication
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
});

const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache({
    typePolicies: {
      // Configure caching policies for banking data
      Organization: {
        fields: {
          controls: {
            merge(existing = [], incoming) {
              return incoming;
            },
          },
        },
      },
      Control: {
        fields: {
          testResults: {
            merge(existing = [], incoming) {
              return incoming;
            },
          },
        },
      },
    },
  }),
  defaultOptions: {
    watchQuery: {
      errorPolicy: 'all',
    },
    query: {
      errorPolicy: 'all',
    },
  },
});

/**
 * Vuetify Material Design Configuration
 *
 * Creates Vuetify instance with custom theme optimized for banking applications.
 * Includes accessibility features and professional color schemes.
 */
const vuetify = createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        dark: false,
        colors: {
          primary: '#1976d2',
          secondary: '#424242',
          accent: '#82b1ff',
          error: '#ff5252',
          info: '#2196f3',
          success: '#4caf50',
          warning: '#ffc107',
        },
      },
    },
  },
  defaults: {
    // Default component props for consistency
    VBtn: {
      variant: 'elevated',
      color: 'primary',
    },
    VCard: {
      variant: 'elevated',
      elevation: 2,
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
    },
  },
});

/**
 * Pinia State Management Store
 *
 * Global state management for user authentication, organization data,
 * and application-wide settings.
 */
const pinia = createPinia();

/**
 * Application Initialization
 *
 * Creates and configures the Vue application instance with all plugins
 * and global configurations required for the banking control platform.
 *
 * Order of plugin registration is important for proper dependency injection.
 */
async function initializeApp() {
  try {
    const app = createApp(App);

    // Register plugins in order of dependency
    app.use(pinia); // State management first
    app.use(router); // Router second
    app.use(vuetify); // UI framework third

    // Provide Apollo client globally
    app.provide(DefaultApolloClient, apolloClient);

    // Global error handler for production monitoring
    app.config.errorHandler = (error, instance, info) => {
      console.error('Global error caught:', error, info);

      // In production, send to monitoring service
      if (import.meta.env.PROD) {
        // TODO: Integrate with error monitoring service (Sentry, etc.)
        console.error('Production error:', { error, info });
      }
    };

    // Global properties for debugging (development only)
    if (import.meta.env.DEV) {
      app.config.globalProperties.$log = console.log;
      app.config.globalProperties.$apolloClient = apolloClient;
    }

    // Mount application to DOM
    app.mount('#app');

    console.log('✅ GoEdu Omicron Banking Platform initialized successfully');
  } catch (error) {
    console.error('❌ Failed to initialize application:', error);

    // Show user-friendly error message
    document.getElementById('app')!.innerHTML = `
      <div style="padding: 2rem; text-align: center; font-family: Arial, sans-serif;">
        <h2 style="color: #d32f2f;">Application Initialization Failed</h2>
        <p>Please refresh the page or contact support if the issue persists.</p>
        <pre style="background: #f5f5f5; padding: 1rem; border-radius: 4px; text-align: left; margin-top: 1rem;">
          ${error}
        </pre>
      </div>
    `;
  }
}

// Initialize the application
initializeApp();
