module "drivers_jobs" {
  source = "../cron-cloudrun"

  for_each = tomap({
    sports_car = {
      id   = "sports-car"
      cron = "0 3 * * *"
    }
    oval = {
      id   = "oval"
      cron = "5 3 * * *"
    }
    formula_car = {
      id   = "formula-car"
      cron = "10 3 * * *"
    }
    road = {
      id   = "road"
      cron = "15 3 * * *"
    }
    dirt_oval = {
      id   = "dirt-oval"
      cron = "20 3 * * *"
    }
    dirt_road = {
      id   = "dirt-road"
      cron = "25 3 * * *"
    }
  })

  name           = "drivers-downloader-job-${each.value.id}"
  short_name     = "dd-job-${each.value.id}"
  region         = var.region
  project        = var.project
  project_number = var.project_number
  schedule       = each.value.cron

  env = {
    IRACING_EMAIL    = var.iracing_email
    IRACING_PASSWORD = var.iracing_password
    DB_USER          = google_sql_user.drivers_downloader.name
    DB_PASS          = google_sql_user.drivers_downloader.password
    DB_NAME          = google_sql_database.database.name
    DB_HOST          = "/cloudsql/${var.db_connection_name}"
    CAR_CLASS        = each.key
  }

  image = "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/drivers-downloader:latest"

  db_connection_name = var.db_connection_name
}
