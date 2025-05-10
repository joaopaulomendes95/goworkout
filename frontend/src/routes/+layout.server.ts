import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	return {
		authenticated: locals.authenticated,
		user: locals.user,
		token: locals.token,
	};
};
