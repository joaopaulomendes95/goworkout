<script lang="ts">
	import '../app.css'; // Import global styles
	import { page } from '$app/stores';
	import type { LayoutData } from './$types';
	import { clientLogout } from '$lib/stores/auth_store';
	import Nav from '$lib/components/Nav.svelte';
	import Footer from '$lib/components/Footer.svelte';

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
	<!-- Main Content -->
	<Nav {authenticated} {user} />

	<main class="flex-grow mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<slot />
	</main>

	<Footer />
</div>
