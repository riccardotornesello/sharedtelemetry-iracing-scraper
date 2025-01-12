module "drivers_downloader_function" {
  source = "../pubsub-cloudrun"

  name       = "drivers-downloader"
  short_name = "dd"
  location   = "europe-west3"
  project    = "sharedtelemetryapp"
  image      = "europe-west3-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/drivers-downloader:latest"
  env = {
    IRACING_EMAIL : var.iracing_email,
    IRACING_PASSWORD : var.iracing_password,
    DB_USER : google_sql_user.drivers_downloader.name,
    DB_PASS : google_sql_user.drivers_downloader.password,
    DB_NAME : google_sql_database.database.name,
    DB_HOST : "/cloudsql/${var.db_connection_name}",
  }
  db_connection_name = var.db_connection_name
  cron = [{
    schedule = "0 * * * *"
    data     = "{\"carClass\":\"sports_car\"}"
  }]
}
