<script lang="ts">
	import type { PageData } from './$types';
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	import Ranking from './ranking.svelte';
	import RankingCrew from './ranking-crew.svelte';

	dayjs.extend(duration);

	function formatDate(date: string): string {
		return dayjs(date).format('DD/MM');
	}

	let { data }: { data: PageData } = $props();

	let overallBest = {};
	for (const eventGroup of data.eventGroups) {
		overallBest[eventGroup.id] = Math.min(
			...data.ranking.map((r) => Object.values(r.results?.[eventGroup.id] || {})).flat()
		);
	}

	let rankingType = 'overall';
</script>

<h1 class="p-5 text-center text-2xl font-bold text-gray-200">
	{data.competition.name}
</h1>

<table class="w-full table-auto text-left text-sm text-gray-400 rtl:text-right">
	<thead class="bg-gray-700 text-center text-xs uppercase text-gray-400">
		<tr>
			<th scope="col" class="px-6 py-3" rowspan="2" colspan="2">Pilota</th>
			<th scope="col" class="px-6 py-3" rowspan="2">Somma</th>
			{#each data.eventGroups as eventGroup}
				<th scope="col" class="px-6 py-3" colspan={eventGroup.dates?.length}>{eventGroup.name}</th>
			{/each}
		</tr>
		<tr>
			{#each data.eventGroups as eventGroup}
				{#each eventGroup.dates as date}
					<th scope="col" class="px-6 py-3"> {formatDate(date)}</th>
				{/each}
			{/each}
		</tr>
	</thead>

	<tbody>
		{#if rankingType === 'overall'}
			<Ranking
				ranking={data.ranking}
				drivers={data.drivers}
				eventGroups={data.eventGroups}
				{overallBest}
			/>
		{:else if rankingType === 'crew'}
			{#each data.crews as crew}
				<RankingCrew {crew} drivers={data.drivers} eventGroups={data.eventGroups} {overallBest} />
			{/each}
		{/if}
	</tbody>
</table>

{#snippet ranking(ranking)}{/snippet}
