resource "google_cloud_scheduler_job" "job" {
  name     = "${var.name}-scheduler"
  schedule = var.schedule
  region   = var.region
  project  = var.project

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project_number}/jobs/${var.job_name}:run"

    oauth_token {
      service_account_email = google_service_account.invoker.email
    }
  }
}
