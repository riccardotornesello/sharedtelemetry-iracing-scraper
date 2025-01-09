resource "google_cloud_scheduler_job" "default" {
  count = var.cron != null ? 1 : 0

  name        = "${var.name}-job"
  description = "${var.name} job"
  schedule    = var.cron

  pubsub_target {
    topic_name = google_pubsub_topic.default.id
    data       = base64encode("Ok")
  }
}
