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
}

module "events" {
  source = "./modules/events"

  iracing_email      = var.iracing_email
  iracing_password   = var.iracing_password
  db_password        = var.db_password
  db_instance_name   = google_sql_database_instance.sharedtelemetry.name
  db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
  region             = var.region
}

module "qualify_results" {
  source = "./modules/qualify-results"
  domain = var.domain

  events_db_user     = module.events.db_user.name
  events_db_password = var.db_password
  events_db_name     = module.events.db.name

  drivers_db_user     = module.drivers.db_user.name
  drivers_db_password = var.db_password
  drivers_db_name     = module.drivers.db.name

  db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
  region             = var.region
}
