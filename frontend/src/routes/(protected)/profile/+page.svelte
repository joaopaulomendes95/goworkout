<script lang="ts">
    import { page } from '$app/stores';
    import type { BackendUser } from '$lib/types'; // Corrected import if app.d.ts is path aliased
    import type { ActionData, PageData } from './$types';
    import { onMount } from 'svelte';
    import { apiRequest } from '$lib/api'; // Assuming apiRequest is set up

    let { data, form }: { data: PageData, form: ActionData } = $props();
    let workouts = $derived(data.workouts); // From server load

    // User data will be initially undefined from $page.data.user if hooks doesn't fetch it.
    // We can try a client-side fetch if absolutely necessary, but it's not ideal.
    let clientFetchedUser = $state<BackendUser | undefined>(undefined);
    let isLoadingUser = $state(false);

    // Try to get user from $page.data first (if hooks somehow got it)
    let user = $derived(($page.data.user as BackendUser | undefined) ?? clientFetchedUser);

    // This demonstrates a client-side fetch for user, which is fallback behavior
    // if you don't have a /users/me called in hooks.server.ts.
    // It's generally better to load essential data server-side via hooks or page/layout load functions.
    onMount(async () => {
        if (!$page.data.user && $page.data.authenticated && !clientFetchedUser) {
            isLoadingUser = true;
            try {
                // THIS STILL REQUIRES AN ENDPOINT LIKE /api/users/me
                // If you don't have it, you can't fetch user details client-side either
                // based on the HttpOnly cookie alone.
                // For this example to work, we assume /api/users/me exists.
                // If not, you can only display "Authenticated" but not user details.
                // const { data: userDataWrapper } = await apiRequest('users/me');
                // if (userDataWrapper && userDataWrapper.user) {
                // clientFetchedUser = userDataWrapper.user;
                // }
                // ELSE, if no /users/me, you can't get user details here.
                // You would remove this onMount block for user fetching.
            } catch (error) {
                console.error("Failed to fetch user details client-side:", error);
            } finally {
                isLoadingUser = false;
            }
        }
    });
</script>

<div class="space-y-4">
    <h1 class="text-2xl font-semibold">Profile</h1>
    {#if $page.data.authenticated}
        {#if user}
            <div class="p-4 bg-gray-800 rounded-md shadow">
                <p><strong>Username:</strong> {user.username}</p>
                <p><strong>Email:</strong> {user.email}</p>
                <p><strong>Bio:</strong> {user.bio || 'Not set'}</p>
                <p><strong>User ID:</strong> {user.id}</p>
            </div>
        {:else if isLoadingUser}
            <p>Loading user data...</p>
        {:else}
            <p>User details not available (user might be authenticated, but details couldn't be fetched).</p>
        {/if}

        <h2 class="text-xl font-semibold mt-6">Your Workouts</h2>
        {#if data.error}
            <p class="text-red-400">{data.error}</p>
        {:else if workouts && workouts.length > 0}
            <ul class="space-y-2">
                {#each workouts as workout (workout.id)}
                    <li class="p-3 bg-gray-700 rounded">{workout.title} - {workout.description}</li>
                {/each}
            </ul>
        {:else}
            <p>No workouts found.</p>
        {/if}
    {:else}
        <p>You are not logged in. Please <a href="/login" class="text-blue-400 hover:underline">log in</a>.</p>
    {/if}
</div>
