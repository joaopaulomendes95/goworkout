<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		type?: 'success' | 'error' | 'warning' | 'info';
		children: Snippet;
	}

	let { type = 'info', children }: Props = $props();

	const types = {
		success: 'border-green-600 bg-green-700 text-green-100',
		error: 'border-red-600 bg-red-700 text-red-100',
		warning: 'border-yellow-600 bg-yellow-700 text-yellow-100',
		info: 'border-blue-600 bg-blue-700 text-blue-100'
	};
</script>

<div class="mb-4 rounded border px-4 py-3 {types[type]}" role="alert">
	{@render children()}
</div>
