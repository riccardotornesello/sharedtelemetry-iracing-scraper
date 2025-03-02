import { Module, DynamicModule, Provider } from '@nestjs/common';
import { Firestore, Settings } from '@google-cloud/firestore';
import {
  FirestoreDatabaseProvider,
  FirestoreOptionsProvider,
} from './firestore.providers';

type FirestoreModuleOptions = {
  imports?: any[];
  useFactory: (...args: any[]) => Settings;
  inject?: any[];
};

type FirestoreFeatureOptions = {
  collections: string[];
};

@Module({})
export class FirestoreModule {
  static forRoot(options: FirestoreModuleOptions): DynamicModule {
    const optionsProvider = {
      provide: FirestoreOptionsProvider,
      useFactory: options.useFactory,
      inject: options.inject,
    };

    const dbProvider = {
      provide: FirestoreDatabaseProvider,
      useFactory: (config) => new Firestore(config),
      inject: [FirestoreOptionsProvider],
    };

    return {
      global: true,
      module: FirestoreModule,
      imports: options.imports,
      providers: [optionsProvider, dbProvider],
      exports: [dbProvider],
    };
  }

  static forFeature(options: FirestoreFeatureOptions): DynamicModule {
    const collectionProviders: Provider[] = options.collections.map(
      (collectionName) => ({
        provide: collectionName,
        useFactory: (db: Firestore) => db.collection(collectionName),
        inject: [FirestoreDatabaseProvider],
      }),
    );

    return {
      module: FirestoreModule,
      providers: collectionProviders,
      exports: collectionProviders,
    };
  }
}
