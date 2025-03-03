<script lang="ts">
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	import type { PageData } from './$types';
	import * as m from '$lib/paraglide/messages.js';
	import Ranking from './ranking.svelte';
	import RankingCrew from './ranking-crew.svelte';

	dayjs.extend(duration);

	function formatDate(date: string): string {
		return dayjs(date).format('DD/MM');
	}

	let { data }: { data: PageData } = $props();

	let competitionWithCrews = data.competitionRanking.competition.crewDriversCount > 1;
	let showOverall = $state(!competitionWithCrews);
</script>

<svelte:head>
	<title>{data.competitionRanking.competition.name} - Results</title>
</svelte:head>

<h1 class="w-full p-5 text-center text-2xl font-bold text-gray-200">
	{data.competitionRanking.competition.name}
</h1>

{#if data.competitionRanking.competition.crewDriversCount > 1}
	<label class="inline-flex cursor-pointer items-center">
		<input
			type="checkbox"
			value={showOverall}
			onchange={() => (showOverall = !showOverall)}
			class="peer hidden"
		/>
		<div
			class="peer relative h-6 w-11 rounded-full bg-gray-200 after:absolute after:start-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rtl:peer-checked:after:-translate-x-full dark:border-gray-600 dark:bg-gray-700 dark:peer-checked:bg-blue-600 dark:peer-focus:ring-blue-800"
		></div>
		<span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300">
			{m.general_ranking()}
		</span>
	</label>
{/if}

<div class="w-full overflow-x-auto">
	<table class="min-w-full table-auto text-left text-sm text-gray-400 rtl:text-right">
		<thead class="bg-gray-700 text-center text-xs uppercase text-gray-400">
			<tr>
				<th scope="col" class="px-6 py-3" rowspan="2" colspan="3">{m.driver()}</th>
				<th scope="col" class="px-6 py-3" rowspan="2">{m.team()}</th>
				<th scope="col" class="px-6 py-3" rowspan="2">{m.sum()}</th>
				{#each data.competitionRanking.eventGroups as eventGroup}
					<th scope="col" class="px-6 py-3" colspan={eventGroup.dates?.length}>
						{eventGroup.name}
					</th>
				{/each}
			</tr>
			<tr>
				{#each data.competitionRanking.eventGroups as eventGroup}
					{#each eventGroup.dates as date}
						<th scope="col" class="px-6 py-3"> {formatDate(date)}</th>
					{/each}
				{/each}
			</tr>
		</thead>

		<tbody>
			{#if showOverall === true}
				<Ranking
					ranking={data.competitionRanking.ranking}
					drivers={data.competitionRanking.drivers}
					eventGroups={data.competitionRanking.eventGroups}
					overallBest={data.overallBest}
					classes={data.classes}
					showCrew={competitionWithCrews}
				/>
			{:else}
				{#each data.crews as crew, index}
					<RankingCrew
						position={index + 1}
						{crew}
						drivers={data.competitionRanking.drivers}
						eventGroups={data.competitionRanking.eventGroups}
						overallBest={data.overallBest}
						classes={data.classes}
					/>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
