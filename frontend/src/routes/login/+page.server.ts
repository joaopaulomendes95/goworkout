import { redirect, fail, isRedirect } from '@sveltejs/kit'; // Import isRedirect
import type { Actions } from './$types';

const GO_API_URL = 'http://app:8080'; // Adjust if different

export const actions: Actions = {
  login: async ({ request, cookies, fetch, locals, url: pageUrl }) => {
    const data = await request.formData();
    const username = data.get('username')?.toString() || '';
    const password = data.get('password')?.toString() || '';

    console.log('[Login Action] Username from form data:', username);
    console.log('[Login Action] Password from form data exists:', !!password);

    if (!username || !password) {
      return fail(400, { username, message: 'Username and password are required.' });
    }

    try {
      const response = await fetch(`${GO_API_URL}/tokens/authentication`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      const result = await response.json();
      console.log('[Login Action] API Response Status:', response.status);
      console.log('[Login Action] API Response Body:', result);


      if (response.ok && result.token) {
        locals.token = result.token;
        if (result.user) { // Assuming login response includes user details
            locals.user = result.user as App.Locals['user'];
        }
        locals.authenticated = true;
        
        cookies.set('auth_token', result.token, {
          path: '/',
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          sameSite: 'strict',
          maxAge: 60 * 60 * 24 // 1 day
        });
        
        const redirectTo = pageUrl.searchParams.get('redirectTo') || '/workouts';
        console.log('[Login Action] Successful login, redirecting to:', redirectTo);
        throw redirect(303, redirectTo); // This will be caught and re-thrown if it's a redirect
      } else {
        console.warn('[Login Action] Login failed. API status:', response.status, 'API result:', result);
        return fail(response.status || 401, {
          username,
          message: result.error || result.message || 'Invalid username or password.'
        });
      }
    } catch (error: any) {
      // Check if the caught error is a SvelteKit redirect.
      // If so, re-throw it to let SvelteKit handle the client-side redirection.
      if (isRedirect(error)) {
        console.log('[Login Action] Caught redirect, re-throwing:', error.location);
        throw error; // Re-throw the redirect
      }

      // If it's not a redirect, it's a genuine network or other unexpected error
      console.error("[Login Action] Unexpected error during login:", error);
      return fail(500, {
        username,
        message: 'An unexpected error occurred during login. Please try again later.'
      });
    }
  }
};
