import { Controller, Get, NotFoundException } from '@nestjs/common';
import { CompetitionsService } from './competitions.service';

type RankingItem = {
  custId: number;
  sum: number;
  isValid: boolean;
  results: Record<string, number>[];
  position: number;
};

@Controller('competitions')
export class CompetitionsController {
  constructor(private readonly competitionsService: CompetitionsService) {}

  @Get()
  async getHello() {
    const competition =
      await this.competitionsService.getCompetitionBySlug('test');
    if (!competition) {
      throw new NotFoundException('Competition not found');
    }

    const bestResults =
      await this.competitionsService.getCompetitionBestResults(competition);

    const driverIds = competition.teams.map((team) =>
      team.crews.map((crew) => crew.drivers.map((driver) => driver.iRacingId)),
    );

    const driversRanking: RankingItem[] = [];

    for (const custId of driverIds.flat(3)) {
      let totalMs = 0;
      let isValid = true;

      for (
        let eventGroupIndex = 0;
        eventGroupIndex < competition.eventGroups.length;
        eventGroupIndex++
      ) {
        const eventGroupResults: Record<string, number> | undefined =
          bestResults[custId.toString()]?.[eventGroupIndex];
        if (!eventGroupResults) {
          isValid = false;
        } else {
          const bestGroupResult = Math.min(
            ...Object.values(eventGroupResults || {}),
          );
          totalMs += bestGroupResult;
        }
      }

      driversRanking.push({
        position: 0,
        custId: custId,
        sum: totalMs,
        isValid,
        results: bestResults[custId.toString()],
      });
    }

    // Sort the driversRanking. Primary sort by isValid, secondary by sum
    driversRanking.sort((a, b) => {
      if (a.isValid && !b.isValid) {
        return -1;
      }
      if (!a.isValid && b.isValid) {
        return 1;
      }
      if (a.sum === 0) {
        return 1;
      }
      if (b.sum === 0) {
        return -1;
      }
      return a.sum - b.sum;
    });

    // Add the position to the driversRanking
    driversRanking.forEach((item, index) => {
      item.position = index + 1;
    });

    // TODO: competition serializer
    return { driversRanking, competition };
  }
}
