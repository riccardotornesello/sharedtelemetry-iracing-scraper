import { Module } from '@nestjs/common';
import { FirestoreModule } from 'src/firestore/firestore.module';
import { CompetitionDocument } from './documents/competition.document';
import { CompetitionsController } from './competitions.controller';
import { CompetitionsService } from './competitions.service';
import { IRacingSessionDocument } from './documents/iracing_session.document';

@Module({
  imports: [
    FirestoreModule.forFeature({
      collections: [
        CompetitionDocument.collectionName,
        IRacingSessionDocument.collectionName,
      ],
    }),
  ],
  controllers: [CompetitionsController],
  providers: [CompetitionsService],
})
export class CompetitionsModule {}
