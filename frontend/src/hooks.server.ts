import { type Handle } from '@sveltejs/kit';

const GO_API_URL = process.env.PRIVATE_GO_API_URL || 'http://app:8080';


export const handle: Handle = async ({ event, resolve }) => {
    const token = event.cookies.get('auth_token');

    if (token) {
        event.locals.token = token;
        event.locals.authenticated = true;
        let user =  await fetch(`${GO_API_URL}/users/me`, {
			headers: {
				Authorization: `Bearer ${token}`
			}
        });
        event.locals.user = await user.json();
        console.log('[Handle] User:', event.locals.user);

    } else {
      event.locals.token = undefined;
      event.locals.authenticated = false;
      event.locals.user = undefined; 
    }

    const response = await resolve(event);

    return response;
};
