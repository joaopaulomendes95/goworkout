import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

const API_URL = 'http://localhost:8080/api';

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const username = data.get('username')?.toString() || '';
		const password = data.get('password')?.toString() || '';

		if (!username || !password) {
			return fail(400, { username, error: 'Username and password are required' });
		}

		try {
			const response = await fetch(`${API_URL}/tokens/authentication`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, password })
			});

			const result = await response.json();

			if (!response.ok) {
				return fail(response.status, { username, error: result.error || 'Login failed' });
			}

			cookies.set('auth_token', result.token, {
				path: '/',
				httpOnly: true,
				secure: process.env.NODE_ENV === 'production',
				sameSite: 'strict',
				maxAge: 60 * 60 * 24 * 30 // 30 days
			});

			throw redirect(303, '/protected');
		} catch (error) {
			if (error instanceof Response || (error && typeof error === 'object' && 'status' in error)) {
				throw error;
			}
			return fail(500, { username, error: 'An unexpected error occurred' });
		}
	}
};
