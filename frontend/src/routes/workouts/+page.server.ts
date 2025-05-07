import { get } from 'svelte/store';
import { token } from '$lib/stores/auth';

const API = 'http://app:8080/workouts';
const ID = '1'; // Replace with the actual ID you want to fetch


type workout = {
    id: number;
    userID: number;
    title: string;
    description: string;
    durationMinutes: number;
    caloriesBurned: number;
    entries: workoutEntry[];
}

type workoutEntry = {
    id: number;
    exerciseName: string;
    sets: number;
    reps: number | null;
    durationSeconds: number | null;
    weight: number | null;
    notes: string;
    oriderIndex: number;
}

export async function apiRequest({ endpoint: string, options: RequestInit = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    const currentToken = get(token);
    if (currentToken) {
        headers['Authorization'] = `Bearer ${currentToken}`;
    }

    const response = await fetch(`${API}${endpoint}`, {
        ...options,
        headers,
        credentials: 'include'
    });

    if (!response.ok) {

    }

    return response;
}

export const load = async () => {
    // fetch workout id
    const response = await fetch(`${API}/workouts/${ID}`);
    const workouts = ((await response.json()) as workout[])
    // const id = response.;
    return {
        workouts: workouts.map
    }
}

export const actions = {
    add_workout: async ({ request }) => {
        const data = await request.formData();
        const id = data.get('id');
        const userID = data.get('user_id');
        const title = data.get('title');
        const description = data.get('description');
        const durationMinutes = data.get('duration_minutes');
        const caloriesBurned = data.get('calorties_burned');
        const entries = data.get('entries');

        //debugggin
        console.log("FormData: ", data)


        try {
            const response = await fetch(`${API}/tokens/authentication`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    id,
                    userID,
                    title,
                    description,
                    durationMinutes,
                    caloriesBurned,
                    entries
                })
            });

            console.log("Response: ", response);
            const result = await response.json();
            console.log("Result: ", result);


            if (response.ok) {
                return {
                    success: true,
                    message: 'Login successful',
                    token: result.token,
                };
            } else {
                return {
                    success: false,
                    message: result.message || 'couldnt get workout',
                    userID, title, description, durationMinutes, caloriesBurned // Return data to repopulate form
                };
            }

        } catch (e) {
            return {
                success: false,
                message: 'Network error',
                userID, title, description, durationMinutes, caloriesBurned // Return data to repopulate form
            };
        }
    }
};
