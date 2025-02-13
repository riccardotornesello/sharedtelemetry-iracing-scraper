import { sequence } from '@sveltejs/kit/hooks';
import * as Sentry from '@sentry/sveltekit';
import type { Handle } from '@sveltejs/kit';
import { i18n } from '$lib/i18n';

Sentry.init({
	dsn: 'https://afb7dcacd7a773a95b447f879945d32f@o4508811750408192.ingest.de.sentry.io/4508811753947216',
	tracesSampleRate: 1,
	release: 'demo'
});

const handleParaglide: Handle = i18n.handle();

export const handle = sequence(Sentry.sentryHandle(), handleParaglide);
export const handleError = Sentry.handleErrorWithSentry();
