import type { PageServerLoad, PageServerLoadEvent } from './$types';
import type { BackendWorkout } from '$lib/types';

export const load: PageServerLoad = async (event: PageServerLoadEvent) => {
    // Authentication is handled by the (protected)/+layout.server.ts
    // User data is available from the root +layout.server.ts via await event.parent() or in $page.data
    // We only need to fetch profile-specific data here, like workouts for this example.

    if (!event.locals.authenticated || !event.locals.token) {
        // This should be caught by (protected) layout, but good for robustness
        return { workouts: [], error: 'Not authenticated', user: event.locals.user };
    }

    try {
        const response = await event.fetch(`/api/workouts/`, { // Relative path
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${event.locals.token}` // Explicitly sending token
            }
        });

        if (!response.ok) {
            const errorResult = await response.json().catch(() => ({ error: `API error (${response.status}) fetching workouts.` }));
            console.error("[Profile Load] API Error fetching workouts:", response.status, errorResult);
            return {
                workouts: [],
                error: errorResult.error || errorResult.message || `Failed to load workouts (status: ${response.status}).`,
                user: event.locals.user // Pass user from locals
            };
        }

        const data = await response.json();
        if (data && Array.isArray(data.workouts)) {
            return {
                workouts: data.workouts as BackendWorkout[],
                user: event.locals.user // Pass user from locals
            };
        } else {
            console.error("[Profile Load] Unexpected data structure from API:", data);
            return {
                workouts: [],
                error: "Unexpected data structure from API.",
                user: event.locals.user
            };
        }
    } catch (e: any) {
        console.error("[Profile Load] Network or unexpected error:", e.message);
        return {
            workouts: [],
            error: 'A network error occurred while fetching workouts.',
            user: event.locals.user
        };
    }
};
