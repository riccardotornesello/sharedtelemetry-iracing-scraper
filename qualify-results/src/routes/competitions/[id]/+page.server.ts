import type { PageServerLoad } from './$types';

const BASE_URL = 'http://localhost:8080'; // TODO: get from env

const CREW_DRIVERS_COUNT = 2; // TODO: get from competition info

export const load: PageServerLoad = async ({ params }) => {
	const { id } = params;

	const res = await fetch(`${BASE_URL}/competitions/${id}/ranking`);
	const { ranking, drivers, eventGroups, competition } = await res.json();

	const crews = {};
	Object.values(ranking).forEach((rank) => {
		const driver = drivers[rank.custId];

		if (!crews[driver.crew.id]) {
			crews[driver.crew.id] = {
				id: driver.crew.id,
				name: driver.crew.name,
				team: driver.crew.team,
				ranking: []
			};
		}

		crews[driver.crew.id].ranking.push(rank);
	});

	for (const crewId of Object.keys(crews)) {
		const driversWithTime = crews[crewId].ranking.filter((rank) => rank.sum);
		if (driversWithTime.length >= CREW_DRIVERS_COUNT) {
			let sum = 0;
			for (let i = 0; i < CREW_DRIVERS_COUNT; i++) {
				sum += driversWithTime[i].sum;
			}
			crews[crewId].sum = sum;
		} else {
			crews[crewId].sum = 0;
		}
	}

	const sortedCrews = Object.values(crews).sort((a, b) => {
		if (a.sum === 0) {
			return 1;
		} else if (b.sum === 0) {
			return -1;
		} else {
			return a.sum - b.sum;
		}
	});

	return {
		ranking,
		drivers,
		eventGroups,
		competition,
		crews: sortedCrews
	};
};
