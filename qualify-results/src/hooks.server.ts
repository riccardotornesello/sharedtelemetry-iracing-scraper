import type { Handle } from '@sveltejs/kit';
import { i18n } from '$lib/i18n';
import { connectToDB } from '$lib/db';
import { sequence } from '@sveltejs/kit/hooks';

const handleParaglide: Handle = i18n.handle();

export const handleDb: Handle = async (params) => {
	const [dbConnEvents, dbConnDrivers] = await connectToDB();
	params.event.locals.dbConnEvents = dbConnEvents;
	params.event.locals.dbConnDrivers = dbConnDrivers;

	const response = await params.resolve(params.event);

	dbConnEvents.release();
	dbConnDrivers.release();

	return response;
};

export const handle = sequence(handleParaglide, handleDb);
