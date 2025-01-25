resource "google_cloud_run_v2_service" "default" {
  name     = var.name
  location = var.location

  depends_on = [google_project_iam_member.runner]

  deletion_protection = false

  template {
    service_account                  = google_service_account.runner.email
    max_instance_request_concurrency = 1
    timeout                          = var.timeout

    scaling {
      max_instance_count = var.max_instance_count
    }

    containers {
      image = var.image

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      dynamic "env" {
        for_each = var.env
        content {
          name  = env.key
          value = env.value
        }
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

output "cloud_run" {
  value = google_cloud_run_v2_service.default
}
