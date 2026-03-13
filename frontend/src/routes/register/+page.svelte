<script lang="ts">
	import { enhance } from '$app/forms';

	interface FormState {
		username?: string;
		email?: string;
		bio?: string;
		error?: string;
	}

	let { form }: { form: FormState | null } = $props();

	let currentUsername = $state('');
	let currentEmail = $state('');
	let currentPassword = $state('');
	let currentBio = $state('');

	$effect(() => {
		if (form?.username) currentUsername = form.username;
		if (form?.email) currentEmail = form.email;
		if (form?.bio) currentBio = form.bio;
	});
</script>

<div
	class="mx-auto mt-10 w-full max-w-md space-y-8 rounded-2xl bg-white p-8 shadow-xl dark:bg-gray-900"
>
	<div class="text-center">
		<h1 class="text-3xl font-bold tracking-tight text-gray-900 dark:text-white">Create Account</h1>
		<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Join Go Svelte Fullstack today!</p>
	</div>

	{#if form?.error}
		<div
			class="mb-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-red-700 transition-all duration-300 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300"
			role="alert"
		>
			<div class="flex items-center">
				<span><strong class="font-medium">Registration Failed:</strong> {form.error}</span>
			</div>
		</div>
	{/if}

	<form method="POST" use:enhance class="space-y-6">
		<div class="space-y-2">
			<label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
				>Username</label
			>
			<div class="relative">
				<input
					bind:value={currentUsername}
					id="username"
					name="username"
					type="text"
					placeholder="Choose a username"
					class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-gray-900 transition-all duration-200 outline-none placeholder:text-gray-400 focus:border-transparent focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:placeholder:text-gray-600"
				/>
			</div>
		</div>

		<div class="space-y-2">
			<label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
				>Email</label
			>
			<div class="relative">
				<input
					bind:value={currentEmail}
					id="email"
					name="email"
					type="email"
					placeholder="you@example.com"
					class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-gray-900 transition-all duration-200 outline-none placeholder:text-gray-400 focus:border-transparent focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:placeholder:text-gray-600"
				/>
			</div>
		</div>

		<div class="space-y-2">
			<label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
				>Password</label
			>
			<div class="relative">
				<input
					bind:value={currentPassword}
					id="password"
					name="password"
					type="password"
					placeholder="Create a strong password"
					class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-gray-900 transition-all duration-200 outline-none placeholder:text-gray-400 focus:border-transparent focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:placeholder:text-gray-600"
				/>
			</div>
		</div>

		<div class="space-y-2">
			<div class="flex items-center justify-between">
				<label for="bio" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
					>Bio</label
				>
				<span class="text-xs text-gray-500 dark:text-gray-400">Optional</span>
			</div>
			<textarea
				bind:value={currentBio}
				id="bio"
				name="bio"
				rows="3"
				placeholder="Tell us a bit about yourself"
				class="w-full resize-none rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-gray-900 transition-all duration-200 outline-none placeholder:text-gray-400 focus:border-transparent focus:ring-2 focus:ring-blue-500 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:placeholder:text-gray-600"
			></textarea>
		</div>

		<button
			type="submit"
			class="flex w-full justify-center rounded-xl border border-transparent bg-blue-600 px-4 py-3 text-sm font-medium text-white shadow-sm transition-all duration-200 hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none"
		>
			Create Account
		</button>
	</form>

	<div class="mt-6">
		<p class="text-center text-sm text-gray-500 dark:text-gray-400">
			Already have an account?
			<a
				href="/login"
				class="font-medium text-blue-600 transition-colors duration-200 hover:text-blue-500 dark:text-blue-400 dark:hover:text-blue-300"
			>
				Log in here
			</a>
		</p>
	</div>
</div>
