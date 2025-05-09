import { redirect, fail, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoadEvent } from './$types'; // Use PageServerLoadEvent for event type in actions

export const actions: Actions = {
  user_login: async (event: PageServerLoadEvent) => { // event here is RequestEvent
    const { request, cookies, locals, url: pageUrl, fetch: eventFetch } = event;
    const data = await request.formData();
    const username = data.get('username')?.toString() || '';
    const password = data.get('password')?.toString() || '';

    if (!username || !password) {
      return fail(400, { username, message: 'Username and password are required.' });
    }

    try {
      // Use eventFetch (SvelteKit's fetch) for API calls in actions
      const response = await eventFetch(`/api/tokens/authentication`, { // Relative path
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      const result = await response.json();

      if (response.ok && result.auth_token && result.auth_token.token) {
        // locals are only for the current request, setting them here won't persist for the redirect
        // The cookie is the primary way to establish the session for subsequent requests.
        // The hook will re-populate locals based on the new cookie.

        cookies.set('auth_token', result.auth_token.token, {
          path: '/',
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production', // SvelteKit/Vite sets this
          sameSite: 'strict',
          maxAge: 60 * 60 * 24 * 7 // 7 days
        });

        const redirectTo = pageUrl.searchParams.get('redirectTo') || '/workouts';
        throw redirect(303, redirectTo);
      } else {
        console.warn('[Login Action] Login failed. API status:', response.status, 'API result:', result);
        return fail(response.status || 401, {
          username,
          message: result.error || result.message || 'Invalid username or password.'
        });
      }
    } catch (error: any) {
      if (isRedirect(error)) {
        throw error;
      }
      console.error("[Login Action] Unexpected error during login:", error);
      return fail(500, {
        username,
        message: 'An unexpected error occurred. Please try again later.'
      });
    }
  }
};
