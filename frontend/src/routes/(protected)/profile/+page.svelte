<script lang="ts">
	import { page } from '$app/stores';
	// import type { BackendUser } from '$lib/types'; // Not strictly needed if user is always undefined here

	// $page.data.user will be undefined from the root layout in "no backend changes"
	// const user = $derived($page.data.user as BackendUser | undefined);
</script>

<div class="space-y-6 p-4 md:p-6">
	<h1 class="text-3xl font-bold text-white">Profile</h1>

	{#if $page.data.authenticated}
		<div class="rounded-lg bg-gray-800 p-6 shadow-md">
			<p class="text-lg text-white">You are logged in.</p>
			<p class="mt-2 text-gray-400">
				This is your profile page. Your GoWorkout journey starts here!
			</p>
			<p class="mt-4">
				<a href="/workouts" class="text-blue-400 hover:text-blue-300">
					Manage your Workouts →
				</a>
			</p>
		</div>
	{:else}
		<!-- This case should ideally not be reached due to (protected) layout -->
		<p class="text-red-400">You are not logged in. Please <a href="/login" class="underline">login</a>.</p>
	{/if}
</div>
