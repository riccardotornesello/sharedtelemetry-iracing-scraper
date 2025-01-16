import pg from 'pg';
import { env } from '$env/dynamic/private';

const { Pool } = pg;

const poolEvents = new Pool({
	database: env.EVENTS_DB_NAME,
	user: env.EVENTS_DB_USER,
	host: env.EVENTS_DB_HOST,
	password: env.EVENTS_DB_PASSWORD,
	port: Number(env.EVENTS_DB_PORT || 5432)
});

const poolDrivers = new Pool({
	database: env.DRIVERS_DB_NAME,
	user: env.DRIVERS_DB_USER,
	host: env.DRIVERS_DB_HOST,
	password: env.DRIVERS_DB_PASSWORD,
	port: Number(env.DRIVERS_DB_PORT || 5432)
});

export const connectToDB = () => Promise.all([poolEvents.connect(), poolDrivers.connect()]);
