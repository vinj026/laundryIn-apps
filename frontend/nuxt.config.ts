// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  runtimeConfig: {
    public: {
      // Use environment variables for production, fallback to proxy/local dev
      // Default to /api so the proxy in routeRules (and Vercel rewrite) is used
      apiBase: process.env.NUXT_PUBLIC_API_BASE_URL || '/api',
      wsBase: process.env.NUXT_PUBLIC_WS_BASE_URL || 'ws://localhost:8080/api/v1/ws/connect'
    }
  },
  devtools: { enabled: true },
  future: {
    compatibilityVersion: 4,
  },
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  // Nuxt Proxy: Standardizing on SSR Proxy for production to handle /api/v1 mapping and avoid CORS.
  routeRules: {
    '/api/**': {
      proxy: (process.env.BACKEND_URL || 'http://localhost:8080').replace(/\/+$/, '') + '/api/v1/**'
    }
  },
  css: ['~/assets/css/main.css'],
  app: {
    pageTransition: { name: 'page', mode: 'out-in' },
    head: {
      title: 'LaundryIn — Modern Laundry Platform',
      meta: [
        { name: 'description', content: 'Smart online laundry platform with real-time tracking and management dashboard.' },
        { name: 'theme-color', content: '#0a0a0a' },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Roboto+Mono:wght@400;500;700&display=swap' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200&display=swap' }
      ]
    }
  },
  vite: {
    esbuild: {
      // Temporarily disabled console drop for remote debugging (Phase 2)
      drop: []
    }
  }
})
