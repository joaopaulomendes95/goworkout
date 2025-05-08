import { type Handle } from '@sveltejs/kit';

// This should be the URL your SvelteKit server uses to reach the Go backend.
// If running in Docker with Nginx, this is likely the service name.
const GO_API_URL = 'http://app:8080'; // Change if needed

export const handle: Handle = async ({ event, resolve }) => {
    const token = event.cookies.get('auth_token');

    if (token) {
      event.locals.token = token;
      event.locals.authenticated = true; // If user data can't be fetched, treat as not authenticated for session.
      event.locals.user = undefined; // If user data can't be fetched, treat as not authenticated for session.
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
