resource "google_cloud_run_v2_service" "qualify_results_frontend" {
  name     = "qualify-results-frontend"
  location = "europe-west3"

  depends_on = [google_project_iam_member.qualify_results_frontend_runner]

  deletion_protection = false

  template {
    service_account = google_service_account.qualify_results_frontend_runner.email

    scaling {
      max_instance_count = 1
    }

    containers {
      image = "europe-west3-docker.pkg.dev/sharedtelemetryapp/qualify-results/qualify-results-front:latest" # TODO: variable

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name  = "DB_USER"
        value = var.db_user
      }
      env {
        name  = "DB_PASSWORD"
        value = var.db_password
      }
      env {
        name  = "DB_NAME"
        value = var.db_name
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

# resource "google_cloud_run_domain_mapping" "qualify_results_frontend_domain" {
#   location = "europe-west3"
#   name     = "verified-domain.com" # TODO

#   spec {
#     route_name = google_cloud_run_v2_service.qualify_results_frontend.name
#   }
# }
