import { redirect, fail, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types'; // Added PageServerLoad

const GO_API_URL = process.env.PRIVATE_GO_API_URL || 'http://app:8080';

// Optional: If you want to redirect logged-in users away from register page
export const load: PageServerLoad = async ({ locals }) => {
	if (locals.authenticated) {
		throw redirect(303, '/workouts');
	}
	return {};
};

export const actions: Actions = {
	user_register: async ({ request, fetch }) => {
		const data = await request.formData();
		const username = data.get('username')?.toString() || '';
		const email = data.get('email')?.toString() || '';
		const password = data.get('password')?.toString() || '';
		const bio = data.get('bio')?.toString() || ''; // Bio is optional

		const formValues = { username, email, bio }; // For repopulating form

		if (!username || !email || !password) {
			return fail(400, { ...formValues, message: 'Username, email, and password are required.' });
		}
		if (password.length < 3) { // Example additional validation
            return fail(400, { ...formValues, message: 'Password must be at least 3 characters long.' });
		}


		try {
			const response = await fetch(`${GO_API_URL}/users`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, email, password, bio })
			});

			if (response.ok) { // Backend returns 201 Created on success
				throw redirect(303, '/login?registered=true');
			} else {
				const result = await response.json().catch(() => ({ error: 'Registration failed and could not parse error response.' }));
				return fail(response.status || 400, {
					...formValues,
					message: result.error || result.message || 'Registration failed. Please try again.'
				});
			}
		} catch (error: any) {
			if (isRedirect(error)) {
				throw error;
			}
			console.error("Register action unexpected error: ", error);
			return fail(500, {
				...formValues, message: 'An unexpected error occurred during registration.'
			});
		}
	}
};
