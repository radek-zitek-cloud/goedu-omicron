# GoEdu Frontend

A modern Vue.js 3 application built with TypeScript, Vite, and the Composition API for the GoEdu control testing platform.

## 🚀 Technology Stack

- **Vue.js 3** - Progressive JavaScript framework with Composition API
- **TypeScript** - Type-safe JavaScript development
- **Vite** - Fast build tool and development server
- **ESLint** - Code linting and error prevention
- **Prettier** - Code formatting and style consistency

## 📋 Project Features

- ✅ Vue 3 with Composition API (`<script setup>` syntax)
- ✅ Full TypeScript configuration
- ✅ Vite build system for fast development and optimized production builds
- ✅ ESLint configuration with Vue 3 + TypeScript support
- ✅ Prettier configuration for consistent code formatting
- ✅ Hot Module Replacement (HMR) for instant development feedback
- ✅ Production-ready build optimization

## 🛠️ Development Setup

### Prerequisites

- Node.js 18+
- npm 8+

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Open browser to http://localhost:5173
```

### Available Scripts

```bash
# Development
npm run dev          # Start development server with HMR

# Building
npm run build        # Build for production
npm run preview      # Preview production build locally

# Code Quality
npm run lint         # Run ESLint with auto-fix
npm run format       # Format code with Prettier
npm run format:check # Check code formatting
npm run type-check   # Run TypeScript compiler checks
```

## 📁 Project Structure

```
fe/
├── src/
│   ├── components/     # Vue components
│   ├── assets/        # Static assets (images, styles)
│   ├── App.vue        # Root component
│   └── main.ts        # Application entry point
├── public/            # Public static files
├── dist/             # Production build output
├── .vscode/          # VS Code settings
├── eslint.config.ts  # ESLint configuration
├── .prettierrc.json  # Prettier configuration
├── tsconfig.json     # TypeScript configuration
├── vite.config.ts    # Vite configuration
└── package.json      # Dependencies and scripts
```

## 🔧 Configuration Details

### TypeScript Configuration

The project uses a multi-file TypeScript configuration:

- `tsconfig.json` - Main configuration
- `tsconfig.app.json` - Application-specific settings
- `tsconfig.node.json` - Node.js/build tools settings

### ESLint Configuration

ESLint is configured with:

- Vue 3 essential rules
- TypeScript recommended rules
- Prettier integration for consistent formatting

### Vite Configuration

- Vue.js plugin for Single File Component support
- Vue DevTools integration for development
- Path aliases (`@` -> `src/`)
- Optimized production builds

## 🎯 Development Guidelines

### Composition API Usage

This project uses Vue 3's Composition API with the `<script setup>` syntax:

```vue
<script setup lang="ts">
import { ref, computed } from 'vue';

// Reactive state
const count = ref(0);

// Computed properties
const doubleCount = computed(() => count.value * 2);

// Props with TypeScript
defineProps<{
  title: string;
  optional?: boolean;
}>();

// Emits
const emit = defineEmits<{
  update: [value: string];
}>();
</script>
```

### Code Style

- Use TypeScript for all new files
- Follow Prettier formatting rules
- Use ESLint recommended practices
- Prefer Composition API over Options API
- Use `<script setup>` syntax for cleaner code

## 🚀 Production Build

The production build is optimized with:

- Tree shaking for smaller bundle sizes
- Code splitting for better caching
- Asset optimization and minification
- Gzip compression analysis

Build output is generated in the `dist/` directory and ready for deployment.

## 🔗 Integration Points

This frontend is designed to integrate with:

- GraphQL API backend (Go/gqlgen)
- OAuth 2.0 authentication system
- MongoDB data layer via GraphQL
- File upload/management system

## 📚 Next Steps

Future enhancements planned:

1. Pinia state management setup (TASK-011)
2. UI component library with Material Design (TASK-012)
3. Authentication integration (TASK-013)
4. GraphQL client setup
5. Responsive design framework

## 🤝 Contributing

1. Follow the established code style
2. Run linting and formatting before commits
3. Ensure TypeScript compilation passes
4. Test builds successfully
5. Write clear commit messages

---

This frontend application follows modern Vue.js 3 best practices and is production-ready for the GoEdu control testing platform.
