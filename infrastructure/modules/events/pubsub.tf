resource "google_pubsub_topic" "leagues_parser_topic" {
  name = "leagues_parser_topic"
}

resource "google_pubsub_topic" "season_parser_topic" {
  name = "season_parser_topic"
}

resource "google_pubsub_topic" "sessions_downloader_topic" {
  name = "sessions_downloader_topic"
}
