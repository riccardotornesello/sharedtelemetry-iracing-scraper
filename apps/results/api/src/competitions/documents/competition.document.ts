import { Timestamp } from '@google-cloud/firestore';

export class CompetitionDocument {
  static collectionName = 'results_competitions';

  leagueId: number;
  seasonId: number;

  name: string;
  slug: string;
  crewDriversCount: number;

  createdAt: Timestamp;
  updatedAt: Timestamp;

  classes: Record<string, CompetitionClass[]>;
  teams: CompetitionTeam[];
}

class CompetitionClass {
  name: string;
  color: string;
}

class CompetitionTeam {
  name: string;
  pictureUrl: string;

  crews: CompetitionCrew[];
}

class CompetitionCrew {
  name: string;
  iRacingCarId: number;

  class: string;

  drivers: CompetitionDriver[];
}

class CompetitionDriver {
  firstName: string;
  lastName: string;

  iRacingId: number;
}
