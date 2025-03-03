export type Crew = {
	// API info
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

	// Calculated info
	ranking: any[]; // TODO
	sum: number;
};
