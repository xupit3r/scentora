# UI/UX Specification

**Last Updated**: January 31, 2026  
**Version**: 2.0 (Accord System)

---

## Overview

This document specifies the user interface design, interactions, and user experience patterns for the Scentora accord management system.

---

## Design Principles

1. **Clarity**: Information is easy to find and understand
2. **Efficiency**: Common tasks require minimal clicks
3. **Consistency**: Patterns repeat across the application
4. **Feedback**: Users always know what's happening
5. **Flexibility**: Supports different workflows and preferences

---

## Color Palette

### Primary Colors

**Pyramid Position Colors**:
- **Top Notes**: Yellow gradient
  - Light: `#FFD93D`
  - Dark: `#FFA800`
  - Used for: Top position badges, top section highlights
  
- **Middle Notes**: Purple gradient
  - Light: `#B565D8`
  - Dark: `#8B5CF6`
  - Used for: Middle position badges, middle section highlights
  
- **Base Notes**: Brown gradient
  - Light: `#A0826D`
  - Dark: `#6B4423`
  - Used for: Base position badges, base section highlights

**Accent Colors**:
- **Primary Action**: Teal `#14B8A6` - Buttons, links, highlights
- **Success**: Green `#10B981` - Success messages, confirmations
- **Warning**: Orange `#F59E0B` - Low stock warnings
- **Error**: Red `#EF4444` - Error messages, critical alerts
- **Info**: Blue `#3B82F6` - Info messages, tooltips

### Neutral Colors
- **Background**: `#F5F5F5` - Page background
- **Surface**: `#FFFFFF` - Cards, modals, panels
- **Border**: `#E5E7EB` - Card borders, dividers
- **Text Primary**: `#111827` - Headings, important text
- **Text Secondary**: `#6B7280` - Body text, labels
- **Text Tertiary**: `#9CA3AF` - Placeholder, disabled

---

## Typography

### Font Family
- **Primary**: Inter, system-ui, sans-serif
- **Monospace**: 'Courier New', monospace (for volume numbers)

### Font Sizes
- **Display**: 48px - Page titles
- **H1**: 32px - Section headers
- **H2**: 24px - Card headers
- **H3**: 20px - Subsection headers
- **Body**: 16px - Regular text
- **Small**: 14px - Secondary text
- **Tiny**: 12px - Labels, captions

### Font Weights
- **Bold**: 700 - Headers, emphasis
- **Semibold**: 600 - Subheaders
- **Medium**: 500 - Buttons, labels
- **Regular**: 400 - Body text

---

## Layout

### Grid System
- **Desktop**: 12-column grid, max-width 1280px
- **Tablet**: 8-column grid, max-width 768px
- **Mobile**: 4-column grid, full width

### Spacing Scale
- `xs`: 4px
- `sm`: 8px
- `md`: 16px
- `lg`: 24px
- `xl`: 32px
- `2xl`: 48px
- `3xl`: 64px

### Container
- Max width: 1280px
- Horizontal padding: 24px (desktop), 16px (mobile)
- Centered on page

---

## Components

### AccordCard

**Purpose**: Display accord in grid view with essential information.

**Desktop Layout** (280x320px):
```
┌─────────────────────────────────┐
│ [Top Badge]          [⋮ Menu]   │
│                                  │
│   Citrus Fresh Accord            │ ← Name (H3)
│                                  │
│   25.5 ml · 510 drops            │ ← Volume (monospace)
│   Perfumer's Apprentice          │ ← Supplier (small text)
│                                  │
│   [fresh] [citrus] [energetic]   │ ← Tags (chips)
│   [summer]                       │
│                                  │
│   Last updated: Jan 20           │ ← Timestamp (tiny)
└─────────────────────────────────┘
```

**Hover State**:
- Subtle shadow elevation
- Border color changes to accent color
- Quick action buttons appear (Edit, View, Delete)

**Low Stock Indicator**:
- Orange badge (< 5ml): "Low Stock"
- Red badge (< 1ml): "Critical"
- Displayed in top-right corner

**Mobile Layout**: Full width, stacked content

---

### AccordForm

**Purpose**: Create or edit an accord with all properties.

**Modal Dialog** (600px wide):

**Tabs**:
1. **Basic Info** (default tab)
   - Name (text input)
   - Pyramid Position (dropdown)
   - Volume in ML (number input with increment buttons)
   - Volume in Drops (calculated, read-only, gray)

2. **Inventory**
   - Supplier (text input with autocomplete)
   - Purchase Date (date picker)
   - Dilution Percentage (number input, 0-100, with % symbol)

3. **Tags**
   - Tag selector component (see TagSelector)
   - Selected tags displayed as removable chips

4. **Notes**
   - Textarea (200px height, resizable)
   - Character count (optional)

**Buttons** (bottom right):
- Cancel (secondary)
- Save (primary, teal)

**Validation**:
- Real-time validation on blur
- Error messages below fields
- Required fields marked with *
- Submit button disabled until valid

**Mobile**: Full screen modal with sticky header

---

### AccordDetail

**Purpose**: Full-page view of a single accord.

**Desktop Layout**:

```
┌─────────────────────────────────────────────────────┐
│ [← Back]                          [Edit] [Delete]    │
├─────────────────────────────────────────────────────┤
│                                                      │
│   Citrus Fresh Accord [Top Badge]                   │ ← H1
│                                                      │
│   ┌──────────────┬──────────────┐                  │
│   │ Volume       │ Supplier     │                  │
│   │ 25.5 ml      │ Perfumer's   │                  │
│   │ 510 drops    │ Apprentice   │                  │
│   ├──────────────┼──────────────┤                  │
│   │ Purchased    │ Dilution     │                  │
│   │ Dec 15, 2025 │ 10%          │                  │
│   └──────────────┴──────────────┘                  │
│                                                      │
│   Tags                                              │
│   [fresh] [citrus] [energetic] [summer] [+Add]     │
│                                                      │
│   Notes                                             │
│   Very bright and zesty. Works well in summer       │
│   blends. Strong projection.                        │
│                                                      │
│   Created: Dec 15, 2025 · Updated: Jan 20, 2026     │
│                                                      │
└─────────────────────────────────────────────────────┘
```

**Actions**:
- Back button: Returns to collection
- Edit button: Opens AccordForm in edit mode
- Delete button: Shows confirmation dialog
- Add tag: Opens tag selector inline

---

### AccordFilters

**Purpose**: Filter and search the accord collection.

**Desktop** (Left sidebar, 280px wide):

```
┌─────────────────────────────┐
│ Filters                     │
│                             │
│ [Search input]              │
│                             │
│ Pyramid Position            │
│ ☑ Top Notes (15)            │
│ ☑ Middle Notes (18)         │
│ ☑ Base Notes (9)            │
│                             │
│ Volume Range                │
│ [━━━●═══════●═] 0-500 ml   │
│   5ml        250ml          │
│                             │
│ Tags                        │
│ [Tag multiselect ▼]         │
│ Selected: [fresh] [citrus]  │
│                             │
│ Supplier                    │
│ [Dropdown ▼]                │
│                             │
│ ☐ Low stock only (< 5ml)    │
│                             │
│ [Clear All Filters]         │
│                             │
│ 42 accords found            │
└─────────────────────────────┘
```

**Mobile**: Bottom sheet or slide-in panel
- Floating filter button (bottom right)
- Tapping opens full-screen filter panel
- Apply/Cancel buttons at bottom

**Filter Chips** (above grid on mobile):
- Show active filters as removable chips
- Example: `[Position: Top ×] [Tag: fresh ×]`

---

### TagSelector

**Purpose**: Select tags with autocomplete and grouping.

**Dropdown Interface**:

```
┌─────────────────────────────────────┐
│ Search tags...                  [×] │
├─────────────────────────────────────┤
│ CHARACTER                           │
│   fresh                             │
│   warm                              │
│   cool                              │
├─────────────────────────────────────┤
│ MOOD                                │
│   energetic                         │
│   calming                           │
├─────────────────────────────────────┤
│ YOUR CUSTOM TAGS                    │
│   my-favorite                       │
│   needs-testing                     │
├─────────────────────────────────────┤
│ ✚ Create "new-tag"                  │
└─────────────────────────────────────┘
```

**Features**:
- Fuzzy search filters all categories
- Keyboard navigation (arrow keys, enter)
- Grouped by category (collapsible)
- Custom tags in separate section
- "Create new tag" option at bottom
- Selected tags not shown in dropdown

**Selected Tags Display**:
```
[fresh ×] [citrus ×] [energetic ×] [my-favorite ×]
```
- Colored chips with category colors
- Remove button (×) on each
- Wrap to multiple lines

---

### Confirmation Dialog

**Purpose**: Confirm destructive actions (delete).

**Modal** (400px wide):

```
┌─────────────────────────────────────┐
│ Delete Accord?                      │
├─────────────────────────────────────┤
│                                     │
│ Are you sure you want to delete     │
│ "Citrus Fresh Accord"?              │
│                                     │
│ This action cannot be undone.       │
│                                     │
│         [Cancel]  [Delete]          │
└─────────────────────────────────────┘
```

**Delete Button**: Red, requires 2-second hold on mobile

---

## Views

### Home (Accord Collection)

**URL**: `/`

**Layout**:

```
┌──────────────────────────────────────────────────┐
│ [Logo] Accord Inventory         [👤 User Menu]   │
├──────────────────────────────────────────────────┤
│                                                   │
│   Accord Inventory                                │ ← H1
│   [+ New Accord]                   [Sort: ▼]     │
│                                                   │
│   ┌─────────────┐  ┌───────────────┬───────────┐│
│   │             │  │               │           ││
│   │  Filters    │  │  AccordCard   │ AccordCard││
│   │             │  │               │           ││
│   │  (sidebar)  │  ├───────────────┼───────────┤│
│   │             │  │  AccordCard   │ AccordCard││
│   │             │  │               │           ││
│   └─────────────┘  └───────────────┴───────────┘│
│                                                   │
└──────────────────────────────────────────────────┘
```

**Features**:
- 3-4 column grid (desktop)
- Left sidebar filters
- Sort dropdown (top right)
- "New Accord" button (top, primary action)
- Empty state when no accords
- Pagination if > 50 accords

**Empty State**:
```
     🏺
     
  No accords yet
  
  Add your first accord to
  start building your collection
  
  [+ Add Accord]
```

---

### Accord Detail View

**URL**: `/accord/:id`

**Layout**: See AccordDetail component above

**Navigation**:
- Back button (top left) → returns to Home
- Breadcrumb: Home > Accord Name

---

### Statistics View

**URL**: `/statistics`

**Layout**:

```
┌──────────────────────────────────────────────────┐
│ [Logo] Statistics                [👤 User Menu]   │
├──────────────────────────────────────────────────┤
│                                                   │
│   Collection Statistics                           │ ← H1
│                                                   │
│   ┌───────────┐ ┌───────────┐ ┌───────────┐    │
│   │ 42        │ │ 1,250.5   │ │ Top: 15   │    │
│   │ Accords   │ │ Total ml  │ │ Mid: 18   │    │
│   │           │ │           │ │ Base: 9   │    │
│   └───────────┘ └───────────┘ └───────────┘    │
│                                                   │
│   Volume by Position                              │
│   ┌─────────────────────────────────────────┐   │
│   │ Top    ████████████ 420.5 ml            │   │
│   │ Middle ████████████████ 550.0 ml        │   │
│   │ Base   ████████ 280.0 ml                │   │
│   └─────────────────────────────────────────┘   │
│                                                   │
│   Most Used Tags                                  │
│   [fresh: 12] [warm: 8] [floral: 7] [citrus: 6] │
│                                                   │
│   Low Stock Alerts (< 5ml)                       │
│   • Vanilla Base - 2.5 ml                        │
│   • Rose Absolute - 1.2 ml                       │
│                                                   │
└──────────────────────────────────────────────────┘
```

**Not Prominent**: Statistics are accessible but not the main focus.

---

## Interactions

### Adding an Accord

**Flow**:
1. Click "New Accord" button → Modal opens
2. Fill Basic Info tab (required fields)
3. Switch to Inventory tab (optional)
4. Switch to Tags tab, add tags
5. Switch to Notes tab, add description
6. Click Save → Accord created, modal closes
7. Success toast: "Accord created successfully"
8. New accord appears in grid

**Keyboard Shortcut**: `Ctrl/Cmd + N`

---

### Editing an Accord

**Flow**:
1. Hover over AccordCard → Edit button appears
2. Click Edit → Modal opens with pre-filled data
3. Modify fields
4. Click Save → Accord updated
5. Success toast: "Accord updated"

**Alternative**: Click card → Navigate to detail view → Click Edit button

---

### Deleting an Accord

**Flow**:
1. Hover over AccordCard → Delete button appears (or menu ⋮)
2. Click Delete → Confirmation dialog appears
3. Confirm deletion → Accord deleted
4. Success toast: "Accord deleted"
5. Accord removed from grid

**Mobile**: Swipe left on card to reveal Delete action

---

### Filtering Accords

**Flow**:
1. Open filter sidebar (always visible on desktop)
2. Check pyramid position checkboxes
3. Adjust volume slider
4. Select tags from multiselect
5. Results update in real-time
6. Filter chips appear above grid
7. Click "Clear All" to reset

**Mobile**: Tap filter button → Bottom sheet opens → Apply filters

---

### Searching Accords

**Flow**:
1. Type in search box (debounced, 300ms)
2. Results filter in real-time
3. Searches name, notes, supplier
4. Highlights matching text (optional)

---

### Adding Tags

**Flow**:
1. In AccordForm or AccordDetail, click tag input
2. Dropdown opens with grouped tags
3. Type to filter, or scroll categories
4. Click tag to add → Chip appears
5. Click "Create new" to add custom tag
6. Save form to persist tags

**Keyboard**: Arrow keys to navigate, Enter to select, Escape to close

---

### Volume Conversion

**Automatic**:
- When user enters ML, drops calculated automatically
- Display: `25.5 ml · 510 drops`
- Formula: drops = ml × 20

**Manual Entry** (future):
- Toggle between ML and Drops input
- Convert on save

---

## Responsive Breakpoints

### Desktop (1024px+)
- 3-4 column grid
- Left sidebar filters (always visible)
- Full-width modals (600px)

### Tablet (768px - 1023px)
- 2-3 column grid
- Collapsible filter sidebar
- Full-width modals

### Mobile (< 768px)
- Single column grid
- Bottom sheet filters
- Full-screen modals
- Swipe gestures (delete, edit)

---

## Animations & Transitions

### Micro-interactions
- **Button Hover**: Scale 1.05, 200ms ease
- **Card Hover**: Shadow elevation, 300ms ease
- **Modal Open**: Fade in + scale from 0.95, 200ms ease-out
- **Modal Close**: Fade out + scale to 0.95, 150ms ease-in
- **Toast**: Slide in from top, 300ms ease-out
- **Filter Panel**: Slide in from left/bottom, 250ms ease

### Loading States
- **Skeleton Screen**: Pulsing gray rectangles while loading
- **Spinner**: Rotating circle for actions (saves, deletes)

---

## Error Handling

### Form Validation Errors
- Display below field in red
- Icon: ⚠️
- Example: "Name is required"

### API Errors
- Toast notification (red background)
- Icon: ❌
- Auto-dismiss after 5 seconds
- Example: "Failed to save accord. Please try again."

### Empty States
- Friendly illustration
- Clear message
- Call-to-action button
- Example: "No accords found. Try adjusting your filters."

---

## Accessibility

### Keyboard Navigation
- All interactive elements focusable
- Visible focus indicators (outline)
- Tab order logical (top to bottom, left to right)
- Escape closes modals and dropdowns
- Enter submits forms

### Screen Readers
- Semantic HTML (header, nav, main, article)
- ARIA labels for icon buttons
- ARIA live regions for dynamic content
- Alt text for images

### Color Contrast
- Text meets WCAG AA standards (4.5:1)
- Colorblind-friendly palette
- Don't rely on color alone (use icons too)

---

## Mobile Optimizations

### Touch Targets
- Minimum 44×44 px
- Spacing between targets: 8px

### Gestures
- Swipe left: Delete accord (with undo toast)
- Pull to refresh: Reload collection
- Pinch to zoom: Not needed (fixed layout)

### Performance
- Lazy load images
- Virtual scrolling for large lists
- Debounced search input

---

## Dark Mode (Future)

Planned but not implemented in Phase 8.

**Color Adjustments**:
- Background: `#1F2937` (dark gray)
- Surface: `#374151` (lighter gray)
- Text: `#F9FAFB` (off-white)
- Borders: `#4B5563`

---

## Browser Support

- Chrome 90+ ✅
- Firefox 88+ ✅
- Safari 14+ ✅
- Edge 90+ ✅

---

## Notes

- All measurements in pixels (px)
- Use CSS Grid and Flexbox for layouts
- Tailwind CSS or custom CSS
- Component library: Consider Headless UI or Radix UI
- Icons: Heroicons or Lucide
- Date picker: Native HTML5 or vue-datepicker
- Toast notifications: vue-toastification
