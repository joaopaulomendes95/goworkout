import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

// writable store for the token
export const token = writable<string | null>(browser ? localStorage.getItem('token') : null);

if (browser) {
    token.subscribe(value => {
        if (value) {
            localStorage.setItem('token', value);
        } else {
            localStorage.removeItem('token');
        }
    });
}

export const isLoggedIn = derived(token, ($token) => !!$token);

export function setToken(newToken: string | null) {
    token.set(newToken);
}

export function logout(): void {
    token.set(null);
}
