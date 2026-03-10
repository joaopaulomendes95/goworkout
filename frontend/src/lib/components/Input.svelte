<script lang="ts">
	interface Props {
		type?: 'text' | 'email' | 'password' | 'number';
		name: string;
		id?: string;
		label: string;
		placeholder?: string;
		required?: boolean;
		min?: number;
		step?: string;
		value?: string | number;
		bind:value?: any;
		oninput?: (e: Event) => void;
		disabled?: boolean;
		error?: string;
	}

	let {
		type = 'text',
		name,
		id,
		label,
		placeholder = '',
		required = false,
		min,
		step,
		value,
		bind:value: bindValue,
		oninput,
		disabled = false,
		error
	}: Props = $props();
</script>

<div>
	<label for={id || name} class="mb-1 block text-sm font-medium text-gray-300">
		{label}
	</label>
	<input
		{type}
		{name}
		id={id || name}
		{placeholder}
		{required}
		{min}
		{step}
		{value}
		bind:value={bindValue}
		oninput={oninput}
		{disabled}
		class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
	/>
	{#if error}
		<p class="mt-1 text-sm text-red-400">{error}</p>
	{/if}
</div>
