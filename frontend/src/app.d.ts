// frontend/src/app.d.ts
// See https://svelte.dev/docs/kit/types#app.d.ts

import type { User } from '$lib/types';

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
		}
	}
}

export {};
