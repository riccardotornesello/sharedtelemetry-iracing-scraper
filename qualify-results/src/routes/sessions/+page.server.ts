import type { PageServerLoad } from './$types';
import dayjs from 'dayjs';

export const load: PageServerLoad = async ({ locals }) => {
	const dates = {
		'track 1': ['2024-09-11', '2024-09-14'],
		'track 2': ['2024-09-12', '2024-09-15']
	};
	const allDates = Object.values(dates).flat();

	const query = `SELECT ss.simsession_number, ss.subsession_id, s.launch_at, s.track_id
		FROM session_simsessions ss
		LEFT JOIN sessions s ON ss.subsession_id = s.subsession_id
		WHERE s.launch_at::date = ANY($1)
		ORDER BY s.launch_at ASC;`;

	const res = await locals.dbConn.query(query, [allDates]);

	return {
		sessions: res.rows
	};
};
