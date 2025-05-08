<script lang="ts">
  import '../app.css'; // Your global styles
  import { page } from '$app/stores';
  import { clientLogout } from '$lib/stores/auth_store'; // We'll create this

  let { children } = $props();
</script>

  <div class="min-h-screen bg-gray-900 text-white">
  <header class="bg-gray-800 shadow-md">
    <nav class="container mx-auto px-4 py-3 flex justify-between items-center">
      <a href="/" class="text-xl font-bold text-blue-400 hover:text-blue-300">Go Workout</a>
      <div>
        {#if $page.data.authenticated}
          <a href="/workouts" class="px-3 py-2 hover:text-blue-300">Workouts</a>
          <a href="/profile" class="px-3 py-2 hover:text-blue-300">Profile</a>
          <form method="POST" action="/logout" class="inline">
            <button type="submit" class="px-3 py-2 hover:text-red-400">Logout</button>
          </form>
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

  <footer class="text-center text-xs text-gray-600 py-4 mt-8 border-t border-gray-700">
    GoWorkout App © {new Date().getFullYear()}
  </footer>
</div>
