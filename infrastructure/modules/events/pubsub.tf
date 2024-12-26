
resource "google_pubsub_topic" "sessions_downloader_topic" {
  name = "sessions_downloader_topic"
}

resource "google_pubsub_topic" "season_parser_topic" {
  name = "season_parser_topic"
}


resource "google_pubsub_subscription" "sessions_downloader_subscription" {
  name  = "sessions_downloader_subscription"
  topic = google_pubsub_topic.sessions_downloader_topic.name

  depends_on = [google_cloud_run_v2_service.sessions_downloader_function]

  push_config {
    push_endpoint = google_cloud_run_v2_service.sessions_downloader_function.uri
    oidc_token {
      service_account_email = google_service_account.sessions_downloader_invoker.email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}

resource "google_pubsub_subscription" "season_parser_subscription" {
  name  = "season_parser_subscription"
  topic = google_pubsub_topic.season_parser_topic.name

  depends_on = [google_cloud_run_v2_service.season_parser_function]

  push_config {
    push_endpoint = google_cloud_run_v2_service.season_parser_function.uri
    oidc_token {
      service_account_email = google_service_account.season_parser_invoker.email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}
