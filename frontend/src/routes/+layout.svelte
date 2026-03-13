<script lang="ts">
	import '../app.css';
	let { data, children } = $props();
</script>

<div class="min-h-screen bg-gray-900 text-white">
	<header class="bg-gray-800 shadow-md">
		<nav class="container mx-auto flex items-center justify-between px-4 py-3">
			<a href="/" class="text-xl font-bold text-blue-400 hover:text-blue-300">Go Svelte Fullstack</a>
			<div>
				{#if data.authenticated}
					<a href="/protected" class="px-3 py-2 hover:text-blue-300">Dashboard</a>
					<a href="/logout" class="px-3 py-2 hover:text-red-400">Logout</a>
				{:else}
					<a href="/login" class="px-3 py-2 hover:text-blue-300">Login</a>
					<a href="/register" class="px-3 py-2 hover:text-blue-300">Register</a>
				{/if}
			</div>
		</nav>
	</header>

	<main class="container mx-auto p-4">
		{@render children()}
	</main>

	<footer class="mt-8 border-t border-gray-700 py-4 text-center text-xs text-gray-600">
		Go Svelte Fullstack © {new Date().getFullYear()}
	</footer>
</div>
