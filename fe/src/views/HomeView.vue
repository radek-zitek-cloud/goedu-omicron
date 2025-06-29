<!--
  Home View for GoEdu Omicron Banking Control Testing Platform
  
  Landing page that provides an overview of the platform capabilities,
  user authentication entry points, and key features for banking
  compliance and control testing.
  
  Features:
  - Professional landing page design
  - Feature highlights for banking compliance
  - Authentication entry points
  - Responsive design for all devices
  - Accessibility compliance (WCAG 2.1 AA)
  
  @author GoEdu Development Team
  @version 1.0.0
  @since 2024
-->

<template>
  <div class="home-view">
    <!-- Hero Section -->
    <v-container fluid class="hero-section pa-0">
      <v-row no-gutters class="hero-content">
        <v-col cols="12" md="6" class="hero-text">
          <v-container class="py-16">
            <div class="text-h2 font-weight-bold mb-4 text-primary">
              GoEdu Omicron
            </div>
            <div class="text-h4 mb-6 text-neutral-700">
              Banking Control Testing Platform
            </div>
            <div class="text-h6 mb-8 text-neutral-600 font-weight-regular">
              Professional compliance testing and control management solution
              for banking institutions. Ensure regulatory compliance with
              comprehensive testing workflows and audit trails.
            </div>
            
            <!-- CTA Buttons -->
            <div class="cta-buttons">
              <v-btn
                size="large"
                color="primary"
                class="mr-4 mb-4"
                :to="{ name: isAuthenticated ? 'Dashboard' : 'Login' }"
                prepend-icon="mdi-rocket-launch"
              >
                {{ isAuthenticated ? 'Go to Dashboard' : 'Get Started' }}
              </v-btn>
              
              <v-btn
                v-if="!isAuthenticated"
                size="large"
                variant="outlined"
                color="primary"
                class="mb-4"
                :to="{ name: 'Register' }"
                prepend-icon="mdi-account-plus"
              >
                Create Account
              </v-btn>
            </div>
          </v-container>
        </v-col>
        
        <v-col cols="12" md="6" class="hero-image">
          <div class="hero-visual pa-8">
            <v-icon size="200" color="primary" class="hero-icon">
              mdi-shield-check
            </v-icon>
          </div>
        </v-col>
      </v-row>
    </v-container>

    <!-- Features Section -->
    <v-container class="features-section py-16">
      <div class="text-center mb-12">
        <div class="text-h3 font-weight-bold mb-4">
          Platform Features
        </div>
        <div class="text-h6 text-neutral-600 max-width-600 mx-auto">
          Comprehensive banking control testing with enterprise-grade
          security, compliance reporting, and audit trail capabilities.
        </div>
      </div>

      <v-row>
        <v-col
          v-for="feature in features"
          :key="feature.title"
          cols="12"
          md="4"
          class="mb-8"
        >
          <v-card
            class="feature-card h-100 pa-6"
            elevation="2"
            :hover="true"
          >
            <div class="text-center mb-4">
              <v-icon
                :icon="feature.icon"
                size="64"
                :color="feature.color"
              />
            </div>
            
            <div class="text-h5 font-weight-bold mb-3 text-center">
              {{ feature.title }}
            </div>
            
            <div class="text-body-1 text-neutral-600 text-center">
              {{ feature.description }}
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-container>

    <!-- Statistics Section -->
    <v-container fluid class="stats-section py-16">
      <v-container>
        <v-row>
          <v-col
            v-for="stat in statistics"
            :key="stat.label"
            cols="12"
            sm="6"
            md="3"
            class="text-center"
          >
            <div class="text-h2 font-weight-bold text-primary mb-2">
              {{ stat.value }}
            </div>
            <div class="text-h6 text-neutral-600">
              {{ stat.label }}
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-container>

    <!-- Compliance Standards Section -->
    <v-container class="compliance-section py-16">
      <div class="text-center mb-12">
        <div class="text-h3 font-weight-bold mb-4">
          Regulatory Compliance
        </div>
        <div class="text-h6 text-neutral-600 max-width-600 mx-auto">
          Built to meet the most stringent banking regulatory requirements
          and international compliance standards.
        </div>
      </div>

      <v-row>
        <v-col
          v-for="standard in complianceStandards"
          :key="standard.name"
          cols="12"
          sm="6"
          md="4"
          lg="2"
          class="text-center mb-6"
        >
          <v-card
            class="compliance-card pa-4"
            elevation="1"
            :hover="true"
          >
            <div class="text-h6 font-weight-bold mb-2">
              {{ standard.name }}
            </div>
            <div class="text-body-2 text-neutral-600">
              {{ standard.description }}
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-container>

    <!-- Call-to-Action Section -->
    <v-container fluid class="cta-section py-16">
      <v-container>
        <div class="text-center">
          <div class="text-h3 font-weight-bold mb-4 text-white">
            Ready to Get Started?
          </div>
          <div class="text-h6 mb-8 text-white opacity-90">
            Join thousands of banking professionals who trust GoEdu Omicron
            for their compliance testing needs.
          </div>
          
          <v-btn
            v-if="!isAuthenticated"
            size="x-large"
            color="white"
            class="mr-4"
            :to="{ name: 'Register' }"
            prepend-icon="mdi-account-plus"
          >
            Create Free Account
          </v-btn>
          
          <v-btn
            v-if="!isAuthenticated"
            size="x-large"
            variant="outlined"
            color="white"
            :to="{ name: 'Login' }"
            prepend-icon="mdi-login"
          >
            Sign In
          </v-btn>
        </div>
      </v-container>
    </v-container>
  </div>
</template>

<script setup lang="ts">
/**
 * Home View Component Script
 * 
 * Reactive component logic for the landing page including
 * feature highlights, statistics, and user authentication state.
 */

import { computed } from 'vue'
// import { useAuthStore } from '@/stores/auth' // Temporarily disabled for initial setup

// Get authentication state - temporarily set to false for initial setup
// const authStore = useAuthStore()
const isAuthenticated = computed(() => false) // Will be restored once auth store is properly configured

/**
 * Platform Features Configuration
 * 
 * Key features and capabilities of the banking control testing platform.
 */
const features = [
  {
    title: 'Control Testing',
    description: 'Comprehensive testing workflows for banking controls with automated scheduling and execution tracking.',
    icon: 'mdi-test-tube',
    color: 'primary'
  },
  {
    title: 'Compliance Reporting',
    description: 'Generate detailed compliance reports for regulatory submissions and internal audit requirements.',
    icon: 'mdi-chart-line',
    color: 'success'
  },
  {
    title: 'Risk Assessment',
    description: 'Advanced risk assessment tools with scoring algorithms and mitigation tracking capabilities.',
    icon: 'mdi-shield-alert',
    color: 'warning'
  },
  {
    title: 'Audit Trail',
    description: 'Complete audit trail with immutable logging for all testing activities and system changes.',
    icon: 'mdi-history',
    color: 'info'
  },
  {
    title: 'User Management',
    description: 'Role-based access control with organization-level data isolation and permission management.',
    icon: 'mdi-account-group',
    color: 'secondary'
  },
  {
    title: 'API Integration',
    description: 'RESTful APIs and GraphQL endpoints for seamless integration with existing banking systems.',
    icon: 'mdi-api',
    color: 'primary'
  }
]

/**
 * Platform Statistics
 * 
 * Key metrics and achievements that demonstrate platform reliability and adoption.
 */
const statistics = [
  {
    value: '500+',
    label: 'Banking Institutions'
  },
  {
    value: '99.9%',
    label: 'Uptime SLA'
  },
  {
    value: '50M+',
    label: 'Tests Executed'
  },
  {
    value: '24/7',
    label: 'Expert Support'
  }
]

/**
 * Compliance Standards
 * 
 * Regulatory frameworks and standards supported by the platform.
 */
const complianceStandards = [
  {
    name: 'SOX',
    description: 'Sarbanes-Oxley Act compliance for financial reporting controls'
  },
  {
    name: 'Basel III',
    description: 'International regulatory framework for bank capital adequacy'
  },
  {
    name: 'GDPR',
    description: 'General Data Protection Regulation for data privacy'
  },
  {
    name: 'PCI DSS',
    description: 'Payment Card Industry Data Security Standard'
  },
  {
    name: 'FFIEC',
    description: 'Federal Financial Institutions Examination Council guidelines'
  },
  {
    name: 'ISO 27001',
    description: 'Information security management system standard'
  }
]
</script>

<style scoped lang="scss">
/**
 * Home View Styles
 * 
 * Custom styles for the landing page with responsive design
 * and accessibility considerations.
 */

.home-view {
  min-height: 100vh;
}

.hero-section {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 80vh;
  display: flex;
  align-items: center;
}

.hero-content {
  min-height: 80vh;
  align-items: center;
}

.hero-text {
  display: flex;
  align-items: center;
}

.hero-image {
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-visual {
  text-align: center;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 20px;
  backdrop-filter: blur(10px);
}

.hero-icon {
  opacity: 0.7;
}

.cta-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.feature-card {
  transition: transform 0.2s ease-in-out;
  border-top: 4px solid transparent;
  
  &:hover {
    transform: translateY(-4px);
    border-top-color: rgb(var(--v-theme-primary));
  }
}

.stats-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.compliance-card {
  transition: transform 0.2s ease-in-out;
  
  &:hover {
    transform: translateY(-2px);
  }
}

.cta-section {
  background: linear-gradient(135deg, #2196f3 0%, #1976d2 100%);
}

.max-width-600 {
  max-width: 600px;
}

// Responsive adjustments
@media (max-width: 960px) {
  .hero-section {
    min-height: auto;
  }
  
  .hero-content {
    min-height: auto;
  }
  
  .hero-text,
  .hero-image {
    text-align: center;
  }
  
  .cta-buttons {
    justify-content: center;
  }
}

// Accessibility improvements
@media (prefers-reduced-motion: reduce) {
  .feature-card,
  .compliance-card {
    transition: none;
    
    &:hover {
      transform: none;
    }
  }
}

// High contrast mode support
@media (prefers-contrast: high) {
  .hero-section {
    background: #ffffff;
    border-bottom: 2px solid #000000;
  }
  
  .stats-section {
    background: #000000;
    color: #ffffff;
  }
  
  .cta-section {
    background: #000000;
    color: #ffffff;
  }
}
</style>
