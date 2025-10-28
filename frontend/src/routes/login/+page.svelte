<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/state';

	let { form } = $props();

	const justRegistered = $derived(page.url.search === 'registered=true');

	// For repopulating form fields on error
	let currentUsername = $state(form?.username || '');
	let currentPassword = $state(''); // Don't repopulate password
</script>

<div class="mx-auto mt-10 w-full max-w-md space-y-6 rounded-lg bg-gray-800 p-8 shadow-xl">
	<div class="text-center">
		<h1 class="text-3xl font-bold tracking-tight">Log in</h1>
		<p class="mt-2 text-sm text-gray-400">Access your GoWorkout account</p>
	</div>

	{#if justRegistered}
		<div
			class="relative mb-4 rounded border border-green-600 bg-green-700 px-4 py-3 text-green-100"
			role="alert"
		>
			<strong class="font-bold">Registration successful!</strong>
			<span class="block sm:inline"> Please log in.</span>
		</div>
	{/if}

	{#if form && form.message}
		<div
			class="relative mb-4 rounded border border-red-600 bg-red-700 px-4 py-3 text-red-100"
			role="alert"
		>
			<strong class="font-bold">Login Failed:</strong>
			<span class="block sm:inline"> {form.message}</span>
		</div>
	{/if}

	<form method="POST" action="?/user_login" use:enhance class="space-y-6">
		<div>
			<label for="username" class="block text-sm font-medium text-gray-300">Username</label>
			<input
				bind:value={currentUsername}
				id="username"
				name="username"
				type="text"
				required
				class="mt-1 block w-full rounded-md border border-gray-600 bg-gray-700 px-3 py-2 text-white placeholder-gray-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500 sm:text-sm"
				placeholder="Your username"
			/>
		</div>
		<div>
			<label for="password" class="block text-sm font-medium text-gray-300">Password</label>
			<input
				bind:value={currentPassword}
				id="password"
				name="password"
				type="password"
				required
				class="mt-1 block w-full rounded-md border border-gray-600 bg-gray-700 px-3 py-2 text-white placeholder-gray-500 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500 sm:text-sm"
				placeholder="••••••••"
			/>
		</div>
		<button
			type="submit"
			class="flex w-full justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-800"
		>
			Log in
		</button>
	</form>
	<p class="mt-4 text-center text-sm text-gray-400">
		Don't have an account?
		<a href="/register" class="font-medium text-blue-400 hover:text-blue-300"> Register here </a>
	</p>
</div>
