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

<h1 class="p-5 text-center text-2xl font-bold text-gray-800 dark:text-gray-200">
	{data.competition.name}
</h1>

<table class="w-full table-auto text-left text-sm text-gray-500 rtl:text-right dark:text-gray-400">
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
		{#each data.ranking as rank, index}
			<tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800">
				<td class="px-6 py-4 text-center">P{index + 1}</td>
				<td class="px-6 py-4">{data.drivers[rank.custId]?.Name || rank.custId}</td>
				<td class="px-2 py-2">
					<TimeCard time={rank.sum} />
				</td>
				{#each data.eventGroups as eventGroup}
					{#each eventGroup.dates as date}
						<td class="px-2 py-2">
							<TimeCard
								isPersonalBest={rank.results?.[eventGroup.id]?.[date] ==
									Math.min(...Object.values(rank.results?.[eventGroup.id] || {}))}
								isOverallBest={rank.results?.[eventGroup.id]?.[date] ==
									Math.min(
										...data.ranking
											.map((r) => Object.values(r.results?.[eventGroup.id] || {}))
											.flat()
									)}
								time={rank.results?.[eventGroup.id]?.[date]}
							/>
						</td>
					{/each}
				{/each}
			</tr>
		{/each}
	</tbody>
</table>
