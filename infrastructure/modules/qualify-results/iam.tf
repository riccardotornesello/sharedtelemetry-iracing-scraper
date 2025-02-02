resource "google_service_account" "qualify_results_frontend_runner" {
  account_id   = "gcr-qrf-runner"
  display_name = "Cloud Run Qualify Results Frontend Runner"
}

resource "google_cloud_run_service_iam_binding" "qualify_results_frontend" {
  location = google_cloud_run_v2_service.qualify_results_frontend.location
  service  = google_cloud_run_v2_service.qualify_results_frontend.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}
