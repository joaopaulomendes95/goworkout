<script lang="ts">
	import { enhance } from '$app/forms';
	let { form } = $props();

	let currentUsername = $state(form?.username || '');
	let currentEmail = $state(form?.email || '');
	let currentPassword = $state('');
	let currentBio = $state(form?.bio || '');


  const inputClasses = "mt-1 block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md shadow-sm placeholder-gray-500 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm text-white";
  const labelClasses = "block text-sm font-medium text-gray-300";
  const buttonClasses = "w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-blue-500";
</script>

<div
	class="mx-auto mt-10 w-full max-w-md space-y-8 rounded-2xl bg-white p-8 shadow-xl dark:bg-gray-800"
>
	<div class="text-center">
		<h1 class="text-3xl font-bold tracking-tight text-gray-900 dark:text-white">Create Account</h1>
		<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Join GoWorkout today!</p>
	</div>

	{#if form && form.message}
		<div
			class="mb-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-red-700 transition-all duration-300 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300"
			role="alert"
		>
			<div class="flex items-center">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mr-2 h-5 w-5"
					viewBox="0 0 20 20"
					fill="currentColor"
				>
					<path
						fill-rule="evenodd"
						d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
						clip-rule="evenodd"
					/>
				</svg>
				<span><strong class="font-medium">Registration Failed:</strong> {form.message}</span>
			</div>
		</div>
	{/if}

	<form method="POST" action="?/user_register" class="space-y-6">
		<div class="space-y-2">
			<label for="username" class={labelClasses}
				>Username</label
			>
			<div class="relative">
				<input
					bind:value={currentUsername}
					id="username"
					name="username"
					type="text"
					placeholder="Choose a username"
					class={inputClasses}
				/>
			</div>
		</div>

		<div class="space-y-2">
			<label for="email" class={labelClasses}
				>Email</label
			>
			<div class="relative">
				<input
					bind:value={currentEmail}
					id="email"
					name="email"
					type="email"
					placeholder="you@example.com"
					class={inputClasses}
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
