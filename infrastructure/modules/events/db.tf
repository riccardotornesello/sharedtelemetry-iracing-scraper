resource "google_sql_database" "database" {
  name     = "iracing_events"
  instance = var.db_instance_name
}

resource "google_sql_user" "events_parser" {
  name     = "events_parser"
  instance = var.db_instance_name
  password = var.db_password
}
