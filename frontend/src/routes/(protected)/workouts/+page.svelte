<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData, PageData } from './$types';
	import type { BackendWorkout, BackendWorkoutEntry } from '$lib/types';
	import { invalidateAll } from '$app/navigation'; // Removed navigating
	import { page } from '$app/stores';
    import { browser } from '$app/environment';

	const { data, form } = $props<{ data: PageData; form?: ActionData }>();

	const workouts = $derived(data.workouts || []);
	const pageLoadError = $derived(
		data.error 
			? (typeof data.error === 'string' ? data.error : data.error.message) 
			: undefined
	);
	
	let actionFeedback = $state<{ type: 'success' | 'error'; message: string } | undefined>(undefined);

	// Form state for adding a new workout
	let newWorkoutTitle = $state(form?.title || '');
	let newWorkoutDescription = $state(form?.description || '');
	let newWorkoutDuration = $state(form?.durationMinutes ? parseInt(form.durationMinutes) : 0);
	let newWorkoutCalories = $state(form?.caloriesBurned ? parseInt(form.caloriesBurned) : 0);
	let newWorkoutEntries = $state<Partial<BackendWorkoutEntry>[]>(
		form?.entries && typeof form.entries === 'string'
			? JSON.parse(form.entries)
			: [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }]
	);

    // --- State to track processed form ---
    let lastProcessedFormTimestamp = $state<number | undefined>(undefined);
    // --- End State to track processed form ---

	// Effect for handling messages from redirects (e.g., after edit)
	$effect(() => {
        if (browser) { 
            const successMessage = $page.url.searchParams.get('message');
            if (successMessage) {
                actionFeedback = { type: 'success', message: successMessage.replace(/_/g, ' ') };
                const newUrl = new URL(window.location.href); 
                newUrl.searchParams.delete('message');
                window.history.replaceState(window.history.state, '', newUrl); 
            }
        }
	});
	
	// Effect for handling form action results
	$effect(() => {
		const currentForm = form; // Capture current form prop for this effect run

        // Only process if 'form' is defined and it's a new submission
        // (form.timestamp is a hypothetical property SvelteKit might add for this,
        // but since it doesn't, we use our own lastProcessedFormTimestamp)
        // A simpler check: if form is present and we haven't processed its specific content yet.
        // For this, we can check if the form object reference has changed or if its content is new.
        // However, SvelteKit re-creates the form object. So, we check if it has data.
		if (!currentForm || (currentForm.timestamp && currentForm.timestamp === lastProcessedFormTimestamp)) {
            // If form is undefined or we've already processed this specific form instance (hypothetical timestamp)
            // A more practical check: if form is undefined, or if it's the same actionFeedback message we just set.
            // This is still tricky. The core idea is to not re-process the *exact same* form result.
            return;
        }
        
        // If we are here, it's potentially a new form result to process.
        // Let's assume `form` itself being non-null means it's a new action result.
        // The problem is if invalidateAll() doesn't lead to `form` becoming undefined before this effect re-runs.

		if (currentForm.success && currentForm.message) {
			actionFeedback = { type: 'success', message: currentForm.message };
			if (currentForm.message?.includes('added successfully')) { 
                newWorkoutTitle = '';
                newWorkoutDescription = '';
                newWorkoutDuration = 0;
                newWorkoutCalories = 0;
                newWorkoutEntries = [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }];
            }
			
            // Mark this form instance as processed (using a hypothetical timestamp or a unique ID if available)
            // Since SvelteKit doesn't provide a unique ID on the form prop for this,
            // this approach is difficult.
            // The primary reliance is that `invalidateAll` should lead to `form` being undefined
            // in the *next* distinct render cycle after data reloads.

			if (browser) { 
				invalidateAll(); 
			}

		} else if (currentForm.formError) {
			actionFeedback = { type: 'error', message: currentForm.formError };
			newWorkoutTitle = currentForm.title || '';
			newWorkoutDescription = currentForm.description || '';
			newWorkoutDuration = currentForm.durationMinutes ? parseInt(currentForm.durationMinutes) : 0;
			newWorkoutCalories = currentForm.caloriesBurned ? parseInt(currentForm.caloriesBurned) : 0;
			if (currentForm.entries && typeof currentForm.entries === 'string') {
				try { newWorkoutEntries = JSON.parse(currentForm.entries); } catch (e) { /* ignore */ }
			}
		} else if (currentForm.message && !currentForm.success) {
            actionFeedback = { type: 'error', message: currentForm.message };
        }


		if (actionFeedback) {
			const timer = setTimeout(() => { 
                actionFeedback = undefined; 
            }, 4000);
			return () => clearTimeout(timer);
		}
	});


	function addEntryRow() {
		newWorkoutEntries.push(
			{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: newWorkoutEntries.length + 1, weight: null }
		);
	}

	function removeEntryRow(index: number) {
		newWorkoutEntries.splice(index, 1);
		newWorkoutEntries.forEach((e, i) => e.order_index = i + 1);
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

<!-- ... rest of the template remains the same ... -->
<div class="space-y-8 p-4 md:p-6">
	<h1 class="text-3xl font-bold text-white">My Workouts</h1>

	{#if actionFeedback}
		<div class={`mb-3 text-sm p-3 rounded-md ${actionFeedback.type === 'success' ? 'bg-green-600 text-green-100 border border-green-500' : 'bg-red-600 text-red-100 border border-red-500'}`}>
			{actionFeedback.message}
		</div>
	{/if}

	<!-- Add Workout Form -->
	<section class="rounded-lg bg-gray-800 p-6 shadow-md">
		<h2 class="mb-4 text-2xl font-semibold text-white">Add New Workout</h2>
		<form method="POST" action="?/addWorkout" use:enhance class="space-y-4">
			<div>
				<label for="title" class={labelClasses}>Title</label>
				<input bind:value={newWorkoutTitle} type="text" id="title" name="title" required class={inputClasses} placeholder="e.g., Morning Strength" />
			</div>
			<div>
				<label for="description" class={labelClasses}>Description</label>
				<input bind:value={newWorkoutDescription} type="text" id="description" name="description" required class={inputClasses} placeholder="e.g., Full body workout" />
			</div>
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<div>
					<label for="durationMinutes" class={labelClasses}>Duration (minutes)</label>
					<input bind:value={newWorkoutDuration} type="number" id="durationMinutes" name="durationMinutes" min="0" required class={inputClasses} placeholder="60" />
				</div>
				<div>
					<label for="caloriesBurned" class={labelClasses}>Calories Burned</label>
					<input bind:value={newWorkoutCalories} type="number" id="caloriesBurned" name="caloriesBurned" min="0" required class={inputClasses} placeholder="300" />
				</div>
			</div>

			<h3 class="border-t border-gray-700 pt-4 text-xl font-semibold text-white">Entries</h3>
			<input type="hidden" name="entries" value={JSON.stringify(newWorkoutEntries.map((e, idx) => ({
				exercise_name: e.exercise_name,
				sets: e.sets,
				reps: e.reps && e.reps > 0 ? e.reps : null,
				duration_seconds: e.duration_seconds && e.duration_seconds > 0 ? e.duration_seconds : null,
				weight: e.weight != null && e.weight >= 0 ? e.weight : null,
				notes: e.notes || '',
				order_index: idx + 1
			})))} />

			<div class="space-y-3">
				{#each newWorkoutEntries as entry, index (index)}
					<div class="relative space-y-2 rounded-md border border-gray-700 p-3">
						<button type="button" onclick={() => removeEntryRow(index)} class="absolute right-1 top-1 {smallButtonClasses} bg-red-700 text-white hover:bg-red-600">×</button>
						<div>
							<label for="entry_name_{index}" class={labelClasses}>Exercise Name</label>
							<input bind:value={entry.exercise_name} type="text" id="entry_name_{index}" required class={inputClasses} placeholder="Bench Press" />
						</div>
						<div class="grid grid-cols-2 gap-2 md:grid-cols-3">
							<div>
								<label for="entry_sets_{index}" class={labelClasses}>Sets</label>
								<input bind:value={entry.sets} type="number" id="entry_sets_{index}" min="1" required class={inputClasses} placeholder="3" />
							</div>
							<div>
								<label for="entry_reps_{index}" class={labelClasses}>Reps</label>
								<input bind:value={entry.reps} type="number" id="entry_reps_{index}" min="0" class={inputClasses} placeholder="10" oninput={() => handleRepsInput(entry)} />
							</div>
							<div>
								<label for="entry_duration_{index}" class={labelClasses}>Duration (s)</label>
								<input bind:value={entry.duration_seconds} type="number" id="entry_duration_{index}" min="0" class={inputClasses} placeholder="60" oninput={() => handleDurationInput(entry)} />
							</div>
						</div>
						<div>
							<label for="entry_weight_{index}" class={labelClasses}>Weight (kg)</label>
							<input bind:value={entry.weight} type="number" step="any" id="entry_weight_{index}" min="0" class={inputClasses} placeholder="60.5" />
						</div>
						<div>
							<label for="entry_notes_{index}" class={labelClasses}>Notes</label>
							<input bind:value={entry.notes} type="text" id="entry_notes_{index}" class={inputClasses} placeholder="e.g., Tempo 2-0-1-0" />
						</div>
					</div>
				{/each}
			</div>
			<button type="button" onclick={addEntryRow} class="{smallButtonClasses} border border-blue-400 text-blue-400 hover:bg-blue-900">
				+ Add Entry
			</button>

			<div class="pt-4">
				<button type="submit" class={buttonClasses}>Save Workout</button>
			</div>
		</form>
	</section>

	<section class="rounded-lg bg-gray-800 p-6 shadow-md">
		<h2 class="mb-4 text-2xl font-semibold text-white">Workout Log</h2>
		{#if pageLoadError}
			<p class="mb-3 rounded bg-red-900 p-2 text-sm text-red-400">{pageLoadError}</p>
		{/if}

		{#if workouts.length > 0}
			<div class="space-y-4">
				{#each workouts as workout (workout.id)}
					<div class="rounded-md bg-gray-700 p-4 shadow-sm">
						<div class="flex items-center justify-between">
							<h3 class="text-xl font-semibold text-blue-400">{workout.title}</h3>
							<div class="flex space-x-2">
								<a href={`/workouts/${workout.id}/edit`} class="{smallButtonClasses} bg-yellow-500 text-white hover:bg-yellow-600">Edit</a>
								<form method="POST" action="?/deleteWorkout" use:enhance>
									<input type="hidden" name="workoutId" value={workout.id} />
									<button type="submit" class="{smallButtonClasses} bg-red-600 text-white hover:bg-red-700">Delete</button>
								</form>
							</div>
						</div>
						<p class="mt-1 text-sm text-gray-300">{workout.description}</p>
						<p class="mt-1 text-xs text-gray-400">
							Duration: {workout.duration_minutes} mins | Calories: {workout.calories_burned}
						</p>
						{#if workout.entries && workout.entries.length > 0}
							<div class="mt-3 border-l-2 border-gray-600 pl-4">
								<h4 class="mb-1 text-sm font-medium text-gray-200">Entries:</h4>
								<ul class="list-disc list-inside space-y-1 text-xs">
									{#each workout.entries as entry (entry.id || entry.order_index)}
										<li>
											{entry.exercise_name}: {entry.sets} sets
											{#if entry.reps} of {entry.reps} reps{/if}
											{#if entry.duration_seconds} for {entry.duration_seconds}s{/if}
											{#if entry.weight} at {entry.weight}kg{/if}
											{#if entry.notes}
												<span class="italic text-gray-500">- "{entry.notes}"</span>
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
			<p class="text-gray-400">You haven't added any workouts yet. Use the form above to add one!</p>
		{/if}
	</section>
</div>
