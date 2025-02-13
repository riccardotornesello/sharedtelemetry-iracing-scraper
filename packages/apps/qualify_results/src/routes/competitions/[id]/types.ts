// TODO: move types and components to lib
// TODO: check types

export type Crew = {
	id: number;
	name: string;
	team: string;
	ranking: DriverRanking[];
	sum: number;
};

export type DriverRanking = {
	pos: number;
	custId: number;
	sum: number;
	results: Record<number, Record<string, number>> | null;
};
