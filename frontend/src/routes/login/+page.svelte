<script lang="ts">
	const { data } = $props();

	console.log('Page Data Login', data);

	// user info
	// let username = '';
	// let email = '';
	// let password = '';
	// let bio = '';
	// idk, gpt
	// let error: string | null = null;
	// let loading = false;
</script>

<div class="md:hidden"></div>
<div
	class="relative container hidden h-[800px] flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0"
>
	<button class="absolute top-4 right-4 md:top-8 md:right-8"> Login </button>
	<div class="bg-muted relative hidden h-full flex-col p-10 text-white lg:flex dark:border-r">
		<div class="absolute inset-0 bg-cover"></div>
		<div class="relative z-20 flex items-center text-lg font-medium">
			<!-- <Command class="mr-2 h-6 w-6" /> -->
			Go Workout
		</div>
		<div class="relative z-20 mt-auto">
			<blockquote class="space-y-2">
				<p class="text-lg">&ldquo;Postgres is better than Tailwind&rdquo;</p>
				<footer class="text-sm">Bill Gates</footer>
			</blockquote>
		</div>
	</div>
	<div class="lg:p-8">
		<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
			<div class="flex flex-col space-y-2 text-center">
				<h1 class="text-2xl font-semibold tracking-tight">Create an account</h1>
				<p class="text-muted-foreground text-sm">Enter your email below to create your account</p>
			</div>
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
</div>
