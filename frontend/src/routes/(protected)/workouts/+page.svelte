<script lang="ts">
  import { enhance } from '$app/forms';
  import type { ActionData, PageData } from './$types';
  import type { BackendWorkout, BackendWorkoutEntry } from '$lib/types';

  let { data, form }: { data: PageData, form: ActionData } = $props();

  // Workouts loaded from the server
  let serverLoadedWorkouts = $derived(data.workouts || []);
  let errorMessage = $derived(data.error);


  // For client-side list of workouts added in this session (optional, if you want immediate UI update without reload)
  let sessionWorkouts = $state<BackendWorkout[]>([]);

  // Form state for adding a new workout
  let newWorkoutTitle = $state(form?.title || '');
  let newWorkoutDescription = $state(form?.description || '');
  let newWorkoutDuration = $state(form?.durationMinutes || 0);
  let newWorkoutCalories = $state(form?.caloriesBurned || 0);
  let newWorkoutEntries = $state<Partial<BackendWorkoutEntry>[]>(
    form?.entries ? JSON.parse(form.entries) : [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }]
  );


  function addEntryRow() {
    newWorkoutEntries.push({ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: newWorkoutEntries.length + 1, weight: null });
    newWorkoutEntries = [...newWorkoutEntries];
  }

  function removeEntryRow(index: number) {
    newWorkoutEntries.splice(index, 1);
    newWorkoutEntries.forEach((e, i) => e.order_index = i + 1);
    newWorkoutEntries = [...newWorkoutEntries];
  }

  $effect(() => {
    if (form?.success && form?.message?.includes("added successfully")) {
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
      newWorkoutEntries = [{ exercise_name: '', sets: 1, reps: null, duration_seconds: null, notes: '', order_index: 1, weight: null }];
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

  const inputClasses = "w-full p-2 rounded bg-gray-700 border border-gray-600 focus:ring-blue-500 focus:border-blue-500 text-white placeholder-gray-400";
  const labelClasses = "block text-sm font-medium mb-1 text-gray-300";
  const buttonClasses = "px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md font-medium text-white disabled:opacity-50";
  const smallButtonClasses = "text-xs px-2 py-1 rounded";

  // Combine server-loaded and session-added workouts for display
  // This is a simple way; you might want more sophisticated state management for larger apps
  let allWorkouts = $derived([...serverLoadedWorkouts, ...sessionWorkouts].sort((a,b) => b.id - a.id)); // Example sort

</script>

<div class="space-y-8">
  <h1 class="text-3xl font-bold">Manage Workouts</h1>

  <!-- Display Existing Workouts -->
  <section class="p-6 bg-gray-800 rounded-lg shadow-md">
    <h2 class="text-2xl font-semibold mb-4">Your Workouts</h2>
    {#if errorMessage}
      <p class="mb-3 text-sm text-red-400 bg-red-900 p-2 rounded">{errorMessage}</p>
    {/if}

    {#if allWorkouts.length > 0}
      <div class="space-y-4">
        {#each allWorkouts as workout (workout.id)}
          <div class="p-4 bg-gray-700 rounded-md shadow-sm">
            <h3 class="text-xl font-semibold text-blue-400">{workout.title}</h3>
            <p class="text-sm text-gray-300 mt-1">{workout.description}</p>
            {#if workout.entries && workout.entries.length > 0}
              <div class="mt-3 pl-4 border-l-2 border-gray-600">
                <h4 class="text-sm font-medium text-gray-200 mb-1">Entries:</h4>
                <ul class="list-disc list-inside space-y-1 text-xs">
                  {#each workout.entries as entry (entry.id)}
                    <li>
                      {entry.exercise_name}: {entry.sets} sets
                      {#if entry.reps} of {entry.reps} reps{/if}
                      {#if entry.duration_seconds} for {entry.duration_seconds}s{/if}
                      {#if entry.weight} at {entry.weight}kg{/if}
                      {#if entry.notes} <span class="text-gray-500 italic">- "{entry.notes}"</span>{/if}
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
      <p class="text-gray-400">You haven't added any workouts yet. Use the form below to add one!</p>
    {/if}
  </section>

  <!-- Add Workout Form (remains largely the same) -->
  <section class="p-6 bg-gray-800 rounded-lg shadow-md">
    <h2 class="text-2xl font-semibold mb-4">Add New Workout</h2>
    {#if form?.formError}
        <p class="mb-3 text-sm text-red-400 bg-red-900 p-2 rounded">{form.formError}</p>
    {/if}
    {#if form?.success && form?.message && !form.formError}
        <p class="mb-3 text-sm text-green-400 bg-green-900 p-2 rounded">{form.message}</p>
    {/if}

    <form method="POST" action="?/addWorkout" use:enhance class="space-y-4">
      <!-- Form fields for title, description, duration, calories -->
      <div>
        <label for="title" class={labelClasses}>Title</label>
        <input bind:value={newWorkoutTitle} type="text" id="title" name="title" required class={inputClasses} placeholder="e.g., Morning Strength"/>
      </div>
      <div>
        <label for="description" class={labelClasses}>Description</label>
        <input bind:value={newWorkoutDescription} type="text" id="description" name="description" required class={inputClasses} placeholder="e.g., Full body workout"/>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="durationMinutes" class={labelClasses}>Duration (minutes)</label>
          <input bind:value={newWorkoutDuration} type="number" id="durationMinutes" name="durationMinutes" min="0" required class={inputClasses} placeholder="60"/>
        </div>
        <div>
          <label for="caloriesBurned" class={labelClasses}>Calories Burned</label>
          <input bind:value={newWorkoutCalories} type="number" id="caloriesBurned" name="caloriesBurned" min="0" required class={inputClasses} placeholder="300"/>
        </div>
      </div>

      <h3 class="text-xl font-semibold pt-4 border-t border-gray-700">Entries</h3>
      <input type="hidden" name="entries" value={JSON.stringify(newWorkoutEntries.map((e, idx) => ({
          ...e,
          order_index: idx + 1,
          reps: e.reps && e.reps > 0 ? e.reps : null,
          duration_seconds: e.duration_seconds && e.duration_seconds > 0 ? e.duration_seconds : null,
          weight: e.weight && e.weight > 0 ? e.weight : null
      })))} />

      <div class="space-y-3">
        {#each newWorkoutEntries as entry, index (index)}
          <div class="p-3 border border-gray-700 rounded-md space-y-2 relative">
            <!-- svelte-ignore event_directive_deprecated -->
            <button type="button" on:click={() => removeEntryRow(index)}
                    class="absolute top-1 right-1 {smallButtonClasses} bg-red-700 hover:bg-red-600 text-white">×</button>
            <div>
              <label for={`entry_name_${index}`} class={labelClasses}>Exercise Name</label>
              <input bind:value={entry.exercise_name} type="text" id={`entry_name_${index}`} required class={inputClasses} placeholder="Bench Press"/>
            </div>
            <div class="grid grid-cols-2 md:grid-cols-3 gap-2">
              <div>
                <label for={`entry_sets_${index}`} class={labelClasses}>Sets</label>
                <input bind:value={entry.sets} type="number" id={`entry_sets_${index}`} min="1" required class={inputClasses} placeholder="3"/>
              </div>
              <div>
                <label for={`entry_reps_${index}`} class={labelClasses}>Reps (optional)</label>
                <!-- svelte-ignore event_directive_deprecated -->
                <input bind:value={entry.reps} type="number" id={`entry_reps_${index}`} min="0" class={inputClasses} placeholder="10" on:input={() => { if (entry.reps && entry.reps > 0) entry.duration_seconds = null; }}/>
              </div>
              <div>
                <label for={`entry_duration_${index}`} class={labelClasses}>Duration (s, optional)</label>
                <!-- svelte-ignore event_directive_deprecated -->
                <input bind:value={entry.duration_seconds} type="number" id={`entry_duration_${index}`} min="0" class={inputClasses} placeholder="e.g. 60" on:input={() => { if (entry.duration_seconds && entry.duration_seconds > 0) entry.reps = null; }}/>
              </div>
            </div>
             <div>
                <label for={`entry_weight_${index}`} class={labelClasses}>Weight (kg, optional)</label>
                <input bind:value={entry.weight} type="number" step="any" id={`entry_weight_${index}`} min="0" class={inputClasses} placeholder="60.5"/>
            </div>
            <div>
              <label for={`entry_notes_${index}`} class={labelClasses}>Notes</label>
              <input bind:value={entry.notes} type="text" id={`entry_notes_${index}`} class={inputClasses} placeholder="e.g., Tempo 2-0-1-0"/>
            </div>
          </div>
        {/each}
      </div>
      <!-- svelte-ignore event_directive_deprecated -->
      <button type="button" on:click={addEntryRow} class="text-sm text-blue-400 hover:text-blue-300 {smallButtonClasses} border border-blue-400 hover:bg-blue-900">
        + Add Entry Row
      </button>

      <div class="pt-4">
        <button type="submit" class={buttonClasses}>Save Workout</button>
      </div>
    </form>
  </section>
</div>
