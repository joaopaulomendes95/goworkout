import { derived } from 'svelte/store';
import { page } from '$app/stores';
import type { BackendUser } from '../../app.d'; // Import your user type

// This store reflects the user data passed from the server (via +layout.server.ts)
export const currentUser = derived<typeof page, BackendUser | undefined>(
  page,
  ($pageStore) => $pageStore.data.user
);

// This store reflects the authentication state passed from the server
export const isAuthenticated = derived<typeof page, boolean>(
  page,
  ($pageStore) => !!$pageStore.data.authenticated
);

// Client-side logout function (if you prefer JS-driven logout over a simple form POST)
// This would typically call a SvelteKit action or a dedicated /logout POST endpoint.
export async function clientLogout() {
  const response = await fetch('/logout', { method: 'POST' }); // Calls the SvelteKit action
  if (response.ok && response.redirected) {
    // SvelteKit's form enhancement will usually handle navigation or you can use:
    // import { goto } from '$app/navigation';
    // await goto(response.url); // Navigate to where the server redirected
  } else if (response.ok) {
     // import { invalidateAll } from '$app/navigation';
     // await invalidateAll(); // If no redirect, invalidate to update page data
     // await goto('/login');
  }
  else {
    console.error("Client-side logout attempt failed.");
    // You might want to display an error to the user
  }
}
