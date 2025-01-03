resource "google_service_account" "invoker" {
  account_id   = "gcr-${var.short_name}-invoker"
  display_name = "Cloud Run ${var.name} Invoker"
}

resource "google_cloud_run_service_iam_binding" "default" {
  location = google_cloud_run_v2_service.default.location
  service  = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.invoker.email}"]
}

output "invoker_email" {
  value = google_service_account.invoker.email
}
