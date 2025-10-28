// src/hooks.server.ts
import { type Handle } from '@sveltejs/kit';
import { API } from '$env/static/private';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('auth_token');

	if (token) {
		try {
			// Verify token with Go backend
			const response = await fetch(`${API}users/me`, {
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'application/json'
				}
			});

			if (response.ok) {
				const user = await response.json();
				event.locals.authenticated = true;
				event.locals.user = user;
			} else {
				// Token is invalid, clear the cookie
				event.cookies.delete('auth_token', { path: '/' });
			}
		} catch (error) {
			console.error('Auth verification failed:', error);
			// Backend unreachable or error - treat as unauthenticated
		}
	}

	return resolve(event);
};
