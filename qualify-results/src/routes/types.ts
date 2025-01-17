export interface ResultsLap {
	time: number;
	isBest?: boolean;
}

export interface DriverResult {
	custId: number;
	name: string;
	sum: number;
	isValid: boolean;
	laps: Record<string, ResultsLap>;
}
