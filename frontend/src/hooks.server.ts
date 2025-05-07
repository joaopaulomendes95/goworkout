import { type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    // Get the token from cookies
    const token = event.cookies.get('auth_token');

    // if the tokens exists, set in locals
    if (token) {
        event.locals.token = token;
        event.locals.authenticated = true;
    } else {
        event.locals.authenticated = false;
    }

    return await resolve(event);
}
