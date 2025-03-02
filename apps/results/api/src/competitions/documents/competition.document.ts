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
  eventGroups: EventGroup[];
}

class CompetitionClass {
  name: string;
  color: string;
}

class CompetitionTeam {
  name: string;
  pictureUrl?: string | null;

  crews: CompetitionCrew[];
}

class CompetitionCrew {
  name: string;
  iRacingCarId: number;

  class?: string | null;

  drivers: CompetitionDriver[];
}

class CompetitionDriver {
  firstName: string;
  lastName: string;

  iRacingId: number;
}

class EventGroup {
  name: string;
  iRacingTrackId: number;
  sessions: EventSession[];
}

class EventSession {
  fromTime: Timestamp;
  toTime: Timestamp;
}
