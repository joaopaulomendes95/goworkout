// See https://svelte.dev/docs/kit/types#app.d.ts

import type { User, Workout } from '$lib/types';

declare global {
	namespace App {
		interface Locals {
			authenticated: boolean;
			token?: string;
			user?: User;
		}
		interface PageData {
			user?: User;
			authenticated?: boolean;
			workouts?: Workout[];
			error?: { message: string };
		}
		// interface Error {}
		// interface Platform {}
	}
}
