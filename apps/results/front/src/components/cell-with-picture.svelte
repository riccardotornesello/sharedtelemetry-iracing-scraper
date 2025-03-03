<script lang="ts">
	interface Props {
		picture?: string;
		alt?: string;
		children: any;
	}

	let { picture, alt, children }: Props = $props();
</script>

<div class="flex items-center space-x-2">
	{#if picture}
		<img src={picture} {alt} class="mr-4 h-6 w-6 rounded-full object-cover" />
	{/if}

	{@render children()}
</div>
