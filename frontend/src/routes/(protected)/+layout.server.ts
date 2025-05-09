import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ url, parent }) => {
	const parentData = await parent(); // Gets { authenticated, user } from root layout

	if (!parentData.authenticated) {
		// Redirect to login page and store intended destination
		throw redirect(303, `/login?redirectTo=${encodeURIComponent(url.pathname + url.search)}`);
	}

	// User is authenticated, proceed.
	// Data like 'authenticated' and 'user' (even if undefined) is already available via parentData.
	return {
		// No need to return parentData.authenticated or parentData.user again,
		// as they are inherited by child layouts/pages.
	};
};
