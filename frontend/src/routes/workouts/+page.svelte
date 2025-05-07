<script lang="ts">
	const data = $props();

	$effect(() => {
		const getData = async () => {
			const response = await fetch('/api/workouts');
			const data = await response.json();
			return data;
		};
	});
</script>

<main class="flex min-h-[100vh] flex-col items-center justify-center">
	<h1>Workout Crud</h1>
	<div class="workout_with_entries flex min-h-[100vh] flex-col items-center justify-center">
		<div class="workout flex min-h-[100vh] flex-col items-center justify-center">
			<form id="workout" method="post">
				<label for="title">Title</label>
				<input required type="text" id="title" name="title" placeholder="Add workout title" />
				<label for="description">Description</label>
				<input
					required
					type="text"
					id="description"
					name="description"
					placeholder="Add workout description"
				/>
				<label for="duration_minutes">Duration</label>
				<input
					required
					type="number"
					id="duration_minutes"
					name="duration_minutes"
					placeholder="Add workout duration"
				/>
				<label for="calories_burned">Calories Burned</label>
				<input
					required
					type="number"
					id="calories_burned"
					name="calories_burned"
					placeholder="Add calories burned"
				/>
			</form>
		</div>

		<!-- Each workout has many workout_entries -->
		<!-- Maybe make this repeat a couple of times ->
		<!-- If a user wants to add more entries -->
		<!-- It should generate a new table row -->
		<div class="workout_entries flex min-h-[100vh] flex-col items-center justify-center">
			<form id="workout" method="post">
				<table>
					<thead>
						<tr>
							<th>Sets</th>
							<th>Reps</th>
							<th>Duration</th>
							<th>Weight</th>
							<th>Notes</th>
							<th>Index Order</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>
								<label for="sets">Sets</label>
								<input
									required
									type="number"
									id="sets"
									name="sets"
									placeholder="Add number of sets"
								/>
							</td>
							<td>
								<label for="reps">Reps</label>
								<input
									required
									type="number"
									id="reps"
									name="reps"
									placeholder="Add number of reps"
								/>
							</td>
							<td>
								<label for="duration_seconds">Duration</label>
								<input
									required
									type="number"
									id="duration_seconds"
									name="duration_seconds"
									placeholder="Add number of seconds"
								/>
							</td>
							<td>
								<label for="weight">Weight</label>
								<input required type="number" id="weight" name="weight" placeholder="Add weight" />
							</td>
							<td>
								<label for="notes">Notes</label>
								<input required type="text" id="notes" name="notes" placeholder="Add some nomes" />
							</td>
							<td>
								<label for="order_index">Order index</label>
								<input
									required
									type="number"
									id="order_index"
									name="order_index"
									placeholder="Add index order"
								/>
							</td>
						</tr>
					</tbody>
					<!-- placeholer for a footer in the table -->
					<!-- <tfoot><tr></tr></tfoot> -->
				</table>
				<button>Save/Edit/Delete</button>
			</form>

			<div>
				{#each data.workouts as workout (workout.id)}
					<div>
						<h2>{workout.title}</h2>
						<p>{workout.description}</p>
						<p>Duration: {workout.duration_minutes} minutes</p>
						<p>Calories Burned: {workout.calories_burned}</p>
					</div>
				{/each}
			</div>
		</div>
	</div>
</main>
