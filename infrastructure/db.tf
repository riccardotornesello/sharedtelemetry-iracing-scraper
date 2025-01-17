resource "google_sql_database_instance" "sharedtelemetry" {
  # TODO: certificate

  name             = "sessions-db"
  database_version = "POSTGRES_17"

  settings {
    edition           = "ENTERPRISE"
    tier              = "db-f1-micro"
    availability_type = "ZONAL"
    disk_autoresize   = true

    ip_configuration {
      ipv4_enabled = true

      dynamic "authorized_networks" {
        for_each = var.db_whitelist
        content {
          value = authorized_networks.value
        }
      }
    }
  }
}
