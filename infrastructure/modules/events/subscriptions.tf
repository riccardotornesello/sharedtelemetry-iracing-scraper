resource "google_pubsub_subscription" "leagues_parser_subscription" {
  name  = "leagues_parser_subscription"
  topic = google_pubsub_topic.leagues_parser_topic.name

  push_config {
    push_endpoint = module.leagues_parser_function.uri
    oidc_token {
      service_account_email = module.leagues_parser_function.invoker_email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}

resource "google_pubsub_subscription" "season_parser_subscription" {
  name  = "season_parser_subscription"
  topic = google_pubsub_topic.season_parser_topic.name

  push_config {
    push_endpoint = module.season_parser_function.uri
    oidc_token {
      service_account_email = module.season_parser_function.invoker_email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}

resource "google_pubsub_subscription" "sessions_downloader_subscription" {
  name  = "sessions_downloader_subscription"
  topic = google_pubsub_topic.sessions_downloader_topic.name

  push_config {
    push_endpoint = module.sessions_downloader_function.uri
    oidc_token {
      service_account_email = module.sessions_downloader_function.invoker_email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}
