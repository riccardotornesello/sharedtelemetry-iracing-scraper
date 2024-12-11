SELECT cust_id,
    MIN(avg)
FROM (
        SELECT laps.cust_id,
            laps.event_session_id,
            SUM(laps.lap_time) / 3 AS avg
        FROM laps
            LEFT JOIN event_sessions ON event_sessions.id = laps.event_session_id
        WHERE event_sessions.simsession_name = 'QUALIFY'
            AND laps.incident = false
            AND laps.lap_time > 0
            AND laps.lap_number > 0
            AND NOT laps.lap_events LIKE '%off track%'
            AND NOT laps.lap_events LIKE '%pitted%'
            AND NOT laps.lap_events LIKE '%invalid%'
        GROUP BY laps.cust_id,
            laps.event_session_id
        HAVING COUNT(*) = 3
    )
GROUP BY cust_id;