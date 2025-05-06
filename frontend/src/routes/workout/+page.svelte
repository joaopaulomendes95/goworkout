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

<!-- Todos.svelte -->
<div class="todoapp stack-large">
	<!-- NewTodo -->
	<form>
		<h2 class="label-wrapper">
			<label for="todo-0" class="label__lg"> What needs to be done? </label>
		</h2>
		<input type="text" id="todo-0" autocomplete="off" class="input input__lg" />
		<button type="submit" disabled="" class="btn btn__primary btn__lg"> Add </button>
	</form>

	<!-- Filter -->
	<div class="filters btn-group stack-exception">
		<button class="btn toggle-btn" aria-pressed="true">
			<span class="visually-hidden">Show</span>
			<span>All</span>
			<span class="visually-hidden">tasks</span>
		</button>
		<button class="btn toggle-btn" aria-pressed="false">
			<span class="visually-hidden">Show</span>
			<span>Active</span>
			<span class="visually-hidden">tasks</span>
		</button>
		<button class="btn toggle-btn" aria-pressed="false">
			<span class="visually-hidden">Show</span>
			<span>Completed</span>
			<span class="visually-hidden">tasks</span>
		</button>
	</div>

	<!-- TodosStatus -->
	<h2 id="list-heading">2 out of 3 items completed</h2>

	<!-- Todos -->
	<ul role="list" class="todo-list stack-large" aria-labelledby="list-heading">
		<!-- todo-1 (editing mode) -->
		<li class="todo">
			<div class="stack-small">
				<form class="stack-small">
					<div class="form-group">
						<label for="todo-1" class="todo-label">
							New name for 'Create a Svelte starter app'
						</label>
						<input type="text" id="todo-1" autocomplete="off" class="todo-text" />
					</div>
					<div class="btn-group">
						<button class="btn todo-cancel" type="button">
							Cancel
							<span class="visually-hidden">renaming Create a Svelte starter app</span>
						</button>
						<button class="btn btn__primary todo-edit" type="submit">
							Save
							<span class="visually-hidden">new name for Create a Svelte starter app</span>
						</button>
					</div>
				</form>
			</div>
		</li>

		<!-- todo-2 -->
		<li class="todo">
			<div class="stack-small">
				<div class="c-cb">
					<input type="checkbox" id="todo-2" checked />
					<label for="todo-2" class="todo-label"> Create your first component </label>
				</div>
				<div class="btn-group">
					<button type="button" class="btn">
						Edit
						<span class="visually-hidden">Create your first component</span>
					</button>
					<button type="button" class="btn btn__danger">
						Delete
						<span class="visually-hidden">Create your first component</span>
					</button>
				</div>
			</div>
		</li>

		<!-- todo-3 -->
		<li class="todo">
			<div class="stack-small">
				<div class="c-cb">
					<input type="checkbox" id="todo-3" />
					<label for="todo-3" class="todo-label"> Complete the rest of the tutorial </label>
				</div>
				<div class="btn-group">
					<button type="button" class="btn">
						Edit
						<span class="visually-hidden">Complete the rest of the tutorial</span>
					</button>
					<button type="button" class="btn btn__danger">
						Delete
						<span class="visually-hidden">Complete the rest of the tutorial</span>
					</button>
				</div>
			</div>
		</li>
	</ul>
	<table>
		<tbody>
			<tr>
				<th> </th><td>
					<label for="name">Table</label>
				</td>
				<th> </th><td>
					<input type="text" name="name" id="name" size="60" value="name" required />
				</td>
			</tr>
		</tbody>
	</table>
	<!-- MoreActions -->
	<div class="btn-group">
		<button type="button" class="btn btn__primary">Check all</button>
		<button type="button" class="btn btn__primary">Remove completed</button>
	</div>
</div>
