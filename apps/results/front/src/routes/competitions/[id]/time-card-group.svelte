<script lang="ts">
	import TimeCard from '../../../components/time-card.svelte';
	import type { EventGroup } from '$lib/api/competition';
	import dayjs from 'dayjs';

	let {
		event,
		eventResults,
		eventOverallBest
	}: {
		event: EventGroup;
		eventResults: Record<string, number>;
		eventOverallBest: number;
	} = $props();

	let personalBest = Math.min(...Object.values(eventResults || {}).filter(Boolean));
</script>

{#each event.sessions as session}
	<td class="px-2 py-1">
		<TimeCard
			isPersonalBest={eventResults?.[dayjs(session.fromTime).format('YYYY-MM-DD')] == personalBest}
			isOverallBest={eventResults?.[dayjs(session.fromTime).format('YYYY-MM-DD')] ==
				eventOverallBest}
			time={eventResults?.[dayjs(session.fromTime).format('YYYY-MM-DD')]}
		/>
	</td>
{/each}
