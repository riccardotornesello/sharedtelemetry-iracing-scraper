resource "google_cloud_run_v2_service" "qualify_results_frontend" {
  name     = "qualify-results-frontend"
  location = var.region

  deletion_protection = false

  template {
    service_account = google_service_account.qualify_results_frontend_runner.email

    scaling {
      max_instance_count = 1
    }

    containers {
      image = "europe-west1-docker.pkg.dev/sharedtelemetryapp/qualify-results/qualify-results-front:latest" # TODO: variable

      env {
        name  = "API_BASE_URL"
        value = "api.${var.domain}"
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
