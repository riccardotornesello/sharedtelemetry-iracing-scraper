variable "classes" {
  type = map(object({
    id   = string
    cron = string
  }))

  default = {
    # sports_car = {
    #   id   = "sports-car"
    #   cron = "0 2 * * *"
    # }
    # oval = {
    #   id   = "oval"
    #   cron = "10 2 * * *"
    # }
    # formula_car = {
    #   id   = "formula-car"
    #   cron = "20 2 * * *"
    # }
    # road = {
    #   id   = "road"
    #   cron = "30 2 * * *"
    # }
    # dirt_oval = {
    #   id   = "dirt-oval"
    #   cron = "40 2 * * *"
    # }
    # dirt_road = {
    #   id   = "dirt-road"
    #   cron = "50 2 * * *"
    # }
  }
}

module "drivers_jobs" {
  source = "../cloudrun-job"

  for_each = var.classes

  name           = "drivers-downloader-job-${each.value.id}"
  short_name     = "dd-job-${each.value.id}"
  region         = var.region
  project        = var.project
  project_number = var.project_number

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

module "drivers_jobs_cron" {
  source = "../cron"

  for_each = var.classes

  name           = "drivers-downloader-job-${each.value.id}"
  short_name     = "dd-job-${each.value.id}"
  region         = var.region
  project        = var.project
  project_number = var.project_number
  schedule       = each.value.cron
  job_name       = module.drivers_jobs[each.key].job.name
}
