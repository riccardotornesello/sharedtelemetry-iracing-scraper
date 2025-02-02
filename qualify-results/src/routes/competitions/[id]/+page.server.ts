import type { PageServerLoad } from './$types';

const BASE_URL = 'https://api.results.sharedtelemetry.com';

export const load: PageServerLoad = async ({ params }) => {
	const { id } = params;

	const res = await fetch(`${BASE_URL}/competitions/${id}/ranking`);
	const { ranking, drivers, eventGroups } = await res.json();

	return {
		ranking,
		drivers,
		eventGroups
	};
};
