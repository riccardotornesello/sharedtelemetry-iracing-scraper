resource "google_pubsub_topic" "drivers_downloader_topic" {
  name = "drivers_downloader_topic"
}

resource "google_pubsub_subscription" "drivers_downloader_subscription" {
  name                 = "drivers_downloader_subscription"
  topic                = google_pubsub_topic.drivers_downloader_topic.name
  ack_deadline_seconds = 600

  depends_on = [google_cloud_run_v2_service.drivers_downloader_function]

  push_config {
    push_endpoint = google_cloud_run_v2_service.drivers_downloader_function.uri
    oidc_token {
      service_account_email = google_service_account.drivers_downloader_invoker.email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}
