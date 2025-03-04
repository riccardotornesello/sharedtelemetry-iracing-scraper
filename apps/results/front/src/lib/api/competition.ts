import { API_BASE_URL } from './common';

export interface RankingItem {
	custId: number;
	sum: number;
	isValid: boolean;
	results: Record<string, number>[];
	position: number;
}

export interface Competition {
	leagueId: number;
	seasonId: number;
	name: string;
	slug: string;
	crewDriversCount: number;
	classes: Record<string, CompetitionClass[]>;
	teams: CompetitionTeam[];
	eventGroups: EventGroup[];
}

export interface CompetitionClass {
	name: string;
	color: string;
}

export interface CompetitionTeam {
	name: string;
	pictureUrl?: string | null;

	crews: CompetitionCrew[];
}

export interface CompetitionCrew {
	name: string;
	iRacingCarId: number;

	class?: string | null;

	drivers: CompetitionDriver[];
}

export interface CompetitionDriver {
	firstName: string;
	lastName: string;

	iRacingId: number;
}

export interface EventGroup {
	name: string;
	iRacingTrackId: number;
	sessions: EventSession[];
}

export interface EventSession {
	fromTime: string;
	toTime: string;
}

export type CompetitionRankingResponse = {
	driversRanking: RankingItem[];
	competition: Competition;
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
