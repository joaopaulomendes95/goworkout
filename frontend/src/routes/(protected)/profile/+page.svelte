<script lang="ts">
	import { page } from '$app/stores';

	const {data} = $props();
	console.log('Profile data: ', data);

	let username = $state(data.user?.user.username);
	let bio = $state(data.user?.user.bio);

	$effect(() => {
		console.log('Username: ', username);
	});

</script>

<div class="space-y-6 p-4 md:p-6">
	<h1 class="text-3xl font-bold text-white">Profile</h1>

	{#if $page.data.authenticated}
		<div class="rounded-lg bg-gray-800 p-6 shadow-md">
			<p class="text-lg text-white">
				You are logged in as {data.user?.user.username}.
			</p>
			<p class="mt-2 text-gray-400">
				This is your profile page. Your GoWorkout journey starts here!
			</p>
			<form method="POST" action="?/update_profile" class="mt-4">
				<label for="username" class="mt-4 block text-sm font-medium text-gray-300">
					Username:</label>
					<input
					 	name="username"
						type="text"
						bind:value={username}
						class="mt-1 block w-full rounded-md border-gray-300 bg-gray-700 text-white shadow-sm focus:border-blue-500 focus:ring-blue-500"
					/>
				<label for="bio" class="mt-4 block text-sm font-medium text-gray-300">
					Bio:</label>
					<textarea
					 	name="bio"
						placeholder="Tell us about yourself"
						bind:value={bio}
						class="mt-1 block w-full rounded-md border-gray-300 bg-gray-700 text-white shadow-sm focus:border-blue-500 focus:ring-blue-500"></textarea>
				<button type="submit" class="rounded bg-red-600 px-4 py-2 text-white hover:bg-red-700">
					Update
				</button>
			</form>
			<p class="mt-4">
				<a href="/workouts" class="text-blue-400 hover:text-blue-300">
					Manage your Workouts →
				</a>
			</p>
		</div>
	{:else}
		<p class="text-red-400">You are not logged in. Please <a href="/login" class="underline">login</a>.</p>
	{/if}
</div>
