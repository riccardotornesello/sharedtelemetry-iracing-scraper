import { Timestamp } from '@google-cloud/firestore';

export class IRacingSessionDocument {
  static collectionName = 'iracing_sessions';

  parsed: boolean;

  leagueId: number;
  seasonId: number;
  launchAt: Timestamp;

  simsessions: SimSession[];
}

class SimSession {
  simsessionNumber: number;
  simsessionType: number;
  simsessionName: string;

  participants: Participant[];
}

class Participant {
  custId: number;
  carId: number;

  laps: Lap[];
}

export class Lap {
  lapEvents: string[];
  incident: boolean;
  lapTime: number;
  lapNumber: number;
}
