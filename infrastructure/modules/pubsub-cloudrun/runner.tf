resource "google_service_account" "runner" {
  account_id   = "gcr-${var.name}-runner"
  display_name = "Cloud Run ${var.name} Runner"
}

resource "google_project_iam_member" "runner" {
  for_each = var.db_connection_name != null ? ["roles/cloudsql.client"] : []

  project = var.project
  role    = each.key
  member  = "serviceAccount:${google_service_account.runner.email}"
}
