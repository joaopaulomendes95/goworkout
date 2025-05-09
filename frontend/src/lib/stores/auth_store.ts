import { derived, type Readable } from 'svelte/store';
import { page } from '$app/stores';

// This store reflects the authentication state passed from the server
export const isAuthenticated = Readable<boolean> = derived(
  page,
  ($pageStore) => !!$pageStore.data.authenticated
);

// Client-side logout function
export async function clientLogout() {
  const response = await fetch('/logout', { method: 'POST' });
  if (response.ok && response.redirected) {
    window.location.href = response.url; // Redirect to the URL provided by the server
  } else if (response.ok) {
    window.location.href = '/login'; // Fallback  if no redirect
  } else {
    console.error("Client-side logout attempt failed.");
  }
}
