/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: 'var(--primary-color)',
        secondary: 'var(--secondary-color)',
        accent: 'var(--accent-color)',
        'text-primary': 'var(--text-color)',
        'bg-primary': 'var(--bg-color)',
        'bg-secondary': 'var(--bg-secondary)',
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
  // Оптимизации
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true,
  },
  purge: {
    // Удаляем неиспользуемые стили
    enabled: process.env.NODE_ENV === 'production',
    content: [
      './src/**/*.vue',
      './src/**/*.js',
      './src/**/*.ts',
      './src/**/*.jsx',
      './src/**/*.tsx',
      './public/**/*.html',
    ],
  },
} 