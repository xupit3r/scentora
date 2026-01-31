# Notion-Inspired UI/UX Redesign Plan - Summary

**Date:** January 31, 2026  
**Status:** Planning Complete ✅  
**Phase:** 8.9  
**Timeline:** 1-2 weeks  

## Executive Summary

Created a comprehensive plan to transform Scentora's UI from its current functional state into a clean, minimalist interface inspired by Notion. The plan includes framework selection, design system specifications, and a 7-phase implementation roadmap.

## What Was Documented

### 1. Main PLAN.md Updates
Added Phase 8.9 section (~200 lines) covering:
- Framework selection rationale
- Design principles to implement
- 7 implementation sub-phases with detailed checklists
- Success criteria
- Migration strategy
- Risk mitigation

**Location:** `/home/joe/code/scentora/PLAN.md` (lines 705-900+)

### 2. UI/UX Specification Updates
Expanded `specs/ui-ux-spec.md` with (~400 lines):
- Complete Notion-inspired design system
- New color palette (text, backgrounds, borders, accents)
- Typography scale and font weights
- 8px spacing grid system
- Shadow and elevation definitions
- Border radius values
- Transition and animation specs
- Sidebar navigation structure
- Component redesign specifications
- Mobile responsive guidelines
- Accessibility requirements
- Performance targets

**Location:** `/home/joe/code/scentora/specs/ui-ux-spec.md` (lines 650-1100+)

### 3. README.md Updates
Updated project status and features section to reflect:
- Current phase (8.2-8.3 complete)
- Next phase (8.9 UI redesign)
- Updated tech stack (added Naive UI + Tailwind)
- New features list focused on accord management
- Planned UI/UX improvements

**Location:** `/home/joe/code/scentora/README.md`

### 4. Session Planning Document
Created detailed implementation plan in session workspace:
- 10-part comprehensive plan
- Framework comparison analysis
- Current state analysis
- Design specifications (colors, typography, spacing)
- Component mockups
- File structure
- Detailed timeline

**Location:** `~/.copilot/session-state/aee8b5eb-f1f6-4f0b-89e1-2a1452bb4424/plan.md`

## Framework Selection

### Evaluation Summary

Evaluated 6 Vue UI frameworks for Notion-like design:

| Framework | Score | Notes |
|-----------|-------|-------|
| **Naive UI** | 9/10 | ✅ RECOMMENDED - Clean, TypeScript-first, minimal |
| Headless UI | 10/10 | Maximum control, more work |
| PrimeVue | 6/10 | Requires heavy theming |
| Element Plus | 5/10 | Too opinionated |
| Vuetify | 4/10 | Material Design conflicts |
| Ant Design Vue | 5/10 | Different design language |

### Selected Stack

**Primary**: Naive UI  
**Styling**: Tailwind CSS  
**Approach**: Hybrid (framework for basics + custom for unique features)

**Why Naive UI?**
- Clean, minimalist design out of the box
- Excellent TypeScript support (written in TypeScript)
- Lightweight and performant
- Tree-shakeable components
- Customizable theme system
- Active development

## Design System Highlights

### Color Palette (Notion-inspired)

**Text**: #37352F (primary), #787774 (secondary), #9B9A97 (tertiary)  
**Backgrounds**: #FFFFFF, #FAFAFA, #F7F6F3  
**Borders**: #E9E9E7, #D9D9D7  
**Accent**: #0F766E (teal for actions)  

**Position Colors** (Refined - Subtle):
- Top: #FEF3C7 background, #92400E text
- Middle: #E9D5FF background, #6B21A8 text
- Base: #F5E6D3 background, #78350F text

### Typography

**Font Stack**: System fonts (-apple-system, BlinkMacSystemFont, Segoe UI...)  
**Scale**: 12, 14, 16, 20, 24, 32, 40px  
**Weights**: 400 (normal), 500 (medium), 600 (semibold)  
**Line Heights**: 1.2 (tight), 1.5 (normal), 1.7 (relaxed)

### Spacing (8px Grid)

4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 80px

**Guidelines**:
- Container padding: 48-64px
- Section gaps: 24-32px
- Related items: 8-16px
- Inline spacing: 4-8px

## Implementation Roadmap

### Phase 8.9.1: Foundation & Setup (2-3 days)
- Install Naive UI and Tailwind CSS
- Create design system (tokens.ts)
- Configure Tailwind with custom theme
- Set up global styles

### Phase 8.9.2: Navigation & Layout (2-3 days)
- Create AppSidebar.vue component
- Restructure App.vue for sidebar layout
- Implement responsive sidebar
- Add breadcrumb navigation

### Phase 8.9.3: Core Components Redesign (3-4 days)
- Redesign AccordCard.vue
- Redesign AccordForm.vue
- Redesign TagSelector.vue
- Redesign AccordFilters.vue

### Phase 8.9.4: Advanced Interactions (2-3 days)
- Implement inline editing
- Add keyboard shortcuts
- Standardize hover effects
- Implement skeleton screens

### Phase 8.9.5: Typography & Spacing Refinement (1-2 days)
- Define type scale
- Review spacing consistency
- Replace vibrant colors with subtle tones
- Apply 8px grid throughout

### Phase 8.9.6: Empty States & Feedback (1 day)
- Create beautiful empty states
- Implement toast notifications
- Add loading indicators
- Create form validation feedback

### Phase 8.9.7: Polish & Testing (2-3 days)
- Responsive testing (desktop, tablet, mobile)
- Accessibility audit
- Performance optimization
- Browser testing
- Lighthouse audit

**Total Timeline**: 10-15 working days

## Key Features to Implement

### Sidebar Navigation
- Fixed sidebar (280px wide)
- Collapsible to icon-only (64px)
- Mobile: Drawer from left
- Navigation items with icons
- User profile at bottom

### Inline Editing
- Click to edit accord name, volume
- Auto-save after 2s delay
- Optimistic UI updates
- Visual feedback

### Keyboard Shortcuts
- `/` - Focus search
- `N` - New accord
- `?` - Show shortcuts help
- `Esc` - Close modals
- `Cmd/Ctrl + K` - Command palette (future)

### Improved Cards
- Cleaner design with subtle shadows
- Hover-only actions
- Smooth lift effect
- Left border for position indicator

## Success Criteria

### Visual Quality ✅
- Clean, uncluttered interface
- Consistent spacing throughout
- Professional typography
- Subtle, purposeful colors
- Smooth animations (no jank)

### Usability ✅
- Intuitive navigation
- Fast, responsive interactions
- Clear feedback for actions
- Helpful empty states
- Keyboard accessible

### Technical ✅
- TypeScript type safety maintained
- Bundle size < 500KB (gzipped)
- Lighthouse score > 90
- No accessibility violations
- Works on all modern browsers

## Migration Strategy

**Approach**: Gradual Migration
1. Set up framework alongside existing components
2. Migrate App.vue first (navigation/layout)
3. Redesign one component at a time
4. Keep old components until new ones are ready
5. Remove old code once migration complete

**Risk Mitigation**:
- Feature flags for new UI
- Git branches for each phase
- Regular user testing
- Rollback plan if issues arise
- Performance monitoring

## File Structure Changes

```
frontend/src/
├── design-system/          # NEW
│   ├── tokens.ts
│   ├── theme.ts
│   └── tailwind.config.js
│
├── components/
│   ├── layout/             # NEW
│   │   ├── AppSidebar.vue
│   │   ├── AppHeader.vue
│   │   └── AppBreadcrumbs.vue
│   │
│   ├── accord/
│   │   ├── AccordCard.vue     # REDESIGN
│   │   ├── AccordForm.vue     # REDESIGN
│   │   ├── AccordGrid.vue     # NEW
│   │   └── AccordDetail.vue   # NEW
│   │
│   ├── ui/
│   │   ├── TagSelector.vue    # REDESIGN
│   │   ├── FilterPanel.vue    # REDESIGN
│   │   ├── EmptyState.vue     # NEW
│   │   └── SkeletonCard.vue   # NEW
│   │
│   └── common/                # NEW
│       ├── Button.vue
│       ├── Input.vue
│       └── Modal.vue
│
└── composables/               # NEW
    ├── useKeyboard.ts
    ├── useToast.ts
    └── useTheme.ts
```

## Next Steps

### Immediate Actions
1. ✅ Review plan with stakeholders
2. ⏳ Approve framework choice (Naive UI + Tailwind)
3. ⏳ Set up development environment
4. ✅ **Create design mockups** - COMPLETE! View at `/mockups/`
5. ⏳ Begin Phase 8.9.1 implementation

### Questions to Answer
- [x] Do we need design mockups first? **YES - Complete!** (HTML/SVG format)
- [ ] Should we do a proof-of-concept with 1-2 components first?
- [ ] Any specific Notion features to prioritize?
- [ ] Dark mode support needed?
- [ ] Accessibility standards (WCAG AA or AAA)?

## Design Mockups

✅ **Status**: Complete (2026-01-31)

Created comprehensive mockups for Phase 8.9 redesign:
- **8 wireframes** (low-fidelity layouts)
- **5 high-fidelity mockups** (polished designs)
- **Interactive gallery** at `mockups/index.html`
- **Design system CSS** with all tokens and components

**View mockups**: Open `mockups/index.html` in a browser or visit the mockups directory.

**Mockups include**:
- Authentication (login, register)
- Main layout with sidebar navigation
- Accord inventory grid and cards
- Accord detail views
- Create/edit modals
- Empty states
- Statistics dashboard

## Commit Details

**Commit**: `f13bb02`  
**Message**: "docs: Add Notion-inspired UI/UX redesign plan (Phase 8.9)"  
**Files Changed**:
- PLAN.md (+200 lines)
- specs/ui-ux-spec.md (+400 lines)
- README.md (updated status)

**Total**: 841 insertions, 16 deletions

## Resources

**Documentation Locations**:
- Main plan: `/home/joe/code/scentora/PLAN.md`
- UI/UX specs: `/home/joe/code/scentora/specs/ui-ux-spec.md`
- Detailed plan: `~/.copilot/session-state/*/plan.md`
- This summary: `/home/joe/code/scentora/PHASE8_9_PLAN_SUMMARY.md`

**External References**:
- Naive UI: https://www.naiveui.com/
- Tailwind CSS: https://tailwindcss.com/
- Notion: https://www.notion.so/ (design inspiration)
- WCAG 2.1: https://www.w3.org/WAI/WCAG21/quickref/

## Conclusion

Planning complete! The redesign will transform Scentora into a professional, Notion-inspired accord management system with:
- Clean, minimal aesthetic
- Improved usability
- Better performance
- Enhanced accessibility
- Delightful interactions

Ready to begin implementation when approved.

---

**Next Action**: Await approval and begin Phase 8.9.1 (Foundation & Setup)
