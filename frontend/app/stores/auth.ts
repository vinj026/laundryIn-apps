import { defineStore } from 'pinia'

interface User {
  id: string
  name: string
  phone: string
  email?: string
  role: 'owner' | 'customer'
}

interface AuthState {
  token: string | null
  user: User | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: null,
    user: null
  }),

  getters: {
    isLoggedIn: (state): boolean => !!state.token,
    isOwner: (state): boolean => state.user?.role === 'owner',
    isCustomer: (state): boolean => state.user?.role === 'customer',
    authHeader: (state): string =>
      state.token ? `Bearer ${state.token}` : ''
  },

  actions: {
    setAuth(token: string, user: User) {
      this.token = token
      this.user = user
      if (process.client) {
        localStorage.setItem('laundryin_token', token)
        localStorage.setItem('laundryin_user', JSON.stringify(user))
      }
    },

    logout() {
      this.token = null
      this.user = null
      if (process.client) {
        localStorage.removeItem('laundryin_token')
        localStorage.removeItem('laundryin_user')
      }
    },

    restoreSession() {
      if (process.client) {
        const token = localStorage.getItem('laundryin_token')
        const userStr = localStorage.getItem('laundryin_user')
        if (token && userStr) {
          try {
            this.token = token
            this.user = JSON.parse(userStr)
          } catch {
            this.logout()
          }
        }
      }
    }
  }
})
