# Icon System

Simple, flat SVG icons for Scentora mockups and implementation.

## Icons Available

### Navigation Icons
- `home.svg` - Home/Dashboard
- `inventory.svg` - Accord Inventory
- `formulations.svg` - Formulations/Flask
- `stats.svg` - Statistics/Analytics
- `settings.svg` - Settings/Preferences
- `users.svg` / `invitations.svg` - Users/Invitations

### UI Action Icons
- `search.svg` - Search
- `bell.svg` - Notifications
- `user.svg` - User profile
- `edit.svg` - Edit/Pencil
- `trash.svg` - Delete/Trash
- `more-vertical.svg` - More options (vertical dots)
- `plus.svg` - Add/Create
- `x.svg` - Close/Remove
- `menu.svg` - Menu/Hamburger

### Feature Icons
- `bottle.svg` - Perfume bottle/Accord
- `tag.svg` - Tag/Label
- `sparkles.svg` - New/Featured
- `book.svg` - Documentation/Notes
- `info.svg` - Information
- `file.svg` - File/Document

### UI Navigation
- `check.svg` / `check-simple.svg` - Checkmark/Complete
- `chevron-left.svg` / `chevron-right.svg` - Navigation arrows

## Usage

### In HTML

**Method 1: Inline SVG (recommended for mockups)**
```html
<span class="icon icon-md">
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
    <polyline points="9 22 9 12 15 12 15 22"/>
  </svg>
</span>
```

**Method 2: Via img tag**
```html
<img src="../assets/icons/home.svg" alt="Home" class="icon icon-md">
```

**Method 3: Via CSS background (not recommended - loses color control)**
```css
.icon-home {
  background-image: url('../assets/icons/home.svg');
}
```

### Sizes

Use the size classes to control icon dimensions:

```html
<span class="icon icon-sm">...</span>  <!-- 16x16px -->
<span class="icon icon-md">...</span>  <!-- 20x20px - default -->
<span class="icon icon-lg">...</span>  <!-- 24x24px -->
<span class="icon icon-xl">...</span>  <!-- 32x32px -->
```

### Colors

Icons inherit `currentColor` from their parent by default. Use color classes:

```html
<span class="icon icon-primary">...</span>   <!-- Primary text color -->
<span class="icon icon-secondary">...</span> <!-- Secondary text color -->
<span class="icon icon-accent">...</span>    <!-- Accent color -->
<span class="icon icon-white">...</span>     <!-- White -->
```

Or apply color directly:
```html
<span class="icon" style="color: #0F766E;">...</span>
```

## Implementation Notes

### For Production (Vue.js)

When implementing in Vue.js, consider:

1. **Component-based approach:**
```vue
<template>
  <span class="icon" :class="sizeClass">
    <svg v-html="iconSvg" />
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import homeIcon from '@/assets/icons/home.svg?raw'

const props = defineProps<{
  name: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
}>()

const iconSvg = computed(() => {
  // Icon mapping logic
  return homeIcon
})

const sizeClass = computed(() => `icon-${props.size || 'md'}`)
</script>
```

2. **Icon library:**
   - Use `vite-svg-loader` to import SVGs as Vue components
   - Or use `unplugin-icons` for automatic icon loading

3. **Recommended packages:**
   - `@iconify/vue` - Large icon collection
   - `unplugin-icons` - Auto-import icons
   - Or create custom icon component wrapper

### Design Principles

1. **Consistent stroke width**: All icons use `stroke-width="2"`
2. **24x24 viewBox**: Standard viewBox for easy scaling
3. **Stroke-based**: Icons use strokes, not fills (cleaner, scalable)
4. **Round caps**: `stroke-linecap="round"` for softer look
5. **Minimal**: Simple shapes, no unnecessary details

### Customization

To modify icons:
1. Edit SVG files directly
2. Maintain 24x24 viewBox
3. Keep stroke-width consistent
4. Test at different sizes (16px, 20px, 24px, 32px)

## Icon Sources

Icons are custom-designed for Scentora, inspired by:
- Lucide Icons (https://lucide.dev) - Clean, minimal style
- Feather Icons (https://feathericons.com) - Stroke-based simplicity
- Tailored to match Notion-inspired aesthetic

## Adding New Icons

To add new icons:

1. **Create SVG file** in `mockups/assets/icons/`
2. **Follow the format:**
   ```svg
   <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
     <!-- Icon paths here -->
   </svg>
   ```
3. **Test at multiple sizes**
4. **Document in this README**

## License

Icons are part of the Scentora project and follow the MIT license.
