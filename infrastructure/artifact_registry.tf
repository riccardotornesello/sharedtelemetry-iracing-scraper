resource "google_artifact_registry_repository" "sessions_downloader_repository" {
  repository_id = "sessions-downloader"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository" "qualify_results_repository" {
  repository_id = "qualify-results"
  format        = "DOCKER"
}
