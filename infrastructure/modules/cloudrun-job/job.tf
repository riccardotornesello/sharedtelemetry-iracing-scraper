resource "google_cloud_run_v2_job" "default" {
  provider = google-beta

  name     = var.name
  location = var.region

  deletion_protection   = false
  start_execution_token = var.run_after_deploy ? formatdate("YYYYMMDDhhmmss", timestamp()) : null

  template {
    annotations = {
      "deploy-time" = timestamp()
    }

    template {
      service_account = google_service_account.runner.email

      volumes {
        name = "cloudsql"
        cloud_sql_instance {
          instances = [var.db_connection_name]
        }
      }

      containers {
        image = var.image
        args  = var.args

        dynamic "env" {
          for_each = var.env
          content {
            name  = env.key
            value = env.value
          }
        }

        volume_mounts {
          name       = "cloudsql"
          mount_path = "/cloudsql"
        }
      }
    }
  }
}

resource "google_service_account" "runner" {
  account_id = "gcr-${var.short_name}-runner"
}

resource "google_project_iam_member" "runner" {
  project = var.project
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.runner.email}"
}
