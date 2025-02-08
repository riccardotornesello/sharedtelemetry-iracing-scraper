resource "google_cloud_run_v2_job" "default" {
  depends_on = [google_project_iam_member.runner]

  name     = var.name
  location = var.region

  deletion_protection = false

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

resource "google_service_account" "invoker" {
  account_id = "gcr-${var.short_name}-invoker"
}

resource "google_cloud_run_v2_job_iam_binding" "invoker" {
  project  = var.project
  location = var.region
  name     = google_cloud_run_v2_job.default.name
  role     = "roles/viewer"
  members  = ["serviceAccount:${google_service_account.invoker.email}"]
}

resource "google_cloud_scheduler_job" "job" {
  name     = "${var.name}-scheduler"
  schedule = var.schedule
  region   = var.region
  project  = var.project

  http_target {
    http_method = "POST"
    uri         = "https://${google_cloud_run_v2_job.default.location}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project_number}/jobs/${google_cloud_run_v2_job.default.name}:run"

    oauth_token {
      service_account_email = google_service_account.invoker.email
    }
  }

  depends_on = [google_cloud_run_v2_job.default, google_cloud_run_v2_job_iam_binding.invoker]
}
