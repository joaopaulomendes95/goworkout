<script lang="ts">
	const { data } = $props();

	console.log('Page Data Login', data);
</script>

<div class="main">
	<div class="absolute top-4 left-4">Go Workout</div>
	<button class="absolute top-4 right-4">Login</button>
	<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
		<div class="cabecalho">
			<div class="flex flex-col space-y-2 text-center">
				<div class="flex flex-col space-y-2 text-center">
					<h1 class="text-2xl font-semibold tracking-tight">Login to your account</h1>
					<p class="text-muted-foreground text-sm">Enter your email below to login your account</p>
				</div>
				<form action="user_login" method="post">
					<div class="flex flex-col text-center">
						<label for="username"><b>Username</b></label>
						<input type="text" placeholder="Enter Username" name="username" required />
						<label for="password"><b>Password</b></label>
						<input type="password" placeholder="Enter password" name="password" required />
						<button type="submit">Login</button>
						<label>
							<input type="checkbox" checked="checked" name="remember" /> Remember me
						</label>
					</div>
					<div class="container">
						<button type="button" class="cancelbtn">Cancel</button>
						<span class="psw">Forgot <a href="placeholder">password?</a></span>
					</div>
				</form>
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
		</div>
	</div>
</div>
