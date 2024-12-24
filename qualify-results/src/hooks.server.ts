import type { Handle } from '@sveltejs/kit';
import { i18n } from '$lib/i18n';
import { connectToDB } from '$lib/db';
import { sequence } from '@sveltejs/kit/hooks';

const handleParaglide: Handle = i18n.handle();

export const handleDb: Handle = async (params) => {
	const dbConn = await connectToDB();
	params.event.locals.dbConn = dbConn;
	const response = await params.resolve(params.event);
	dbConn.release();
	return response;
};

export const handle = sequence(handleParaglide, handleDb);
