<script lang="ts">
	import '../app.css'; // Your global styles
	import Navbar from '$lib/components/Navbar.svelte';
	import Footer from '$lib/components/Footer.svelte';

	let { children } = $props();
</script>

<div class="min-h-screen bg-gray-900 text-white">
	<Navbar></Navbar>

	<main class="container mx-auto p-4">
		{@render children()}
	</main>

	<Footer></Footer>
</div>
