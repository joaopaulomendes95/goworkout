<script lang="ts">
	import { clientLogout } from '$lib/stores/auth_store';
  import type { BackendUser } from '$lib/types';

  export let authenticated: boolean;
  export let user: BackendUser;

	async function handleLogout() {
		await clientLogout();
	}

</script>

	<nav class="bg-gray-800 shadow-md">
		<div class=" mx-auto px-6 py-3 flex justify-between items-center">
			<a href="/" class="text-2xl font-bold text-indigo-400 hover:text-indigo-300">
				Go Workout
			</a>
			<div class="space-x-4">
				<a href="/" class="hover:text-indigo-300">Home</a>
				{#if authenticated && user}
					<a href="/workouts" class="hover:text-indigo-300">My Workouts</a>
					<a href="/profile" class="hover:text-indigo-300">Profile ({user.user.username})</a>
					<button on:click={handleLogout} class="bg-red-500 hover:bg-red-600 text-white py-2 px-4 rounded-md text-sm">
						Logout
					</button>
				{:else}
					<a href="/login" class="hover:text-indigo-300">Login</a>
					<a href="/register" class="bg-indigo-500 hover:bg-indigo-600 text-white py-2 px-4 rounded-md text-sm">
						Register
					</a>
				{/if}
			</div>
		</div>
	</nav>
