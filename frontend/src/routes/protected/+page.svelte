<script lang="ts">
	let { data } = $props();
</script>

<div class="mx-auto mt-10 max-w-2xl">
	<h1 class="mb-4 text-2xl font-bold">Dashboard</h1>
	<p class="text-gray-400">Welcome! You're logged in {data.user?.username} .</p>
</div>
