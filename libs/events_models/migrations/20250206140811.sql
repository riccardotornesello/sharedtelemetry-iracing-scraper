-- Create "leagues" table
CREATE TABLE "public"."leagues" (
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "league_id" bigserial NOT NULL,
  PRIMARY KEY ("league_id")
);
-- Create index "idx_leagues_deleted_at" to table: "leagues"
CREATE INDEX "idx_leagues_deleted_at" ON "public"."leagues" ("deleted_at");
-- Create "league_seasons" table
CREATE TABLE "public"."league_seasons" (
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "league_id" bigint NOT NULL,
  "season_id" bigint NOT NULL,
  PRIMARY KEY ("league_id", "season_id"),
  CONSTRAINT "fk_league_seasons_league" FOREIGN KEY ("league_id") REFERENCES "public"."leagues" ("league_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_league_seasons_deleted_at" to table: "league_seasons"
CREATE INDEX "idx_league_seasons_deleted_at" ON "public"."league_seasons" ("deleted_at");
-- Create "competitions" table
CREATE TABLE "public"."competitions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "league_id" bigint NOT NULL,
  "season_id" bigint NOT NULL,
  "name" text NOT NULL,
  "slug" text NOT NULL,
  "crew_drivers_count" bigint NOT NULL DEFAULT 1,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_competitions_slug" UNIQUE ("slug"),
  CONSTRAINT "fk_competitions_league_season" FOREIGN KEY ("league_id", "season_id") REFERENCES "public"."league_seasons" ("league_id", "season_id") ON UPDATE SET NULL ON DELETE SET NULL
);
-- Create index "idx_competitions_deleted_at" to table: "competitions"
CREATE INDEX "idx_competitions_deleted_at" ON "public"."competitions" ("deleted_at");
-- Create "competition_teams" table
CREATE TABLE "public"."competition_teams" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "competition_id" bigint NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_competition_teams_competition" FOREIGN KEY ("competition_id") REFERENCES "public"."competitions" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_competition_teams_deleted_at" to table: "competition_teams"
CREATE INDEX "idx_competition_teams_deleted_at" ON "public"."competition_teams" ("deleted_at");
-- Create "competition_crews" table
CREATE TABLE "public"."competition_crews" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "team_id" bigint NOT NULL,
  "name" text NOT NULL,
  "i_racing_car_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_competition_crews_team" FOREIGN KEY ("team_id") REFERENCES "public"."competition_teams" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_competition_crews_deleted_at" to table: "competition_crews"
CREATE INDEX "idx_competition_crews_deleted_at" ON "public"."competition_crews" ("deleted_at");
-- Create "competition_drivers" table
CREATE TABLE "public"."competition_drivers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "crew_id" bigint NOT NULL,
  "name" text NOT NULL,
  "i_racing_cust_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_competition_drivers_crew" FOREIGN KEY ("crew_id") REFERENCES "public"."competition_crews" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_competition_drivers_deleted_at" to table: "competition_drivers"
CREATE INDEX "idx_competition_drivers_deleted_at" ON "public"."competition_drivers" ("deleted_at");
-- Create "event_groups" table
CREATE TABLE "public"."event_groups" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "competition_id" bigint NOT NULL,
  "name" text NOT NULL,
  "i_racing_track_id" bigint NOT NULL,
  "dates" text[] NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_event_groups_competition" FOREIGN KEY ("competition_id") REFERENCES "public"."competitions" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_event_groups_deleted_at" to table: "event_groups"
CREATE INDEX "idx_event_groups_deleted_at" ON "public"."event_groups" ("deleted_at");
-- Create "sessions" table
CREATE TABLE "public"."sessions" (
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "subsession_id" bigserial NOT NULL,
  "league_id" bigint NULL,
  "season_id" bigint NULL,
  "launch_at" timestamptz NULL,
  "track_id" bigint NULL,
  PRIMARY KEY ("subsession_id")
);
-- Create index "idx_sessions_deleted_at" to table: "sessions"
CREATE INDEX "idx_sessions_deleted_at" ON "public"."sessions" ("deleted_at");
-- Create index "idx_sessions_launch_at" to table: "sessions"
CREATE INDEX "idx_sessions_launch_at" ON "public"."sessions" ("launch_at");
-- Create "session_simsessions" table
CREATE TABLE "public"."session_simsessions" (
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "subsession_id" bigint NOT NULL,
  "simsession_number" bigint NOT NULL,
  "simsession_type" bigint NULL,
  "simsession_name" text NULL,
  PRIMARY KEY ("subsession_id", "simsession_number"),
  CONSTRAINT "fk_session_simsessions_session" FOREIGN KEY ("subsession_id") REFERENCES "public"."sessions" ("subsession_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_session_simsessions_deleted_at" to table: "session_simsessions"
CREATE INDEX "idx_session_simsessions_deleted_at" ON "public"."session_simsessions" ("deleted_at");
-- Create "session_simsession_participants" table
CREATE TABLE "public"."session_simsession_participants" (
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "subsession_id" bigint NOT NULL,
  "simsession_number" bigint NOT NULL,
  "cust_id" bigint NOT NULL,
  "car_id" bigint NULL,
  PRIMARY KEY ("subsession_id", "simsession_number", "cust_id"),
  CONSTRAINT "fk_session_simsession_participants_session_simsession" FOREIGN KEY ("subsession_id", "simsession_number") REFERENCES "public"."session_simsessions" ("subsession_id", "simsession_number") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_session_simsession_participants_deleted_at" to table: "session_simsession_participants"
CREATE INDEX "idx_session_simsession_participants_deleted_at" ON "public"."session_simsession_participants" ("deleted_at");
-- Create "laps" table
CREATE TABLE "public"."laps" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "subsession_id" bigint NOT NULL,
  "simsession_number" bigint NOT NULL,
  "cust_id" bigint NOT NULL,
  "lap_events" text[] NULL,
  "incident" boolean NULL,
  "lap_time" bigint NULL,
  "lap_number" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_laps_session_simsession_participant" FOREIGN KEY ("subsession_id", "simsession_number", "cust_id") REFERENCES "public"."session_simsession_participants" ("subsession_id", "simsession_number", "cust_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_laps_deleted_at" to table: "laps"
CREATE INDEX "idx_laps_deleted_at" ON "public"."laps" ("deleted_at");
