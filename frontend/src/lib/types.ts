// general types to be used in the client,
// this are supossed to mirror the go backend entities

export interface BackendUser {
    id: number;
    username: string;
    email: string;
    bio?: string;
}

export interface BackendWorkoutEntry {
    id?: number;
    exercise_name: string;
    sets: number;
    reps?: number | null;
    duration_seconds?: number | null;
    weight?: number | null;
    notes: string;
    order_index: number;
}

export interface BackendWorkout {
    id: number;
    user_id: number;
    title: string;
    description: string;
    duration_minutes: number;
    calories_burned: number;
    entries: BackendWorkoutEntry[];
}
