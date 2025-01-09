module "sessions_downloader_function" {
  source = "../pubsub-cloudrun"

  name       = "sessions-downloader"
  short_name = "sd"
  location   = "europe-west3"
  project    = "sharedtelemetryapp"
  image      = "europe-west3-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/sessions-downloader:latest"
  env = {
    IRACING_EMAIL : var.iracing_email,
    IRACING_PASSWORD : var.iracing_password,
    DB_USER : google_sql_user.default.name,
    DB_PASS : google_sql_user.default.password,
    DB_NAME : google_sql_database.default.name,
    DB_HOST : "/cloudsql/${var.db_connection_name}",
  }
  db_connection_name = var.db_connection_name
}

module "season_parser_function" {
  source = "../pubsub-cloudrun"

  name       = "season-parser"
  short_name = "sp"
  location   = "europe-west3"
  project    = "sharedtelemetryapp"
  image      = "europe-west3-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/season-parser:latest"
  env = {
    IRACING_EMAIL : var.iracing_email,
    IRACING_PASSWORD : var.iracing_password,
    DB_USER : google_sql_user.default.name,
    DB_PASS : google_sql_user.default.password,
    DB_NAME : google_sql_database.default.name,
    DB_HOST : "/cloudsql/${var.db_connection_name}",
    PUBSUB_PROJECT : "sharedtelemetryapp",
    PUBSUB_TOPIC : module.sessions_downloader_function.pubsub_topic.name
  }
  db_connection_name = var.db_connection_name
  pubsub_client      = true
}

module "leagues_parser_function" {
  source = "../pubsub-cloudrun"

  name       = "leagues-parser"
  short_name = "lp"
  location   = "europe-west3"
  project    = "sharedtelemetryapp"
  image      = "europe-west3-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/leagues-parser:latest"
  env = {
    IRACING_EMAIL : var.iracing_email,
    IRACING_PASSWORD : var.iracing_password,
    DB_USER : google_sql_user.default.name,
    DB_PASS : google_sql_user.default.password,
    DB_NAME : google_sql_database.default.name,
    DB_HOST : "/cloudsql/${var.db_connection_name}",
    PUBSUB_PROJECT : "sharedtelemetryapp",
    PUBSUB_TOPIC : module.season_parser_function.pubsub_topic.name
  }
  db_connection_name = var.db_connection_name
  pubsub_client      = true
}
