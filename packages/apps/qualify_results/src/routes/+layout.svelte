<script lang="ts">
	import { i18n } from '$lib/i18n';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import '../app.css';
	import Footer from '../components/footer.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<title>Sharedtelemetry Results</title>
</svelte:head>

<ParaglideJS {i18n}>
	<div class="dark flex min-h-screen flex-col bg-[#111a22] text-white">
		<main class="flex-1">
			{@render children()}
		</main>

		<Footer />
	</div>
</ParaglideJS>
