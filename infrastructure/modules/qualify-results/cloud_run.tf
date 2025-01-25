resource "google_cloud_run_v2_service" "qualify_results_frontend" {
  name     = "qualify-results-frontend"
  location = var.region

  depends_on = [google_project_iam_member.qualify_results_frontend_runner]

  deletion_protection = false

  template {
    service_account = google_service_account.qualify_results_frontend_runner.email

    scaling {
      max_instance_count = 1
    }

    containers {
      image = "europe-west1-docker.pkg.dev/sharedtelemetryapp/qualify-results/qualify-results-front:latest" # TODO: variable

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name  = "EVENTS_DB_USER"
        value = var.events_db_user
      }
      env {
        name  = "EVENTS_DB_PASSWORD"
        value = var.events_db_password
      }
      env {
        name  = "EVENTS_DB_NAME"
        value = var.events_db_name
      }
      env {
        name  = "EVENTS_DB_HOST"
        value = "/cloudsql/${var.db_connection_name}"
      }

      env {
        name  = "DRIVERS_DB_USER"
        value = var.drivers_db_user
      }
      env {
        name  = "DRIVERS_DB_PASSWORD"
        value = var.drivers_db_password
      }
      env {
        name  = "DRIVERS_DB_NAME"
        value = var.drivers_db_name
      }
      env {
        name  = "DRIVERS_DB_HOST"
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

resource "google_cloud_run_domain_mapping" "default" {
  location = var.region
  name     = var.domain

  metadata {
    namespace = "sharedtelemetryapp"
  }

  spec {
    route_name = google_cloud_run_v2_service.qualify_results_frontend.name
  }
}
