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

module "test" {
  source = "./modules/public-cloud-function"

  region         = var.region
  project        = "sharedtelemetryapp"
  source_archive = "../apps/results/api/dist/api.zip"
  name           = "api"
  short_name     = "api"
  entry_point    = "apiNEST"
  runtime        = "nodejs22"
  environment_variables = {
    "ENVIRONMENT" = "production"
  }
  domain = "api.results.sharedtelemetry.com"
  roles  = ["roles/firestore.serviceAgent", "roles/datastore.viewer"]
}

# module "drivers" {
#   source = "./modules/drivers"

#   iracing_email      = var.iracing_email
#   iracing_password   = var.iracing_password
#   db_password        = var.db_password
#   db_instance_name   = google_sql_database_instance.sharedtelemetry.name
#   db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
#   region             = var.region
#   project            = var.project
#   project_number     = var.project_number
# }

# module "cars" {
#   source = "./modules/cars"

#   iracing_email      = var.iracing_email
#   iracing_password   = var.iracing_password
#   db_password        = var.db_password
#   db_instance_name   = google_sql_database_instance.sharedtelemetry.name
#   db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
#   region             = var.region
#   project            = var.project
#   project_number     = var.project_number
# }

# module "events" {
#   source = "./modules/events"

#   iracing_email      = var.iracing_email
#   iracing_password   = var.iracing_password
#   db_password        = var.db_password
#   db_instance_name   = google_sql_database_instance.sharedtelemetry.name
#   db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name
#   region             = var.region
#   project            = var.project
#   project_number     = var.project_number
# }

module "qualify_results" {
  source = "./modules/qualify-results"
  domain = var.domain

  region = var.region
}

# module "api" {
#   source = "./modules/api"
#   domain = "api.${var.domain}"

#   db_connection_name = google_sql_database_instance.sharedtelemetry.connection_name

#   events_db_user     = module.events.db_user.name
#   events_db_password = var.db_password
#   events_db_name     = module.events.db.name

#   cars_db_user     = module.cars.db_user.name
#   cars_db_password = var.db_password
#   cars_db_name     = module.cars.db.name

#   region = var.region
# }
