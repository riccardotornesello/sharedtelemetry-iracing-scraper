output "job" {
  value = google_cloud_run_v2_job.default
}

output "runner" {
  value = google_service_account.runner
}
