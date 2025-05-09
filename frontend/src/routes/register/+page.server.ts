import { redirect, fail, isRedirect } from '@sveltejs/kit';
import type { Actions } from './$types';

// Adjust to match GO backend URL
const GO_API_URL = 'http://app:8080';

export const actions: Actions = {
  // Renamed from user_register to default for simpler form action
  user_register: async ({ request, fetch }) => { 
    const data = await request.formData();
    const username = data.get('username')?.toString() || '';
    const email = data.get('email')?.toString() || '';
    const password = data.get('password')?.toString() || '';
    const bio = data.get('bio')?.toString() || ''; // Bio is optional in Go backend

    // Basic client-side validation (backend will also validate)
    if (!username || !email || !password) {
      return fail(400, { username, email, bio, message: 'Username, email, and password are required.' });
    }

    try {
      const response = await fetch(`${GO_API_URL}/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password, bio })
      });

      if (response.ok) { // Backend returns 201 Created on success
        // const result = await response.json(); // Contains created user if needed
        throw redirect(303, '/login?registered=true');
      } else {
        const result = await response.json().catch(() => ({error: 'Registration failed and could not parse error.'}));
        return fail(response.status || 400, {
          username, email, bio, // Repopulate form
          message: result.error || result.message || 'Registration failed. Please try again.'
        });
      }
    } catch (error: any) {
      // Check if error is Svelte redirect
      if (isRedirect(error)) {
        // let Svelte handle it
        throw error;
      }
      
      console.error("Register action unexpected error: ", error);
      return fail(500, {
        username, email, bio, message: 'An unexpected error occurred.'
      });
    }
  }
};
