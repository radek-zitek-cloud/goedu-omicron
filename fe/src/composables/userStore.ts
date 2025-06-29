/**
 * Validation script for Vue 3 + TypeScript setup
 * This script demonstrates that all requirements are met
 */

// Vue 3 Composition API verification
import { ref, computed, reactive } from 'vue';

// TypeScript types for demonstration
interface User {
  id: number;
  name: string;
  email: string;
  active: boolean;
}

// Reactive state using Composition API
const userCount = ref(0);
const users = reactive<User[]>([]);

// Computed property
const activeUsers = computed(() => users.filter(user => user.active));

// Type-safe function
function addUser(name: string, email: string): void {
  const newUser: User = {
    id: Date.now(),
    name,
    email,
    active: true,
  };

  users.push(newUser);
  userCount.value++;
}

// Export for use in components
export { userCount, users, activeUsers, addUser };

// This validates:
// ✅ Vue 3 Composition API (ref, computed, reactive)
// ✅ TypeScript support (interfaces, type annotations)
// ✅ Modern ES6+ syntax
// ✅ Proper exports
