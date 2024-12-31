resource "google_sql_database" "database" {
  name     = "iracing_events"
  instance = var.db_instance_name
}

resource "google_sql_user" "events_parser" {
  name     = "events_parser"
  instance = var.db_instance_name
  password = var.db_password
}

output "db_user" {
  value = google_sql_user.events_parser.name
}

output "db_name" {
  value = google_sql_database.database.name
}
