import type { RequestEvent } from '@sveltejs/kit';

// Base API URL - Vite/Nginx will proxy /api/ to the backend
const API_BASE_URL = '/api/';

export async function apiRequest(
	endpoint: string,
	options: RequestInit = {},
	event?: RequestEvent // Pass SvelteKit's event object when calling from load/actions
) {
	const url = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint.substring(1) : endpoint}`;

	const headers: HeadersInit = {
		'Content-Type': 'application/json',
		...options.headers
	};


	// Use SvelteKit's fetch if event is provided (for server-side load/actions)
	const fetchFn = event?.fetch || globalThis.fetch;

	try {
		const response = await fetchFn(url, {
			...options,
			headers
		});

		let data = null;
		// Handle 204 No Content, which has no body
		if (response.status === 204) {
			return { response, data };
		}

		// Attempt to parse JSON if content-type suggests it
		const contentType = response.headers.get('content-type');
		if (contentType && contentType.includes('application/json')) {
			data = await response.json();
		} else if (!response.ok) {
			// If not OK and not JSON, try to get text for error message
			const errorText = await response.text();
			console.error(`API Error (${response.status}) for ${url}: ${errorText}`);
			// Re-throw an error that includes the status and message if possible
			throw Object.assign(new Error(errorText || `API Error: ${response.status}`), { status: response.status });
		} else {
            // If OK but not JSON (e.g. plain text response)
            data = await response.text();
        }

		// If response was not OK, but we parsed JSON, the error might be in `data`
		if (!response.ok && data && (data.error || data.message)) {
			throw Object.assign(new Error(data.error || data.message || `API Error: ${response.status}`), { status: response.status, details: data });
		}
		if (!response.ok) { // Fallback error if no JSON error message
			throw Object.assign(new Error(`API Error: ${response.status}`), { status: response.status });
		}


		return { response, data };
	} catch (error: any) {
		// Log and re-throw, or handle more gracefully
		console.error(`API request failed for ${url}:`, error.message, error.details || '');
		// To ensure the error can be caught and handled by SvelteKit's fail mechanism or try/catch blocks
		throw error;
	}
}

// --- Authentication API functions ---
export const authApi = {
	login: async (username: string, password: string, event?: RequestEvent) => {
		return apiRequest('tokens/authentication', { // No leading slash for endpoint
			method: 'POST',
			body: JSON.stringify({ username, password })
		}, event);
	},

	register: async (userData: any, event?: RequestEvent) => {
		return apiRequest('users', { // No leading slash
			method: 'POST',
			body: JSON.stringify(userData)
		}, event);
	}
	// No explicit logout API needed here if using SvelteKit server action for cookie deletion
};

// --- Workout API functions ---
export const workoutApi = {
	getAll: async (event?: RequestEvent) => {
		return apiRequest('workouts/', {}, event); // Added trailing slash if your Go router needs it
	},

	getById: async (id: number, event?: RequestEvent) => {
		return apiRequest(`workouts/${id}`, {}, event);
	},

	create: async (workout: any, event?: RequestEvent) => {
		return apiRequest('workouts/', { // Added trailing slash
			method: 'POST',
			body: JSON.stringify(workout)
		}, event);
	},

	update: async (id: number, workout: any, event?: RequestEvent) => {
		return apiRequest(`workouts/${id}`, {
			method: 'PUT',
			body: JSON.stringify(workout)
		}, event);
	},

	delete: async (id: number, event?: RequestEvent) => {
		return apiRequest(`workouts/${id}`, { // Corrected typo from wokrouts
			method: 'DELETE'
		}, event);
	}
};

// --- Workout Entries API functions (example) ---
export const workoutEntriesApi = {
	getByWorkoutId: async (workoutId: number, event?: RequestEvent) => {
		return apiRequest(`workouts/${workoutId}/entries`, {}, event);
	},

	create: async (workoutId: number, entry: any, event?: RequestEvent) => {
		return apiRequest(`workouts/${workoutId}/entries`, {
			method: 'POST',
			body: JSON.stringify(entry)
		}, event);
	}
};
