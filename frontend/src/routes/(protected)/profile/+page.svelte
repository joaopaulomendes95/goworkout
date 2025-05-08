<script lang="ts">
  // This page relies on the (protected)/+layout.server.ts for auth check
  // and the root +layout.server.ts to load user data into $page.data.user
  import { page } from '$app/stores';
	import type { BackendUser } from '../../../app';

  // User data comes from the layout
  let user = $derived(($page.data.user as BackendUser | undefined));
</script>

<div class="space-y-4">
    <h1 class="text-2xl font-semibold">Profile</h1>
    {#if user}
        <div class="p-4 bg-gray-800 rounded-md shadow">
            <p><strong>Username:</strong> {user.username}</p>
            <p><strong>Email:</strong> {user.email}</p>
            <p><strong>Bio:</strong> {user.bio || 'Not set'}</p>
        </div>
    {:else if $page.data.authenticated}
        <p>Loading user data...</p>
    {:else}
        <p>You are not logged in.</p> 
    {/if}
</div>
