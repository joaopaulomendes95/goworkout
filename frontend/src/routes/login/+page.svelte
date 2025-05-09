<script lang="ts">
  import { enhance } from '$app/forms';
  import { page } from '$app/stores';

  let { form } = $props();

  const justRegistered = $derived($page.url.searchParams.get('registered') === 'true');
  
  // For repopulating form fields on error
  let currentUsername = $state(form?.username || '');
  let currentPassword = $state(''); // Don't repopulate password

  const inputClasses = "mt-1 block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md shadow-sm placeholder-gray-500 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm text-white";
  const labelClasses = "block text-sm font-medium text-gray-300";
  const buttonClasses = "w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-blue-500";
</script>

<div class="w-full max-w-md p-8 space-y-6 bg-gray-800 rounded-lg shadow-xl mx-auto mt-10">
    <div class="text-center">
        <h1 class="text-3xl font-bold tracking-tight">Log in</h1>
        <p class="mt-2 text-sm text-gray-400">Access your GoWorkout account</p>
    </div>

    {#if justRegistered}
      <div class="bg-green-700 border border-green-600 text-green-100 px-4 py-3 rounded relative mb-4" role="alert">
        <strong class="font-bold">Registration successful!</strong>
        <span class="block sm:inline"> Please log in.</span>
      </div>
    {/if}

    {#if form && form.message }
      <div class="bg-red-700 border border-red-600 text-red-100 px-4 py-3 rounded relative mb-4" role="alert">
        <strong class="font-bold">Login Failed:</strong>
        <span class="block sm:inline"> {form.message}</span>
      </div>
    {/if}

    <form method="POST" action="?/user_login" use:enhance class="space-y-6">
        <div>
            <label for="username" class={labelClasses}>Username</label>
            <input
                bind:value={currentUsername}
                id="username"
                name="username"
                type="text"
                required
                class={inputClasses}
                placeholder="Your username"
            />
        </div>
        <div>
            <label for="password" class={labelClasses}>Password</label>
            <input
                bind:value={currentPassword}
                id="password"
                name="password"
                type="password"
                required
                class={inputClasses}
                placeholder="••••••••"
            />
        </div>
        <button type="submit" class={buttonClasses}>
            Log in
        </button>
    </form>
    <p class="mt-4 text-center text-sm text-gray-400">
        Don't have an account?
        <a href="/register" class="font-medium text-blue-400 hover:text-blue-300">
            Register here
        </a>
    </p>
</div>
