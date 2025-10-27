<script lang="ts">
	let { data } = $props();
	let workout = $derived(data);
</script>

<div
	class="w-full max-w-sm rounded-lg border border-gray-200 bg-white p-4 shadow-sm sm:p-6 dark:border-gray-700 dark:bg-gray-800"
>
	<h5 class="mb-3 text-base font-semibold text-gray-900 md:text-xl dark:text-white">
		{workout.data.}
	</h5>
	<p class="text-sm font-normal text-gray-500 dark:text-gray-400">Information about the workout</p>
	<ul class="my-4 space-y-3">
		<li>
			<a
				href="#"
				class="group flex items-center rounded-lg bg-gray-50 p-3 text-base font-bold text-gray-900 hover:bg-gray-100 hover:shadow dark:bg-gray-600 dark:text-white dark:hover:bg-gray-500"
			>
				<span class="ms-3 flex-1 whitespace-nowrap">MetaMask</span>
				<span
					class="ms-3 inline-flex items-center justify-center rounded-sm bg-gray-200 px-2 py-0.5 text-xs font-medium text-gray-500 dark:bg-gray-700 dark:text-gray-400"
					>Popular</span
				>
			</a>
		</li>
		<li>
			<a
				href="#"
				class="group flex items-center rounded-lg bg-gray-50 p-3 text-base font-bold text-gray-900 hover:bg-gray-100 hover:shadow dark:bg-gray-600 dark:text-white dark:hover:bg-gray-500"
			>
				>
				<span class="ms-3 flex-1 whitespace-nowrap">Coinbase Wallet</span>
			</a>
		</li>
		<li>
			<a
				href="#"
				class="group flex items-center rounded-lg bg-gray-50 p-3 text-base font-bold text-gray-900 hover:bg-gray-100 hover:shadow dark:bg-gray-600 dark:text-white dark:hover:bg-gray-500"
			>
				<span class="ms-3 flex-1 whitespace-nowrap">Opera Wallet</span>
			</a>
		</li>
		<li>
			<a
				href="#"
				class="group flex items-center rounded-lg bg-gray-50 p-3 text-base font-bold text-gray-900 hover:bg-gray-100 hover:shadow dark:bg-gray-600 dark:text-white dark:hover:bg-gray-500"
			>
				<span class="ms-3 flex-1 whitespace-nowrap">WalletConnect</span>
			</a>
		</li>
		<li>
			<a
				href="#"
				class="group flex items-center rounded-lg bg-gray-50 p-3 text-base font-bold text-gray-900 hover:bg-gray-100 hover:shadow dark:bg-gray-600 dark:text-white dark:hover:bg-gray-500"
			>
				<span class="ms-3 flex-1 whitespace-nowrap">Fortmatic</span>
			</a>
		</li>
	</ul>
	<div>
		<a
			href="#"
			class="inline-flex items-center text-xs font-normal text-gray-500 hover:underline dark:text-gray-400"
		>
			Why do I need to connect with my wallet?</a
		>
	</div>
</div>
