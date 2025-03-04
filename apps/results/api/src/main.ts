import { NestFactory, Reflector } from '@nestjs/core';
import { AppModule } from './app.module';
import { ClassSerializerInterceptor } from '@nestjs/common';
import { http } from '@google-cloud/functions-framework';
import { Logger } from '@nestjs/common';
import * as express from 'express';
import { ExpressAdapter } from '@nestjs/platform-express';

const server = express();

async function bootstrap(expressInstance) {
  const app = await NestFactory.create(
    AppModule,
    new ExpressAdapter(expressInstance),
  );

  // Serialization
  app.useGlobalInterceptors(new ClassSerializerInterceptor(app.get(Reflector)));

  return app.init();
}

bootstrap(server).then(async (app) => {
  // Start the server
  if (process.env.ENVIRONMENT === 'production') {
    Logger.log('🚀 Starting production server...');
  } else {
    const port = process.env.PORT || 3000;

    Logger.log(`🚀 Starting development server on http://localhost:${port}`);

    await app.listen(port);
  }
});

http('apiNEST', server);
