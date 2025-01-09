resource "google_pubsub_topic" "default" {
  name = "${var.name}-topic"
}

output "pubsub_topic" {
  value = google_pubsub_topic.default
}
