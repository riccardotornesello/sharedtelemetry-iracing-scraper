<script lang="ts">
	import TimeCard from '../../../components/time-card.svelte';
	import type {
		CompetitionDriver,
		EventGroup,
		CompetitionCrew,
		CompetitionTeam,
		RankingItem
	} from '$lib/api/competition';
	import TimeCardGroup from './time-card-group.svelte';

	let {
		rankingItem,
		driver,
		events,
		crew,
		team,
		overallBest
	}: {
		rankingItem: RankingItem;
		driver: CompetitionDriver;
		events: EventGroup[];
		crew: CompetitionCrew;
		team: CompetitionTeam;
		overallBest: Record<string, number>;
	} = $props();
</script>

<tr class="border-b border-gray-700 bg-gray-800">
	<td class="px-6 py-4 text-center">P{rankingItem.position}</td>

	<td class="px-6 py-4">
		<div class="flex items-center space-x-2">
			{driver.firstName}
			{driver.lastName}
		</div>
	</td>

	<td class="px-6 py-4">
		<div class="flex items-center space-x-2">
			<div class="flex flex-row">
				{#if team.pictureUrl}
					<img src={team.pictureUrl} alt={team.name} class="mr-4 h-6 w-6 object-cover" />
				{/if}
				<span>
					{team.name}
				</span>
			</div>
		</div>
	</td>

	<td class="px-2 py-2">
		<TimeCard time={rankingItem.sum} />
	</td>

	{#each events as event, eventIndex}
		<TimeCardGroup
			{event}
			eventResults={rankingItem.results[eventIndex]}
			eventOverallBest={overallBest[eventIndex]}
		/>
	{/each}
</tr>
