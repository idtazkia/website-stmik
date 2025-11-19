/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#e6eaf3',
          100: '#ccd5e7',
          200: '#99abcf',
          300: '#6681b7',
          400: '#33579f',
          500: '#194189', // Logo blue
          600: '#14346e',
          700: '#0f2752',
          800: '#0a1a37',
          900: '#050d1b',
        },
        secondary: {
          50: '#fef3e9',
          100: '#fde7d3',
          200: '#fbcfa7',
          300: '#f9b77b',
          400: '#f79f4f',
          500: '#EE7B1D', // Logo orange
          600: '#be6217',
          700: '#8f4a11',
          800: '#5f310c',
          900: '#301906',
        },
      },
      fontFamily: {
        sans: [
          'system-ui',
          '-apple-system',
          'BlinkMacSystemFont',
          '"Segoe UI"',
          'Roboto',
          '"Helvetica Neue"',
          'Arial',
          'sans-serif',
        ],
      },
    },
  },
  plugins: [],
};
