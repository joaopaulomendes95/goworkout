<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData, PageData } from './$types';
	import type { BackendWorkoutEntry } from '$lib/types';
    import { browser } from '$app/environment'; // For potential browser-only logic if needed

	const { data, form } = $props<{ data: PageData; form?: ActionData }>();

	// These $state variables are for the form inputs
	let workoutTitle = $state('');
	let workoutDescription = $state('');
	let workoutDuration = $state(0);
	let workoutCalories = $state(0);
	let workoutEntries = $state<Partial<BackendWorkoutEntry>[]>([]);
	
	let actionFeedback = $state<{ type: 'error'; message: string } | undefined>(undefined);
    let isInitialized = $state(false); // Flag to ensure initial setup from props happens once

	// Effect to initialize form fields from `data.workout` on first load,
    // and then from `form` if a submission fails (to repopulate with user's attempted changes).
	$effect(() => {
        const workoutDataFromLoad = data.workout;
        const formDataFromAction = form;

        console.log('[Edit Page Svelte] Effect for state init/update. Form:', formDataFromAction, 'Data.workout:', workoutDataFromLoad, 'Initialized:', isInitialized);

        if (formDataFromAction && formDataFromAction.formError) {
            // If form has errors (failed submission), prioritize repopulating from form data
            console.log('[Edit Page Svelte] Repopulating from failed form submission.');
            workoutTitle = formDataFromAction.title || '';
            workoutDescription = formDataFromAction.description || '';
            workoutDuration = formDataFromAction.durationMinutes ? parseInt(formDataFromAction.durationMinutes) : 0;
            workoutCalories = formDataFromAction.caloriesBurned ? parseInt(formDataFromAction.caloriesBurned) : 0;
            if (formDataFromAction.entries && typeof formDataFromAction.entries === 'string') {
                try { 
                    const parsedEntries = JSON.parse(formDataFromAction.entries);
                    workoutEntries = Array.isArray(parsedEntries) && parsedEntries.length > 0 ? parsedEntries : createDefaultEntries();
                } catch { workoutEntries = createDefaultEntries(); }
            } else {
                 workoutEntries = createDefaultEntries();
            }
            actionFeedback = { type: 'error', message: formDataFromAction.formError };
            isInitialized = true; // Mark as initialized even on form error to prevent data.workout override
        } else if (workoutDataFromLoad && !isInitialized) {
            // If no form error and not yet initialized, populate from data.workout
            console.log('[Edit Page Svelte] Initializing from data.workout.');
            workoutTitle = workoutDataFromLoad.title || '';
            workoutDescription = workoutDataFromLoad.description || '';
            workoutDuration = workoutDataFromLoad.duration_minutes || 0;
            workoutCalories = workoutDataFromLoad.calories_burned || 0;
            workoutEntries = workoutDataFromLoad.entries?.map(e => ({ ...e })) || createDefaultEntries();
            if (workoutEntries.length === 0) {
                workoutEntries = createDefaultEntries();
            }
            isInitialized = true;
        } else if (!workoutDataFromLoad && !formDataFromAction && !isInitialized) {
            // Case: No workout data from load (e.g. 404 or error) and no form submission yet
            console.log('[Edit Page Svelte] No workout data from load and no form data.');
            if (data.error) {
                 actionFeedback = { type: 'error', message: typeof data.error === 'string' ? data.error : data.error.message };
            }
            workoutEntries = createDefaultEntries(); // Ensure entries is an array
            isInitialized = true; // Mark initialized to prevent re-triggering this path
        }
	});

    function createDefaultEntries() {
        return [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }];
    }

	// Effect for clearing actionFeedback message after a timeout
	$effect(() => {
		if (actionFeedback) {
			const timer = setTimeout(() => { 
                console.log('[Edit Page Svelte] Clearing actionFeedback.');
                actionFeedback = undefined; 
            }, 4000);
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

	{#if (!data.workout && !form?.title) && !actionFeedback && !isInitialized}
		<p class="text-yellow-400">Loading workout data...</p>
	{/if}

	{#if actionFeedback}
		<div class={`mb-3 text-sm p-3 rounded-md ${actionFeedback.type === 'error' ? 'bg-red-600 text-red-100 border border-red-500' : 'bg-green-600 text-green-100 border border-green-500'}`}>
			{actionFeedback.message}
		</div>
	{/if}

	{#if (data.workout || form?.title) || (!data.workout && isInitialized) } <!-- Render form if data exists, or if form error, or if initialized after error -->
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
	{:else if !actionFeedback} <!-- Show loading/error only if not already showing an actionFeedback -->
        <p class="text-yellow-400">Workout data not available or an error occurred.</p>
    {/if}
</div>
