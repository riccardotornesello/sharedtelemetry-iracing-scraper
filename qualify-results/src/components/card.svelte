<script lang="ts">
	interface CardProps {
		title?: string;
		link: string;
		children?: () => any;
	}

	let { title, link, children }: CardProps = $props();
</script>

<a
	href={link}
	class="block rounded-lg border border-gray-200 bg-white p-6 shadow hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-700"
>
	{#if title}
		<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
			{title}
		</h5>
	{/if}

	<div>
		{@render children?.()}
	</div>
</a>
