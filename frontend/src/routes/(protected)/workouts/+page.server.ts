import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad, PageServerLoadEvent } from './$types';
import type { BackendWorkout, BackendWorkoutEntry } from '$lib/types';


export const load: PageServerLoad = async (event: PageServerLoadEvent) => {
    // Auth handled by (protected) layout. User data from root layout.
    if (!event.locals.authenticated || !event.locals.token) {
        return { workouts: [], error: 'Not authenticated', user: event.locals.user };
    }

    try {
        const response = await event.fetch(`/api/workouts/`, { // Relative path
            method: 'GET',
            headers: { 'Authorization': `Bearer ${event.locals.token}` }
        });

        if (!response.ok) {
            const errorResult = await response.json().catch(() => ({ error: `API error (${response.status})` }));
            return {
                workouts: [],
                error: errorResult.error || errorResult.message || `Failed to load workouts (status: ${response.status}).`,
                user: event.locals.user
            };
        }
        const data = await response.json();
        if (data && Array.isArray(data.workouts)) {
            return {
                workouts: data.workouts as BackendWorkout[],
                user: event.locals.user
            };
        }
        return { workouts: [], error: "Unexpected data structure.", user: event.locals.user };
    } catch (e: any) {
        return {
            workouts: [],
            error: 'A network error occurred while fetching workouts.',
            user: event.locals.user
        };
    }
};

export const actions: Actions = {
  addWorkout: async (event: PageServerLoadEvent) => {
    const { request, locals, fetch: eventFetch } = event;
    if (!locals.token) {
      return fail(401, { formError: 'Authentication required.' });
    }
    const formData = await request.formData();
    const title = formData.get('title')?.toString();
    const description = formData.get('description')?.toString();
    const durationMinutesStr = formData.get('durationMinutes')?.toString();
    const caloriesBurnedStr = formData.get('caloriesBurned')?.toString();
    const entriesString = formData.get('entries')?.toString();

    // ... (rest of your validation logic for form fields) ...
    if (!title || !description || !durationMinutesStr || !caloriesBurnedStr || !entriesString) {
      return fail(400, { /* ... form values for repopulation ... */ formError: 'All fields required.' });
    }
     let entries: Partial<BackendWorkoutEntry>[];
    try {
      entries = JSON.parse(entriesString);
      // ... (entries validation) ...
    } catch (e:any) {
      return fail(400, { /* ... */ formError: `Invalid entries: ${e.message}`});
    }

    const newWorkoutPayload = {
      title,
      description,
      duration_minutes: parseInt(durationMinutesStr, 10),
      calories_burned: parseInt(caloriesBurnedStr, 10),
      entries: entries.map((e, idx) => ({ /* ... map entry data ... */ }))
    };

    try {
      const response = await eventFetch(`/api/workouts/`, { // Relative path
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${locals.token}`
        },
        body: JSON.stringify(newWorkoutPayload)
      });

      if (!response.ok) {
        const errorResult = await response.json().catch(() => ({ error: `API error ${response.status}` }));
        return fail(response.status, { /* ...form values... */ formError: errorResult.error || errorResult.message || 'Failed to add workout.' });
      }
      const createdWorkout = await response.json();
      return { success: true, message: "Workout added successfully!", createdWorkout: createdWorkout.workout }; // assuming backend wraps in 'workout'
    } catch (e: any) {
      console.error("[Workouts Action] Add workout error:", e.message);
      return fail(500, { /* ...form values... */ formError: 'A network error occurred.' });
    }
  },

  deleteWorkout: async (event: PageServerLoadEvent) => {
    const { request, locals, fetch: eventFetch } = event;
    if (!locals.token) return fail(401, { message: 'Auth required.' });

    const formData = await request.formData();
    const workoutId = formData.get('workoutId')?.toString();
    if (!workoutId) return fail(400, { message: 'Workout ID missing.' });

    try {
        const response = await eventFetch(`/api/workouts/${workoutId}`, { // Relative path
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${locals.token}` }
        });
        if (response.ok || response.status === 204) { // 204 is also success for DELETE
            return { success: true, deletedId: workoutId, message: "Workout deleted." };
        }
        const errResult = await response.json().catch(() => ({}));
        return fail(response.status, { message: errResult.error || errResult.message || "Failed to delete." });
    } catch (e:any) {
        console.error("[Workouts Action] Delete workout error:", e.message);
        return fail(500, { message: "Network error deleting workout." });
    }
  }
};
