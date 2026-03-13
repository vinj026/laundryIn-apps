/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./app.vue",
    "./error.vue"
  ],
  darkMode: "class",
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Inter"', 'system-ui', 'sans-serif'],
        mono: ['"Roboto Mono"', 'ui-monospace', 'monospace'],
      },
      colors: {
        primary: {
          DEFAULT: '#2dd4bf',       // teal-400 — softer, more premium
          light: '#5eead4',         // teal-300
          dark: '#0d9488',          // teal-600
          container: '#042f2e',     // teal-950
          onContainer: '#99f6e4',   // teal-200
          text: '#021a19',          // near-black teal
        },
        secondary: {
          DEFAULT: '#a78bfa',       // violet-400
          container: '#1e1b4b',     // indigo-950
          onContainer: '#c4b5fd',   // violet-300
        },
        success: {
          DEFAULT: '#4ade80',
          muted: 'rgba(74, 222, 128, 0.15)',
        },
        warning: {
          DEFAULT: '#fbbf24',
          muted: 'rgba(251, 191, 36, 0.15)',
        },
        danger: {
          DEFAULT: '#f87171',
          muted: 'rgba(248, 113, 113, 0.15)',
        },
        surface: {
          DEFAULT: '#0a0a0a',       // near-black base
          raised: '#111111',        // slightly lifted
          container: '#161616',     // cards, panels
          containerHigh: '#1e1e1e', // elevated cards
          overlay: '#262626',       // modals, popovers
          onSurface: '#f0f0f0',     // primary text
          onSurfaceVariant: '#8a8a8a', // muted text
        },
        border: {
          DEFAULT: '#262626',
          subtle: '#1e1e1e',
          hover: '#404040',
        },
        outline: '#404040',
      },
      borderRadius: {
        'xl': '0.875rem',
        '2xl': '1rem',
        '3xl': '1.25rem',
        '4xl': '1.5rem',
      },
      transitionDuration: {
        fast: '150ms',
        normal: '250ms',
        slow: '400ms',
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.35s ease-out',
        'pulse-soft': 'pulseSoft 2s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(12px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        pulseSoft: {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.6' },
        },
      },
    },
  },
  plugins: [],
}
