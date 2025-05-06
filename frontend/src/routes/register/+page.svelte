<script lang="ts">
	const { form } = $props();

	console.log('Page Data Register', form);
</script>

<div class="main">
	<div class="absolute top-4 left-4">Go Workout</div>
	<button class="absolute top-4 right-4">Register</button>
	<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px] md:w-[600px]">
		<div class="cabecalho">
			<div class="flex flex-col space-y-2 text-center">
				<h1 class="text-2xl font-semibold tracking-tight">Create account</h1>
				<p class="text-muted-foreground text-sm">Enter the details below to create a account</p>
			</div>
			<form action="?/user_register" method="POST">
				<div class="flex flex-col text-center text-black">
					<label for="username"><b>Username</b></label>
					<input name="username" type="text" placeholder="Enter Username" required />
					<label for="email"><b>Email</b></label>
					<input name="email" type="text" placeholder="Enter Email" required />
					<label for="password"><b>Password</b></label>
					<input name="password" type="password" placeholder="Enter Password" required />
					<br />
					<label for="bio"><b>Bio</b></label>
					<textarea name="bio" id="bio" placeholder="Enter your bio (optional)" required></textarea>
					<br />
					<button type="submit">Register</button>
					<br />
					<label>
						<input type="checkbox" name="remember" /> Remember me
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
