import { type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    const token = event.cookies.get('auth_token');
    // Check if token exists
    if (token) {
      event.locals.token = token;
      event.locals.authenticated = true;
      event.locals.user = undefined; // Initialize user to null
    } else {
      event.locals.token = undefined;
      event.locals.authenticated = false;
      event.locals.user = undefined; 
    }

    return await resolve(event);
};
