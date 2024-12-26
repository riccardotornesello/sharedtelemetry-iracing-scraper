resource "google_service_account" "sessions_downloader_invoker" {
  account_id   = "gcr-sd-invoker"
  display_name = "Cloud Run Sessions Downloader Invoker"
}

resource "google_service_account" "sessions_downloader_runner" {
  account_id   = "gcr-sd-runner"
  display_name = "Cloud Run Sessions Downloader Runner"
}

# Allow sessions downloader to connecto to cloud sql
resource "google_project_iam_member" "sessions_downloader_runner" {
  project = "sharedtelemetryapp"
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.sessions_downloader_runner.email}"
}

# Allow pub sub to invoke sessions downloader
resource "google_cloud_run_service_iam_binding" "sessions_downloader" {
  location = google_cloud_run_v2_service.sessions_downloader_function.location
  service  = google_cloud_run_v2_service.sessions_downloader_function.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.sessions_downloader_invoker.email}"]
}



resource "google_service_account" "season_parser_invoker" {
  account_id   = "gcr-sp-invoker"
  display_name = "Cloud Run Sessions Downloader Invoker"
}

resource "google_service_account" "season_parser_runner" {
  account_id   = "gcr-sp-runner"
  display_name = "Cloud Run Season Parser Runner"
}

# Allow the season parser to connect to the db and to push new tasks to the queue
resource "google_project_iam_member" "season_parser_runner" {
  for_each = toset([
    "roles/cloudsql.client",
    "roles/pubsub.publisher",
  ])

  project = "sharedtelemetryapp"
  role    = each.key
  member  = "serviceAccount:${google_service_account.season_parser_runner.email}"
}

# Allow pubsub to invoke the season parser
resource "google_cloud_run_service_iam_binding" "season_parser" {
  location = google_cloud_run_v2_service.season_parser_function.location
  service  = google_cloud_run_v2_service.season_parser_function.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.season_parser_invoker.email}"]
}
