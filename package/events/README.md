# Events package

## Commands

### Leagues parser

No payload

### Season parser

Payload:

- leagueId
- seasonId

### Sessions downloader

Payload:

- subsessionId
- launchAt

## Database

```mermaid
erDiagram
    LEAGUE["LEAGUE (iRacing)"] {
        int LeagueID PK
        string Name
    }
    LEAGUE_SEASON["LEAGUE_SEASON (iRacing's group of sessions)"] {
        int LeagueID PK,FK
        int SeasonID PK
    }

    COMPETITION["COMPETITION (platform's season with specific rules)"]
    COMPETITION_DRIVER
    COMPETITION_TEAM
    COMPETITION_CREW
    EVENT_GROUP

    SESSION["SESSION (iRacing's subsession)"] {
        int SubsessionID PK
        int LeagueID
        int SeasonID
        datetime LaunchAt
        int TrackID
    }
    SESSION_SIMSESSION["SESSION_SIMSESSION (event's part, like qualifying, race...)"] {
        int SubsessionID PK,FK
        int SimsessionNumber PK
        int SimsessionType
        string SimsessionName
    }
    SESSION_SIMSESSION_PARTICIPANT["SESSION_SIMSESSION_PARTICIPANT (user who participated in a specific session)"] {
        int SubsessionID PK,FK
        int SimsessionNumber PK,FK
        int CustID PK
    }
    LAP {
        int ID PK
        int SubsessionID FK
        int SimsessionNumber FK
        int CustID FK
        string[] LapCOMPETITIONs
        bool Incident
        int LapTime
        int LapNumber
    }

    LEAGUE ||--|{ LEAGUE_SEASON: ""

    LEAGUE_SEASON ||--|{ COMPETITION: ""
    COMPETITION ||--|{ COMPETITION_TEAM: ""
    COMPETITION_TEAM ||--|{ COMPETITION_CREW: ""
    COMPETITION_CREW ||--|{ COMPETITION_DRIVER: ""

    COMPETITION ||--|{ EVENT_GROUP: ""

    LEAGUE_SEASON ||--|{ SESSION: "-------"
    SESSION ||--|{ SESSION_SIMSESSION: ""
    SESSION_SIMSESSION ||--|{ SESSION_SIMSESSION_PARTICIPANT: ""
    SESSION_SIMSESSION_PARTICIPANT ||--|{ LAP: ""
```

## Architecture

```mermaid
sequenceDiagram
    box yellow
    participant CRON as 🕒 CRON - Every 5 minutes
    end
    box cyan
    participant PubSub1 as 📢 Topic PubSub - Trigger League Parser
    end
    box grey
    participant LeagueParser as ☁️ Cloud Run - leagues_parser
    end
    box cyan
    participant PubSub2 as 📢 Topic PubSub - Trigger Season Parser
    end
    box grey
    participant SeasonParser as ☁️ Cloud Run - season_parser
    end
    box cyan
    participant PubSub3 as 📢 Topic PubSub - Trigger Sessions Downloader
    end
    box grey
    participant SessionsDownloader as ☁️ Cloud Run - sessions_downloader
    end
    box purple
    participant DB as Database 📃
    end
    box purple
    participant IRacing as iRacing 📃
    end

    CRON->>PubSub1: Start operation

    PubSub1->>LeagueParser: Trigger 🛠
    LeagueParser->>DB: Get leagues
    DB->>LeagueParser: 
    loop Each league
        LeagueParser->>PubSub2: League id and season id
    end

    PubSub2->>SeasonParser: Trigger 🛠
    SeasonParser->>IRacing: Get sessions
    IRacing->>SeasonParser: 
    loop Each session
        SeasonParser->>PubSub3: Subsession id and date
    end

    PubSub3->>SessionsDownloader: Trigger 🛠
    SessionsDownloader->>IRacing: Get results
    IRacing->>SessionsDownloader: 
    SessionsDownloader->>DB: Store session 📄
```
