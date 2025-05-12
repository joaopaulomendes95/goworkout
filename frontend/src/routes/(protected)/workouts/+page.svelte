<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData, PageData } from './$types';
	import type { BackendWorkoutEntry } from '$lib/types';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { fade, fly, scale, slide } from 'svelte/transition';
	import { quintOut, backOut, elasticOut } from 'svelte/easing';
	import { flip } from 'svelte/animate';

	const { data, form } = $props<{ data: PageData; form?: ActionData }>();

	const workouts = $derived(data.workouts || []);
	const pageLoadError = $derived(
		data.error ? (typeof data.error === 'string' ? data.error : data.error.message) : undefined
	);

	// Enhanced action feedback with more details and auto-dismissal
	let actionFeedback = $state<{ type: 'success' | 'error'; message: string } | undefined>(
		undefined
	);

	// State for UI interactions
	let activeTab = $state<'log' | 'add'>('log'); // Default to showing workout log
	let expandedWorkoutId = $state<string | null>(null);
	let formProcessed = $state(false);
	let showStats = $state(true);
	let workoutToDelete = $state<string | null>(null); // Track which workout is being deleted

	// Form state for adding a new workout
	let newWorkoutTitle = $state('');
	let newWorkoutDescription = $state('');
	let newWorkoutDuration = $state(0);
	let newWorkoutCalories = $state(0);
	let newWorkoutEntries = $state<Partial<BackendWorkoutEntry>[]>([
		{
			exercise_name: '',
			sets: 1,
			reps: null,
			duration_seconds: null,
			notes: '',
			order_index: 1,
			weight: null
		}
	]);

	// Calculated stats for the dashboard
	const stats = $derived(
		workouts.length === 0
			? { totalWorkouts: 0, totalMinutes: 0, totalCalories: 0, avgDuration: 0 }
			: {
					totalWorkouts: workouts.length,
					totalMinutes: workouts.reduce((sum, w) => sum + (w.duration_minutes || 0), 0),
					totalCalories: workouts.reduce((sum, w) => sum + (w.calories_burned || 0), 0),
					avgDuration: Math.round(
						workouts.reduce((sum, w) => sum + (w.duration_minutes || 0), 0) / workouts.length
					)
				}
	);

	// Handling form repopulation when there's an error
	$effect(() => {
		if (form && form.formError && form.title !== undefined) {
			newWorkoutTitle = form.title || '';
			newWorkoutDescription = form.description || '';
			newWorkoutDuration = form.durationMinutes ? parseInt(form.durationMinutes) : 0;
			newWorkoutCalories = form.caloriesBurned ? parseInt(form.caloriesBurned) : 0;
			if (form.entries && typeof form.entries === 'string') {
				try {
					const parsedEntries = JSON.parse(form.entries);
					newWorkoutEntries =
						Array.isArray(parsedEntries) && parsedEntries.length > 0
							? parsedEntries
							: [createEmptyEntry(0)];
				} catch (e) {
					newWorkoutEntries = [createEmptyEntry(0)];
				}
			}
		}
	});

	// handling redirect messages from the server
	$effect(() => {
		if (browser) {
			const messageFromRedirect = $page.url.searchParams.get('message');
			if (messageFromRedirect) {
				actionFeedback = { type: 'success', message: messageFromRedirect.replace(/_/g, ' ') };
				const newUrl = new URL(window.location.href);
				newUrl.searchParams.delete('message');
				window.history.replaceState(window.history.state, '', newUrl);
			}
		}
	});

	// handling form feedback and show appropriate toast notifications
	$effect(() => {
		const currentForm = form;
		if (currentForm && !formProcessed) {
			formProcessed = true;
			if (currentForm.success && currentForm.message) {
				actionFeedback = { type: 'success', message: currentForm.message };
				if (currentForm.message?.includes('added successfully')) {
					resetForm();
					// Switch to log tab to see the new workout
					activeTab = 'log';
				}
				if (browser) invalidateAll();
			} else if (currentForm.formError) {
				actionFeedback = { type: 'error', message: currentForm.formError };
			} else if (currentForm.message && !currentForm.success) {
				actionFeedback = { type: 'error', message: currentForm.message };
			}
			if (actionFeedback) {
				const timer = setTimeout(() => {
					actionFeedback = undefined;
				}, 5000);
				return () => clearTimeout(timer);
			}
		} else if (!currentForm && formProcessed) {
			formProcessed = false;
		}
	});

	// creating a template for empty exercise entries
	function createEmptyEntry(index: number) {
		return {
			exercise_name: '',
			sets: 1,
			reps: null,
			duration_seconds: null,
			notes: '',
			order_index: index + 1,
			weight: null
		};
	}

	// Reset form to initial state
	function resetForm() {
		newWorkoutTitle = '';
		newWorkoutDescription = '';
		newWorkoutDuration = 0;
		newWorkoutCalories = 0;
		newWorkoutEntries = [createEmptyEntry(0)];
	}

	// Adding new exercise rows to the workout form
	function addEntryRow() {
		newWorkoutEntries = [...newWorkoutEntries, createEmptyEntry(newWorkoutEntries.length)];
	}

	// Removing exercise rows and update indices
	function removeEntryRow(index: number) {
		newWorkoutEntries = newWorkoutEntries.filter((_, i) => i !== index);
		// Update order indices
		newWorkoutEntries = newWorkoutEntries.map((e, i) => ({ ...e, order_index: i + 1 }));
	}

    // Handling input changes for reps
	function handleRepsInput(entry: Partial<BackendWorkoutEntry>) {
		if (entry.reps && entry.reps > 0) entry.duration_seconds = null;
	}

    // Handling input changes for duration
	function handleDurationInput(entry: Partial<BackendWorkoutEntry>) {
		if (entry.duration_seconds && entry.duration_seconds > 0) entry.reps = null;
	}

	// Toggle workout details expansion
	function toggleWorkoutDetails(id: string) {
		expandedWorkoutId = expandedWorkoutId === id ? null : id;
	}

	// Managing the delete confirmation flow
	function confirmDelete(id: string) {
		workoutToDelete = id;
	}

    // If the user cancels the delete action
    // reset workoutToDelete to null
	function cancelDelete() {
		workoutToDelete = null;
	}

	// Form submission with animation
	function handleSubmit() {
		formProcessed = false;
	}

	// Calculate workout intensity based on calories and duration
	function getIntensity(calories: number, duration: number): 'low' | 'medium' | 'high' {
		const caloriesPerMinute = calories / duration;
		if (caloriesPerMinute < 5) return 'low';
		if (caloriesPerMinute < 10) return 'medium';
		return 'high';
	}

	// Format date from ISO string
	function formatDate(dateString: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-PT', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	// UI utility classes with modern design system
	const baseInputClasses =
		'w-full rounded-lg border bg-gray-800/90 px-4 py-3 text-white transition duration-200 placeholder:text-gray-400 focus:border-orange-500 focus:ring-2 focus:ring-orange-500/50 focus:ring-offset-1 focus:ring-offset-gray-900';
	const inputClasses = `${baseInputClasses} border-gray-700`;
	const labelClasses = 'mb-1.5 block text-sm font-medium text-gray-200';
	const buttonClasses =
		'flex items-center justify-center gap-2 rounded-lg px-4 py-2.5 font-medium shadow-lg transition duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-900 disabled:opacity-60';
	const primaryButtonClasses = `${buttonClasses} bg-orange-600 text-white hover:bg-orange-500 active:bg-orange-700 focus:ring-orange-500`;
	const secondaryButtonClasses = `${buttonClasses} bg-gray-700 text-white hover:bg-gray-600 active:bg-gray-800 focus:ring-gray-500`;
	const dangerButtonClasses = `${buttonClasses} bg-red-600 text-white hover:bg-red-700 active:bg-red-800 focus:ring-red-500`;
	const smallButtonClasses =
		'rounded-md px-2.5 py-1.5 text-xs font-medium transition focus:outline-none focus:ring-2';
</script>

<!-- Toast Notification -->
{#if actionFeedback}
	<div
		in:fly={{ y: 32, duration: 300, easing: backOut }}
		out:fade={{ duration: 200 }}
		class="fixed right-6 bottom-6 z-50 flex max-w-md min-w-72 items-center gap-3 rounded-lg border px-4 py-3 shadow-2xl md:right-8 md:bottom-8
			{actionFeedback.type === 'success'
			? 'border-green-500 bg-gradient-to-r from-green-900/90 to-green-800/90 text-green-50'
			: 'border-red-500 bg-gradient-to-r from-red-900/90 to-red-800/90 text-red-50'}"
		role="status"
		aria-live="polite"
	>
		<div class="flex-1">
			<p class="font-medium">
				{actionFeedback.type === 'success' ? '✓ Success' : '✕ Error'}
			</p>
			<p class="text-sm opacity-90">{actionFeedback.message}</p>
		</div>
		<button
			class="self-start rounded-full p-1 opacity-80 transition hover:bg-white/10 hover:opacity-100"
			on:click={() => (actionFeedback = undefined)}
			aria-label="Dismiss notification"
		>
			<span class="text-lg">×</span>
		</button>
	</div>
{/if}

<!-- Modal for delete confirmation -->
{#if workoutToDelete}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
		in:fade={{ duration: 200 }}
		out:fade={{ duration: 150 }}
	>
		<div
			class="relative w-full max-w-md rounded-xl border border-gray-700 bg-gray-900 p-6 shadow-2xl"
			in:scale={{ duration: 250, start: 0.9, easing: quintOut }}
		>
			<h3 class="mb-3 text-xl font-semibold text-white">Confirm Deletion</h3>
			<p class="mb-5 text-gray-300">
				Are you sure you want to delete this workout? This action cannot be undone.
			</p>

			<div class="flex justify-end gap-3">
				<button type="button" class={secondaryButtonClasses} on:click={cancelDelete}>
					Cancel
				</button>

				<form
					method="POST"
					action="?/deleteWorkout"
					use:enhance={() => {
						return ({ update }) => {
							workoutToDelete = null;
							update();
						};
					}}
				>
					<input type="hidden" name="workoutId" value={workoutToDelete} />
					<button type="submit" class={dangerButtonClasses}> Delete Workout </button>
				</form>
			</div>
		</div>
	</div>
{/if}

<!-- Displaying a stats dashboard with toggle arrow at the bottom -->
<div class="mx-auto max-w-4xl space-y-6 p-4 md:p-8">
	<section
		class="relative overflow-hidden rounded-xl border border-orange-900/30 bg-gradient-to-br from-gray-900 via-gray-900 to-gray-900/95 shadow-xl"
	>
		<div class="p-5 {showStats ? 'block' : 'hidden'}">
			<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
				<div
					class="flex flex-col items-center justify-center space-y-2 rounded-lg bg-black/20 p-4 text-center"
				>
					<span class="text-xl text-orange-400">🏆</span>
					<h3 class="text-sm font-medium text-orange-300/90">Total Workouts</h3>
					<p class="text-2xl font-bold text-white">{stats.totalWorkouts}</p>
				</div>
				<div
					class="flex flex-col items-center justify-center space-y-2 rounded-lg bg-black/20 p-4 text-center"
				>
					<span class="text-xl text-orange-400">⏱️</span>
					<h3 class="text-sm font-medium text-orange-300/90">Total Minutes</h3>
					<p class="text-2xl font-bold text-white">{stats.totalMinutes}</p>
					<p class="text-xs text-orange-300/80">Avg: {stats.avgDuration} min/workout</p>
				</div>
				<div
					class="flex flex-col items-center justify-center space-y-2 rounded-lg bg-black/20 p-4 text-center"
				>
					<span class="text-xl text-orange-400">🔥</span>
					<h3 class="text-sm font-medium text-orange-300/90">Calories Burned</h3>
					<p class="text-2xl font-bold text-white">{stats.totalCalories}</p>
				</div>
			</div>
		</div>

		<!-- Toggle button at the bottom of the stats panel -->
		<button
			class="flex w-full items-center justify-center border-t border-orange-900/30 bg-black/30 py-1.5 text-sm text-orange-400 transition hover:bg-black/40"
			on:click={() => (showStats = !showStats)}
			aria-label={showStats ? 'Hide stats' : 'Show stats'}
		>
			<span>{showStats ? 'Hide Stats' : 'Show Stats'}</span>
			<span class="ml-1 inline-block transition duration-300 {showStats ? 'rotate-180' : ''}"
				>▼</span
			>
		</button>
	</section>

	<!-- Tab Navigation -->
	<div class="flex space-x-1 rounded-lg bg-gray-800/60 p-1">
		<button
			class="flex-1 rounded-md px-4 py-2.5 text-center font-medium transition {activeTab === 'log'
				? 'bg-gray-700 text-white'
				: 'text-gray-400 hover:bg-gray-700/40 hover:text-gray-300'}"
			on:click={() => (activeTab = 'log')}
		>
			My Workouts
		</button>
		<button
			class="flex-1 rounded-md px-4 py-2.5 text-center font-medium transition {activeTab === 'add'
				? 'bg-gray-700 text-white'
				: 'text-gray-400 hover:bg-gray-700/40 hover:text-gray-300'}"
			on:click={() => (activeTab = 'add')}
		>
			Add New Workout
		</button>
	</div>

	<!-- Workout Log Tab -->
	{#if activeTab === 'log'}
		<section
			in:fade={{ duration: 250 }}
			class="min-h-[400px] rounded-xl border border-gray-800 bg-gray-900/90 p-6 shadow-xl"
		>
			<h2 class="mb-6 text-2xl font-semibold text-white">Workout History</h2>

			{#if pageLoadError}
				<div
					in:fade={{ duration: 200 }}
					class="mb-5 rounded-lg border border-red-800 bg-red-900/40 p-4 text-red-300"
				>
					<p class="text-sm">{pageLoadError}</p>
				</div>
			{/if}

			{#if workouts.length > 0}
				<div class="space-y-4">
					{#each workouts as workout (workout.id)}
						{@const isExpanded = expandedWorkoutId === workout.id}
						{@const intensity = getIntensity(
							workout.calories_burned || 0,
							workout.duration_minutes || 1
						)}

						<div
							in:scale={{ duration: 200, easing: quintOut, start: 0.95 }}
							animate:flip={{ duration: 200 }}
							class="overflow-hidden rounded-lg border border-gray-700/50 bg-gradient-to-r from-gray-800/80 to-gray-800/50 shadow-lg transition-all duration-300 hover:border-gray-600/70 {isExpanded
								? 'shadow-xl'
								: ''}"
						>
							<div
								class="relative flex flex-col p-5 sm:flex-row sm:items-center sm:justify-between"
							>
								<!-- Left content: Title and main info -->
								<div class="flex-1 pr-4">
									<div class="flex items-center justify-between gap-2.5">
										<div class="flex items-center gap-2">
											<span
												class="text-lg {intensity === 'low'
													? 'text-green-400'
													: intensity === 'medium'
														? 'text-yellow-400'
														: 'text-red-400'}"
											>
												💪
											</span>
											<h3 class="text-xl font-semibold text-orange-300">{workout.title}</h3>
										</div>

										<!-- Improving the details toggle button -->
										<button
											class="flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-gray-400 transition hover:bg-gray-800 hover:text-white"
											on:click={() => toggleWorkoutDetails(workout.id)}
											aria-expanded={isExpanded}
											aria-label={isExpanded
												? 'Collapse workout details'
												: 'Expand workout details'}
										>
											{isExpanded ? 'Hide Details' : 'Show Details'}
											<span class="text-xs">{isExpanded ? '▲' : '▼'}</span>
										</button>
									</div>

									<p class="mt-1.5 line-clamp-1 text-sm text-gray-300">{workout.description}</p>

									<div class="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs">
										<span class="text-gray-400">
											⏱️ {workout.duration_minutes} minutes
										</span>
										<span class="text-gray-400">
											🔥 {workout.calories_burned} calories
										</span>
										{#if workout.created_at}
											<span class="text-gray-500">
												{formatDate(workout.created_at)}
											</span>
										{/if}
									</div>
								</div>

								<!-- Right content: Actions -->
								<div class="mt-4 flex items-center gap-2 sm:mt-0">
									<a
										href={`/workouts/${workout.id}/edit`}
										class="{smallButtonClasses} bg-amber-600/90 text-white hover:bg-amber-500"
									>
										<span class="hidden sm:inline">Edit</span>
										<span class="sm:hidden">✏️</span>
									</a>

									<button
										type="button"
										class="{smallButtonClasses} bg-red-600/90 text-white hover:bg-red-500"
										on:click={() => confirmDelete(workout.id)}
									>
										<span class="hidden sm:inline">Delete</span>
										<span class="sm:hidden">🗑️</span>
									</button>
								</div>
							</div>

							<!-- Details section -->
							{#if isExpanded && workout.entries && workout.entries.length > 0}
								<div
									in:slide={{ duration: 200, easing: quintOut }}
									class="border-t border-gray-700/50 bg-gray-800/60 px-5 py-4"
								>
									<h4 class="mb-3 font-medium text-gray-300">Exercise Details</h4>
									<ul class="grid grid-cols-1 gap-3 md:grid-cols-2">
										{#each workout.entries as entry, i (entry.id || entry.order_index)}
											<li
												class="relative rounded-md border border-gray-700/50 bg-gray-700/40 p-3 pl-8 text-sm leading-relaxed text-gray-300"
											>
												<!-- Exercise number -->
												<div
													class="absolute top-3 left-2 flex h-5 w-5 items-center justify-center rounded-full bg-orange-600/80 text-xs font-bold text-white"
												>
													{i + 1}
												</div>

												<div class="font-medium text-white">
													{entry.exercise_name}
												</div>
												<div class="mt-1 flex flex-wrap gap-x-4 gap-y-1 text-xs">
													<span>{entry.sets} sets</span>
													{#if entry.reps}
														<span>{entry.reps} reps</span>
													{/if}
													{#if entry.duration_seconds}
														<span>{entry.duration_seconds}s</span>
													{/if}
													{#if entry.weight}
														<span>{entry.weight}kg</span>
													{/if}
												</div>
												{#if entry.notes}
													<div class="mt-1.5 text-xs text-gray-400 italic">"{entry.notes}"</div>
												{/if}
											</li>
										{/each}
									</ul>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{:else if !pageLoadError}
				<div
					class="flex min-h-[200px] flex-col items-center justify-center rounded-lg border border-dashed border-gray-700 bg-gray-800/30 p-8 text-center"
				>
					<span class="mb-2 text-4xl text-gray-600">💪</span>
					<p class="mb-4 text-gray-400">You haven't added any workouts yet</p>
					<button class="{primaryButtonClasses} text-sm" on:click={() => (activeTab = 'add')}>
						<span>➕ Add Your First Workout</span>
					</button>
				</div>
			{/if}
		</section>
	{/if}

	<!-- Add Workout form -->
	{#if activeTab === 'add'}
		<section
			in:fade={{ duration: 250 }}
			class="overflow-hidden rounded-xl border border-gray-800 bg-gray-900/90 shadow-xl"
		>
			<!-- Header with Nike-inspired design -->
			<div class="relative bg-gradient-to-r from-orange-700 to-orange-500 px-6 py-5">
				<h2 class="text-2xl font-bold text-white">Create New Workout</h2>
				<p class="mt-1 text-sm text-orange-100/80">Track your progress, achieve your goals</p>

				<!-- Decorative elements -->
				<div class="absolute top-4 -right-4 h-20 w-20 rounded-full bg-orange-400/20"></div>
				<div class="absolute -top-6 -right-6 h-16 w-16 rounded-full bg-orange-400/10"></div>
			</div>

			<form
				method="POST"
				action="?/addWorkout"
				use:enhance={() => {
					handleSubmit();
					return ({ update }) => {
						update({ reset: false });
					};
				}}
				class="p-6"
			>
				<!-- Workout details card -->
				<div class="mb-6 rounded-lg border border-gray-700/50 bg-black/20 p-5">
					<h3 class="mb-4 flex items-center text-lg font-semibold text-white">
						<span class="mr-2 text-orange-400">📋</span> Workout Details
					</h3>

					<div class="grid grid-cols-1 gap-5 md:grid-cols-2">
						<div>
							<label for="title" class={labelClasses}>Workout Title</label>
							<input
								bind:value={newWorkoutTitle}
								type="text"
								id="title"
								name="title"
								required
								class={inputClasses}
								placeholder="e.g., Morning Strength"
								autocomplete="off"
							/>
						</div>
						<div>
							<label for="description" class={labelClasses}>Description</label>
							<input
								bind:value={newWorkoutDescription}
								type="text"
								id="description"
								name="description"
								required
								class={inputClasses}
								placeholder="e.g., Full body workout"
								autocomplete="off"
							/>
						</div>
					</div>

					<div class="mt-5 grid grid-cols-1 gap-5 md:grid-cols-2">
						<div>
							<label for="durationMinutes" class={labelClasses}>Duration (minutes)</label>
							<input
								bind:value={newWorkoutDuration}
								type="number"
								id="durationMinutes"
								name="durationMinutes"
								min="0"
								required
								class={inputClasses}
								placeholder="60"
							/>
						</div>
						<div>
							<label for="caloriesBurned" class={labelClasses}>Calories Burned</label>
							<input
								bind:value={newWorkoutCalories}
								type="number"
								id="caloriesBurned"
								name="caloriesBurned"
								min="0"
								required
								class={inputClasses}
								placeholder="300"
							/>
						</div>
					</div>
				</div>

				<!-- Exercises Section -->
				<div class="space-y-5">
					<div class="flex items-center justify-between">
						<h3 class="flex items-center text-lg font-semibold text-white">
							<span class="mr-2 text-orange-400">💪</span> Exercises
						</h3>
						<button
							type="button"
							on:click={addEntryRow}
							class="rounded-full bg-orange-600 px-4 py-2 text-sm font-medium text-white shadow-lg transition hover:bg-orange-500 focus:ring-2 focus:ring-orange-500 focus:ring-offset-2 focus:ring-offset-gray-900 focus:outline-none"
						>
							Add Exercise
						</button>
					</div>

					<!-- Hidden form field for entries data -->
					<input
						type="hidden"
						name="entries"
						value={JSON.stringify(
							newWorkoutEntries.map((e, idx) => ({
								exercise_name: e.exercise_name,
								sets: e.sets,
								reps: e.reps && e.reps > 0 ? e.reps : null,
								duration_seconds:
									e.duration_seconds && e.duration_seconds > 0 ? e.duration_seconds : null,
								weight: e.weight != null && e.weight >= 0 ? e.weight : null,
								notes: e.notes || '',
								order_index: idx + 1
							}))
						)}
					/>

					<div class="max-h-[600px] space-y-5 overflow-y-auto pr-1">
						{#each newWorkoutEntries as entry, index (index)}
							<div
								in:scale={{ duration: 180, start: 0.95 }}
								out:fade={{ duration: 120 }}
								class="relative rounded-lg border border-gray-700/60 bg-gray-800/50 shadow-md"
							>
								<!-- Exercise number positioning -->
								<div
									class="absolute -top-2 left-4 flex h-8 min-w-8 items-center justify-center rounded-md bg-orange-600 px-2 text-sm font-bold text-white shadow-lg"
								>
									#{index + 1}
								</div>

								<div class="p-5 pt-7">
									<button
										type="button"
										on:click={() => removeEntryRow(index)}
										class="absolute top-3 right-3 flex h-7 w-7 items-center justify-center rounded-full bg-red-600/80 text-white transition hover:bg-red-500"
										aria-label="Remove exercise"
									>
										×
									</button>

									<div class="mb-4">
										<label for="entry_name_{index}" class={labelClasses}>Exercise Name</label>
										<input
											bind:value={entry.exercise_name}
											type="text"
											id="entry_name_{index}"
											required
											class={inputClasses}
											placeholder="Bench Press"
											autocomplete="off"
										/>
									</div>

									<div class="grid grid-cols-3 gap-4">
										<div>
											<label for="entry_sets_{index}" class={labelClasses}>Sets</label>
											<input
												bind:value={entry.sets}
												type="number"
												id="entry_sets_{index}"
												min="1"
												required
												class={inputClasses}
												placeholder="3"
											/>
										</div>
										<div>
											<label for="entry_reps_{index}" class={labelClasses}>Reps</label>
											<div class="relative">
												<input
													bind:value={entry.reps}
													type="number"
													id="entry_reps_{index}"
													min="0"
													class="{inputClasses} {entry.duration_seconds
														? 'border-gray-600 bg-gray-900 '
														: ''}"
													placeholder="10"
													on:input={() => handleRepsInput(entry)}
													disabled={!!entry.duration_seconds}
												/>
												{#if entry.duration_seconds}
													<div
														class="absolute inset-0 flex cursor-pointer items-center justify-center rounded-lg bg-gray-900 text-xs text-gray-400"
														on:click={() => {
															entry.duration_seconds = null;
														}}
													>
														Click to use reps instead
													</div>
												{/if}
											</div>
										</div>
										<div>
											<label for="entry_duration_{index}" class={labelClasses}>Duration (s)</label>
											<div class="relative">
												<input
													bind:value={entry.duration_seconds}
													type="number"
													id="entry_duration_{index}"
													min="0"
													class="{inputClasses} {entry.reps ? 'border-gray-600 bg-gray-900 ' : ''}"
													placeholder="60"
													on:input={() => handleDurationInput(entry)}
													disabled={!!entry.reps}
												/>
												{#if entry.reps}
													<div
														class="absolute inset-0 flex cursor-pointer items-center justify-center rounded-lg bg-gray-900 text-xs text-gray-400"
														on:click={() => {
															entry.reps = null;
														}}
													>
														Click to use duration instead
													</div>
												{/if}
											</div>
										</div>
									</div>

									<div class="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
										<div>
											<label for="entry_weight_{index}" class={labelClasses}>Weight (kg)</label>
											<input
												bind:value={entry.weight}
												type="number"
												step="any"
												id="entry_weight_{index}"
												min="0"
												class={inputClasses}
												placeholder="60,5"
											/>
										</div>
										<div>
											<label for="entry_notes_{index}" class={labelClasses}>Notes</label>
											<input
												bind:value={entry.notes}
												type="text"
												id="entry_notes_{index}"
												class={inputClasses}
												placeholder="e.g., Tempo 2-0-1-0"
											/>
										</div>
									</div>
								</div>

								<!-- Exercise card footer -->
								<div
									class="border-t border-gray-700/40 bg-black/20 px-5 py-2 text-right text-xs text-gray-500"
								>
									{entry.exercise_name ? entry.exercise_name : 'New exercise'}
								</div>
							</div>
						{/each}
					</div>
				</div>

				<!-- Form actions bar with Nike-inspired design -->
				<div
					class="mt-6 flex items-center justify-between rounded-lg bg-gradient-to-r from-gray-800 to-gray-900 p-4"
				>
					<button
						type="button"
						class="flex items-center gap-1 rounded-md border border-gray-600 bg-transparent px-3 py-2 text-sm font-medium text-gray-300 transition hover:bg-gray-800 hover:text-white"
						on:click={resetForm}
					>
						<span>Reset</span>
					</button>

					<button
						type="submit"
						class="group relative overflow-hidden rounded-md bg-gradient-to-r from-orange-600 to-orange-500 px-5 py-2.5 text-sm font-medium text-white shadow-lg transition hover:from-orange-500 hover:to-orange-400 focus:ring-2 focus:ring-orange-500 focus:ring-offset-2 focus:ring-offset-gray-900 focus:outline-none"
					>
						<span class="relative z-10">Save Workout</span>
						<!-- Save button animation -->
						<span
							class="absolute inset-0 -translate-x-full transform bg-gradient-to-r from-orange-400/30 via-orange-300/20 to-transparent transition-transform duration-700 ease-in-out group-hover:translate-x-full"
						></span>
					</button>
				</div>
			</form>
		</section>
	{/if}
</div>
