import { redirect, fail } from '@sveltejs/kit';
import type { Actions } from './$types';

const API_URL = 'http://localhost:8080/api';

export const actions: Actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const username = data.get('username')?.toString() || '';
		const email = data.get('email')?.toString() || '';
		const password = data.get('password')?.toString() || '';
		const bio = data.get('bio')?.toString() || '';

		if (!username || !email || !password) {
			return fail(400, { username, email, bio, error: 'All fields are required' });
		}

		try {
			const response = await fetch(`${API_URL}/users`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, email, password, bio })
			});

			const result = await response.json();

			if (!response.ok) {
				return fail(response.status, { username, email, bio, error: result.error || 'Registration failed' });
			}

			throw redirect(303, '/login?registered=true');
		} catch (error) {
			if (error instanceof Response || (error && typeof error === 'object' && 'status' in error)) {
				throw error;
			}
			return fail(500, { username, email, bio, error: 'An unexpected error occurred' });
		}
	}
};
