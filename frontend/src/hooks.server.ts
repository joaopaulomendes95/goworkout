import { type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('auth_token');

	if (token) {
		try {
			const response = await fetch('http://localhost:8080/api/me', {
				headers: { Authorization: `Bearer ${token}` }
			});
			if (response.ok) {
				const data = await response.json();
				event.locals.user = data.user;
				event.locals.authenticated = true;
			}
		} catch {
			event.locals.authenticated = false;
		}
	} else {
		event.locals.authenticated = false; // If user data can't be fetched, treat as not authenticated for session.
		event.locals.token = undefined;
		event.locals.user = undefined; // If user data can't be fetched, treat as not authenticated for session.
	}

	// If /users/me endpoint doesn't exist, and you only want to check for token presence:
	// if (token && !event.locals.user) { // If token exists but user wasn't fetched (e.g. no /users/me)
	//   event.locals.authenticated = true; // Optimistically set, API calls will validate
	// }

	return await resolve(event);
};
