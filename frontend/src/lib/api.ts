import { API } from '$env/static/private';
import type { Workout, WorkoutEntry } from './types';

// Base API request function with authentication
export async function apiRequest(endpoint: string, options: RequestInit = {}) {
	const url = `${API}${endpoint.startsWith('/') ? endpoint.substring(1) : endpoint}`;

	const headers: HeadersInit = {
		'Content-Type': 'application/json',
		...options.headers
	};

	try {
		const response = await fetch(url, {
			...options,
			headers,
			credentials: 'include' // Important for cookies!
		});

		const data = await response.json();
		return { response, data };
	} catch (error) {
		console.error(`API request error for ${endpoint}:`, error);
		throw error;
	}
}

// Authentication API functions
export const authApi = {
	login: async (username: string, password: string) => {
		return apiRequest('/tokens/authentication', {
			method: 'POST',
			body: JSON.stringify({ username, password })
		});
	},
	register: async (userData: { username: string; password: string; email?: string }) => {
		return apiRequest('/users', {
			method: 'POST',
			body: JSON.stringify(userData)
		});
	}
};

// Workout API functions
export const workoutApi = {
	getAll: async () => {
		return apiRequest('/workouts/');
	},
	getById: async (id: number) => {
		return apiRequest(`/workouts/${id}`);
	},
	create: async (workout: Workout) => {
		return apiRequest('/workouts/', {
			method: 'POST',
			body: JSON.stringify(workout)
		});
	},
	update: async (id: number, workout: Workout) => {
		return apiRequest(`/workouts/${id}`, {
			method: 'PUT',
			body: JSON.stringify(workout)
		});
	},
	delete: async (id: number) => {
		return apiRequest(`/workouts/${id}`, {
			method: 'DELETE'
		});
	}
};

// Workout entries API functions
export const workoutEntriesApi = {
	getByWorkoutId: async (workoutId: number) => {
		return apiRequest(`/workouts/${workoutId}/entries`);
	},
	create: async (workoutId: number, entry: WorkoutEntry) => {
		return apiRequest(`/workouts/${workoutId}/entries`, {
			method: 'POST',
			body: JSON.stringify(entry)
		});
	}
};
