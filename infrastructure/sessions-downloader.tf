resource "google_storage_bucket" "source_bucket" {
  name                        = "sharedtelemetryapp-gcf-source"
  location                    = "EU"
  uniform_bucket_level_access = true
}

data "archive_file" "default" {
  type        = "zip"
  output_path = "/tmp/function-source.zip"
  source_dir  = "../package/"
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.source_bucket.name
  source = "/tmp/function-source.zip"
}

resource "google_project_iam_member" "build_account_editor" {
  # TODO: only permission to get the bucket

  project = "sharedtelemetryapp"
  role    = "roles/editor"
  member  = "serviceAccount:${google_service_account.build_account.email}"
}

resource "google_service_account" "build_account" {
  account_id   = "gcf-build-sa"
  display_name = "GCF Build Service Account"
}

resource "google_cloudfunctions2_function" "function" {
  name        = "function-v2"
  location    = "us-central1"
  description = "a new function"

  depends_on = [google_project_service.api_run, google_project_iam_member.build_account_editor]

  build_config {
    runtime         = "go122"
    entry_point     = "SessionsDownloader"
    service_account = google_service_account.build_account.id
    source {
      storage_source {
        bucket = google_storage_bucket.source_bucket.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    max_instance_count = 1
    available_memory   = "256M"
    timeout_seconds    = 60
  }

  event_trigger {
    event_type   = "google.cloud.pubsub.topic.v1.messagePublished"
    pubsub_topic = google_pubsub_topic.default.id
    retry_policy = "RETRY_POLICY_RETRY"
  }
}

resource "google_pubsub_topic" "default" {
  name = "pubsub_topic"
}

resource "google_pubsub_subscription" "default" {
  name  = "pubsub_subscription"
  topic = google_pubsub_topic.default.name
}

resource "google_cloud_scheduler_job" "default" {
  name        = "test-job"
  description = "test job"
  schedule    = "* * * * *"

  pubsub_target {
    topic_name = google_pubsub_topic.default.id
    data       = base64encode("Hello world!")
  }
}
