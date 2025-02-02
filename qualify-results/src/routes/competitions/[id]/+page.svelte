<script lang="ts">
	import type { PageData } from './$types';
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	import TimeCard from '../../../components/time-card.svelte';

	dayjs.extend(duration);
	
	function formatDate(date: string): string {
		return dayjs(date).format('DD/MM');
	}
	
	let { data }: { data: PageData } = $props();
</script>

<table class="w-full table-auto text-left text-sm text-gray-500 rtl:text-right dark:text-gray-400">
	<thead class="bg-gray-700 text-center text-xs uppercase text-gray-400">
		<tr>
			<th scope="col" class="px-6 py-3" rowspan="2" colspan="2">Pilota</th>
			<th scope="col" class="px-6 py-3" rowspan="2">Somma</th>
			{#each data.eventGroups as eventGroup}
				<th scope="col" class="px-6 py-3" colspan={eventGroup.Dates.length}>{eventGroup.Name}</th>
			{/each}
		</tr>
		<tr>
			{#each data.eventGroups as eventGroup}
				{#each eventGroup.Dates as date}
					<th scope="col" class="px-6 py-3"> {formatDate(date)}</th>
				{/each}
			{/each}
		</tr>
	</thead>

	<tbody>
		{#each data.ranking as rank, index}
			<tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800">
				<td class="px-6 py-4 text-center">P{index + 1}</td>
				<td class="px-6 py-4">{data.drivers[rank.custId]?.Name || rank.custId}</td>
				<td class="px-2 py-2">
					<TimeCard time={rank.sum} />
				</td>
				{#each data.eventGroups as eventGroup}
					{#each eventGroup.Dates as date}
						<td class="px-2 py-2">
							<TimeCard
								isPersonalBest={rank.results?.[eventGroup.ID]?.[date] == Math.min(...Object.values(rank.results?.[eventGroup.ID] || {}))}
								isOverallBest={rank.results?.[eventGroup.ID]?.[date] == Math.min(...data.ranking.map(r => Object.values(r.results?.[eventGroup.ID] || {})).flat())}
								time={rank.results?.[eventGroup.ID]?.[date]}
							/>
						</td>
					{/each}
				{/each}
			</tr>
		{/each}
	</tbody>
</table>
