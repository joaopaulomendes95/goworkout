<script lang="ts">
	const { data } = $props();

	console.log('page Data', data);

	let health = data.health;
	console.log('health', health);
</script>

<main class="flex min-h-[100vh] flex-col items-center justify-center">
	<h1>Homepage</h1>

	<nav>
		<a href="/login">Login</a>
		<a href="/register">Register</a>
	</nav>
</main>

<!-- <main class="flex min-h-[100vh] flex-col items-center justify-center"> -->
<!-- 	<h1 class="text-2xl text-orange-300">Status: {health.status}</h1> -->
<!-- 	<h1 class="text-2xl text-orange-300">Message: {health.message}</h1> -->
<!-- </main> -->
<!---->
<!-- <main class="flex min-h-[100vh] flex-col items-center justify-center"> -->
<!-- 	<form method="POST" class="flex flex-col items-center justify-center"> -->
<!-- 		<label for="username" class="text-xl font-bold">Username</label> -->
<!-- 		<input bind:value={username} name="username" type="text" class="rounded-xl text-black" /> -->
<!---->
<!-- 		<label for="email" class="text-xl font-bold">Email</label> -->
<!-- 		<input bind:value={email} name="email" type="text" class="rounded-xl text-black" /> -->
<!---->
<!-- 		<label for="password" class="text-xl font-bold">Password</label> -->
<!-- 		<input bind:value={password} name="password" type="password" class="rounded-xl text-black" /> -->
<!---->
<!-- 		<label for="bio" class="text-xl font-bold">Bio</label> -->
<!-- 		<textarea bind:value={bio} name="bio" class="rounded-xl text-black" rows="4" cols="50" -->
<!-- 		></textarea> -->
<!---->
<!-- 		<button type="submit" class="rounded-xl bg-orange-300 p-2">Register</button> -->
<!-- 	</form> -->
<!-- </main> -->
