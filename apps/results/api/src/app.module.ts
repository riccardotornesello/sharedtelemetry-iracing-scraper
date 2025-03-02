import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { FirestoreModule } from './firestore/firestore.module';
import { CompetitionsModule } from './competitions/competitions.module';

@Module({
  imports: [
    FirestoreModule.forRoot({
      useFactory: () => ({
        projectId: 'sharedtelemetryapp', // TODO: variable
      }),
    }),
    CompetitionsModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
