import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import type { BackendUser } from '$lib/types';

export const load: LayoutServerLoad = async ({
	url,
	parent
}: {
	url: URL;
	parent: () => Promise<{ authenticated: boolean; user: BackendUser | undefined }>;
}) => {
	const parentData = await parent();

	if (!parentData.authenticated) {
		throw redirect(303, `/login?redirectTo=${encodeURIComponent(url.pathname)}`);
	}

	return {};
};
