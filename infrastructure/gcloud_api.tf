resource "google_project_service" "cloudrun_api" {
  project = "sharedtelemetryapp"
  service = "run.googleapis.com"
}
