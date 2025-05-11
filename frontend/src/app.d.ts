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
        // interface Error {}
        // interface Platform {}
    }
}
export {};
