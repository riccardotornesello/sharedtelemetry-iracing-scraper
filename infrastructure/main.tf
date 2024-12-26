provider "google" {
  project         = "sharedtelemetryapp"
  region          = "europe-west3"
  zone            = "europe-west3-a"
  request_timeout = "60s"
}

module "drivers" {
  source = "./modules/drivers"

  iracing_email      = var.iracing_email
  iracing_password   = var.iracing_password
  db_password        = var.db_password
  db_instance_name   = google_sql_database_instance.sharedtelemetry.name
  db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
}

module "events" {
  source = "./modules/events"

  iracing_email      = var.iracing_email
  iracing_password   = var.iracing_password
  db_password        = var.db_password
  db_instance_name   = google_sql_database_instance.sharedtelemetry.name
  db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
}
