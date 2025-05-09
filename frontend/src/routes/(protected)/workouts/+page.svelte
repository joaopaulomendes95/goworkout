<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData, PageData } from './$types';
	import type { BackendWorkout, BackendWorkoutEntry } from '$lib/types';
	import { invalidateAll } from '$app/navigation';
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
	let newWorkoutTitle = $state('');
	let newWorkoutDescription = $state('');
	let newWorkoutDuration = $state(0);
	let newWorkoutCalories = $state(0);
	let newWorkoutEntries = $state<Partial<BackendWorkoutEntry>[]>([
		{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }
	]);

    // --- State to track if the current `form` prop has been processed by the effect ---
    let formProcessed = $state(false);
    // ---

    // Effect to repopulate "Add Workout" form on failure
    $effect(() => {
        if (form && form.formError && form.title !== undefined) {
            console.log('[Workouts Page Svelte] Repopulating Add form from failed submission:', form);
            newWorkoutTitle = form.title || '';
            newWorkoutDescription = form.description || '';
            newWorkoutDuration = form.durationMinutes ? parseInt(form.durationMinutes) : 0;
            newWorkoutCalories = form.caloriesBurned ? parseInt(form.caloriesBurned) : 0;
            if (form.entries && typeof form.entries === 'string') {
                try { 
                    const parsedEntries = JSON.parse(form.entries);
                    newWorkoutEntries = Array.isArray(parsedEntries) && parsedEntries.length > 0 
                        ? parsedEntries 
                        : [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }];
                } catch (e) { 
                    console.error("Error parsing entries from form prop:", e);
                    newWorkoutEntries = [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }];
                }
            }
        }
    });

	// Effect for handling messages from redirects (e.g., after successful edit)
	$effect(() => {
        if (browser) { 
            const messageFromRedirect = $page.url.searchParams.get('message');
            if (messageFromRedirect) {
                console.log('[Workouts Page Svelte] Effect: Detected messageFromRedirect:', messageFromRedirect);
                actionFeedback = { type: 'success', message: messageFromRedirect.replace(/_/g, ' ') };
                const newUrl = new URL(window.location.href); 
                newUrl.searchParams.delete('message');
                window.history.replaceState(window.history.state, '', newUrl); 
            }
        }
	});
	
    // This effect runs when `form` prop changes.
    // We use `formProcessed` to ensure we only act on a new `form` value once.
	$effect(() => {
		const currentForm = form; 
        console.log('[Workouts Page Svelte] Form Processing Effect. Current Form:', currentForm, 'Form Processed Flag:', formProcessed);

		if (currentForm && !formProcessed) { // Only process if form exists and not yet processed
            console.log('[Workouts Page Svelte] Processing new form data:', currentForm);
            formProcessed = true; // Mark as processed immediately

			if (currentForm.success && currentForm.message) {
				actionFeedback = { type: 'success', message: currentForm.message };
				if (currentForm.message?.includes('added successfully')) { 
                    newWorkoutTitle = '';
                    newWorkoutDescription = '';
                    newWorkoutDuration = 0;
                    newWorkoutCalories = 0;
                    newWorkoutEntries = [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }];
                }
				
				if (browser) { 
					console.log('[Workouts Page Svelte] Form success: Calling invalidateAll().');
					invalidateAll(); 
				}

			} else if (currentForm.formError) {
				actionFeedback = { type: 'error', message: currentForm.formError };
				// Repopulation is handled by the other effect now
			} else if (currentForm.message && !currentForm.success) {
	            actionFeedback = { type: 'error', message: currentForm.message };
	        }

            // Timer to clear the feedback message
            if (actionFeedback) {
                const timer = setTimeout(() => { 
                    console.log('[Workouts Page Svelte] Clearing actionFeedback.');
                    actionFeedback = undefined; 
                }, 4000);
                return () => clearTimeout(timer); // Cleanup for this specific effect run
            }
		} else if (!currentForm && formProcessed) {
            // If form becomes undefined (e.g., after invalidateAll and reload), reset the flag
            console.log('[Workouts Page Svelte] Form is now undefined, resetting formProcessed flag.');
            formProcessed = false;
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

<!-- TEMPLATE REMAINS THE SAME -->
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
