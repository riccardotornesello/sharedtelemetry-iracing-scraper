<script lang="ts">
	import type {
		CompetitionRankingResponseDriver,
		CompetitionRankingEventGroup,
		CompetitionRankingResponseClass
	} from '$lib/api/competition';
	import type { Crew } from './types';
	import TimeCard from '../../../components/time-card.svelte';
	import Ranking from './ranking.svelte';

	let {
		position,
		crew,
		drivers,
		eventGroups,
		overallBest,
		classes
	}: {
		position: number;
		crew: Crew;
		drivers: Record<number, CompetitionRankingResponseDriver>;
		eventGroups: CompetitionRankingEventGroup[];
		overallBest: Record<number, number>;
		classes: Record<number, CompetitionRankingResponseClass>;
	} = $props();

	let isOpen = $state(true);

	let carClass = classes[crew.classId];
	let carModel = crew.carModel;
	let carBrandIcon = crew.carBrandIcon;

	let team = crew.team;
	let teamPicture = team.picture;
	let teamName = team.name;
</script>

<tr class="border-b border-gray-700 bg-gray-900">
	<td class="px-6 py-4 text-center">
		<button onclick={() => (isOpen = !isOpen)}>
			P{position}
			{isOpen ? '▼' : '►'}
		</button>
	</td>

	<td
		class="{carClass && 'min-w-10 px-4 py-2'} text-center text-white"
		style:background-color={carClass?.color}
	>
		{carClass?.name}
	</td>

	<td class="px-6 py-4">
		<div class="flex items-center space-x-2">
			{#if carBrandIcon}
				<img src={carBrandIcon} alt={carModel} class="mr-4 h-6 w-6 object-cover" />
			{/if}

			{crew.name}
		</div>
	</td>

	<td class="px-6 py-4">
		<div class="flex items-center space-x-2">
			{#if teamPicture}
				<img src={teamPicture} alt={teamName} class="mr-4 h-6 w-6 object-cover" />
			{/if}

			{crew.team.name}
		</div>
	</td>

	<td class="px-2 py-2"><TimeCard time={crew.sum} /></td>

	{#each eventGroups as eventGroup}
		{#each eventGroup.dates as date}
			<td></td>
		{/each}
	{/each}
</tr>

{#if isOpen}
	<Ranking
		ranking={crew.ranking}
		{drivers}
		{eventGroups}
		{overallBest}
		{classes}
		showTeam={false}
		showClass={false}
		showCar={false}
	/>
{/if}
