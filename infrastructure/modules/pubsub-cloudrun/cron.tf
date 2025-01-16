resource "google_cloud_scheduler_job" "default" {
  for_each = { for job in var.cron : base64encode(job.data) => job }

  name        = "${var.name}-job-${index(var.cron, each.value)}"
  description = "${var.name} job"
  schedule    = each.value.schedule

  pubsub_target {
    topic_name = google_pubsub_topic.default.id
    data       = base64encode(each.value.data)
  }
}
