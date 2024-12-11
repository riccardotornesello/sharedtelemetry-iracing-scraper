provider "google" {
  project         = "sharedtelemetryapp"
  region          = "europe-west3"
  zone            = "europe-west3-a"
  request_timeout = "60s"
}
