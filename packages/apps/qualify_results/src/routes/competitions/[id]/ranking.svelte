<script lang="ts">
	import type {
		CompetitionRankingResponseDriverRank,
		CompetitionRankingResponseDriver,
		CompetitionRankingEventGroup,
		CompetitionRankingResponseClass
	} from '$lib/api/rank';
	import RankingRow from './ranking-row.svelte';

	let {
		ranking,
		drivers,
		eventGroups,
		overallBest,
		classes,
		showTeam = true,
		showClass = true,
		showCar = true,
		showCrew = false
	}: {
		ranking: CompetitionRankingResponseDriverRank[];
		drivers: Record<number, CompetitionRankingResponseDriver>;
		eventGroups: CompetitionRankingEventGroup[];
		overallBest: Record<number, number>;
		classes: Record<number, CompetitionRankingResponseClass>;
		showTeam?: boolean;
		showClass?: boolean;
		showCar?: boolean;
		showCrew?: boolean;
	} = $props();
</script>

{#each ranking as rank}
	<RankingRow
		driverRank={rank}
		{drivers}
		{eventGroups}
		{overallBest}
		{classes}
		{showTeam}
		{showClass}
		{showCar}
		{showCrew}
	/>
{/each}
