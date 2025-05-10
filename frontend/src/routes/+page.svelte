<script lang="ts">
	import type { PageData } from './$types';

	const { data } = $props<{ data: PageData }>();

	const health = $derived(data.health);
	const authenticated = $derived(data.authenticated);

	console.log('Page data:', data);
</script>

<main class="flex min-h-[calc(100vh-150px)] flex-col items-center justify-center text-center">
	<h1 class="mb-6 text-4xl font-bold">Welcome to GoWorkout!</h1>
	<h1 class="mb-6 text-4xl font-bold">{ data.user.user.username }</h1>

	{#if health}
		<div class="mb-6 rounded-lg bg-gray-800 p-4 shadow">
			<h2 class="mb-2 text-xl font-semibold text-blue-400">Backend Health Status</h2>
			{#if health.status === 'up'}
				<p class="text-green-400">Backend is healthy and running!</p>
				{#if health.message}<p class="text-xs text-gray-500">Message: {health.message}</p>{/if}
			{:else}
				<p class="text-red-400">Backend is down or unreachable.</p>
				{#if health.error}
					<p class="mt-1 text-xs text-red-500">Details: {health.error}</p>
				{/if}
			{/if}
		</div>
	{:else}
		<p class="text-gray-500">Checking backend health...</p>
	{/if}

	<p class="mb-8 text-lg text-gray-300">
		Track your fitness journey, log workouts, and achieve your goals.
	</p>

	<nav class="space-x-4">
		{#if authenticated}
			<a href="/workouts" class="rounded bg-blue-600 px-6 py-3 text-white hover:bg-blue-700">
				Go to Workouts
			</a>
		{:else}
			<a href="/login" class="rounded bg-blue-600 px-6 py-3 text-white hover:bg-blue-700">
				Login
			</a>
			<a href="/register" class="rounded bg-gray-600 px-6 py-3 text-white hover:bg-gray-700">
				Register
			</a>
		{/if}
	</nav>
</main>
