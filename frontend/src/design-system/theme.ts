/**
 * Naive UI Theme Configuration
 * Customizes Naive UI components to match Notion-inspired design
 */

import { GlobalThemeOverrides } from 'naive-ui';
import { colors, typography, spacing, borderRadius, shadows } from './tokens';

export const naiveTheme: GlobalThemeOverrides = {
  common: {
    fontFamily: typography.fontFamily.sans,
    fontSize: typography.fontSize.base,
    fontWeightStrong: typography.fontWeight.semibold.toString(),
    
    // Primary colors
    primaryColor: colors.accent.primary,
    primaryColorHover: colors.accent.hover,
    primaryColorPressed: colors.accent.hover,
    primaryColorSuppl: colors.accent.light,
    
    // Text colors
    textColorBase: colors.text.primary,
    textColor1: colors.text.primary,
    textColor2: colors.text.secondary,
    textColor3: colors.text.tertiary,
    
    // Background colors
    bodyColor: colors.bg.primary,
    cardColor: colors.bg.primary,
    modalColor: colors.bg.primary,
    popoverColor: colors.bg.primary,
    
    // Border colors
    borderColor: colors.border.light,
    dividerColor: colors.border.light,
    
    // Border radius
    borderRadius: borderRadius.base,
    borderRadiusSmall: borderRadius.sm,
    
    // Shadows
    boxShadow1: shadows.sm,
    boxShadow2: shadows.base,
    boxShadow3: shadows.md,
  },
  
  Button: {
    fontWeight: typography.fontWeight.medium.toString(),
    borderRadius: borderRadius.base,
    paddingMedium: `${spacing['2']} ${spacing['4']}`,
    heightMedium: '40px',
    
    // Primary button
    colorPrimary: colors.accent.primary,
    colorHoverPrimary: colors.accent.hover,
    colorPressedPrimary: colors.accent.hover,
    
    // Default button
    color: colors.bg.secondary,
    colorHover: colors.bg.tertiary,
    colorPressed: colors.bg.tertiary,
    
    textColor: colors.text.primary,
    border: `1px solid ${colors.border.light}`,
  },
  
  Input: {
    borderRadius: borderRadius.base,
    heightMedium: '40px',
    paddingMedium: spacing['3'],
    fontSizeMedium: typography.fontSize.sm,
    
    border: `1px solid ${colors.border.light}`,
    borderHover: `1px solid ${colors.border.medium}`,
    borderFocus: `1px solid ${colors.accent.primary}`,
    
    color: colors.bg.primary,
    colorFocus: colors.bg.primary,
    
    textColor: colors.text.primary,
    placeholderColor: colors.text.tertiary,
  },
  
  Card: {
    borderRadius: borderRadius.md,
    paddingMedium: spacing['6'],
    borderColor: colors.border.light,
    color: colors.bg.primary,
    
    boxShadow: shadows.sm,
  },
  
  Dialog: {
    borderRadius: borderRadius.lg,
    padding: spacing['6'],
    
    color: colors.bg.primary,
    textColor: colors.text.primary,
    titleFontSize: typography.fontSize.lg,
    titleFontWeight: typography.fontWeight.semibold.toString(),
  },
  
  Tag: {
    borderRadius: borderRadius.sm,
    padding: `${spacing['1']} ${spacing['3']}`,
    fontSize: typography.fontSize.xs,
    fontWeight: typography.fontWeight.medium.toString(),
  },
  
  Select: {
    peers: {
      InternalSelection: {
        borderRadius: borderRadius.base,
        heightMedium: '40px',
        
        border: `1px solid ${colors.border.light}`,
        borderHover: `1px solid ${colors.border.medium}`,
        borderActive: `1px solid ${colors.accent.primary}`,
        borderFocus: `1px solid ${colors.accent.primary}`,
      },
    },
  },
  
  Menu: {
    borderRadius: borderRadius.base,
    itemHeight: '40px',
    
    itemTextColor: colors.text.primary,
    itemTextColorHover: colors.text.primary,
    itemTextColorActive: colors.accent.primary,
    itemTextColorActiveHover: colors.accent.primary,
    
    itemColorHover: colors.bg.hover,
    itemColorActive: colors.accent.light,
    itemColorActiveHover: colors.accent.light,
    
    itemIconColor: colors.text.secondary,
    itemIconColorHover: colors.text.primary,
    itemIconColorActive: colors.accent.primary,
  },
  
  Notification: {
    borderRadius: borderRadius.md,
    padding: spacing['4'],
    
    color: colors.bg.primary,
    textColor: colors.text.primary,
    
    boxShadow: shadows.lg,
  },
};
