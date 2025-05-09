<script lang="ts">
    // This page relies on the (protected)/+layout.server.ts for auth check
    // and the root +layout.server.ts to load user data into $page.data.user
    import { page } from '$app/stores';
    import type { BackendUser } from '../../../app';
	import type { ActionData, PageData } from './$types';

    // User data comes from the form
    let { data, form }: { data: PageData, form: ActionData } = $props();

    let workouts = data.workouts;

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
            <p><strong>Teste With user id:</strong> {user.id.toString()}</p>
        </div>
    {:else if $page.data.authenticated}
        <p>Loading user data...</p>
            <p><strong>Teste With user id:</strong></p>
    {:else}
        <p>You are not logged in.</p> 
    {/if}
</div>
