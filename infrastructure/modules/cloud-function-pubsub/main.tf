resource "random_id" "bucket_prefix" {
  byte_length = 8
}


data "archive_file" "source_code" {
  type        = "zip"
  output_path = "/tmp/gcf-source-${var.function_name}-${random_id.bucket_prefix.hex}.zip"
  source_dir  = var.source_dir
}


resource "google_project_service" "cloudbuild" {
  service            = "cloudbuild.googleapis.com"
  disable_on_destroy = false
}
resource "google_project_service" "cloudfunctions" {
  service            = "cloudfunctions.googleapis.com"
  disable_on_destroy = false
}
resource "google_project_service" "eventarc" {
  service            = "eventarc.googleapis.com"
  disable_on_destroy = false
}
resource "google_project_service" "cloudrun" {
  service            = "run.googleapis.com"
  disable_on_destroy = false
}

resource "google_storage_bucket" "source" {
  name                        = "gcf-source-${var.function_name}-${random_id.bucket_prefix.hex}" # Every bucket name must be globally unique
  location                    = var.location
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_storage_bucket_object" "default" {
  name   = "gcf-source-${data.archive_file.source_code.output_md5}.zip"
  bucket = google_storage_bucket.source.name
  source = data.archive_file.source_code.output_path
}

resource "google_service_account" "runner" {
  account_id = "gcf-sa-${var.function_name}-runner"
}

resource "google_service_account" "cloudbuild" {
  account_id = "gcf-sa-${var.function_name}-build"
}

resource "google_service_account" "invoker" {
  account_id = "gcf-sa-${var.function_name}-invoker"
}

resource "google_project_iam_member" "builder_log_writer" {
  project = google_service_account.cloudbuild.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}

resource "google_project_iam_member" "builder_artifact_registry_writer" {
  project = google_service_account.cloudbuild.project
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}

resource "google_project_iam_member" "builder_storage_object_admin" {
  project = google_service_account.cloudbuild.project
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}

resource "google_project_iam_member" "runner_pubsub_publisher" {
  project = google_service_account.runner.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.runner.email}"
}

resource "google_project_iam_member" "runner_firestore_writer" {
  project = google_service_account.runner.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.runner.email}"
}

resource "google_project_iam_member" "runner_log_writer" {
  project = google_service_account.runner.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.runner.email}"
}

resource "time_sleep" "wait_builder_permissions" {
  create_duration = "30s"

  depends_on = [
    google_project_iam_member.builder_log_writer,
    google_project_iam_member.builder_artifact_registry_writer,
    google_project_iam_member.builder_storage_object_admin,
  ]
}

resource "time_sleep" "wait_runner_permissions" {
  create_duration = "30s"

  depends_on = [
    google_project_iam_member.runner_pubsub_publisher,
    google_project_iam_member.runner_firestore_writer,
    google_project_iam_member.runner_log_writer,
  ]
}

resource "google_cloudfunctions2_function" "default" {
  name     = var.function_name
  location = var.location

  depends_on = [
    google_project_service.cloudbuild,
    google_project_service.cloudfunctions,
    google_project_service.eventarc,
    google_project_service.cloudrun,
    google_storage_bucket_object.default,
    time_sleep.wait_builder_permissions,
    time_sleep.wait_runner_permissions,
  ]

  build_config {
    runtime         = var.runtime
    entry_point     = var.entrypoint
    service_account = google_service_account.cloudbuild.id

    source {
      storage_source {
        bucket = google_storage_bucket.source.name
        object = google_storage_bucket_object.default.name
      }
    }
  }

  service_config {
    max_instance_count             = 1
    min_instance_count             = 0
    available_memory               = "256M"
    timeout_seconds                = 540
    environment_variables          = var.environment_variables
    ingress_settings               = "ALLOW_INTERNAL_ONLY"
    all_traffic_on_latest_revision = true
    service_account_email          = google_service_account.runner.email
  }

  event_trigger {
    trigger_region        = var.location
    event_type            = "google.cloud.pubsub.topic.v1.messagePublished"
    pubsub_topic          = var.pubsub_topic_id
    retry_policy          = "RETRY_POLICY_RETRY"
    service_account_email = google_service_account.invoker.email
  }
}

resource "google_cloudfunctions2_function_iam_member" "invoker" {
  project        = google_cloudfunctions2_function.default.project
  location       = google_cloudfunctions2_function.default.location
  cloud_function = google_cloudfunctions2_function.default.name
  role           = "roles/cloudfunctions.invoker"
  member         = "serviceAccount:${google_service_account.invoker.email}"
}

resource "google_cloud_run_service_iam_member" "cloud_run_invoker" {
  project  = google_cloudfunctions2_function.default.project
  location = google_cloudfunctions2_function.default.location
  service  = google_cloudfunctions2_function.default.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.invoker.email}"
}

resource "google_cloud_scheduler_job" "cron" {
  for_each = { for idx, val in var.cron_schedule : idx => val }

  name     = "gcf-${var.function_name}-cron-${each.key}"
  schedule = each.value.schedule
  project  = google_cloudfunctions2_function.default.project
  region   = google_cloudfunctions2_function.default.location

  pubsub_target {
    topic_name = var.pubsub_topic_id
    data       = base64encode(each.value.payload)
  }
}
