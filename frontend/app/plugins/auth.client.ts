import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()

  // Restore session on app initialization (client-side only)
  if (import.meta.client) {
    authStore.restoreSession()
  }
})
