resource "google_cloud_run_domain_mapping" "default" {
  location = var.region
  name     = var.domain

  metadata {
    namespace = "sharedtelemetryapp"
  }

  spec {
    route_name = google_cloudfunctions2_function.function.name
  }
}
