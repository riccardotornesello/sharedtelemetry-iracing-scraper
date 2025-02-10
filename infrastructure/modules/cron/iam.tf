resource "google_service_account" "invoker" {
  account_id = "gcr-${var.short_name}-invoker"
}

resource "google_cloud_run_v2_job_iam_binding" "invoker" {
  project  = var.project
  location = var.region
  name     = var.job_name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.invoker.email}"]
}
