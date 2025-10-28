export interface User {
	id: number;
	username: string;
	email: string;
	bio?: string;
	createdAt: string;
	updatedAt: string;
}

export interface WorkoutEntry {
	id?: number;
	exercise_name: string;
	sets: number;
	reps?: number | null;
	duration_seconds?: number | null;
	weight?: number | null;
	notes: string;
	order_index: number;
}

export interface Workout {
	id: number;
	user_id: number;
	title: string;
	description: string;
	duration_minutes: number;
	calories_burned: number;
	entries: WorkoutEntry[];
}
