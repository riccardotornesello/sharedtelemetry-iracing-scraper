<script lang="ts">
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	import type { PageData } from './$types';
	import * as m from '$lib/paraglide/messages.js';
	import Ranking from './ranking.svelte';

	dayjs.extend(duration);

	function formatDate(date: string): string {
		return dayjs(date).format('DD/MM');
	}

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>{data.competition.name} - Results</title>
</svelte:head>

<h1 class="w-full p-5 text-center text-2xl font-bold text-gray-200">
	{data.competition.name}
</h1>

<div class="w-full overflow-x-auto">
	<table class="min-w-full table-auto text-left text-sm text-gray-400 rtl:text-right">
		<thead class="bg-gray-700 text-center text-xs uppercase text-gray-400">
			<tr>
				<th scope="col" class="px-6 py-2" rowspan="2" colspan="2">{m.driver()}</th>
				<th scope="col" class="px-6 py-2" rowspan="2">{m.team()}</th>
				<th scope="col" class="px-6 py-2" rowspan="2">{m.sum()}</th>
				{#each data.competition.eventGroups as eventGroup}
					<th scope="col" class="px-6 py-2" colspan={eventGroup.sessions.length}>
						{eventGroup.name}
					</th>
				{/each}
			</tr>
			<tr>
				{#each data.competition.eventGroups as eventGroup}
					{#each eventGroup.sessions as session}
						<th scope="col" class="px-6 py-2"> {formatDate(session.fromTime)}</th>
					{/each}
				{/each}
			</tr>
		</thead>

		<tbody>
			<Ranking
				ranking={data.driversRanking}
				drivers={data.driversMap}
				events={data.competition.eventGroups}
				driversCrewMap={data.driversCrewMap}
				driversTeamMap={data.driversTeamMap}
				overallBest={data.overallBest}
			/>
		</tbody>
	</table>
</div>
