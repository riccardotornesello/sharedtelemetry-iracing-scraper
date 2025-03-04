import { Timestamp } from '@google-cloud/firestore';
import { Exclude, Transform, Type } from 'class-transformer';
import { TimestampTransform } from 'src/utils/decorators/transform';

export class CompetitionDocument {
  static collectionName = 'results_competitions';

  leagueId: number;
  seasonId: number;

  name: string;
  slug: string;

  crewDriversCount: number;

  @Exclude()
  @TimestampTransform()
  createdAt: Timestamp;

  @Exclude()
  @TimestampTransform()
  updatedAt: Timestamp;

  @Type(() => CompetitionClass)
  classes: Record<string, CompetitionClass[]>;

  @Type(() => CompetitionTeam)
  teams: CompetitionTeam[];

  @Type(() => EventGroup)
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

  @Type(() => EventSession)
  sessions: EventSession[];
}

class EventSession {
  @TimestampTransform()
  fromTime: Timestamp;

  @TimestampTransform()
  toTime: Timestamp;
}
