<script lang="ts">
	import { i18n } from '$lib/i18n';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import '../app.css';
	let { children } = $props();
</script>

<ParaglideJS {i18n}>
	<div class="dark min-h-screen bg-[#111a22] text-white">
		{@render children()}
	</div>
</ParaglideJS>
