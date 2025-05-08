// frontend/src/routes/logout/+server.ts
import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ cookies, locals }) => {
  cookies.delete('auth_token', { path: '/' });
  locals.user = undefined;
  locals.token = undefined;
  locals.authenticated = false;

  // Redirect to login page after logout
  throw redirect(303, '/login');
};
