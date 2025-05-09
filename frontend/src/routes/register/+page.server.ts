import { redirect, fail, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoadEvent } from './$types';

export const actions: Actions = {
  user_register: async (event: PageServerLoadEvent) => {
    const { request, fetch: eventFetch } = event;
    const data = await request.formData();
    const username = data.get('username')?.toString() || '';
    const email = data.get('email')?.toString() || '';
    const password = data.get('password')?.toString() || '';
    const bio = data.get('bio')?.toString() || '';

    if (!username || !email || !password) {
      return fail(400, { username, email, bio, message: 'Username, email, and password are required.' });
    }

    try {
      const response = await eventFetch(`/api/users`, { // Relative path
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password, bio })
      });

      if (response.ok) {
        throw redirect(303, '/login?registered=true');
      } else {
        const result = await response.json().catch(() => ({ error: 'Registration failed and could not parse error.' }));
        return fail(response.status || 400, {
          username, email, bio,
          message: result.error || result.message || 'Registration failed. Please try again.'
        });
      }
    } catch (error: any) {
      if (isRedirect(error)) {
        throw error;
      }
      console.error("Register action unexpected error: ", error);
      return fail(500, {
        username, email, bio, message: 'An unexpected error occurred.'
      });
    }
  }
};
