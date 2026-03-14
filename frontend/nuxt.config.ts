// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  runtimeConfig: {
    public: {
      // Use HTTPS/WSS for production, fallback to localhost for development
      apiBase: process.env.NUXT_PUBLIC_API_BASE_URL ||
        (process.env.VERCEL ? 'https://laundryin-backend-production.up.railway.app/api/v1' : 'http://localhost:8080/api/v1'),
      wsBase: process.env.NUXT_PUBLIC_WS_BASE_URL ||
        (process.env.VERCEL ? 'wss://laundryin-backend-production.up.railway.app/api/v1/ws/connect' : 'ws://localhost:8080/api/v1/ws/connect')
    }
  },
  devtools: { enabled: true },
  future: {
    compatibilityVersion: 4,
  },
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  // SSR Proxy is only needed for local dev to avoid CORS. 
  // In production, we hit the API directly via useApiFetch.
  routeRules: process.env.VERCEL ? {} : {
    '/api/**': { proxy: 'http://localhost:8080/api/v1/**' }
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
      drop: process.env.NODE_ENV === 'production' ? ['console', 'debugger'] : []
    }
  }
})
