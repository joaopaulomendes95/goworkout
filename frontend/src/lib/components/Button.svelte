<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		type?: 'submit' | 'button' | 'reset';
		variant?: 'primary' | 'secondary' | 'danger';
		disabled?: boolean;
		class?: string;
		onclick?: (e: MouseEvent) => void;
		children: Snippet;
	}

	let {
		type = 'button',
		variant = 'primary',
		disabled = false,
		class: className = '',
		onclick,
		children
	}: Props = $props();

	const variants = {
		primary: 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500',
		secondary: 'bg-gray-600 hover:bg-gray-700 focus:ring-gray-500',
		danger: 'bg-red-600 hover:bg-red-700 focus:ring-red-500'
	};
</script>

<button
	{type}
	{disabled}
	{onclick}
	class="rounded-md px-4 py-2 font-medium text-white focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 disabled:opacity-50 {variants[variant]} {className}"
>
	{@render children()}
</button>
