<script lang="ts">
	import '../app.css'; // Import global styles
	import { page } from '$app/stores';
	import type { LayoutData } from './$types';
	import { clientLogout } from '$lib/stores/auth_store';

	export let data: LayoutData; // Data from +layout.server.ts (user, authenticated, etc.)
	
	// Reactive statement to update user when data changes
	$: user = data.user;
	$: authenticated = data.authenticated;

	async function handleLogout() {
		await clientLogout();
	}
	
	console.log('Layout data:', data);
</script>

<div class="min-h-screen bg-gray-900 text-gray-100 flex flex-col">
	<!-- Navbar -->
	<nav class="bg-gray-800 shadow-md">
		<div class="container mx-auto px-6 py-3 flex justify-between items-center">
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

	<!-- Main Content -->
	<main class="flex-grow container mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<slot />
	</main>

	<!-- Footer -->
	<footer class="bg-gray-800 text-gray-400 text-center p-4 mt-auto">
		<p>© {new Date().getFullYear()} GoWorkout. Your Fitness Journey.</p>
	</footer>
</div>
