resource "google_service_account" "runner" {
  account_id   = "gcf-${var.short_name}-sa"
  display_name = "Test Service Account"
}

resource "google_project_iam_member" "runner" {
  for_each = toset(var.roles)

  project = google_service_account.runner.project
  role    = each.value
  member  = "serviceAccount:${google_service_account.runner.email}"
}
