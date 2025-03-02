import { Module } from '@nestjs/common';
import { FirestoreModule } from 'src/firestore/firestore.module';
import { CompetitionDocument } from './documents/competition.document';
import { CompetitionsController } from './competitions.controller';
import { CompetitionsService } from './competitions.service';

@Module({
  imports: [
    FirestoreModule.forFeature({
      collections: [CompetitionDocument.collectionName],
    }),
  ],
  controllers: [CompetitionsController],
  providers: [CompetitionsService],
})
export class CompetitionsModule {}
