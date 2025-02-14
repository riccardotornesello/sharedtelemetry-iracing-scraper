module "cars_jobs" {
  source = "../cloudrun-job"

  name           = "cars-downloader-job"
  short_name     = "cd-job"
  region         = var.region
  project        = var.project
  project_number = var.project_number

  env = {
    IRACING_EMAIL    = var.iracing_email
    IRACING_PASSWORD = var.iracing_password
    DB_USER          = google_sql_user.cars_downloader.name
    DB_PASS          = google_sql_user.cars_downloader.password
    DB_NAME          = google_sql_database.database.name
    DB_HOST          = "/cloudsql/${var.db_connection_name}"
  }

  image = "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/cars-downloader:latest"

  db_connection_name = var.db_connection_name
}
