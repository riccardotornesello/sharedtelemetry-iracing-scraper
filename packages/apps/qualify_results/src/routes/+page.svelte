<script>
	import { Mail } from 'lucide-svelte';
</script>

<div class="relative min-h-screen w-full overflow-hidden bg-black text-white">
	<!-- Background image with overlay -->
	<div
		class="absolute inset-0 z-0 opacity-40"
		style="
        background-size: cover;
        background-position: center;
      "
	></div>

	<!-- Content -->
	<div class="relative z-10 mx-auto max-w-6xl px-4 py-20 sm:px-6 lg:px-8">
		<div class="space-y-8 text-center">
			<h1 class="text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
				iRacing Championship Results Platform
			</h1>
			<p class="mx-auto max-w-3xl text-xl text-gray-300 sm:text-2xl">
				Manage and track your racing competitions with precision. Free platform for organizing
				races, qualifying sessions, and monitoring lap times.
			</p>
			<div class="grid grid-cols-1 gap-8 py-12 md:grid-cols-3">
				<div class="rounded-lg bg-gray-900/50 p-6 backdrop-blur">
					<h3 class="mb-2 text-xl font-semibold">Race Management</h3>
					<p class="text-gray-400">Organize and manage your racing events with ease</p>
				</div>
				<div class="rounded-lg bg-gray-900/50 p-6 backdrop-blur">
					<h3 class="mb-2 text-xl font-semibold">Qualifying Sessions</h3>
					<p class="text-gray-400">Track qualifying times and grid positions</p>
				</div>
				<div class="rounded-lg bg-gray-900/50 p-6 backdrop-blur">
					<h3 class="mb-2 text-xl font-semibold">Lap Time Tracking</h3>
					<p class="text-gray-400">Monitor and compare lap times across all tracks</p>
				</div>
			</div>
			<div class="mt-12">
				<a
					href="mailto:dev@sharedtelemetry.com"
					class="inline-flex items-center rounded-full bg-white px-8 py-4 text-lg font-medium text-black transition-colors hover:bg-gray-100"
				>
					<Mail class="mr-2 h-5 w-5" />
					Contact Us to Get Started
				</a>
				<p class="mt-4 text-sm text-gray-400">
					Email us at dev@sharedtelemetry.com to add your competition
				</p>
			</div>
		</div>
	</div>
</div>
