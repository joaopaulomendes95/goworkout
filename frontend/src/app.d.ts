// frontend/src/app.d.ts
// See https://svelte.dev/docs/kit/types#app.d.ts

import type { BackendUser, BackendWorkout } from '$lib/types';

declare global {
    namespace App {
        interface Locals {
            authenticated: boolean;
            token?: string;
            user?: BackendUser; // Use imported type
        }
        interface PageData {
            user?: BackendUser; // Use imported type
            authenticated?: boolean;
            workouts?: BackendWorkout[]; // Use imported type
            error?: { message: string };
        }
        // interface Error {}
        // interface Platform {}
    }
}