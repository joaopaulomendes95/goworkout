<main class="flex min-h-[100vh] flex-col items-center justify-center">
	<h1>Workout Crud</h1>
	<div class="workout_with_entries">
		<div class="workout">
			<form id="workout" method="post">
				<label for="title">Title</label>
				<input required type="text" id="title" name="title" placeholder="Add workout title" />
				<br />
				<label for="description">Description</label>
				<input
					required
					type="text"
					id="description"
					name="description"
					placeholder="Add workout description"
				/>
				<br />
				<label for="duration_minutes">Duration</label>
				<input
					required
					type="number"
					id="duration_minutes"
					name="duration_minutes"
					placeholder="Add workout duration"
				/>
				<br />
				<label for="calories_burned">Calories Burned</label>
				<input
					required
					type="number"
					id="calories_burned"
					name="calories_burned"
					placeholder="Add calories burned"
				/>
				<br />
			</form>
		</div>

		<!-- Each workout has many workout_entries -->
		<div class="workout_entries">
			<form id="workout" method="post">
				<label for="exercise_name">Exercise Name</label>
				<input
					required
					type="text"
					id="exercise_name"
					name="exercise_name"
					placeholder="Add exercise name"
				/>
				<br />
				<label for="sets">Sets</label>
				<input required type="number" id="sets" name="sets" placeholder="Add number of sets" />
				<br />
				<label for="reps">Reps</label>
				<input required type="number" id="reps" name="reps" placeholder="Add number of reps" />
				<br />
				<label for="duration_seconds">Duration</label>
				<input
					required
					type="number"
					id="duration_seconds"
					name="duration_seconds"
					placeholder="Add number of seconds"
				/>
				<br />
				<label for="weight">Weight</label>
				<input required type="number" id="weight" name="weight" placeholder="Add weight" />
				<br />
				<label for="notes">Notes</label>
				<input required type="text" id="notes" name="notes" placeholder="Add some nomes" />
				<br />
				<label for="order_index">Order index</label>
				<input
					required
					type="number"
					id="order_index"
					name="order_index"
					placeholder="Add index order"
				/>
				<br />
				<!--workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index -->
				<button>Save/Edit/Delete</button>
			</form>
		</div>
	</div>
</main>
