<script lang="ts">
	// This page relies on the (protected)/+layout.server.ts for auth check
	// and the root +layout.server.ts to load user data into $page.data.user
	import { page } from '$app/state';
	import type { User } from '$lib/types';

	let user = $derived(page.data.user) as User;
	let authenticated = $derived(page.data.authenticated);
</script>

<div class="space-y-4">
	<h1 class="text-2xl font-semibold">Profile</h1>
	{#if user}
		<div class="rounded-md bg-gray-800 p-4 shadow">
			<p><strong>Username:</strong> {user.username}</p>
			<p><strong>Email:</strong> {user.email}</p>
			<p><strong>Bio:</strong> {user.bio || 'Not set'}</p>
			<p><strong>Teste With user id:</strong> {user.id.toString()}</p>
		</div>
	{:else if authenticated}
		<p>Loading user data...</p>
		<p><strong>Teste With user id:</strong></p>
	{:else}
		<p>You are not logged in.</p>
	{/if}
</div>
