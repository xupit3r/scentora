# Scentora Design System

Notion-inspired design system for Scentora built with **Naive UI** and **Tailwind CSS**.

## Structure

```
design-system/
├── tokens.ts        # Design tokens (colors, typography, spacing, etc.)
├── theme.ts         # Naive UI theme overrides
└── README.md        # This file
```

## Design Tokens

All design tokens are exported from `tokens.ts` and can be imported throughout the application:

```typescript
import { colors, typography, spacing, shadows } from '@/design-system/tokens';
```

### Colors

**Text Colors:**
- `colors.text.primary` - `#37352F` - Main text
- `colors.text.secondary` - `#787774` - Secondary text
- `colors.text.tertiary` - `#9B9A97` - Tertiary/muted text

**Background Colors:**
- `colors.bg.primary` - `#FFFFFF` - Main background
- `colors.bg.secondary` - `#FAFAFA` - Secondary background
- `colors.bg.tertiary` - `#F7F6F3` - Tertiary background

**Accent Colors:**
- `colors.accent.primary` - `#0F766E` - Primary action color (teal)
- `colors.accent.hover` - `#0D6460` - Hover state

**Position Colors (Refined):**
- Top Note: `#FEF3C7` background, `#92400E` text
- Middle Note: `#E9D5FF` background, `#6B21A8` text
- Base Note: `#F5E6D3` background, `#78350F` text

### Typography

**Font Family:**
- Sans-serif: `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, ...`
- Monospace: `"SF Mono", Monaco, "Cascadia Code", ...`

**Font Sizes:**
- `xs`: 12px
- `sm`: 14px
- `base`: 16px
- `lg`: 20px
- `xl`: 24px
- `2xl`: 32px
- `3xl`: 40px

**Font Weights:**
- `normal`: 400
- `medium`: 500
- `semibold`: 600

### Spacing (8px Grid)

- `1`: 4px
- `2`: 8px
- `3`: 12px
- `4`: 16px
- `5`: 20px
- `6`: 24px
- `8`: 32px
- `10`: 40px
- `12`: 48px
- `16`: 64px
- `20`: 80px

### Border Radius

- `sm`: 4px
- `base`: 8px
- `md`: 12px
- `lg`: 16px
- `full`: 9999px

### Shadows

- `sm`: Light shadow for subtle elevation
- `base`: Default card shadow
- `md`: Medium elevation
- `lg`: Large elevation for modals
- `hover`: Hover state shadow

## Naive UI Theme

The Naive UI theme is configured in `theme.ts` and applied globally via `NConfigProvider` in `App.vue`.

**Usage:**
```vue
<script setup>
import { naiveTheme } from '@/design-system/theme';
</script>

<template>
  <n-config-provider :theme-overrides="naiveTheme">
    <!-- Your app -->
  </n-config-provider>
</template>
```

## Tailwind CSS

Tailwind is configured with custom tokens matching the design system in `tailwind.config.js`.

**Custom Classes:**
```html
<!-- Text colors -->
<div class="text-text-primary">Primary text</div>
<div class="text-text-secondary">Secondary text</div>

<!-- Backgrounds -->
<div class="bg-bg-secondary">Secondary background</div>

<!-- Accent -->
<button class="bg-accent hover:bg-accent-hover">Action</button>

<!-- Position indicators -->
<div class="bg-position-top-bg text-position-top-text">Top Note</div>
<div class="bg-position-middle-bg text-position-middle-text">Middle Note</div>
<div class="bg-position-base-bg text-position-base-text">Base Note</div>

<!-- Spacing -->
<div class="p-6 gap-4">Content with 24px padding and 16px gap</div>

<!-- Shadows -->
<div class="shadow-md hover:shadow-hover">Card</div>
```

## Global Styles

Global styles are defined in `src/style.css` and include:

- **Base styles**: Reset, typography, links
- **Component classes**: `.card-hover`, `.position-top`, `.empty-state`
- **Utility classes**: `.truncate-2`, `.transition-smooth`
- **Animations**: `fadeIn`, `slideInRight`, `pulse`

## Components

### Using Naive UI Components

Import and use Naive UI components directly:

```vue
<script setup>
import { NButton, NCard, NInput } from 'naive-ui';
</script>

<template>
  <n-card>
    <n-input placeholder="Type here..." />
    <n-button type="primary">Submit</n-button>
  </n-card>
</template>
```

Available components are auto-registered in `main.ts`.

## Design Principles

1. **Clean & Minimal**: Less is more, focus on content
2. **Consistent Spacing**: 8px grid for all spacing
3. **Subtle Colors**: Muted, professional color palette
4. **Clear Hierarchy**: Typography scale for visual hierarchy
5. **Smooth Interactions**: 200ms transitions for all interactions
6. **Purposeful Shadows**: Use shadows sparingly for elevation

## Best Practices

1. **Use design tokens**: Always use tokens instead of hard-coded values
2. **Follow spacing grid**: Use multiples of 8px (4, 8, 12, 16, 24, 32...)
3. **Consistent typography**: Use defined font sizes and weights
4. **Semantic colors**: Use semantic colors (success, error) for feedback
5. **Accessible**: Maintain WCAG AA contrast ratios minimum

## Examples

### Card Component
```vue
<n-card class="card-hover" :bordered="true">
  <h3 class="text-lg font-semibold mb-4">Card Title</h3>
  <p class="text-text-secondary">Card content goes here.</p>
</n-card>
```

### Button Group
```vue
<n-space>
  <n-button type="primary">Primary Action</n-button>
  <n-button>Secondary Action</n-button>
</n-space>
```

### Form
```vue
<n-form>
  <n-form-item label="Name">
    <n-input placeholder="Enter name..." />
  </n-form-item>
  <n-form-item label="Notes">
    <n-input type="textarea" placeholder="Enter notes..." />
  </n-form-item>
</n-form>
```

## Resources

- [Naive UI Documentation](https://www.naiveui.com/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
- [Notion Design Inspiration](https://www.notion.so/)
