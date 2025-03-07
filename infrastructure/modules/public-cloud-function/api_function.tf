resource "google_storage_bucket" "bucket" {
  name                        = "gcf-source-${var.project}-${var.name}"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source#${filemd5(var.source_archive)}.zip"
  bucket = google_storage_bucket.bucket.name
  source = var.source_archive
}

resource "google_cloudfunctions2_function" "function" {
  name     = var.name
  location = var.region

  build_config {
    runtime     = var.runtime
    entry_point = var.entry_point
    source {
      storage_source {
        bucket = google_storage_bucket.bucket.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    min_instance_count    = 0
    max_instance_count    = 1
    available_memory      = "256M"
    timeout_seconds       = 60
    environment_variables = var.environment_variables
    service_account_email = google_service_account.runner.email
  }
}

resource "google_cloud_run_service_iam_binding" "default" {
  location = google_cloudfunctions2_function.function.location
  service  = google_cloudfunctions2_function.function.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}
