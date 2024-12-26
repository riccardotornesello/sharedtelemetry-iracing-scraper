resource "google_cloud_run_v2_service" "sessions_downloader_function" {
  name     = "sessions-downloader"
  location = "europe-west3"

  depends_on = [google_project_iam_member.sessions_downloader_runner]

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
        value = google_sql_user.events_parser.name
      }
      env {
        name  = "DB_PASS"
        value = google_sql_user.events_parser.password
      }
      env {
        name  = "DB_NAME"
        value = google_sql_database.database.name
      }
      env {
        name  = "DB_HOST"
        value = "/cloudsql/${var.db_connection_name}"
      }
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [var.db_connection_name]
      }
    }
  }
}

resource "google_cloud_run_v2_service" "season_parser_function" {
  name     = "season-parser"
  location = "europe-west3"

  depends_on = [google_project_iam_member.season_parser_runner]

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
        value = google_sql_user.events_parser.name
      }
      env {
        name  = "DB_PASS"
        value = google_sql_user.events_parser.password
      }
      env {
        name  = "DB_NAME"
        value = google_sql_database.database.name
      }
      env {
        name  = "DB_HOST"
        value = "/cloudsql/${var.db_connection_name}"
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
        instances = [var.db_connection_name]
      }
    }
  }
}
