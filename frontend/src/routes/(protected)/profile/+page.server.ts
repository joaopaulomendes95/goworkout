import { fail } from '@sveltejs/kit'; 
import type { Actions, PageServerLoad } from './$types';
import type { BackendWorkout, BackendWorkoutEntry } from '$lib/types';

const GO_API_URL = 'http://app:8080';


export const load: PageServerLoad = async ({ locals }) => {
    if (!locals.authenticated || !locals.token) {
    return { workouts: [], error: 'Not authenticated' }; // Or throw redirect
  }

  try {
    const response = await fetch(`${GO_API_URL}/workouts/`, { // Ensure trailing slash matches Go route
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${locals.token}`
      }
    });

    if (!response.ok) {
      const errorResult = await response.json().catch(() => ({ error: `API error (${response.status}) fetching workouts.` }));
      console.error("[Workouts Load] API Error fetching workouts:", response.status, errorResult);
      return {
        workouts: [], // Return empty array on error
        error: errorResult.error || errorResult.message || `Failed to load workouts (status: ${response.status}).`
      };
    }

    // Go API returns { workouts: [...] }
    const data = await response.json(); 
    
    // Ensure the data structure matches what you expect
    if (data && Array.isArray(data.workouts)) {
      return {
        workouts: data.workouts as BackendWorkout[] // Pass the workouts array to the page
      };
    } else {
      console.error("[Workouts Load] Unexpected data structure from API:", data);
      return {
        workouts: [],
        error: "Unexpected data structure from API."
      };
    }

  } catch (e: any) {
    console.error("[Workouts Load] Network or unexpected error:", e);
    return {
      workouts: [],
      error: 'A network error occurred while fetching workouts. Please try again.'
    };
  }
};