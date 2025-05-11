// frontend/src/app.d.ts
// See https://svelte.dev/docs/kit/types#app.d.ts

import type { BackendUser, BackendWorkout } from '$lib/types';

declare global {
    namespace App {
        interface Locals {
            authenticated: boolean;
            token?: string;
            user?: BackendUser;
        }
        interface PageData {
            authenticated: boolean;
            user?: BackendUser;
            workouts?: BackendWorkout[];
            workout?: BackendWorkout;
            error?: { message: string };
            [key: string]: any;
        }
        // interface Error {}
        // interface Platform {}
    }
}
export {};
