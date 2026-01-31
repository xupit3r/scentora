# Mobile Responsive Design Guidelines

Guidelines for creating responsive, mobile-friendly interfaces in Scentora.

## 📱 Breakpoints

We follow a mobile-first approach with three main breakpoints:

| Device | Breakpoint | Strategy |
|--------|------------|----------|
| **Mobile** | 320px - 767px | Default styles, single column |
| **Tablet** | 768px - 1023px | 2-column layouts, expanded spacing |
| **Desktop** | 1024px+ | Multi-column, full sidebar |

```css
/* Mobile-first (default) */
.element { /* mobile styles */ }

/* Tablet and up */
@media (min-width: 768px) {
  .element { /* tablet styles */ }
}

/* Desktop and up */
@media (min-width: 1024px) {
  .element { /* desktop styles */ }
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

**Mobile:**
- Single column layouts
- Stacked cards
- Reduced padding (16px)
- Smaller font sizes for titles
- Hidden secondary information

**Desktop:**
- Multi-column grids (2-4 columns)
- Side-by-side layouts
- Generous padding (48-64px)
- Full information displayed

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

**Mobile:**
- Full width (single column)
- Compact padding (16px)
- Stacked information
- Essential info only

**Desktop:**
- Grid layout (2-4 columns)
- Standard padding (20px)
- Expanded information
- Hover effects

```css
.card-grid {
  display: grid;
  gap: var(--spacing-4);
  grid-template-columns: 1fr;
}

@media (min-width: 640px) {
  .card-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .card-grid {
    grid-template-columns: repeat(3, 1fr);
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

Reduce spacing on mobile to maximize content area:

```css
.content-wrapper {
  padding: var(--spacing-4); /* 16px on mobile */
}

@media (min-width: 768px) {
  .content-wrapper {
    padding: var(--spacing-8); /* 32px on tablet */
  }
}

@media (min-width: 1024px) {
  .content-wrapper {
    padding: var(--spacing-12); /* 48px on desktop */
  }
}
```

## 🔄 Responsive Utilities

Use utility classes for show/hide patterns:

```html
<!-- Show only on mobile -->
<div class="show-mobile hide-tablet">...</div>

<!-- Show on tablet and desktop -->
<div class="hide-mobile show-tablet">...</div>

<!-- Show only on desktop -->
<div class="hide-mobile hide-tablet show-desktop">...</div>
```

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
