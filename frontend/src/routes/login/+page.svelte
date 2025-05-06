<script lang="ts">
	const { data, form } = $props();

	let username = $state('');
	let password = $state('');

	console.log('Page Data Login', data);
	$inspect(form);
</script>

<div class="main">
	<div class="absolute top-4 left-4">Go Workout</div>
	<button class="absolute top-4 right-4">Login</button>
	<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
		<div class="cabecalho">
			<div class="flex flex-col space-y-2 text-center">
				<div class="flex flex-col space-y-2 text-center">
					<h1 class="text-2xl font-semibold tracking-tight">Log in to your account</h1>
					<p class="text-muted-foreground text-sm">Enter your email below to login your account</p>
				</div>
				<form action="?/user_login" method="post">
					<div class="flex flex-col text-center">
						<label for="username"><b>Username</b></label>
						<input
							id="username"
							placeholder="Enter Username"
							name="username"
							required
							bind:value={username}
						/>
						<label for="password"><b>Password</b></label>
						<input
							id="password"
							name="password"
							type="password"
							placeholder="Enter password"
							required
							bind:value={password}
						/>
						<button type="submit">Login</button>
						<label>
							<input type="checkbox" checked={false} name="remember" /> Remember me
						</label>
					</div>
					<div class="container">
						<button type="button" class="cancelbtn">Cancel</button>
						<span class="psw">Forgot <a href="placeholder">password?</a></span>
					</div>
				</form>

				{#if form?.success === false}
					<p class="mt-4 text-red-500">{form.message}</p>
				{/if}

				<p class="text-muted-foreground px-8 text-center text-sm">
					Dont have account yet?
					<a href="/register" class="hover:text-primary underline underline-offset-4">
						Click here to register
					</a>
				</p>
			</div>
			<div class="absolute bottom-4 left-4">
				<blockquote class="space-y-2">
					<p class="text-lg">&ldquo;Postgres is better than Tailwind&rdquo;</p>
					<footer class="text-sm">Bill Gates</footer>
				</blockquote>
			</div>
			<div>
				<h1 class="text-2xl text-red-500">{form?.message}</h1>
				<h1 class="text-2xl text-red-500">{form?.token}</h1>
			</div>
		</div>
	</div>
</div>
