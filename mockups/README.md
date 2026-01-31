# Scentora Design Mockups

This directory contains design mockups and wireframes for **Phase 8.9** - the Notion-inspired UI/UX redesign of Scentora's accord management system.

## 📂 Directory Structure

```
mockups/
├── index.html                  # Interactive mockup gallery (START HERE!)
├── assets/
│   └── design-system.css      # Design tokens and reusable styles
├── wireframes/                # Low-fidelity layout sketches
│   ├── 01-login.html
│   ├── 02-register.html
│   ├── 03-layout-sidebar.html
│   ├── 04-accord-grid.html
│   ├── 05-accord-detail.html
│   ├── 06-create-modal.html
│   ├── 07-empty-state.html
│   └── 08-statistics.html
└── high-fidelity/            # Polished, fully-styled mockups
    ├── 01-login.html
    ├── 02-register.html
    ├── 03-main-layout.html
    ├── 04-accord-cards.html
    └── 05-empty-state.html
```

## 🚀 Quick Start

### View Mockups in Browser

1. **Option 1: Start from the gallery**
   ```bash
   # Open the main gallery page
   open mockups/index.html
   ```

2. **Option 2: View individual mockups**
   ```bash
   # View a specific mockup
   open mockups/wireframes/01-login.html
   open mockups/high-fidelity/03-main-layout.html
   ```

3. **Option 3: Use a local server** (recommended for best experience)
   ```bash
   # Using Python
   cd mockups
   python3 -m http.server 8080
   # Then visit: http://localhost:8080

   # Using Node.js
   cd mockups
   npx serve
   ```

## 📐 Wireframes vs. High-Fidelity

### Wireframes (Low-Fidelity)
- **Purpose**: Show layout structure and content hierarchy
- **Style**: Simple, minimal styling with placeholders
- **Use case**: Early-stage planning and layout validation
- **Located in**: `mockups/wireframes/`

### High-Fidelity Mockups
- **Purpose**: Demonstrate final look and feel
- **Style**: Full design system with colors, typography, and interactions
- **Use case**: Design approval and developer handoff
- **Located in**: `mockups/high-fidelity/`

## 🎨 Design System

The mockups implement the **Notion-inspired design system** defined in Phase 8.9 planning:

### Color Palette

| Purpose | Color | Hex |
|---------|-------|-----|
| Primary Text | Dark gray | `#37352F` |
| Secondary Text | Medium gray | `#787774` |
| Accent | Teal | `#0F766E` |
| Top Note | Amber | `#FEF3C7` / `#92400E` |
| Middle Note | Purple | `#E9D5FF` / `#6B21A8` |
| Base Note | Brown | `#F5E6D3` / `#78350F` |

### Typography

- **Font Stack**: System fonts (-apple-system, BlinkMacSystemFont, 'Segoe UI', etc.)
- **Sizes**: 12, 14, 16, 20, 24, 32, 40px
- **Weights**: 400 (normal), 500 (medium), 600 (semibold)

### Spacing

- **Grid**: 8px base unit
- **Scale**: 4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 80px

### Shadows & Borders

- **Shadows**: sm, md, lg, xl (subtle, progressive elevation)
- **Border Radius**: 3px (sm), 6px (md), 12px (lg)
- **Transitions**: 150ms (fast), 250ms (normal), 350ms (slow)

## 📄 Mockup Inventory

### Authentication
- ✅ Login page
- ✅ Register page (with invitation code)

### Main Application
- ✅ Main layout with sidebar navigation
- ✅ Accord inventory grid
- ✅ Accord card components
- ✅ Accord detail view
- ✅ Create accord modal
- ✅ Empty states
- ✅ Statistics dashboard

### Layout States
- ✅ Sidebar (expanded)
- ⏳ Sidebar (collapsed) - planned
- ⏳ Mobile view - planned

## 🛠️ Design System File

All mockups reference a shared design system CSS file:

```
mockups/assets/design-system.css
```

This file contains:
- CSS custom properties (design tokens)
- Base styles and resets
- Reusable component styles (buttons, inputs, cards, etc.)
- Utility classes
- Notion-inspired design patterns

## 📝 Implementation Notes

### For Designers
- All mockups are HTML/CSS and can be opened directly in a browser
- Colors, spacing, and typography match the Phase 8.9 specifications
- Use these as reference for creating production components

### For Developers
- Component patterns demonstrated in mockups should be replicated in Vue.js
- Design tokens in `design-system.css` should be ported to:
  - `frontend/src/design-system/tokens.ts` (TypeScript)
  - Tailwind config (if using Tailwind)
  - Naive UI theme customization
- Hover effects and transitions are shown for interaction reference

## 🔗 Related Documentation

- **Phase 8.9 Plan**: `/home/joe/code/scentora/PLAN.md` (Phase 8.9 section)
- **UI/UX Specifications**: `/home/joe/code/scentora/specs/ui-ux-spec.md`
- **Phase 8.9 Summary**: `/home/joe/code/scentora/PHASE8_9_PLAN_SUMMARY.md`
- **Main README**: `/home/joe/code/scentora/README.md`

## ✨ Key Features Demonstrated

### Sidebar Navigation
- Fixed width (280px), collapsible to icon-only (64px)
- Section-based organization
- Active state highlighting
- User profile in footer

### Card Design
- Clean with subtle shadows
- Left border indicates position (top/middle/base)
- Hover effects reveal actions
- Position-specific color coding (subtle)

### Empty States
- Welcoming and informative
- Clear call-to-action
- Feature highlights
- Beautiful iconography

### Forms & Modals
- Clear labeling and validation
- Multi-step form support
- Tag selector component
- Inline editing patterns

## 🚦 Status

- ✅ **Wireframes**: 8 pages complete
- ✅ **High-Fidelity**: 5 core pages complete
- ✅ **Design System**: CSS file complete
- ✅ **Gallery**: Interactive index page complete

## 📞 Feedback

These mockups are part of the Phase 8.9 planning process. For questions or feedback:
1. Review the mockups in the gallery (`index.html`)
2. Compare against Phase 8.9 specifications
3. Provide feedback on layout, styling, and interactions
4. Suggest additional pages or components needed

---

**Last Updated**: 2026-01-31  
**Phase**: 8.9 (Notion-Inspired UI/UX Redesign)  
**Repository**: https://github.com/xupit3r/scentora
