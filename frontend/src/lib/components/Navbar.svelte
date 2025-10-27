<script lang="ts">
	import { page } from '$app/state';
	let authenticated = $state(page.data.authenticated);
</script>

<header class="bg-gray-800 shadow-md">
	<nav class="container mx-auto flex items-center justify-between px-4 py-3">
		<a href="/" class="text-xl font-bold text-blue-400 hover:text-blue-300">Go Workout</a>
		<div>
			{#if authenticated}
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
