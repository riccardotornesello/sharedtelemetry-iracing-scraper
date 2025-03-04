import { error } from '@sveltejs/kit';
import {
	getCompetitionRanking,
	type CompetitionCrew,
	type CompetitionDriver,
	type CompetitionTeam
} from '$lib/api/competition';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
	const { id } = params;

	const competitionRankingResponse = await getCompetitionRanking(id);
	if (!competitionRankingResponse) {
		return error(404, 'Competition not found');
	}

	// Calculate the overall best time for each event group
	const overallBest = {} as Record<string, number>;
	competitionRankingResponse.driversRanking.forEach((driverRanking) => {
		Object.entries(driverRanking.results).forEach(([eventGroupId, results]) => {
			Object.values(results).forEach((result) => {
				if (result > 0 && (!overallBest[eventGroupId] || result < overallBest[eventGroupId])) {
					overallBest[eventGroupId] = result;
				}
			});
		});
	});

	// Extract the drivers
	let driversMap: Record<number, CompetitionDriver> = {};
	let crewsMap: Record<number, CompetitionCrew> = {};
	let teamsMap: Record<number, CompetitionTeam> = {};
	let driversCrewMap: Record<number, CompetitionCrew> = {};
	let driversTeamMap: Record<number, CompetitionTeam> = {};

	competitionRankingResponse.competition.teams.forEach((team, teamIndex) => {
		teamsMap[teamIndex] = team;
		team.crews.forEach((crew, crewIndex) => {
			crewsMap[crewIndex] = crew;
			crew.drivers.forEach((driver) => {
				driversMap[driver.iRacingId] = driver;
				driversCrewMap[driver.iRacingId] = crew;
				driversTeamMap[driver.iRacingId] = team;
			});
		});
	});

	return {
		driversRanking: competitionRankingResponse.driversRanking,
		competition: competitionRankingResponse.competition,
		overallBest,
		driversMap,
		crewsMap,
		teamsMap,
		driversCrewMap,
		driversTeamMap
	};
};
