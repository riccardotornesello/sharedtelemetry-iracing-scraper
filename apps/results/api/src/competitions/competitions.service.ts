import { Injectable, Inject } from '@nestjs/common';
import { CollectionReference } from '@google-cloud/firestore';
import { CompetitionDocument } from './documents/competition.document';

@Injectable()
export class CompetitionsService {
  constructor(
    @Inject(CompetitionDocument.collectionName)
    private competitionsCollection: CollectionReference<CompetitionDocument>,
  ) {}

  async findAll(): Promise<CompetitionDocument[]> {
    const snapshot = await this.competitionsCollection.get();
    const competitions: CompetitionDocument[] = [];
    snapshot.forEach((doc) => competitions.push(doc.data()));
    return competitions;
  }
}
