# Mobile Responsive Design Guidelines

Guidelines for creating responsive, mobile-friendly interfaces in Scentora.

## 📱 Breakpoints

We follow a mobile-first approach with **8 granular breakpoints** to support all device form factors gracefully:

| Name | Breakpoint | Devices | Grid Columns |
|------|------------|---------|--------------|
| **xs** | 320px - 479px | Small mobile phones | 1 column |
| **sm** | 480px - 639px | Large mobile phones | 1 column |
| **md** | 640px - 767px | Small tablets (portrait) | 2 columns |
| **lg** | 768px - 1023px | Tablets | 2 columns |
| **xl** | 1024px - 1279px | Small desktops/laptops | 3 columns |
| **2xl** | 1280px - 1535px | Standard desktops | 4 columns |
| **3xl** | 1536px - 1919px | Large desktops | 5 columns |
| **4xl** | 1920px+ | Ultra-wide / 4K displays | 6 columns |

### Usage in CSS

```css
/* Mobile-first (xs): Default styles */
.element { 
  padding: 16px;
  font-size: 16px;
}

/* sm: Large phones (480px+) */
@media (min-width: 480px) {
  .element { padding: 20px; }
}

/* md: Small tablets (640px+) */
@media (min-width: 640px) {
  .element { padding: 24px; }
}

/* lg: Tablets (768px+) */
@media (min-width: 768px) {
  .element { 
    padding: 24px;
    font-size: 18px;
  }
}

/* xl: Small desktops (1024px+) */
@media (min-width: 1024px) {
  .element { padding: 32px; }
}

/* 2xl: Standard desktops (1280px+) */
@media (min-width: 1280px) {
  .element { 
    padding: 40px;
    font-size: 20px;
  }
}

/* 3xl: Large desktops (1536px+) */
@media (min-width: 1536px) {
  .element { padding: 48px; }
}

/* 4xl: Ultra-wide (1920px+) */
@media (min-width: 1920px) {
  .element { 
    padding: 64px;
    font-size: 24px;
  }
}
```

## 🎯 Mobile Design Principles

### 1. Touch-Friendly Targets
- **Minimum size**: 44x44px for all interactive elements
- **Spacing**: 8px minimum between touch targets
- **Visual feedback**: Active states for all buttons/links

```css
.touch-target {
  min-width: 44px;
  min-height: 44px;
  padding: var(--spacing-3);
}
```

### 2. Navigation Patterns

**Mobile (< 768px):**
- Hamburger menu button
- Slide-out drawer from left
- Full-screen overlay when drawer is open
- Close button in drawer header

**Tablet/Desktop (≥ 768px):**
- Fixed sidebar (280px wide)
- Always visible
- Collapsible to icon-only (future)

### 3. Content Adaptation

**Extra Small Mobile (xs: 320-479px):**
- Minimal padding (16px)
- Smallest font sizes
- Compact spacing
- Essential content only
- Single column

**Small Mobile (sm: 480-639px):**
- Slightly increased padding (20px)
- Single column maintained
- Better spacing between elements

**Medium (md: 640-767px):**
- Two-column grids where appropriate
- Standard padding (24px)
- More comfortable spacing

**Large / Tablet (lg: 768-1023px):**
- Two-column layouts standard
- Expanded spacing
- Sidebar can be collapsible
- More detailed information shown

**Extra Large / Desktop (xl: 1024-1279px):**
- Three-column grids
- Fixed sidebar (280px)
- Full feature set
- Standard desktop spacing

**2X Large / Standard Desktop (2xl: 1280-1535px):**
- Four-column grids
- Generous whitespace
- Comfortable reading widths
- Full navigation always visible

**3X Large / Large Desktop (3xl: 1536-1919px):**
- Five-column grids
- Maximum information density
- Wide content areas
- Enhanced spacing

**4X Large / Ultra-wide (4xl: 1920px+):**
- Six-column grids
- Optimal for multi-tasking
- Extra padding (64px+)
- Larger typography scale

## 📐 Component Patterns

### Header

**Mobile:**
```html
<header class="mobile-header">
  <button class="menu-button">☰</button>
  <h1 class="header-title">Page Title</h1>
  <button class="header-action">🔍</button>
</header>
```
- Height: 56px
- Centered title
- Hamburger menu on left
- Single action on right

**Desktop:**
```html
<header class="desktop-header">
  <nav class="breadcrumbs">...</nav>
  <div class="search-bar">...</div>
  <div class="actions">...</div>
</header>
```
- Height: 64px
- Breadcrumbs on left
- Search bar in center
- Multiple actions on right

### Navigation Drawer

**Mobile Implementation:**
```html
<!-- Overlay -->
<div class="drawer-overlay" onclick="closeDrawer()"></div>

<!-- Drawer -->
<nav class="mobile-drawer">
  <div class="drawer-header">
    <div class="logo">...</div>
    <button class="close-drawer">×</button>
  </div>
  <div class="drawer-nav">...</div>
  <div class="drawer-footer">...</div>
</nav>
```

**CSS:**
```css
.mobile-drawer {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 280px;
  transform: translateX(-100%);
  transition: transform 250ms ease-in-out;
}

.mobile-drawer.active {
  transform: translateX(0);
}
```

### Cards

**Mobile (xs-sm: < 640px):**
- Full width (single column)
- Compact padding (16px)
- Stacked information
- Essential info only

**Small Tablet (md: 640-767px):**
- Two-column grid
- Standard padding (20px)

**Tablet (lg: 768-1023px):**
- Two-column grid
- Expanded information
- Hover effects enabled

**Small Desktop (xl: 1024-1279px):**
- Three-column grid
- Standard padding (20px)

**Standard Desktop (2xl: 1280-1535px):**
- Four-column grid
- Full information
- Enhanced hover effects

**Large Desktop (3xl: 1536-1919px):**
- Five-column grid
- Comfortable spacing

**Ultra-wide (4xl: 1920px+):**
- Six-column grid
- Maximum density
- Extra spacing between cards

```css
.card-grid {
  display: grid;
  gap: var(--spacing-4);
  grid-template-columns: 1fr; /* xs: 1 col */
}

@media (min-width: 480px) {
  .card-grid {
    gap: var(--spacing-5); /* sm: better spacing */
  }
}

@media (min-width: 640px) {
  .card-grid {
    grid-template-columns: repeat(2, 1fr); /* md: 2 cols */
  }
}

@media (min-width: 768px) {
  .card-grid {
    gap: var(--spacing-6); /* lg: expanded spacing */
  }
}

@media (min-width: 1024px) {
  .card-grid {
    grid-template-columns: repeat(3, 1fr); /* xl: 3 cols */
  }
}

@media (min-width: 1280px) {
  .card-grid {
    grid-template-columns: repeat(4, 1fr); /* 2xl: 4 cols */
  }
}

@media (min-width: 1536px) {
  .card-grid {
    grid-template-columns: repeat(5, 1fr); /* 3xl: 5 cols */
  }
}

@media (min-width: 1920px) {
  .card-grid {
    grid-template-columns: repeat(6, 1fr); /* 4xl: 6 cols */
    gap: var(--spacing-8);
  }
}
```

### Floating Action Button (FAB)

**Mobile Only:**
- Fixed position: bottom-right
- 56x56px circular button
- Prominent shadow
- Primary action (e.g., "Add Accord")

```css
.fab {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  box-shadow: 0 4px 12px rgba(15, 118, 110, 0.4);
}

@media (min-width: 768px) {
  .fab {
    display: none; /* Use toolbar button instead */
  }
}
```

### Forms

**Mobile:**
- Single column
- Full-width inputs
- Larger input heights (44px)
- Stacked labels

**Desktop:**
- Multi-column where appropriate
- Standard input heights
- Inline labels for short fields

## 🎨 Typography Scaling

Scale font sizes appropriately for mobile:

```css
:root {
  --font-size-xs: 12px;
  --font-size-sm: 14px;
  --font-size-base: 16px;
  --font-size-lg: 18px;   /* Reduced from 20px on mobile */
  --font-size-xl: 22px;   /* Reduced from 24px on mobile */
  --font-size-2xl: 28px;  /* Reduced from 32px on mobile */
}

@media (min-width: 768px) {
  :root {
    --font-size-lg: 20px;
    --font-size-xl: 24px;
    --font-size-2xl: 32px;
  }
}
```

## 📏 Spacing

Spacing scales fluidly across breakpoints to maintain optimal density:

```css
.content-wrapper {
  padding: 16px; /* xs: mobile */
}

@media (min-width: 480px) {
  .content-wrapper {
    padding: 20px; /* sm: large phone */
  }
}

@media (min-width: 640px) {
  .content-wrapper {
    padding: 24px; /* md: small tablet */
  }
}

@media (min-width: 768px) {
  .content-wrapper {
    padding: 24px; /* lg: tablet */
  }
}

@media (min-width: 1024px) {
  .content-wrapper {
    padding: 32px; /* xl: small desktop */
  }
}

@media (min-width: 1280px) {
  .content-wrapper {
    padding: 40px; /* 2xl: desktop */
  }
}

@media (min-width: 1536px) {
  .content-wrapper {
    padding: 48px; /* 3xl: large desktop */
  }
}

@media (min-width: 1920px) {
  .content-wrapper {
    padding: 64px; /* 4xl: ultra-wide */
  }
}
```

### Container Max-Widths

Use the `.container` class for centered, max-width constrained content:

```css
.container {
  width: 100%;
  margin: 0 auto;
  padding: 0 16px;
  max-width: 100%; /* xs-sm: full width */
}

@media (min-width: 640px) {
  .container {
    max-width: 640px; /* md */
  }
}

@media (min-width: 768px) {
  .container {
    max-width: 768px; /* lg */
    padding: 0 24px;
  }
}

@media (min-width: 1024px) {
  .container {
    max-width: 1024px; /* xl */
  }
}

@media (min-width: 1280px) {
  .container {
    max-width: 1280px; /* 2xl */
  }
}

@media (min-width: 1536px) {
  .container {
    max-width: 1536px; /* 3xl */
  }
}

@media (min-width: 1920px) {
  .container {
    max-width: 1920px; /* 4xl */
    padding: 0 32px;
  }
}
```

## 🔄 Responsive Utilities

Use utility classes for show/hide patterns across all breakpoints:

```html
<!-- Show only on extra small mobile -->
<div class="show-xs hide-sm">Visible only on phones < 480px</div>

<!-- Show only on small mobile -->
<div class="hide-xs show-sm hide-md">Visible only on phones 480-639px</div>

<!-- Show only on small tablets -->
<div class="hide-sm show-md hide-lg">Visible only on 640-767px</div>

<!-- Show on tablets and up -->
<div class="hide-xs hide-sm show-lg">Visible on tablets and larger</div>

<!-- Show only on desktop -->
<div class="hide-mobile show-desktop">Desktop only (1024px+)</div>

<!-- Show on large desktops only -->
<div class="hide-xl show-2xl">Visible on 1280px+ screens</div>

<!-- Show on ultra-wide only -->
<div class="hide-3xl show-4xl">Ultra-wide only (1920px+)</div>
```

### Available Utility Classes

| Class | Breakpoint | Behavior |
|-------|------------|----------|
| `.show-xs` / `.hide-xs` | < 480px | Extra small mobile |
| `.show-sm` / `.hide-sm` | 480px+ | Small mobile and up |
| `.show-md` / `.hide-md` | 640px+ | Small tablets and up |
| `.show-lg` / `.hide-lg` | 768px+ | Tablets and up |
| `.show-xl` / `.hide-xl` | 1024px+ | Small desktops and up |
| `.show-2xl` / `.hide-2xl` | 1280px+ | Standard desktops and up |
| `.show-3xl` / `.hide-3xl` | 1536px+ | Large desktops and up |
| `.show-4xl` / `.hide-4xl` | 1920px+ | Ultra-wide and up |
| `.show-mobile` / `.hide-mobile` | Legacy (< 768px) | Mobile devices |
| `.show-tablet` / `.hide-tablet` | Legacy (768-1023px) | Tablets |
| `.show-desktop` / `.hide-desktop` | Legacy (1024px+) | Desktops |

## ⚡ Performance Considerations

### 1. Touch Optimization
- Use `-webkit-overflow-scrolling: touch;` for smooth scrolling
- Implement passive event listeners
- Avoid hover-only interactions

### 2. Loading States
- Show skeleton screens on mobile
- Implement infinite scroll carefully
- Lazy load images below fold

### 3. Gestures
- Support swipe gestures where appropriate
- Pull-to-refresh for lists
- Swipe-to-delete for items

## 📱 Mobile-Specific Features

### Bottom Sheet (Modal Alternative)
On mobile, consider bottom sheets instead of center modals:

```css
.bottom-sheet {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  border-radius: 16px 16px 0 0;
  transform: translateY(100%);
  transition: transform 300ms ease-out;
}

.bottom-sheet.active {
  transform: translateY(0);
}
```

### Safe Areas (for iOS notch)
```css
.mobile-header {
  padding-top: env(safe-area-inset-top);
}

.fab {
  bottom: calc(24px + env(safe-area-inset-bottom));
}
```

## ✅ Mobile Checklist

### Visual
- [ ] All text is readable at default zoom
- [ ] Touch targets are minimum 44x44px
- [ ] Adequate spacing between interactive elements
- [ ] Content doesn't overflow horizontally
- [ ] Images scale appropriately

### Navigation
- [ ] Hamburger menu opens/closes smoothly
- [ ] Drawer slides in from left
- [ ] Overlay closes drawer when tapped
- [ ] Active nav item is clearly indicated
- [ ] Back navigation is available

### Interaction
- [ ] All buttons/links have active states
- [ ] Forms are easy to fill on mobile
- [ ] Input fields auto-focus appropriately
- [ ] Keyboard pushes content up (not covered)
- [ ] No hover-only interactions

### Performance
- [ ] Smooth 60fps animations
- [ ] Fast initial load time
- [ ] Responsive to touch immediately
- [ ] No layout shifts during load

### Accessibility
- [ ] Sufficient color contrast
- [ ] Focus indicators visible
- [ ] Screen reader friendly
- [ ] Supports text scaling
- [ ] Keyboard accessible (for tablet)

## 🚀 Implementation in Vue.js

### Breakpoint Composable
```typescript
// composables/useBreakpoint.ts
import { ref, onMounted, onUnmounted } from 'vue'

export function useBreakpoint() {
  const isMobile = ref(window.innerWidth < 768)
  const isTablet = ref(window.innerWidth >= 768 && window.innerWidth < 1024)
  const isDesktop = ref(window.innerWidth >= 1024)

  const updateBreakpoint = () => {
    const width = window.innerWidth
    isMobile.value = width < 768
    isTablet.value = width >= 768 && width < 1024
    isDesktop.value = width >= 1024
  }

  onMounted(() => {
    window.addEventListener('resize', updateBreakpoint)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', updateBreakpoint)
  })

  return { isMobile, isTablet, isDesktop }
}
```

### Usage in Components
```vue
<script setup lang="ts">
import { useBreakpoint } from '@/composables/useBreakpoint'

const { isMobile, isDesktop } = useBreakpoint()
</script>

<template>
  <nav v-if="isDesktop" class="sidebar">...</nav>
  <nav v-else class="mobile-drawer">...</nav>
</template>
```

## 📚 Resources

- [Material Design - Touch Targets](https://material.io/design/usability/accessibility.html#layout-and-typography)
- [iOS Human Interface Guidelines - Layout](https://developer.apple.com/design/human-interface-guidelines/ios/visual-design/adaptivity-and-layout/)
- [Web.dev - Responsive Design](https://web.dev/responsive-web-design-basics/)
- [MDN - Mobile Web Best Practices](https://developer.mozilla.org/en-US/docs/Web/Guide/Mobile)

---

**Last Updated**: 2026-01-31  
**Phase**: 8.9 (Notion-Inspired UI/UX Redesign - Mobile Responsive)
