#################################
# Artifact Registry
#################################

resource "google_artifact_registry_repository" "sessions_downloader_repository" {
  repository_id = "sessions-downloader"
  format        = "DOCKER"
}

#################################
# Pub/Sub
#################################

resource "google_pubsub_topic" "sessions_downloader_topic" {
  name = "sessions_downloader_topic"
}

resource "google_pubsub_subscription" "subscription" {
  name  = "pubsub_subscription"
  topic = google_pubsub_topic.sessions_downloader_topic.name

  depends_on = [google_cloud_run_v2_service.sessions_downloader_function]

  push_config {
    push_endpoint = google_cloud_run_v2_service.sessions_downloader_function.uri
    oidc_token {
      service_account_email = google_service_account.invoker.email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}

#################################
# Google Cloud Run
#################################

resource "google_cloud_run_v2_service" "sessions_downloader_function" {
  name     = "sessions-downloader"
  location = "europe-west3"

  depends_on = [google_project_service.cloudrun_api, google_project_iam_member.runner]

  deletion_protection = false

  template {
    service_account = google_service_account.runner.email

    containers {
      image = "europe-west3-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/sessions-downloader:latest" # TODO: variable

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name  = "IRACING_EMAIL"
        value = var.iracing_email
      }
      env {
        name  = "IRACING_PASSWORD"
        value = var.iracing_password
      }
      env {
        name  = "DB_USER"
        value = google_sql_user.sessions_downloader_user.name
      }
      env {
        name  = "DB_PASS"
        value = google_sql_user.sessions_downloader_user.password
      }
      env {
        name  = "DB_NAME"
        value = google_sql_database.database.name
      }
      env {
        name  = "DB_HOST"
        value = "/cloudsql/${google_sql_database_instance.sharedtelemetry.connection_name}"
      }
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.sharedtelemetry.connection_name]
      }
    }
  }
}

#################################
# Google Cloud SQL
#################################

resource "google_sql_user" "sessions_downloader_user" {
  name       = "sessions-downloader"
  instance   = google_sql_database_instance.sharedtelemetry.name
  password   = var.db_password
}

#################################
# IAM
#################################

resource "google_service_account" "invoker" {
  account_id   = "gcr-sessions-invoker"
  display_name = "Cloud Run Sessions Downloader Invoker"
}

resource "google_service_account" "runner" {
  account_id   = "gcr-sessions-runner"
  display_name = "Cloud Run Sessions Downloader Runner"
}

resource "google_project_iam_member" "runner" {
  project = "sharedtelemetryapp"
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.runner.email}"
}

resource "google_cloud_run_service_iam_binding" "invoker" {
  location = google_cloud_run_v2_service.sessions_downloader_function.location
  service  = google_cloud_run_v2_service.sessions_downloader_function.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.invoker.email}"]
}
