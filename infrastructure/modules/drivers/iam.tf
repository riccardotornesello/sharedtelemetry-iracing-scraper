resource "google_service_account" "drivers_downloader_invoker" {
  account_id   = "gcr-drivers-downloader-invoker"
  display_name = "Cloud Run Sessions Downloader Invoker"
}

resource "google_service_account" "drivers_downloader_runner" {
  account_id   = "gcr-drivers-sd-runner"
  display_name = "Cloud Run Drivers Downloader Runner"
}

# Allow the drivers downloader user to connect to cloudsql
resource "google_project_iam_member" "drivers_downloader_runner" {
  project = "sharedtelemetryapp"
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.drivers_downloader_runner.email}"
}

# Allow the invoker to run the function
resource "google_cloud_run_service_iam_binding" "drivers_downloader" {
  location = google_cloud_run_v2_service.drivers_downloader_function.location
  service  = google_cloud_run_v2_service.drivers_downloader_function.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.drivers_downloader_invoker.email}"]
}
