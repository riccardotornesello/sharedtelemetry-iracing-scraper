module "events_migration" {
  source = "../cloudrun-job"

  name           = "migration-events"
  short_name     = "me"
  region         = var.region
  project        = var.project
  project_number = var.project_number

  env = {
    DB_USER = google_sql_user.default.name
    DB_PASS = google_sql_user.default.password
    DB_NAME = google_sql_database.default.name
    DB_HOST = "/cloudsql/${var.db_connection_name}"
  }

  run_after_deploy = true
  image            = "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/events-models:latest"
  args             = ["migrate", "apply", "--url", "postgres://${google_sql_user.default.name}:${google_sql_user.default.password}@/${google_sql_database.default.name}?host=/cloudsql/${var.db_connection_name}"]

  db_connection_name = var.db_connection_name
}
