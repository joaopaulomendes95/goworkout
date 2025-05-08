// frontend/src/routes/(protected)/workouts/+page.server.ts
import { fail } from '@sveltejs/kit'; // No redirect/error from load, layout handles auth
import type { Actions, PageServerLoad } from './$types';
 import type { BackendWorkoutEntry } from '$lib/types';

const GO_API_URL = 'http://app:8080';

// This load function might not fetch workouts anymore if no list endpoint exists.
// It will mainly ensure the user is authenticated (via the (protected) layout).
export const load: PageServerLoad = async ({ locals }) => {
  // Authentication is handled by (protected)/+layout.server.ts
  // If we reach here, the user is authenticated.
  // locals.token will be available.
  
  // If there's no GET /workouts endpoint to list workouts, we return no workout data here.
  // The page will primarily be for creating workouts.
  return {
    // workouts: [] // Or simply nothing if the page doesn't display a list initially
  };
};

export const actions: Actions = {
  addWorkout: async ({ request, fetch, locals }) => {
    if (!locals.token) {
      return fail(401, { formError: 'Authentication required.' });
    }

    const formData = await request.formData();
    const title = formData.get('title')?.toString();
    const description = formData.get('description')?.toString();
    const durationMinutesStr = formData.get('durationMinutes')?.toString();
    const caloriesBurnedStr = formData.get('caloriesBurned')?.toString();
    const entriesString = formData.get('entries')?.toString();

    if (!title || !description || !durationMinutesStr || !caloriesBurnedStr || !entriesString) {
      return fail(400, {
        formError: 'All workout fields and at least one entry are required.',
        // Pass back values for repopulation
        title, description, durationMinutes: durationMinutesStr, caloriesBurned: caloriesBurnedStr, entries: entriesString
      });
    }

    let entries: Partial<BackendWorkoutEntry>[]; // Use Partial if IDs are not sent for new entries
    try {
      entries = JSON.parse(entriesString);
      if (!Array.isArray(entries) || entries.length === 0) {
        throw new Error("Entries must be a non-empty array.");
      }
      // Validate each entry if needed
      for (const entry of entries) {
        if (!entry.exercise_name || !entry.sets || (entry.reps == null && entry.duration_seconds == null)) {
            throw new Error("Each entry must have exercise name, sets, and either reps or duration.");
        }
      }
    } catch (e: any) {
      return fail(400, {
        formError: `Invalid entries data: ${e.message}`,
        title, description, durationMinutes: durationMinutesStr, caloriesBurned: caloriesBurnedStr, entries: entriesString
      });
    }

    const newWorkoutPayload = {
      title,
      description,
      duration_minutes: parseInt(durationMinutesStr, 10),
      calories_burned: parseInt(caloriesBurnedStr, 10),
      entries: entries.map((e, idx) => {
        // Determine if it's a rep-based or duration-based entry
        let finalReps = e.reps != null && e.reps > 0 ? e.reps : null;
        let finalDuration = e.duration_seconds != null && e.duration_seconds > 0 ? e.duration_seconds : null;

        // If both are somehow provided (and positive), prioritize one or make it an error earlier
        // For now, let's assume UI or prior validation ensures only one is meaningfully filled.
        // If reps are provided (and > 0), nullify duration.
        // If duration is provided (and > 0), nullify reps.
        // If both are 0 or null, the backend constraint will catch it.

        if (finalReps !== null && finalDuration !== null) {
          // This case should ideally be prevented by UI logic.
          // For now, let's say we prioritize reps if both are somehow positive.
          // Or, better, your UI should only allow one to be entered.
          // For the constraint: if reps has a value, duration MUST be null.
          finalDuration = null;
        }


        return {
          exercise_name: e.exercise_name,
          sets: e.sets,
          reps: finalReps,
          duration_seconds: finalDuration,
          weight: e.weight != null ? e.weight : null, // Send null if not applicable
          notes: e.notes || "",
          order_index: e.order_index !== undefined ? e.order_index : idx + 1,
        };
      })
    };

    console.log("[Workouts Action] Sending payload to POST /workouts:", JSON.stringify(newWorkoutPayload, null, 2));

    try {
      const response = await fetch(`${GO_API_URL}/workouts/`, { // POST /workouts
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${locals.token}`
        },
        body: JSON.stringify(newWorkoutPayload)
      });

      if (!response.ok) { // e.g., 201 Created is ok
        const errorResult = await response.json().catch(() => ({ error: `API error (${response.status}) and failed to parse error response.` }));
        console.error("[Workouts Action] API Error creating workout:", response.status, errorResult);
        return fail(response.status, {
          formError: errorResult.error || errorResult.message || `Failed to add workout.`,
          title, description, durationMinutes: durationMinutesStr, caloriesBurned: caloriesBurnedStr, entries: entriesString
        });
      }

      const createdWorkout = await response.json(); // Backend should return the created workout
      console.log("[Workouts Action] Workout created successfully:", createdWorkout);

      // Instead of invalidating a list, we can return the created workout
      // or just a success message. If you want to show a list client-side,
      // you'd add this `createdWorkout` to a client-side store.
      return { success: true, message: "Workout added successfully!", createdWorkout };
    } catch (e: any) {
      console.error("[Workouts Action] Network or unexpected error:", e);
      return fail(500, {
        formError: 'A network error occurred. Please try again.',
        title, description, durationMinutes: durationMinutesStr, caloriesBurned: caloriesBurnedStr, entries: entriesString
      });
    }
  },

  // deleteWorkout action would still be relevant if you can list/identify workouts to delete
  // but how would the user know which workout ID to delete without a list?
  // This might be for a different view (e.g., a workout detail page).
  deleteWorkout: async ({ request, fetch, locals }) => {
    // ... (implementation for DELETE /workouts/{id})
    // This action makes sense if you navigate to a page showing a single workout
    // and then delete it.
    if (!locals.token) return fail(401, { message: 'Auth required.' });
    const formData = await request.formData();
    const workoutId = formData.get('workoutId')?.toString();
    if (!workoutId) return fail(400, { message: 'Workout ID missing.' });

    try {
        const response = await fetch(`${GO_API_URL}/workouts/${workoutId}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${locals.token}` }
        });
        if (response.ok || response.status === 204) {
            return { success: true, deletedId: workoutId, message: "Workout deleted." };
        }
        const errResult = await response.json().catch(() => ({}));
        return fail(response.status, { message: errResult.error || errResult.message || "Failed to delete." });
    } catch (e) {
        return fail(500, { message: "Network error deleting workout." });
    }
  }
};
