# Shared telemetry app

Shared telemetry app is a project with the goal of creating an ecosystem of useful tools for managing teams and events on the major simracing simulators.

This project is very young and started as a personal experiment, so initially there will be sporadic updates and a lot of confusion.

## Architecture

The architectural cornerstones for this project are:

- overengineering
- highly varied technologies
- chaos

Why? Because this project is about learning new things, so it is useful to complicate life.

The languages and frameworks used are:

- Sveltekit for the frontend
- Plain Go for the backend tasks
- Gin for the API

## DB schema

```mermaid
erDiagram
    %%{init: {
    "theme": "default",
    "themeCSS": [
        ".er.relationshipLabel { fill: black; }",
        ".er.relationshipLabelBox { fill: white; }",
        ".er.entityBox { fill: lightgray; }",
        "[id^=entity-BRAND] .er.entityBox { fill: lightgreen;} ",
        "[id^=entity-CAR_CLASS] .er.entityBox { fill: lightgreen;} ",
        "[id^=entity-CAR_IN_CLASS] .er.entityBox { fill: lightgreen;} ",
        "[id^=entity-CAR] .er.entityBox { fill: lightgreen;} ",
        "[id^=entity-DRIVER_STATS] .er.entityBox { fill: pink;} ",
        "[id^=entity-DRIVER] .er.entityBox { fill: pink;} "
    ]
    }}%%

    %%%%%%%%%%%%%%%%%%%%%%%

    BRAND {
        string name PK
        string icon
    }
    CAR_CLASS {
        int id PK
        string name
        string short_name
    }
    CAR_IN_CLASS {
        int id PK
        int car_id FK
        int car_class_id FK
    }
    CAR {
        int id PK
        string name
        string name_abbreviated
        string brand
        string logo
        string small_image
        string sponsor_logo
    }

    CAR_IN_CLASS }|--|| CAR: ""
    CAR_IN_CLASS }|--|| CAR_CLASS: ""
    CAR }|..|| BRAND: ""

    %%%%%%%%%%%%%%%%%%%%%%%%%%%

    DRIVER_STATS {
        int cust_id PK
        int car_category PK
        string license
        int i_rating
    }
    DRIVER {
        int cust_id PK
        string name
        string location
    }
    DRIVER_STATS }|..|| DRIVER: ""

    %%%%%%%%%%%%%%%%%%%%%%%%%%%

    COMPETITION_CLASS {
        int id PK
        int competition_id FK
        string name
        string color
        int index
    }
    COMPETITION_CREW {
        int id PK
        int team_id FK
        int class_id FK
        string name
        int i_racing_car_id
    }
    COMPETITION_DRIVER {
        int id PK
        int crew_id FK
        int i_racing_cust_id
        string first_name
        string last_name
    }
    COMPETITION_TEAM {
        int id PK
        int competition_id FK
        string name
        string picture
    }
    COMPETITION {
        int id PK
        int league_id FK
        int season_id FK
        string name
        string slug
        int crew_drivers_count
    }
    EVENT_GROUP {
        int id PK
        int competition_id FK
        string name
        int i_racing_track_id
        string[] dates
    }
    LAP {
        int id PK
        int subsession_id FK
        int simsession_number FK
        int cust_id FK
        string[] lap_events
        bool incident
        int lap_time
        int lap_number
    }
    LEAGUE_SEASON {
        int league_id PK, FK
        int season_id PK
    }
    LEAGUE {
        int league_id PK
    }
    SESSION_SIMSESSION_PARTICIPANT {
        int subsession_id PK, FK
        int simsession_number PK, FK
        int cust_id PK
        int car_id
    }
    SESSION_SIMSESSION {
        int subsession_id PK, FK
        int simsession_number PK
        int simsession_type
        string simsession_name
    }
    SESSION {
        int subsession_id PK
        int league_id
        int season_id
        time launch_at
        int track_id
    }

    COMPETITION_CLASS }|--|| COMPETITION: ""
    COMPETITION_CREW }|--|| COMPETITION_TEAM: ""
    COMPETITION_CREW }o--o| COMPETITION_CLASS: ""
    COMPETITION_DRIVER }|--|| COMPETITION_CREW: ""
    COMPETITION_TEAM }|--|| COMPETITION: ""
    COMPETITION }|--|| LEAGUE_SEASON: ""
    EVENT_GROUP }|--|| COMPETITION: ""
    LAP }|--|| SESSION_SIMSESSION_PARTICIPANT: ""
    LEAGUE_SEASON }|--|| LEAGUE: ""
    SESSION_SIMSESSION_PARTICIPANT }|--|| SESSION_SIMSESSION: ""
    SESSION_SIMSESSION }|--|| SESSION: ""
    LEAGUE_SEASON }|..|| SESSION: ""
```
