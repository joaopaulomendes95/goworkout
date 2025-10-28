import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async (event) => {
	return {
		authenticated: event.locals.authenticated,
		user: event.locals.user
	};
};
