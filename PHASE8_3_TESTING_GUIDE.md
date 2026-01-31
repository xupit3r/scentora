# Phase 8.3: Frontend Testing Guide

**Date**: January 31, 2026  
**Status**: Backend Complete ✅ | Frontend Ready for Testing  
**Phase**: 8.3 - Frontend Accord Integration

---

## Quick Start

### 1. Server Status
Both servers should be running:
- **Backend**: http://localhost:3000 ✅ Running
- **Frontend**: http://localhost:5173 ✅ Running
- **PostgreSQL**: port 5435 ✅ Running

### 2. Test Credentials
```
Email: demo@scentora.com
Password: demo1234
```

### 3. Sample Data
6 accords pre-created:
- Bergamot Essential Oil (top, 25.5ml)
- Lemon Essential Oil (top, 40ml)
- Lavender Absolute (middle, 15ml)
- Rose Otto (middle, 5.5ml)
- Sandalwood Essential Oil (base, 30ml)
- Patchouli Dark (base, 20ml)

---

## Backend API Verification ✅

All backend endpoints tested and working:

### Authentication
- ✅ `POST /api/auth/register` - User registration with invitation
- ✅ `POST /api/auth/login` - User authentication
- ✅ `GET /api/auth/me` - Current user info

### Accords
- ✅ `GET /api/accords` - List accords (returns 6 accords)
- ✅ `POST /api/accords` - Create accord
- ✅ `GET /api/accords/:id` - Get single accord
- ✅ `PUT /api/accords/:id` - Update accord
- ✅ `DELETE /api/accords/:id` - Delete accord

### Tags
- ✅ `GET /api/tags` - Get all predefined tags (57 tags across 9 categories)
- ✅ `GET /api/tags/search?q=...` - Search tags
- ✅ `GET /api/tags/grouped` - Tags grouped by category

---

## Frontend Testing Checklist

### Phase 1: Login & Authentication

#### Test Steps:
1. Open http://localhost:5173 in browser
2. Should redirect to `/login` if not authenticated
3. Enter credentials:
   - Email: `demo@scentora.com`
   - Password: `demo1234`
4. Click "Login"

#### Expected Results:
- ✅ Login form displays correctly
- ✅ Form validation works (email format, password required)
- ✅ Login succeeds and redirects to home page
- ✅ Auth token stored in localStorage
- ✅ User menu shows "demo" username in header

---

### Phase 2: Accord List View

#### Test Steps:
1. After login, you should see the Accord Inventory page
2. Check that sample accords are displayed

#### Expected Results:
- ✅ Page title shows "Accord Inventory"
- ✅ "New Accord" button visible in top right
- ✅ "Filters" button visible
- ✅ All 6 sample accords displayed in grid
- ✅ Each accord card shows:
  - Name
  - Pyramid position badge (top/middle/base)
  - Volume (ml and drops)
  - Supplier
  - Tags
  - Edit and Delete buttons

#### Visual Check:
- Cards should have proper spacing
- Position badges should have colors:
  - Top: Yellow/amber tones
  - Middle: Purple tones
  - Base: Brown/earth tones
- Tags displayed as small badges
- Hover effects on cards

---

### Phase 3: Create New Accord

#### Test Steps:
1. Click "New Accord" button
2. Modal/form should open
3. Fill out form:
   - **Name**: "Vanilla Absolute"
   - **Pyramid Position**: Select "base"
   - **Volume (ml)**: 12.5
   - **Supplier**: "Eden Botanicals"
   - **Tags**: Try typing and selecting tags
   - **Notes**: "Rich, creamy vanilla scent"
4. Click "Save" or "Create"

#### Expected Results:
- ✅ Modal opens smoothly
- ✅ Form has all fields
- ✅ Tag autocomplete works (shows suggestions as you type)
- ✅ Volume validation (must be positive number)
- ✅ Position dropdown works
- ✅ Form submits successfully
- ✅ New accord appears in grid immediately
- ✅ Modal closes after success

#### Edge Cases to Test:
- Try leaving required fields empty (should show validation)
- Try entering negative volume (should prevent or show error)
- Try adding custom tags (type a tag not in predefined list)

---

### Phase 4: Edit Existing Accord

#### Test Steps:
1. Click "Edit" button on any accord card
2. Modal should open with existing data pre-filled
3. Modify some fields:
   - Change volume
   - Add/remove tags
   - Update notes
4. Click "Save"

#### Expected Results:
- ✅ Edit modal opens with current data
- ✅ All fields editable
- ✅ Tags pre-selected
- ✅ Changes save successfully
- ✅ Accord card updates in grid
- ✅ Volume drops recalculated automatically

---

### Phase 5: Delete Accord

#### Test Steps:
1. Click "Delete" button on an accord card
2. Confirmation modal should appear
3. Click "Cancel" first (should close without deleting)
4. Click "Delete" again
5. This time click "Confirm" or "Delete"

#### Expected Results:
- ✅ Confirmation modal appears
- ✅ Shows accord name being deleted
- ✅ Warning message clear
- ✅ "Cancel" button works
- ✅ "Delete" button removes accord
- ✅ Accord disappears from grid
- ✅ Database updated (accord gone on refresh)

---

### Phase 6: Filtering

#### Test Steps:
1. Click "Filters" button to open filter sidebar
2. Test each filter:

**Position Filter**:
- Select "top" → should show only top notes (Bergamot, Lemon)
- Select "middle" → should show only middle notes (Lavender, Rose)
- Select "base" → should show only base notes (Sandalwood, Patchouli)

**Volume Filter**:
- Set min volume: 20 → should filter out accords < 20ml
- Set max volume: 30 → combined with min, narrow range

**Supplier Filter**:
- Type "Eden" → should show only Eden Botanicals accords

**Tag Filter**:
- Select tag "citrus" → should show citrus-tagged accords
- Add multiple tags → should show accords matching any/all tags

3. Click "Clear Filters" → should show all accords again

#### Expected Results:
- ✅ Filters sidebar opens/closes smoothly
- ✅ Each filter updates grid immediately
- ✅ Multiple filters work together (AND logic)
- ✅ Result count updates
- ✅ "Clear Filters" resets everything
- ✅ Empty state shows if no results

---

### Phase 7: Search

#### Test Steps:
1. Look for search input (might be in filters or header)
2. Type search terms:
   - "lemon" → should find Lemon Essential Oil
   - "lavender" → should find Lavender Absolute
   - "woody" → should find accords with woody tag

#### Expected Results:
- ✅ Search works in real-time or on Enter
- ✅ Searches name, notes, and potentially tags
- ✅ Case-insensitive
- ✅ Clears with filter reset

---

### Phase 8: Empty States

#### Test Steps:
1. Apply filters that return no results
2. Observe empty state message

#### Expected Results:
- ✅ Shows helpful empty state
- ✅ Message like "No accords match your filters"
- ✅ Button to clear filters
- ✅ Nice icon or illustration

---

### Phase 9: Tag Autocomplete

#### Test Steps:
1. Open create/edit form
2. Click in tags field
3. Start typing:
   - "flo" → should suggest "floral"
   - "war" → should suggest "warm"
   - "cit" → should suggest "citrus"
4. Try selecting from dropdown
5. Try typing custom tag (not in predefined list)

#### Expected Results:
- ✅ Dropdown appears while typing
- ✅ Shows relevant suggestions
- ✅ Grouped by category
- ✅ Click to select works
- ✅ Multiple tags can be added
- ✅ Tags can be removed (X button)
- ✅ Custom tags can be added

---

### Phase 10: Responsive Design

#### Test Steps:
1. Resize browser window to different widths:
   - Mobile: 375px
   - Tablet: 768px
   - Desktop: 1200px+

#### Expected Results:
- ✅ Grid adjusts columns (1 col mobile, 2-3 tablet, 3-4 desktop)
- ✅ Filters become drawer/modal on mobile
- ✅ Navigation stays accessible
- ✅ Forms are usable on small screens
- ✅ Text is readable at all sizes

---

### Phase 11: Loading States

#### Test Steps:
1. Throttle network in browser DevTools
2. Reload page
3. Create/edit/delete accord
4. Observe loading indicators

#### Expected Results:
- ✅ Initial load shows spinner or skeleton
- ✅ Buttons show "Loading..." or spinner during actions
- ✅ Forms disabled during submit
- ✅ No double-submit possible

---

### Phase 12: Error Handling

#### Test Steps:
1. Stop backend server (kill the process)
2. Try to create an accord
3. Try to load accords
4. Start backend again
5. Try operations again

#### Expected Results:
- ✅ Network errors show error message
- ✅ Error messages are user-friendly
- ✅ Retry button or auto-retry
- ✅ No crashes or white screen
- ✅ Works again when backend returns

---

## Known Issues to Watch For

### Frontend Issues to Check:
1. **API Response Format**: Frontend expects `data.accords` but backend returns `{ accords: [] }` - verify this works
2. **Tag Array Handling**: Backend returns tags as array, frontend should display them
3. **Date Formatting**: purchaseDate is optional, check null handling
4. **Volume Display**: Both ml and drops should show
5. **Position Badge Colors**: Should match design (yellow/purple/brown)

### Backend Issues (Already Fixed):
- ✅ List endpoint returns accords correctly
- ✅ Tags are loaded for each accord
- ✅ Authentication works
- ✅ CORS configured properly

---

## Browser Console Checks

Open DevTools (F12) and check:

### Console Tab:
- No errors on page load
- No errors on login
- No errors on accord operations
- API calls logged (if debug enabled)

### Network Tab:
- `GET /api/accords` returns 200 with accord array
- `POST /api/accords` returns 201 with created accord
- `PUT /api/accords/:id` returns 200 with updated accord
- `DELETE /api/accords/:id` returns 204 (no content)
- `GET /api/tags` returns 200 with tags array

### Application Tab:
- localStorage has `accessToken`
- localStorage has `refreshToken`
- localStorage has user info (if stored)

---

## Performance Checks

### Page Load:
- Initial page load < 2 seconds
- Accord list loads < 500ms
- Filters apply < 300ms

### Interactions:
- Forms open instantly
- Tag autocomplete < 200ms
- Create/edit/delete < 1 second

---

## Accessibility Checks

### Keyboard Navigation:
- Tab through all interactive elements
- Enter key submits forms
- Escape key closes modals
- Arrow keys work in dropdowns

### Screen Reader:
- Form labels are read correctly
- Buttons have aria-labels
- Error messages announced
- Status updates announced

---

## What to Document

As you test, document:

### ✅ Works Perfectly:
- List features that work without issues

### ⚠️ Minor Issues:
- Visual glitches
- Confusing UX
- Missing validation
- Performance issues

### ❌ Broken:
- Features that don't work
- Errors thrown
- Data not saving
- Navigation broken

### 💡 Improvements:
- UX suggestions
- Missing features noticed
- Better error messages needed

---

## Testing Complete Criteria

Phase 8.3 is complete when:

- [ ] Login/authentication works
- [ ] All 6 sample accords display
- [ ] Create new accord works
- [ ] Edit accord works
- [ ] Delete accord works (with confirmation)
- [ ] Filters work (position, volume, supplier, tags)
- [ ] Search works
- [ ] Tag autocomplete works
- [ ] Responsive on mobile/tablet/desktop
- [ ] Loading states display
- [ ] Error states display
- [ ] No console errors
- [ ] API calls succeed
- [ ] Data persists (survives refresh)

---

## Next Steps After Testing

Once testing is complete:

1. **Document Results**: Create test results file
2. **Fix Issues**: Address any bugs found
3. **Update Components**: Make UI/UX improvements
4. **Phase Completion**: Mark Phase 8.3 complete
5. **Move to Phase 8.4**: Statistics view or Phase 8.9 UI redesign

---

## Quick Debug Commands

If you need to check backend:

```bash
# Check backend is running
curl http://localhost:3000/health

# List accords (replace TOKEN)
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:3000/api/accords

# Check database
docker exec scentora-postgres psql -U admin -d scentora -c "SELECT id, name, pyramid_position FROM accords;"

# Check predefined tags
curl http://localhost:3000/api/tags | python3 -m json.tool | grep '"tag"' | wc -l
# Should return 57
```

---

## Support

If you encounter issues:

1. Check browser console for errors
2. Check network tab for failed requests
3. Check backend logs: `/tmp/backend.log`
4. Verify database has data
5. Try logging out and back in
6. Try hard refresh (Ctrl+Shift+R)

---

**Happy Testing!** 🧪✨

Report back with results and we'll iterate on any issues found.
