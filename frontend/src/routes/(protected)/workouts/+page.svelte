<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData, PageData } from './$types';
	import type { BackendWorkout, BackendWorkoutEntry } from '$lib/types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	// Workouts loaded from the server
	let serverLoadedWorkouts = $derived(data.workouts || []);
	let errorMessage = $derived(data.error);

	// For client-side list of workouts added in this session (optional, if you want immediate UI update without reload)
	let sessionWorkouts = $state<BackendWorkout[]>([]);

	const defaultEntry = {
		exercise_name: '',
		sets: 1,
		reps: null,
		duration_seconds: null,
		notes: '',
		order_index: 1,
		weight: null
	};

	let newWorkoutTitle = $state('');
	let newWorkoutDescription = $state('');
	let newWorkoutDuration = $state(0);
	let newWorkoutCalories = $state(0);
	let newWorkoutEntries = $state<Partial<BackendWorkoutEntry>[]>([defaultEntry]);

	function addEntryRow() {
		newWorkoutEntries.push({
			exercise_name: '',
			sets: 1,
			reps: null,
			duration_seconds: null,
			notes: '',
			order_index: newWorkoutEntries.length + 1,
			weight: null
		});
		newWorkoutEntries = [...newWorkoutEntries];
	}

	function removeEntryRow(index: number) {
		newWorkoutEntries.splice(index, 1);
		newWorkoutEntries.forEach((e, i) => (e.order_index = i + 1));
		newWorkoutEntries = [...newWorkoutEntries];
	}

	$effect(() => {
		if (form?.success && form?.message?.includes('added successfully')) {
			if (form.createdWorkout) {
				// Add to session workouts for immediate display
				sessionWorkouts.push(form.createdWorkout as BackendWorkout);
				sessionWorkouts = [...sessionWorkouts];
				// Optionally, you could also add it to serverLoadedWorkouts if you want to avoid a full reload
				// but typically after a successful action, you might invalidate data to refetch from server
				// For now, sessionWorkouts handles immediate UI feedback for new items.
			}
			// Reset form
			newWorkoutTitle = '';
			newWorkoutDescription = '';
			newWorkoutDuration = 0;
			newWorkoutCalories = 0;
			newWorkoutEntries = [{ ...defaultEntry }];
			// To refresh the list from the server after adding:
			// import { invalidate } from '$app/navigation';
			// invalidate('app:workouts'); // Or a more specific identifier if you have one
			// Or, if your action returns the new list or the new item, you can update client-side state.
		} else if (form && !form.success && form.title !== undefined) {
			newWorkoutTitle = form.title;
			newWorkoutDescription = form.description;
			newWorkoutDuration = parseInt(form.durationMinutes || '0');
			newWorkoutCalories = parseInt(form.caloriesBurned || '0');
			if (form.entries) newWorkoutEntries = JSON.parse(form.entries);
		}
	});

	let allWorkouts = $derived(
		[...serverLoadedWorkouts, ...sessionWorkouts].sort((a, b) => b.id - a.id)
	);
</script>

<div class="space-y-8">
	<h1 class="text-3xl font-bold">Manage Workouts</h1>

	<!-- Display Existing Workouts -->
	<section class="rounded-lg bg-gray-800 p-6 shadow-md">
		<h2 class="mb-4 text-2xl font-semibold">Your Workouts</h2>
		{#if errorMessage}
			<p class="mb-3 rounded bg-red-900 p-2 text-sm text-red-400">{errorMessage}</p>
		{/if}

		{#if allWorkouts.length > 0}
			<div class="space-y-4">
				{#each allWorkouts as workout (workout.id)}
					<div class="rounded-md bg-gray-700 p-4 shadow-sm">
						<h3 class="text-xl font-semibold text-blue-400">{workout.title}</h3>
						<p class="mt-1 text-sm text-gray-300">{workout.description}</p>
						{#if workout.entries && workout.entries.length > 0}
							<div class="mt-3 border-l-2 border-gray-600 pl-4">
								<h4 class="mb-1 text-sm font-medium text-gray-200">Entries:</h4>
								<ul class="list-inside list-disc space-y-1 text-xs">
									{#each workout.entries as entry (entry.id)}
										<li>
											{entry.exercise_name}: {entry.sets} sets
											{#if entry.reps}
												of {entry.reps} reps{/if}
											{#if entry.duration_seconds}
												for {entry.duration_seconds}s{/if}
											{#if entry.weight}
												at {entry.weight}kg{/if}
											{#if entry.notes}
												<span class="text-gray-500 italic">- "{entry.notes}"</span>{/if}
										</li>
									{/each}
								</ul>
							</div>
						{/if}
						<!-- TODO: Add Edit/Delete buttons here, perhaps linking to a workout detail page or using actions -->
					</div>
				{/each}
			</div>
		{:else if !errorMessage}
			<p class="text-gray-400">
				You haven't added any workouts yet. Use the form below to add one!
			</p>
		{/if}
	</section>

	<!-- Add Workout Form (remains largely the same) -->
	<section class="rounded-lg bg-gray-800 p-6 shadow-md">
		<h2 class="mb-4 text-2xl font-semibold">Add New Workout</h2>
		{#if form?.formError}
			<p class="mb-3 rounded bg-red-900 p-2 text-sm text-red-400">{form.formError}</p>
		{/if}
		{#if form?.success && form?.message && !form.formError}
			<p class="mb-3 rounded bg-green-900 p-2 text-sm text-green-400">{form.message}</p>
		{/if}

		<form method="POST" action="?/addWorkout" use:enhance class="space-y-4">
			<!-- Form fields for title, description, duration, calories -->
			<div>
				<label for="title" class="mb-1 block text-sm font-medium text-gray-300">Title</label>
				<input
					bind:value={newWorkoutTitle}
					type="text"
					id="title"
					name="title"
					required
					class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
					placeholder="e.g., Morning Strength"
				/>
			</div>
			<div>
				<label for="description" class="mb-1 block text-sm font-medium text-gray-300"
					>Description</label
				>
				<input
					bind:value={newWorkoutDescription}
					type="text"
					id="description"
					name="description"
					required
					class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
					placeholder="e.g., Full body workout"
				/>
			</div>
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<div>
					<label for="durationMinutes" class="mb-1 block text-sm font-medium text-gray-300"
						>Duration (minutes)</label
					>
					<input
						bind:value={newWorkoutDuration}
						type="number"
						id="durationMinutes"
						name="durationMinutes"
						min="0"
						required
						class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
						placeholder="60"
					/>
				</div>
				<div>
					<label for="caloriesBurned" class="mb-1 block text-sm font-medium text-gray-300"
						>Calories Burned</label
					>
					<input
						bind:value={newWorkoutCalories}
						type="number"
						id="caloriesBurned"
						name="caloriesBurned"
						min="0"
						required
						class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
						placeholder="300"
					/>
				</div>
			</div>

			<h3 class="border-t border-gray-700 pt-4 text-xl font-semibold">Entries</h3>
			<input
				type="hidden"
				name="entries"
				value={JSON.stringify(
					newWorkoutEntries.map((e, idx) => ({
						...e,
						order_index: idx + 1,
						reps: e.reps && e.reps > 0 ? e.reps : null,
						duration_seconds:
							e.duration_seconds && e.duration_seconds > 0 ? e.duration_seconds : null,
						weight: e.weight && e.weight > 0 ? e.weight : null
					}))
				)}
			/>

			<div class="space-y-3">
				{#each newWorkoutEntries as entry, index (index)}
					<div class="relative space-y-2 rounded-md border border-gray-700 p-3">
						<button
							type="button"
							onclick={() => removeEntryRow(index)}
							class="absolute top-1 right-1 rounded bg-red-700 px-2 py-1 text-xs text-white hover:bg-red-600"
							>×</button
						>
						<div>
							<label
								for={`entry_name_${index}`}
								class="mb-1 block text-sm font-medium text-gray-300">Exercise Name</label
							>
							<input
								bind:value={entry.exercise_name}
								type="text"
								id={`entry_name_${index}`}
								required
								class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
								placeholder="Bench Press"
							/>
						</div>
						<div class="grid grid-cols-2 gap-2 md:grid-cols-3">
							<div>
								<label
									for={`entry_sets_${index}`}
									class="mb-1 block text-sm font-medium text-gray-300">Sets</label
								>
								<input
									bind:value={entry.sets}
									type="number"
									id={`entry_sets_${index}`}
									min="1"
									required
									class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
									placeholder="3"
								/>
							</div>
							<div>
								<label
									for={`entry_reps_${index}`}
									class="mb-1 block text-sm font-medium text-gray-300">Reps (optional)</label
								>
								<input
									bind:value={entry.reps}
									type="number"
									id={`entry_reps_${index}`}
									min="0"
									class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
									placeholder="10"
									oninput={() => {
										if (entry.reps && entry.reps > 0) entry.duration_seconds = null;
									}}
								/>
							</div>
							<div>
								<label
									for={`entry_duration_${index}`}
									class="mb-1 block text-sm font-medium text-gray-300">Duration (s, optional)</label
								>
								<input
									bind:value={entry.duration_seconds}
									type="number"
									id={`entry_duration_${index}`}
									min="0"
									class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
									placeholder="e.g. 60"
									oninput={() => {
										if (entry.duration_seconds && entry.duration_seconds > 0) entry.reps = null;
									}}
								/>
							</div>
						</div>
						<div>
							<label
								for={`entry_weight_${index}`}
								class="mb-1 block text-sm font-medium text-gray-300">Weight (kg, optional)</label
							>
							<input
								bind:value={entry.weight}
								type="number"
								step="any"
								id={`entry_weight_${index}`}
								min="0"
								class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
								placeholder="60.5"
							/>
						</div>
						<div>
							<label
								for={`entry_notes_${index}`}
								class="mb-1 block text-sm font-medium text-gray-300">Notes</label
							>
							<input
								bind:value={entry.notes}
								type="text"
								id={`entry_notes_${index}`}
								class="w-full rounded border border-gray-600 bg-gray-700 p-2 text-white placeholder-gray-400 focus:border-blue-500 focus:ring-blue-500"
								placeholder="e.g., Tempo 2-0-1-0"
							/>
						</div>
					</div>
				{/each}
			</div>
			<button
				type="button"
				onclick={addEntryRow}
				class="rounded border border-blue-400 px-2 py-1 text-sm text-blue-400 hover:bg-blue-900 hover:text-blue-300"
			>
				+ Add Entry Row
			</button>

			<div class="pt-4">
				<button
					type="submit"
					class="rounded-md bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700 disabled:opacity-50"
					>Save Workout</button
				>
			</div>
		</form>
	</section>
</div>
