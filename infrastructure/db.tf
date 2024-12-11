resource "google_sql_database_instance" "sessions_db" {
  name             = "sessions-db"
  database_version = "POSTGRES_17"

  depends_on = [google_service_networking_connection.private_vpc_connection]

  settings {
    edition           = "ENTERPRISE"
    tier              = "db-f1-micro"
    availability_type = "ZONAL"
    disk_autoresize   = true

    ip_configuration {
      ipv4_enabled    = true # TODO: Check if it's possible to allow only connections from specific IP addresses
      private_network = google_compute_network.main.self_link
    }
  }
}
