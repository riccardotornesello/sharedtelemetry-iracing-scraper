<script lang="ts">
	import type { PageData } from './$types';
	import Card from '../../components/card.svelte';
	import SessionCard from '../../components/session-card.svelte';

	let { data }: { data: PageData } = $props();
</script>

{#each data.sessions as session}
	<SessionCard
		subsessionId={session.subsession_id}
		simsessionNumber={session.simsession_number}
		launchAt={session.launch_at}
		trackId={session.track_id}
	/>
{/each}
