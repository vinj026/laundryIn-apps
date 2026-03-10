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
        sans: ['"Roboto"', 'sans-serif'],
      },
      colors: {
        primary: {
          DEFAULT: '#00e5ff',
          container: '#004f5e',
          onContainer: '#82f3ff',
          text: '#003641',
        },
        secondary: {
          container: '#334b4f',
          onContainer: '#cde7ec',
        },
        surface: {
          dim: '#0c1012',
          DEFAULT: '#131b1f',
          container: '#1c252a',
          containerHigh: '#273136',
          onSurface: '#e1e3e4',
          onSurfaceVariant: '#9ba8ab',
          outlineVariant: '#3f484a',
        },
        outline: '#899294',
      },
      borderRadius: {
        '3xl': '1.75rem',
        '4xl': '2.5rem',
      }
    },
  },
  plugins: [],
}
