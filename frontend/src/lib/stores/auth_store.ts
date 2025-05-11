import { derived, type Readable } from 'svelte/store';
import { page } from '$app/stores';

// This store reflects the authentication state passed from the server
export const isAuthenticated: Readable<boolean> = derived(
  page,
  ($pageStore) => !!$pageStore.data.authenticated
);

// Client-side logout function
export async function clientLogout() {
  const response = await fetch('/logout', { method: 'POST' });
  if (response.ok && response.redirected) {
    // Redirect to the URL provided by the server
    window.location.href = response.url; 
  } else if (response.ok) {
    // Fallback  if no redirect
    window.location.href = '/login'; 
  } else {
    console.error("Client-side logout attempt failed.");
  }
}
