<script lang="ts">
	import type {
		CompetitionRankingResponseDriverRank,
		CompetitionRankingResponseDriver,
		CompetitionRankingEventGroup,
		CompetitionRankingResponseClass
	} from '$lib/api/rank';
	import TimeCard from '../../../components/time-card.svelte';

	let {
		driverRank,
		drivers,
		eventGroups,
		overallBest,
		classes,
		showTeam = true,
		showClass = true,
		showCar = true,
		showCrew = false
	}: {
		driverRank: CompetitionRankingResponseDriverRank;
		drivers: Record<number, CompetitionRankingResponseDriver>;
		eventGroups: CompetitionRankingEventGroup[];
		overallBest: Record<number, number>;
		classes: Record<number, CompetitionRankingResponseClass>;
		showTeam?: boolean;
		showClass?: boolean;
		showCar?: boolean;
		showCrew?: boolean;
	} = $props();

	let driver = drivers[driverRank.custId];

	let crew = driver?.crew;
	let crewName = crew?.name;

	let team = crew?.team;
	let teamPicture = team?.picture;
	let teamName = team?.name;

	let carClass = classes[crew?.classId];
	let carModel = crew?.carModel;
	let carBrandIcon = crew?.carBrandIcon;
</script>

<tr class="border-b border-gray-700 bg-gray-800">
	<!-- TODO: position in class -->
	<td class="px-6 py-4 text-center">P{driverRank.pos}</td>

	{#if showClass}
		<td
			class="{carClass && 'min-w-10 px-4 py-2'} text-center text-white"
			style:background-color={carClass?.color}
		>
			{carClass?.name}
		</td>
	{:else}
		<td></td>
	{/if}

	<td class="px-6 py-4">
		<div class="flex items-center space-x-2">
			{#if carBrandIcon && showCar}
				<img src={carBrandIcon} alt={carModel} class="mr-4 h-6 w-6 object-cover" />
			{/if}

			{driver ? `${driver.firstName} ${driver.lastName}` : driverRank.custId}
		</div>
	</td>

	<td class="px-6 py-4">
		{#if showTeam}
			<div class="flex items-center space-x-2">
				{#if teamPicture}
					<img src={teamPicture} alt={teamName} class="mr-4 h-6 w-6 object-cover" />
				{/if}

				<div class="flex flex-col">
					<span>
						{teamName}
					</span>

					{#if showCrew && crewName}
						<span>
							{crewName}
						</span>
					{/if}
				</div>
			</div>
		{/if}
	</td>

	<td class="px-2 py-2">
		<TimeCard time={driverRank.sum} />
	</td>

	{#each eventGroups as eventGroup}
		{#each eventGroup.dates as date}
			<td class="px-2 py-2">
				<TimeCard
					isPersonalBest={driverRank.results?.[eventGroup.id]?.[date] ==
						Math.min(...Object.values(driverRank.results?.[eventGroup.id] || {}))}
					isOverallBest={driverRank.results?.[eventGroup.id]?.[date] == overallBest[eventGroup.id]}
					time={driverRank.results?.[eventGroup.id]?.[date]}
				/>
			</td>
		{/each}
	{/each}
</tr>
