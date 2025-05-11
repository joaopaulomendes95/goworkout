import { redirect, fail, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import type { BackendUser } from '$lib/types';

const GO_API_URL = process.env.PRIVATE_GO_API_URL || 'http://app:8080';

// Optional: If you want to redirect logged-in users away from login page
export const load: PageServerLoad = async ({ locals, url }) => {
	if (locals.authenticated) {
		const redirectTo = url.searchParams.get('redirectTo') || '/workouts';
		throw redirect(303, redirectTo);
	}
	return {};
};

export const actions: Actions = {
	user_login: async ({ request, cookies, fetch, locals, url: pageUrl }) => {
		const data = await request.formData();
		const username = data.get('username')?.toString() || '';
		const password = data.get('password')?.toString() || '';

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

			if (response.ok && result.auth_token?.token) {
				locals.token = result.auth_token.token;
				locals.authenticated = true;
				locals.user = result.user;

				cookies.set('auth_token', result.auth_token.token, {
					path: '/',
					httpOnly: true,
					secure: process.env.NODE_ENV === 'production',
					sameSite: 'strict',
					maxAge: 60 * 60 * 24 // 1 day
				});

				const redirectTo = pageUrl.searchParams.get('redirectTo') || '/workouts';
				throw redirect(303, redirectTo);
			} else {
				return fail(response.status || 401, {
					username,
					message: result.error || result.message || 'Invalid username or password.'
				});
			}
		} catch (error: any) {
			if (isRedirect(error)) {
				throw error;
			}
			console.error('[Login Action] Unexpected error:', error);
			return fail(500, {
				username,
				message: 'An unexpected error occurred. Please try again.'
			});
		}
	}
};
