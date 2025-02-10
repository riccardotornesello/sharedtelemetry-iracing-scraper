provider "google-beta" {
  project         = "sharedtelemetryapp"
  region          = var.region
  zone            = "${var.region}-a"
  request_timeout = "60s"
}

provider "google" {
  project         = "sharedtelemetryapp"
  region          = var.region
  zone            = "${var.region}-a"
  request_timeout = "60s"
}

module "drivers" {
  source = "./modules/drivers"

  iracing_email      = var.iracing_email
  iracing_password   = var.iracing_password
  db_password        = var.db_password
  db_instance_name   = google_sql_database_instance.sharedtelemetry.name
  db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
  region             = var.region
  project            = var.project
  project_number     = var.project_number
}

module "events" {
  source = "./modules/events"

  iracing_email      = var.iracing_email
  iracing_password   = var.iracing_password
  db_password        = var.db_password
  db_instance_name   = google_sql_database_instance.sharedtelemetry.name
  db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
  region             = var.region
  project            = var.project
  project_number     = var.project_number
}

module "qualify_results" {
  source = "./modules/qualify-results"
  domain = var.domain

  region = var.region
}

module "api" {
  source = "./modules/api"
  domain = "api.${var.domain}"

  db_user            = module.events.db_user.name
  db_password        = var.db_password
  db_name            = module.events.db.name
  db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name

  region = var.region
}
