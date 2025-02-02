import type { PageServerLoad } from './$types';

const BASE_URL = 'http://localhost:8080';

export const load: PageServerLoad = async ({ params }) => {
	const { id } = params;

	const res = await fetch(`${BASE_URL}/competitions/${id}/ranking`);
	const { ranking, drivers, eventGroups, competition } = await res.json();

	return {
		ranking,
		drivers,
		eventGroups,
		competition
	};
};
