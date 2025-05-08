<script lang="ts">
	const { data } = $props();

	console.log('page Data', data);

	let health = data.health;
	console.log('health', health);
</script>

<main class="flex min-h-[100vh] flex-col items-center justify-center">
	<h1>Homepage</h1>

	<nav>
		<a href="/login">Login</a>
		<a href="/register">Register</a>
	</nav>
</main>
