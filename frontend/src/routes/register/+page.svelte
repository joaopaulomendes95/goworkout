<script lang="ts">
	const { data } = $props();

	console.log('Page Data Register', data);
</script>

<div class="main">
	<div class="absolute top-4 left-4">Go Workout</div>
	<button class="absolute top-4 right-4">Register</button>
	<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
		<div class="cabecalho">
			<div class="flex flex-col space-y-2 text-center">
				<h1 class="text-2xl font-semibold tracking-tight">Create account</h1>
				<p class="text-muted-foreground text-sm">Enter the details below to create a account</p>
			</div>
			<form action="user_register" method="post">
				<div class="flex flex-col text-center">
					<label for="username"><b>Username</b></label>
					<input type="text" placeholder="Enter Username" name="username" required />
					<label for="email"><b>Email</b></label>
					<input type="text" placeholder="Enter Email" name="email" required />
					<label for="password"><b>Password</b></label>
					<input type="password" placeholder="Enter Password" name="password" required />
					<br />
					<button type="submit">Register</button>
					<br />
					<label>
						<input type="checkbox" checked="checked" name="remember" /> Remember me
					</label>
				</div>
			</form>
			<p class="text-muted-foreground px-8 text-center text-sm">
				By clicking continue, you agree to our
				<a href="/terms" class="hover:text-primary underline underline-offset-4">
					Terms of Service
				</a>
				and
				<a href="/privacy" class="hover:text-primary underline underline-offset-4">
					Privacy Policy
				</a>
				.
			</p>
		</div>
	</div>
	<div class="absolute bottom-4 left-4">
		<blockquote class="quote">
			<p class="text">&ldquo; You miss 100% of the shots you don't take. -Wayne Gretzky.&rdquo;</p>
			<footer class="author">Michael Scott</footer>
		</blockquote>
	</div>
</div>
