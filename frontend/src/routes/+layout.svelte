<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';

	// Svelte 5: $props() is implicitly available for top-level component data.
	// The default slot for children is handled by <slot />.
</script>

<div class="min-h-screen bg-gray-900 text-white">
	<header class="bg-gray-800 shadow-md">
		<nav class="container mx-auto flex items-center justify-between px-4 py-3">
			<a href="/" class="text-xl font-bold text-blue-400 hover:text-blue-300">Go Workout</a>
			<div>
				{#if $page.data.authenticated}
					<a href="/workouts" class="px-3 py-2 hover:text-blue-300">Workouts</a>
					<a href="/profile" class="px-3 py-2 hover:text-blue-300">Profile</a>
					<form method="POST" action="/logout" class="inline">
						<button type="submit" class="px-3 py-2 hover:text-red-400">Logout</button>
					</form>
				{:else}
					<a href="/login" class="px-3 py-2 hover:text-blue-300">Login</a>
					<a href="/register" class="px-3 py-2 hover:text-blue-300">Register</a>
				{/if}
			</div>
		</nav>
	</header>

	<main class="container mx-auto p-4">
		<slot /> <!-- Correct way to render child content (default slot) -->
	</main>

	<footer class="mt-8 border-t border-gray-700 py-4 text-center text-xs text-gray-600">
		GoWorkout App © {new Date().getFullYear()}
	</footer>
</div>
