import { error } from '@sveltejs/kit';
import { getCompetitionRanking, type CompetitionRankingResponseClass } from '$lib/api/rank';
import type { PageServerLoad } from './$types';
import type { Crew } from './types';

export const load: PageServerLoad = async ({ params }) => {
	const { id } = params;

	const competitionRanking = await getCompetitionRanking(id);
	if (!competitionRanking) {
		return error(404, 'Competition not found');
	}

	// Extract the crews from the competition ranking
	const crews: Record<number, Crew> = {};
	competitionRanking.ranking.forEach((rank) => {
		const driver = competitionRanking.drivers[rank.custId];
		if (!driver) {
			throw new Error(`Driver with id ${rank.custId} not found`);
		}

		// If the crew is not already in the crews object, add it
		if (!crews[driver.crew.id]) {
			crews[driver.crew.id] = {
				...driver.crew,
				ranking: [],
				sum: 0
			};
		}

		// Add the driver's rank to the crew's ranking
		crews[driver.crew.id].ranking.push(rank);
	});

	// Calculate the sum for each crew
	for (const crewId in crews) {
		// Check that the crew has enough drivers with a time
		const driversWithTime = crews[crewId].ranking.filter((rank) => rank.sum);

		if (driversWithTime.length >= competitionRanking.competition.crewDriversCount) {
			// Sum the times of the first `competitionRanking.competition.crewDriversCount` drivers
			let sum = 0;
			for (let i = 0; i < competitionRanking.competition.crewDriversCount; i++) {
				sum += driversWithTime[i].sum;
			}
			crews[crewId].sum = sum;
		} else {
			crews[crewId].sum = 0;
		}
	}

	// Sort the crews by sum
	const sortedCrews = Object.values(crews).sort((a, b) => {
		if (a.sum === 0) {
			return 1;
		} else if (b.sum === 0) {
			return -1;
		} else {
			return a.sum - b.sum;
		}
	});

	// Convert the classes array to a map
	const classesMap = competitionRanking.classes.reduce(
		(acc, cls) => {
			acc[cls.id] = cls;
			return acc;
		},
		{} as Record<number, CompetitionRankingResponseClass>
	);

	// Calculate the overall best time for each event group
	const overallBest = competitionRanking.eventGroups.reduce(
		(acc, eventGroup) => {
			acc[eventGroup.id] = Math.min(
				...competitionRanking.ranking.map((r) => Object.values(r.results?.[eventGroup.id] || {})).flat()
			);
			return acc;
		},
		{} as Record<number, number>
	);

	return {
		competitionRanking,
		crews: sortedCrews,
		classes: classesMap,
		overallBest,
	};
};
