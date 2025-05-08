<!-- frontend/src/routes/(protected)/workouts/+page.svelte -->
<script lang="ts">
  import { enhance } from '$app/forms';
  import type { ActionData, PageData } from './$types';
   import type { BackendWorkout, BackendWorkoutEntry } from '$lib/types';

  let { data, form }: { data: PageData, form: ActionData } = $props();

  // For client-side list of workouts added in this session (optional)
  let sessionWorkouts = $state<BackendWorkout[]>([]);

  // Form state for adding a new workout
  let newWorkoutTitle = $state(form?.title || ''); // Repopulate on error
  let newWorkoutDescription = $state(form?.description || '');
  let newWorkoutDuration = $state(form?.durationMinutes || 0);
  let newWorkoutCalories = $state(form?.caloriesBurned || 0);
  let newWorkoutEntries = $state<Partial<BackendWorkoutEntry>[]>(
    form?.entries ? JSON.parse(form.entries) : [{ exercise_name: '', sets: 1, reps: 10, notes: '', order_index: 1 }]
  );


  function addEntryRow() {
    newWorkoutEntries.push({ exercise_name: '', sets: 1, reps: 1, notes: '', order_index: newWorkoutEntries.length + 1 });
    newWorkoutEntries = [...newWorkoutEntries];
  }

  function removeEntryRow(index: number) {
    newWorkoutEntries.splice(index, 1);
    newWorkoutEntries.forEach((e, i) => e.order_index = i + 1); // Re-index from 1
    newWorkoutEntries = [...newWorkoutEntries];
  }

  // Handle form action result
  $effect(() => {
    if (form?.success && form?.message?.includes("added successfully")) {
      console.log("Workout added form details:", form);
      if (form.createdWorkout) { // If the action returned the created workout
        sessionWorkouts.push(form.createdWorkout as BackendWorkout);
        sessionWorkouts = [...sessionWorkouts];
      }
      // Reset form
      newWorkoutTitle = '';
      newWorkoutDescription = '';
      newWorkoutDuration = 0;
      newWorkoutCalories = 0;
      newWorkoutEntries = [{ exercise_name: '', sets: 1, reps: 10, notes: '', order_index: 1 }];
      // No need to invalidate('app:workouts') if there's no server-side list to fetch
    } else if (form && !form.success && form.title !== undefined) { // Repopulate on error
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
</script>

<div class="space-y-8">
  <h1 class="text-3xl font-bold">Manage Workouts</h1>

  <!-- Add Workout Form -->
  <section class="p-6 bg-gray-800 rounded-lg shadow-md">
    <h2 class="text-2xl font-semibold mb-4">Add New Workout</h2>
    {#if form?.formError}
        <p class="mb-3 text-sm text-red-400 bg-red-900 p-2 rounded">{form.formError}</p>
    {/if}
    {#if form?.success && form?.message}
        <p class="mb-3 text-sm text-green-400 bg-green-900 p-2 rounded">{form.message}</p>
    {/if}

    <form method="POST" action="?/addWorkout" use:enhance class="space-y-4">
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
      <input type="hidden" name="entries" value={JSON.stringify(newWorkoutEntries.map((e, idx) => ({...e, order_index: idx + 1})))} />

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
                <input bind:value={entry.reps} type="number" id={`entry_reps_${index}`} min="0" class={inputClasses} placeholder="10"/>
              </div>
              <div>
                <label for={`entry_duration_${index}`} class={labelClasses}>Duration (s, optional)</label>
                <input bind:value={entry.duration_seconds} type="number" id={`entry_duration_${index}`} min="0" class={inputClasses} placeholder="e.g. 60"/>
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

  <!-- Optional: Display workouts added in this session (client-side list) -->
  {#if sessionWorkouts.length > 0}
  <section class="mt-10">
    <h2 class="text-2xl font-semibold mb-4">Workouts Added This Session</h2>
    <div class="space-y-6">
      {#each sessionWorkouts as workout (workout.id)}
        <div class="bg-gray-800 p-5 rounded-lg shadow">
          <h3 class="text-xl font-semibold text-blue-400">{workout.title}</h3>
          <p class="text-gray-400 text-sm">{workout.description}</p>
          <!-- ... display more workout details and entries ... -->
        </div>
      {/each}
    </div>
  </section>
  {/if}
</div>
