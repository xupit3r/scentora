/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Notion-inspired color palette
        text: {
          primary: '#37352F',
          secondary: '#787774',
          tertiary: '#9B9A97',
        },
        bg: {
          primary: '#FFFFFF',
          secondary: '#FAFAFA',
          tertiary: '#F7F6F3',
        },
        border: {
          light: '#E9E9E7',
          medium: '#D9D9D7',
        },
        accent: {
          DEFAULT: '#0F766E',
          hover: '#0D6460',
        },
        position: {
          top: {
            bg: '#FEF3C7',
            text: '#92400E',
          },
          middle: {
            bg: '#E9D5FF',
            text: '#6B21A8',
          },
          base: {
            bg: '#F5E6D3',
            text: '#78350F',
          },
        },
      },
      fontFamily: {
        sans: [
          '-apple-system',
          'BlinkMacSystemFont',
          '"Segoe UI"',
          'Roboto',
          '"Helvetica Neue"',
          'Arial',
          'sans-serif',
        ],
      },
      fontSize: {
        xs: '12px',
        sm: '14px',
        base: '16px',
        lg: '20px',
        xl: '24px',
        '2xl': '32px',
        '3xl': '40px',
      },
      spacing: {
        '1': '4px',
        '2': '8px',
        '3': '12px',
        '4': '16px',
        '5': '20px',
        '6': '24px',
        '8': '32px',
        '10': '40px',
        '12': '48px',
        '16': '64px',
        '20': '80px',
      },
      borderRadius: {
        'sm': '4px',
        'DEFAULT': '8px',
        'md': '12px',
        'lg': '16px',
      },
      boxShadow: {
        'sm': '0 1px 2px rgba(0, 0, 0, 0.04)',
        'DEFAULT': '0 1px 3px rgba(0, 0, 0, 0.06)',
        'md': '0 4px 6px rgba(0, 0, 0, 0.06)',
        'lg': '0 10px 15px rgba(0, 0, 0, 0.08)',
        'hover': '0 8px 16px rgba(0, 0, 0, 0.1)',
      },
    },
  },
  plugins: [],
}
