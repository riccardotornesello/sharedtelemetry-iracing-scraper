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

resource "google_pubsub_topic" "season_parser_topic" {
  name = "season_parser_topic"
}

resource "google_pubsub_topic" "drivers_downloader_topic" {
  name = "drivers_downloader_topic"
}

resource "google_pubsub_subscription" "sessions_downloader_subscription" {
  name  = "sessions_downloader_subscription"
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

resource "google_pubsub_subscription" "season_parser_subscription" {
  name  = "season_parser_subscription"
  topic = google_pubsub_topic.season_parser_topic.name

  depends_on = [google_cloud_run_v2_service.season_parser_function]

  push_config {
    push_endpoint = google_cloud_run_v2_service.season_parser_function.uri
    oidc_token {
      service_account_email = google_service_account.invoker.email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }
}

resource "google_pubsub_subscription" "drivers_downloader_subscription" {
  name                 = "drivers_downloader_subscription"
  topic                = google_pubsub_topic.drivers_downloader_topic.name
  ack_deadline_seconds = 600

  depends_on = [google_cloud_run_v2_service.drivers_downloader_function]

  push_config {
    push_endpoint = google_cloud_run_v2_service.drivers_downloader_function.uri
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

  depends_on = [google_project_service.cloudrun_api, google_project_iam_member.sessions_downloader_runner]

  deletion_protection = false

  template {
    service_account                  = google_service_account.sessions_downloader_runner.email
    max_instance_request_concurrency = 50

    scaling {
      max_instance_count = 1
    }

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

resource "google_cloud_run_v2_service" "season_parser_function" {
  name     = "season-parser"
  location = "europe-west3"

  depends_on = [google_project_service.cloudrun_api, google_project_iam_member.season_parser_runner]

  deletion_protection = false

  template {
    service_account                  = google_service_account.season_parser_runner.email
    max_instance_request_concurrency = 50

    scaling {
      max_instance_count = 1
    }

    containers {
      image = "europe-west3-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/season-parser:latest" # TODO: variable

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
      env {
        name  = "PROJECT_ID"
        value = "sharedtelemetryapp"
      }
      env {
        name  = "TOPIC_ID"
        value = google_pubsub_topic.sessions_downloader_topic.name
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

resource "google_cloud_run_v2_service" "drivers_downloader_function" {
  name     = "drivers-downloader"
  location = "europe-west3"

  depends_on = [google_project_service.cloudrun_api, google_project_iam_member.drivers_downloader_runner]

  deletion_protection = false

  template {
    service_account                  = google_service_account.drivers_downloader_runner.email
    max_instance_request_concurrency = 50
    timeout                          = "600s"

    scaling {
      max_instance_count = 1
    }

    containers {
      image = "europe-west3-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/drivers-downloader:latest" # TODO: variable

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
  name     = "sessions-downloader"
  instance = google_sql_database_instance.sharedtelemetry.name
  password = var.db_password
}

#################################
# IAM
#################################

resource "google_service_account" "invoker" {
  account_id   = "gcr-sessions-invoker"
  display_name = "Cloud Run Sessions Downloader Invoker"
}

#################################

resource "google_service_account" "sessions_downloader_runner" {
  account_id   = "gcr-sessions-sd-runner"
  display_name = "Cloud Run Sessions Downloader Runner"
}

resource "google_project_iam_member" "sessions_downloader_runner" {
  project = "sharedtelemetryapp"
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.sessions_downloader_runner.email}"
}

resource "google_cloud_run_service_iam_binding" "sessions_downloader" {
  location = google_cloud_run_v2_service.sessions_downloader_function.location
  service  = google_cloud_run_v2_service.sessions_downloader_function.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.invoker.email}"]
}

#################################

resource "google_service_account" "season_parser_runner" {
  account_id   = "gcr-sessions-sp-runner"
  display_name = "Cloud Run Season Parser Runner"
}

resource "google_project_iam_member" "season_parser_runner" {
  for_each = toset([
    "roles/cloudsql.client",
    "roles/pubsub.publisher",
  ])

  project = "sharedtelemetryapp"
  role    = each.key
  member  = "serviceAccount:${google_service_account.season_parser_runner.email}"
}

resource "google_cloud_run_service_iam_binding" "season_parser" {
  location = google_cloud_run_v2_service.season_parser_function.location
  service  = google_cloud_run_v2_service.season_parser_function.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.invoker.email}"]
}

#################################

resource "google_service_account" "drivers_downloader_runner" {
  account_id   = "gcr-drivers-sd-runner"
  display_name = "Cloud Run Drivers Downloader Runner"
}

resource "google_project_iam_member" "drivers_downloader_runner" {
  project = "sharedtelemetryapp"
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.drivers_downloader_runner.email}"
}

resource "google_cloud_run_service_iam_binding" "drivers_downloader" {
  location = google_cloud_run_v2_service.drivers_downloader_function.location
  service  = google_cloud_run_v2_service.drivers_downloader_function.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.invoker.email}"]
}
