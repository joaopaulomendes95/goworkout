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
            authenticated: boolean;
            user?: BackendUser;
            authenticated?: boolean;
            workouts?: BackendWorkout[]; // Use imported type
            workout?: BackendWorkout; // Use imported type
            error?: { message: string };
            [key: string]: any; // Allow any other properties
        }
        // interface Error {}
        // interface Platform {}
    }
}
export {};
