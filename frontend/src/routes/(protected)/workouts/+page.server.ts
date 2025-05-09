// frontend/src/routes/(protected)/workouts/+page.server.ts
import { redirect, error as svelteKitError, type HttpError, isRedirect, fail } from '@sveltejs/kit'; // <<<--- ADDED 'fail' HERE
import type { Actions, PageServerLoad } from './$types';
import type { BackendWorkout } from '$lib/types';

const GO_API_URL = process.env.PRIVATE_GO_API_URL || 'http://app:8080';

export const load: PageServerLoad = async ({ locals, fetch: svelteKitFetch, cookies, url }) => {
	if (!locals.authenticated || !locals.token) {
		throw redirect(303, `/login?redirectTo=${encodeURIComponent(url.pathname + url.search)}`);
	}

	try {
		// console.log('[Workouts Load] Fetching from:', `${GO_API_URL}/workouts/`); // Optional: for debugging
		const response = await svelteKitFetch(`${GO_API_URL}/workouts/`, {
			headers: {
				Authorization: `Bearer ${locals.token}`
			}
		});
		// console.log('[Workouts Load] API Response Status:', response.status); // Optional: for debugging

		if (response.status === 401 || response.status === 403) {
			cookies.delete('auth_token', { path: '/' });
			locals.authenticated = false;
			locals.token = undefined;
			throw redirect(303, `/login?reason=session_expired&redirectTo=${encodeURIComponent(url.pathname + url.search)}`);
		}

		if (!response.ok) {
			let errorDetail = `API error (${response.status}) fetching workouts.`;
			try {
				const errorPayload = await response.json();
				errorDetail = errorPayload.error || errorPayload.message || errorDetail;
				// console.error('[Workouts Load] API Error Payload:', errorPayload); // Optional: for debugging
			} catch (jsonParseError) {
				// console.error('[Workouts Load] Failed to parse error JSON, using status text:', response.statusText); // Optional
				errorDetail = response.statusText || errorDetail;
			}
			throw svelteKitError(response.status, errorDetail);
		}

		const responseData = await response.json();
		// console.log('[Workouts Load] API Response Data:', responseData); // Optional: for debugging

		if (responseData && Array.isArray(responseData.workouts)) {
			return {
				workouts: responseData.workouts as BackendWorkout[],
				messageFromRedirect: url.searchParams.get('message')
			};
		} else {
			console.error('[Workouts Load] Unexpected data structure from API:', responseData);
			throw svelteKitError(500, 'Unexpected data structure from API when fetching workouts.');
		}
	} catch (e: unknown) { 
		if (isRedirect(e)) {
			throw e;
		}
		if (typeof e === 'object' && e !== null && 'status' in e) {
            const httpError = e as { status: number; body?: { message: string } }; 
            if (httpError.status && httpError.status >= 400 && httpError.status < 600) {
                throw e; 
            }
		}
		
		console.error('[Workouts Load] Catch-all - Network or unexpected SSR error:', e);
        let errorMessage = 'A server-side error occurred while fetching workouts.';
        if (e instanceof Error) {
            errorMessage = e.message;
        } else if (typeof e === 'string') {
            errorMessage = e;
        }
		throw svelteKitError(500, errorMessage);
	}
};

export const actions: Actions = {
	addWorkout: async ({ request, fetch: svelteKitFetch, locals, cookies }) => {
		if (!locals.token) {
			return fail(401, { formError: 'Authentication required.' });
		}

		const formData = await request.formData();
		const title = formData.get('title')?.toString();
		const description = formData.get('description')?.toString();
		const durationMinutesStr = formData.get('durationMinutes')?.toString();
		const caloriesBurnedStr = formData.get('caloriesBurned')?.toString();
		const entriesString = formData.get('entries')?.toString();

		const formValues = { title, description, durationMinutes: durationMinutesStr, caloriesBurned: caloriesBurnedStr, entries: entriesString };

		if (!title || !description || !durationMinutesStr || !caloriesBurnedStr || !entriesString) {
			return fail(400, { ...formValues, formError: 'All workout fields and at least one entry are required.' });
		}

		let entries: Partial<import('$lib/types').BackendWorkoutEntry>[];
		try {
			entries = JSON.parse(entriesString);
			if (!Array.isArray(entries) || entries.length === 0) {
				throw new Error('Entries must be a non-empty array.');
			}
			for (const entry of entries) {
				if (!entry.exercise_name || typeof entry.sets !== 'number' || entry.sets < 1 || (entry.reps == null && entry.duration_seconds == null)) {
					throw new Error('Each entry must have exercise name, valid sets, and either reps or duration.');
				}
                 if (entry.reps != null && entry.duration_seconds != null) {
                    throw new Error('An entry cannot have both reps and duration. Please provide one or the other.');
                }
			}
		} catch (e: any) {
			return fail(400, { ...formValues, formError: `Invalid entries data: ${e.message}` });
		}

		const newWorkoutPayload = {
			title,
			description,
			duration_minutes: parseInt(durationMinutesStr, 10),
			calories_burned: parseInt(caloriesBurnedStr, 10),
			entries: entries.map((e, idx) => ({
				exercise_name: e.exercise_name,
				sets: e.sets,
				reps: e.reps && e.reps > 0 ? e.reps : null,
				duration_seconds: e.duration_seconds && e.duration_seconds > 0 ? e.duration_seconds : null,
				weight: e.weight != null && e.weight >= 0 ? e.weight : null,
				notes: e.notes || '',
				order_index: e.order_index !== undefined ? e.order_index : idx + 1
			}))
		};

		try {
			const response = await svelteKitFetch(`${GO_API_URL}/workouts/`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${locals.token}`
				},
				body: JSON.stringify(newWorkoutPayload)
			});

			if (response.status === 401 || response.status === 403) {
				cookies.delete('auth_token', { path: '/' });
				return fail(response.status, { ...formValues, formError: "Session expired or invalid. Please log in again." });
			}

			if (!response.ok) {
				const errorResult = await response.json().catch(() => ({}));
				return fail(response.status, {
					...formValues,
					formError: errorResult.error || errorResult.message || `Failed to add workout.`
				});
			}
			return { success: true, message: 'Workout added successfully!' };
		} catch (e: any) {
			console.error('[AddWorkout Action] Network or unexpected error:', e);
			return fail(500, { ...formValues, formError: 'A network error occurred. Please try again.' });
		}
	},

	deleteWorkout: async ({ request, fetch: svelteKitFetch, locals, cookies }) => {
		if (!locals.token) return fail(401, { message: 'Authentication required.' }); // This fail was correctly imported implicitly

		const formData = await request.formData();
		const workoutId = formData.get('workoutId')?.toString();
		if (!workoutId) return fail(400, { message: 'Workout ID missing.' });

		try {
			const response = await svelteKitFetch(`${GO_API_URL}/workouts/${workoutId}`, {
				method: 'DELETE',
				headers: { Authorization: `Bearer ${locals.token}` }
			});

			if (response.status === 401 || response.status === 403) {
				cookies.delete('auth_token', { path: '/' });
				return fail(response.status, { message: "Session expired or invalid. Please log in again." });
			}
			if (response.status === 404) {
                return fail(404, { message: "Workout not found, it might have already been deleted." });
            }
			if (response.ok || response.status === 204) {
				return { success: true, deletedId: workoutId, message: 'Workout deleted successfully.' };
			}
			const errResult = await response.json().catch(() => ({}));
			return fail(response.status, { message: errResult.error || errResult.message || 'Failed to delete workout.' });
		} catch (e: any) {
			console.error('[DeleteWorkout Action] Network or unexpected error:', e);
			return fail(500, { message: 'A network error occurred while deleting the workout.' });
		}
	}
};
