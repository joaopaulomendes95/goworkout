import { redirect, fail, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

const GO_API_URL = process.env.PRIVATE_GO_API_URL || 'http://app:8080';

export const load: PageServerLoad = async ({ locals }) => {
	console.log('[Profile Load Page] Current usser :', locals.user);
	return { user: locals.user };
};

export const actions: Actions = {
update_profile: async ({ request, cookies, locals }) => {
	const formData = await request.formData();
	console.log('Form Data Update:', formData);	
	const username = formData.get('username')?.toString() || '';
	const bio = formData.get('bio')?.toString() || '';

	const formValues = { username, bio, }; // For repopulating form

	if (!username) {
		return fail(400, {
			username, bio, message: 'username is required.'
		});
	}

	const token = cookies.get('auth_token');
	console.log('Token:', token);

	try {
		const response = await fetch(`${GO_API_URL}/users/update`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
				'Authization': `Bearer ${token}` // Include the token in the request headers
			},
			body: JSON.stringify({ username, bio })
		});

		if (response.ok) { // Backend returns 201 Updated with success
			throw redirect(303, '/profile?updated=true');
		} else {
			const result = await response.json().catch(() => ({ error: 'Update failed and could not parse error response.' }));
			return fail(response.status || 400, {
				...formValues,
				message: result.error || result.message || 'Update failed. Please try again.'
			});
		}
	} catch (error: any) {
		if (isRedirect(error)) {
			throw error;
		}
		console.error("Register action unexpected error: ", error);
		return fail(500, {
			...formValues, message: 'An unexpected error occurred during update.'
		});
	}
}};
