SELECT cust_id,
    MIN(avg) * interval '1 sec'
FROM (
        SELECT laps.cust_id,
            laps.subsession_id,
            SUM(laps.lap_time) / 3 / 10000 AS avg
        FROM laps
            LEFT JOIN event_sessions ON event_sessions.subsession_id = laps.subsession_id
            AND laps.simsession_number = event_sessions.simsession_number
            LEFT JOIN events ON events.subsession_id = event_sessions.subsession_id
        WHERE event_sessions.simsession_name = 'QUALIFY'
            AND laps.incident = false
            AND laps.lap_time > 0
            AND laps.lap_number > 0
            AND NOT laps.lap_events && array ['off track', 'pitted', 'invalid']
            AND (
                events.launch_at BETWEEN '2024-09-11 00:00:00' AND '2024-09-11 23:59:59'
                OR events.launch_at BETWEEN '2024-09-14 00:00:00' AND '2024-09-14 23:59:59'
            )
        GROUP BY laps.cust_id,
            laps.subsession_id
        HAVING COUNT(*) = 3
    )
GROUP BY cust_id
ORDER BY MIN(avg) ASC;