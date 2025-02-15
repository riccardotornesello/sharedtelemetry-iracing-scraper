import { API_BASE_URL } from './common';

export type CompetitionRankingResponse = {
	classes: CompetitionRankingResponseClass[];
	ranking: CompetitionRankingResponseDriverRank[];
	drivers: Record<number, CompetitionRankingResponseDriver>;
	eventGroups: CompetitionRankingEventGroup[];
	competition: {
		id: number;
		name: string;
		crewDriversCount: number;
	};
};

export type CompetitionRankingResponseDriverRank = {
	pos: number;
	custId: number;
	sum: number;
	results: Record<number, Record<string, number>> | null;
};

export type CompetitionRankingResponseDriver = {
	custId: number;
	firstName: string;
	lastName: string;
	crew: {
		id: number;
		name: string;
		carId: number;
		team: {
			id: number;
			name: string;
			picture: string;
		};
		classId: number;
		carModel: string;
		carBrandIcon: string;
	};
};

export type CompetitionRankingResponseClass = {
	id: number;
	name: string;
	color: string;
	index: number;
};

export type CompetitionRankingEventGroup = {
	id: number;
	name: string;
	trackId: number;
	dates: string[];
};

export async function getCompetitionRanking(
	competitionSlug: string
): Promise<CompetitionRankingResponse | null> {
	const res = await fetch(`${API_BASE_URL}/competitions/${competitionSlug}/ranking`);

	// If the response status is 404, return null
	if (res.status === 404) {
		return null;
	}

	// If the response status is not ok, throw an error
	if (!res.ok) {
		throw new Error('Failed to fetch competition ranking');
	}

	const data: CompetitionRankingResponse = await res.json();
	return data;
}
