# 🧪 Testing Infrastructure Setup Guide

## 📋 Current State: No Testing Infrastructure

The GoEdu Omicron project currently has **zero testing infrastructure**, which is critical for a banking/financial application. This guide provides step-by-step instructions to establish comprehensive testing.

## 🎯 Testing Requirements for Banking Applications

### **Minimum Requirements**
- **Unit Test Coverage**: 80%+ backend, 70%+ frontend
- **Integration Tests**: Database, API, external services
- **Security Tests**: Authentication, authorization, input validation
- **End-to-End Tests**: Complete user workflows
- **Performance Tests**: Load testing for control processing
- **Compliance Tests**: Audit logging, data retention

## 🏗️ Backend Testing Setup (Go)

### **Step 1: Add Testing Dependencies**

```bash
cd be

# Add testing dependencies to go.mod
go mod tidy

# Install testing tools
go get -t ./...
```

### **Step 2: Create Test Structure**

```bash
# Create test directories
mkdir -p {cmd,internal,pkg}/*/
mkdir -p test/{fixtures,helpers,integration,e2e}

# Create test files structure:
# internal/
# ├── config/
# │   ├── config.go
# │   └── config_test.go
# ├── models/
# │   ├── models.go
# │   └── models_test.go
# ├── services/
# │   ├── interfaces.go
# │   ├── control_service.go
# │   └── control_service_test.go
# └── repositories/
#     ├── interfaces.go
#     ├── mongo_control_repository.go
#     └── mongo_control_repository_test.go
```

### **Step 3: Create Test Utilities**

Create `test/helpers/database.go`:
```go
package helpers

import (
    "context"
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/mongodb"
)

// TestDatabase provides test database utilities
type TestDatabase struct {
    Container *mongodb.MongoDBContainer
    URI       string
}

// NewTestDatabase creates a test MongoDB container
func NewTestDatabase(t *testing.T) *TestDatabase {
    ctx := context.Background()
    
    mongodbContainer, err := mongodb.RunContainer(ctx,
        testcontainers.WithImage("mongo:6"),
    )
    if err != nil {
        t.Fatal(err)
    }
    
    uri, err := mongodbContainer.ConnectionString(ctx)
    if err != nil {
        t.Fatal(err)
    }
    
    return &TestDatabase{
        Container: mongodbContainer,
        URI:       uri,
    }
}

// Cleanup terminates the test database
func (td *TestDatabase) Cleanup(t *testing.T) {
    if err := td.Container.Terminate(context.Background()); err != nil {
        t.Errorf("failed to terminate container: %s", err)
    }
}
```

### **Step 4: Example Unit Test**

Create `internal/models/models_test.go`:
```go
package models

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestOrganization_Validate(t *testing.T) {
    tests := []struct {
        name    string
        org     Organization
        wantErr bool
    }{
        {
            name: "valid organization",
            org: Organization{
                Name:         "Test Bank",
                Slug:         "test-bank",
                ContactEmail: "admin@testbank.com",
                Status:       "active",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            org: Organization{
                Name:         "Test Bank",
                Slug:         "test-bank",
                ContactEmail: "invalid-email",
                Status:       "active",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.org.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestControl_CalculateRiskScore(t *testing.T) {
    control := &Control{
        RiskLevel:  "High",
        Importance: "Critical",
    }

    score := control.CalculateRiskScore()
    assert.Equal(t, 9, score) // High risk + Critical importance
}
```

### **Step 5: Example Integration Test**

Create `test/integration/control_service_test.go`:
```go
package integration

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/radek-zitek-cloud/goedu-omicron/be/internal/models"
    "github.com/radek-zitek-cloud/goedu-omicron/be/internal/services"
    "github.com/radek-zitek-cloud/goedu-omicron/be/test/helpers"
)

func TestControlService_Integration(t *testing.T) {
    // Setup test database
    testDB := helpers.NewTestDatabase(t)
    defer testDB.Cleanup(t)

    // Initialize service with test database
    service := services.NewControlService(testDB.URI)
    ctx := context.Background()

    t.Run("CreateControl", func(t *testing.T) {
        input := &services.CreateControlInput{
            OrganizationID: "test-org",
            ControlID:      "TEST-001",
            Title:          "Test Control",
            Description:    "Test control description",
            Framework:      "COSO",
            Category:       "IT Controls",
            RiskLevel:      "Medium",
            Importance:     "High",
        }

        control, err := service.CreateControl(ctx, input)
        require.NoError(t, err)
        assert.Equal(t, input.Title, control.Title)
        assert.NotEmpty(t, control.ID)
    })
}
```

### **Step 6: Add Make Targets**

Update `be/Makefile`:
```makefile
# Testing targets
test: ## Run unit tests
	@echo "🧪 Running unit tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

test-integration: ## Run integration tests
	@echo "🧪 Running integration tests..."
	@go test -v -race -tags=integration ./test/integration/...

test-all: test test-integration ## Run all tests

benchmark: ## Run benchmark tests
	@echo "🏃 Running benchmarks..."
	@go test -bench=. -benchmem ./...
```

## 🎭 Frontend Testing Setup (Vue.js)

### **Step 1: Install Testing Dependencies**

```bash
cd fe

# Install Vitest and testing utilities
npm install --save-dev \
  vitest \
  @vitest/ui \
  jsdom \
  @vue/test-utils \
  @testing-library/vue \
  happy-dom \
  msw

# For E2E testing
npm install --save-dev \
  playwright \
  @playwright/test
```

### **Step 2: Configure Vitest**

Create `fe/vitest.config.ts`:
```typescript
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
})
```

### **Step 3: Create Test Setup**

Create `fe/src/test/setup.ts`:
```typescript
import { expect, afterEach } from 'vitest'
import { cleanup } from '@testing-library/vue'
import * as matchers from '@testing-library/jest-dom/matchers'

// Extends Vitest's expect with Testing Library matchers
expect.extend(matchers)

// Cleanup after each test
afterEach(() => {
  cleanup()
})

// Mock IntersectionObserver for Vuetify components
global.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  unobserve() {}
}
```

### **Step 4: Example Component Test**

Create `fe/src/components/__tests__/HelloWorld.test.ts`:
```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import HelloWorld from '../HelloWorld.vue'

describe('HelloWorld', () => {
  it('renders properly', () => {
    const wrapper = mount(HelloWorld, { props: { msg: 'Hello Vitest' } })
    expect(wrapper.text()).toContain('Hello Vitest')
  })

  it('increments counter when button clicked', async () => {
    const wrapper = mount(HelloWorld, { props: { msg: 'Test' } })
    const button = wrapper.find('button')
    
    await button.trigger('click')
    expect(wrapper.text()).toContain('count is 1')
  })
})
```

### **Step 5: Example Store Test**

Create `fe/src/stores/__tests__/auth.test.ts`:
```typescript
import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../auth'

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with correct default state', () => {
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
    expect(store.user).toBeNull()
    expect(store.token).toBeNull()
  })

  it('logs in user correctly', async () => {
    const store = useAuthStore()
    const mockUser = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User'
    }

    await store.login('test@example.com', 'password')
    
    expect(store.isAuthenticated).toBe(true)
    expect(store.user).toMatchObject(mockUser)
  })
})
```

### **Step 6: Update Package.json**

```json
{
  "scripts": {
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:run": "vitest run",
    "test:coverage": "vitest run --coverage",
    "test:e2e": "playwright test",
    "test:e2e:headed": "playwright test --headed"
  }
}
```

## 🌐 End-to-End Testing Setup

### **Step 1: Initialize Playwright**

```bash
cd fe
npx playwright install
```

### **Step 2: Create E2E Test**

Create `fe/tests/e2e/auth.spec.ts`:
```typescript
import { test, expect } from '@playwright/test'

test.describe('Authentication Flow', () => {
  test('user can login successfully', async ({ page }) => {
    await page.goto('/')
    
    // Navigate to login
    await page.click('text=Login')
    
    // Fill login form
    await page.fill('[data-testid=email]', 'test@example.com')
    await page.fill('[data-testid=password]', 'password')
    await page.click('[data-testid=login-button]')
    
    // Verify successful login
    await expect(page).toHaveURL('/dashboard')
    await expect(page.locator('[data-testid=user-menu]')).toBeVisible()
  })

  test('shows error for invalid credentials', async ({ page }) => {
    await page.goto('/login')
    
    await page.fill('[data-testid=email]', 'invalid@example.com')
    await page.fill('[data-testid=password]', 'wrongpassword')
    await page.click('[data-testid=login-button]')
    
    await expect(page.locator('[data-testid=error-message]')).toBeVisible()
    await expect(page.locator('[data-testid=error-message]')).toContainText('Invalid credentials')
  })
})
```

## 🚀 CI/CD Pipeline Integration

### **GitHub Actions Workflow**

Create `.github/workflows/test.yml`:
```yaml
name: Test Suite

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    services:
      mongodb:
        image: mongo:6
        ports:
          - 27017:27017
      redis:
        image: redis:7
        ports:
          - 6379:6379

    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      run: |
        cd be
        go mod download
    
    - name: Run tests
      run: |
        cd be
        go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./be/coverage.out

  frontend-tests:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
        cache: 'npm'
        cache-dependency-path: fe/package-lock.json
    
    - name: Install dependencies
      run: |
        cd fe
        npm ci

    - name: Run unit tests
      run: |
        cd fe
        npm run test:run

    - name: Run E2E tests
      run: |
        cd fe
        npm run test:e2e

  security-tests:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Run Snyk security scan
      uses: snyk/actions/golang@master
      env:
        SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
      with:
        args: --severity-threshold=high
```

## 📊 Test Coverage Goals

### **Backend Coverage Targets**
- **Unit Tests**: 85%+ line coverage
- **Integration Tests**: All API endpoints
- **Repository Tests**: All CRUD operations
- **Service Tests**: All business logic

### **Frontend Coverage Targets**
- **Component Tests**: 75%+ coverage
- **Store Tests**: 90%+ coverage
- **Utility Tests**: 95%+ coverage
- **E2E Tests**: Critical user paths

## 🔧 Running Tests

### **Backend**
```bash
cd be

# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
make test-integration

# Run benchmarks
make benchmark
```

### **Frontend**
```bash
cd fe

# Run unit tests
npm run test

# Run with UI
npm run test:ui

# Run with coverage
npm run test:coverage

# Run E2E tests
npm run test:e2e
```

## 📈 Next Steps

1. **Week 1**: Set up basic unit testing infrastructure
2. **Week 2**: Add integration tests for core services
3. **Week 3**: Implement E2E tests for critical workflows
4. **Week 4**: Add performance and security tests
5. **Ongoing**: Maintain 80%+ coverage for all new code

---

**Priority**: Implement testing infrastructure before adding new features.
**Timeline**: Complete basic setup within 1 week.
**Success Metric**: All tests pass in CI/CD pipeline.