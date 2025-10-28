<script lang="ts">
	import Workout from '$lib/components/Workout.svelte';

	const { data } = $props();

	console.log('page Data', data);

	let health = data.health;
	console.log('health', health);
</script>

<main class="flex min-h-screen flex-col items-center justify-center">
	<h1>Homepage</h1>
	<Workout {data} />

	<nav>
		<a href="/login">Login</a>
		<a href="/register">Register</a>
	</nav>
</main>
