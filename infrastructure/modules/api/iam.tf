resource "google_service_account" "api_runner" {
  account_id   = "gcr-api-runner"
  display_name = "Cloud Run Qualify Results Frontend Runner"
}

resource "google_project_iam_member" "api_runner" {
  project = "sharedtelemetryapp"
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.api_runner.email}"
}

resource "google_cloud_run_service_iam_binding" "api" {
  location = google_cloud_run_v2_service.api.location
  service  = google_cloud_run_v2_service.api.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}
