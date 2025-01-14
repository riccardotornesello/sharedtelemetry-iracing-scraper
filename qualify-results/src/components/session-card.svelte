<script lang="ts">
	import Card from './card.svelte';

	interface Props {
		subsessionId: number;
		simsessionNumber: number;
		launchAt: string;
		trackId: number;
	}

	let { subsessionId, simsessionNumber, launchAt, trackId }: Props = $props();
</script>

<Card link="#">
	<p>{subsessionId} | {simsessionNumber} | {launchAt} | {trackId}</p>
</Card>
