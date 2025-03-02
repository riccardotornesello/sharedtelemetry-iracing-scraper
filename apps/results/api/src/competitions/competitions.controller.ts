import { Controller, Get, NotFoundException } from '@nestjs/common';
import { CompetitionsService } from './competitions.service';
import * as dayjs from 'dayjs';

@Controller('competitions')
export class CompetitionsController {
  constructor(private readonly competitionsService: CompetitionsService) {}

  @Get()
  async getHello() {
    const competition = await this.competitionsService.getCompetition(
      'ywbHYk74cFdDz1ZmIPay',
    );
    if (!competition) {
      throw new NotFoundException('Competition not found');
    }

    const bestResults =
      await this.competitionsService.getCompetitionBestResults(competition);

    const ranking: any[] = [];

    // TODO: loop on drivers
    for (const custId in bestResults) {
      let totalMs = 0;
      let isValid = true;

      for (
        let eventGroupIndex = 0;
        eventGroupIndex < competition.eventGroups.length;
        eventGroupIndex++
      ) {
        const eventGroupResults: Record<string, number> | undefined =
          bestResults[custId][eventGroupIndex];
        if (!eventGroupResults) {
          isValid = false;
        } else {
          const bestGroupResult = Math.min(
            ...Object.values(eventGroupResults || {}),
          );
          totalMs += bestGroupResult;
        }
      }

      ranking.push({
        custId: parseInt(custId),
        sum: totalMs,
        isValid,
        results: bestResults[custId],
      });
    }

    return ranking;
  }
}
