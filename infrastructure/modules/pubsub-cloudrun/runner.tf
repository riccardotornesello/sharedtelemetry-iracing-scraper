resource "google_service_account" "runner" {
  account_id   = "gcr-${var.short_name}-runner"
  display_name = "Cloud Run ${var.name} Runner"
}

resource "google_project_iam_member" "runner" {
  # If var.db_connection_name is not null add the cloudsql.client role
  # Also, if var.pubsub_client is true add the pubsub.publisher role
  for_each = {
    for role in [
      var.db_connection_name != null ? "roles/cloudsql.client" : null,
      var.pubsub_client ? "roles/pubsub.publisher" : null
    ] : role => role if role != null
  }

  project = var.project
  role    = each.key
  member  = "serviceAccount:${google_service_account.runner.email}"
}
