import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	// The (protected)/+layout.server.ts ensures authentication.
	// locals.user will be undefined in this "no backend change" setup.
	// The profile page will primarily confirm the user is logged in.
	// No specific data needs to be loaded here for the profile page itself.
	return {
		// user: locals.user // This is available from root layout data, will be undefined
	};
};
