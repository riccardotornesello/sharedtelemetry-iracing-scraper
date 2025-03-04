<script lang="ts">
	import type {
		CompetitionDriver,
		EventGroup,
		CompetitionCrew,
		CompetitionTeam,
		RankingItem
	} from '$lib/api/competition';
	import RankingRow from './ranking-row.svelte';

	let {
		ranking,
		drivers,
		events,
		driversCrewMap,
		driversTeamMap,
		overallBest
	}: {
		ranking: RankingItem[];
		drivers: Record<number, CompetitionDriver>;
		events: EventGroup[];
		driversCrewMap: Record<number, CompetitionCrew>;
		driversTeamMap: Record<number, CompetitionTeam>;
		overallBest: Record<string, number>;
	} = $props();
</script>

{#each ranking as rankingItem}
	<RankingRow
		{rankingItem}
		driver={drivers[rankingItem.custId]}
		{events}
		crew={driversCrewMap[rankingItem.custId]}
		team={driversTeamMap[rankingItem.custId]}
		{overallBest}
	/>
{/each}
