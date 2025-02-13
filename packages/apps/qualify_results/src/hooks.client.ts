import { handleErrorWithSentry, replayIntegration } from "@sentry/sveltekit";
import * as Sentry from '@sentry/sveltekit';

Sentry.init({
  dsn: 'https://afb7dcacd7a773a95b447f879945d32f@o4508811750408192.ingest.de.sentry.io/4508811753947216',
  tracesSampleRate: 1.0,
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1.0,
  integrations: [replayIntegration()],
});

export const handleError = handleErrorWithSentry();
