/**
 * Design System Tokens
 * Notion-inspired design tokens for Scentora
 */

export const colors = {
  // Text colors
  text: {
    primary: '#37352F',
    secondary: '#787774',
    tertiary: '#9B9A97',
    inverse: '#FFFFFF',
  },

  // Background colors
  bg: {
    primary: '#FFFFFF',
    secondary: '#FAFAFA',
    tertiary: '#F7F6F3',
    hover: '#F7F6F3',
  },

  // Border colors
  border: {
    light: '#E9E9E7',
    medium: '#D9D9D7',
  },

  // Accent colors
  accent: {
    primary: '#0F766E',
    hover: '#0D6460',
    light: '#CCFBF1',
  },

  // Position colors (subtle, refined)
  position: {
    top: {
      bg: '#FEF3C7',
      text: '#92400E',
      border: '#FDE68A',
    },
    middle: {
      bg: '#E9D5FF',
      text: '#6B21A8',
      border: '#D8B4FE',
    },
    base: {
      bg: '#F5E6D3',
      text: '#78350F',
      border: '#E4C9A0',
    },
  },

  // Semantic colors
  success: {
    bg: '#D1FAE5',
    text: '#065F46',
    border: '#6EE7B7',
  },
  warning: {
    bg: '#FEF3C7',
    text: '#92400E',
    border: '#FDE68A',
  },
  error: {
    bg: '#FEE2E2',
    text: '#991B1B',
    border: '#FCA5A5',
  },
  info: {
    bg: '#DBEAFE',
    text: '#1E40AF',
    border: '#93C5FD',
  },
} as const;

export const typography = {
  fontFamily: {
    sans: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
    mono: '"SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas, "Courier New", monospace',
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

  fontWeight: {
    normal: 400,
    medium: 500,
    semibold: 600,
  },

  lineHeight: {
    tight: 1.2,
    normal: 1.5,
    relaxed: 1.7,
  },
} as const;

export const spacing = {
  '0': '0px',
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
} as const;

export const borderRadius = {
  none: '0px',
  sm: '4px',
  base: '8px',
  md: '12px',
  lg: '16px',
  full: '9999px',
} as const;

export const shadows = {
  none: 'none',
  sm: '0 1px 2px rgba(0, 0, 0, 0.04)',
  base: '0 1px 3px rgba(0, 0, 0, 0.06)',
  md: '0 4px 6px rgba(0, 0, 0, 0.06)',
  lg: '0 10px 15px rgba(0, 0, 0, 0.08)',
  hover: '0 8px 16px rgba(0, 0, 0, 0.1)',
} as const;

export const transitions = {
  fast: '100ms ease',
  base: '200ms ease',
  slow: '300ms ease',
} as const;

export const zIndex = {
  base: 1,
  dropdown: 1000,
  sticky: 1100,
  fixed: 1200,
  modalBackdrop: 1300,
  modal: 1400,
  popover: 1500,
  tooltip: 1600,
} as const;

export const breakpoints = {
  mobile: '640px',
  tablet: '768px',
  desktop: '1024px',
  wide: '1280px',
} as const;

export const layout = {
  sidebarWidth: {
    expanded: '280px',
    collapsed: '64px',
  },
  headerHeight: '64px',
  containerMaxWidth: '1400px',
  contentPadding: {
    mobile: spacing['4'],
    tablet: spacing['6'],
    desktop: spacing['12'],
  },
} as const;
