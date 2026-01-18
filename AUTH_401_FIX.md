# Auth 401 Error - Diagnosis and Fix

## Problems Identified

### 1. **Axios Interceptor Not Initialized Early**
The axios interceptor for handling token refresh was defined in the auth store, but the store wasn't initialized until first use. This meant:
- If a user had an expired token in localStorage
- The first API call would fail with 401
- The interceptor wouldn't be set up yet to handle it

### 2. **No Token Expiry Check on App Load**
When the app loaded with a token in localStorage:
- It would set the Authorization header
- But wouldn't check if the token was expired
- First API call would fail with 401

### 3. **Potential Infinite Loop on Refresh**
If the refresh token itself returned 401, the interceptor might try to refresh it again, causing a loop.

## Fixes Applied

### Fix 1: Early Store Initialization
**File:** `frontend/src/main.ts`

```typescript
import { useAuthStore } from './stores/auth';

// ... after app.use(pinia)
useAuthStore(); // Initialize auth store early
app.mount('#app');
```

This ensures the axios interceptor is set up before any API calls.

### Fix 2: Token Expiry Check Function
**File:** `frontend/src/stores/auth.ts`

Added `isTokenExpired()` function:
```typescript
function isTokenExpired(token: string): boolean {
  try {
    const parts = token.split('.');
    const payload = JSON.parse(atob(parts[1]));
    const exp = payload.exp * 1000;
    const now = Date.now();
    
    // Token is expired if it expires in less than 30 seconds
    return exp - now < 30000;
  } catch (e) {
    return true;
  }
}
```

### Fix 3: Proactive Token Refresh on Init
**File:** `frontend/src/stores/auth.ts`

Modified initialization code:
```typescript
if (accessToken.value) {
  if (isTokenExpired(accessToken.value) && refreshToken.value) {
    // Token is expired, refresh immediately
    axios.post(`${API_BASE}/auth/refresh`, {
      refreshToken: refreshToken.value,
    })
    .then((response) => {
      storeTokens(response.data.accessToken, response.data.refreshToken);
    })
    .catch(() => {
      clearTokens();
    });
  } else {
    setAuthHeader(accessToken.value);
  }
}
```

### Fix 4: Prevent Refresh Loop
**File:** `frontend/src/stores/auth.ts`

Updated interceptor to skip refresh requests:
```typescript
// Don't try to refresh if this was the refresh request itself
if (originalRequest.url?.includes('/auth/refresh')) {
  return Promise.reject(error);
}
```

## Testing

### Quick Test Steps

1. **Start the application:**
   ```bash
   ./scentora.sh start
   ```

2. **Run the auth flow test:**
   ```bash
   ./test-auth-flow.sh
   ```

3. **Test in browser:**
   - Open http://localhost:5173
   - Register a new user or login
   - Navigate around the app
   - Check browser console for errors
   - Check Network tab for 401 responses

### Manual Test in Browser Console

Open browser console on http://localhost:5173 and run:

```javascript
// Check tokens
console.log('Access Token:', localStorage.getItem('scentora_access_token'));
console.log('Refresh Token:', localStorage.getItem('scentora_refresh_token'));

// Decode and check expiry
const token = localStorage.getItem('scentora_access_token');
if (token) {
  const payload = JSON.parse(atob(token.split('.')[1]));
  console.log('Expires:', new Date(payload.exp * 1000));
  console.log('Expired:', new Date() > new Date(payload.exp * 1000));
}

// Test API call
fetch('http://localhost:3000/api/perfumes', {
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('scentora_access_token')
  }
})
.then(r => r.json())
.then(data => console.log('Perfumes:', data))
.catch(err => console.error('Error:', err));
```

## Expected Behavior Now

1. **On app load with valid token:**
   - Token is checked for expiry
   - If valid, sets Authorization header
   - API calls work normally

2. **On app load with expired token:**
   - Detects token is expired
   - Automatically calls refresh endpoint
   - Gets new tokens
   - Continues normally

3. **During usage when token expires:**
   - API call returns 401
   - Interceptor catches it
   - Automatically refreshes token
   - Retries original request
   - User doesn't notice

4. **When refresh token is invalid:**
   - Refresh fails
   - Clears all tokens
   - Redirects to /login
   - No infinite loop

## Debugging 401 Errors

If you still see 401 errors:

1. **Check browser console for errors**
2. **Check Network tab:**
   - Look for failed requests
   - Check if Authorization header is present
   - Check if refresh is being called
3. **Check localStorage:**
   - Ensure tokens are being stored
   - Check if they're valid JWTs
4. **Backend logs:**
   - Check if backend is receiving tokens
   - Check for JWT verification errors
5. **Clear localStorage and login again:**
   ```javascript
   localStorage.clear();
   window.location.href = '/login';
   ```

## Token Lifetimes

Current configuration:
- **Access Token:** 15 minutes
- **Refresh Token:** 7 days

You can adjust these in `backend/.env`:
```
JWT_ACCESS_EXPIRES_IN=15m
JWT_REFRESH_EXPIRES_IN=7d
```

## Additional Notes

- Tokens are stored in localStorage (keys: `scentora_access_token`, `scentora_refresh_token`)
- Axios default headers are set automatically with Bearer token
- Router guards check for token before allowing access to protected routes
- Rate limiting on auth endpoints: 5 requests per 15 minutes per IP
