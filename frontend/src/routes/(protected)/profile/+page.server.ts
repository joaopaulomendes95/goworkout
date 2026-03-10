import type { PageServerLoad } from './$types';
import type { BackendWorkout } from '$lib/types';

const GO_API_URL = 'http://app:8080';

export const load: PageServerLoad = async ({ locals, parent }) => {
	const parentData = await parent();

	if (!locals.authenticated || !locals.token) {
		return { workouts: [], error: 'Not authenticated', user: null, authenticated: false };
	}

	try {
		const response = await fetch(`${GO_API_URL}/workouts/`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${locals.token}`
			}
		});

		if (!response.ok) {
			const errorResult: { error?: string; message?: string } = await response
				.json()
				.catch(() => ({ error: `API error (${response.status}) fetching workouts.` }));
			console.error('[Workouts Load] API Error fetching workouts:', response.status, errorResult);
			return {
				workouts: [],
				user: parentData.user,
				authenticated: parentData.authenticated,
				error:
					errorResult.error ||
					errorResult.message ||
					`Failed to load workouts (status: ${response.status}).`
			};
		}

		const data = await response.json();

		if (data && Array.isArray(data.workouts)) {
			return {
				workouts: data.workouts as BackendWorkout[],
				user: parentData.user,
				authenticated: parentData.authenticated
			};
		} else {
			console.error('[Workouts Load] Unexpected data structure from API:', data);
			return {
				workouts: [],
				user: parentData.user,
				authenticated: parentData.authenticated,
				error: 'Unexpected data structure from API.'
			};
		}
	} catch (e) {
		console.error('[Workouts Load] Network or unexpected error:', e);
		return {
			workouts: [],
			user: parentData.user,
			authenticated: parentData.authenticated,
			error: 'A network error occurred while fetching workouts. Please try again.'
		};
	}
};
