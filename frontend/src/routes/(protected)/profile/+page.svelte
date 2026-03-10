<script lang="ts">
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let user = $derived(data.user);
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
	{:else if data.authenticated}
		<p>Loading user data...</p>
		<p><strong>Teste With user id:</strong></p>
	{:else}
		<p>You are not logged in.</p>
	{/if}
</div>
