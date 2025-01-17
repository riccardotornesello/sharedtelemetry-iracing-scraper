import type { PageServerLoad } from './$types';
import dayjs from 'dayjs';
import type { DriverResult } from './types';

type DriverNameRow = {
	cust_id: string;
	name: string;
};

type DriverNamesMap = Record<number, string>;

type LapsRow = {
	cust_id: string;
	launch_date: string;
	best_lap: string;
};

type DriverLaps = Record<string, Record<string, number>>;

export const load: PageServerLoad = async ({ locals }) => {
	const dates = {
		'Algarve International Circuit': ['2024-09-11', '2024-09-14'],
		'Autodromo Nazionale Monza': ['2024-09-12', '2024-09-15']
	};
	const allDates = Object.values(dates).flat();

	const bestPerTrack: Record<string, number> = {};

	const query = `SELECT cust_id,
			launch_date,
			MIN(avg) AS best_lap
		FROM (
					SELECT laps.cust_id,
							(sessions.launch_at AT TIME ZONE 'CET')::date AS launch_date,
							laps.subsession_id,
							SUM(laps.lap_time) / 3 / 10000 AS avg
					FROM laps
							LEFT JOIN session_simsessions ON session_simsessions.subsession_id = laps.subsession_id
							AND laps.simsession_number = session_simsessions.simsession_number
							LEFT JOIN sessions ON sessions.subsession_id = session_simsessions.subsession_id
					WHERE session_simsessions.simsession_name = 'QUALIFY'
							AND laps.incident = false
							AND laps.lap_time > 0
							AND laps.lap_number > 0
							AND NOT laps.lap_events && array ['off track', 'pitted', 'invalid']
							AND (sessions.launch_at AT TIME ZONE 'CET')::date = ANY($1)
					GROUP BY laps.cust_id,
							(sessions.launch_at AT TIME ZONE 'CET')::date,
							laps.subsession_id
					HAVING COUNT(*) = 3
			) AS daily_laps
		GROUP BY cust_id,
			launch_date
		ORDER BY launch_date ASC,
			best_lap ASC;`;

	const res = await locals.dbConnEvents.query(query, [allDates]);

	const drivers: DriverLaps = (res.rows as LapsRow[]).reduce((acc, row) => {
		acc[row.cust_id] = acc[row.cust_id] || {};
		acc[row.cust_id][dayjs(row.launch_date).format('YYYY-MM-DD')] = parseFloat(row.best_lap);
		return acc;
	}, {} as DriverLaps);

	const driverIds = Object.keys(drivers);

	const driverNames = await locals.dbConnDrivers.query(
		`SELECT cust_id, name FROM drivers WHERE cust_id = ANY($1);`,
		[driverIds]
	);

	const driverNamesMap: DriverNamesMap = (driverNames.rows as DriverNameRow[]).reduce(
		(acc, row) => {
			acc[parseInt(row.cust_id)] = row.name;
			return acc;
		},
		{} as DriverNamesMap
	);

	const results: DriverResult[] = [];
	for (const [custId, driverResults] of Object.entries(drivers) as any) {
		const result: DriverResult = {
			custId,
			name: driverNamesMap[custId],
			sum: 0,
			isValid: true,
			laps: {}
		};

		for (const [track, trackDates] of Object.entries(dates)) {
			let minLap = 0;
			let minLapDate = '';
			let atLeastOneLap = false;

			for (const date of trackDates) {
				if (driverResults[date]) {
					if (bestPerTrack[track] === undefined || driverResults[date] < bestPerTrack[track]) {
						bestPerTrack[track] = driverResults[date];
					}
					if (minLap === 0 || driverResults[date] < minLap) {
						minLap = driverResults[date];
						minLapDate = date;
					}
					result.laps[date] = { time: driverResults[date] };
					atLeastOneLap = true;
				}
			}

			if (atLeastOneLap) {
				result.sum += minLap;
				result.laps[minLapDate].isBest = true;
			} else {
				result.isValid = false;
			}
		}

		results.push(result);
	}

	// Sort results by sum, first the valid ones and then the invalid
	const sortedResults = results.sort((a, b) => {
		if (a.isValid && !b.isValid) return -1;
		if (!a.isValid && b.isValid) return 1;
		return a.sum - b.sum;
	});

	return {
		bestPerTrack,
		results: sortedResults,
		dates
	};
};
