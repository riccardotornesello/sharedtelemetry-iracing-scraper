<script lang="ts">
	import TimeCard from '../../../components/time-card.svelte';
	import Ranking from './ranking.svelte';

	let { crew, drivers, eventGroups, overallBest } = $props();

	let isOpen = $state(false)
</script>

<tr class="border-b border-gray-700 bg-gray-900">
	<td>
		<button class="px-6 py-4 text-center" onclick={() => (isOpen = !isOpen)}>
			{isOpen ? '▼' : '►'}
		</button>
	</td>
	<td class="px-6 py-4">{crew.name}</td>
	<td class="px-2 py-2"><TimeCard time={crew.sum} /></td>
	{#each eventGroups as eventGroup}
		{#each eventGroup.dates as date}
			<td></td>
		{/each}
	{/each}
</tr>

{#if isOpen}
	<Ranking ranking={crew.ranking} {drivers} {eventGroups} {overallBest} />
{/if}
