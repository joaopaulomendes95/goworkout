<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData, PageData } from './$types';
	import type { BackendWorkoutEntry } from '$lib/types';

	const { data, form } = $props<{ data: PageData; form?: ActionData }>();

	// Initialize $state with potentially undefined values first
	let workoutTitle = $state<string>('');
	let workoutDescription = $state<string>('');
	let workoutDuration = $state<number>(0);
	let workoutCalories = $state<number>(0);
	let workoutEntries = $state<Partial<BackendWorkoutEntry>[]>([]);
	
	let actionFeedback = $state<{ type: 'error'; message: string } | undefined>(undefined);

	// This effect will run when `data.workout` or `form` changes.
	// It correctly prioritizes `form` data (e.g., after a failed submission with repopulated values)
	// over the initial `data.workout` from the load function.
	$effect(() => {
		const source = form?.title !== undefined ? form : data.workout;

		if (source) {
			workoutTitle = source.title || '';
			workoutDescription = source.description || '';
			
			// Ensure numeric conversion if values might come from form as strings
			const durMinutes = source.durationMinutes ?? source.duration_minutes;
			workoutDuration = typeof durMinutes === 'string' ? parseInt(durMinutes, 10) || 0 : (durMinutes || 0);

			const calBurned = source.caloriesBurned ?? source.calories_burned;
			workoutCalories = typeof calBurned === 'string' ? parseInt(calBurned, 10) || 0 : (calBurned || 0);

			let entriesSource = source.entries;
			if (typeof entriesSource === 'string') {
				try { entriesSource = JSON.parse(entriesSource); } catch { entriesSource = []; }
			}
			
			workoutEntries = Array.isArray(entriesSource) ? entriesSource.map((e: any) => ({ ...e })) : [];

			if (workoutEntries.length === 0) {
				workoutEntries = [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }];
			}
		} else if (data.workout === null && !form?.title) {
            actionFeedback = { type: 'error', message: "Workout data could not be loaded or does not exist." };
        } else if (data.error && !form?.title) {
			actionFeedback = { type: 'error', message: typeof data.error === 'string' ? data.error : data.error.message };
		}
	});

	// Effect for handling form action error messages (from `fail`)
	$effect(() => {
        let currentForm = form; // Capture current form prop for this effect run
		if (currentForm?.formError) { // This is the primary error message from `fail`
			actionFeedback = { type: 'error', message: currentForm.formError };
		} else if (currentForm?.message && !currentForm.success) { // General message if `fail` was used with just `message`
            actionFeedback = { type: 'error', message: currentForm.message };
        }
		// Success messages are handled by redirecting from the action in this app.

		if (actionFeedback) {
			const timer = setTimeout(() => { actionFeedback = undefined; }, 4000);
			return () => clearTimeout(timer); // Cleanup function for the effect
		}
	});

	function addEntryRow() {
		workoutEntries.push(
			{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: workoutEntries.length + 1, weight: null }
		);
	}

	function removeEntryRow(index: number) {
		workoutEntries.splice(index, 1);
		// Re-index after removal
		workoutEntries.forEach((e, i) => e.order_index = i + 1);
	}

	function handleRepsInput(entry: Partial<BackendWorkoutEntry>) {
		if (entry.reps && entry.reps > 0) {
			entry.duration_seconds = null;
		}
	}

	function handleDurationInput(entry: Partial<BackendWorkoutEntry>) {
		if (entry.duration_seconds && entry.duration_seconds > 0) {
			entry.reps = null;
		}
	}

	const inputClasses = "w-full p-2 rounded bg-gray-700 border border-gray-600 focus:ring-blue-500 focus:border-blue-500 text-white placeholder-gray-400";
	const labelClasses = "block text-sm font-medium mb-1 text-gray-300";
	const buttonClasses = "px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md font-medium text-white disabled:opacity-50";
	const smallButtonClasses = "text-xs px-2 py-1 rounded";
</script>

<div class="space-y-8 p-4 md:p-6">
	<h1 class="text-3xl font-bold text-white">Edit Workout</h1>

	{#if (!data.workout && !form?.title) && !actionFeedback}
		<p class="text-yellow-400">Loading workout data...</p>
		{#if data.error}
		    <p class="text-red-400">Error: {typeof data.error === 'string' ? data.error : data.error.message}</p>
		{/if}
	{/if}

	{#if actionFeedback}
		<div class={`mb-3 text-sm p-3 rounded-md ${actionFeedback.type === 'error' ? 'bg-red-600 text-red-100 border border-red-500' : 'bg-green-600 text-green-100 border border-green-500'}`}>
			{actionFeedback.message}
		</div>
	{/if}

	{#if data.workout || form?.title} <!-- Render form if we have initial data or form resubmission data -->
		<form method="POST" action="?/updateWorkout" use:enhance class="space-y-4 rounded-lg bg-gray-800 p-6 shadow-md">
			<div>
				<label for="title" class={labelClasses}>Title</label>
				<input bind:value={workoutTitle} type="text" id="title" name="title" required class={inputClasses} />
			</div>
			<div>
				<label for="description" class={labelClasses}>Description</label>
				<input bind:value={workoutDescription} type="text" id="description" name="description" required class={inputClasses} />
			</div>
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<div>
					<label for="durationMinutes" class={labelClasses}>Duration (minutes)</label>
					<input bind:value={workoutDuration} type="number" id="durationMinutes" name="durationMinutes" min="0" required class={inputClasses} />
				</div>
				<div>
					<label for="caloriesBurned" class={labelClasses}>Calories Burned</label>
					<input bind:value={workoutCalories} type="number" id="caloriesBurned" name="caloriesBurned" min="0" required class={inputClasses} />
				</div>
			</div>

			<h3 class="border-t border-gray-700 pt-4 text-xl font-semibold text-white">Entries</h3>
			<input type="hidden" name="entries" value={JSON.stringify(workoutEntries.map((e, idx) => ({
				id: e.id, 
				exercise_name: e.exercise_name,
				sets: e.sets,
				reps: e.reps && e.reps > 0 ? e.reps : null,
				duration_seconds: e.duration_seconds && e.duration_seconds > 0 ? e.duration_seconds : null,
				weight: e.weight != null && e.weight >= 0 ? e.weight : null,
				notes: e.notes || '',
				order_index: idx + 1 
			})))} />
			
			<div class="space-y-3">
				{#each workoutEntries as entry, index (entry.id || index)} 
					<div class="relative space-y-2 rounded-md border border-gray-700 p-3">
						<button type="button" onclick={() => removeEntryRow(index)} class="absolute right-1 top-1 {smallButtonClasses} bg-red-700 text-white hover:bg-red-600">×</button>
						<div>
							<label for="entry_name_{index}" class={labelClasses}>Exercise Name</label>
							<input bind:value={entry.exercise_name} type="text" id="entry_name_{index}" required class={inputClasses} />
						</div>
						<div class="grid grid-cols-2 gap-2 md:grid-cols-3">
							<div>
								<label for="entry_sets_{index}" class={labelClasses}>Sets</label>
								<input bind:value={entry.sets} type="number" id="entry_sets_{index}" min="1" required class={inputClasses} />
							</div>
							<div>
								<label for="entry_reps_{index}" class={labelClasses}>Reps</label>
								<input bind:value={entry.reps} type="number" id="entry_reps_{index}" min="0" class={inputClasses} oninput={() => handleRepsInput(entry)} />
							</div>
							<div>
								<label for="entry_duration_{index}" class={labelClasses}>Duration (s)</label>
								<input bind:value={entry.duration_seconds} type="number" id="entry_duration_{index}" min="0" class={inputClasses} oninput={() => handleDurationInput(entry)} />
							</div>
						</div>
						<div>
							<label for="entry_weight_{index}" class={labelClasses}>Weight (kg)</label>
							<input bind:value={entry.weight} type="number" step="any" id="entry_weight_{index}" min="0" class={inputClasses} />
						</div>
						<div>
							<label for="entry_notes_{index}" class={labelClasses}>Notes</label>
							<input bind:value={entry.notes} type="text" id="entry_notes_{index}" class={inputClasses} />
						</div>
					</div>
				{/each}
			</div>
			<button type="button" onclick={addEntryRow} class="{smallButtonClasses} border border-blue-400 text-blue-400 hover:bg-blue-900">
				+ Add Entry
			</button>

			<div class="flex space-x-4 pt-4">
				<button type="submit" class={buttonClasses}>Update Workout</button>
				<a href="/workouts" class="rounded-md bg-gray-600 px-4 py-2 font-medium text-white hover:bg-gray-500">Cancel</a>
			</div>
		</form>
	{/if}
</div>
