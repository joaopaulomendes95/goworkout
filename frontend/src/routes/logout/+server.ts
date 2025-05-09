import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ cookies, locals }) => {
	cookies.delete('auth_token', { path: '/' });
	locals.user = undefined;
	locals.token = undefined;
	locals.authenticated = false;

	throw redirect(303, '/login?logged_out=true'); // Redirect to login page after logout
};
