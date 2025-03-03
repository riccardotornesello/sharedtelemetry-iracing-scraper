import { Injectable, Inject } from '@nestjs/common';
import { CollectionReference } from '@google-cloud/firestore';
import { CompetitionDocument } from './documents/competition.document';
import {
  IRacingSessionDocument,
  Lap,
} from './documents/iracing_session.document';
import { Timestamp } from '@google-cloud/firestore';
import * as dayjs from 'dayjs';
import { plainToInstance } from 'class-transformer';

@Injectable()
export class CompetitionsService {
  constructor(
    @Inject(CompetitionDocument.collectionName)
    private competitionsCollection: CollectionReference<CompetitionDocument>,

    @Inject(IRacingSessionDocument.collectionName)
    private iRacingSessionsCollection: CollectionReference<IRacingSessionDocument>,
  ) {}

  async getCompetition(id: string): Promise<CompetitionDocument | null> {
    const snapshot = await this.competitionsCollection.doc(id).get();
    return snapshot.data() || null;
  }

  async getCompetitionBySlug(
    slug: string,
  ): Promise<CompetitionDocument | null> {
    const snapshot = await this.competitionsCollection
      .where('slug', '==', slug)
      .limit(1)
      .get();

    // TODO: move to utils file
    const transformFirestoreTypes = (obj: any): any => {
      Object.keys(obj).forEach((key) => {
        if (!obj[key]) return;
        if (typeof obj[key] === 'object' && 'toDate' in obj[key]) {
          obj[key] = obj[key].toDate().toISOString();
        } else if (obj[key].constructor.name === 'GeoPoint') {
          const { latitude, longitude } = obj[key];
          obj[key] = { latitude, longitude };
        } else if (obj[key].constructor.name === 'DocumentReference') {
          const { id, path } = obj[key];
          obj[key] = { id, path };
        } else if (typeof obj[key] === 'object') {
          transformFirestoreTypes(obj[key]);
        }
      });
      return obj;
    };

    const data = snapshot.docs[0]?.data();
    return data
      ? plainToInstance(
          CompetitionDocument,
          transformFirestoreTypes({
            ...data,
          }),
          {
            ignoreDecorators: true,
          },
        )
      : null;
  }

  async getCompetitionBestResults(competition: CompetitionDocument) {
    const bestResults = {}; //  Customer ID, Group, Date, average ms

    for (
      let eventGroupIndex = 0;
      eventGroupIndex < competition.eventGroups.length;
      eventGroupIndex++
    ) {
      const eventGroup = competition.eventGroups[eventGroupIndex];

      for (const eventSession of eventGroup.sessions) {
        const sessionResults = await this.getGroupSessions({
          leagueId: competition.leagueId,
          seasonId: competition.seasonId,
          trackId: eventGroup.iRacingTrackId,
          fromTime: eventSession.fromTime,
          toTime: eventSession.toTime,
        });

        for (const custId in sessionResults) {
          if (!bestResults[custId]) {
            bestResults[custId] = {};
          }

          if (!bestResults[custId][eventGroupIndex]) {
            bestResults[custId][eventGroupIndex] = {};
          }

          const dateString = dayjs(eventSession.fromTime.toMillis()).format(
            'YYYY-MM-DD',
          );

          bestResults[custId][eventGroupIndex][dateString] =
            sessionResults[custId];
        }
      }
    }

    return bestResults;
  }

  async getGroupSessions({
    leagueId,
    seasonId,
    trackId,
    fromTime,
    toTime,
  }: {
    leagueId: number;
    seasonId: number;
    trackId: number;
    fromTime: Timestamp;
    toTime: Timestamp;
  }) {
    const sessions = await this.iRacingSessionsCollection
      .where('launchAt', '>=', fromTime)
      .where('launchAt', '<', toTime)
      .where('trackId', '==', trackId)
      .where('leagueId', '==', leagueId)
      .where('seasonId', '==', seasonId)
      .get();

    const groupDriverResults = {};

    sessions.forEach((s) => {
      const session = s.data();

      for (const simsession of session.simsessions) {
        // TODO: variable allowed simsession types
        if (simsession.simsessionName !== 'QUALIFY') {
          continue;
        }

        for (const participant of simsession.participants) {
          // TODO: check if the car is allowed

          const avgLapTime = this.extractAverageLapTime(participant.laps);
          if (
            avgLapTime &&
            (!groupDriverResults[participant.custId] ||
              groupDriverResults[participant.custId] > avgLapTime)
          ) {
            groupDriverResults[participant.custId] = avgLapTime;
          }
        }
      }
    });

    return groupDriverResults;
  }

  extractAverageLapTime(laps: Lap[]): number | null {
    let validLaps = 0;
    let totalLapTime = 0;

    for (const lap of laps) {
      if (this.isLapPitted(lap)) {
        // If the driver already started a stint, end it
        if (validLaps > 0) {
          return null;
        } else {
          continue;
        }
      }

      if (!this.isLapValid(lap)) {
        return null;
      }

      validLaps++;
      totalLapTime += lap.lapTime;

      // TODO: variable stint length
      const stintLength = 3;
      if (validLaps == stintLength) {
        return totalLapTime / stintLength / 10;
      }
    }

    return null;
  }

  isLapPitted(lap: Lap): boolean {
    return lap.lapEvents.includes('pitted');
  }

  isLapValid(lap: Lap): boolean {
    if (lap.lapNumber <= 0) {
      return false;
    }

    if (lap.lapTime <= 0) {
      return false;
    }

    if (lap.incident) {
      return false;
    }

    if (lap.lapEvents.length > 0) {
      return false;
    }

    return true;
  }
}
