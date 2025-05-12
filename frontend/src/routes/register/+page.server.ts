import { redirect, fail, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

const GO_API_URL = process.env.PRIVATE_GO_API_URL || 'http://app:8080';

// If user already authenticated, redirect to /workouts
export const load: PageServerLoad = async ({ locals }) => {
	if (locals.authenticated) {
		throw redirect(303, '/workouts');
	}
	return {};
};

// Register action
export const actions: Actions = {
	user_register: async ({ request, fetch }) => {
		const data = await request.formData();
		const username = data.get('username')?.toString() || '';
		const email = data.get('email')?.toString() || '';
		const password = data.get('password')?.toString() || '';
		const bio = data.get('bio')?.toString() || '';

		// For repopulating form
		const formValues = { username, email, bio };

		// Validations
		if (!username || !email || !password) {
			return fail(400, { ...formValues, message: 'Username, email, and password are required.' });
		}
		if (username.length < 3) {
			return fail(400, { ...formValues, message: 'Username must be at least 3 characters long.' });
		}
		if (username.length > 20) {
			return fail(400, { ...formValues, message: 'Username must be at least 20 characters long.' });
		}
		if (password.length < 3) {
			return fail(400, { ...formValues, message: 'Password must be at least 3 characters long.' });
		}
		if (password.length > 20) {
			return fail(400, { ...formValues, message: 'Password must be at most 20 characters long.' });
		}

		// Try to register the user
		try {
			const response = await fetch(`${GO_API_URL}/users`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, email, password, bio })
			});

			// If backend response is ok
			if (response.ok) {
				throw redirect(303, '/login?registered=true');
			} else {
				const result = await response.json().catch(() => ({ error: 'Registration failed and could not parse error response.' }));
				return fail(response.status || 400, {
					...formValues,
					message: result.error || result.message || 'Registration failed. Please try again.'
				});
			}

		// Handle unexpected errors
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
