resource "google_sql_database" "default" {
  name     = "iracing_events"
  instance = var.db_instance_name
}

resource "google_sql_user" "default" {
  name     = "iracing_events"
  instance = var.db_instance_name
  password = var.db_password
}
