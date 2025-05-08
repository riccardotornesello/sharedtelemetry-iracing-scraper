terraform { 
  cloud { 
    
    organization = "sharedtelemetry" 

    workspaces { 
      name = "sharedtelemetry-iracing-scraper-prod" 
    } 
  } 
}

provider "google" {
  project         = "sharedtelemetryapp"
  region          = var.region
  request_timeout = "60s"
}

resource "google_pubsub_topic" "cars_parse_trigger" {
  name = "iracing_cars_parse_trigger"
}

resource "google_pubsub_topic" "drivers_parse_trigger" {
  name = "iracing_drivers_parse_trigger"
}

resource "google_pubsub_topic" "leagues_parse_trigger" {
  name = "iracing_leagues_parse_trigger"
}

resource "google_pubsub_topic" "season_parse_trigger" {
  name = "iracing_season_parse_trigger"
}

resource "google_pubsub_topic" "session_parse_trigger" {
  name = "iracing_session_parse_trigger"
}

locals {
  tasks = [
    {
      name         = "cars"
      source_dir   = "${path.module}/../apps/cars"
      pubsub_topic = google_pubsub_topic.cars_parse_trigger.id
      environment_variables = {
        FIRESTORE_PROJECT_ID = "sharedtelemetryapp"
        IRACING_EMAIL        = var.iracing_email
        IRACING_PASSWORD     = var.iracing_password
      }
    },
    {
      name         = "drivers"
      source_dir   = "${path.module}/../apps/drivers"
      pubsub_topic = google_pubsub_topic.drivers_parse_trigger.id
      environment_variables = {
        FIRESTORE_PROJECT_ID = "sharedtelemetryapp"
        IRACING_EMAIL        = var.iracing_email
        IRACING_PASSWORD     = var.iracing_password
      }
    },
    {
      name         = "leagues"
      source_dir   = "${path.module}/../apps/leagues"
      pubsub_topic = google_pubsub_topic.leagues_parse_trigger.id
      environment_variables = {
        FIRESTORE_PROJECT_ID = "sharedtelemetryapp"
        PUBSUB_PROJECT_ID    = "sharedtelemetryapp"
        PUBSUB_TOPIC_ID      = google_pubsub_topic.season_parse_trigger.name
      }
    },
    {
      name         = "season"
      source_dir   = "${path.module}/../apps/season"
      pubsub_topic = google_pubsub_topic.season_parse_trigger.id
      environment_variables = {
        FIRESTORE_PROJECT_ID = "sharedtelemetryapp"
        IRACING_EMAIL        = var.iracing_email
        IRACING_PASSWORD     = var.iracing_password
        PUBSUB_PROJECT_ID    = "sharedtelemetryapp"
        PUBSUB_TOPIC_ID      = google_pubsub_topic.session_parse_trigger.name
      }
    },
    {
      name         = "sessions"
      source_dir   = "${path.module}/../apps/sessions"
      pubsub_topic = google_pubsub_topic.session_parse_trigger.id
      environment_variables = {
        FIRESTORE_PROJECT_ID = "sharedtelemetryapp"
        IRACING_EMAIL        = var.iracing_email
        IRACING_PASSWORD     = var.iracing_password
      }
    }
  ]
}

module "tasks" {
  for_each = { for task in local.tasks : task.name => task }

  source = "./modules/cloud-function-pubsub"

  location              = var.region
  function_name         = "ir-${each.key}"
  source_dir            = each.value.source_dir
  pubsub_topic_id       = each.value.pubsub_topic
  environment_variables = each.value.environment_variables
}
