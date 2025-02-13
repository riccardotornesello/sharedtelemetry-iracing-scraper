<script lang="ts">
	import CellWithPicture from '../../../components/cell-with-picture.svelte';
	import TimeCard from '../../../components/time-card.svelte';
	import type { DriverRanking } from './types';

	let {
		ranking,
		drivers,
		eventGroups,
		overallBest
	}: {
		ranking: DriverRanking[];
		drivers: any;
		eventGroups: any;
		overallBest: any;
	} = $props();
</script>

{#each ranking as rank}
	<tr class="border-b border-gray-700 bg-gray-800">
		<td class="px-6 py-4 text-center">P{rank.pos}</td>
		<td class="px-6 py-4">
			<CellWithPicture picture={drivers[rank.custId]?.crew.carBrandPicture}>
				{drivers[rank.custId]?.name || rank.custId}
			</CellWithPicture>
		</td>
		<td class="px-6 py-4">
			<CellWithPicture picture={drivers[rank.custId]?.crew.team.picture}>
				{drivers[rank.custId]?.crew.team.name}
			</CellWithPicture>
		</td>
		<td class="px-2 py-2">
			<TimeCard time={rank.sum} />
		</td>
		{#each eventGroups as eventGroup}
			{#each eventGroup.dates as date}
				<td class="px-2 py-2">
					<TimeCard
						isPersonalBest={rank.results?.[eventGroup.id]?.[date] ==
							Math.min(...Object.values(rank.results?.[eventGroup.id] || {}))}
						isOverallBest={rank.results?.[eventGroup.id]?.[date] == overallBest[eventGroup.id]}
						time={rank.results?.[eventGroup.id]?.[date]}
					/>
				</td>
			{/each}
		{/each}
	</tr>
{/each}
