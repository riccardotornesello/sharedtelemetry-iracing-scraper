resource "google_cloud_run_v2_service" "drivers_downloader_function" {
  name     = "drivers-downloader"
  location = "europe-west3"

  depends_on = [google_project_iam_member.drivers_downloader_runner]

  deletion_protection = false

  template {
    service_account                  = google_service_account.drivers_downloader_runner.email
    max_instance_request_concurrency = 50
    timeout                          = "600s"

    scaling {
      max_instance_count = 1
    }

    containers {
      image = "europe-west3-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/drivers-downloader:latest" # TODO: variable

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name  = "IRACING_EMAIL"
        value = var.iracing_email
      }
      env {
        name  = "IRACING_PASSWORD"
        value = var.iracing_password
      }
      env {
        name  = "DB_USER"
        value = google_sql_user.drivers_downloader.name
      }
      env {
        name  = "DB_PASS"
        value = google_sql_user.drivers_downloader.password
      }
      env {
        name  = "DB_NAME"
        value = google_sql_database.database.name
      }
      env {
        name  = "DB_HOST"
        value = "/cloudsql/${var.db_connection_name}"
      }
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [var.db_connection_name]
      }
    }
  }
}
