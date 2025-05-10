import { redirect, error as svelteKitError, type HttpError, isRedirect, fail } from '@sveltejs/kit'; // <<<--- ADDED 'fail' HERE
import type { Actions, PageServerLoad } from './$types';
import type { BackendWorkout, BackendWorkoutEntry } from '$lib/types';

const GO_API_URL = process.env.PRIVATE_GO_API_URL || 'http://app:8080';

export const load: PageServerLoad = async ({ locals, fetch, params, cookies, url }) => {
	if (!locals.authenticated || !locals.token) {
		throw redirect(303, `/login?redirectTo=${encodeURIComponent(url.pathname + url.search)}`);
	}

	const workoutId = params.id;
	if (!workoutId || isNaN(parseInt(workoutId))) {
		throw svelteKitError(400, 'Valid Workout ID is required.');
	}

	try {
		const response = await fetch(`${GO_API_URL}/workouts/${workoutId}`, {
			headers: {
				Authorization: `Bearer ${locals.token}`
			}
		});

		if (response.status === 401 || response.status === 403) {
			cookies.delete('auth_token', { path: '/' });
			locals.authenticated = false;
			locals.token = undefined;
			throw redirect(303, `/login?reason=session_expired&redirectTo=${encodeURIComponent(url.pathname + url.search)}`);
		}

		if (response.status === 404) {
			throw svelteKitError(404, 'Workout not found.');
		}
		if (!response.ok) {
			let errorDetail = `API error (${response.status}) loading workout.`;
			try {
				const errorPayload = await response.json();
				errorDetail = errorPayload.error || errorPayload.message || errorDetail;
			} catch (jsonParseError) { /* Use default errorDetail */ }
			throw svelteKitError(response.status, errorDetail);
		}

		const data = await response.json();
		if (data && data.workout) {
			return {
				workout: data.workout as BackendWorkout
			};
		} else {
			throw svelteKitError(500, 'Unexpected data structure from API for workout details.');
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
		console.error('[Edit Workout Load] Network or unexpected SSR error:', e);
        let errorMessage = 'A server-side error occurred while fetching workout details.';
        if (e instanceof Error) {
            errorMessage = e.message;
        } else if (typeof e === 'string') {
            errorMessage = e;
        }
		throw svelteKitError(500, errorMessage);
	}
};

export const actions: Actions = {
	updateWorkout: async ({ request, fetch, locals, params, cookies }) => {
		if (!locals.token) {
			return fail(401, { formError: 'Authentication required.' });
		}

		const workoutId = params.id;
		const formData = await request.formData();
		const title = formData.get('title')?.toString();
		const description = formData.get('description')?.toString();
		const durationMinutesStr = formData.get('durationMinutes')?.toString();
		const caloriesBurnedStr = formData.get('caloriesBurned')?.toString();
		const entriesString = formData.get('entries')?.toString();

		const formValues = { workoutId, title, description, durationMinutes: durationMinutesStr, caloriesBurned: caloriesBurnedStr, entries: entriesString };

		if (!title || !description || !durationMinutesStr || !caloriesBurnedStr || !entriesString) {
			return fail(400, { ...formValues, formError: 'All workout fields and at least one entry are required.' });
		}

		let entries: Partial<BackendWorkoutEntry>[];
		try {
			entries = JSON.parse(entriesString);
			if (!Array.isArray(entries) || entries.length === 0) throw new Error('Entries must be a non-empty array.');
			for (const entry of entries) {
				if (!entry.exercise_name || typeof entry.sets !== 'number' || entry.sets < 1 || (entry.reps == null && entry.duration_seconds == null)) {
					throw new Error('Each entry must have exercise name, valid sets, and either reps or duration.');
				}
                if (entry.reps != null && entry.duration_seconds != null) {
                    throw new Error('An entry cannot have both reps and duration.');
                }
			}
		} catch (e: any) {
			return fail(400, { ...formValues, formError: `Invalid entries data: ${e.message}` });
		}

		const updatedWorkoutPayload = {
			title,
			description,
			duration_minutes: parseInt(durationMinutesStr, 10),
			calories_burned: parseInt(caloriesBurnedStr, 10),
			entries: entries.map((e, idx) => ({
				id: e.id, 
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
			const response = await fetch(`${GO_API_URL}/workouts/${workoutId}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${locals.token}`
				},
				body: JSON.stringify(updatedWorkoutPayload)
			});

			if (response.status === 401 || response.status === 403) {
				cookies.delete('auth_token', { path: '/' });
				return fail(response.status, { ...formValues, formError: "Session expired or invalid. Please log in again." });
			}
			if (response.status === 404) {
				return fail(404, { ...formValues, formError: "Workout not found. It might have been deleted." });
			}
			if (!response.ok) {
				const errorResult = await response.json().catch(() => ({}));
				return fail(response.status, {
					...formValues,
					formError: errorResult.error || errorResult.message || `Failed to update workout.`
				});
			}
			throw redirect(303, '/workouts?message=Workout_updated_successfully');
		} catch (e: unknown) { 
			if (isRedirect(e)) throw e; 
			console.error('[UpdateWorkout Action] Network or unexpected error:', e);
            let errorMessage = 'A network error occurred while updating the workout.';
            if (e instanceof Error) {
                errorMessage = e.message;
            } else if (typeof e === 'string') {
                errorMessage = e;
            }
			return fail(500, { ...formValues, formError: errorMessage });
		}
	}
};
