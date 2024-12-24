import pg from 'pg';
import { env } from '$env/dynamic/private';

const { Pool } = pg;

const pool = new Pool({
	database: env.DB_NAME,
	user: env.DB_USER,
	host: env.DB_HOST,
	password: env.DB_PASSWORD,
	port: Number(env.DB_PORT || 5432)
});

export const connectToDB = async () => await pool.connect();
