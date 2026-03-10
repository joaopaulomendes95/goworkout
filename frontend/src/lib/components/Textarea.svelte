<script lang="ts">
	interface Props {
		name: string;
		id?: string;
		label: string;
		placeholder?: string;
		required?: boolean;
		rows?: number;
		value?: string;
		bind:value?: any;
		error?: string;
	}

	let {
		name,
		id,
		label,
		placeholder = '',
		required = false,
		rows = 3,
		value,
		bind:value: bindValue,
		error
	}: Props = $props();
</script>

<div>
	<label for={id || name} class="mb-1 block text-sm font-medium text-gray-300">
		{label}
	</label>
	<textarea
		{name}
		id={id || name}
		{placeholder}
		{required}
		{rows}
		{value}
		bind:value={bindValue}
		class="w-full resize-none rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-blue-500"
	></textarea>
	{#if error}
		<p class="mt-1 text-sm text-red-400">{error}</p>
	{/if}
</div>
