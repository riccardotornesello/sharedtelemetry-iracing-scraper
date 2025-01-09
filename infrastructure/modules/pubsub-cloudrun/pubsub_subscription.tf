resource "google_pubsub_subscription" "default" {
  name  = "${var.name}-subscription"
  topic = google_pubsub_topic.default.name

  ack_deadline_seconds = var.ack_deadline_seconds

  push_config {
    push_endpoint = google_cloud_run_v2_service.default.uri
    oidc_token {
      service_account_email = google_service_account.invoker.email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}
