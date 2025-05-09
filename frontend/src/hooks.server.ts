import { type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    const token = event.cookies.get('auth_token');

    if (token) {
      event.locals.token = token;
      event.locals.authenticated = true;
      event.locals.user = undefined;
    } else {
      event.locals.token = undefined;
      event.locals.authenticated = false;
      event.locals.user = undefined; 
    }

    const response = await resolve(event);

    return response;
};
