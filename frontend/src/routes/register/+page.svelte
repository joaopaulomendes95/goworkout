<script lang="ts">
  import { enhance } from '$app/forms';
  let { form } = $props();

  let currentUsername = $state(form?.username || '');
  let currentEmail = $state(form?.email || '');
  let currentPassword = $state('');
  let currentBio = $state(form?.bio || '');

</script>

<div class="w-full max-w-md p-8 space-y-6 bg-gray-800 rounded-lg shadow-xl mx-auto mt-10">
    <div class="text-center">
        <h1 class="text-3xl font-bold tracking-tight">Create Account</h1>
        <p class="mt-2 text-sm text-gray-400">Join GoWorkout today!</p>
    </div>

    {#if form && form.message }
      <div class="bg-red-700 border border-red-600 text-red-100 px-4 py-3 rounded relative mb-4" role="alert">
        <strong class="font-bold">Registration Failed:</strong>
        <span class="block sm:inline"> {form.message}</span>
      </div>
    {/if}

    <form method="POST" use:enhance class="space-y-6">
        <div>
            <label for="username">Username</label>
            <input bind:value={currentUsername} id="username" name="username" type="text" placeholder="Choose a username"/>
        </div>
        <div>
            <label for="email">Email</label>
            <input bind:value={currentEmail} id="email" name="email" type="email" placeholder="you@example.com"/>
        </div>
        <div>
            <label for="password">Password</label>
            <input bind:value={currentPassword} id="password" name="password" type="password" placeholder="Create a strong password"/>
        </div>
        <div>
            <label for="bio">Bio <span class="text-xs text-gray-500">(Optional)</span></label>
            <textarea bind:value={currentBio} id="bio" name="bio" rows="3" placeholder="Tell us a bit about yourself"></textarea>
        </div>
        <button type="submit">
            Register
        </button>
    </form>
    <p class="mt-4 text-center text-sm text-gray-400">
        Already have an account?
        <a href="/login" class="font-medium text-blue-400 hover:text-blue-300">
            Log in here
        </a>
    </p>
</div>
